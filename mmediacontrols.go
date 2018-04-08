package webostv

func (tv *Tv) MediaControlsFastForward() (err error) {
	_, err = tv.Request("ssap://media.controls/fastForward", nil)
	return err
}

func (tv *Tv) MediaControlsPause() (err error) {
	_, err = tv.Request("ssap://media.controls/pause", nil)
	return err
}

func (tv *Tv) MediaControlsPlay() (err error) {
	_, err = tv.Request("ssap://media.controls/play", nil)
	return err
}

func (tv *Tv) MediaControlsRewind() (err error) {
	_, err = tv.Request("ssap://media.controls/rewind", nil)
	return err
}

func (tv *Tv) MediaControlsStop() (err error) {
	_, err = tv.Request("ssap://media.controls/stop", nil)
	return err
}
