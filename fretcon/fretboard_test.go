package fretcon

import (
	"fmt"
	"strings"
	"testing"
)

type Case struct {
	label       string
	begin       int
	end         int
	strs        []string
	input       []string
	errExpected bool
	errMessage  string
	expected    string
}

type Cases []Case

func TestNewFretboard(t *testing.T) {
	t.Parallel()
	cases := Cases{
		{
			"triggers 'string range error' error",
			9,
			5,
			[]string{},
			[]string{},
			true,
			ErrStringRange.Error(),
			"",
		},
		{
			"triggers 'no strings' error",
			0,
			12,
			[]string{},
			[]string{},
			true,
			ErrNoStrings.Error(),
			"",
		},
		{
			"returns expected instance",
			0,
			25,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{},
			false,
			"",
			"",
		},
		{
			"returns expected instance",
			0,
			25,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{},
			false,
			ErrNoStrings.Error(),
			"",
		},
		{
			"returns expected instance for partial fretboard",
			4,
			11,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{},
			false,
			"",
			"",
		},
	}

	for _, c := range cases {
		fmt.Println("CASE", c.label)
		f, err := NewFretboard(c.begin, c.end, c.strs...)
		if c.errExpected {
			if err == nil || err.Error() != c.errMessage {
				t.Errorf("NewFretboard(...) returned incorrect or nil error %q, expected %q",
					err.Error(), c.errMessage)
			}
			continue
		}

		if err != nil {
			t.Errorf("NewFretboard(...) returned an unexpected error %q", err.Error())
			continue
		}

		cfgTest(t, f, c.begin, c.end, c.strs...)

		if len(f.data) != len(f.cfg.Strings) {
			t.Errorf("len(data) returned %v, expected %v", len(f.data), len(c.strs))
		}

		if len(f.data[0]) != f.cfg.End-f.cfg.Begin {
			t.Errorf("len(data[0]) returned %v, expected %v", len(f.data[0]), c.end)
		}

		if len(f.buff) != len(f.cfg.Strings)+2 {
			t.Errorf("len(buff) returned %v, expected %v", len(f.buff), len(c.strs))
		}

		if len(f.buff[0]) != f.cfg.End-f.cfg.Begin {
			t.Errorf("len(data[0]) returned %v, expected %v", len(f.buff[0]), c.end)
		}
	}
}

func TestDraw(t *testing.T) {
	t.Parallel()
	cases := Cases{
		{
			"empty standard guitar fretboard, full length",
			0,
			25,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{},
			false,
			"",
			`
-
....0....1...2...3...4...5...6...7...8...9...10..11..12..13..14..15..16..17..18..19..20..21..22..23..24.
e'.___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
b..___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
g..___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
d..___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
a..___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
e..___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
.................o.......o.......o.......o...........o...........o.......o.......o.......o.......o......
-`,
		},
		{
			"full length standard guitar fretboard, valid data",
			0,
			25,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{
				"1", "0", "xyz",
				"2", "1", "1b",
				"3", "2", "2",
				"4", "3", "4",
				"5", "3", "",
				"6", "1", "1babc",
			},
			false,
			"",
			`
-
....0....1...2...3...4...5...6...7...8...9...10..11..12..13..14..15..16..17..18..19..20..21..22..23..24.
e'.xyz||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
b..___||_1b|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
g..___||___|_2_|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
d..___||___|___|_4_|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
a..___||___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
e..___||1ba|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___|___}
.................o.......o.......o.......o...........o...........o.......o.......o.......o.......o......
-`,
		},
		{
			"partial standard guitar fretboard from the nut, valid data",
			0,
			5,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{
				"1", "4", "4",
				"2", "1", "1",
				"3", "2", "2",
				"4", "3", "4",
				"5", "3", "3",
				"6", "1", "1babc",
			},
			false,
			"",
			`
-
....0....1...2...3...4..
e'.___||___|___|___|_4_|
b..___||_1_|___|___|___\
g..___||___|_2_|___|___/
d..___||___|___|_4_|___\
a..___||___|___|_3_|___/
e..___||1ba|___|___|___\
.................o......
-
			`,
		},
		{
			"partial middle range guitar fretboard, valid data",
			3,
			10,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{
				"1", "3", "1b",
				"2", "4", "2",
				"3", "5", "3",
				"4", "3", "x",
				"5", "3", "x",
				"6", "3", "1b",
			},
			false,
			"",
			`
-
.....3...4...5...6...7...8...9..
e'.\_1b|___|___|___|___|___|___/
b../___|_2_|___|___|___|___|___\
g..\___|___|_3_|___|___|___|___/
d../_x_|___|___|___|___|___|___\
a..\_x_|___|___|___|___|___|___/
e../_1b|___|___|___|___|___|___\
.....o.......o.......o.......o..
-
			`,
		},
		{
			"partial middle range standard guitar fretboard, no data",
			3,
			11,
			[]string{"e'", "b", "g", "d", "a", "e"},
			[]string{},
			false,
			"",
			`
-
.....3...4...5...6...7...8...9...10.
e'.\___|___|___|___|___|___|___|___/
b../___|___|___|___|___|___|___|___\
g..\___|___|___|___|___|___|___|___/
d../___|___|___|___|___|___|___|___\
a..\___|___|___|___|___|___|___|___/
e../___|___|___|___|___|___|___|___\
.....o.......o.......o.......o......
-
			`,
		},
		{
			"triggers 'string number out of range' error",
			0,
			20,
			[]string{"e'", "b"},
			[]string{
				"0", "3", "4",
			},
			true,
			ErrStringNum.Error(),
			"",
		},
		{
			"triggers 'fret number out of range' error",
			0,
			20,
			[]string{"e'", "b"},
			[]string{
				"1", "22", "4",
			},
			true,
			ErrFretNum.Error(),
			"",
		},
		{
			"triggers 'invalid argument number' error",
			0,
			20,
			[]string{"e'", "b"},
			[]string{
				"1", "2",
			},
			true,
			ErrInvalidArgNum.Error(),
			"",
		},
	}

	for _, c := range cases {
		fmt.Println("CASE", c.label)
		expected := strings.TrimSpace(
			strings.Replace(
				strings.Replace(c.expected, ".", " ", -1),
				"\r", "", -1))

		f, _ := NewFretboard(c.begin, c.end, c.strs...)
		res, err := f.Draw(c.input...)

		if c.errExpected {
			if err == nil || err.Error() != c.errMessage {
				t.Errorf("f.Draw(%v) returned incorrect or nil error %q, expected %q",
					c.input, err.Error(), c.errMessage)
			}
			continue
		}

		if err != nil {
			t.Errorf("f.Draw(...) returned an unexpected error %q", err.Error())
			continue
		}

		actual := strings.TrimSpace("-\n" + res + "-")
		if actual != expected {
			t.Errorf("f.Draw(...) expected:\n\"%v\"\nreceived:\n\"%v\"",
				strings.Replace(expected, " ", ".", -1),
				strings.Replace(actual, " ", ".", -1))
		}
	}
}

func TestNewGuitar(t *testing.T) {
	t.Parallel()
	f := NewGuitar()
	cfgTest(t, f, 0, 25, "e'", "b", "g", "d", "a", "e")
}

func TestNewShortGuitar(t *testing.T) {
	t.Parallel()
	f := NewShortGuitar()
	cfgTest(t, f, 0, 13, "e'", "b", "g", "d", "a", "e")
}

func TestNewBass4(t *testing.T) {
	t.Parallel()
	f := NewBass4()
	cfgTest(t, f, 0, 25, "g", "d", "a", "e")
}

func TestNewBass5(t *testing.T) {
	t.Parallel()
	f := NewBass5()
	cfgTest(t, f, 0, 25, "g", "d", "a", "e", "b")
}

func cfgTest(t *testing.T, f *fretboard, b, e int, strs ...string) {
	if f.cfg.Begin != b {
		t.Errorf("cfg.Begin returned %v, expected %v", f.cfg.Begin, b)
	}
	if f.cfg.End != e {
		t.Errorf("cfg.End returned %v, expected %v", f.cfg.End, e)
	}
	for i, e := range f.cfg.Strings {
		if e != strs[i] {
			t.Errorf("cfg.Strings[%d] returned %q, expected %q", i, e, strs[i])
		}
	}
}
