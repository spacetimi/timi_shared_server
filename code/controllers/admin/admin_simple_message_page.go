package admin

import (
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "net/http"
)

func showSimpleMessagePage(httpResponseWriter http.ResponseWriter, request *http.Request, pageObject AdminSimpleMessageObject) {
    templates, err := template.ParseGlob(config.GetSharedTemplateFilesPath() + "/admin_tool/*")
    if err != nil {
        logger.LogError("error parsing templates" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = templates.ExecuteTemplate(httpResponseWriter, "simple_message_template.html", pageObject)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.String() +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
    return
}
