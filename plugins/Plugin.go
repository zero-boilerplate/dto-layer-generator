package plugins

import (
	"github.com/francoishill/dto-layer-generator/setup"
)

type Plugin interface {
	GenerateCode(logger Logger, dtoSetup *setup.DTOSetup) []byte
}
