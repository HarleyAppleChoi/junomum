package main

import (
	"os"

	"github.com/forbole/flowJuno/cmd/parse"

	"github.com/forbole/flowJuno/modules/messages"
	"github.com/forbole/flowJuno/modules/registrar"

	"github.com/forbole/flowJuno/cmd"
)

func main() {
	// Config the runner
	config := cmd.NewConfig("flowjuno").
		WithParseConfig(parse.NewConfig().
			WithRegistrar(registrar.NewDefaultRegistrar(
				messages.CosmosMessageAddressesParser,
			)),
		)

	// Run the commands and panic on any error
	exec := cmd.BuildDefaultExecutor(config)
	err := exec.Execute()
	if err != nil {
		os.Exit(1)
	}
}