package shared_init

type IAppInitializer interface {
	AppInit() error
}

type RequiredServices struct {
	Services []string
}

