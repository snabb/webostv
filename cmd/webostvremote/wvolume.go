package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type volume struct {
	*Slider
}

func newVolume() *volume {
	w := NewSlider()
	w.SetBorder(true)
	w.SetTitle("Volume")
	return &volume{w}
}

func (v *volume) update(volume int) {
	v.SetPercent(volume)
}

func (v *volume) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		key := event.Key()
		var kr rune
		if key == tcell.KeyRune {
			kr = event.Rune()
		}
		switch {
		case key == tcell.KeyRight || (key == tcell.KeyRune && kr == '+'):
			go tv.AudioVolumeUp()
			// XXX check err
		case key == tcell.KeyLeft || (key == tcell.KeyRune && kr == '-'):
			go tv.AudioVolumeDown()
			// XXX check err
		}
	})
}
