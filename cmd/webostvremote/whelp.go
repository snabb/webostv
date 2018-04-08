package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func newHelp() tview.Primitive {
	w := tview.NewTextView()
	w.SetBorder(true)
	w.SetScrollable(true)
	w.SetTitle("Help")
	w.SetWrap(false)
	fmt.Fprintln(w, "Tab ⇥ / ⇤ next / prev")
	fmt.Fprintln(w, "V         volume")
	fmt.Fprintln(w, "C         channels")
	fmt.Fprintln(w, "I         inputs")
	fmt.Fprintln(w, "A         apps")
	fmt.Fprintln(w, "Enter     select")
	fmt.Fprintln(w, "arrows    move")
	fmt.Fprintln(w, "Q / Esc   quit")
	fmt.Fprintln(w, "Ctrl+X    turn off+quit\n")
	fmt.Fprintln(w, "webostvremote © J.Snabb 2018")
	fmt.Fprint(w, "github.com/snabb/webostv")
	w.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w.ScrollToBeginning()
		return x + 1, y + 1, width - 2, height - 2
	})
	return w
}
