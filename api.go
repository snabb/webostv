package webostv

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	Timeout               = time.Second * 5
	RegisterTimeout       = time.Second * 30
	ErrTimeout            = errors.New("timeout")
	ErrNoResponse         = errors.New("no response")
	ErrRegistrationFailed = errors.New("registration failed")
)

type Tv struct {
	Address      string
	ws           *websocket.Conn
	wsWriteMutex sync.Mutex
	respCh       map[string]chan<- Msg
	respChMutex  sync.Mutex
	debugFunc    func(string)
}

type Dialer struct {
	DisableTLS      bool
	WebsocketDialer *websocket.Dialer
}

var DefaultDialer = Dialer{
	DisableTLS: false,
	WebsocketDialer: &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // TV has a self signed certificate
		},
		NetDial: (&net.Dialer{
			Timeout:   time.Second * 5,
			KeepAlive: time.Second * 30, // ensure we notice if the TV goes away
		}).Dial,
	},
}

func (dialer *Dialer) Dial(address string) (tv *Tv, err error) {
	var url string
	if dialer.DisableTLS {
		url = "ws://" + address + ":3000"
	} else {
		url = "wss://" + address + ":3001"
	}
	wsDialer := dialer.WebsocketDialer
	if wsDialer == nil {
		wsDialer = websocket.DefaultDialer
	}
	ws, resp, err := wsDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &Tv{
		Address: address,
		ws:      ws,
	}, nil
}

func (tv *Tv) debug(str string, buf []byte) {
	if tv.debugFunc != nil {
		if buf != nil {
			tv.debugFunc(str + string(buf))
		} else {
			tv.debugFunc(str)
		}
	}
}

func (tv *Tv) SetDebug(debugFunc func(string)) {
	tv.debugFunc = debugFunc
}

func helloPayload() Payload {
	return Payload{
		"forcePairing": false,
		"pairingType":  "PROMPT",
		"manifest": map[string]interface{}{
			"manifestVersion": 1,
			"appVersion":      "1.1",
			"signed": map[string]interface{}{
				"created":  "20140509",
				"appId":    "com.lge.test",
				"vendorId": "com.lge",
				"localizedAppNames": map[string]string{
					"":       "LG Remote App",
					"ko-KR":  "리모컨 앱",
					"zxx-XX": "ЛГ Rэмotэ AПП",
				},
				"localizedVendorNames": map[string]string{
					"": "LG Electronics",
				},
				"permissions": []string{
					"TEST_SECURE",
					"CONTROL_INPUT_TEXT",
					"CONTROL_MOUSE_AND_KEYBOARD",
					"READ_INSTALLED_APPS",
					"READ_LGE_SDX",
					"READ_NOTIFICATIONS",
					"SEARCH",
					"WRITE_SETTINGS",
					"WRITE_NOTIFICATION_ALERT",
					"CONTROL_POWER",
					"READ_CURRENT_CHANNEL",
					"READ_RUNNING_APPS",
					"READ_UPDATE_INFO",
					"UPDATE_FROM_REMOTE_APP",
					"READ_LGE_TV_INPUT_EVENTS",
					"READ_TV_CURRENT_TIME",
				},
				"serial": "2f930e2d2cfe083771f68e4fe7bb07",
			},
			"permissions": []string{
				"LAUNCH",
				"LAUNCH_WEBAPP",
				"APP_TO_APP",
				"CLOSE",
				"TEST_OPEN",
				"TEST_PROTECTED",
				"CONTROL_AUDIO",
				"CONTROL_DISPLAY",
				"CONTROL_INPUT_JOYSTICK",
				"CONTROL_INPUT_MEDIA_RECORDING",
				"CONTROL_INPUT_MEDIA_PLAYBACK",
				"CONTROL_INPUT_TV",
				"CONTROL_POWER",
				"READ_APP_STATUS",
				"READ_CURRENT_CHANNEL",
				"READ_INPUT_DEVICE_LIST",
				"READ_NETWORK_STATE",
				"READ_RUNNING_APPS",
				"READ_TV_CHANNEL_LIST",
				"WRITE_NOTIFICATION_TOAST",
				"READ_POWER_STATE",
				"READ_COUNTRY_INFO",
			},
			"signatures": []map[string]interface{}{
				map[string]interface{}{
					"signatureVersion": 1,
					"signature":        "eyJhbGdvcml0aG0iOiJSU0EtU0hBMjU2Iiwia2V5SWQiOiJ0ZXN0LXNpZ25pbmctY2VydCIsInNpZ25hdHVyZVZlcnNpb24iOjF9.hrVRgjCwXVvE2OOSpDZ58hR+59aFNwYDyjQgKk3auukd7pcegmE2CzPCa0bJ0ZsRAcKkCTJrWo5iDzNhMBWRyaMOv5zWSrthlf7G128qvIlpMT0YNY+n/FaOHE73uLrS/g7swl3/qH/BGFG2Hu4RlL48eb3lLKqTt2xKHdCs6Cd4RMfJPYnzgvI4BNrFUKsjkcu+WD4OO2A27Pq1n50cMchmcaXadJhGrOqH5YmHdOCj5NSHzJYrsW0HPlpuAx/ECMeIZYDh6RMqaFM2DXzdKX9NmmyqzJ3o/0lkk/N97gfVRLW5hA29yeAwaCViZNCP8iC9aO0q9fQojoa7NQnAtw==",
				},
			},
		},
	}
}

func (tv *Tv) Register(key string) (newKey string, err error) {
	helloMsg := Msg{
		Type:    "register",
		Id:      makeId(),
		Payload: helloPayload(),
	}
	ch := make(chan Msg, 1)
	tv.registerRespCh(helloMsg.Id, ch)
	defer tv.unregisterRespCh(helloMsg.Id)

	if key != "" {
		helloMsg.Payload["client-key"] = key
	}

	err = tv.writeJSON(&helloMsg)
	if err != nil {
		return "", err
	}

	var respMsg Msg
	var ok bool

	for {
		select {
		case respMsg, ok = <-ch:
			if !ok {
				return "", ErrNoResponse
			}
			err = checkResponse(respMsg)
			if err != nil {
				return "", err
			}

		case <-time.After(RegisterTimeout):
			return "", ErrTimeout
		}
		if respMsg.Type != "response" {
			break
		}
	}
	if respMsg.Type != "registered" {
		return "", ErrRegistrationFailed
	}
	if tmp, ok := respMsg.Payload["client-key"]; ok {
		tmp, ok := tmp.(string)
		if !ok {
			return "", errors.New("client-key from TV is not a string")
		}
		newKey = tmp
	}
	return newKey, nil
}

func (tv *Tv) writeJSON(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "JSON marshal error")
	}
	tv.debug("write: ", buf)
	tv.wsWriteMutex.Lock()
	err = tv.ws.WriteMessage(websocket.TextMessage, buf)
	tv.wsWriteMutex.Unlock()
	if err != nil {
		return errors.Wrap(err, "websocket write error")
	}
	return nil
}

func (tv *Tv) registerRespCh(id string, ch chan<- Msg) {
	tv.respChMutex.Lock()
	if tv.respCh == nil {
		tv.respCh = make(map[string]chan<- Msg)
	}
	tv.respCh[id] = ch
	tv.respChMutex.Unlock()
}

func (tv *Tv) unregisterRespCh(id string) {
	tv.respChMutex.Lock()
	if ch, ok := tv.respCh[id]; ok {
		close(ch)
		delete(tv.respCh, id)
	}
	tv.respChMutex.Unlock()
}

func (tv *Tv) Close() (err error) {
	err = tv.ws.Close()
	return err
}

type Payload map[string]interface{}

func (p *Payload) UnmarshalJSON(b []byte) (err error) {
	var tmp map[string]interface{}
	err = json.Unmarshal(b, &tmp)
	if err == nil {
		*p = tmp
	} else {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			err = nil
			*p = nil
		}
	}
	return err
}

type Msg struct {
	Type    string  `json:"type,omitempty"`
	Id      string  `json:"id,omitempty"`
	Uri     string  `json:"uri,omitempty"`
	Payload Payload `json:"payload,omitempty"`
	Error   string  `json:"error,omitempty"`
}

func (tv *Tv) MessageHandler() (err error) {
	defer func() {
		// close the channels to indicate that the reader is exiting
		tv.respChMutex.Lock()
		for _, ch := range tv.respCh {
			close(ch)
		}
		tv.respCh = nil
		tv.respChMutex.Unlock()
	}()

	for {
		messageType, p, err := tv.ws.ReadMessage()
		if err != nil {
			return err
		}
		tv.debug("read: ", p)
		if messageType != websocket.TextMessage {
			tv.debug("non-text message type, ignored", nil)
			continue
		}
		var msg Msg
		err = json.Unmarshal(p, &msg)
		if err != nil {
			tv.debug("invalid json in message, ignored", nil)
			continue
		}
		tv.respChMutex.Lock()
		ch := tv.respCh[msg.Id]
		tv.respChMutex.Unlock()
		ch <- msg
	}
	// not reached
}

func (tv *Tv) RequestResponseParam(uri string, req Payload, resp interface{}) (err error) {
	r, err := tv.Request(uri, req)
	if err != nil {
		return err
	}
	return mapstructure.Decode(r, resp)
}

func (tv *Tv) Request(uri string, req Payload) (resp Payload, err error) {
	var msg Msg
	msg.Type = "request"
	msg.Id = makeId()
	msg.Uri = uri
	msg.Payload = req

	ch := make(chan Msg, 1)
	tv.registerRespCh(msg.Id, ch)
	defer tv.unregisterRespCh(msg.Id)

	err = tv.writeJSON(&msg)
	if err != nil {
		return nil, err
	}
	select {
	case respMsg, ok := <-ch:
		if !ok {
			return nil, ErrNoResponse
		}
		err = checkResponse(respMsg)
		return respMsg.Payload, err
	case <-time.After(Timeout):
		return nil, ErrTimeout
	}
}

func checkResponse(r Msg) (err error) {
	switch r.Type {
	case "error":
		var err2 error
		if _, ok := r.Payload["returnValue"]; ok {
			err2 = checkPayloadReturnValue(r.Payload)
		}
		if err2 == nil {
			return errors.Errorf("API error: %s", r.Error)
		}
		return errors.Errorf("API error: %s - %s", r.Error, err2)
	case "response":
		return checkPayloadReturnValue(r.Payload)
	case "registered":
		if r.Payload == nil {
			return errors.New("nil payload")
		}
		return nil
	default:
		return errors.Errorf("unexpeced API response type: %s", r.Type)
	}
}

func checkPayloadReturnValue(p Payload) (err error) {
	if p == nil {
		return errors.New("nil payload")
	}
	returnValueI, ok := p["returnValue"]
	if !ok {
		return errors.New("returnValue missing")
	}

	returnValue, ok := returnValueI.(bool)
	if !ok {
		return errors.New("returnValue type is not bool")
	}
	if !returnValue {
		if p["errorCode"] != nil {
			return errors.Errorf("error %v: %v", p["errorCode"], p["errorText"])
		} else {
			return errors.New("returnValue: false, errorCode: nil")
		}
	}
	return nil
}

func (tv *Tv) Subscribe(uri string, req Payload, msgCh chan<- Msg) (id string, err error) {
	var msg Msg
	msg.Type = "subscribe"
	msg.Id = makeId()
	msg.Uri = uri
	msg.Payload = req

	tv.registerRespCh(msg.Id, msgCh)

	err = tv.writeJSON(&msg)
	if err != nil {
		tv.unregisterRespCh(msg.Id)
		return "", err
	}
	return msg.Id, nil
}

func (tv *Tv) Unsubscribe(uri string, id string, req Payload) error {
	var msg Msg
	msg.Type = "unsubscribe"
	msg.Id = id
	msg.Uri = uri
	msg.Payload = req

	tv.unregisterRespCh(msg.Id)

	return tv.writeJSON(&msg)
}

func (tv *Tv) MonitorStatus(uri string, req Payload, processPayload func(Payload) error, quit <-chan struct{}) (err error) {
	msgCh := make(chan Msg, 1)

	id, err := tv.Subscribe(uri, req, msgCh)
	if err != nil {
		return err
	}
	defer tv.Unsubscribe(uri, id, nil)

	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				return nil
			}
			if msg.Payload == nil {
				continue
			}
			err = processPayload(msg.Payload)
			if err != nil {
				return err
			}
		case <-quit:
			return nil
		}
	}
	// not reached
}

func makeId() string {
	return randSeq(8)
}

var randSeqLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = randSeqLetters[rand.Intn(len(randSeqLetters))]
	}
	return string(b)
}
