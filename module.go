package hive

type Module struct {
	controllers []Controller
	config      ModuleConfig
}

func CreateModule() (module Module) {
	module = Module{}
	return module
}

func (m *Module) AddController(controller Controller) *Module {
	m.controllers = append(m.controllers, controller)
	return m
}

func (m *Module) SetConfig(config ModuleConfig) *Module {
	m.config = config
	return m
}

type ModuleConfig struct {
	Prefix string
}
