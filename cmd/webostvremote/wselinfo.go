package main

import (
	"github.com/rivo/tview"
)

type selInfo struct {
	*tview.TextView
}

func newSelInfo() *selInfo {
	w := tview.NewTextView()
	w.SetBorder(true)
	w.SetScrollable(true)
	w.SetTitle("Selection Info")
	w.SetWrap(true)
	w.SetWordWrap(true)

	s := &selInfo{TextView: w}
	return s
}

func (s *selInfo) update(str string) {
	s.SetText(str)
	s.ScrollToBeginning()
}
