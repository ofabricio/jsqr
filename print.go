package jsqr

import (
	"fmt"
	"io"
	"os"
)

func print(a ast) {
	writeast(os.Stdout, a)
}

func writeast(w io.Writer, a ast) {
	write(w, a, 0)
}

func write(w io.Writer, a ast, depth int) {
	for i := 0; i < depth; i++ {
		fmt.Fprint(w, "    ")
	}
	switch t := a.(type) {
	case astOr:
		fmt.Fprintln(w, "[Or]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astAnd:
		fmt.Fprintln(w, "[And]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astEQ:
		fmt.Fprintln(w, "[EQ]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astEQI:
		fmt.Fprintln(w, "[EQI]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astNE:
		fmt.Fprintln(w, "[NE]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astNEI:
		fmt.Fprintln(w, "[NEI]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astGTE:
		fmt.Fprintln(w, "[GTE]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astGT:
		fmt.Fprintln(w, "[GT]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astLTE:
		fmt.Fprintln(w, "[LTE]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astLT:
		fmt.Fprintln(w, "[LT]")
		write(w, t.A, depth+1)
		write(w, t.B, depth+1)
	case astPath:
		fmt.Fprintln(w, "[Path]")
		for _, v := range t.V {
			write(w, v, depth+1)
		}
	case astNumFloat:
		fmt.Fprintf(w, "[Num] %v\n", t.V)
	case astKey:
		fmt.Fprintf(w, "[Key]\n")
		write(w, t.V, depth+1)
	case astIdent:
		fmt.Fprintf(w, "[Idn] %s\n", t.V)
	case astStr:
		fmt.Fprintf(w, "[Str] %s\n", t.V)
	case astBool:
		fmt.Fprintf(w, "[Bool] %t\n", t.V)
	case astNil:
		fmt.Fprintln(w, "[Null]")
	case astThis:
		fmt.Fprintln(w, "[This]")
	case astArr:
		fmt.Fprintln(w, "[Array]")
		write(w, t.V, depth+1)
	case astArrIdx:
		fmt.Fprintf(w, "[Index] %v\n", t.V)
	case astUpper:
		fmt.Fprintln(w, "[Fun] upper")
	case astLower:
		fmt.Fprintln(w, "[Fun] lower")
	case astExists:
		fmt.Fprintln(w, "[Fun] exists")
	}
}
