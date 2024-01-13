package linter

import (
	"github.com/MikeMwita/go-strict/datamodels"
	"github.com/MikeMwita/go-strict/services/complexity"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"testing"
)

func TestLinterService_LintFiles(t *testing.T) {
	type fields struct {
		config     *datamodels.LintConfig
		complexity *complexity.ComplexityService
		fileCount  int
		funcCount  int
	}
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*datamodels.LintResult
		wantErr bool
	}{
		{
			name: "Test valid Go file",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				files: []string{"./testdata/valid.go"},
			},
			want: []*datamodels.LintResult{
				{
					File:     "./testdata/valid.go",
					Message:  "No issues found",
					Severity: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "Test invalid Go file",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				files: []string{"./testdata/invalid.go"},
			},
			want: []*datamodels.LintResult{
				{
					File:     "./testdata/invalid.go",
					Message:  "expected 'package', found 'EOF'",
					Severity: "error",
				},
			},
			wantErr: true,
		},
		{
			name: "Test directory with Go files",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				files: []string{"./testdata"},
			},
			want: []*datamodels.LintResult{
				{
					File:     "./testdata/valid.go",
					Message:  "No issues found",
					Severity: "info",
				},
				{
					File:     "./testdata/invalid.go",
					Message:  "expected 'package', found 'EOF'",
					Severity: "error",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinterService{
				config:     tt.fields.config,
				complexity: tt.fields.complexity,
				fileCount:  tt.fields.fileCount,
				funcCount:  tt.fields.funcCount,
			}
			got, err := ls.LintFiles(tt.args.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("LintFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LintFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinterService_LintFunctions(t *testing.T) {
	type fields struct {
		config     *datamodels.LintConfig
		complexity *complexity.ComplexityService
		fileCount  int
		funcCount  int
	}
	type args struct {
		functions []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*datamodels.LintResult
		wantErr bool
	}{
		{
			name: "Test valid function",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				functions: []string{"func Add(x, y int) int { return x + y }"},
			},
			want: []*datamodels.LintResult{
				{
					File:     "tmpfile.go",
					Function: "Add",
					Message:  "No issues found",
					Severity: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "Test invalid function",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				functions: []string{"func Subtract(x, y int) int { return x - y }", "func Multiply(x, y int) int { return x * y }"},
			},
			want: []*datamodels.LintResult{
				{
					File:     "tmpfile.go",
					Function: "Subtract",
					Message:  "No issues found",
					Severity: "info",
				},
				{
					File:     "tmpfile.go",
					Function: "Multiply",
					Message:  "Line too long: func Multiply(x, y int) int { return x * y }",
					Severity: "warning",
				},
			},
			wantErr: false,
		},
		{
			name: "Test empty function",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				functions: []string{""},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinterService{
				config:     tt.fields.config,
				complexity: tt.fields.complexity,
				fileCount:  tt.fields.fileCount,
				funcCount:  tt.fields.funcCount,
			}
			got, err := ls.LintFunctions(tt.args.functions)
			if (err != nil) != tt.wantErr {
				t.Errorf("LintFunctions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LintFunctions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinterService_lintFile(t *testing.T) {
	type fields struct {
		config     *datamodels.LintConfig
		complexity *complexity.ComplexityService
		fileCount  int
		funcCount  int
	}
	type args struct {
		fset *token.FileSet
		f    *ast.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*datamodels.LintResult
		wantErr bool
	}{
		{
			name: "Test valid Go file",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				fset: token.NewFileSet(),
				f:    parseFile("./testdata/valid.go"),
			},
			want: []*datamodels.LintResult{
				{
					File:     "./testdata/valid.go",
					Function: "Add",
					Message:  "No issues found",
					Severity: "info",
				},
				{
					File:     "./testdata/valid.go",
					Function: "Subtract",
					Message:  "No issues found",
					Severity: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "Test invalid Go file",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				fset: token.NewFileSet(),
				f:    parseFile("./testdata/invalid.go"),
			},
			want: []*datamodels.LintResult{
				{
					File:     "./testdata/invalid.go",
					Function: "Multiply",
					Message:  "Line too long: func Multiply(x, y int) int { return x * y }",
					Severity: "warning",
				},
			},
			wantErr: false,
		},
		{
			name: "Test empty file name",
			fields: fields{
				config:     &datamodels.LintConfig{MaxComplexity: 10, MaxLineLength: 80},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				fset: token.NewFileSet(),
				f:    &ast.File{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinterService{
				config:     tt.fields.config,
				complexity: tt.fields.complexity,
				fileCount:  tt.fields.fileCount,
				funcCount:  tt.fields.funcCount,
			}
			got, err := ls.lintFile(tt.args.fset, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("lintFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lintFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// parseFile is a helper function that parses a Go file and returns an *ast.File
func parseFile(filename string) *ast.File {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}

func TestLinterService_lintFunction(t *testing.T) {
	type fields struct {
		config     *datamodels.LintConfig
		complexity *complexity.ComplexityService
		fileCount  int
		funcCount  int
	}
	type args struct {
		fset     *token.FileSet
		funcDecl *ast.FuncDecl
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *datamodels.LintResult
		wantErr bool
	}{
		{
			name: "Test simple function",
			fields: fields{
				config:     &datamodels.LintConfig{Threshold: 10},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				fset:     token.NewFileSet(),
				funcDecl: parseFunction("func Add(x, y int) int { return x + y }"),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Test complex function",
			fields: fields{
				config:     &datamodels.LintConfig{Threshold: 10},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				fset:     token.NewFileSet(),
				funcDecl: parseFunction("func Fibonacci(n int) int { if n <= 1 { return n } return Fibonacci(n-1) + Fibonacci(n-2) }"),
			},
			want: &datamodels.LintResult{
				Line:     1,
				Severity: "warning",
				Message:  "function has a cognitive complexity of 11 which is higher than the threshold of 10\n+ 1 (found at line: 1)\n+ 1 (found 'if' at line: 1)\n+ 9 (found at line: 2)",
			},
			wantErr: false,
		},
		{
			name: "Test nil function",
			fields: fields{
				config:     &datamodels.LintConfig{Threshold: 10},
				complexity: &complexity.ComplexityService{},
				fileCount:  0,
				funcCount:  0,
			},
			args: args{
				fset:     token.NewFileSet(),
				funcDecl: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LinterService{
				config:     tt.fields.config,
				complexity: tt.fields.complexity,
				fileCount:  tt.fields.fileCount,
				funcCount:  tt.fields.funcCount,
			}
			got, err := ls.lintFunction(tt.args.fset, tt.args.funcDecl)
			if (err != nil) != tt.wantErr {
				t.Errorf("lintFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lintFunction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// parseFunction -->  parses a function declaration and returns an *ast.FuncDecl
func parseFunction(code string) *ast.FuncDecl {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		panic(err)
	}
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			return funcDecl
		}
	}
	return nil
}

func TestNewLinterService(t *testing.T) {
	type args struct {
		config     *datamodels.LintConfig
		complexity *complexity.ComplexityService
	}
	tests := []struct {
		name string
		args args
		want *LinterService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLinterService(tt.args.config, tt.args.complexity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLinterService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createTempFile(t *testing.T) {
	type args struct {
		functions []string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createTempFile(tt.args.functions)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTempFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTempFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
