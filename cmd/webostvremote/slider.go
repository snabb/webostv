package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"sync"
)

type Slider struct {
	*tview.Box
	percent             int
	percentMutex        sync.Mutex
	BarColor            tcell.Color
	NoBarColor          tcell.Color
	IndicatorOnBarColor tcell.Color
	IndicatorNoBarColor tcell.Color
	IndicatorFormat     string
	IndicatorAlign      int
}

func NewSlider() *Slider {
	return &Slider{
		Box:                 tview.NewBox(),
		BarColor:            tview.Styles.GraphicsColor,
		NoBarColor:          tview.Styles.ContrastBackgroundColor,
		IndicatorOnBarColor: tview.Styles.ContrastSecondaryTextColor,
		IndicatorNoBarColor: tview.Styles.SecondaryTextColor,
		IndicatorFormat:     "%d%%",
		IndicatorAlign:      tview.AlignCenter,
	}
}

func (s *Slider) SetPercent(v int) {
	s.percentMutex.Lock()
	s.percent = v
	s.percentMutex.Unlock()
}

func (s *Slider) GetPercent() (v int) {
	s.percentMutex.Lock()
	v = s.percent
	s.percentMutex.Unlock()
	return v
}

func (s *Slider) Draw(screen tcell.Screen) {
	s.Box.Draw(screen)
	x, y, width, height := s.GetInnerRect()

	if height < 1 {
		return
	}
	percent := s.GetPercent()

	var indicator []rune
	if s.IndicatorFormat != "" {
		indicator = []rune(fmt.Sprintf(s.IndicatorFormat, percent))
	}
	var ipos int
	if indicator != nil {
		switch s.IndicatorAlign {
		case tview.AlignLeft:
			ipos = 0
		case tview.AlignCenter:
			ipos = (width - len(indicator)) / 2
		case tview.AlignRight:
			ipos = width - len(indicator) - 1
		}
	}

	for i := 0; i < width; i++ {
		c := ' '
		var style tcell.Style
		if i*100/width >= percent {
			style = tcell.StyleDefault.
				Foreground(s.IndicatorNoBarColor).
				Background(s.NoBarColor)
		} else {
			style = tcell.StyleDefault.
				Foreground(s.IndicatorOnBarColor).
				Background(s.BarColor)
		}
		if indicator != nil && i >= ipos && i-ipos < len(indicator) {
			c = indicator[i-ipos]
		}
		screen.SetContent(x+i, y, c, nil, style)
	}
}
