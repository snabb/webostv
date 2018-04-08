package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/snabb/webostv"
	"sync"
)

type inputs struct {
	*tview.Table
	inputs      []webostv.TvExternalInput
	inputsMutex sync.Mutex
	updateInfo  func(str string)
}

func newInputs() *inputs {
	w := tview.NewTable()
	w.SetBorder(true)
	w.SetTitle("Inputs")
	w.SetSelectable(false, false)

	i := &inputs{Table: w}
	w.SetSelectedFunc(i.selected)
	w.SetSelectionChangedFunc(i.selectionChanged)

	return i
}

func (i *inputs) selected(row, column int) {
	var sel webostv.TvExternalInput
	set := false
	i.inputsMutex.Lock()
	if i.inputs != nil && row < len(i.inputs) {
		sel = i.inputs[row]
		set = true
	}
	i.inputsMutex.Unlock()
	if set {
		go tv.TvSwitchInput(sel.Id)
		// XXX check err
	}
}

func (i *inputs) selectionChanged(row, column int) {
	if i.updateInfo == nil {
		return
	}
	var sel webostv.TvExternalInput
	set := false
	i.inputsMutex.Lock()
	if i.inputs != nil && row < len(i.inputs) {
		sel = i.inputs[row]
		set = true
	}
	i.inputsMutex.Unlock()
	if !set {
		i.updateInfo("")
		return
	}
	i.updateInfo(fmt.Sprintf("Input label: %s\nconnected: %v, favorite: %v, autoav: %v\nid: %s, appId: %s", sel.Label, sel.Connected, sel.Favorite, sel.Autoav, sel.Id, sel.AppId))
}

func (i *inputs) updateFromTv() (err error) {
	tvInputs, err := tv.TvGetExternalInputList()
	if err != nil {
		return errors.Wrap(err, "error updating inputs from TV")
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
	i.update(tvInputs)
	return nil
}

func (i *inputs) update(tvInputs []webostv.TvExternalInput) {
	i.inputsMutex.Lock()
	i.inputs = tvInputs
	i.inputsMutex.Unlock()

	i.Clear()
	for row, input := range tvInputs {
		i.SetCell(row, 0, tview.NewTableCell(input.Label))
		i.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", input.Connected)))
	}
	i.ScrollToBeginning()
}
