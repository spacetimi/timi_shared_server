package admin

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_factory"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "html/template"
    "io"
    "net/http"
    "path/filepath"
    "regexp"
    "sort"
    "strings"
    "time"
)

const kMetadataRoute_SelectSpace = "METADATA_SELECT_SPACE"
const kMetadataRoute_AppOverview = "METADATA_APP_OVERVIEW"
const kMetadataRoute_SharedOverview = "METADATA_SHARED_OVERVIEW"
const kMetadataRoute_AppEditVersion = "METADATA_APP_EDIT_VERSION"
const kMetadataRoute_SharedEditVersion = "METADATA__SHARED_EDIT_VERSION"
const kMetadataRoute_AppViewMetadata = "METADATA_APP_VIEW_METADATA"
const kMetadataRoute_SharedViewMetadata = "METADATA_SHARED_VIEW_METADATA"
const kMetadataRoute_AppDownload = "METADATA_APP_DOWNLOAD"
const kMetadataRoute_SharedDownload = "METADATA_SHARED_DOWNLOAD"
const kMetadataRoute_AppUpload = "METADATA_APP_UPLOAD"
const kMetadataRoute_SharedUpload = "METADATA_SHARED_UPLOAD"

var kAdminMetadataRoutes = map[string]string{
    "/admin/metadata$": kMetadataRoute_SelectSpace,
    "/admin/metadata/app$": kMetadataRoute_AppOverview,
    "/admin/metadata/app/setCurrentVersions$": kMetadataRoute_AppOverview,
    "/admin/metadata/app/editVersion/[0-9]+\\.[0-9]+$": kMetadataRoute_AppEditVersion,
    "/admin/metadata/app/view/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_AppViewMetadata,
    "/admin/metadata/app/download/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_AppDownload,
    "/admin/metadata/app/upload/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_AppUpload,
    "/admin/metadata/shared$": kMetadataRoute_SharedOverview,
    "/admin/metadata/shared/setCurrentVersions$": kMetadataRoute_SharedOverview,
    "/admin/metadata/shared/editVersion/[0-9]+\\.[0-9]+$": kMetadataRoute_SharedEditVersion,
    "/admin/metadata/shared/view/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_SharedViewMetadata,
    "/admin/metadata/shared/download/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_SharedDownload,
    "/admin/metadata/shared/upload/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_SharedUpload,
}

var kAdminMetadataRouteRegexToRouteName map[*regexp.Regexp]string

/** Package init **/
func init() {
    kAdminMetadataRouteRegexToRouteName = make(map[*regexp.Regexp]string, len(kAdminMetadataRoutes))
    for route, routeName := range kAdminMetadataRoutes {
        reg, err := regexp.Compile(route)
        if err != nil {
            logger.LogError("bad route regex in admin metadata controller" +
                            "|regex=" + route +
                            "|error=" + err.Error())
            continue
        }
        kAdminMetadataRouteRegexToRouteName[reg] = routeName
    }
}

func showAdminMetadataPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {
    // Add link for back navigation
    adminPageObject.NavBackLinks = append(adminPageObject.NavBackLinks,
                                          NavBackLink{
                                              LinkName: "metadata",
                                              Href: "/admin/metadata",
                                          })

    matchingRoute := getRouteNameForRequest(kAdminMetadataRouteRegexToRouteName, request.URL.Path)

    switch matchingRoute {

    case kMetadataRoute_SelectSpace:
        showMetadataSelectPage(httpResponseWriter, request, adminPageObject)
        return

    case kMetadataRoute_AppOverview:
        showMetadataOverviewPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_AppEditVersion:
        showMetadataEditVersionPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_AppViewMetadata:
        showMetadataViewOrDownloadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP, true)
        return

    case kMetadataRoute_AppDownload:
        showMetadataViewOrDownloadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP, false)
        return

    case kMetadataRoute_AppUpload:
        showMetadataUploadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_SharedOverview:
        showMetadataOverviewPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case kMetadataRoute_SharedEditVersion:
        showMetadataEditVersionPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case kMetadataRoute_SharedViewMetadata:
        showMetadataViewOrDownloadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED, true)
        return

    case kMetadataRoute_SharedDownload:
        showMetadataViewOrDownloadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED, false)
        return

    case kMetadataRoute_SharedUpload:
        showMetadataUploadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
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
    pageObject.Space = space.String()

    // Parse url for editing version
    versionString := filepath.Base(request.URL.Path)
    version, err := core.GetAppVersionFromString(versionString)
    if err != nil {
        logger.LogError("error parsing editing version from url" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }

    validVersion, err := metadata_service.Instance().IsVersionValid(version.String(), space)
    if !validVersion {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: pageObject.AdminPageObject,
            SimpleMessage: "Invalid version: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    // Add links for back navigation
    pageObject.NavBackLinks = append(pageObject.NavBackLinks,
                                     NavBackLink{
                                         LinkName: space.String(),
                                         Href: "/admin/metadata/" + space.String(),
                                     })
    pageObject.NavBackLinks = append(pageObject.NavBackLinks,
                                     NavBackLink{
                                         LinkName: "editVersion (" +  version.String() + ")",
                                         Href: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
                                     })


    pageObject.Version = version.String()

    metadataItems, err := metadata_service.Instance().GetMetadataItemsInVersion(version.String(), space)
    if err != nil {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: pageObject.AdminPageObject,
            SimpleMessage: "Error finding metadata items: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }
    for _, metadataItem := range metadataItems {
        pageObject.Items = append(pageObject.Items, AdminMetadataItem{
                                                        Key:metadataItem.MetadataKey,
                                                        Hash:metadataItem.Hash,
                                                    })
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

func showMetadataViewOrDownloadPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace, viewOnly bool) {

    // Parse url for editing version and metadata item key
    tokens := strings.Split(request.URL.Path, "/")
    if len(tokens) < 2 {
        logger.LogError("malformed request url in metadata download request" +
                        "|request url=" + request.URL.Path)
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }
    metadataItemKey := tokens[len(tokens) - 1]
    versionString := tokens[len(tokens) - 2]

    version, err := core.GetAppVersionFromString(versionString)
    if err != nil {
        logger.LogError("error parsing editing version from url" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }

    validVersion, err := metadata_service.Instance().IsVersionValid(version.String(), space)
    if !validVersion {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Invalid version: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }


    metadataItem, err := metadata_factory.InstantiateMetadataItem(metadataItemKey)
    if err != nil {
        logger.LogError("failed to instantiate metadata item in admin metadata download" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Failed to instantiate metadata item: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    err = metadata_service.Instance().GetMetadataItem(metadataItem, version)
    if err != nil {
        logger.LogError("failed to fetch metadata item in admin metadata download" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Failed to fetch metadata item: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    content, err := json.Marshal(metadataItem)
    if err != nil {
        logger.LogError("failed to serialize metadata item in admin metadata download" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Failed to serialize metadata item: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    if viewOnly {
        // Just write out the serialized metadata item
        _, err := fmt.Fprintln(httpResponseWriter, string(content))
        if err != nil {
            logger.LogError("error writing metadata json" +
                "|metadata item key=" + metadataItemKey +
                "|error=" + err.Error())
            httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        }

    } else {
        // Mark the returned content as downloadable to the browser
        httpResponseWriter.Header().Add("Content-Disposition", "Attachment")

        http.ServeContent(httpResponseWriter, request, metadataItemKey + ".json", time.Now(), bytes.NewReader([]byte(content)))
    }
}

func showMetadataUploadPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace) {

    // Parse url for editing version and metadata item key
    tokens := strings.Split(request.URL.Path, "/")
    if len(tokens) < 2 {
        logger.LogError("malformed request url in metadata download request" +
                        "|request url=" + request.URL.Path)
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }
    metadataItemKey := tokens[len(tokens) - 1]
    versionString := tokens[len(tokens) - 2]

    version, err := core.GetAppVersionFromString(versionString)
    if err != nil {
        logger.LogError("error parsing editing version from url" +
                        "|request url=" + request.URL.Path +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }

    validVersion, err := metadata_service.Instance().IsVersionValid(version.String(), space)
    if !validVersion {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Invalid version: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }


    err = request.ParseMultipartForm(32 << 20) // limit your max input length!
    if err != nil {
        logger.LogError("error parsing request for uploading metadata item" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    var buffer bytes.Buffer
    uploadedFile, _, err := request.FormFile(metadataItemKey)
    if err != nil {
        logger.LogError("error getting file from request for uploading metadata item" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer func() {
        err = uploadedFile.Close()
    }()

    _, err = io.Copy(&buffer, uploadedFile)
    if err != nil {
        logger.LogError("error copying file contents from request for uploading metadata item" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    fileContents := buffer.String()

    metadataItem, err := metadata_factory.InstantiateMetadataItem(metadataItemKey)
    if err != nil {
        logger.LogError("error instantiating metadata item while processing request for uploading metadata item" +
                        "|metadata item key=" + metadataItemKey +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }
    if space != metadataItem.GetMetadataSpace() {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Wrong metadata space for " + metadataItemKey +
                           ". Expected=" + metadataItem.GetMetadataSpace().String() +
                           ". Got=" + space.String() + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    err = json.Unmarshal([]byte(fileContents), metadataItem)
    if err != nil {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Error deserializing " + metadataItemKey +
                           ". Error=" + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    defer metadata_service.ReleaseInstanceRW()
    err = metadata_service.InstanceRW().SetMetadataItem(metadataItem, version)
    if err != nil {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Error saving " + metadataItemKey +
                           ". Error=" + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    simpleMessagePageObject := AdminSimpleMessageObject{
        AdminPageObject: adminPageObject,
        SimpleMessage: "Successfully saved metadata for " + metadataItemKey,
        BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
    }


    logger.LogInfo("Updated metadata item" +
                   "|metadata space=" + space.String() +
                   "|version=" + version.String() +
                   "|metadata item key=" + metadataItemKey)

    showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
    return
}
