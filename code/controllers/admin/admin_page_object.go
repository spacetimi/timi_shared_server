package admin

type AdminPageObject struct {
    IsLoggedIn bool
    LoggedInUser string

    AppName string
    AppEnvironment string

    HasError bool
    ErrorString string
}
