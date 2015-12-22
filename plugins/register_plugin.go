package plugins

var registeredPlugins map[string]Plugin = make(map[string]Plugin)

func RegisterPlugin(alias string, plugin Plugin) {
	registeredPlugins[alias] = plugin
}
