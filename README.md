webostv - Go package for controlling LG WebOS TV
================================================

[![GoDoc](https://godoc.org/github.com/snabb/webostv?status.svg)](https://godoc.org/github.com/snabb/webostv)

This is Go library and a terminal application for remote control of
LG WebOS smart televisions. Works on Linux and Windows and probably
on OS X as well. It has been tested with LG 42LB650V-ZN television.

Simple example of using the library to turn off the TV:

```Go
package main

import "github.com/snabb/webostv"

func main() {
	tv, err := webostv.DefaultDialer.Dial("LGsmartTV.lan")
	if err != nil {
		panic(err)
	}
	defer tv.Close()
	go tv.MessageHandler()

	_, err = tv.Register("")
	if err != nil {
		panic(err)
	}

	err = tv.SystemTurnOff()
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
(Check https://golang.org/dl/ for newer version.)

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
- Add missing subscriptions?
- Play media? 

### webostvremote 

- Documentation.
- Volume mute.
- Make "store" a separate generic package (or find pre-existing one).
- Make it look better.
  * Colors.
  * Clarify channel + program info etc.
- Error popup?
- Program guide.
- TV mouse control.
- TV keyboard input.
- Play media?


License
-------

MIT
