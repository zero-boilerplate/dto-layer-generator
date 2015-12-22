package main

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/ghodss/yaml"
	"io/ioutil"

	"github.com/francoishill/dto-layer-generator/plugins"
	"github.com/francoishill/dto-layer-generator/setup"
)

type iplugin interface {
	GenerateCode(dtoSetup *setup.DTOSetup) []byte
}

func main() {
	fileBytes, err := ioutil.ReadFile("example.yml")
	CheckError(err)

	d := &setup.DTOSetup{}
	err = yaml.Unmarshal(fileBytes, d)
	CheckError(err)

	pluginsWithFilePaths := map[iplugin]string{
		new(plugins.GoPlugin): "out/out.go",
	}

	for plugin, filePath := range pluginsWithFilePaths {
		err = ioutil.WriteFile(filePath, plugin.GenerateCode(d), 0600)
		CheckError(err)
	}
}
