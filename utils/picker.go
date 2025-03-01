package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// FuzzyListPicker lets the user search and navigate the strings in ss. It returns the
// currently selected string on <enter>.
func FuzzyListPicker(t, p string, ss []string) (string, error) {

	SetUnbufferedInput()
	defer SetBufferedInput()

	input := ""
	buf := make([]byte, 2)
	fss := ss
	current := 0

	for {
		if current >= len(fss) {
			current = len(fss) - 1
		}

		RenderClear()

		fmt.Printf("%s", t)

		for i, f := range fss {
			marker := " "
			col := NORM
			if i == current {
				marker = NOTICE + ">" + NORM
				col = SELECTED
			}

			fmt.Printf(" %s %s%v\n"+NORM, marker, col, f)
		}

		fmt.Printf("%s > %s", SUBTITLE+p+NORM, SELECTED+input+NORM)

		_, _ = os.Stdin.Read(buf)
		r := rune(buf[0])

		switch {
		// TODO: Should add `?` for help
		case r == 65: // up
			current--
			if current < 0 {
				current = 0
			}
		case r == 66: // down
			current++
			if current >= len(fss) {
				current = len(fss) - 1
			}
		case r == '\n': // enter
			if len(fss) < 1 {
				return "", errors.New("nothing selected")
			}
			return fss[current], nil
		case r == '\x7f': // bksp
			current = 0
			fss = ss
			if len(input) > 0 {
				input = string(input[0 : len(input)-1])
			}
		case r == 27: // esc
			// do nothing by design (up and down arrows are spooky)
		default:
			input = input + string(r)
		}

		fss = filter(fss, input)
	}
}

func filter(ss []string, m string) []string {

	re, err := regexp.Compile(strings.Join(strings.Split(strings.ToLower(m), ""), ".*"))
	if err != nil {
		return []string{}
	}

	ms := []string{}

	for _, s := range ss {
		if re.MatchString(strings.ToLower(s)) {
			ms = append(ms, s)
		}
	}

	return ms
}

func ListPicker(p string, ss []string) (int, string) {

	for i, f := range ss {
		fmt.Printf("%4d: %v\n", i+1, f)
	}

	fmt.Printf("%s > ", p)

	r := bufio.NewReader(os.Stdin)
	input, _, _ := r.ReadLine()

	n, err := strconv.Atoi(string(input))
	if err != nil || n < 1 || n > len(ss) {
		return -1, ""
	}

	return n, ss[n-1]
}
