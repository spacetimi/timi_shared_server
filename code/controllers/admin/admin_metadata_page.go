package admin

import (
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "net/http"
    "sort"
    "strings"
)

func showAdminMetadataPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {
    // Add link for back navigation
    adminPageObject.NavBackLinks = append(adminPageObject.NavBackLinks,
                                          NavBackLink{
                                              LinkName: "metadata",
                                              Href: "/admin/metadata",
                                          })
    switch request.URL.Path {

    case "/admin/metadata":
        showMetadataSelectPage(httpResponseWriter, request, adminPageObject)
        return

    case "/admin/metadata/app":
        showMetadataOverviewPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case "/admin/metadata/app/setCurrentVersions":
        showMetadataOverviewPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case "/admin/metadata/app/editVersion":
        showMetadataEditVersionPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case "/admin/metadata/shared":
        showMetadataOverviewPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case "/admin/metadata/shared/setCurrentVersions":
        showMetadataOverviewPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case "/admin/metadata/shared/editVersion":
        showMetadataEditVersionPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    default:
        logger.LogWarning("Unknown metadata route request|request url=" + request.URL.Path)
    }
}


func showMetadataSelectPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {
    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    err = templates.ExecuteTemplate(httpResponseWriter, "metadata_select_page_template.html", adminPageObject)

    if err != nil {
        logger.LogError("Error executing templates" +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

func showMetadataOverviewPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace) {
    pageObject := AdminMetadataPageObject{}
    pageObject.AdminPageObject = adminPageObject

    // Add link for back navigation
    pageObject.NavBackLinks = append(pageObject.NavBackLinks,
                                     NavBackLink{
                                         LinkName: space.String(),
                                         Href: "/admin/metadata/" + space.String(),
                                     })

    allVersions := metadata_service.Instance().GetAllVersions(space)
    allVersionsSorted := make([]string, len(allVersions))
    copy(allVersionsSorted, allVersions)
    sort.Strings(allVersionsSorted)
    sort.Sort(sort.Reverse(sort.StringSlice(allVersionsSorted)))

    pageObject.MetadataInfo = MetadataInfo{
        Space:space.String(),
        CurrentVersions: metadata_service.Instance().GetCurrentVersions(space),
        CurrentVersionsCSV: strings.Join(metadata_service.Instance().GetCurrentVersions(space), ","),
        AllVersions: allVersionsSorted,
    }

    // Check post arguments
    err := request.ParseForm()
    if err != nil {
        logger.LogError("error parsing form for metadata request" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    newCurrentVersionsCSV := request.Form.Get("currentVersionsCSV")
    // If new current csv arguments are sent, try to update and redirect to show success / failure
    if newCurrentVersionsCSV != "" {
        err := updateNewCurrentVersions(space, newCurrentVersionsCSV)
        messageToShow := "Successfully updated current versions."
        if err != nil {
            messageToShow = "Something went wrong updating current versions."
            pageObject.HasError = true
            pageObject.ErrorString = err.Error()
        }

        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: pageObject.AdminPageObject,
            SimpleMessage: messageToShow,
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    if err != nil {
        logger.LogError("error parsing templates" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = templates.ExecuteTemplate(httpResponseWriter, "metadata_overview_template.html", pageObject)
    if err != nil {
        logger.LogError("Error executing templates" +
            "|request url=" + request.URL.String() +
            "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func updateNewCurrentVersions(space metadata_typedefs.MetadataSpace, newCurrentVersionsCSV string) error {
    newCurrentVersions := strings.Split(strings.Replace(newCurrentVersionsCSV, " ", "", -1), ",")

    defer metadata_service.ReleaseInstanceRW()
    err := metadata_service.InstanceRW().SetCurrentVersions(newCurrentVersions, space)
    if err != nil {
        logger.LogError("error updating metadata current version" +
                        "|metadata space=" + space.String() +
                        "|new current versions=" + newCurrentVersionsCSV +
                        "|error=" + err.Error())
        return err
    }

    logger.LogInfo("Updated current versions for metadata space: " + space.String() +
                   " to: " + newCurrentVersionsCSV)

    return nil
}

func showMetadataEditVersionPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace) {
    pageObject := AdminEditMetadataPageObject{}
    pageObject.AdminPageObject = adminPageObject

    // Add links for back navigation
    pageObject.NavBackLinks = append(pageObject.NavBackLinks,
                                     NavBackLink{
                                         LinkName: space.String(),
                                         Href: "/admin/metadata/" + space.String(),
                                     })
    pageObject.NavBackLinks = append(pageObject.NavBackLinks,
                                     NavBackLink{
                                         LinkName: "editVersion",
                                         Href: "/admin/metadata/" + space.String() + "/editVersion",
                                     })

    pageObject.Space = space.String()

    // Check post arguments
    err := request.ParseForm()
    if err != nil {
        logger.LogError("error parsing form for metadata request" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    version := request.Form.Get("version")
    // If version argument is set, show metadata for that version
    if version != "" {
        pageObject.Version = version
    } else {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: pageObject.AdminPageObject,
            SimpleMessage: "Edit version: Invalid version",
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    templates, err := template.ParseGlob(config.GetTemplateFilesPath() + "/admin_tool/*")
    if err != nil {
        logger.LogError("error parsing templates" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = templates.ExecuteTemplate(httpResponseWriter, "metadata_edit_version_template.html", pageObject)
    if err != nil {
        logger.LogError("Error executing templates" +
                        "|request url=" + request.URL.String() +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    return
}

