webostv - Go package for controlling LG WebOS Smart TV
======================================================

[![GoDoc](https://godoc.org/github.com/snabb/webostv?status.svg)](https://godoc.org/github.com/snabb/webostv)

This is Go library and a terminal application for remote control of
LG WebOS televisions.

Installation
------------

Install Go compiler if you do not have it:
```
curl https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz | sudo tar xzC /usr/local
PATH=$PATH:/usr/local/go/bin
```

Compile and install:
```
go install github.com/snabb/webostv/cmd/webostvremote
```
The resulting binary is at: ~/go/bin/webostvremote


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
