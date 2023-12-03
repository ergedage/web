package gateway

import "code.byted.org/bytelingoi18n/ai_tutor_api/model"

type CommonArgs struct {
	UserID int64

	DeviceID int64
	AppID    int32
	OsName   model.OS

	Timestamp      int64
	Language       string
	TimeZoneName   string
	TimeZoneOffset int32
}

func newCommonArgs(requester *Requester) *CommonArgs {
	args := &CommonArgs{}

	args.UserID = requester.GetUserID()
	args.AppID = requester.GetAppID()
	args.DeviceID = requester.GetDeviceID()
	args.OsName = requester.GetOSName()
	args.Timestamp = requester.GetTimestamp()
	args.Language = requester.GetLanguage()
	args.TimeZoneName = requester.GetTimeZoneName()
	args.TimeZoneOffset = requester.GetTimeZoneOffset()

	return args
}
