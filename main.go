package main

import "C"
import (
	"log"
	"time"

	. "github.com/rthornton128/goncurses"
)

const (
	WindowSize = 60
)

type pos struct{ x, y int }

type window struct {
	win        *Window
	windowSize pos
}

type food struct{ foodPos pos }

type snake struct {
	snakePos  pos
	direction pos
	length    int
	width     int
}

func main() {

	s := &snake{
		snakePos:  pos{0, 0},
		length:    0,
		width:     0,
		direction: pos{0, 1},
	}

	w := &window{windowSize: pos{WindowSize * 0.7, WindowSize}}

	f := &food{foodPos: pos{w.windowSize.x / 2, w.windowSize.y / 2}}

	stdscr, err := Init()
	if err != nil {
		panic(err)
	}
	defer End()

	Raw(true)
	Echo(true)
	Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	my, mx := stdscr.MaxYX()
	log.Printf("my :%d\nmx: %d", my, mx)

	win, _ := NewWindow(w.windowSize.x, w.windowSize.y, 0, 0)
	win.Keypad(true)
	w.win = win

	stdscr.Refresh()

	printGrid(s, w, f)

	ch := make(chan Key)

	go func() {
		for {
			ch <- stdscr.GetChar()
		}
	}()

	for {
		// receive the value from the channel
		key := <-ch
		switch key {

		case KEY_UP:
			s.direction = pos{0, 1}

		case KEY_DOWN:
			s.direction = pos{0, -1}

		case KEY_LEFT:
			s.direction = pos{-1, 0}

		case KEY_RIGHT:
			s.direction = pos{1, 0}

		case 'q':
			return
		}
		printGrid(s, w, f)
		time.Sleep(time.Millisecond)

		stdscr.MovePrintf(my-2, 0, "Snake Details: %v", s)
		stdscr.Refresh()
	}
}

func printGrid(s *snake, w *window, f *food) {
	if s.snakePos.x == f.foodPos.x && s.snakePos.y == f.foodPos.y {
		s.length++
	}

	s.snakePos.x = (s.snakePos.x + s.direction.x) % w.windowSize.x
	s.snakePos.y = (s.direction.y + s.snakePos.y) % w.windowSize.y

	for i := 0; i < w.windowSize.x; i++ {
		for j := 0; j < w.windowSize.y; j++ {

			if i == f.foodPos.x && j == f.foodPos.y {
				w.win.MovePrint(j, i, "*")
			} else if i == s.snakePos.x && j == s.snakePos.y {
				w.win.MovePrint(j, i, "=")
			} else {
				w.win.MovePrint(j, i, ".")
			}
		}
	}
	w.win.Refresh()
}
