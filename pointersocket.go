package webostv

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type PointerSocket struct {
	Address string
	ws      *websocket.Conn
	sync.Mutex
}

func (dialer *Dialer) DialPointerSocket(address string) (ps *PointerSocket, err error) {
	wsDialer := dialer.WebsocketDialer
	if wsDialer == nil {
		wsDialer = websocket.DefaultDialer
	}
	ws, resp, err := wsDialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &PointerSocket{
		Address: address,
		ws:      ws,
	}, nil
}

func (tv *Tv) NewPointerSocket() (ps *PointerSocket, err error) {
	socketPath, err := tv.GetPointerInputSocket()
	if err != nil {
		return nil, err
	}
	return DefaultDialer.DialPointerSocket(socketPath)
}

func (ps *PointerSocket) MessageHandler() (err error) {
	for {
		_, _, err = ps.ws.ReadMessage()
		if err != nil {
			return err
		}
	}
	// not reached
}

func (ps *PointerSocket) Close() (err error) {
	ps.Lock()
	defer ps.Unlock()
	if ps.ws != nil {
		err = ps.ws.Close()
		ps.ws = nil
	}
	return err
}

func (ps *PointerSocket) writeMessage(messageType int, data []byte) error {
	ps.Lock()
	defer ps.Unlock()
	return ps.ws.WriteMessage(messageType, data)
}

func (ps *PointerSocket) Input(btype, bname string) (err error) {
	msg := "type:" + btype + "\n" + "name:" + bname + "\n\n"
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

// UP DOWN LEFT RIGHT HOME BACK DASH  and numbers

func (ps *PointerSocket) Move(dx, dy int) (err error) {
	msg := fmt.Sprintf("type:move\ndx:%d\ndy:%d\ndown:0\n\n", dx, dy)
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

func (ps *PointerSocket) Scroll(dx, dy int) (err error) {
	msg := fmt.Sprintf("type:scroll\ndx:%d\ndy:%d\ndown:0\n\n", dx, dy)
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}

func (ps *PointerSocket) Click() (err error) {
	msg := "type: click\n\n"
	return ps.writeMessage(websocket.TextMessage, []byte(msg))
}
