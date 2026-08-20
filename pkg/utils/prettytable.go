package utils

import (
	"bytes"
	"github.com/kluctl/kluctl/lib/term"
	"sort"
	"strings"
	"unicode/utf8"
)

type Str []rune
type Row []Str

type PrettyTable struct {
	rows []Row
}

func (t *PrettyTable) AddRow(c ...string) {
	formatted := make([]Str, len(c))
	for i, str := range c {
		formatted[i] = Str(str)
	}
	t.rows = append(t.rows, formatted)
}

func (t *PrettyTable) SortRows(col int) {
	sort.SliceStable(t.rows[1:], func(i, j int) bool {
		return string(t.rows[i+1][col]) < string(t.rows[j+1][col])
	})
}

func (t *PrettyTable) Render(limitWidths []int) string {
	cols := len(t.rows[0])

	maxWidth := func(col int, maxW int) int {
		w := 0
		for _, l := range t.rows {
			count := 0
			for _, c := range l[col] {
				count++
				if c == '\n' {
					if maxW != -1 && count > maxW {
						return maxW
					} else if count > w {
						w = count
					}
					count = 0
				}
			}
			if maxW != -1 && count > maxW {
				return maxW
			} else if count > w {
				w = count
			}
		}
		return w
	}
	subStr := func(str Str, s int, e int) string {
		if s > len(str) {
			s = len(str)
		}
		if e > len(str) {
			e = len(str)
		}
		return string(str[s:e])
	}

	widths := make([]int, cols)
	widthSum := 0
	for i := 0; i < cols; i++ {
		w := -1
		if i < len(limitWidths) {
			w = limitWidths[i]
		}
		widths[i] = maxWidth(i, w)
		if i != cols-1 {
			widthSum += widths[i]
		}
	}

	if len(limitWidths) < cols {
		tw := term.GetWidth()
		// last column should use all remaining space
		tw = tw - widthSum - (cols-1)*3 - 4
		if tw <= 0 {
			tw = 1
		}
		widths[len(limitWidths)] = tw
	}

	hsep := "+-"
	for i := 0; i < cols; i++ {
		hsep += strings.Repeat("-", widths[i])
		if i != cols-1 {
			hsep += "-+-"
		}
	}
	hsep += "-+\n"

	buf := bytes.NewBuffer(nil)
	buf.WriteString(hsep)
	pos := make([]int, cols)
	for _, l := range t.rows {
		for i := 0; i < cols; i++ {
			pos[i] = 0
		}

		for {
			anyLess := false
			for i := 0; i < cols; i++ {
				if pos[i] < len(l[i]) {
					anyLess = true
				}
			}
			if !anyLess {
				break
			}

			buf.WriteString("| ")
			for i := 0; i < cols; i++ {
				x := subStr(l[i], pos[i], pos[i]+widths[i])
				newLine := strings.IndexRune(x, '\n')
				if newLine != -1 {
					x = x[:newLine]
					pos[i] += 1
				}
				lenX := utf8.RuneCountInString(x)
				pos[i] += lenX
				buf.WriteString(x)
				buf.WriteString(strings.Repeat(" ", widths[i]-lenX))
				if i != cols-1 {
					buf.WriteString(" | ")
				}
			}
			buf.WriteString(" |\n")
		}
		buf.WriteString(hsep)
	}
	return buf.String()
}
