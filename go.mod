module github.com/funkycode/tinygo-corne

go 1.19

require (
	github.com/sago35/tinygo-keyboard v0.0.0-00010101000000-000000000000
	tinygo.org/x/tinyfont v0.4.0
)

require (
	github.com/bgould/tinygo-rotary-encoder v0.0.0-20221224155058-c26fcc9a3d20 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/muka/go-bluetooth v0.0.0-20221213043340-85dc80edc4e1 // indirect
	github.com/saltosystems/winrt-go v0.0.0-20230921082907-2ab5b7d431e1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tinygo-org/cbgo v0.0.4 // indirect
	golang.org/x/sys v0.11.0 // indirect
	tinygo.org/x/bluetooth v0.8.0

	// github.com/funkycode/bluetooth v0.0.0
	tinygo.org/x/drivers v0.25.0 // indirect
)

replace github.com/bgould/tinygo-rotary-encoder => github.com/akif999/tinygo-rotary-encoder v0.0.0-20230411081648-5d87ee99295e

replace github.com/sago35/tinygo-keyboard => github.com/funkycode/tinygo-keyboard v0.0.0-20231030191234-2c7935fb51b9

replace tinygo.org/x/bluetooth => /home/zogg/dev/bluetooth
