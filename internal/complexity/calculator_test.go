package complexity

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/MikeMwita/go-strict/config"
)

// fixture parses src and returns the Calculator and the parsed *ast.File.
func fixture(t *testing.T, src string) (*Calculator, *ast.File, *token.FileSet) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return NewCalculator(fset, f, config.DefaultConfig()), f, fset
}

// firstFunc returns the first *ast.FuncDecl in f.
func firstFunc(f *ast.File) *ast.FuncDecl {
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			return fd
		}
	}
	return nil
}

func TestSimpleFunction(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Add(x, y int) int { return x + y }
`)
	r := calc.Calculate(firstFunc(f))
	// base=1, return at nesting=0 does not count → complexity=1
	if r.Complexity != 1 {
		t.Errorf("Add: want 1, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestIfStatement(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo(x int) int {
	if x > 0 {
		return x
	}
	return 0
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// if(nesting=0)=+1 → 2
	// return inside if(nesting=1)=+2 → 4
	if r.Complexity != 4 {
		t.Errorf("if stmt: want 4, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestNestedIfInFor(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo(xs []int) {
	for _, x := range xs {
		if x > 0 {
			return
		}
	}
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// range(nesting=0)=+1 → 2
	// if(nesting=1)=+2   → 4
	// return(nesting=2)=+3 → 7
	if r.Complexity != 7 {
		t.Errorf("nested if-in-for: want 7, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestSwitchWithCases(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo(x int) string {
	switch x {
	case 1:
		return "one"
	case 2:
		return "two"
	default:
		return "other"
	}
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// switch(nesting=0)=+1 → 2
	// case(flat)=+1 → 3; return(nesting=1)=+2 → 5
	// case(flat)=+1 → 6; return(nesting=1)=+2 → 8
	// default(flat)=+1 → 9; return(nesting=1)=+2 → 11
	if r.Complexity != 11 {
		t.Errorf("switch 3 cases: want 11, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestLogicalOperators(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo(a, b, c bool) bool {
	return a && b && c
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1, one && sequence = +1 → 2
	if r.Complexity != 2 {
		t.Errorf("single && sequence: want 2, got %d", r.Complexity)
	}
}

func TestMixedLogicalOperators(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo(a, b, c bool) bool {
	return a && b || c
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1, two sequences (&& then ||) = +2 → 3
	if r.Complexity != 3 {
		t.Errorf("mixed &&/||: want 3, got %d", r.Complexity)
	}
}

func TestClosure(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo() {
	fn := func() {
		if true {
		}
	}
	_ = fn
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// closure(nesting=0)=+1 → 2
	// if inside closure(nesting=1)=+2 → 4
	if r.Complexity != 4 {
		t.Errorf("closure: want 4, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestElseIfChain(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Grade(score int) string {
	if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	} else {
		return "F"
	}
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// missing comment (10 lines)=+1→2
	// if(n=0)=+1→3; return(n=1)=+2→5
	// else(flat)=+1→6; return inside else-if body(n=1)=+2→8
	// else(flat)=+1→9; return inside else-if body(n=1)=+2→11
	// else(flat)=+1→12; return inside else body(n=1)=+2→14
	if r.Complexity != 14 {
		t.Errorf("else-if chain: want 14, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestBranchInsideSwitch(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Foo(x int) {
	switch x {
	case 1:
		break
	}
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// switch(n=0)=+1→2
	// case(flat)=+1→3
	// branch(n=1)=+2→5
	if r.Complexity != 5 {
		t.Errorf("branch in switch: want 5, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestSelectStatement(t *testing.T) {
	calc, f, _ := fixture(t, `package p
import "time"
func Foo(ch chan int) int {
	select {
	case v := <-ch:
		return v
	case <-time.After(0):
		return 0
	}
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// select(n=0)=+1→2
	// case(flat)=+1→3; return(n=1)=+2→5
	// case(flat)=+1→6; return(n=1)=+2→8
	if r.Complexity != 8 {
		t.Errorf("select: want 8, got %d (increments=%v)", r.Complexity, r.Increments)
	}
}

func TestExceededFlag(t *testing.T) {
	calc, f, _ := fixture(t, `package p
func Heavy(a, b, c, d bool) bool {
	if a {
		if b {
			if c {
				if d {
					return true
				}
			}
		}
	}
	return false
}
`)
	r := calc.Calculate(firstFunc(f))
	// base=1
	// missing comment (11 lines)=+1→2
	// if(n=0)=+1→3
	// if(n=1)=+2→5
	// if(n=2)=+3→8
	// if(n=3)=+4→12
	// return(n=4)=+5→17
	if r.Complexity != 17 {
		t.Errorf("deeply nested: want 17, got %d (increments=%v)", r.Complexity, r.Increments)
	}
	cfg := calc.cfg
	if r.Complexity > cfg.Threshold && !r.Exceeded {
		t.Error("Exceeded flag should be true")
	}
}
