
# Cognitive Complexity Linter for Golang

This project is a linter for Golang code that measures and reports the cognitive complexity of functions and statements. Cognitive complexity is a metric that quantifies how difficult it is for humans to understand and maintain a piece of code. It is based on the number and nesting of control structures, such as loops, conditionals, switches, etc. It helps developers write simple and readable code that avoids excessive branching and nesting.

## Installation

To install this project, you need to have Go installed on your system. You can download it from [here](^1^). Then, you can clone this repository using the following command:

```
git clone https://github.com/MikeMwita/cognitive-complexity-linter.git
```

## Usage

To use this project, you can run the following command in the root directory of the project:

```
go run cmd/main.go [options] [files]
```

The options are:

- `-h` or `--help`: show the help message and exit
- `-v` or `--version`: show the version number and exit
- `-c` or `--config`: specify the path to the configuration file
- `-o` or `--output`: specify the output format (`text`, `json`, or `xml`)

The files are the paths to the Go files or directories that you want to lint. If no files are given, the current directory is used.

The output will show the cognitive complexity score for each function and statement, along with the line number and the file name. For example:

```
Function main (main.go:5): 3
  IfStmt (main.go:7): 1
  ForStmt (main.go:12): 1
    IfStmt (main.go:14): 1
Function foo (main.go:20): 5
  SwitchStmt (main.go:22): 1
    CaseClause (main.go:23): 1
    CaseClause (main.go:25): 1
    CaseClause (main.go:27): 1
  IfStmt (main.go:30): 1
    IfStmt (main.go:31): 1
```

## Results

The results of running this project on a sample directory are as follows:

- Number of files: 124
- Number of functions: 245
- Highest complexity: 18
- Overall average complexity per function: 2.79
- Number of complex lines: 1

Example content of the `output.txt` file is:

```
./temp/services/app-auth/internal/core/models/http.go:17:1 - ResponseFromError has complexity: 18
  complexity = 1
  + 1 (found 'missing or wrong comment for function with more that 10 lines' at line: 17, complexity = 2)
  + 1 (found 'switch' at line: 23, complexity = 3)
  + 1 (found 'case' at line: 24, complexity = 4)
    + 2 (found 'branch' at line: 27, complexity = 6)
  + 1 (found 'case' at line: 28, complexity = 7)
    + 2 (found 'branch' at line: 31, complexity = 9)
  + 1 (found 'case' at line: 32, complexity = 10)
    + 2 (found 'branch' at line: 35, complexity = 12)
  + 1 (found 'case' at line: 36, complexity = 13)
    + 2 (found 'branch' at line: 39, complexity = 15)
  + 1 (found 'case' at line: 40, complexity = 16)
    + 2 (found 'branch' at line: 45, complexity = 18)
./temp/services/app-auth/internal/routes/middleware/cors.go:8:1 - Cors has complexity: 14
  complexity = 1
  + 1 (found 'missing or wrong comment for function with more that 10 lines' at line: 8, complexity = 2)
  + 1 (found 'range' at line: 18, complexity = 3)
    + 2 (found 'if' at line: 19, complexity = 5)
      + 3 (found 'branch' at line: 21, complexity = 8)
  + 1 (found 'if' at line: 25, complexity = 9)
  + 1 (found 'if with lines >= 10' at line: 25, complexity = 10)
  + 1 (found 'else' at line: 39, complexity = 11)
  + 1 (found 'else with lines >= 10' at line: 39, complexity = 12)
    + 2 (found 'if' at line: 32, complexity = 14)
```

## License

This project is licensed under the MIT License. See the [LICENSE](^2^) file for details..

