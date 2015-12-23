package plugins

import (
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

type Plugin interface {
	GenerateCode(logger Logger, dtoSetup *setup.DTOSetup) []byte
}
