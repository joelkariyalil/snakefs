package main

import (
	"log"
	"math/rand"
	"time"

	. "github.com/rthornton128/goncurses"
)

const (
	WindowSize = 50
)

type pos struct{ x, y int }

type window struct {
	win        *Window
	windowSize pos
}

type food struct{ foodPos pos }

type snake struct {
	body         []pos            // ordered: tail -> head
	snakePos     map[pos]struct{} // O(1) lookup
	direction    pos
	headPos      pos
	isGrowLength bool
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	start := pos{0, 0}
	s := &snake{
		body: []pos{start},
		snakePos: map[pos]struct{}{
			start: {},
		},
		headPos:   start,
		direction: pos{0, 1},
	}

	w := &window{windowSize: pos{WindowSize, WindowSize}}

	stdscr, err := Init()
	if err != nil {
		panic(err)
	}
	defer End()

	Raw(true)
	Echo(false)
	Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	my, mx := stdscr.MaxYX()
	log.Printf("my :%d\nmx: %d", my, mx)

	win, _ := NewWindow(w.windowSize.y*7/10, w.windowSize.x, 0, 0)
	win.Keypad(true)
	w.win = win

	f := &food{foodPos: randomFoodPos(w.windowSize, s.snakePos, rng)}

	stdscr.Refresh()

	ch := make(chan Key)
	go func() {
		for {
			ch <- stdscr.GetChar()
		}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case key := <-ch:
			switch key {
			case KEY_DOWN:
				s.direction = pos{0, 1}
			case KEY_UP:
				s.direction = pos{0, -1}
			case KEY_LEFT:
				s.direction = pos{-1, 0}
			case KEY_RIGHT:
				s.direction = pos{1, 0}
			case KEY_BACKSPACE:
				time.Sleep(time.second * 100)
			case 'q':
				return
			}

		case <-ticker.C:
			updateSnake(s, w, f, rng)
			printGrid(s, w, f)

			stdscr.MovePrintf(
				my-2,
				0,
				"Snake head=%v len=%d dir=%v food=%v",
				s.headPos,
				len(s.body),
				s.direction,
				f.foodPos,
			)
			stdscr.Refresh()
		}
	}
}

func updateSnake(s *snake, w *window, f *food, rng *rand.Rand) {
	// eat food
	if s.headPos == f.foodPos {
		s.isGrowLength = true
		f.foodPos = randomFoodPos(w.windowSize, s.snakePos, rng)
	}

	head := s.body[len(s.body)-1]
	newHead := pos{head.x + s.direction.x, head.y + s.direction.y}

	// wrap around boundaries
	if newHead.x < 0 {
		newHead.x = w.windowSize.x - 1
	}
	if newHead.x >= w.windowSize.x {
		newHead.x = 0
	}
	if newHead.y < 0 {
		newHead.y = w.windowSize.y - 1
	}
	if newHead.y >= w.windowSize.y {
		newHead.y = 0
	}

	// add new head
	s.body = append(s.body, newHead)
	s.snakePos[newHead] = struct{}{}
	s.headPos = newHead

	if !s.isGrowLength {
		tail := s.body[0]
		s.body = s.body[1:]
		delete(s.snakePos, tail)
	} else {
		s.isGrowLength = false
	}
}

func printGrid(s *snake, w *window, f *food) {
	for i := 0; i < w.windowSize.x; i++ {
		for j := 0; j < w.windowSize.y; j++ {
			p := pos{i, j}

			if p == f.foodPos {
				w.win.MovePrint(j, i, "*")
				continue
			}

			if _, ok := s.snakePos[p]; ok {
				if p == s.headPos {
					w.win.MovePrint(j, i, "@")
				} else {
					w.win.MovePrint(j, i, "m")
				}
			} else {
				w.win.MovePrint(j, i, " ")
			}
		}
	}
	w.win.Refresh()
}

func randomFoodPos(size pos, occupied map[pos]struct{}, rng *rand.Rand) pos {
	for {
		p := pos{
			x: rng.Intn(size.x),
			y: rng.Intn(size.y),
		}
		if _, ok := occupied[p]; !ok {
			return p
		}
	}
}
