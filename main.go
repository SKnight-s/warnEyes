package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	width       float32 = 300
	height      float32 = 200
	speedDraw   int     = 1
	defaultTime string  = "30"
)

type Timer struct {
	timeFocus    int // in seconds
	timeRemain   int // in seconds
	timeChan     chan string
	stopChan     chan bool
	timeProgress chan float64
	launched     bool
}

func (t *Timer) SetTimeFocus(time int) {
	t.timeFocus = time
	t.timeRemain = time
}

func (t *Timer) launchTimer() {
	tickTime := time.Duration(speedDraw) * time.Second
	log.Println("timer launched for", t.timeFocus, "seconds")
	timer := time.NewTimer(tickTime)
	t.timeChan <- "Stop (" + strconv.Itoa(t.timeRemain/60) + " min remaining)"
	t.timeProgress <- 0
	for {
		select {
		case <-t.stopChan:
			t.SetTimeFocus(t.timeFocus)
			t.timeChan <- "Start"
			return
		case <-timer.C:
			timer.Reset(tickTime)
			t.timeRemain -= speedDraw
			if t.timeRemain <= 0 {
				notify()
				t.timeChan <- "Finished !"
				t.timeRemain = t.timeFocus
				go t.launchTimer()
				return
			} else {
				if t.timeRemain%60 == 0 {
					t.timeChan <- "Stop (" + strconv.Itoa(t.timeRemain/60) + " min remaining)"
				}
				t.timeProgress <- (float64(t.timeFocus-t.timeRemain) / float64(t.timeFocus))

			}
		}
	}
}

func (t *Timer) LaunchTimer() {
	if !t.launched {
		t.launched = true
		go t.launchTimer()
	} else {
		t.StopTimer()
		t.LaunchTimer()
	}
}

func (t *Timer) StopTimer() {
	log.Println("timer stopped")
	t.stopChan <- true
	t.launched = false
	t.timeRemain = t.timeFocus
	t.timeChan <- "Start"
	t.timeProgress <- 0
}

func main() {
	a := app.New()
	w := a.NewWindow("WarnEyes")
	w.Resize(fyne.NewSize(width, height))
	icon, err := fyne.LoadResourceFromPath("eye.png")
	if err != nil {
		log.Println("error while loading icon")
	}
	w.SetIcon(icon)

	timer := Timer{timeFocus: 30, timeChan: make(chan string), stopChan: make(chan bool), timeProgress: make(chan float64)}
	var working bool

	focusTimeInput := widget.NewEntry()
	focusTimeInput.SetPlaceHolder("min")
	focusTimeInput.SetText(defaultTime)
	focusTimeInput.OnSubmitted = func(s string) {
		working = true
		timef, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
		if err != nil {
			focusTimeInput.SetText(defaultTime)
		}
		timer.SetTimeFocus(int(timef) * 60)
		timer.LaunchTimer()
	}

	button := widget.NewButton("Start", func() {
		working = !working
		if working { // we launch the timer
			timef, err := strconv.ParseInt(strings.TrimSpace(focusTimeInput.Text), 10, 32)
			if err != nil {
				focusTimeInput.SetText(defaultTime)
			}
			timer.SetTimeFocus(int(timef) * 60)
			timer.LaunchTimer()
		} else { // we stop the timer
			timer.StopTimer()
		}
	})

	progress := widget.NewProgressBar()

	go func() {
		for {
			select {
			case s := <-timer.timeChan:
				button.SetText(s)
			case p := <-timer.timeProgress:
				progress.SetValue(p)
			}
		}
	}()

	content := container.New(&appLayout{})
	content.Add(focusTimeInput)
	content.Add(progress)
	content.Add(button)
	w.SetContent(content)
	w.Canvas().Focus(focusTimeInput)
	w.ShowAndRun()
}

type appLayout struct{}

func (a *appLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(200, 200)
}

func (a *appLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var x, y, w, h float32
	for _, o := range objects {
		switch o.(type) {
		case *widget.Entry:
			minSize := o.MinSize()
			w = minSize.Width
			h = minSize.Height
			x = size.Width/2 - w/2
			y = size.Height * 0.1
			o.Resize(fyne.NewSize(w, h))
			o.Move(fyne.NewPos(x, y))
		case *widget.ProgressBar:
			w = 200
			h = 20
			x = size.Width/2 - w/2
			y = size.Height * 0.4
			o.Resize(fyne.NewSize(w, h))
			o.Move(fyne.NewPos(x, y))
		case *widget.Button:
			w = 250
			h = 50
			x = size.Width/2 - w/2
			y = size.Height * 0.6
			o.Resize(fyne.NewSize(w, h))
			o.Move(fyne.NewPos(x, y))
		}
	}

}

func min(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
