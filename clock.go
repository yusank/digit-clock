package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	globalWidth  = 16
	globalHeight = 9
	timeFormat   = defTimeFormat
)

const (
	dateFormat    = "2006/01/02 15:04:05"
	secondFormat  = "15:04:05"
	defTimeFormat = "15:04"
)

type pos struct {
	x, y int
}

func main() {
	logFile, err := os.OpenFile("stderr.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)
	var (
		stop       = make(chan bool, 0)
		timeStr    = make(chan string, 0)
		isShowSec  bool
		isShowDate bool
	)

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.BoolVar(&isShowSec, "s", false, "is show second")
	flagSet.BoolVar(&isShowDate, "t", false, "is show date")
	_ = flagSet.Parse(os.Args[1:])

	if isShowSec {
		timeFormat = secondFormat
	}

	if isShowDate {
		timeFormat = dateFormat
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for t := range ticker.C {
			timeStr <- t.Format(timeFormat)
		}
	}()

	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyCtrlQ, termbox.KeyCtrlC:
					stop <- true
					return
				}
			}
		}
	}()

loop:
	for {
		select {
		case str := <-timeStr:
			draw(str)
		case <-stop:
			break loop
		}

	}
}

func checkPosition(s rune, p pos) bool {
	block := 2
	if globalWidth < 9 {
		block = 1
	}
	var (
		midX = (globalWidth - block) / 2
		midY = globalHeight / 2

		lastX   = globalWidth - block - 2
		lastY   = globalHeight - 1
		lineWid = block
	)

	// 留 2 列空
	if p.x > lastX {
		return false
	}

	switch s {
	case '0':
		if p.x < lineWid {
			return true
		}

		if p.x > lastX-lineWid {
			return true
		}

		return p.y == 0 || p.y == lastY
	case '1':
		return p.x > lastX-lineWid
	case '2':
		if p.y == 0 || p.y == midY || p.y == lastY {
			return true
		}

		if p.y < midY {
			return p.x > lastX-lineWid
		}

		return p.y > midY && p.x < lineWid
	case '3':
		if p.y == 0 || p.y == midY || p.y == lastY {
			return true
		}

		return p.x > lastX-lineWid
	case '4':
		if p.y == midY {
			return true
		}

		if p.x > lastX-lineWid {
			return true
		}

		return p.y < midY && p.x < lineWid
	case '5':
		if p.y == 0 || p.y == midY || p.y == lastY {
			return true
		}

		if p.y < midY {
			return p.x < lineWid
		}

		return p.y > midY && p.x > lastX-lineWid
	case '6':
		if p.y == 0 || p.y == midY || p.y == lastY {
			return true
		}

		if p.x < lineWid {
			return true
		}

		return p.y > midY && p.x > lastX-lineWid
	case '7':
		if p.y == 0 {
			return true
		}
		return p.x > lastX-lineWid
	case '8':
		if p.y == 0 || p.y == midY || p.y == lastY {
			return true
		}

		return p.x < lineWid || p.x > lastX-lineWid
	case '9':
		if p.y == 0 || p.y == midY || p.y == lastY {
			return true
		}
		if p.x > lastX-lineWid {
			return true
		}
		return p.y < midY && p.x < lineWid
	case ':':
		if p.x > midX || p.x < midX-lineWid {
			return false
		}

		return p.y == midY+1 || p.y == midY-1
	case '/':
		dot := (p.x + p.y) / 2
		log.Println(dot, midX)
		return dot == midX
	}

	return false
}

func draw(val string) {
	w, _ := termbox.Size()
	if w/len(timeFormat) < globalWidth {
		globalWidth = w / len(timeFormat)
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	var width = globalWidth
	for i, r := range val {
		tempW := width * (i + 1)
		startX := width * i
		for y := 0; y < globalHeight; y++ {
			for x := startX; x < tempW; x++ {
				p := pos{x: x % width, y: y}
				if checkPosition(r, p) {
					termbox.SetCell(x, y+5, '▒', termbox.ColorDefault,
						termbox.ColorCyan)
				}
			}
		}
	}
	termbox.Flush()
}
