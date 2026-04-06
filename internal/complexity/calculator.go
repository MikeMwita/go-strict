// Package complexity implements a cognitive complexity calculator for Go source
// code. The algorithm is based on the SonarSource Cognitive Complexity
// whitepaper adapted for Go-specific constructs.
//
// Core scoring rules:
//   - Control-flow structures (if, for, range, switch, select, defer, go,
//     BranchStmt, ReturnStmt inside a branch): +1 + nesting_level
//   - Continuation clauses (else, case): +1 flat
//   - if/else blocks whose body spans >= blockLengthThreshold lines: +1 extra
//   - Logical operator sequences (contiguous &&/||): +1 per sequence
//   - Function-level rules (missing comment, parameter count, function length,
//     generic type params): +1 each (flat)
//
// Nesting increments: nesting increases by 1 when entering the body of any
// control-flow structure, else clause, case body, or function literal. The
// nesting value is passed by value so sibling statements share the same depth.
package complexity

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/MikeMwita/go-strict/models"
)

// Calculator computes the cognitive complexity of a single Go function.
type Calculator struct {
	fset     *token.FileSet
	file     *ast.File
	cfg      *models.LintConfig
	funcName string
}

// NewCalculator creates a Calculator for the given file set, source file, and
// lint configuration.
func NewCalculator(fset *token.FileSet, file *ast.File, cfg *models.LintConfig) *Calculator {
	return &Calculator{fset: fset, file: file, cfg: cfg}
}

// Calculate returns the full FunctionResult for funcDecl.
func (c *Calculator) Calculate(funcDecl *ast.FuncDecl) models.FunctionResult {
	if funcDecl == nil || funcDecl.Body == nil {
		return models.FunctionResult{}
	}
	c.funcName = funcDecl.Name.Name

	pos := c.fset.Position(funcDecl.Pos())
	result := models.FunctionResult{
		File:     pos.Filename,
		Line:     pos.Line,
		Col:      pos.Column,
		Function: c.funcName,
	}

	// Base complexity of 1 for the function itself.
	running := 1
	result.Increments = append(result.Increments, models.Increment{
		Kind:         "function",
		Line:         pos.Line,
		Col:          pos.Column,
		Score:        1,
		Nesting:      0,
		RunningTotal: running,
	})

	// Function-level rules applied before walking the body.
	running = c.applyFunctionRules(funcDecl, running, &result.Increments)

	// Walk the function body with nesting=0.
	running = c.walkBody(funcDecl.Body, 0, running, &result.Increments)

	result.Complexity = running
	result.Exceeded = running > c.cfg.Threshold
	return result
}

// applyFunctionRules checks rules that apply to the function as a whole.
func (c *Calculator) applyFunctionRules(funcDecl *ast.FuncDecl, running int, incs *[]models.Increment) int {
	startLine := c.fset.Position(funcDecl.Pos()).Line
	endLine := c.fset.Position(funcDecl.End()).Line
	funcLines := endLine - startLine

	// missing_comment rule
	if r := c.cfg.Rule("missing_comment"); r.Enabled {
		minLines := r.MinLines
		if minLines == 0 {
			minLines = 10
		}
		if funcLines >= minLines && !hasCorrectComment(funcDecl) {
			running++
			*incs = append(*incs, models.Increment{
				Kind:         "missing or wrong comment for function with more than 10 lines",
				Line:         startLine,
				Score:        1,
				Nesting:      0,
				RunningTotal: running,
			})
		}
	}

	// parameter_count rule
	if r := c.cfg.Rule("parameter_count"); r.Enabled {
		maxParams := r.MaxParams
		if maxParams == 0 {
			maxParams = 5
		}
		paramCount := countParams(funcDecl)
		if paramCount > maxParams {
			w := scoreWeight(r.Weight)
			running += w
			*incs = append(*incs, models.Increment{
				Kind:         "too many parameters",
				Line:         startLine,
				Score:        w,
				Nesting:      0,
				RunningTotal: running,
			})
		}
	}

	// function_length rule
	if r := c.cfg.Rule("function_length"); r.Enabled {
		maxLines := r.MaxLines
		if maxLines == 0 {
			maxLines = 50
		}
		if funcLines > maxLines {
			w := scoreWeight(r.Weight)
			running += w
			*incs = append(*incs, models.Increment{
				Kind:         "function too long",
				Line:         startLine,
				Score:        w,
				Nesting:      0,
				RunningTotal: running,
			})
		}
	}

	// generics rule
	if r := c.cfg.Rule("generics"); r.Enabled {
		if funcDecl.Type.TypeParams != nil && len(funcDecl.Type.TypeParams.List) > 0 {
			w := scoreWeight(r.Weight)
			running += w
			*incs = append(*incs, models.Increment{
				Kind:         "generic type parameters",
				Line:         startLine,
				Score:        w,
				Nesting:      0,
				RunningTotal: running,
			})
		}
	}

	return running
}

// walkBody iterates over statements in a block, threading nesting by value.
func (c *Calculator) walkBody(body *ast.BlockStmt, nesting, running int, incs *[]models.Increment) int {
	if body == nil {
		return running
	}
	for _, stmt := range body.List {
		running = c.walkStmt(stmt, nesting, running, incs)
	}
	return running
}

// walkStmt is the central dispatcher. nesting is passed by value so each
// recursive call gets an independent copy — this is the key correctness
// invariant that the original implementation violated.
func (c *Calculator) walkStmt(stmt ast.Stmt, nesting, running int, incs *[]models.Increment) int {
	if stmt == nil {
		return running
	}

	switch s := stmt.(type) {

	case *ast.IfStmt:
		running = c.walkIfStmt(s, nesting, running, incs)

	case *ast.ForStmt:
		score := 1 + nesting
		line := c.fset.Position(s.Pos()).Line
		running += score
		*incs = append(*incs, models.Increment{
			Kind: "for", Line: line, Score: score,
			Nesting: nesting, RunningTotal: running,
		})
		running = c.walkBody(s.Body, nesting+1, running, incs)

	case *ast.RangeStmt:
		score := 1 + nesting
		line := c.fset.Position(s.Pos()).Line
		running += score
		*incs = append(*incs, models.Increment{
			Kind: "range", Line: line, Score: score,
			Nesting: nesting, RunningTotal: running,
		})
		running = c.walkBody(s.Body, nesting+1, running, incs)

	case *ast.SwitchStmt:
		running = c.walkSwitchStmt(s, nesting, running, incs)

	case *ast.TypeSwitchStmt:
		running = c.walkTypeSwitchStmt(s, nesting, running, incs)

	case *ast.SelectStmt:
		running = c.walkSelectStmt(s, nesting, running, incs)

	case *ast.BranchStmt:
		score := 1 + nesting
		line := c.fset.Position(s.Pos()).Line
		running += score
		*incs = append(*incs, models.Increment{
			Kind: "branch", Line: line, Score: score,
			Nesting: nesting, RunningTotal: running,
		})

	case *ast.ReturnStmt:
		// Return inside a control structure counts as a branch.
		if nesting > 0 {
			score := 1 + nesting
			line := c.fset.Position(s.Pos()).Line
			running += score
			*incs = append(*incs, models.Increment{
				Kind: "branch", Line: line, Score: score,
				Nesting: nesting, RunningTotal: running,
			})
		}
		// Walk return expressions for embedded closures/logical ops.
		for _, expr := range s.Results {
			running = c.walkExpr(expr, nesting, running, incs)
		}

	case *ast.GoStmt:
		score := 1 + nesting
		line := c.fset.Position(s.Pos()).Line
		running += score
		*incs = append(*incs, models.Increment{
			Kind: "go", Line: line, Score: score,
			Nesting: nesting, RunningTotal: running,
		})
		running = c.walkCallExpr(s.Call, nesting+1, running, incs)

	case *ast.DeferStmt:
		score := 1 + nesting
		line := c.fset.Position(s.Pos()).Line
		running += score
		*incs = append(*incs, models.Increment{
			Kind: "defer", Line: line, Score: score,
			Nesting: nesting, RunningTotal: running,
		})
		running = c.walkCallExpr(s.Call, nesting+1, running, incs)

	case *ast.LabeledStmt:
		line := c.fset.Position(s.Pos()).Line
		running++
		*incs = append(*incs, models.Increment{
			Kind: "label", Line: line, Score: 1,
			Nesting: 0, RunningTotal: running,
		})
		if s.Stmt != nil {
			running = c.walkStmt(s.Stmt, nesting, running, incs)
		}

	case *ast.AssignStmt:
		for _, rhs := range s.Rhs {
			running = c.walkExpr(rhs, nesting, running, incs)
		}

	case *ast.ExprStmt:
		running = c.walkExpr(s.X, nesting, running, incs)

	case *ast.SendStmt:
		running = c.walkExpr(s.Value, nesting, running, incs)

	case *ast.BlockStmt:
		running = c.walkBody(s, nesting, running, incs)

	case *ast.IncDecStmt:
		running = c.walkExpr(s.X, nesting, running, incs)

	case *ast.DeclStmt:
		running = c.walkDecl(s.Decl, nesting, running, incs)
	}

	return running
}

// walkIfStmt handles if / else if / else chains with correct nesting.
// Per the cognitive complexity spec, else/else-if are continuations and do NOT
// increase the nesting level for their own scoring — only the code inside
// each branch runs at nesting+1.
func (c *Calculator) walkIfStmt(s *ast.IfStmt, nesting, running int, incs *[]models.Increment) int {
	score := 1 + nesting
	line := c.fset.Position(s.Pos()).Line
	running += score
	*incs = append(*incs, models.Increment{
		Kind: "if", Line: line, Score: score,
		Nesting: nesting, RunningTotal: running,
	})

	// Block-length bonus for the if body.
	running = c.maybeBlockLengthBonus("if", s.Body, line, running, incs)

	// Walk the if-body at nesting+1.
	running = c.walkBody(s.Body, nesting+1, running, incs)

	// Handle the else clause as a flat continuation chain.
	if s.Else != nil {
		running = c.walkElseClause(s.Else, nesting, running, incs)
	}

	return running
}

// walkElseClause handles the else / else-if continuation of an if statement.
// It emits +1 flat for each else/else-if and walks their bodies at nesting+1.
// Critically, else-if does NOT recursively call walkIfStmt (which would add
// another +1+nesting) — only the body is walked at the original nesting+1.
func (c *Calculator) walkElseClause(else_ ast.Stmt, nesting, running int, incs *[]models.Increment) int {
	elseLine := c.fset.Position(else_.Pos()).Line

	switch e := else_.(type) {
	case *ast.IfStmt:
		// else if: +1 flat (continuation, no nesting penalty).
		running++
		*incs = append(*incs, models.Increment{
			Kind: "else", Line: elseLine, Score: 1,
			Nesting: 0, RunningTotal: running,
		})
		running = c.maybeBlockLengthBonus("else", e.Body, elseLine, running, incs)
		// Walk the else-if body at the same nesting+1 as the original if body.
		running = c.walkBody(e.Body, nesting+1, running, incs)
		// Continue the chain for any further else/else-if.
		if e.Else != nil {
			running = c.walkElseClause(e.Else, nesting, running, incs)
		}

	case *ast.BlockStmt:
		// Pure else: +1 flat.
		running++
		*incs = append(*incs, models.Increment{
			Kind: "else", Line: elseLine, Score: 1,
			Nesting: 0, RunningTotal: running,
		})
		running = c.maybeBlockLengthBonus("else", e, elseLine, running, incs)
		running = c.walkBody(e, nesting+1, running, incs)
	}

	return running
}

// walkSwitchStmt handles switch statements.
func (c *Calculator) walkSwitchStmt(s *ast.SwitchStmt, nesting, running int, incs *[]models.Increment) int {
	score := 1 + nesting
	line := c.fset.Position(s.Pos()).Line
	running += score
	*incs = append(*incs, models.Increment{
		Kind: "switch", Line: line, Score: score,
		Nesting: nesting, RunningTotal: running,
	})
	for _, stmt := range s.Body.List {
		if cc, ok := stmt.(*ast.CaseClause); ok {
			running = c.walkCaseClause(cc, nesting+1, running, incs)
		}
	}
	return running
}

// walkTypeSwitchStmt handles type-switch statements.
func (c *Calculator) walkTypeSwitchStmt(s *ast.TypeSwitchStmt, nesting, running int, incs *[]models.Increment) int {
	score := 1 + nesting
	line := c.fset.Position(s.Pos()).Line
	running += score
	*incs = append(*incs, models.Increment{
		Kind: "type switch", Line: line, Score: score,
		Nesting: nesting, RunningTotal: running,
	})
	for _, stmt := range s.Body.List {
		if cc, ok := stmt.(*ast.CaseClause); ok {
			running = c.walkCaseClause(cc, nesting+1, running, incs)
		}
	}
	return running
}

// walkSelectStmt handles select statements.
func (c *Calculator) walkSelectStmt(s *ast.SelectStmt, nesting, running int, incs *[]models.Increment) int {
	score := 1 + nesting
	line := c.fset.Position(s.Pos()).Line
	running += score
	*incs = append(*incs, models.Increment{
		Kind: "select", Line: line, Score: score,
		Nesting: nesting, RunningTotal: running,
	})
	for _, stmt := range s.Body.List {
		if cc, ok := stmt.(*ast.CommClause); ok {
			running = c.walkCommClause(cc, nesting+1, running, incs)
		}
	}
	return running
}

// walkCaseClause handles a case inside switch/type-switch.
// The case itself scores +1 flat; its body runs at nesting.
func (c *Calculator) walkCaseClause(cc *ast.CaseClause, nesting, running int, incs *[]models.Increment) int {
	line := c.fset.Position(cc.Pos()).Line
	running++
	*incs = append(*incs, models.Increment{
		Kind: "case", Line: line, Score: 1,
		Nesting: 0, RunningTotal: running,
	})
	for _, stmt := range cc.Body {
		running = c.walkStmt(stmt, nesting, running, incs)
	}
	return running
}

// walkCommClause handles a case inside select.
func (c *Calculator) walkCommClause(cc *ast.CommClause, nesting, running int, incs *[]models.Increment) int {
	line := c.fset.Position(cc.Pos()).Line
	running++
	*incs = append(*incs, models.Increment{
		Kind: "case", Line: line, Score: 1,
		Nesting: 0, RunningTotal: running,
	})
	for _, stmt := range cc.Body {
		running = c.walkStmt(stmt, nesting, running, incs)
	}
	return running
}

// walkExpr walks an expression looking for function literals (closures) and
// logical operator sequences.
func (c *Calculator) walkExpr(expr ast.Expr, nesting, running int, incs *[]models.Increment) int {
	if expr == nil {
		return running
	}
	switch e := expr.(type) {
	case *ast.FuncLit:
		score := 1 + nesting
		line := c.fset.Position(e.Pos()).Line
		running += score
		*incs = append(*incs, models.Increment{
			Kind: "closure", Line: line, Score: score,
			Nesting: nesting, RunningTotal: running,
		})
		running = c.walkBody(e.Body, nesting+1, running, incs)

	case *ast.BinaryExpr:
		if r := c.cfg.Rule("logical_operators"); r.Enabled {
			seqs := countLogicalSequences(e)
			if seqs > 0 {
				line := c.fset.Position(e.Pos()).Line
				w := scoreWeight(r.Weight) * seqs
				running += w
				*incs = append(*incs, models.Increment{
					Kind: "logical operators", Line: line, Score: w,
					Nesting: 0, RunningTotal: running,
				})
			}
		}
		// After counting at this level, descend only into non-logical sub-exprs
		// to avoid double-counting.
		running = c.walkExprNoLogical(e.X, nesting, running, incs)
		running = c.walkExprNoLogical(e.Y, nesting, running, incs)

	case *ast.CallExpr:
		running = c.walkCallExpr(e, nesting, running, incs)

	case *ast.CompositeLit:
		for _, elt := range e.Elts {
			running = c.walkExpr(elt, nesting, running, incs)
		}

	case *ast.UnaryExpr:
		running = c.walkExpr(e.X, nesting, running, incs)

	case *ast.ParenExpr:
		running = c.walkExpr(e.X, nesting, running, incs)

	case *ast.KeyValueExpr:
		running = c.walkExpr(e.Value, nesting, running, incs)

	case *ast.IndexExpr:
		running = c.walkExpr(e.X, nesting, running, incs)
		running = c.walkExpr(e.Index, nesting, running, incs)
	}
	return running
}

// walkExprNoLogical descends into an expression without re-counting logical
// operators at this level (they were already counted by the parent BinaryExpr).
func (c *Calculator) walkExprNoLogical(expr ast.Expr, nesting, running int, incs *[]models.Increment) int {
	if expr == nil {
		return running
	}
	switch e := expr.(type) {
	case *ast.FuncLit:
		return c.walkExpr(e, nesting, running, incs)
	case *ast.CallExpr:
		return c.walkCallExpr(e, nesting, running, incs)
	case *ast.BinaryExpr:
		// Only descend if NOT a logical operator (arithmetic etc. are fine).
		if e.Op != token.LAND && e.Op != token.LOR {
			running = c.walkExprNoLogical(e.X, nesting, running, incs)
			running = c.walkExprNoLogical(e.Y, nesting, running, incs)
		}
	case *ast.ParenExpr:
		// Parenthesised sub-expression starts a fresh logical sequence context.
		running = c.walkExpr(e.X, nesting, running, incs)
	case *ast.UnaryExpr:
		running = c.walkExprNoLogical(e.X, nesting, running, incs)
	}
	return running
}

// walkCallExpr handles a call expression, looking for function-literal args.
func (c *Calculator) walkCallExpr(call *ast.CallExpr, nesting, running int, incs *[]models.Increment) int {
	if call == nil {
		return running
	}
	running = c.walkExpr(call.Fun, nesting, running, incs)
	for _, arg := range call.Args {
		running = c.walkExpr(arg, nesting, running, incs)
	}
	return running
}

// walkDecl handles var/const declaration blocks.
func (c *Calculator) walkDecl(decl ast.Decl, nesting, running int, incs *[]models.Increment) int {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return running
	}
	for _, spec := range genDecl.Specs {
		if vs, ok := spec.(*ast.ValueSpec); ok {
			for _, val := range vs.Values {
				running = c.walkExpr(val, nesting, running, incs)
			}
		}
	}
	return running
}

// maybeBlockLengthBonus adds +1 if the block body is >= blockLengthThreshold
// lines (configurable via the block_length rule).
func (c *Calculator) maybeBlockLengthBonus(kind string, block *ast.BlockStmt, line, running int, incs *[]models.Increment) int {
	r := c.cfg.Rule("block_length")
	if !r.Enabled || block == nil {
		return running
	}
	minLines := r.MinLines
	if minLines == 0 {
		minLines = 10
	}
	bodyLines := c.fset.Position(block.End()).Line - c.fset.Position(block.Pos()).Line
	if bodyLines >= minLines {
		w := scoreWeight(r.Weight)
		running += w
		*incs = append(*incs, models.Increment{
			Kind:         kind + " with lines >= " + itoa(minLines),
			Line:         line,
			Score:        w,
			Nesting:      0,
			RunningTotal: running,
		})
	}
	return running
}

// --- Helpers -----------------------------------------------------------------

// hasCorrectComment returns true if funcDecl has a doc comment whose first
// word matches the function name (Go convention: // FuncName ...).
func hasCorrectComment(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
		return false
	}
	text := funcDecl.Doc.List[0].Text
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSpace(text)
	return strings.HasPrefix(text, funcDecl.Name.Name)
}

// countParams returns the total number of parameters across all param groups.
func countParams(funcDecl *ast.FuncDecl) int {
	if funcDecl.Type.Params == nil {
		return 0
	}
	count := 0
	for _, field := range funcDecl.Type.Params.List {
		if len(field.Names) == 0 {
			count++ // unnamed param
		} else {
			count += len(field.Names)
		}
	}
	return count
}

// countLogicalSequences counts contiguous sequences of && or || in a binary
// expression tree. Parenthesised sub-expressions start a new context.
// e.g. "a && b && c"         → 1 sequence
//
//	"a && b || c"         → 2 sequences (operator change)
//	"(a && b) || (c && d)"→ 2 sequences
func countLogicalSequences(expr ast.Expr) int {
	var ops []token.Token
	flattenLogical(expr, &ops)
	if len(ops) == 0 {
		return 0
	}
	seqs := 1
	for i := 1; i < len(ops); i++ {
		if ops[i] != ops[i-1] {
			seqs++
		}
	}
	return seqs
}

// flattenLogical extracts the left-recursive &&/|| chain from expr into ops.
// It stops at parenthesised sub-expressions so they form their own context.
func flattenLogical(expr ast.Expr, ops *[]token.Token) {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok {
		return
	}
	if bin.Op != token.LAND && bin.Op != token.LOR {
		return
	}
	flattenLogical(bin.X, ops)
	*ops = append(*ops, bin.Op)
	// Do NOT descend into bin.Y if it is parenthesised — that starts a new
	// logical context and will be counted separately by walkExpr.
	if _, paren := bin.Y.(*ast.ParenExpr); !paren {
		flattenLogical(bin.Y, ops)
	}
}

// scoreWeight converts a float64 rule weight to an integer score delta.
// A weight of 0 is treated as 1 (default).
func scoreWeight(w float64) int {
	if w <= 0 {
		return 1
	}
	return int(w)
}

// itoa converts an int to its decimal string representation without importing
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n >= 10 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	pos--
	buf[pos] = byte('0' + n)
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
