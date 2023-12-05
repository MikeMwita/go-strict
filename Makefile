## tools: install/update all project tools
tools:
	@echo
	@echo "> delve : debugger"
	@echo -n "before: "
	@dlv version
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@echo -n "after: "
	@dlv version
	@echo

	@echo
	@echo "> golangci-lint : a linter for checking source code"
	@echo -n "before: "
	@golangci-lint version
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo -n "after: "
	@golangci-lint version
	@echo

	@echo "> identypo : finds typos in identifiers"
	@go install github.com/alexkohler/identypo/cmd/identypo@latest
	@echo

	@echo "> nakedret : a naked return linter"
	@go install github.com/alexkohler/nakedret/cmd/nakedret@latest
	@echo

	@echo "> gosec : a Security Checker"
	@echo -n "before: "
	@gosec -version
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo -n "after: "
	@gosec -version
	@echo

	@echo "> staticcheck : a static analysis linter"
	@echo -n "before: "
	@staticcheck -version
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo -n "after: "
	@staticcheck -version
	@echo

	@echo "> revive : a fast linter"
	@echo -n "before: "
	@revive -version
	@go install github.com/mgechev/revive@latest
	@echo -n "after: "
	@revive -version
	@echo

	@echo "> go-safer : a linter for reporting reflect.SliceHeader and reflect.StringHeader misuses"
	@go install github.com/jlauinger/go-safer@latest
	@echo

	@echo "> govulncheck : a vulnerability linter"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo

	@echo "> go-consistent : a code consistency checker"
	@go install github.com/quasilyte/go-consistent@latest
	@echo

	# @echo "> godoc : for looking at generated inline documentation"
	# @go install golang.org/x/tools/cmd/godoc@latest
	# @echo

	@echo "> pkgsite : for looking at generated inline documentation"
	@go install golang.org/x/pkgsite/cmd/pkgsite@latest
	@echo

	@echo "> gocoverstats : for looking at code coverage statistics"
	@go install gitlab.com/fgmarand/gocoverstats@latest
	@echo

	@echo "> go-hasdefault : linter for switch checking for default"
	@go install github.com/nathants/go-hasdefault@latest
	@echo

	@echo "> go-coverage : code coverage viewer"
	@go install github.com/gojekfarm/go-coverage@latest
	@echo

	@echo "> gocritic : linter"
	@go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	@echo

	@echo "> nilness : inspects the control-flow graph of an SSA function"
	@GO111MODULE=on go install golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness@latest
	@echo

	@echo "> ineffassign : Detect ineffectual assignments"
	@go install github.com/gordonklaus/ineffassign@latest
	@echo

	@echo "> gocritic : source code linter"
	@go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	@echo

	@echo "> goreleaser : builds various Go binaries"
	@go install github.com/goreleaser/goreleaser@latest
	@echo

	@echo "> mockgen : for generating mock files"
	@echo -n "before: "
	@mockgen -version
	@go install github.com/golang/mock/mockgen@latest
	@echo -n "after: "
	@mockgen -version
	@echo

	@echo "> smrcptr : detects mixing pointer and value method receivers"
	@go install github.com/nikolaydubina/smrcptr@latest
	@echo

	@echo "> go-cleanarch : checks for clean architecture"
	@go install github.com/roblaszczak/go-cleanarch@latest

	@echo "> gofumpt : source code formatter"
	@go install mvdan.cc/gofumpt@latest

	@echo "> ireturn : Accept Interfaces, Return Concrete Types"
	@go install github.com/butuzov/ireturn/cmd/ireturn@latest



linters:
	@echo -n "──── "
	@golangci-lint version
	@golangci-lint run --fix
	@echo "──── vet"
	@go vet ./...
	@echo -n "──── revive "
	@revive -version
	@revive -formatter friendly -exclude ./vendor/... ./...
	@echo "──── smrcptr"
	@smrcptr ./...
	@echo "──── nilness"
	@nilness ./...
	@echo "──── identypo"
	@identypo ./...
	@echo "──── nakedret"
	@nakedret ./...
	@echo -n "──── "
	@staticcheck -version
	@staticcheck -checks all -tests=false ./...
	@echo "──── go-safer"
	@go-safer ./...
	@echo "──── ineffassign"
	@ineffassign ./...
	@echo "──── gocritic"
	@gocritic check ./...
	@echo -n "──── gosec "
	@gosec -version
	@gosec -fmt=golint -quiet ./...
	@echo "──── govulncheck"
	@govulncheck ./...
	@echo "──── go-consistent"
	@go-consistent ./...