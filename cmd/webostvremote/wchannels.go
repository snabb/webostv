package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/snabb/webostv"
	"sort"
	"sync"
	"time"
)

type channels struct {
	*tview.Table
	channels                                   []webostv.TvChannel
	channelsMutex                              sync.Mutex
	updateInfo                                 func(str string)
	cancelPreviousGetChannelCurrentProgramInfo CancelPrevious
}

func newChannels() *channels {
	w := tview.NewTable()
	w.SetBorder(true)
	w.SetTitle("Channels")
	w.SetSelectable(false, false)

	c := &channels{Table: w}
	w.SetSelectedFunc(c.selected)
	w.SetSelectionChangedFunc(c.selectionChanged)

	return c
}

func (c *channels) cancelTasks() {
	c.cancelPreviousGetChannelCurrentProgramInfo.Cancel()
}

func (c *channels) selected(row, column int) {
	var sel webostv.TvChannel
	set := false
	c.channelsMutex.Lock()
	if c.channels != nil && row < len(c.channels) {
		sel = c.channels[row]
		set = true
	}
	c.channelsMutex.Unlock()
	if set {
		go tv.TvOpenChannelId(sel.ChannelId)
		// XXX check err
	}
}

func (c *channels) selectionChanged(row, column int) {
	if c.updateInfo == nil {
		return
	}
	var sel webostv.TvChannel
	set := false
	c.channelsMutex.Lock()
	if c.channels != nil && row < len(c.channels) {
		sel = c.channels[row]
		set = true
	}
	c.channelsMutex.Unlock()
	if !set {
		c.updateInfo("")
		return
	}
	updBasic := fmt.Sprintf("Channel number: %s, name: %s\ntype: %s, hdtv: %v\nsignal id: %s, frequency: %d", sel.ChannelNumber, sel.ChannelName, sel.ChannelType, sel.HDTV, sel.SignalChannelId, sel.Frequency)

	c.updateInfo(updBasic)

	go func() {
		cancelCh := c.cancelPreviousGetChannelCurrentProgramInfo.NewCancel()
		// throttle
		select {
		case <-cancelCh:
			app.logger.Debug("canceled")
			return
		case <-time.After(time.Millisecond * 500):
		}

		info, err := tv.TvGetChannelCurrentProgramInfo(sel.ChannelId)
		if err != nil {
			app.logger.Error("TvGetChannelCurrentProgramInfo error", "channelId", sel.ChannelId, "err", err)
			return
		}
		c.cancelPreviousGetChannelCurrentProgramInfo.Lock()
		defer c.cancelPreviousGetChannelCurrentProgramInfo.Unlock()
		select {
		case <-cancelCh:
			app.logger.Debug("canceled")
			return
		default:
		}
		c.updateInfo(updBasic + fmt.Sprintf("\ncurrent program: %s\nstart: %s, end: %s, duration %d\ndescription: %s", info.ProgramName, info.LocalStartTime, info.LocalEndTime, info.Duration, info.Description))
	}()
}

func (c *channels) updateFromTv() (err error) {
	tvChannels, err := tv.TvGetChannelList()
	if err != nil {
		return errors.Wrap(err, "error updating channels from TV")
	}

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
	c.update(tvChannels)
	return nil
}

func (c *channels) update(tvChannels []webostv.TvChannel) {
	c.channelsMutex.Lock()
	c.channels = tvChannels
	c.channelsMutex.Unlock()

	c.Clear()
	for row, ch := range tvChannels {
		var info string
		if ch.HDTV {
			info = "HDTV"
		} else if ch.TV {
			info = "TV"
		} else if ch.Radio {
			info = "Radio"
		}

		c.SetCell(row, 0, tview.NewTableCell(ch.ChannelNumber).SetAlign(tview.AlignRight))
		c.SetCell(row, 1, tview.NewTableCell(ch.ChannelName))
		c.SetCell(row, 2, tview.NewTableCell(info))
	}
	c.ScrollToBeginning()
}
