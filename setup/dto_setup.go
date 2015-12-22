package setup

type PluginName string

func (p PluginName) String() string { return string(p) }

type OutputFilePath string

func (o OutputFilePath) String() string { return string(o) }

type DTOSetup struct {
	Name   string
	Url    string
	Output struct {
		Placeholder string
		Plugins     map[PluginName]OutputFilePath
	}
	Fields []*DTOField
}
