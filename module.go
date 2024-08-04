package hive

type Module struct {
	controllers []Controller
}

func (m *Module) AddController(controller Controller) *Module {
	m.controllers = append(m.controllers, controller)
	return m
}
