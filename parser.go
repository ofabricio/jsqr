package jsqr

import (
	"strconv"
	"strings"
)

type parser struct{}

func (p *parser) parse(expr ast, jsn []byte) Json {
	j := Json{jsn}
	j.Spaces()
	return p.parseExpr(expr, j)
}

func (p *parser) parseExpr(expr ast, j Json) Json {
	switch expr := expr.(type) {
	case astOr:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.IsTrue() || b.IsTrue())
	case astAnd:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.IsTrue() && b.IsTrue())
	case astEQ:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.EQ(b))
	case astNE:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.NE(b))
	case astEQI:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.EQI(b))
	case astNEI:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.NEI(b))
	case astGTE:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.GTE(b))
	case astGT:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.GT(b))
	case astLTE:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.LTE(b))
	case astLT:
		a := p.parseExpr(expr.A, j)
		b := p.parseExpr(expr.B, j)
		return trueOrFalse(a.LT(b))
	case astPath:
		for _, n := range expr.V {
			j = p.parseExpr(n, j)
		}
		return j
	case astArr:
		for _, item := range j.IterArray() {
			if j := p.parseExpr(expr.V, item); j.IsTrue() {
				return item
			}
		}
	case astArrIdx:
		return j.GetIndex(expr.V)
	case astKey:
		return j.GetKey(p.parseExpr(expr.V, j).String())
	case astIdent:
		return j.GetKey(expr.V)
	case astStr:
		return Json{[]byte(expr.V)}
	case astNumFloat:
		f := strconv.FormatFloat(expr.V, 'f', -1, 64)
		return Json{[]byte(f)}
	case astBool:
		return trueOrFalse(expr.V)
	case astNil:
		return Json{[]byte("null")}
	case astThis:
		return j
	case astUpper:
		return Json{[]byte(strings.ToUpper(j.String()))}
	case astLower:
		return Json{[]byte(strings.ToLower(j.String()))}
	case astExists:
		return trueOrFalse(j.More())
	}
	return Json{}
}

func trueOrFalse(v bool) Json {
	if v {
		return Json{[]byte("true")}
	}
	return Json{[]byte("false")}
}
