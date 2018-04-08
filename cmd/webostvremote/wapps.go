package main

import (
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/snabb/webostv"
	"sync"
)

type apps struct {
	*tview.Table
	apps       []webostv.LaunchPoint
	appsMutex  sync.Mutex
	updateInfo func(str string)
}

func newApps() *apps {
	w := tview.NewTable()
	w.SetBorder(true)
	w.SetTitle("Apps")
	w.SetSelectable(false, false)

	a := &apps{Table: w}
	w.SetSelectedFunc(a.selected)
	w.SetSelectionChangedFunc(a.selectionChanged)

	return a
}

func (a *apps) appNames() (appNames map[string]string) {
	appNames = make(map[string]string)

	a.appsMutex.Lock()
	for _, app := range a.apps {
		appNames[app.Id] = app.Title
	}
	a.appsMutex.Unlock()
	return appNames
}

func (a *apps) selected(row, column int) {
	var sel webostv.LaunchPoint
	set := false
	a.appsMutex.Lock()
	if a.apps != nil && row < len(a.apps) {
		sel = a.apps[row]
		set = true
	}
	a.appsMutex.Unlock()
	if !set {
		return
	}
	payload := make(webostv.Payload)
	for k, v := range sel.Params {
		payload[k] = v
	}

	go tv.ApplicationManagerLaunch(sel.Id, payload)
	// XXX check err
}

func (a *apps) selectionChanged(row, column int) {
	if a.updateInfo == nil {
		return
	}
	var sel webostv.LaunchPoint
	set := false
	a.appsMutex.Lock()
	if a.apps != nil && row < len(a.apps) {
		sel = a.apps[row]
		set = true
	}
	a.appsMutex.Unlock()
	if !set {
		a.updateInfo("")
		return
	}
	a.updateInfo("App title: " + sel.Title + "\n" +
		"vendor: " + sel.Vendor + "\n" +
		"version: " + sel.Version + "\n" +
		"id: " + sel.Id)
}

func (a *apps) updateFromTv() (err error) {
	tvApps, _, err := tv.ApplicationManagerListLaunchPoints()
	if err != nil {
		return errors.Wrap(err, "error updating apps from TV")
	}

	/*
		sort.Slice(tvChannels, func(i, j int) bool {
			li := len(tvChannels[i].ChannelNumber)
			lj := len(tvChannels[j].ChannelNumber)

			if li < lj {
				return true
			}
			if li > lj {
				return false
			}
			return tvChannels[i].ChannelNumber < tvChannels[j].ChannelNumber
		})
	*/
	a.update(tvApps)
	return nil
}

func (a *apps) update(tvApps []webostv.LaunchPoint) {
	a.appsMutex.Lock()
	a.apps = tvApps
	a.appsMutex.Unlock()

	a.Clear()
	for row, app := range tvApps {
		a.SetCell(row, 0, tview.NewTableCell(app.Title))
		a.SetCell(row, 1, tview.NewTableCell(app.Vendor))
	}
	a.ScrollToBeginning()
}
