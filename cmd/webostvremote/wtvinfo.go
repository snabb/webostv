package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/snabb/webostv"
	"sync"
	"time"
)

type tvInfo struct {
	*tview.TextView
	sync.Mutex
	systemInfo        webostv.SystemInfo
	foregroundAppInfo webostv.ForegroundAppInfo
	appNames          map[string]string
	tvCurrentChannel  webostv.TvCurrentChannel
}

func newTvInfo() *tvInfo {
	w := tview.NewTextView()
	w.SetBorder(true)
	w.SetScrollable(true)
	w.SetTitle("TV Information")
	w.SetWrap(false)
	w.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w.ScrollToBeginning()
		return x + 1, y + 1, width - 2, height - 2
	})

	i := &tvInfo{TextView: w}

	return i
}

func (i *tvInfo) updateFromTv() (err error) {
	systemInfo, err := tv.SystemGetSystemInfo()
	if err != nil {
		return errors.Wrap(err, "error getting system info from TV")
	}

	appNames := make(map[string]string)

	apps, err := tv.ApplicationManagerListApps()
	if err != nil {
		return errors.Wrap(err, "error getting app list from TV")
	}
	for _, app := range apps {
		appNames[app.Id] = app.Title
	}

	i.Lock()
	i.systemInfo = systemInfo
	i.appNames = appNames
	i.update()
	i.Unlock()
	return nil
}

func (i *tvInfo) update() {
	systemInfo := i.systemInfo
	foregroundAppInfo := i.foregroundAppInfo
	tvCurrentChannel := i.tvCurrentChannel

	i.Clear()

	fmt.Fprintf(i, "Address  %s", tv.Address)
	fmt.Fprintf(i, "\nModel    %s", systemInfo.ModelName)
	fmt.Fprintf(i, "\nReceiver %s", systemInfo.ReceiverType)

	var features string
	for k, v := range systemInfo.Features {
		if !v {
			continue
		}
		if features == "" {
			features += k
		} else {
			features += ", " + k
		}
	}
	fmt.Fprintf(i, "\nFeatures %s", features)

	appName := foregroundAppInfo.AppId
	if i.appNames != nil {
		str := i.appNames[appName]
		if str != "" {
			appName = str
		}
	}
	fmt.Fprintf(i, "\nApp      %s", appName)

	if tvCurrentChannel.ChannelNumber != "" {
		fmt.Fprintf(i, "\nChannel  %s • %s", tvCurrentChannel.ChannelNumber, tvCurrentChannel.ChannelName)
		fmt.Fprintf(i, "\n         %s", tvCurrentChannel.ChannelTypeName)
	} else {
		fmt.Fprint(i, "\n\n")
	}

	i.ScrollToBeginning()
}

func (i *tvInfo) monitorTvCurrentInfo(quit chan struct{}) (err error) {
	var channelQuitCh chan struct{}
	errorCh := make(chan error, 1)

	err = tv.ApplicationManagerMonitorForegroundAppInfo(func(info webostv.ForegroundAppInfo) error {
		var startChannelMonitor, stopChannelMonitor bool
		i.Lock()
		if info.IsLiveTv() && !i.foregroundAppInfo.IsLiveTv() {
			startChannelMonitor = true
		} else if !info.IsLiveTv() && i.foregroundAppInfo.IsLiveTv() {
			stopChannelMonitor = true
		}
		i.foregroundAppInfo = info
		i.update()
		i.Unlock()

		if startChannelMonitor {
			app.logger.Debug("starting channel monitor")
			channelQuitCh = make(chan struct{})
			go func() {
				err := i.monitorTvCurrentChannel(channelQuitCh)
				errorCh <- err
				if err != nil {
					close(quit)
				}
			}()
		} else if stopChannelMonitor {
			app.logger.Debug("stopping channel monitor")
			close(channelQuitCh)
			err := <-errorCh
			channelQuitCh = nil
			i.Lock()
			i.tvCurrentChannel = webostv.TvCurrentChannel{}
			i.update()
			i.Unlock()
			if err != nil {
				return err
			}
		}

		return nil
	}, quit)

	if channelQuitCh != nil {
		close(channelQuitCh)
	}
	if err != nil {
		return err
	}
	select {
	case err = <-errorCh:
	default:
	}
	return err
}

var errRetry = errors.New("retry")

func (i *tvInfo) monitorTvCurrentChannel(quit chan struct{}) (err error) {
	for {
		err = tv.TvMonitorCurrentChannel(func(cur webostv.TvCurrentChannel) error {
			app.logger.Debug("got current channel message", "cur", cur)
			if cur.ChannelNumber == "0" && cur.IsSkipped {
				// this happens if we have sent the subscription message too quickly
				return errRetry
			}
			i.Lock()
			i.tvCurrentChannel = cur
			i.update()
			i.Unlock()
			return nil
		}, quit)

		if err == errRetry {
			select {
			case <-time.After(time.Second):
				app.logger.Debug("retrying current channel subscription")
				continue
			case <-quit:
				return nil
			}
		}
		return err
	}
}
