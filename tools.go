// +build tools

package main

// This file replace the bootstrap.json
import (
	_ "golang.org/x/lint/golint",
	_ "github.com/fzipp/gocyclo",
	_ "github.com/gordonklaus/ineffassign",
	_ "github.com/client9/misspell/cmd/misspell",
	_ "golang.org/x/mobile/cmd/gomobile",
	_ "golang.org/x/tools/go/ast/astutil",
	_ "golang.org/x/tools/go/internal",
	_ "golang.org/x/tools/go/gcexportdata",
	_ "github.com/golangci/golangci-lint",
	_ "github.com/shuLhan/go-bindata/cmd/go-bindata"

	_ "github.com/ugorji/go"
)
