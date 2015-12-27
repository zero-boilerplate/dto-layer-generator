package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"

	//Required for plugins to register
	_ "github.com/zero-boilerplate/dto-layer-generator/plugins/client/java_retrofit"
	_ "github.com/zero-boilerplate/dto-layer-generator/plugins/server/go_common_ddd"

	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/plugins"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

type tmpLogger struct {
	errColor   *color.Color
	warnColor  *color.Color
	infoColor  *color.Color
	debugColor *color.Color
}

func (t *tmpLogger) Error(msg string, params ...interface{}) {
	t.errColor.Println(fmt.Sprintf("[E] "+msg, params...))
}

func (t *tmpLogger) Warn(msg string, params ...interface{}) {
	t.warnColor.Println(fmt.Sprintf("[W] "+msg, params...))
}

func (t *tmpLogger) Info(msg string, params ...interface{}) {
	t.infoColor.Println(fmt.Sprintf("[I] "+msg, params...))
}

func (t *tmpLogger) Debug(msg string, params ...interface{}) {
	t.debugColor.Println(fmt.Sprintf("[D] "+msg, params...))
}

func main() {
	logger := &tmpLogger{
		color.New(color.FgHiRed, color.Bold),
		color.New(color.FgMagenta),
		color.New(color.FgHiWhite),
		color.New(color.FgWhite),
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("FAILURE: %#v", r)
			os.Exit(2)
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
