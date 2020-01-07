package admin

import (
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "net/http"
)

func showAdminMetadataPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {
    // Add link for back navigation
    adminPageObject.NavBackLinks = append(adminPageObject.NavBackLinks,
                                          NavBackLink{
                                              LinkName: "metadata",
                                              Href: "/admin/metadata",
                                          })
    metadataPageObject := AdminMetadataPageObject{}
    metadataPageObject.AdminPageObject = adminPageObject

    switch request.URL.Path {

    case "/admin/metadata":
        showMetadataSelectPage(httpResponseWriter, request, metadataPageObject)
        return

    case "/admin/metadata/app":
        showMetadataOverviewPage(httpResponseWriter, request, metadataPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case "/admin/metadata/shared":
        showMetadataOverviewPage(httpResponseWriter, request, metadataPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    default:
        logger.LogWarning("Unknown metadata route request|request url=" + request.URL.Path)
    }
}


func showMetadataSelectPage(httpResponseWriter http.ResponseWriter, request *http.Request, pageObject AdminMetadataPageObject) {
    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    err = templates.ExecuteTemplate(httpResponseWriter, "metadata_select_page_template.html", pageObject)

    if err != nil {
        logger.LogError("Error executing templates" +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

func showMetadataOverviewPage(httpResponseWriter http.ResponseWriter, request *http.Request, pageObject AdminMetadataPageObject, space metadata_typedefs.MetadataSpace) {
    // Add link for back navigation
    pageObject.NavBackLinks = append(pageObject.NavBackLinks,
                                     NavBackLink{
                                         LinkName: space.String(),
                                         Href: "/admin/metadata/" + space.String(),
                                     })

    pageObject.MetadataInfo = MetadataInfo{
        Space:space.String(),
        CurrentVersions: metadata_service.Instance().GetCurrentVersions(space),
    }

    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    err = templates.ExecuteTemplate(httpResponseWriter, "metadata_overview_template.html", pageObject)

    if err != nil {
        logger.LogError("Error executing templates" +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}
