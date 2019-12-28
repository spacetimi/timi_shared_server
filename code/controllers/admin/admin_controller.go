package admin

import (
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "net/http"
)

func AdminController(httpResponseWriter http.ResponseWriter, request *http.Request) {
    adminPageObject := AdminPageObject {
        AppName: config.GetAppName(),
        AppEnvironment: "Unknown",
        IsLoggedIn: false,
        LoggedInUser: "",
    }
    switch config.GetEnvironmentConfiguration().AppEnvironment {
    case config.LOCAL: adminPageObject.AppEnvironment = "Local"
    case config.TEST: adminPageObject.AppEnvironment = "Test"
    case config.STAGING: adminPageObject.AppEnvironment = "Staging"
    case config.PRODUCTION: adminPageObject.AppEnvironment = "Production"
    default: adminPageObject.AppEnvironment = "Unknown"
    }

    switch request.URL.String() {

    case "/admin": showAdminPage(httpResponseWriter, request, adminPageObject)
    case "/admin/": showAdminPage(httpResponseWriter, request, adminPageObject)

    case "/admin/login": showLoginPage(httpResponseWriter, request, adminPageObject)
    case "/admin/login/": showLoginPage(httpResponseWriter, request, adminPageObject)

    default:
        httpResponseWriter.WriteHeader(http.StatusNotFound)
    }
}

func showAdminPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {
    if !adminPageObject.IsLoggedIn {
        // Redirect to login page
        var newUrl string
        switch request.URL.String() {
        case "/admin": newUrl = "admin/login"
        case "/admin/": newUrl = "login"
        default:
            httpResponseWriter.WriteHeader(http.StatusNotFound)
            return
        }

        http.Redirect(httpResponseWriter, request, newUrl, http.StatusSeeOther)
        return
    }

    t := template.New("admin_tool_template.html")
    templateFilePath := config.GetTemplateFilesPath() + "/admin_tool/" + "admin_tool_template.html"

    t, err := t.ParseFiles(templateFilePath)
    if err != nil {
        logger.LogError("Error parsing template file" +
            "|file path=" + templateFilePath +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        fmt.Fprintf(httpResponseWriter, "Error loading page")
    }
    err = t.Execute(httpResponseWriter, adminPageObject)
    if err != nil {
        logger.LogError("Error executing template" +
            "|file path=" + templateFilePath +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        fmt.Fprintf(httpResponseWriter, "Error loading page")
    }
}

func showLoginPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {

    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    err = templates.ExecuteTemplate(httpResponseWriter, "admin_login_template.html", adminPageObject)

    if err != nil {
        logger.LogError("Error executing templates" +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}
