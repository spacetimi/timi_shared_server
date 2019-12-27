package login

import (
	"encoding/json"
	"errors"
	"github.com/spacetimi/timi_shared_server/code/core"
	"strconv"
	"strings"
)

/****************************************/

type LoginRequestParams struct {
	DeviceUID int64
	AppVersionString string
}

func (loginRequestParams *LoginRequestParams) parse() (int64, int64, int64, error) {
	if loginRequestParams.DeviceUID <= 0 {
		return 0, 0, 0, errors.New("No DeviceUID sent")
	}

	if loginRequestParams.AppVersionString == "" {
		return 0, 0, 0, errors.New("No AppVersion sent")
	}
	tokens := strings.Split(loginRequestParams.AppVersionString, ".")
	if len(tokens) != 2 {
		return 0, 0, 0, errors.New("Malformed AppVersion sent|appVersion=" + loginRequestParams.AppVersionString)
	}
	majorVersion, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, 0, 0, errors.New("Invalid major version|error=" + err.Error())
	}
	if majorVersion <= 0 {
		return 0, 0, 0, errors.New("Invalid major version|appVersion=" + loginRequestParams.AppVersionString)
	}
	minorVersion, err := strconv.ParseInt(tokens[1], 10, 32)
	if err != nil {
		return 0, 0, 0, errors.New("Invalid minor version|error=" + err.Error())
	}
	if minorVersion < 0 {
		return 0, 0, 0, errors.New("Invalid minor version|appVersion=" + loginRequestParams.AppVersionString)
	}

	return loginRequestParams.DeviceUID, majorVersion, minorVersion, nil
}

/****************************************/

type LoginRequest struct {
	DeviceUID int64
	AppVersion *core.AppVersion
}

func NewLoginRequest(params *LoginRequestParams) (*LoginRequest, error) {
	lr := LoginRequest{}
	lr.AppVersion = &core.AppVersion{}

	var err error
	lr.DeviceUID, lr.AppVersion.MajorVersion, lr.AppVersion.MinorVersion, err = params.parse()
	if err != nil {
		return nil, err
	}

	return &lr, nil
}

/****************************************/

type LoginResponse struct {
	Success bool
	Body string
	ErrorMessage string
}

func (loginResponse *LoginResponse) String() string {
	j, err := json.Marshal(loginResponse)
	if err != nil {
		return "{}"
	}
	return string(j)
}

/****************************************/

