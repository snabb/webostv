webostv - Go package for controlling LG WebOS TV
================================================

[![GoDoc](https://godoc.org/github.com/snabb/webostv?status.svg)](https://godoc.org/github.com/snabb/webostv)

This is Go library and a terminal application for remote control of
LG WebOS smart televisions. The development has been done on Linux and it
works also on Windows. It has been tested with LG 42LB650V-ZN television.

Simple example of using the library to turn off the TV:

```Go
package main

import "github.com/snabb/webostv"

func main() {
	tv, err := webostv.DefaultDialer.Dial("LGsmartTV.lan")
	if err != nil {
		panic(err)
	}
	go tv.MessageHandler()

	_, err = tv.Register("")
	if err != nil {
		panic(err)
	}

	err = tv.SystemTurnOff()
	if err != nil {
		panic(err)
	}

	err = tv.Close()
	if err != nil {
		panic(err)
	}
}
```

Installing the remote control application
-----------------------------------------

Install Go compiler if you do not have it:
```
curl https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz | sudo tar xzC /usr/local
PATH=$PATH:/usr/local/go/bin
```

Compile and install:
```
go get github.com/snabb/webostv/cmd/webostvremote
```
The resulting binary is at: `~/go/bin/webostvremote`


Unimplemented / TODO
--------------------

### webostv library

- Documentation.
- PIN based pairing.
- UPnP discovery?
- Some subscriptions missing?
- Play media? 

### webostvremote 

- Documentation.
- Make it look better.
  * Colors.
  * Clarify channel + program info etc.
- Program guide.
- TV mouse control.
- TV keyboard input.
- Play media?


License
-------

MIT
