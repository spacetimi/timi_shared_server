package admin

type AdminPageObject struct {
    IsLoggedIn bool
    LoggedInUser string

    AppName string
    AppEnvironment string

    HasError bool
    ErrorString string

    NavBackLinks []NavBackLink
}

type NavBackLink struct {
    LinkName string
    Href string
}

type AdminMetadataPageObject struct {
    AdminPageObject
    MetadataInfo
}

type MetadataInfo struct {
    Space string
}
