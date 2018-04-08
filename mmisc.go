package webostv

type ServiceListEntry struct {
	Name    string
	Version int
}

func (tv *Tv) ApiGetServiceList() (list []ServiceListEntry, err error) {
	// "payload":{"services":[{"name":"api","version":1},{"name":"audio","version":1},{"name":"media.controls","version":1},{"name":"media.viewer","version":1},{"name":"pairing","version":1},{"name":"system","version":1},{"name":"system.launcher","version":1},{"name":"system.notifications","version":1},{"name":"tv","version":1},{"name":"webapp","version":2}],"returnValue":true}
	var resp struct {
		Services []ServiceListEntry
	}
	err = tv.RequestResponseParam("ssap://api/getServiceList", nil, &resp)
	return resp.Services, err
}

// TODO ssap://com.webos.service.apiadapter/audio/changeSoundOutput
// TODO ssap://com.webos.service.apiadapter/audio/getSoundOutput // 404 no such service or method
// TODO ssap://com.webos.service.appstatus/getAppStatus // 404 no such service or method
// TODO ssap://com.webos.service.bluetooth/gap/findDevices
// TODO ssap://com.webos.service.bluetooth/gap/getTrustedDevices
// TODO ssap://com.webos.service.bluetooth/gap/isWiFiOnly
// TODO ssap://com.webos.service.bluetooth/gap/removeTrustedDevice
// TODO ssap://com.webos.service.bluetooth/service/connect
// TODO ssap://com.webos.service.bluetooth/service/disconnect
// TODO ssap://com.webos.service.bluetooth/service/getStates
// TODO ssap://com.webos.service.bluetooth/service/subscribeNotifications
// TODO ssap://com.webos.service.connectionmanager/getinfo // 404 no such service or method

func (tv *Tv) ImeDeleteCharacters(count int) (err error) {
	_, err = tv.Request("ssap://com.webos.service.ime/deleteCharacters",
		Payload{"count": count})
	return err
}

func (tv *Tv) ImeInsertText(text string, replace bool) (err error) {
	_, err = tv.Request("ssap://com.webos.service.ime/insertText",
		Payload{
			"text":    text,
			"replace": replace,
		})
	return err
}

// TODO ssap://com.webos.service.ime/registerRemoteKeyboard

// {"type":"response","id":"nlHxhqwT","payload":{"currentWidget":{"autoCapitalizationEnabled":true,"contentType":"text","correctionEnabled":false,"cursorPosition":13,"focus":true,"hasSurroundingText":true,"hiddenText":false,"predictionEnabled":true,"surroundingTextLength":13},"focusChanged":true}}

func (tv *Tv) ImeSendEnterKey() (err error) {
	_, err = tv.Request("ssap://com.webos.service.ime/sendEnterKey", nil)
	return err
}

// TODO ssap://com.webos.service.miracast/close
// TODO ssap://com.webos.service.miracast/getConnectionStatus
// TODO ssap://com.webos.service.miracast/getP2pState
// TODO ssap://com.webos.service.miracast/setUACSettings
// TODO ssap://com.webos.service.miracast/uibc/getUibcKeyEvent

func (tv *Tv) GetPointerInputSocket() (socketPath string, err error) {
	// "payload":{"returnValue":true,"scenario":"mastervolume_tv_speaker","volume":9,"muted":false}
	var resp struct {
		SocketPath string
	}
	err = tv.RequestResponseParam("ssap://com.webos.service.networkinput/getPointerInputSocket", nil, &resp)
	return resp.SocketPath, err
}

func (tv *Tv) SdxGetHttpHeaderForServiceRequest() (resp map[string]string, err error) {
	// {"clearedForDuty":true,"returnValue":true}
	tmpresp, err := tv.Request("ssap://com.webos.service.sdx/getHttpHeaderForServiceRequest", nil)
	if tmpresp != nil {
		resp = make(map[string]string)
	}

	for k, v := range tmpresp {
		if v, ok := v.(string); ok {
			resp[k] = v
		}
	}
	return resp, err
}

func (tv *Tv) SecondscreenGatewayTestSecure() (clearedForDuty bool, err error) {
	// {"clearedForDuty":true,"returnValue":true}
	var resp struct {
		ClearedForDuty bool
	}
	err = tv.RequestResponseParam("ssap://com.webos.service.secondscreen.gateway/test/secure", nil, &resp)
	return resp.ClearedForDuty, err
}

func (tv *Tv) Get3DStatus() (status bool, pattern string, err error) {
	// {"returnValue":true,"status3D":{"status":true,"pattern":"2dto3d"}
	var resp struct {
		Status3D struct {
			Status  bool
			Pattern string
		}
	}
	err = tv.RequestResponseParam("ssap://com.webos.service.tv.display/get3DStatus", nil, &resp)
	return resp.Status3D.Status, resp.Status3D.Pattern, err
}

func (tv *Tv) Set3DOff() (err error) {
	_, err = tv.Request("ssap://com.webos.service.tv.display/set3DOff", nil)
	return err
}

func (tv *Tv) Set3DOn() (err error) {
	_, err = tv.Request("ssap://com.webos.service.tv.display/set3DOn", nil)
	return err
}

// TODO ssap://com.webos.service.tv.keymanager/listInterestingEvents // error KEYMANAGER_ERROR_0001: Required parameter does not exist - subscribe
// TODO ssap://com.webos.service.tvpower/power/getPowerState // 404 no such service or method
// TODO ssap://com.webos.service.tvpower/power/turnOnScreen // 404 no such service or method

func (tv *Tv) GetCurrentTime() (y, m, d, h, min, s int, err error) {
	var resp struct {
		Year   int
		Month  int
		Day    int
		Hour   int
		Minute int
		Second int
	}
	err = tv.RequestResponseParam("ssap://com.webos.service.tv.time/getCurrentTime", nil, &resp)
	return resp.Year, resp.Month, resp.Day, resp.Hour, resp.Minute, resp.Second, err
}

type CurrentSWInformation struct {
	ProductName   string `mapstructure:"product_name"`   // "product_name":"webOS"
	ModelName     string `mapstructure:"model_name"`     // "model_name":"HE_DTV_WT1M_AFAAABAA"
	SwType        string `mapstructure:"sw_type"`        // "sw_type":"FIRMWARE"
	MajorVer      string `mapstructure:"major_ver"`      // "major_ver":"05"
	MinorVer      string `mapstructure:"minor_ver"`      // "minor_ver":"05.35"
	Country       string `mapstructure:"country"`        // "country":"FI"
	DeviceId      string `mapstructure:"device_id"`      // "device_id":"3c:cd:93:7b:91:9e"
	AuthFlag      string `mapstructure:"auth_flag"`      // "auth_flag":"N"
	IgnoreDisable string `mapstructure:"ignore_disable"` // "ignore_disable":"N"
	EcoInfo       string `mapstructure:"eco_info"`       // "eco_info":"01"
	ConfigKey     string `mapstructure:"config_key"`     // "config_key":"00"
	LanguageCode  string `mapstructure:"language_code"`  // "language_code":"en-GB"}
}

func (tv *Tv) GetCurrentSWInformation() (info CurrentSWInformation, err error) {
	err = tv.RequestResponseParam("ssap://com.webos.service.update/getCurrentSWInformation", nil, &info)
	return info, err
}

// TODO ssap://com.webos.service.update/getProgress
// TODO ssap://com.webos.service.update/getStatus
// TODO ssap://com.webos.service.update/startUpdateByRemoteApp
// TODO ssap://config/getConfigs // 404 no such service or method

// TODO ssap://pairing/setPin
// TODO ssap://settings/getSystemSettings // 404 no such service or method
// TODO ssap://system/getHostMessage // 404 no such service or method

type SystemInfo struct {
	// {"type":"response","id":"e8soy4EW","payload":{"features":{"3d":true,"dvr":true},"receiverType":"dvb","modelName":"42LB650V-ZN","returnValue":true}}
	Features     map[string]bool
	ReceiverType string
	ModelName    string
}

func (tv *Tv) SystemGetSystemInfo() (info SystemInfo, err error) {
	err = tv.RequestResponseParam("ssap://system/getSystemInfo", nil, &info)
	return info, err
}

// TODO ssap://system.notifications/createAlert // 404 no such service or method

func (tv *Tv) SystemNotificationsCreateToast(msg string) (toastId string, err error) {
	// {"toastId":"com.webos.service.apiadapter-1522334066285","returnValue":true}
	var resp struct {
		ToastId string
	}
	err = tv.RequestResponseParam("ssap://system.notifications/createToast",
		Payload{"message": msg}, &resp)

	return resp.ToastId, err
}

func (tv *Tv) SystemTurnOff() (err error) {
	_, err = tv.Request("ssap://system/turnOff", nil)
	return err
}

// TODO ssap://timer/getSettings
// TODO ssap://timer/setSettings
// TODO ssap://user/resetUserInfo
// TODO ssap://user/setUserData
// TODO ssap://user/setUserInfo // 404 no such service or method
// TODO ssap://user/setUserSchedule
// TODO ssap://webapp/closeWebApp
// TODO ssap://webapp/connectToApp
// TODO ssap://webapp/isWebAppPinned
// TODO ssap://webapp/launchWebApp
// TODO ssap://webapp/pinWebApp
// TODO ssap://webapp/removePinnedWebApp
