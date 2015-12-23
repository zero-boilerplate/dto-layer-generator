package main

import (
	"fmt"
	"os"

	//Required for plugins to register
	_ "github.com/zero-boilerplate/dto-layer-generator/plugins/client"
	_ "github.com/zero-boilerplate/dto-layer-generator/plugins/server"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/plugins"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

type tmpLogger struct{}

func (t *tmpLogger) Error(msg string, params ...interface{}) {
	fmt.Println(fmt.Sprintf("[E] "+msg, params...))
}

func (t *tmpLogger) Warn(msg string, params ...interface{}) {
	fmt.Println(fmt.Sprintf("[W] "+msg, params...))
}

func (t *tmpLogger) Info(msg string, params ...interface{}) {
	fmt.Println(fmt.Sprintf("[I] "+msg, params...))
}

func (t *tmpLogger) Debug(msg string, params ...interface{}) {
	fmt.Println(fmt.Sprintf("[D] "+msg, params...))
}

func main() {
	logger := &tmpLogger{}
	defer func() {
		if r := recover(); r != nil {
			logger.Error("FAILURE: %#v", r)
		}
	}()

	if len(os.Args) < 2 {
		panic("First command-line argument must be path to the dto setup YAML file.")
	}

	dtoSetup := setup.MustParseYAMLFile(os.Args[1])

	logger.Debug("Found %d plugins", len(dtoSetup.Output.Plugins))
	for pluginName, outputFilePath := range dtoSetup.Output.Plugins {
		plugin := plugins.ParsePluginFromName(pluginName)
		helpers.InjectContentIntoFilePlaceholder(
			logger,
			outputFilePath.String(),
			dtoSetup.Output.Placeholder,
			string(plugin.GenerateCode(logger, dtoSetup)))
	}
}
