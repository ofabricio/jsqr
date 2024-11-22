package jsqr

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/ofabricio/scan"
)

type compiler struct {
	s scan.Bytes
}

func (c *compiler) Compile(expr string) Expr {
	c.s = scan.Bytes(expr)
	var out ast
	if c.expr(&out) {
		return Expr{expr: out}
	}
	return Expr{}
}

func (c *compiler) expr(out *ast) bool {
	return c.or(out)
}

func (c *compiler) or(out *ast) bool {
	var a, b ast
	if c.and(&a) {
		if c.s.Match("|") && c.or(&b) {
			*out = astOr{a, b}
			return true
		}
		*out = a
		return true
	}
	return false
}

func (c *compiler) and(out *ast) bool {
	var a, b ast
	if c.cnd(&a) {
		if c.s.Match("&") && c.and(&b) {
			*out = astAnd{a, b}
			return true
		}
		*out = a
		return true
	}
	return false
}

func (c *compiler) cnd(out *ast) bool {
	var a, b ast
	if c.val(&a) {
		c.s.Spaces()
		if op := c.cndop(); op != "" && c.cnd(&b) {
			switch op {
			case "==":
				*out = astEQ{a, b}
			case "!=":
				*out = astNE{a, b}
			case ">=":
				*out = astGTE{a, b}
			case ">":
				*out = astGT{a, b}
			case "<=":
				*out = astLTE{a, b}
			case "<":
				*out = astLT{a, b}
			case "eq":
				*out = astEQI{a, b}
			case "ne":
				*out = astNEI{a, b}
			default:
				return false
			}
			return true
		}
		*out = a
		return true
	}
	return false
}

func (c *compiler) val(out *ast) bool {
	c.s.Spaces()
	if c.s.Match("(") && c.expr(out) && c.s.Match(")") {
		return true
	}
	if c.path(out) {
		return true
	}
	if c.str(out) {
		return true
	}
	if c.num(out) {
		return true
	}
	if c.bool(out) {
		return true
	}
	if c.null(out) {
		return true
	}
	return false
}

func (c *compiler) path(out *ast) bool {
	var vs []ast
	var v ast
	for c.key(&v) {
		vs = append(vs, v)
	}
	if len(vs) > 0 {
		*out = astPath{vs}
		return true
	}
	return false
}

func (c *compiler) key(out *ast) bool {
	if c.s.Match(".") {
		if c.str(out) {
			*out = astKey{*out}
			return true
		}
		if c.idn(out) {
			return true
		}
		if c.arr(out) {
			return true
		}
		if c.fun(out) {
			return true
		}
		*out = astThis{}
		return true
	}
	return false
}

func (c *compiler) arr(out *ast) bool {
	var a ast
	if c.s.Match("[") && c.expr(&a) && c.s.Match("]") {
		if n, ok := a.(astNumFloat); ok {
			*out = astArrIdx{int(n.V)}
		} else {
			*out = astArr{a}
		}
		return true
	}
	return false
}

func (c *compiler) fun(out *ast) bool {
	var a ast
	if c.s.Match("(") && c.idn(&a) && c.s.Match(")") {
		switch a.(astIdent).V {
		case "upper":
			*out = astUpper{}
		case "lower":
			*out = astLower{}
		case "exists":
			*out = astExists{}
		default:
			return false
		}
		return true
	}
	return false
}

func (c *compiler) cndop() string {
	if m := c.s.Mark(); c.s.Match("==") ||
		c.s.Match("!=") ||
		c.s.Match(">=") ||
		c.s.Match(">") ||
		c.s.Match("<=") ||
		c.s.Match("<") ||
		c.s.Match("eq") ||
		c.s.Match("ne") {
		return c.s.Text(m)
	}
	return ""
}

func (c *compiler) num(out *ast) bool {
	if m := c.s.Mark(); c.s.MatchNumber() {
		f, _ := strconv.ParseFloat(c.s.Delta(m).String(), 64) // TODO: handle error.
		*out = astNumFloat{f}
		return true
	}
	return false
}

func (c *compiler) str(out *ast) bool {
	if m := c.s.Mark(); c.s.MatchString(`"`) {
		v := strings.ReplaceAll(c.s.Delta(m).String(), `\"`, `"`)
		*out = astStr{v}
		return true
	}
	return false
}

func (c *compiler) idn(out *ast) bool {
	m := c.s.Mark()
	for c.s.MatchFunc(unicode.IsDigit) ||
		c.s.MatchFunc(unicode.IsLetter) ||
		c.s.Match("_") {
	}
	if t := c.s.Text(m); t != "" {
		*out = astIdent{c.s.Text(m)}
		return true
	}
	return false
}

func (c *compiler) bool(out *ast) bool {
	if c.s.Match("true") {
		*out = astBool{true}
		return true
	}
	if c.s.Match("false") {
		*out = astBool{false}
		return true
	}
	return false
}

func (c *compiler) null(out *ast) bool {
	if c.s.Match("null") {
		*out = astNil{}
		return true
	}
	return false
}
