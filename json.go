package jsqr

import (
	"iter"
	"strconv"
	"strings"

	"github.com/ofabricio/scan"
)

type Json struct{ scan.Bytes }

func (j Json) GetKey(key string) Json {
	if len(key) > 0 && key[0] == '"' {
		key = key[1 : len(key)-1]
	}
	for k, v := range j.IterObject() {
		if k[1:len(k)-1] == key {
			return v
		}
	}
	return Json{}
}

func (j Json) GetIndex(idx int) Json {
	for i, v := range j.IterArray() {
		if i == idx {
			return v
		}
	}
	return Json{}
}

func (j *Json) IterArray() iter.Seq2[int, Json] {
	return func(yield func(int, Json) bool) {
		if j.MatchChar("[") {
			j.Spaces()
			for i := 0; !j.MatchChar("]"); i++ {
				if !yield(i, j.value()) {
					return
				}
				if j.Spaces(); !j.MatchChar(",") {
					break
				}
				j.Spaces()
			}
		}
	}
}

func (j *Json) IterObject() iter.Seq2[string, Json] {
	return func(yield func(string, Json) bool) {
		if j.MatchChar("{") {
			for j.Spaces(); !j.MatchChar("}"); j.Spaces() {

				m := j.Mark()
				j.MatchString(`"`)
				k := j.Text(m)

				_, _, _ = j.Spaces(), j.MatchChar(":"), j.Spaces()

				if !yield(k, j.value()) {
					return
				}

				if j.Spaces(); !j.MatchChar(",") {
					break
				}
			}
		}
	}
}

func (a Json) EQ(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() == b.Float64()
	}
	return a.String() == b.String()
}

func (a Json) NE(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() != b.Float64()
	}
	return a.String() != b.String()
}

func (a Json) EQI(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() == b.Float64()
	}
	return strings.EqualFold(a.String(), b.String())
}

func (a Json) NEI(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() != b.Float64()
	}
	return !strings.EqualFold(a.String(), b.String())
}

func (a Json) GTE(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() >= b.Float64()
	}
	return a.String() >= b.String()
}

func (a Json) GT(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() > b.Float64()
	}
	return a.String() > b.String()
}

func (a Json) LTE(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() <= b.Float64()
	}
	return a.String() <= b.String()
}

func (a Json) LT(b Json) bool {
	if a.IsNumber() && b.IsNumber() {
		return a.Float64() < b.Float64()
	}
	return a.String() < b.String()
}

func (a Json) IsNumber() bool {
	b := a.Byte()
	return b == '-' || b >= '0' && b <= '9'
}

func (a Json) IsTrue() bool {
	return a.Byte() == 't'
}

func (a Json) Float64() float64 {
	f, _ := strconv.ParseFloat(a.String(), 64)
	return f
}

func (j *Json) value() Json {
	m := j.Mark()
	_ = j.MatchString(`"`) ||
		j.matchOpenCloseCount('{', '}') ||
		j.matchOpenCloseCount('[', ']') ||
		j.Match("true") || j.Match("false") ||
		j.Match("null") ||
		j.MatchNumber()
	return Json{j.Delta(m)}
}

// matchOpenCloseCount matches open and close by counting them.
func (s *Json) matchOpenCloseCount(open, clos byte) bool {
	if ss := s.Mark(); len(ss) > 0 && ss[0] == open {
		c := 0
		for i := 0; i < len(ss); i++ {
			if ss[i] == open {
				c++
				continue
			}
			if ss[i] == clos {
				if c--; c == 0 {
					s.Move(ss[i+1:])
					break
				}
			}
			// Skip string.
			if ss[i] == '"' {
				for i = i + 1; i < len(ss); i++ {
					if ss[i] == '"' && ss[i-1] != '\\' {
						break
					}
				}
				continue
			}

		}
		return c == 0
	}
	return false
}
