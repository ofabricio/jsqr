package jsqr

// Get applies expr to jsn and returns the resulting JSON.
// The JSON must be valid.
func Get(jsn []byte, expr string) Json {
	exp := Compile(expr)
	return exp.Parse(jsn)
}

// GetStruct applies expr to v and returns the resulting Go value.
func GetStruct(v any, expr string) any {
	exp := Compile(expr)
	return exp.ParseStruct(v)
}

// Compile parses an expression and returns, if successful,
// an [Expr] object that can be used to match against JSON.
func Compile(expr string) Expr {
	c := compiler{}
	return c.Compile(expr)
}

type Expr struct {
	expr ast
}

func (e *Expr) Parse(jsn []byte) Json {
	p := parser{}
	return p.parse(e.expr, jsn)
}

func (e *Expr) ParseStruct(a any) any {
	p := parserStruct{}
	return p.parse(e.expr, a)
}

type ast interface{}
type astOr struct{ A, B ast }
type astAnd struct{ A, B ast }
type astEQ struct{ A, B ast }
type astEQI struct{ A, B ast }
type astNE struct{ A, B ast }
type astNEI struct{ A, B ast }
type astGTE struct{ A, B ast }
type astGT struct{ A, B ast }
type astLTE struct{ A, B ast }
type astLT struct{ A, B ast }
type astPath struct{ V []ast }
type astNumFloat struct{ V float64 }
type astKey struct{ V ast }
type astIdent struct{ V string }
type astStr struct{ V string }
type astBool struct{ V bool }
type astNil struct{}
type astThis struct{}
type astArr struct{ V ast }
type astArrIdx struct{ V int }
type astUpper struct{}
type astLower struct{}
type astExists struct{}
