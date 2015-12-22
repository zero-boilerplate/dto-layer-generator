package main

import (
	"fmt"
	"os"

	"github.com/francoishill/dto-layer-generator/helpers"
	"github.com/francoishill/dto-layer-generator/plugins"
	"github.com/francoishill/dto-layer-generator/setup"
)

type tmpLogger struct{}

func (t *tmpLogger) Error(msg string, params ...interface{}) {
	fmt.Println(fmt.Sprintf("[E] "+msg, params...))
}

func (t *tmpLogger) Warn(msg string, params ...interface{}) {
	fmt.Println(fmt.Sprintf("[W] "+msg, params...))
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

	for pluginName, outputFilePath := range dtoSetup.Output.Plugins {
		plugin := plugins.ParsePluginFromName(pluginName)
		helpers.InjectContentIntoFilePlaceholder(
			outputFilePath.String(),
			dtoSetup.Output.Placeholder,
			string(plugin.GenerateCode(logger, dtoSetup)))
	}
}
