webostv - Go package for controlling LG WebOS TV
================================================

[![GoDoc](https://godoc.org/github.com/snabb/webostv?status.svg)](https://godoc.org/github.com/snabb/webostv)

This is Go library and a terminal application for remote control of
LG WebOS smart televisions. Works on Linux and Windows and probably
on OS X as well. It has been tested with LG 42LB650V-ZN television.


Installing and using the remote control application
---------------------------------------------------

Download pre-built executable binary for your OS/architecture from:
https://github.com/snabb/webostv/releases

Run the downloaded application in a terminal window. On Linux, add execute
permissions after downloading (`chmod 755 webostvremote-linux-amd64`). The
IP address or name of the TV can be given as a command line argument:
```
./webostvremote 192.0.2.123
```
If the address is not supplied, it will try to connect to the default
address `LGsmartTV.lan`.


Building the remote control application from source
---------------------------------------------------

Install Go compiler if you do not have it:
```
curl https://dl.google.com/go/go1.12.linux-amd64.tar.gz | sudo tar xzC /usr/local
PATH=$PATH:/usr/local/go/bin
```
(See https://golang.org/dl/ for newer version and more detailed
instructions.)

Compile:
```
git clone https://github.com/snabb/webostv.git
cd webostv
go build ./cmd/webostvremote
```
The compiled binary `webostvremote` is produced in the current working
directory. If there are errors, try again with up-to-date Go compiler
version.


Simple example of using the library to turn off the TV
------------------------------------------------------

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

Unimplemented / TODO
--------------------

### webostv library

- Documentation.
- Consider the method names, some could be shortened.
- PIN based pairing.
- UPnP discovery?
- Add missing subscriptions?
- Play media? 

### webostvremote application

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
