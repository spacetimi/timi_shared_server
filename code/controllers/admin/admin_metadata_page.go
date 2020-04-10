package admin

import (
    "archive/zip"
    "bytes"
    "encoding/json"
    "errors"
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
    "path"
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
const kMetadataRoute_AppDownloadAll = "METADATA_APP_DOWNLOAD_ALL"
const kMetadataRoute_SharedDownloadAll = "METADATA_SHARED_DOWNLOAD_ALL"
const kMetadataRoute_AppUpload = "METADATA_APP_UPLOAD"
const kMetadataRoute_SharedUpload = "METADATA_SHARED_UPLOAD"
const kMetadataRoute_AppUploadAll = "METADATA_APP_UPLOAD_ALL"
const kMetadataRoute_SharedUploadAll = "METADATA_SHARED_UPLOAD_ALL"
const kMetadataRoute_AppRefresh = "METADATA_APP_REFRESH"
const kMetadataRoute_SharedRefresh = "METADATA_SHARED_REFRESH"
const kMetadataRoute_AppCreateNewVersion = "METADATA_APP_CREATE_NEW_VERSION"
const kMetadataRoute_SharedCreateNewVersion = "METADATA_SHARED_CREATE_NEW_VERSION"

var kAdminMetadataRoutes = map[string]string{
    "/admin/metadata$": kMetadataRoute_SelectSpace,
    "/admin/metadata/app$": kMetadataRoute_AppOverview,
    "/admin/metadata/app/setCurrentVersions$": kMetadataRoute_AppOverview,
    "/admin/metadata/app/editVersion/[0-9]+\\.[0-9]+$": kMetadataRoute_AppEditVersion,
    "/admin/metadata/app/view/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_AppViewMetadata,
    "/admin/metadata/app/download/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_AppDownload,
    "/admin/metadata/app/download_all/[0-9]+\\.[0-9]+$": kMetadataRoute_AppDownloadAll,
    "/admin/metadata/app/upload/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_AppUpload,
    "/admin/metadata/app/upload_all/[0-9]+\\.[0-9]+$": kMetadataRoute_AppUploadAll,
    "/admin/metadata/app/refresh$": kMetadataRoute_AppRefresh,
    "/admin/metadata/app/createNewVersion$": kMetadataRoute_AppCreateNewVersion,
    "/admin/metadata/shared$": kMetadataRoute_SharedOverview,
    "/admin/metadata/shared/setCurrentVersions$": kMetadataRoute_SharedOverview,
    "/admin/metadata/shared/editVersion/[0-9]+\\.[0-9]+$": kMetadataRoute_SharedEditVersion,
    "/admin/metadata/shared/view/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_SharedViewMetadata,
    "/admin/metadata/shared/download/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_SharedDownload,
    "/admin/metadata/shared/download_all/[0-9]+\\.[0-9]+$": kMetadataRoute_SharedDownloadAll,
    "/admin/metadata/shared/upload/[0-9]+\\.[0-9]+/.*$": kMetadataRoute_SharedUpload,
    "/admin/metadata/shared/upload_all/[0-9]+\\.[0-9]+$": kMetadataRoute_SharedUploadAll,
    "/admin/metadata/shared/refresh$": kMetadataRoute_SharedRefresh,
    "/admin/metadata/shared/createNewVersion$": kMetadataRoute_SharedCreateNewVersion,
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

    case kMetadataRoute_AppDownloadAll:
        showMetadataDownloadAllPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_AppUpload:
        showMetadataUploadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_AppUploadAll:
        showMetadataUploadAllPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_AppRefresh:
        refreshMetadata(httpResponseWriter, request, metadata_typedefs.METADATA_SPACE_APP)
        return

    case kMetadataRoute_AppCreateNewVersion:
        showMetadataCreateNewVersionPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_APP)
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

    case kMetadataRoute_SharedDownloadAll:
        showMetadataDownloadAllPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case kMetadataRoute_SharedUpload:
        showMetadataUploadPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case kMetadataRoute_SharedUploadAll:
        showMetadataUploadAllPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case kMetadataRoute_SharedRefresh:
        refreshMetadata(httpResponseWriter, request, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    case kMetadataRoute_SharedCreateNewVersion:
        showMetadataCreateNewVersionPage(httpResponseWriter, request, adminPageObject, metadata_typedefs.METADATA_SPACE_SHARED)
        return

    default:
        logger.LogWarning("Unknown metadata route request|request url=" + request.URL.Path)
    }
}

func refreshMetadata(httpResponseWriter http.ResponseWriter, request *http.Request, fromSpace metadata_typedefs.MetadataSpace) {
    metadata_service.RefreshMetadata()

    http.Redirect(httpResponseWriter, request, "/admin/metadata/" + fromSpace.String(), http.StatusSeeOther)
}

func showMetadataSelectPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject) {
    templates, err := template.ParseGlob(config.GetSharedTemplateFilesPath() + "/admin_tool/*")
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

    pageObject.MetadataInfo = MetadataInfo {
        Space: space.String(),
        CurrentVersions: metadata_service.Instance().GetCurrentVersions(space),
        CurrentVersionsCSV: strings.Join(metadata_service.Instance().GetCurrentVersions(space), ","),
        AllVersions: allVersionsSorted,
        IsUpToDate: metadata_service.CheckIfMetadataUpToDate(space),
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
        var messageToShow string

        metadataUpToDate := metadata_service.CheckIfMetadataUpToDate(space)
        if !metadataUpToDate {
            messageToShow = "Metadata not up to date. Please hit Refresh and try again."
            pageObject.HasError = true
            pageObject.ErrorString = "stale metadata for space: " + space.String()

        } else {
            err := updateNewCurrentVersions(space, newCurrentVersionsCSV)
            messageToShow = "Successfully updated current versions."
            if err != nil {
                messageToShow = "Something went wrong updating current versions."
                pageObject.HasError = true
                pageObject.ErrorString = err.Error()
            }
        }

        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: pageObject.AdminPageObject,
            SimpleMessage: messageToShow,
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    templates, err := template.ParseGlob(config.GetSharedTemplateFilesPath() + "/admin_tool/*")
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
        return err
    }

    err = metadata_service.MarkMetadataAsUpdated(space)
    if err != nil {
        return errors.New("error marking metadata as updated: " + err.Error())
    }
    metadata_service.RefreshLastUpdatedTimestamps()

    logger.LogInfo("Updated current versions for metadata space: " + space.String() +
                   " to: " + newCurrentVersionsCSV)

    return nil
}

func showMetadataCreateNewVersionPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace) {

    metadataUpToDate := metadata_service.CheckIfMetadataUpToDate(space)
    if !metadataUpToDate {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Metadata not up to date. Please hit Refresh and try again.",
            BackLinkHref: "/admin/metadata/" + space.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = "stale metadata"

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
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

    newVersionNumberString := request.Form.Get("newVersionNumberString")
    newVersion, err := core.GetAppVersionFromString(newVersionNumberString)
    if err != nil {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Error parsing new version number: " + newVersionNumberString,
            BackLinkHref: "/admin/metadata/" + space.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = err.Error()

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    isCurrent := request.Form.Get("newVersionIsCurrent") == "true"

    defer metadata_service.ReleaseInstanceRW()
    err = metadata_service.InstanceRW().CreateNewVersion(newVersion, space, isCurrent)

    if err != nil {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Error creating new version: " + newVersion.String(),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = err.Error()

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    err = metadata_service.MarkMetadataAsUpdated(space)
    if err != nil {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Error marking metadata as updated",
            BackLinkHref: "/admin/metadata/" + space.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = err.Error()

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }
    metadata_service.RefreshLastUpdatedTimestamps()

    simpleMessagePageObject := AdminSimpleMessageObject{
        AdminPageObject: adminPageObject,
        SimpleMessage: "Successfully created new version: " + newVersion.String(),
        BackLinkHref: "/admin/metadata/" + space.String(),
    }

    showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
    return
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
            SimpleMessage: "Invalid version",
            BackLinkHref: "/admin/metadata/" + space.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = err.Error()

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

    metadataUpToDate := metadata_service.CheckIfMetadataUpToDate(space)
    if !metadataUpToDate {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: pageObject.AdminPageObject,
            SimpleMessage: "Metadata not up to date. Please hit Refresh and try again.",
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    metadataFactories := metadata_factory.GetRegisteredFactories()
    for _, metadataFactory := range metadataFactories {
        metadataItem := metadataFactory.Instantiate()
        metadataManifestItem, err := metadata_service.Instance().GetMetadataManifestItemInVersion(metadataItem.GetKey(), version, space)
        if err != nil || metadataManifestItem == nil {
            pageObject.Items = append(pageObject.Items, AdminMetadataItem{
                Key:metadataItem.GetKey(),
                Hash:"",
                Defined:false,
            })
        } else {
            pageObject.Items = append(pageObject.Items, AdminMetadataItem{
                Key:metadataManifestItem.MetadataKey,
                Hash:metadataManifestItem.Hash,
                Defined:true,
            })
        }
    }

    sort.Slice(pageObject.Items, func(i, j int) bool {
        return pageObject.Items[i].Key < pageObject.Items[j].Key
    })

    templates, err := template.ParseGlob(config.GetSharedTemplateFilesPath() + "/admin_tool/*")
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

    content, err := metadata_service.Instance().GetMetadataItemRawContent(metadataItemKey, version, space)
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

    if viewOnly {
        // Just write out the serialized metadata item
        _, err := fmt.Fprintln(httpResponseWriter, content)
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

func showMetadataDownloadAllPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace) {

    // Parse url for version
    tokens := strings.Split(request.URL.Path, "/")
    if len(tokens) < 2 {
        logger.LogError("malformed request url in metadata download all request" +
                        "|request url=" + request.URL.Path)
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }
    versionString := tokens[len(tokens) - 1]

    version, err := core.GetAppVersionFromString(versionString)
    if err != nil {
        logger.LogError("error parsing version from url" +
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

    manifestItems, err := metadata_service.Instance().GetMetadataManifestItemsInVersion(versionString, space)
    if len(manifestItems) == 0 {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Couldn't get metadata items for version. Error: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    hasError := false
    var errorList []string

    zipBuffer := new(bytes.Buffer)
    zipWriter := zip.NewWriter(zipBuffer)

    for _, manifestItem := range manifestItems {
        metadataItemJson, err := metadata_service.Instance().GetMetadataItemRawContent(manifestItem.MetadataKey, version, space)
        if err != nil {
            hasError = true
            errorList = append(errorList, "error fetching metadata item for key: " + manifestItem.MetadataKey +
                                          ". Error: " + err.Error())
            continue
        }

        fileWriter, err := zipWriter.Create(manifestItem.MetadataKey + ".json")
        if err != nil {
            hasError = true
            errorList = append(errorList, "error adding file to zip archive for metadata item key: " + manifestItem.MetadataKey +
                                          ". Error: " + err.Error())
            continue
        }
        _, err = fileWriter.Write([]byte(metadataItemJson))
        if err != nil {
            hasError = true
            errorList = append(errorList, "error writing metadata item json to zip archive for metadata item key: " + manifestItem.MetadataKey +
                                          ". Error: " + err.Error())
            continue
        }
    }

    err = zipWriter.Close()
    if err != nil {
        hasError = true
        errorList = append(errorList, "error closing zip writer for download all metadata" +
                                      ". Error: " + err.Error())
    }

    if hasError {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Something went wrong while creating downloadable metadata archive. Errors:\n" + strings.Join(errorList, ", "),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    var downloadableFileName string
    if space == metadata_typedefs.METADATA_SPACE_APP {
        downloadableFileName = config.GetAppName() + ".metadata-" + version.String()
    } else {
        downloadableFileName = space.String() + ".metadata-" + version.String()
    }

    // Mark the returned content as downloadable zip to the browser
    httpResponseWriter.Header().Set("Content-Type", "application/zip")
    httpResponseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", downloadableFileName))

    http.ServeContent(httpResponseWriter, request, downloadableFileName, time.Now(), bytes.NewReader(zipBuffer.Bytes()))
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
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    metadataUpToDate := metadata_service.CheckIfMetadataUpToDate(space)
    if !metadataUpToDate {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Metadata not up to date. Please hit Refresh and try again." + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    err = request.ParseMultipartForm(32 << 20)
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

    err = metadata_service.MarkMetadataAsUpdated(space)
    if err != nil {
         simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Error marking metadata as updated: " + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }
    metadata_service.RefreshLastUpdatedTimestamps()

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

func showMetadataUploadAllPage(httpResponseWriter http.ResponseWriter, request *http.Request, adminPageObject AdminPageObject, space metadata_typedefs.MetadataSpace) {

    // Parse url for version
    tokens := strings.Split(request.URL.Path, "/")
    if len(tokens) < 2 {
        logger.LogError("malformed request url in metadata upload all request" +
                        "|request url=" + request.URL.Path)
        httpResponseWriter.WriteHeader(http.StatusNotFound)
        return
    }
    versionString := tokens[len(tokens) - 1]

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
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = "no such version"

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    metadataUpToDate := metadata_service.CheckIfMetadataUpToDate(space)
    if !metadataUpToDate {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Metadata not up to date. Please hit Refresh and try again." + err.Error(),
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = "metadata not up to date"

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    err = request.ParseMultipartForm(32 << 20)
    if err != nil {
        logger.LogError("error parsing request for uploading metadata items" +
                        "|version=" + version.String() +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
        return
    }

    hasErrors := false
    var errorList []string
    var metadataItems []metadata_typedefs.IMetadataItem

    uploadedFiles := request.MultipartForm.File["uploadedFiles"]
    for _, fileHeader := range uploadedFiles {

        metadataItemKey := strings.TrimSuffix(fileHeader.Filename, path.Ext(fileHeader.Filename))
        if metadataItemKey == "" {
            hasErrors = true
            errorList = append(errorList, "error getting metadata item key from file name: " + fileHeader.Filename)
            continue
        }

        file, err := fileHeader.Open()
        if err != nil {
            hasErrors = true
            errorList = append(errorList, "error opening uploaded file: " + fileHeader.Filename)
            continue
        }

        var buffer bytes.Buffer
        _, err = io.Copy(&buffer, file)
        if err != nil {
            hasErrors = true
            errorList = append(errorList, "error copying file contents from file: " + fileHeader.Filename +
                                                  "(error=" + err.Error() + ")")
            continue
        }
        fileContents := buffer.String()

        metadataItem, err := metadata_factory.InstantiateMetadataItem(metadataItemKey)
        if err != nil {
            hasErrors = true
            errorList = append(errorList, "error instantiating metadata item for key: " + metadataItemKey +
                                                  "(error=" + err.Error() + ")")
            continue
        }

        if space != metadataItem.GetMetadataSpace() {
            hasErrors = true
            errorList = append(errorList, "wrong metadata space for metadata item with key: " + metadataItemKey +
                                                  "(expected=" + space.String() +
                                                  ", actual=" + metadataItem.GetMetadataSpace().String() + ")")
            continue
        }

        err = json.Unmarshal([]byte(fileContents), metadataItem)
        if err != nil {
            hasErrors = true
            errorList = append(errorList, "error deserializing metadata item for key: " + metadataItemKey +
                                                  "(error=" + err.Error() + ")")
            continue
        }

        metadataItems = append(metadataItems, metadataItem)

        _ = file.Close()
    }

    if hasErrors {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Something went wrong while preparing to upload metadata. ",
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = ""
        simpleMessagePageObject.MessageExtras = errorList

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    if len(metadataItems) == 0 {
         simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Found no metadata to upload in request",
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
         }
         simpleMessagePageObject.HasError = true
         simpleMessagePageObject.ErrorString = "no metadata found"

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    var metadataItemKeys []string

    defer metadata_service.ReleaseInstanceRW()
    metadataServiceInstance := metadata_service.InstanceRW()
    for _, metadataItem := range metadataItems {
        err = metadataServiceInstance.SetMetadataItem(metadataItem, version)
        if err != nil {
            hasErrors = true
            errorList = append(errorList, "error uploading metadata item for key: " + metadataItem.GetKey() +
                                                  "(error=" + err.Error() + ")")
            continue
        }
        metadataItemKeys = append(metadataItemKeys, metadataItem.GetKey())
    }

    err = metadata_service.MarkMetadataAsUpdated(space)
    if err != nil {
        hasErrors = true
        errorList = append(errorList, "error marking metadata as updated: " +
            "(error=" + err.Error() + ")")
    }
    metadata_service.RefreshLastUpdatedTimestamps()

    if hasErrors {
        simpleMessagePageObject := AdminSimpleMessageObject{
            AdminPageObject: adminPageObject,
            SimpleMessage: "Something went wrong while uploading metadata. ",
            BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
        }
        simpleMessagePageObject.HasError = true
        simpleMessagePageObject.ErrorString = "<br/>"
        simpleMessagePageObject.MessageExtras = errorList

        showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
        return
    }

    simpleMessagePageObject := AdminSimpleMessageObject{
        AdminPageObject: adminPageObject,
        SimpleMessage: "Successfully saved metadata for version: " + version.String() + ". Metadata items uploaded: ",
        MessageExtras: metadataItemKeys,
        BackLinkHref: "/admin/metadata/" + space.String() + "/editVersion/" + version.String(),
    }

    logger.LogInfo("Updated metadata items" +
                   "|metadata space=" + space.String() +
                   "|version=" + version.String() +
                   "|metadata item keys=" + strings.Join(metadataItemKeys, ", "))

    showSimpleMessagePage(httpResponseWriter, request, simpleMessagePageObject)
    return
}

