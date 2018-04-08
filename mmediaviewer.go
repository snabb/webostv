package webostv

func (tv *Tv) MediaViewerClose(sessionId string) (err error) {
	_, err = tv.Request("ssap://media.viewer/close",
		Payload{"sessionId": sessionId})
	return err
}

func (tv *Tv) MediaViewerOpen(url, title, description, mimeType, iconSrc string, loop bool) (appId, sessionId string, err error) {
	// {"returnValue":true,"id":"com.webos.app.tvsimpleviewer","sessionId":"Y29tLndlYm9zLmFwcC50dnNpbXBsZXZpZXdlcjp1bmRlZmluZWQ="}

	p := make(Payload)
	p["target"] = url
	if title != "" {
		p["title"] = title
	}
	if description != "" {
		p["description"] = description
	}
	if mimeType != "" {
		p["mimeType"] = mimeType
	}
	if iconSrc != "" {
		p["iconSrc"] = iconSrc
	}
	if loop {
		p["loop"] = loop
	}
	var resp struct {
		Id        string
		SessionId string
	}
	err = tv.RequestResponseParam("ssap://media.viewer/open", p, &resp)

	return resp.Id, resp.SessionId, err
}
