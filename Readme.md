
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
go run main.go [options] [files]
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

## Layout

This project follows the standard layout for Go projects, as described [here](https://github.com/golang-standards/project-layout). The main packages and files are:

- `cmd`: contains the main package that runs the linter
    - `main.go`: the entry point of the program
- `pkg`: contains the packages that implement the linter logic
    - `linter`: the core package that defines the linter interface and the common types and functions
        - `linter.go`: defines the linter interface and the lint function
        - `result.go`: defines the result type that represents a linting error or warning
        - `config.go`: defines the config type that represents the configuration for the linter
    - `complexity`: the package that contains the complexity calculator for the linter
        - `complexity.go`: defines the complexity interface and the complexity calculator
        - `assignstmt.go`: implements the complexity calculation for assignment statements
        - `blockstmt.go`: implements the complexity calculation for block statements
        - `branchstmt.go`: implements the complexity calculation for branch statements
        - `caseclause.go`: implements the complexity calculation for case clauses
        - `deferstmt.go`: implements the complexity calculation for defer statements
        - `forstmt.go`: implements the complexity calculation for for statements
        - `gostmt.go`: implements the complexity calculation for go statements
        - `ifstmt.go`: implements the complexity calculation for if statements
        - `labelstmt.go`: implements the complexity calculation for label statements
        - `rangestmt.go`: implements the complexity calculation for range statements
        - `switchstmt.go`: implements the complexity calculation for switch statements
        - `typeswitchstmt.go`: implements the complexity calculation for type switch statements
        - `comment.go`: implements the complexity calculation for comments
        - `params.go`: implements the complexity calculation for the number of parameters
        - `generics.go`: implements the complexity calculation for the use of generics
        - `length.go`: implements the complexity calculation for the length of functions
        - `operators.go`: implements the complexity calculation for the use of logical operators
        - `complexity.go`: implements the complexity calculation for the complexity of lines
- `test`: contains the tests for the project
    - `linter_test.go`: tests the linter package
    - `complexity_test.go`: tests the complexity package
- `build`: contains the scripts and files for building and packaging the project
    - `Dockerfile`: the file for creating a Docker image of the project
    - `Makefile`: the file for automating the build process
- `docs`: contains the documentation for the project
    - `README.md`: this file
    - `LICENSE`: the license file
    - `CHANGELOG.md`: the file that records the changes and updates of the project
    - `OWNERS`: the file that lists the owners and maintainers of the project

## Contributing

We welcome contributions from anyone who is interested in improving this project. Here are some ways you can contribute:

- Report issues or bugs on the [issue tracker](https://github.com/your-username/cognitive-complexity-linter/issues).
- Suggest new features or enhancements on the [issue tracker](https://github.com/your-username/cognitive-complexity-linter/issues).
- Submit pull requests with code changes or bug fixes. Please follow the [pull request template](https://github.com/your-username/cognitive-complexity-linter/blob/main/.github/PULL_REQUEST_TEMPLATE.md) and the [code of conduct](https://github.com/your-username/cognitive-complexity-linter/blob/main/CODE_OF_CONDUCT.md).
- Write tests or improve the test coverage. We use the [testing](https://pkg.go.dev/testing) package and the [testify](https://github.com/stretchr/testify) package for testing.
- Follow the code style and conventions of the Go language and this project. We use [gofmt](https://golang.org/cmd/gofmt/) and [golangci-lint](https://golangci-lint.run/) to format and lint the code.

## Contact

If you have any questions or feedback, you can contact us at:

- Email: michaelmasubo27@gmail.com
- Twitter: [@mike_mwita](https://twitter.com/Mikemwita)

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/your-username/cognitive-complexity-linter/blob/main/LICENSE) file for details.

## Acknowledgements

This project is inspired by or uses the following sources or libraries:

- [Go](https://golang.org/): the programming language we use
- [gofmt](https://golang.org/cmd/gofmt/): the tool for formatting Go code
- [golangci-lint](https://golangci-lint.run/): the tool for linting Go code
- [testing](https://pkg.go.dev/testing): the standard package for testing Go code
- [testify](https://github.com/stretchr/testify): the package for enhancing testing Go code
- [golang/example](https://github.com/golang/example): the repository for Go example projects
- [dave/rebecca](https://github.com/dave/rebecca): the tool for generating README files for Go projects
- [makeareadme](https://www.makeareadme.com/): the website for learning how to make a README file