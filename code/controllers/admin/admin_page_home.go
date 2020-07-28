package admin

import (
	"html/template"
	"net/http"

	"github.com/spacetimi/timi_shared_server/v2/code/config"
	"github.com/spacetimi/timi_shared_server/v2/utils/logger"
)

func showAdminPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {

	templates, err := template.ParseGlob(config.GetSharedTemplateFilesPath() + "/admin_tool/*")
	err = templates.ExecuteTemplate(httpResponseWriter, "admin_page_home_template.html", adminPageObject)

	if err != nil {
		logger.LogError("Error executing templates" +
			"|request url=" + request.URL.String() +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
