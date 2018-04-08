package webostv

import (
	"github.com/mitchellh/mapstructure"
)

func (tv *Tv) TvChannelDown() (err error) {
	_, err = tv.Request("ssap://tv/channelDown", nil)
	return err
}

func (tv *Tv) TvChannelUp() (err error) {
	_, err = tv.Request("ssap://tv/channelUp", nil)
	return err
}

// TODO ssap://tv/getACRAuthToken // 401 insufficient permissions

type TvCurrentProgramInfo struct {
	ProgramId      string // "programId": "0_31_13105_42559",
	ProgramName    string // "programName": "Keno ja Synttärit",
	Description    string // "description": "Illan Keno-arvonnan [..] visailuohjelma. (2')",
	StartTime      string // "startTime": "2018,04,03,17,58,00"
	EndTime        string // "endTime": "2018,04,03,18,00,00",
	LocalStartTime string // "localStartTime": "2018,04,03,20,58,00",
	LocalEndTime   string // "localEndTime": "2018,04,03,21,00,00",
	ChanelId       string // "channelId": "3_32_24_24_31_13105_0",
	ChannelName    string // "channelName": "Nelonen HD",
	ChannelNumber  string // "channelNumber": "24",
	ChannelMode    string // "channelMode": "Cable",
	Duration       int    // "duration": 120,
}

func (tv *Tv) TvGetChannelCurrentProgramInfo(channelId string) (info TvCurrentProgramInfo, err error) {
	var payload Payload
	if channelId != "" {
		payload = Payload{
			"channelId": channelId,
		}

	}
	err = tv.RequestResponseParam("ssap://tv/getChannelCurrentProgramInfo", payload, &info)
	return info, err
}

type TvChannelGroupId struct {
	Id_  string `mapstructure:"_id"`              // "_id": "7d81",
	Id   int    `mapstructure:"channelGroupId"`   // "channelGroupId": 1,
	Name string `mapstructure:"channelGroupName"` // "channelGroupName": "DTV"
}

type TvChannel struct {
	ChannelId       string // "channelId": "3_3_23_23_17_3291_0",
	ChannelMajMinNo string // "channelMajMinNo": "04-00023-000-003",
	ChannelName     string // "channelName": "MTV3 HD",
	ChannelNumber   string // "channelNumber": "23",
	ChannelType     string // "channelType": "Cable Digital TV",
	ChannelTypeId   int    // "channelTypeId": 4,
	ChannelMode     string // "channelMode": "Cable",
	ChannelModeId   int    // "channelModeId": 1,
	SignalChannelId string // "signalChannelId": "17_3291_0",
	ProgramId       string // "programId": "17_3291_0",
	FavoriteGroup   string // "favoriteGroup": "",
	SatelliteName   string // "satelliteName": " ",
	Frequency       int    // "Frequency": 130000,
	Bandwidth       int    // "Bandwidth": 0,
	SourceIndex     int    // "sourceIndex": 3,
	ServiceType     int    // "serviceType": 25,
	ShortCut        int    // "shortCut": 0,
	Handle          int    // "Handle": 0,
	ONID            int    // "ONID": 0,
	SVCID           int    // "SVCID": 3291,
	TSID            int    // "TSID": 17,
	ConfigurationId int    // "configurationId": 0,
	MajorNumber     int    // "majorNumber": 23,
	MinorNumber     int    // "minorNumber": 23,
	PhysicalNumber  int    // "physicalNumber": 3,
	ATV             bool   // "ATV": false,
	DTV             bool   // "DTV": true,
	Data            bool   // "Data": false,
	HDTV            bool   // "HDTV": true,
	Invisible       bool   // "Invisible": false,
	Numeric         bool   // "Numeric": false,
	PrimaryCh       bool   // "PrimaryCh": true,
	Radio           bool   // "Radio": false,
	TV              bool   // "TV": true,
	Descrambled     bool   // "descrambled": true,
	FineTuned       bool   // "fineTuned": false,
	Locked          bool   // "locked": false,
	SatelliteLcn    bool   // "satelliteLcn": false,
	Scrambled       bool   // "scrambled": false,
	Skipped         bool   // "skipped": false,
	SpecialService  bool   // "specialService": false
	GroupIdList     []TvChannelGroupId
	// "CASystemIDList": {}, // ???
	// "CASystemIDListCount": 0, // ???
}

func (tv *Tv) TvGetChannelList() (list []TvChannel, err error) {
	var resp struct {
		ChannelList []TvChannel
	}
	err = tv.RequestResponseParam("ssap://tv/getChannelList", nil, &resp)
	return resp.ChannelList, err
}

type TvProgramRating struct {
	Id           string `mapstructure:"_id"` // "_id": "157ac",
	RatingString string // "ratingString": "",
	RatingValue  int    // "ratingValue": 3,
	Region       int    // "region": 4606286
}

type TvProgram struct {
	ProgramId       string            // "programId": "0_31_13105_42559",
	ProgramName     string            // "programName": "Keno ja Synttärit",
	Description     string            // "description": "Illan Keno-arvonnan [..] visailuohjelma. (2')",
	StartTime       string            // "startTime": "2018,04,03,17,58,00"
	EndTime         string            // "endTime": "2018,04,03,18,00,00",
	LocalStartTime  string            // "localStartTime": "2018,04,03,20,58,00",
	LocalEndTime    string            // "localEndTime": "2018,04,03,21,00,00",
	DSTStartTime    string            // "DSTStartTime": "2018,04,03,20,58,00",
	DSTEndTime      string            // "DSTEndTime": "2018,04,03,21,00,00",
	SignalChannelId string            // "signalChannelId": "31_13105_0",
	Duration        int               // "duration": 120,
	IsPresent       bool              // "isPresent": false,
	Rating          []TvProgramRating // "rating": [
}

func (tv *Tv) TvGetChannelProgramInfo(channelId string) (channel TvChannel, programlist []TvProgram, err error) {
	var resp struct {
		Channel     TvChannel
		ProgramList []TvProgram
	}
	var payload Payload
	if channelId != "" {
		payload = Payload{
			"channelId": channelId,
		}

	}
	err = tv.RequestResponseParam("ssap://tv/getChannelProgramInfo", payload, &resp)
	return resp.Channel, resp.ProgramList, err
}

type TvCurrentChannel struct {
	ChannelId       string // "channelId":"3_32_24_24_31_13105_0"
	SignalChannelId string // "signalChannelId":"31_13105_0"
	ChannelModeId   int    // "channelModeId":1
	ChannelModeName string // "channelModeName":"Cable"
	ChannelTypeId   int    // "channelTypeId":4
	ChannelTypeName string // "channelTypeName":"Cable Digital TV"
	ChannelNumber   string // "channelNumber":"24"
	ChannelName     string // "channelName":"Nelonen HD"
	PhysicalNumber  int    // "physicalNumber":32
	IsSkipped       bool   // "isSkipped":false
	IsLocked        bool   // "isLocked":false
	IsDescrambled   bool   // "isDescrambled":true
	IsScrambled     bool   // "isScrambled":false
	IsFineTuned     bool   // "isFineTuned":false
	IsInvisible     bool   // "isInvisible":false
	// "favoriteGroup":null
	// "hybridtvType":null
	// "dualChannel":{"dualChannelId":null
	// "dualChannelTypeId":null
	// "dualChannelTypeName":null
	// "dualChannelNumber":null
}

func (tv *Tv) TvGetCurrentChannel() (cur TvCurrentChannel, err error) {
	err = tv.RequestResponseParam("ssap://tv/getCurrentChannel", nil, &cur)
	return cur, err
}

func (tv *Tv) TvMonitorCurrentChannel(process func(cur TvCurrentChannel) error, quit <-chan struct{}) error {
	return tv.MonitorStatus("ssap://tv/getCurrentChannel", nil, func(payload Payload) (err error) {
		var cur TvCurrentChannel
		err = mapstructure.Decode(payload, &cur)
		if err == nil {
			err = process(cur)
		}
		return err
	}, quit)
}

type TvExternalInput struct {
	Id              string // "id": "SCART_1",
	Label           string // "label": "AV1",
	Port            int    // "port": 1,
	AppId           string // "appId": "com.webos.app.externalinput.scart",
	Icon            string // "icon": "http://lgsmarttv.lan:3000/resources/d8dd219500f8c1604e548d980c0f60979be5b5a5/scart.png",
	CurrentTVStatus string // "currentTVStatus": "",
	Modified        bool   // "modified": false,
	Autoav          bool   // "autoav": false,
	Connected       bool   // "connected": false,
	Favorite        bool   // "favorite": false
	// "subList": [],
	// "subCount": 0,
}

func (tv *Tv) TvGetExternalInputList() (list []TvExternalInput, err error) {
	var resp struct {
		Devices []TvExternalInput
	}

	err = tv.RequestResponseParam("ssap://tv/getExternalInputList", nil, &resp)
	return resp.Devices, err
}

func (tv *Tv) TvOpenChannelId(channelId string) (err error) {
	_, err = tv.Request("ssap://tv/openChannel",
		Payload{"channelId": channelId})
	return err
}

func (tv *Tv) TvOpenChannelNumber(channelNumber string) (err error) {
	_, err = tv.Request("ssap://tv/openChannel",
		Payload{"channelNumber": channelNumber})
	return err
}

func (tv *Tv) TvSwitchInput(inputId string) (err error) {
	_, err = tv.Request("ssap://tv/switchInput", Payload{"inputId": inputId})
	return err
}
