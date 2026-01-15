package main

import (
	"fmt"
	"os"

	gen "github.com/zllangct/ecs/generator/golang"
	"github.com/zllangct/rockmem/cmd/rockmem/cmds"
	"github.com/zllangct/rockmem/generator/golang"
)

func main() {
	err := golang.RegisterCodeGenerators(gen.NewComponentGenerator())
	if err != nil {
		panic(err)
	}
	if err := cmds.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
