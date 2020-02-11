package login

import (
	"encoding/json"
	"fmt"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
	"github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"net/http"
)

func HandleLogin(httpResponseWriter http.ResponseWriter, request *http.Request) {
	loginResponse := processLoginRequest(request)
	_,err := fmt.Fprintln(httpResponseWriter, &loginResponse)
	if err != nil {
		logger.LogError("Something went wrong sending login response: " + loginResponse.String())
	}
}

func processLoginRequest(request *http.Request) LoginResponse {
	err := request.ParseForm()
	if err != nil {
		return LoginResponse{Success:false, ErrorMessage:"Badly formed login request: " + err.Error()}
	}

	loginParams_json := request.Form.Get("login_params")
	if len(loginParams_json) <= 0 {
		return LoginResponse{Success:false, ErrorMessage:"No login params provided"}
	}

	loginParams := LoginRequestParams{}
	err = json.Unmarshal([]byte(loginParams_json), &loginParams)
	if err != nil {
		return LoginResponse{Success:false, ErrorMessage:"Unable to deserialize login params json: " + err.Error()}
	}

	loginRequest, err := NewLoginRequest(&loginParams)
	if err != nil {
		return LoginResponse{Success:false, ErrorMessage:"Unable to construct login request: " + err.Error()}
	}

	// TODO: krisa: Use Appversion in GetMetaDataItem()
	// and also insert AppVersion and/or LoginRequest into the context from this point downward?

	mdUpToDate, err := metadata_service.Instance().IsMetadataHashUpToDate("MetadataTest", "f5499079646c8e6a3fded46d2903516d", metadata_typedefs.METADATA_SPACE_APP, loginRequest.AppVersion)
	if err != nil {
		logger.LogError("failed to find if metadata up to date|error=" + err.Error())
	}

	if mdUpToDate {
		logger.LogInfo("metadata up to date")
	} else {
		logger.LogInfo("metadata not up to date")
	}

	return LoginResponse{Success:true, ErrorMessage:"", Body:"Successfully logged in to App: " + config.GetAppName()}
}
