package main

import (
	"fmt"
	"fretcon/fretcon"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := append(os.Args[1:])
	if len(args) < 1 {
		panic("Usage: <instr-name-or-config-string> [data...]")
	}

	var f fretcon.Fretboard

	config := strings.Split(args[0], " ")
	cl := len(config)
	switch {
	case cl == 1:
		switch config[0] {
		default:
			panic("unknown instrument - " + args[0])
		case "guitar":
			f = fretcon.NewGuitar()
		case "shortguitar":
			f = fretcon.NewShortGuitar()
		case "bass":
			fallthrough
		case "bass4":
			f = fretcon.NewBass4()
		case "bass5":
			f = fretcon.NewBass5()
		}
	case cl > 3:
		begin, err := strconv.Atoi(config[0])
		e(err)

		end, err := strconv.Atoi(config[1])
		e(err)

		f, err = fretcon.NewFretboard(begin, end, config[2:]...)
		e(err)
	default:
		panic("invalid config argument format")
	}

	data := append(args[1:])
	val, err := f.Draw(data...)
	e(err)

	fmt.Println(val)
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
