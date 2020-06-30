package shared_init

type IAppInitializer interface {
	AppName() string
	AppInit() error
}

type RequiredServices struct {
	Services []string
}
