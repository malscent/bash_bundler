package main

import (
	"flag"
	"fmt"

	bundler "github.com/malscent/bash_bundler/pkg/bundler"
	logging "github.com/malscent/bash_bundler/pkg/log"
	"github.com/thatisuday/commando"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var logOptions logging.Options = logging.Options{}
var log = logf.Log.WithName("main")

func main() {
	// config logging
	logOptions.AddFlagSet(flag.CommandLine)
	logger := logging.New(&logOptions)
	logf.SetLogger(logger)
	// configure commando
	commando.SetExecutableName("bash_bundler").
		SetVersion("1.0.0").
		SetDescription("This simple tool bundles bash files into a single bash file.")

	commando.Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			//nolint:forbidigo
			fmt.Printf("%s\n", commando.DefaultCommandRegistry.Executable)
			//nolint:forbidigo
			fmt.Printf("%s\n", commando.DefaultCommandRegistry.Desc)
			//nolint:forbidigo
			fmt.Printf("Version: %s\n\n", commando.DefaultCommandRegistry.Version)
			//nolint:forbidigo
			fmt.Printf("See -h/--help for more information.")
		})

	commando.Register("bundle").
		SetShortDescription("bundle a bash script").
		SetDescription("Takes an entry bash script and bundles it and all its sources into a single output file.").
		AddFlag("entry,e", "The entrypoint to the bash script to bundle.", commando.String, nil).
		AddFlag("output,o", "The output file to write to", commando.String, nil).
		AddFlag("minify,m", "Minify the output file", commando.Bool, false).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			log.Info("Performing bundling", "entrypoint", flags["entry"].Value, "output", flags["output"].Value)

			entry, err := flags["entry"].GetString()
			bundler.CheckError(err)

			output, err := flags["output"].GetString()
			bundler.CheckError(err)

			min, err := flags["minify"].GetBool()
			bundler.CheckError(err)

			content, err := bundler.Bundle(entry, true)
			bundler.CheckError(err)

			if min {
				content, err = bundler.Minify(content)
				bundler.CheckError(err)
			}
			err = bundler.WriteToFile(output, content)
			bundler.CheckError(err)
		})
	commando.Parse(nil)
}
