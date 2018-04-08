package webostv

import (
	"github.com/mitchellh/mapstructure"
)

type App struct {
	Id                         string      // "id": "com.webos.app.discovery",
	Title                      string      // "title": "LG Store",
	Version                    string      // "version": "1.0.19",
	Vendor                     string      // "vendor": "LGE",
	FolderPath                 string      // "folderPath": "/mnt/otncabi/usr/palm/applications/com.webos.app.discovery",
	DefaultWindowType          string      // "defaultWindowType": "card",
	MediumIcon                 string      // "mediumIcon": "lgstore_80x80.png",
	Miniicon                   string      // "miniicon": "http://lgsmarttv.lan:3000/r[..]e/miniicon.png",
	RequestedWindowOrientation string      // "requestedWindowOrientation": "",
	LargeIcon                  string      // "largeIcon": "lgstore_130x130.png",
	Icon                       string      // "icon": "http://lgsmarttv.lan:3000/r[..]9/lgstore_80x80.png",
	Category                   string      // "category": "",
	TrustLevel                 string      // "trustLevel": "default",
	SplashBackground           string      // "splashBackground": "lgstore_splash.png",
	DeeplinkingParams          string      // "deeplinkingParams": "{\"contentTarget\":\"$CONTENTID\"}",
	RequiredEULA               string      // "requiredEULA": "generalTerms",
	Main                       string      // "main": "index.html",
	Type                       string      // "type": "web",
	BgImage                    string      // "bgImage": "lgstore_preview.png",
	IconColor                  string      // "iconColor": "#cf0652",
	ImageForRecents            string      // "imageForRecents": "RECENTS.png",
	Resolution                 string      // "resolution": "1280x720",
	BgColor                    string      // "bgColor": "#8e191b",
	ContainerCSS               string      // "containerCSS": "build/app1.css",
	EnyoVersion                string      // "enyoVersion": "2.3.0",
	ContainerJS                string      // "containerJS": "build/app1.js",
	Splashicon                 string      // "splashicon": "http://lgsmarttv.lan:3000/r[..]2/splash.png",
	Appsize                    int         // "appsize": 0,
	HardwareFeaturesNeeded     int         // "hardwareFeaturesNeeded": 0,
	Age                        int         // "age": 0,
	BinId                      int         // "binId": 361092,
	RequiredMemory             int         // "requiredMemory": 160,
	Lockable                   bool        // "lockable": true,
	Transparent                bool        // "transparent": false,
	CheckUpdateOnLaunch        bool        // "checkUpdateOnLaunch": true,
	Launchinnewgroup           bool        // "launchinnewgroup": false,
	HandlesRelaunch            bool        // "handlesRelaunch": false,
	Inspectable                bool        // "inspectable": false,
	InAppSetting               bool        // "inAppSetting": false,
	PrivilegedJail             bool        // "privilegedJail": false,
	Visible                    bool        // "visible": true,
	NoWindow                   bool        // "noWindow": false,
	Removable                  bool        // "removable": true,
	DisableBackHistoryAPI      bool        // "disableBackHistoryAPI": true
	InternalInstallationOnly   bool        // "internalInstallationOnly": true,
	NoSplashOnLaunch           bool        // "noSplashOnLaunch": true,
	CustomPlugin               bool        // "customPlugin": true,
	Hidden                     bool        // "hidden": true,
	UIRevision                 interface{} // "uiRevision": 2, // "uiRevision": "2",
	MimeTypes                  []struct {  // "mimeTypes": [
		Mime      string // "mime": "application/vnd.lge.appstore"
		Extension string // "extension": "html",
		Scheme    string // "scheme": "https"
		Stream    bool   // "stream": true,
	}
	Class struct { // "class": {
		Hidden bool // "hidden": true
	}
	OnDeviceSource   map[string]string      // "onDeviceSource": {
	VendorExtension  map[string]interface{} // // "vendorExtension": {
	BootLaunchParams struct {               // "bootLaunchParams": {
		BGMode string // "BGMode": "1"
		Boot   bool   // "boot": true
	}
	// "windowGroup": {
	// "keyFilterTable": [
}

func (tv *Tv) ApplicationManagerGetAppInfo(id string) (info App, err error) {
	// {"type":"response","id":"YokZ11MX","payload":{"mute":false,"returnValue":true}}
	var resp struct {
		AppInfo App
		AppId   string
	}
	err = tv.RequestResponseParam("ssap://com.webos.applicationManager/getAppInfo",
		Payload{"id": id},
		&resp)

	return resp.AppInfo, err
}

type ForegroundAppInfo struct {
	AppId     string
	WindowId  string
	ProcessId string
}

func (i *ForegroundAppInfo) IsLiveTv() bool {
	return i.AppId == "com.webos.app.livetv"
}

func (tv *Tv) ApplicationManagerGetForegroundAppInfo() (info ForegroundAppInfo, err error) {
	err = tv.RequestResponseParam("ssap://com.webos.applicationManager/getForegroundAppInfo", nil, &info)
	// {"type":"response","id":"CyvqdwSl","payload":{"appId":"com.webos.app.hdmi2","returnValue":true,"windowId":"","processId":"n-1059"}}
	return info, err
}

func (tv *Tv) ApplicationManagerMonitorForegroundAppInfo(process func(info ForegroundAppInfo) error, quit <-chan struct{}) error {
	return tv.MonitorStatus("ssap://com.webos.applicationManager/getForegroundAppInfo", nil, func(payload Payload) (err error) {
		var info ForegroundAppInfo
		err = mapstructure.Decode(payload, &info)
		if err == nil {
			err = process(info)
		}
		return err
	}, quit)
}

func (tv *Tv) ApplicationManagerLaunch(id string, params Payload) (processId string, err error) {
	// {"type":"response","id":"IkVU1ZGv","payload":{"returnValue":true,"processId":"1001"}}
	p := make(Payload)
	p["id"] = id
	for k, v := range params {
		p[k] = v
	}
	var resp struct {
		ProcessId string
	}
	err = tv.RequestResponseParam("ssap://com.webos.applicationManager/launch", p, &resp)

	return resp.ProcessId, err
}

func (tv *Tv) ApplicationManagerListApps() (list []App, err error) {
	var resp struct {
		Apps []App
	}
	err = tv.RequestResponseParam("ssap://com.webos.applicationManager/listApps", nil, &resp)
	return resp.Apps, err
}

type LaunchPoint struct {
	Removable       bool              // "removable": false,
	LargeIcon       string            // "largeIcon": "/mnt/otncabi/usr/palm/applications/com.webos.app.discovery/lgstore_130x130.png",
	Vendor          string            // "vendor": "LGE",
	Id              string            // "id": "com.webos.app.discovery",
	Title           string            // "title": "LG Store",
	BgColor         string            // "bgColor": "#8e191b",
	VendorURL       string            // "vendorUrl": "",
	IconColor       string            // "iconColor": "#4b4b4b",
	AppDescription  string            // "appDescription": "",
	Params          map[string]string // "params": { //           "deviceId": "HDMI_2"
	Version         string            // "version": "1.0.19",
	BgImage         string            // "bgImage": "/mnt/otncabi/usr/palm/applications/com.webos.app.discovery/lgstore_preview.png",
	Icon            string            // "icon": "http://lgsmarttv.lan:3000/resources/e1a2afa2ee2c03b7e7c89247d3425a8af8657e5d/lgstore_80x80.png",
	LaunchPointId   string            // "launchPointId": "com.webos.app.discovery_default",
	ImageForRecents string            // "imageForRecents": "/media/cryptofs/apps/usr/palm/applications/netflix/RECENTS.png"
}

type CaseDetail struct {
	Code               int    // "code": 1,
	ServiceCountryCode string // "serviceCountryCode": "FIN",
	// "change": [], // ???
	LocaleCode           string // "localeCode": "en-GB",
	BroadcastCountryCode string // "broadcastCountryCode": "FIN"
}

func (tv *Tv) ApplicationManagerListLaunchPoints() (launchPoints []LaunchPoint, caseDetail CaseDetail, err error) {
	var resp struct {
		Subscribed   bool
		LaunchPoints []LaunchPoint
		CaseDetail   CaseDetail
	}
	err = tv.RequestResponseParam("ssap://com.webos.applicationManager/listLaunchPoints", nil, &resp)
	return resp.LaunchPoints, resp.CaseDetail, err
}
func (tv *Tv) SystemLauncherClose(sessionId string) (err error) {
	_, err = tv.Request("ssap://system.launcher/close",
		Payload{"sessionId": sessionId})
	return err
}

func (tv *Tv) SystemLauncherGetAppState(sessionId string) (running, visible bool, err error) {
	// {"running":true,"visible":true,"returnValue":true}
	var resp struct {
		Running bool
		Visible bool
	}
	err = tv.RequestResponseParam("ssap://system.launcher/getAppState",
		Payload{"sessionId": sessionId}, &resp)

	return resp.Running, resp.Visible, err
}

func (tv *Tv) SystemLauncherLaunch(appId string, params Payload) (sessionId string, err error) {
	// {"returnValue":true,"sessionId":"eW91dHViZS5sZWFuYmFjay52NDp1bmRlZmluZWQ="}
	p := make(Payload)
	p["id"] = appId
	for k, v := range params {
		p[k] = v
	}
	var resp struct {
		SessionId string
	}
	err = tv.RequestResponseParam("ssap://system.launcher/launch", p, &resp)

	return resp.SessionId, err
}

func (tv *Tv) SystemLauncherOpen(url string) (appId, sessionId string, err error) {
	// {"returnValue":true,"id":"com.webos.app.browser","sessionId":"Y29tLndlYm9zLmFwcC5icm93c2VyOnVuZGVmaW5lZA=="}
	var resp struct {
		Id        string
		SessionId string
	}
	err = tv.RequestResponseParam("ssap://system.launcher/open",
		Payload{"target": url}, &resp)

	return resp.Id, resp.SessionId, err
}

/*
func (tv *Tv) LaunchBrowser(url string) (sessionId string, err error) {
	var p Payload
	if url != "" {
		p = Payload{"target": url}
	}
	return tv.SystemLauncherLaunch("com.webos.app.browser", p)
}
*/

func (tv *Tv) LaunchYoutube(videoId string) (sessionId string, err error) {
	var p Payload
	if videoId != "" {
		p = Payload{
			"params": Payload{
				"contentTarget": "http://www.youtube.com/tv?v=" + videoId,
			},
		}
	}
	return tv.SystemLauncherLaunch("youtube.leanback.v4", p)
}

func (tv *Tv) LaunchNetflix(contentId string) (sessionId string, err error) {
	var p Payload
	if contentId != "" {
		p = Payload{
			"contentId": "m=http%3A%2F%2Fapi.netflix.com%2Fcatalog%2Ftitles%2Fmovies%2F" + contentId + "&source_type=4",
		}
	}
	return tv.SystemLauncherLaunch("netflix", p)
}
