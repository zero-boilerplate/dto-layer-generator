package setup

import (
	. "github.com/francoishill/golang-web-dry/errors/checkerror"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func MustParseYAML(xamlBytes []byte) *DTOSetup {
	d := &dtoSetupYAML{}
	err := yaml.Unmarshal(xamlBytes, d)
	CheckError(err)

	return NewDTOSetupFromYAML(d)
}

func MustParseYAMLFile(filePath string) *DTOSetup {
	fileBytes, err := ioutil.ReadFile(filePath)
	CheckError(err)
	return MustParseYAML(fileBytes)
}
