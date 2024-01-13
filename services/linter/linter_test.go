package linter

import (
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
