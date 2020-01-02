package admin

import (
    "errors"
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "net/http"
)

const kCookieName = "jwtTokenForAdminUser"

func AdminController(httpResponseWriter http.ResponseWriter, request *http.Request) {
    adminPageObject := AdminPageObject {
        AppName: config.GetAppName(),
        AppEnvironment: "Unknown",
        IsLoggedIn: false,
        LoggedInUser: "",
        HasError: false,
        ErrorString: "",
    }

    switch config.GetEnvironmentConfiguration().AppEnvironment {
    case config.LOCAL: adminPageObject.AppEnvironment = "Local"
    case config.TEST: adminPageObject.AppEnvironment = "Test"
    case config.STAGING: adminPageObject.AppEnvironment = "Staging"
    case config.PRODUCTION: adminPageObject.AppEnvironment = "Production"
    default: adminPageObject.AppEnvironment = "Unknown"
    }

    // If request is for the login page, just show that

    if request.URL.Path == "/admin/login" {
        showLoginPage(httpResponseWriter, request, adminPageObject)
        return
    }

    // If not, make sure the user is logged in as admin

    hasLoggedIn, username, err := hasUserLoggedIn(request)
    if err != nil {
        logger.LogWarning("error checking if admin user logged in" +
            "|request URL=" + request.URL.Path +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    // If user is not logged in as admin, redirect to login page

    if !hasLoggedIn {
        http.Redirect(httpResponseWriter, request, "/admin/login", http.StatusSeeOther)
        return
    }

    adminPageObject.IsLoggedIn = true
    adminPageObject.LoggedInUser = username

    switch request.URL.String() {

    case "/admin": showAdminPage(httpResponseWriter, request, adminPageObject)
    case "/admin/": showAdminPage(httpResponseWriter, request, adminPageObject)

    default:
        logger.LogWarning("unknown route request for admin controller" +
            "|request URL=" + request.URL.Path)
        httpResponseWriter.WriteHeader(http.StatusNotFound)
    }
}

func hasUserLoggedIn(request *http.Request) (bool, string, error) {

    jwtCookie, err := request.Cookie(kCookieName)

    if err != nil {
        if err == http.ErrNoCookie {
            return false, "", nil
        }
        return false, "", errors.New("error trying to get admin login token cookie: " + err.Error())
    }

    if jwtCookie == nil {
        return false, "", errors.New("unknown error getting admin login token cookie")
    }

    ok, username, err := checkAdminLoginClaim(jwtCookie.Value)
    if err != nil {
        return false, "", errors.New("error validation admin login claim: " + err.Error())
    }
    if ok {
        return true, username, nil
    } else {
        return false, "", errors.New("unknown error validating admin login claim")
    }
}

func showLoginPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {

    err := request.ParseForm()
    if err != nil {
        logger.LogError("error parsing form for login request" +
            "|request url=" + request.URL.Path +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    username := request.Form.Get("username")

    // If the client has sent a {username, password} try to login
    if username != "" {
        password := request.Form.Get("password")
        response, err := tryLoginWithAdminCredentials(&AdminLoginRequest{username:username, password:password})
        if err != nil || response == nil {
            adminPageObject.HasError = true
            if err != nil {
                adminPageObject.ErrorString = err.Error()
            } else {
                adminPageObject.ErrorString = "unknown error occurred"
            }
            logger.LogWarning("problem authenticating admin user" +
                "|username=" + username +
                "|error=" + adminPageObject.ErrorString)
        } else {
            cookie := http.Cookie{Name: kCookieName, Value: response.jwtTokenString, Expires: response.expirationTime}
            http.SetCookie(httpResponseWriter, &cookie)
            http.Redirect(httpResponseWriter, request, "/admin", http.StatusSeeOther)
            return
        }
    }

    // Show login page

    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    err = templates.ExecuteTemplate(httpResponseWriter, "admin_login_template.html", adminPageObject)

    if err != nil {
        logger.LogError("Error executing templates" +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

func showAdminPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {

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
