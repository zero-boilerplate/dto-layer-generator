package plugins

import (
	"github.com/francoishill/dto-layer-generator/setup"
)

func ParsePluginFromName(name setup.PluginName) Plugin {
	if plugin, ok := registeredPlugins[name.String()]; ok {
		return plugin
	}
	panic("Unknown plugin name '" + name + "'. Perhaps it is not registered with `RegisterPlugin()`")
}
