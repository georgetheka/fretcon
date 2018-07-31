package fretcon

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"sync"
)

type Fretboard interface {
	Draw(args ...string) (string, error)
}

type config struct {
	Strings []string
	Begin   int
	End     int
}

type fretboard struct {
	cfg  config
	data [][]string
	buff [][]string
	ll   int
	lock sync.Mutex
}

const (
	sp       = " "
	ret      = "\n"
	dot      = "o"
	fretwire = "|"
	nutstart = "||"
	nutend   = "}"
	se1      = "/"
	se2      = "\\"
	str      = "_"
	lastfret = 21
)

var (
	ErrStringNum     = errors.New("String number out of range")
	ErrFretNum       = errors.New("Fret number out of range")
	ErrStringRange   = errors.New("begin >= end")
	ErrNoStrings     = errors.New("at least one string must be provided")
	ErrInvalidArgNum = errors.New("invalid argument number (must be a multiple of 3)")
)

func NewFretboard(begin, end int, strs ...string) (*fretboard, error) {
	return newFretboard(config{
		Begin:   begin,
		End:     end,
		Strings: strs,
	})
}

func newFretboard(cfg config) (*fretboard, error) {
	if err := validateconfig(cfg); err != nil {
		return nil, err
	}
	f := new(fretboard)
	f.lock = sync.Mutex{}
	f.cfg = cfg

	f.ll = longestString(f.cfg.Strings)
	l := f.length()

	f.data = init2dArray(len(f.cfg.Strings), l)
	f.buff = init2dArray(len(f.cfg.Strings)+2, l)

	f.initBuffer()
	return f, nil
}

func NewGuitar() *fretboard {
	f, _ := NewFretboard(0, 25, "e'", "b", "g", "d", "a", "e")
	return f
}

func NewShortGuitar() *fretboard {
	f, _ := NewFretboard(0, 13, "e'", "b", "g", "d", "a", "e")
	return f
}

func NewBass4() *fretboard {
	f, _ := NewFretboard(0, 25, "g", "d", "a", "e")
	return f
}

func NewBass5() *fretboard {
	f, _ := NewFretboard(0, 25, "g", "d", "a", "e", "b")
	return f
}

func (f *fretboard) length() int {
	return f.cfg.End - f.cfg.Begin
}

func (f *fretboard) initBuffer() {
	l := len(f.buff)
	for i := 0; i < l; i++ {
		switch i {
		default:
			f.renderString(i)
		case 0:
			f.renderNums()
		case l - 1:
			f.renderDots()
		}
	}
}

func (f *fretboard) renderString(snum int) {
	s := f.buff[snum]
	label := f.cfg.Strings[snum-1]
	plabel := label + strings.Repeat(sp, f.ll-len(label))
	fs := strings.Repeat(str, 3)
	for i, j := 0, f.cfg.Begin; i < f.length(); i, j = i+1, j+1 {
		var v string
		if j == 0 {
			v = plabel + fs + f.fretWire(snum, 0)
		} else if i == 0 {
			v = plabel + f.fretWire(snum, 0) + fs + fretwire
		} else {
			v = fs + f.fretWire(snum, j)
		}
		s[i] = v
	}
}

func (f *fretboard) fretWire(snum, num int) string {
	if num == 0 {
		if f.cfg.Begin == 0 {
			return nutstart
		}
		if snum%2 == 0 {
			return se1
		}
		return se2
	}
	if num == f.cfg.End-1 {
		if f.cfg.End > lastfret {
			return nutend
		}
		if snum%2 == 0 {
			return se2
		}
		return se1
	}
	return fretwire
}

func (f *fretboard) renderNums() {
	s := f.buff[0]
	for i, j := 0, f.cfg.Begin; i < f.length(); i, j = i+1, j+1 {
		var v string
		if j == 0 {
			v = strings.Repeat(sp, f.ll) + renderData(strconv.Itoa(j), sp) + sp + sp
		} else if i == 0 {
			v = strings.Repeat(sp, f.ll+1) + renderData(strconv.Itoa(j), sp) + sp
		} else {
			v = renderData(strconv.Itoa(j), sp) + sp
		}
		s[i] = v
	}
}

func (f *fretboard) renderDots() {
	s := f.buff[len(f.buff)-1]
	b := f.cfg.Begin
	for i, j := 0, b; i < f.length(); i, j = i+1, j+1 {
		var data string
		switch j {
		default:
			data = sp
		case 0:
			data = strings.Repeat(sp, f.ll)
		case 3, 5, 7, 9, 12, 15, 17, 19, 21, 23:
			data = dot
		}
		if i == 0 {
			s[i] = strings.Repeat(sp, f.ll+1) + renderData(data, sp) + sp
		} else {
			s[i] = renderData(data, sp) + sp
		}
	}
}

func renderData(data, ch string) string {
	switch len(data) {
	case 0:
		return strings.Repeat(ch, 3)
	case 1:
		return ch + data + ch
	case 2:
		return ch + data
	case 3:
		return data
	default:
		return data[:3]
	}
}

func (f *fretboard) Draw(args ...string) (string, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if len(args)%3 != 0 {
		return "", ErrInvalidArgNum
	}

	for i := 0; i < len(args); i += 3 {
		snum, err := strconv.Atoi(args[i])
		if err != nil {
			return "", err
		}

		fnum, err := strconv.Atoi(args[i+1])
		if err != nil {
			return "", err
		}

		if snum < 1 || snum > len(f.buff) {
			return "", ErrStringNum
		}

		if fnum < f.cfg.Begin || fnum > f.cfg.End {
			return "", ErrFretNum
		}

		f.data[snum-1][fnum-f.cfg.Begin] = args[i+2]
	}
	return f.String(), nil
}

func (f *fretboard) String() string {
	var buffer bytes.Buffer
	l := len(f.buff)
	for i := 0; i < l; i++ {
		line := f.buff[i]
		for j := 0; j < len(line); j++ {
			switch i {
			case 0, l - 1:
				buffer.WriteString(line[j])
			default:
				text := f.data[i-1][j]
				if text == "" {
					buffer.WriteString(line[j])
				} else if j == 0 {
					s := line[j]
					if f.cfg.Begin == 0 { //first fret of a full range neck
						buffer.WriteString(s[:len(s)-5])
						buffer.WriteString(renderData(text, str))
						buffer.WriteString(nutstart)
					} else {
						buffer.WriteString(s[:len(s)-5])
						buffer.WriteString(string(s[len(s)-5]))
						buffer.WriteString(renderData(text, str))
						buffer.WriteString(fretwire)
					}
				} else {
					buffer.WriteString(renderData(text, str))
					buffer.WriteString(fretwire)
				}
			}
		}
		buffer.WriteString(ret)
	}
	return buffer.String()
}

func validateconfig(cfg config) error {
	if cfg.Begin >= cfg.End {
		return ErrStringRange
	}
	if len(cfg.Strings) == 0 {
		return ErrNoStrings
	}
	return nil
}

func init2dArray(x, y int) [][]string {
	arr := make([][]string, x)
	for i := range arr {
		arr[i] = make([]string, y)
	}
	return arr
}

func longestString(strs []string) int {
	ll := len(strs[0]) + 1
	for _, s := range strs {
		if ol := len(s); ol > ll {
			ll = ol
		}
	}
	return ll
}
