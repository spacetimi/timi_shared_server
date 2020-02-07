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

type AdminEditMetadataPageObject struct {
    AdminPageObject
    Space string
    Version string
    Items []AdminMetadataItem
}

type MetadataInfo struct {
    Space string
    CurrentVersions []string
    CurrentVersionsCSV string
    AllVersions []string
    IsUpToDate bool
}

type AdminMetadataItem struct {
    Key string
    Hash string
}

type AdminSimpleMessageObject struct {
    AdminPageObject
    SimpleMessage string
    BackLinkHref string
}