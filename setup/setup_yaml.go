package setup

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func MustParseYAML(xamlBytes []byte) *DTOSetup {
	d := &DTOSetup{}
	err := yaml.Unmarshal(xamlBytes, d)
	CheckError(err)

	return d
}

func MustParseYAMLFile(filePath string) *DTOSetup {
	fileBytes, err := ioutil.ReadFile(filePath)
	CheckError(err)
	return MustParseYAML(fileBytes)
}
