package plugins

import (
	"github.com/zero-boilerplate/dto-layer-generator/helpers"
	"github.com/zero-boilerplate/dto-layer-generator/setup"
)

type Plugin interface {
	GenerateCode(logger helpers.Logger, dtoSetup *setup.DTOSetup) []byte
}
