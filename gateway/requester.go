package gateway

import (
	"code.byted.org/bytelingoi18n/ai_tutor_api/model"
	"code.byted.org/gopkg/env"
	"code.byted.org/gopkg/logs"
	"code.byted.org/lang/gg/gslice"
	"context"
	"strconv"
	"strings"
	"time"

	"code.byted.org/bytelingoi18n/ai_tutor_api/gateway/session"

	"code.byted.org/middleware/hertz/pkg/app"
)

func createRequester(ctx context.Context, request *app.RequestContext) *Requester {
	hdlStartTime := time.Now()
	requester := &Requester{
		RPCtx:        ctx,
		request:      request,
		hdlStartTime: hdlStartTime,
	}
	return requester
}

type Requester struct {
	RPCtx   context.Context
	request *app.RequestContext

	hdlStartTime time.Time // 服务端开始处理的时间
}

func (requester *Requester) GetHdlStartTime() time.Time {
	return requester.hdlStartTime
}

func (requester *Requester) GetModule() string {
	path := requester.Path()
	// URL 遵守规范：https://{domain}/{application}/{module}/{func}/{version}?query=xxx
	parts := strings.Split(path, "/")
	if len(parts) <= 2 {
		return "unknown"
	}
	return parts[2]
}

func (requester *Requester) Path() string {
	return string(requester.request.Path())
}

func (requester *Requester) GetOSName() model.OS {
	device := requester.GetDevicePlatform()
	os := model.ParseOS(device)
	if os == model.OSUnknown {
		logs.CtxWarn(requester.RPCtx, "unknown OS `%s`", device)
	}
	return os
}

func (requester *Requester) GetAppID() int32 {
	appID, _ := strconv.ParseInt(requester.request.Query("app_id"), 10, 32)
	return int32(appID)
}

func (requester *Requester) GetCommonArgs() *CommonArgs {
	commonArgs := newCommonArgs(requester)
	return commonArgs
}

func (requester *Requester) IsLogin() bool {
	if env.IsBoe() || env.IsPPE() || env.IsTesting() {
		loginStr := requester.GetHeader("Aitutor-Mock-Login")
		userIdStr := requester.GetHeader("Aitutor-Mock-Uid")
		if userId, err := strconv.ParseInt(userIdStr, 10, 64); len(loginStr) > 0 && err == nil {
			return loginStr == "true" && userId > 0
		}
	}
	return session.IsLogin(requester.RPCtx, requester.request)
}

func (requester *Requester) GetUserID() int64 {
	if env.IsBoe() || env.IsPPE() || env.IsTesting() {
		userIdStr := requester.GetHeader("Aitutor-Mock-Uid")
		if userId, err := strconv.ParseInt(userIdStr, 10, 64); err == nil {
			return userId
		}
	}
	if requester.IsLogin() {
		return session.GetUserId(requester.RPCtx, requester.request)
	}
	return session.GetOdinId(requester.RPCtx, requester.request)
}

func (requester *Requester) GetDeviceID() int64 {
	deviceID, _ := strconv.ParseInt(requester.request.Query("device_id"), 10, 64)
	return deviceID
}

func (requester *Requester) GetDeviceType() string {
	return requester.request.Query("device_type")
}

func (requester *Requester) GetClientIP() string {
	// from https://bytedance.feishu.cn/wiki/wikcnwUcsCIqKNCFPUBleYDBkjt
	ip := strings.TrimSpace(requester.GetHeader("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(requester.GetHeader("X-Forwarded-For"))
	ips := strings.Split(ip, ",")
	ips = gslice.Map(ips, func(v string) string {
		return strings.TrimSpace(v)
	})
	ips = gslice.Filter(ips, func(v string) bool {
		return v != ""
	})
	if len(ips) != 0 {
		return ips[0]
	}

	ip = strings.TrimSpace(requester.GetHeader("x-alicdn-da-via"))
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(requester.GetHeader("X-Appengine-Remote-Addr"))
	if ip != "" {
		return ip
	}
	logs.CtxWarn(requester.RPCtx, "client ip is null")
	return ""
}

func (requester *Requester) GetDevicePlatform() string {
	return requester.request.Query("device_platform")
}

func (requester *Requester) GetChannel() string {
	return requester.request.Query("channel")
}

func (requester *Requester) GetUncheckAppVersion() string {
	return requester.request.Query("app_version")
}

func (requester *Requester) GetAppVersion() model.AITutorVersion {
	versionStr := requester.GetUncheckAppVersion()
	version, err := model.ParseAITutorVersion(versionStr)
	if err != nil {
		logs.CtxErrorKvs(requester.RPCtx, "ParseAITutorVersion err", err)
	}
	return version
}

func (requester *Requester) getContentType() string {
	return string(requester.request.ContentType())
}

type bodycodec uint8

const (
	bodyUnsupport = iota
	bodyJSON
	bodyProtobuf
)

func (requester *Requester) getBodyCodec() bodycodec {
	ct := requester.getContentType()
	idx := strings.Index(ct, ";")
	if idx != -1 {
		ct = strings.TrimRight(ct[:idx], " ")
	}
	switch ct {
	case "application/json":
		return bodyJSON
	case "application/x-protobuf":
		return bodyProtobuf
	//case "application/x-www-form-urlencoded", "multipart/form-data":
	//	return bodyForm
	default:
		return bodyUnsupport
	}
}
func (requester *Requester) GetUseJsonPb() bool {
	val := requester.GetHeader("Use-Jsonpb")
	if len(val) == 0 {
		logs.CtxInfo(requester.RPCtx, "Use-Jsonpb is null")
		return false
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		logs.CtxWarn(requester.RPCtx, "Use-Jsonpb ParseInt error, value=%+v, err=%+v", val, err)
		return false
	}
	return res == 1
}

func (requester *Requester) GetHeader(name string) string {
	return string(requester.request.GetHeader(name))
}

func (requester *Requester) GetHeaders() map[string]string {
	headers := make(map[string]string)
	requester.request.VisitAllHeaders(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

func (requester *Requester) GetQueries() map[string]string {
	queries := make(map[string]string)
	requester.request.VisitAllQueryArgs(func(key, value []byte) {
		queries[string(key)] = string(value)
	})
	return queries
}

func (requester *Requester) GetTimestamp() int64 {
	val := requester.GetHeader("Aitutor-Req-Timestamp")
	if len(val) == 0 {
		logs.CtxWarn(requester.RPCtx, "Aitutor-Req-Timestamp is null")
		return 0
	}
	ts, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		logs.CtxWarn(requester.RPCtx, "Aitutor-Req-Timestamp ParseInt error, value=%+v, err=%+v", val, err)
		return 0
	}
	return ts
}

func (requester *Requester) GetLanguage() string {
	val := requester.GetHeader("Aitutor-App-Language")
	if len(val) == 0 {
		logs.CtxWarn(requester.RPCtx, "Aitutor-Req-Language is null")
		return ""
	}
	return val
}

func (requester *Requester) GetTimeZoneName() string {
	val := requester.GetHeader("Aitutor-Req-Timezone-Name")
	if len(val) == 0 {
		logs.CtxWarn(requester.RPCtx, "Aitutor-Req-Timezone-Name is null")
		return ""
	}
	return val
}

func (requester *Requester) GetTimeZoneOffset() int32 {
	val := requester.GetHeader("Aitutor-Req-Timezone-Offset")
	if len(val) == 0 {
		logs.CtxWarn(requester.RPCtx, "Aitutor-Req-Timezone-Offset is null")
		return 0
	}
	offset, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		logs.CtxWarn(requester.RPCtx, "Aitutor-Req-Timezone-Offset ParseInt error, value=%+v, err=%+v", val, err)
		return 0
	}
	return int32(offset)
}
