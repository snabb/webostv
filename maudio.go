package webostv

import (
	"github.com/mitchellh/mapstructure"
)

func (tv *Tv) AudioGetMute() (mute bool, err error) {
	// "payload":{"mute":false,"returnValue":true}
	var resp struct {
		Mute bool
	}
	err = tv.RequestResponseParam("ssap://audio/getMute", nil, &resp)
	return resp.Mute, err
}

type AudioStatus struct {
	Scenario string
	Volume   int
	Mute     bool
}

func (tv *Tv) AudioGetStatus() (as AudioStatus, err error) {
	// "payload":{"returnValue":true,"scenario":"mastervolume_tv_speaker","volume":9,"mute":false}
	err = tv.RequestResponseParam("ssap://audio/getStatus", nil, &as)
	return as, err
}

func (tv *Tv) AudioMonitorStatus(process func(as AudioStatus) error, quit <-chan struct{}) error {
	return tv.MonitorStatus("ssap://audio/getStatus", nil, func(payload Payload) (err error) {
		var as AudioStatus
		err = mapstructure.Decode(payload, &as)
		if err == nil {
			err = process(as)
		}
		return err
	}, quit)
}

func (tv *Tv) AudioGetVolume() (scenario string, volume int, muted bool, err error) {
	// "payload":{"returnValue":true,"scenario":"mastervolume_tv_speaker","volume":9,"muted":false}
	var resp struct {
		Scenario string
		Volume   int
		Muted    bool
	}
	err = tv.RequestResponseParam("ssap://audio/getVolume", nil, &resp)
	return resp.Scenario, resp.Volume, resp.Muted, err
}

func (tv *Tv) AudioSetMute(mute bool) (err error) {
	_, err = tv.Request("ssap://audio/setMute",
		Payload{"mute": mute})
	return err
}

func (tv *Tv) AudioSetVolume(volume int) (err error) {
	_, err = tv.Request("ssap://audio/setVolume",
		Payload{"volume": volume})
	return err
}

func (tv *Tv) AudioVolumeDown() (err error) {
	_, err = tv.Request("ssap://audio/volumeDown", nil)
	return err
}

func (tv *Tv) AudioVolumeUp() (err error) {
	_, err = tv.Request("ssap://audio/volumeUp", nil)
	return err
}
