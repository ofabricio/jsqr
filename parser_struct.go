package jsqr

import (
	"reflect"
	"strings"
)

type parserStruct struct{}

func (p *parserStruct) parse(expr ast, v any) any {
	vOf := reflect.ValueOf(v)
	if v := p.parseExpr(expr, vOf); v.IsValid() {
		return v.Interface()
	}
	return nil
}

func (p *parserStruct) parseExpr(expr ast, v reflect.Value) reflect.Value {
	v = reflect.Indirect(v)
	switch expr := expr.(type) {
	case astOr:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(a.Bool() || b.Bool())
	case astAnd:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(a.Bool() && b.Bool())
	case astEQ:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(p.eq(a, b))
	case astNE:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(!p.eq(a, b))
	case astEQI:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(p.eqi(a, b))
	case astNEI:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(!p.eqi(a, b))
	case astGTE:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(p.gte(a, b))
	case astGT:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(p.gt(a, b))
	case astLTE:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(!p.gt(a, b))
	case astLT:
		a := p.parseExpr(expr.A, v)
		b := p.parseExpr(expr.B, v)
		return reflect.ValueOf(!p.gte(a, b))
	case astPath:
		for _, n := range expr.V {
			v = p.parseExpr(n, v)
		}
		return v
	case astArr:
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			for i := range v.Len() {
				item := v.Index(i)
				v := p.parseExpr(expr.V, item)
				if v.IsValid() && v.Kind() == reflect.Bool && v.Interface() == true {
					return item
				}
			}
		}
	case astArrIdx:
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			if i := expr.V; i < v.Len() {
				return v.Index(i)
			}
		}
	case astKey:
		return p.parseExpr(expr.V, v)
	case astIdent:
		if v.Kind() == reflect.Struct {
			if f := v.FieldByName(expr.V); f.IsValid() {
				return f
			}
		}
	case astStr:
		return reflect.ValueOf(expr.V[1 : len(expr.V)-1])
	case astNumFloat:
		return reflect.ValueOf(expr.V)
	case astBool:
		return reflect.ValueOf(expr.V)
	case astNil:
		type customNil int
		return reflect.ValueOf((*customNil)(nil))
	case astThis:
		return v
	case astUpper:
		return reflect.ValueOf(strings.ToUpper(v.String()))
	case astLower:
		return reflect.ValueOf(strings.ToLower(v.String()))
	case astExists:
		return reflect.ValueOf(v.IsValid())
	}
	return reflect.ValueOf(nil)
}

func (p *parserStruct) eq(a, b reflect.Value) bool {
	if a.IsValid() && b.IsValid() && a.Comparable() && b.Comparable() {
		if a.CanConvert(reflect.TypeFor[float64]()) && b.CanConvert(reflect.TypeFor[float64]()) {
			a = a.Convert(reflect.TypeFor[float64]())
			b = b.Convert(reflect.TypeFor[float64]())
			return a.Float() == b.Float()
		}
		if a.Kind() == reflect.Ptr && b.Kind() == reflect.Ptr && a.IsNil() && b.IsNil() {
			return true
		}
		return a.Interface() == b.Interface()
	}
	return false
}

func (p *parserStruct) gte(a, b reflect.Value) bool {
	if a.IsValid() && b.IsValid() && a.Comparable() && b.Comparable() {
		if a.CanConvert(reflect.TypeFor[float64]()) && b.CanConvert(reflect.TypeFor[float64]()) {
			a = a.Convert(reflect.TypeFor[float64]())
			b = b.Convert(reflect.TypeFor[float64]())
			return a.Float() >= b.Float()
		}
		if a.CanConvert(reflect.TypeFor[string]()) && b.CanConvert(reflect.TypeFor[string]()) {
			a = a.Convert(reflect.TypeFor[string]())
			b = b.Convert(reflect.TypeFor[string]())
			return a.String() >= b.String()
		}
	}
	return false
}

func (p *parserStruct) gt(a, b reflect.Value) bool {
	if a.IsValid() && b.IsValid() && a.Comparable() && b.Comparable() {
		if a.CanConvert(reflect.TypeFor[float64]()) && b.CanConvert(reflect.TypeFor[float64]()) {
			a = a.Convert(reflect.TypeFor[float64]())
			b = b.Convert(reflect.TypeFor[float64]())
			return a.Float() > b.Float()
		}
		if a.CanConvert(reflect.TypeFor[string]()) && b.CanConvert(reflect.TypeFor[string]()) {
			a = a.Convert(reflect.TypeFor[string]())
			b = b.Convert(reflect.TypeFor[string]())
			return a.String() > b.String()
		}
	}
	return false
}

func (p *parserStruct) eqi(a, b reflect.Value) bool {
	if a.IsValid() && b.IsValid() && a.Comparable() && b.Comparable() {
		if a.CanConvert(reflect.TypeFor[string]()) && b.CanConvert(reflect.TypeFor[string]()) {
			a = a.Convert(reflect.TypeFor[string]())
			b = b.Convert(reflect.TypeFor[string]())
			return strings.EqualFold(a.String(), b.String())
		}
	}
	return false
}
