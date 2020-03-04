package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/inconshreveable/log15"
	"github.com/ogier/pflag"
	"github.com/rivo/tview"
	"github.com/snabb/webostv"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const DefaultAddress = "LGsmartTV.lan"

type myError struct {
	where string
	err   error
}

type myTv struct {
	*webostv.Tv
	errorCh chan myError
}

var tv myTv

type myApp struct {
	*tview.Application

	wTvInfo   *tvInfo
	wVolume   *volume
	wHelp     tview.Primitive
	wSelInfo  *selInfo
	wChannels *channels
	wInputs   *inputs
	wApps     *apps
	nextFocus map[tview.Primitive]tview.Primitive

	logger log15.Logger
}

var app = myApp{
	Application: tview.NewApplication(),
	logger:      log15.Root(),
}

func setSelectable(widget tview.Primitive, yesno bool) {
	if widget, ok := widget.(interface {
		SetSelectable(bool, bool) *tview.Table
	}); ok {
		widget.SetSelectable(yesno, false)
	}
}

func cancelTasks(widget tview.Primitive) {
	if widget, ok := widget.(interface {
		cancelTasks()
	}); ok {
		widget.cancelTasks()
	}
}

func (app *myApp) changeFocus(currentFocus, newFocus tview.Primitive) {
	if currentFocus != nil {
		setSelectable(currentFocus, false)
		cancelTasks(currentFocus)
	}
	setSelectable(newFocus, true)
	app.SetFocus(newFocus)
	if newFocus, ok := newFocus.(interface {
		selectionChanged(int, int)
		GetSelection() (int, int)
	}); ok {
		newFocus.selectionChanged(newFocus.GetSelection())
	}
}

func (app *myApp) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	currentFocus := app.GetFocus()
	key := event.Key()
	switch key {
	case tcell.KeyTAB:
		if nf, ok := app.nextFocus[currentFocus]; ok {
			app.changeFocus(currentFocus, nf)
		} else {
			app.changeFocus(currentFocus, app.wChannels)
		}
		return nil
	case tcell.KeyBacktab:
		for k, v := range app.nextFocus {
			if v == currentFocus {
				app.changeFocus(currentFocus, k)
				return nil
			} 
		}
		app.changeFocus(currentFocus, app.wVolume)
		return nil
	case tcell.KeyExit, tcell.KeyESC:
		app.Stop()
		return nil
	case tcell.KeyCtrlX:
		err := tv.SystemTurnOff()
		if err != nil {
			app.logger.Error("error turning off", "err", err)
		}
		app.Stop()
		return nil
	}
	if _, ok := currentFocus.(*tview.InputField); ok {
		return event
	}

	switch key {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'v', 'V':
			app.changeFocus(currentFocus, app.wVolume)
			return nil
		case 'c', 'C':
			app.changeFocus(currentFocus, app.wChannels)
			return nil
		case 'i', 'I':
			app.changeFocus(currentFocus, app.wInputs)
			return nil
		case 'a', 'A':
			app.changeFocus(currentFocus, app.wApps)
			return nil
		case 'q', 'Q':
			app.Stop()
			return nil
		}

	}
	return event
}

func (app *myApp) initWidgets() {
	app.wTvInfo = newTvInfo()
	app.wVolume = newVolume()
	app.wHelp = newHelp()

	app.wSelInfo = newSelInfo()

	app.wChannels = newChannels()
	app.wChannels.updateInfo = app.wSelInfo.update

	app.wInputs = newInputs()
	app.wInputs.updateInfo = app.wSelInfo.update

	app.wApps = newApps()
	app.wApps.updateInfo = app.wSelInfo.update
}

func (app *myApp) initLayout() {
	layoutLeft := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.wTvInfo, 0, 2, false).
		AddItem(app.wVolume, 3, 0, false).
		AddItem(app.wHelp, 0, 2, false)

	layoutRight := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(app.wChannels, 0, 4, false).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(app.wInputs, 0, 1, false).
				AddItem(app.wApps, 0, 2, false), 0, 3, false), 0, 3, false).
		AddItem(app.wSelInfo, 0, 1, false)

	layout := tview.NewFlex().
		AddItem(layoutLeft, 0, 1, false).
		AddItem(layoutRight, 0, 2, false)

	app.SetRoot(layout, true)

	app.nextFocus = map[tview.Primitive]tview.Primitive{
		app.wChannels: app.wInputs,
		app.wInputs:   app.wApps,
		app.wApps:     app.wVolume,
		app.wVolume:   app.wChannels,
	}
}

func initTv(address string) {
	store := openMyStore()
	clientKey := store.Get(address)

	var err error
	tv.Tv, err = webostv.DefaultDialer.Dial(address)
	if err != nil {
		fmt.Fprintln(os.Stderr, "TV connection error:", err)
		os.Exit(1)
	}

	tv.errorCh = make(chan myError, 8)
	go func() {
		err := tv.MessageHandler()
		tv.errorCh <- myError{"tv.MessageHandler()", err}
		app.Stop()
	}()

	newKey, err := tv.Register(clientKey)
	if err != nil {
		tv.Close()
		fmt.Fprintln(os.Stderr, "TV registration error:", err)
		os.Exit(1)
	}

	if newKey != clientKey {
		store.Set(address, newKey)
	}
	store.Close()
}

func main() {
	var err error

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "[OPTION]... [ADDRESS]\n")
		fmt.Fprintln(os.Stderr, "ADDRESS is the name or IP address of the LG WebOS TV (default: \""+DefaultAddress+"\").\n")
		fmt.Fprintln(os.Stderr, "The following OPTIONS are available:")
		pflag.PrintDefaults()
	}

	debugLog := pflag.StringP("debug", "d", "", "debug log file name")
	pflag.Parse()

	if *debugLog != "" {
		logHandler, err := log15.FileHandler(*debugLog, log15.LogfmtFormat())
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening debug log:", err)
			os.Exit(1)
		}
		app.logger.SetHandler(logHandler)
	} else {
		app.logger.SetHandler(log15.DiscardHandler())
	}

	var address string
	switch pflag.NArg() {
	case 0:
		address = DefaultAddress
	case 1:
		address = pflag.Arg(0)
	default:
		pflag.Usage()
		os.Exit(1)
	}

	app.logger.Debug("starting")

	rand.Seed(time.Now().UnixNano())

	initTv(address)

	if *debugLog != "" {
		tv.SetDebug(func(str string) {
			app.logger.Debug(str)
		})
	}

	app.initWidgets()
	app.initLayout()
	app.SetInputCapture(app.inputCapture)

	var wg sync.WaitGroup
	quit := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tv.AudioMonitorStatus(func(as webostv.AudioStatus) error {
			app.wVolume.update(as.Volume)
			return nil
		}, quit)
		tv.errorCh <- myError{"AudioMonitorStatus", err}
		app.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.wTvInfo.monitorTvCurrentInfo(quit)
		tv.errorCh <- myError{"monitorTvCurrentInfo", err}
		app.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.wTvInfo.updateFromTv()
		app.Draw()
		app.wChannels.updateFromTv()
		app.Draw()
		app.wInputs.updateFromTv()
		app.Draw()
		app.wApps.updateFromTv()
		app.Draw()
		// XXX check errors ?
	}()

	err = app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		app.logger.Error("app.Run() returned error", "err", err)
	}

	close(quit)
	wg.Wait()

	var err2 error
errorChReadLoop:
	for {
		select {
		case myErr := <-tv.errorCh:
			if myErr.err != nil {
				fmt.Fprintln(os.Stderr, "error in", myErr.where+":", myErr.err)
				app.logger.Error("error", "goroutine", myErr.where, "err", myErr.err)
				err2 = myErr.err
			}
		default:
			break errorChReadLoop
		}
	}
	err3 := tv.Close()
	if err3 != nil {
		fmt.Fprintln(os.Stderr, "error:", err3)
		app.logger.Error("tv.Close() returned error", "err", err3)
	}
	app.logger.Debug("exiting")

	if err != nil || err2 != nil || err3 != nil {
		os.Exit(1)
	}
}

func openMyStore() (store *Store) {
	var name string
	if home := os.Getenv("HOME"); home != "" {
		name = filepath.Join(home, ".webostv.json")
	} else {
		name = ".webostv.json"
	}
	var err error
	store, err = OpenStore(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return store
}
