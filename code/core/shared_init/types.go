package shared_init

type IAppInitializer interface {
	AppName() string
	AppInit() error
}

type RequiredServicesConfig struct {
	Services []string
}

func (rsc *RequiredServicesConfig) OnConfigLoaded() {
	// Nothing to do yet
}
