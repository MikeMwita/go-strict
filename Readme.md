
# Cognitive Complexity Linter for Golang

This project is a linter for Golang code that measures and reports the cognitive complexity of functions and statements. Cognitive complexity is a metric that quantifies how difficult it is for humans to understand and maintain a piece of code. It is based on the number and nesting of control structures, such as loops, conditionals, switches, etc. It helps developers write simple and readable code that avoids excessive branching and nesting.

## Installation

To install this project, you need to have Go installed on your system. You can download it from [here](https://golang.org/dl/).
Then, you can clone this repository using the following command:

```bash
git clone https://github.com/MikeMwita/cognitive-complexity-linter.git
```

## Usage

To use this project, you can run the following command in the root directory of the project:

```bash
go run initializers/main.go [options] [files]
```

The options are:

- `-h` or `--help`: show the help message and exit
- `-v` or `--version`: show the version number and exit
- `-c` or `--config`: specify the path to the configuration file
- `-o` or `--output`: specify the output format (text, json, or xml)

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



## Contact

If you have any questions or feedback, you can contact us at:

- Email: michaelmasubo27@gmail.com
- Twitter: [@mike_mwita](https://twitter.com/Mikemwita)

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/your-username/cognitive-complexity-linter/blob/main/LICENSE) file for details.



