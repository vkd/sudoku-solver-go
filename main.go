package main

import (
	"fmt"
	"strconv"
	"strings"
)

const initPuzzle = `
8..|..6|3.5
.4.|...|.7.
...|...|...
---+---+---
.1.|.38|7.4
...|1.4|...
3..|.7.|29.
---+---+---
...|..3|...
.2.|...|.4.
5.6|8..|..2
`

func main() {
	// s := NewSudoku(`
	// 8..|4.6|3.5
	// .4.|...|.7.
	// ...|...|4..
	// ---+---+---
	// .1.|.38|7.4
	// ...|1.4|...
	// 384|675|291
	// ---+---+---
	// 4..|..3|...
	// .2.|...|.4.
	// 5.6|84.|..2
	// `)
	s := NewSudoku(`
	87.|4.6|325
	.4.|3..|17.
	..3|7..|4..
	---+---+---
	.1.|.38|7.4
	...|1.4|.3.
	384|675|291
	---+---+---
	49.|..3|..7
	.2.|...|.43
	536|847|912
	`)
	fmt.Printf("%s", s.String())
	fmt.Printf("---------------------------------------\n")
	var rs []Ruler = DefaultRules()
	// rs = append(rs, Pair{R: [2]int{2, 4}, C: [2]int{1, 1}, Ds: [2]Digit{5, 6}})
	// rs = append(rs, Pair{R: [2]int{1, 2}, C: [2]int{4, 4}, Ds: [2]Digit{5, 8}})
	var still bool = true
	for still {
		still = false
		for i := range rs {
			var ok bool
			s, rs, ok = rs[i].TryReduce(s, rs)
			if ok {
				still = true
				// fmt.Printf("%s", s.String())
				fmt.Printf("--- still ---\n")
			}
		}
		// if !still {
		// 	for i := 0; i < 9; i++ {
		// 		for j := 0; j < 9; j++ {
		// 			var c int
		// 			var cd Digit
		// 			for d := Digit(1); d <= MaxDigit; d++ {
		// 				test := testDigit(s, rs, i, j, d)
		// 				if test {
		// 					c++
		// 					cd = d
		// 				}
		// 			}
		// 			if c == 1 {
		// 				still = true
		// 				s = s.Set(i, j, cd)
		// 				fmt.Printf("--- ok ---\n")
		// 				// fmt.Printf("%s", s.String(0, 0))
		// 			}
		// 		}
		// 	}
		// }
	}
	fmt.Printf("---------------------------------------\n")
	fmt.Printf("%s", s.String())
}

func appendUniq(rs []Ruler, r Ruler) ([]Ruler, bool) {
	for i := range rs {
		if rs[i] == r {
			return rs, false
		}
	}
	return append(rs, r), true
}

func colsEquals(c1, c2 []int) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i := range c1 {
		if c1[i] != c2[i] {
			return false
		}
	}
	return true
}

type Digit uint8

const MaxDigit Digit = 9

type Sudoku struct {
	grid [9][9]Digit
}

func NewSudoku(s string) (out Sudoku) {
	var row int
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if row == 0 && len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "-") {
			continue
		}
		var col int
		for _, c := range line {
			switch c {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				d, _ := strconv.Atoi(string(c))
				out.grid[row][col] = Digit(d)
				col++
			case '.':
				col++
			case '|':
			}
		}
		row++
	}
	return
}

func (s Sudoku) Set(row, column int, d Digit) Sudoku {
	if s.grid[row][column] != 0 {
		panic("wrong sudoku set")
	}
	s.grid[row][column] = d
	fmt.Printf("reduced: (%dx%d) -> %d\n", row, column, d)
	fmt.Printf("%s", s.StringMark(row, column))
	return s
}

func DefaultRules() (out []Ruler) {
	for i := 0; i < 9; i++ {
		out = append(out, RowRule{Row: i})
		out = append(out, ColumnRule{Column: i})
	}
	for d := Digit(1); d <= MaxDigit; d++ {
		// for i := 0; i < 9; i++ {
		// 	// out = append(out, RowRule{Expected: d, Row: i})
		// 	// out = append(out, ColumnRule{Expected: d, Column: i})
		// }
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				out = append(out, BlockRule{C: i, R: j, Expected: d})
			}
		}
	}
	return
}

func (s Sudoku) String() string {
	return s.StringMark(-1, -1)
}

func (s Sudoku) StringMark(r, c int) string {
	var sb strings.Builder
	for i, row := range s.grid {
		if i > 0 && i%3 == 0 {
			if r > -1 {
				sb.WriteString("  ")
			}
			sb.WriteString("---+---+---\n")
		}
		if r > -1 {
			if i == r {
				sb.WriteString("> ")
			} else {
				sb.WriteString("  ")
			}
		}
		for j, d := range row {
			if j > 0 && j%3 == 0 {
				sb.WriteRune('|')
			}
			switch d {
			case 0:
				sb.WriteRune('.')
			default:
				sb.WriteString(strconv.Itoa(int(d)))
			}
		}
		sb.WriteRune('\n')
	}
	if c > -1 {
		if r > -1 {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.Repeat(" ", c+int(c/3)) + "^\n")
	}
	return sb.String()
}

type RowRule struct {
	// Expected Digit
	Row int
	// SkipColumns [9]bool
	// found bool
}

func (r RowRule) TryReduce(s Sudoku, rs []Ruler) (Sudoku, []Ruler, bool) {
	s, rs, changed := reduce(s, rs, cellRows[r.Row])
	if changed {
		fmt.Printf("reduced by %#v\n", r)
	}
	return s, rs, changed
}

func reduce(s Sudoku, rs []Ruler, cells Cells) (_ Sudoku, _ []Ruler, changed bool) {
	var debug bool
	if (cells[0] == Cell{R: 0, C: 1}) && (cells[1] == Cell{R: 1, C: 1}) {
		// debug = true
	}
	if debug {
		fmt.Printf("----------------------------\n")
		fmt.Printf("cells %v\n", cells)
	}
	var perDigits [MaxDigit][]int
	var perIndex [9][]Digit
	for i, cl := range cells {
		var row = cl.R
		var col = cl.C
		if s.grid[row][col] != 0 {
			continue
		}
		for d := Digit(1); d <= MaxDigit; d++ {
			test := testDigit(s, rs, row, col, d)
			if test {
				perDigits[d-1] = append(perDigits[d-1], i)
				perIndex[i] = append(perIndex[i], d)
			}
		}
	}

	if debug {
		fmt.Printf("perDigits %v\n", perDigits)
		fmt.Printf("perIndex %v\n", perIndex)
		fmt.Printf("----------------------------\n")
	}

	for i, idxes := range perDigits {
		d := Digit(i + 1)
		switch len(idxes) {
		case 0:
		case 1:
			row := cells[idxes[0]].R
			col := cells[idxes[0]].C
			// fmt.Printf("reduced %dx%d -> %d\n", row, col, d)
			s = s.Set(row, col, d)
			changed = true
			return s, rs, true
		case 2, 3:
			var dubls int = 1
			var d2 Digit
			for j := i + 1; j < len(perDigits); j++ {
				if colsEquals(perDigits[j], idxes) {
					dubls++
					switch dubls {
					case 2:
						d2 = Digit(j + 1)
					}
				}
			}

			if dubls == len(idxes) {
				var cs []Cell
				for _, idx := range idxes {
					cs = append(cs, cells[idx])
				}
				switch dubls {
				case 2:
					var ok bool
					rs, ok = appendUniq(rs, Pair{
						Cs: [2]Cell{
							cells[idxes[0]],
							cells[idxes[1]],
						},
						Ds: [2]Digit{d, d2},
					})
					if ok {
						changed = true
						fmt.Printf("pair reducer added %v -> %d,%d\n", cs, d, d2)
					}
				default:
					fmt.Printf("!!! Make pair per digits %v -> %d,%d\n", cs, d, d2)
				}
			}
		}
	}

	for idx, ds := range perIndex {
		switch len(ds) {
		case 1:
			row := cells[idx].R
			col := cells[idx].C
			d := ds[0]
			// fmt.Printf("reduced %dx%d -> %d\n", row, col, d)
			s = s.Set(row, col, d)
			changed = true
			return s, rs, true
			// case 2:
			// 	fmt.Printf("Make pair per columns %dx%d -> %v\n", row, col, ds)
		}
	}

	return s, rs, changed
}

func (r RowRule) TestDigit(s Sudoku, row, column int, d Digit) bool {
	if r.Row != row {
		return true
	}
	// if r.Expected != d {
	// 	return true
	// }
	if s.grid[row][column] != 0 {
		panic("unexpected")
	}
	// var exists bool
	for c := 0; c < 9; c++ {
		if c == column {
			continue
		}
		if s.grid[row][c] == d {
			return false
		}
	}
	return true
}

func testDigit(s Sudoku, rs []Ruler, x, y int, d Digit) bool {
	if s.grid[x][y] != 0 {
		return false
	}
	for _, rule := range rs {
		test := rule.TestDigit(s, x, y, d)
		if !test {
			return false
		}
	}
	return true
}

type ColumnRule struct {
	// Expected Digit
	Column int
	// SkipRows [9]bool
	// found bool
}

func (r ColumnRule) TryReduce(s Sudoku, rs []Ruler) (Sudoku, []Ruler, bool) {
	s, rs, changed := reduce(s, rs, cellCols[r.Column])
	if changed {
		fmt.Printf("reduced by %#v\n", r)
	}
	return s, rs, changed
}

func (r ColumnRule) TestDigit(s Sudoku, row, column int, d Digit) bool {
	if r.Column != column {
		return true
	}
	// if r.Expected != d {
	// 	return true
	// }
	if s.grid[row][column] != 0 {
		panic("unexpected")
	}
	// var exists bool
	for i := 0; i < 9; i++ {
		if i == row {
			continue
		}
		if s.grid[i][column] == d {
			return false
		}
	}
	return true
}

type BlockRule struct {
	C, R     int
	Expected Digit
}

func (r BlockRule) TestDigit(s Sudoku, row, column int, d Digit) bool {
	if r.C*3+2 < column || r.C*3 > column {
		return true
	}
	if r.R*3+2 < row || r.R*3 > row {
		return true
	}
	if r.Expected != d {
		return true
	}
	if s.grid[row][column] != 0 {
		panic("unexpected")
	}
	// var exists bool
	for i := r.R * 3; i < r.R*3+3; i++ {
		for j := r.C * 3; j < r.C*3+3; j++ {
			if i == row && j == column {
				continue
			}
			if s.grid[i][j] == d {
				return false
			}
		}
	}
	return true
}

func (r BlockRule) TryReduce(s Sudoku, rs []Ruler) (Sudoku, []Ruler, bool) {
	// if r.found {
	// 	return s, false
	// }
	var found bool
	var foundColumn int
	var foundRow int

	for i := r.R * 3; i < r.R*3+3; i++ {
		for j := r.C * 3; j < r.C*3+3; j++ {
			if s.grid[i][j] == r.Expected {
				return s, rs, false
			}
			if s.grid[i][j] != 0 {
				continue
			}
			ok := testDigit(s, rs, i, j, r.Expected)
			// if r.R == 0 && r.C == 2 {
			// 	fmt.Printf("%dx%d -> %d: %t\n", i, j, r.Expected, ok)
			// }
			if ok {
				if found {
					return s, rs, false
				}
				found = true
				foundRow = i
				foundColumn = j
			}
		}
	}
	if found {
		// r.found = true
		fmt.Printf("reduced by %#v\n", r)
		return s.Set(foundRow, foundColumn, r.Expected), rs, true
	}
	panic("unexpected")
}

type Ruler interface {
	TryReduce(s Sudoku, rs []Ruler) (Sudoku, []Ruler, bool)
	TestDigit(s Sudoku, row, column int, d Digit) bool
}

type Pair struct {
	Cs [2]Cell
	Ds [2]Digit
}

var _ Ruler = Pair{}

func (p Pair) TryReduce(s Sudoku, rs []Ruler) (Sudoku, []Ruler, bool) {
	for _, cell := range p.Cs {
		var c int
		var cd Digit
		for _, d := range p.Ds {
			ok := testDigit(s, rs, cell.R, cell.C, d)
			if ok {
				c++
				cd = d
			}
		}
		if c == 1 {
			return s.Set(cell.R, cell.C, cd), rs, true
		}
	}
	return s, rs, false
}

func (p Pair) TestDigit(s Sudoku, row, column int, d Digit) bool {
	var idx int = -1
	for i, c := range p.Cs {
		if c.R == row && c.C == column {
			idx = i
		}
	}
	if idx == -1 {
		return true
	}
	var ok bool
	for _, digit := range p.Ds {
		if digit == d {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}

	for _, c := range p.Cs {
		if s.grid[c.R][c.C] == d {
			return false
		}
	}

	return true
}

type Cell struct {
	R, C int
}

type Cells []Cell

var (
	cellRows   [9][]Cell
	cellCols   [9][]Cell
	cellBlocks [3][3][]Cell
)

func init() {
	for i := 0; i < 9; i++ {
		cellRows[i] = rowCells(i)
		cellCols[i] = colCells(i)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cellBlocks[i][j] = blockCells(i, j)
		}
	}
}

func rowCells(row int) []Cell {
	var out = make([]Cell, 9)
	for i := 0; i < 9; i++ {
		out[i].R = row
		out[i].C = i
	}
	return out
}

func colCells(col int) []Cell {
	var out = make([]Cell, 9)
	for i := 0; i < 9; i++ {
		out[i].R = i
		out[i].C = col
	}
	return out
}

func blockCells(r, c int) []Cell {
	var out = make([]Cell, 0, 9)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			out = append(out, Cell{R: r*3 + i, C: c*3 + j})
		}
	}
	return out
}
