//go:build tinygo

package main

import (
	_ "embed"
	"fmt"
	"log"
	"machine"
	kb "machine/usb/hid/keyboard"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	kc "github.com/sago35/tinygo-keyboard/keycodes"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/bluetooth"
)

const (
	KeyPersentage keyboard.Keycode = 0xF000 | 0xC4
)

type keyEvent struct {
	layer, indx int
	state       keyboard.State
}

var rxEvent = make([]byte, 0, 3)
var notified bool

func main() {
	machine.InitSerial()
	var err error
	println("enabling adapter")
	err = adapter.Enable()
	if err != nil {
		log.Fatal(err)
	}
	//niceview.ClearScreen()
	rx, err = connectToSplit()
	if err != nil {
		println("failed to connect to other half:", err.Error())
		log.Fatal(err)
	}

	rx.EnableNotifications(
		func(buf []byte) {
			println("recieved buf len:", len(buf))
			rxEvent = buf
			notified = true
		},
	)
	err = connect()
	if err != nil {
		println("failed to establish LESC connection:", err.Error())
		log.Fatal(err)
	}
	println("esteblished LESC connection")
	registerHID()
	println("registered HID")
	// time.Sleep(10 * time.Second)
	// err = adapter.DefaultAdvertisement().Restart()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //adv.Stop()
	// for {

	// 	time.Sleep(100 * time.Millisecond)
	// }
	err = run()
	if err != nil {
		log.Fatal(err)
	}
}

func connect() error {

	// peerPKey := make([]byte, 0, 64)
	// privLesc, err := ecdh.P256().GenerateKey(rand.Reader)
	// if err != nil {
	// 	return err
	// }
	// lescChan := make(chan struct{})
	bluetooth.SetSecParamsBonding()
	//bluetooth.SetSecParamsLesc()
	bluetooth.SetSecCapabilities(bluetooth.NoneGapIOCapability)
	// time.Sleep(4 * time.Second)
	// println("getting own pub key")
	// var key []byte

	// pk := privLesc.PublicKey().Bytes()
	// pubKey := swapEndinan(pk[1:])
	//bluetooth.SetLesPublicKey(swapEndinan(privLesc.PublicKey().Bytes()[1:]))
	// pubKey = nil
	//println(" key is set")

	// println("register lesc callback")
	// adapter.SetLescRequestHandler(
	// 	func(pubKey []byte) {
	// 		peerPKey = pubKey
	// 		close(lescChan)
	// 	},
	// )

	println("def adv")
	adv := adapter.DefaultAdvertisement()
	println("adv config")
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinygo-corne",
		ServiceUUIDs: []bluetooth.UUID{
			bluetooth.ServiceUUIDDeviceInformation,
			bluetooth.ServiceUUIDBattery,
			bluetooth.ServiceUUIDHumanInterfaceDevice,
		},
	})
	println("adv start")
	return adv.Start()

	// select {
	// case <-lescChan:
	// 	peerPKey = append([]byte{0x04}, swapEndinan(peerPKey)...)
	// 	p, err := ecdh.P256().NewPublicKey(peerPKey)
	// 	if err != nil {
	// 		println("failed on parsing pub:", err.Error())
	// 		return err
	// 	}
	// 	println("calculating ecdh")
	// 	key, err = privLesc.ECDH(p)
	// 	if err != nil {
	// 		println("failed on curving:", err.Error())
	// 		return errfffffff
	// 	}
	// 	println("key len:", len(key))
	// 	return bluetooth.ReplyLesc(swapEndinan(key))
	// }

}

func swapEndinan(in []byte) []byte {
	var reverse = make([]byte, len(in))
	for i, b := range in[:32] {

		reverse[31-i] = b
	}
	if len(in) > 32 {
		for i, b := range in[32:] {
			reverse[63-i] = b
		}
	}

	return reverse
}

func registerHID() {
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDDeviceInformation,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDManufacturerNameString,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("Nice Keyboards"),
			},
			{
				UUID:  bluetooth.CharacteristicUUIDModelNumberString,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("nice!nano"),
			},
			{
				UUID:  bluetooth.CharacteristicUUIDPnPID,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte{0x02, 0x8a, 0x24, 0x66, 0x82, 0x34, 0x36},
				//Value: []byte{0x02, uint8(0x10C4 >> 8), uint8(0x10C4 & 0xff), uint8(0x0001 >> 8), uint8(0x0001 & 0xff)},
			},
		},
	})
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDBattery,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDBatteryLevel,
				Value: []byte{80},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	// gacc
	/*
	   device name r
	   apperance r
	   peripheral prefreed connection

	*/

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDGenericAccess,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDDeviceName,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("tinygo-corne"),
			},
			{

				UUID:  bluetooth.New16BitUUID(0x2A01),
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte{uint8(0x03c4 >> 8), uint8(0x03c4 & 0xff)}, /// []byte(strconv.Itoa(961)),
			},
			// {
			// 	UUID:  bluetooth.CharacteristicUUIDPeripheralPreferredConnectionParameters,
			// 	Flags: bluetooth.CharacteristicReadPermission,
			// 	Value: []byte{0x02},
			// },

			// // 		//
		},
	})

	//v := []byte{0x85, 0x02} // 0x85, 0x02
	reportValue := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	//var reportmap bluetooth.Characteristic

	// hid
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHumanInterfaceDevice,
		/*
			 - hid information r
			 - report map r
			 - report nr
			   - client charecteristic configuration
			   - report reference
			- report nr
			   - client charecteristic configuration
			   - report reference
			- hid control point wnr
		*/
		Characteristics: []bluetooth.CharacteristicConfig{
			// {
			// 	UUID:  bluetooth.CharacteristicUUIDHIDInformation,
			// 	Flags: bluetooth.CharacteristicReadPermission,
			// 	Value: []byte{uint8(0x0111 >> 8), uint8(0x0111 & 0xff), uint8(0x0002 >> 8), uint8(0x0002 & 0xff)},
			// },
			{
				//Handle: &reportmap,
				UUID:  bluetooth.CharacteristicUUIDReportMap,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: reportMap,
			},
			{

				Handle: &reportIn,
				UUID:   bluetooth.CharacteristicUUIDReport,
				Value:  reportValue[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				// protocl mode
				UUID:  bluetooth.New16BitUUID(0x2A4E),
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicReadPermission,
				// Value: []byte{uint8(1)},
				// WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
				// 	print("protocol mode")
				// },
			},
			{
				UUID:  bluetooth.CharacteristicUUIDHIDControlPoint,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission,
				//	Value: []byte{0x02},
			},
		},
	})
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D031,
		machine.D029,
		machine.D002,
		machine.D115,
		machine.D113,
		machine.D111,
	}

	rowPins := []machine.Pin{
		machine.D022,
		machine.D024,
		machine.D100,
		machine.D011,
	}

	matrixKeyCodes := [][]keyboard.Keycode{
		{
			keyboard.Keycode(kb.KeyTab), keyboard.Keycode(kb.KeyQ), keyboard.Keycode(kb.KeyW), keyboard.Keycode(kb.KeyE), keyboard.Keycode(kb.KeyR), keyboard.Keycode(kb.KeyT),
			keyboard.Keycode(kb.KeyLeftCtrl), keyboard.Keycode(kb.KeyA), keyboard.Keycode(kb.KeyS), keyboard.Keycode(kb.KeyD), keyboard.Keycode(kb.KeyF), keyboard.Keycode(kb.KeyG),
			keyboard.Keycode(kb.KeyLeftShift), keyboard.Keycode(kb.KeyZ), keyboard.Keycode(kb.KeyX), keyboard.Keycode(kb.KeyC), keyboard.Keycode(kb.KeyV), keyboard.Keycode(kb.KeyB),
			0, 0, 0, keyboard.Keycode(kb.KeyModifierGUI), kc.KeyMod1, keyboard.Keycode(kb.KeySpace),
		},
		{
			keyboard.Keycode(kb.KeyTab), keyboard.Keycode(kb.KeyA), keyboard.Keycode(kb.KeyW), keyboard.Keycode(kb.KeyE), keyboard.Keycode(kb.KeyR), KeyPersentage,
			keyboard.Keycode(kb.KeyLeftCtrl), keyboard.Keycode(kb.KeyA), keyboard.Keycode(kb.KeyS), keyboard.Keycode(kb.KeyD), keyboard.Keycode(kb.KeyF), keyboard.Keycode(kb.KeyG),
			keyboard.Keycode(kb.KeyLeftShift), keyboard.Keycode(kb.KeyZ), keyboard.Keycode(kb.KeyX), keyboard.Keycode(kb.KeyC), keyboard.Keycode(kb.KeyV), keyboard.Keycode(kb.KeyB),
			0, 0, 0, keyboard.Keycode(kb.KeyModifierGUI), kc.KeyMod1, keyboard.Keycode(kb.KeySpace),
		},
		{
			keyboard.Keycode(kb.KeyTab), keyboard.Keycode(kb.Key1), keyboard.Keycode(kb.Key2), keyboard.Keycode(kb.Key3), keyboard.Keycode(kb.Key4), keyboard.Keycode(kb.Key5),
			keyboard.Keycode(kb.KeyLeftCtrl), keyboard.Keycode(kb.KeyA), keyboard.Keycode(kb.KeyS), keyboard.Keycode(kb.KeyD), keyboard.Keycode(kb.KeyF), keyboard.Keycode(kb.KeyG),
			keyboard.Keycode(kb.KeyLeftShift), keyboard.Keycode(kb.KeyZ), keyboard.Keycode(kb.KeyX), keyboard.Keycode(kb.KeyC), keyboard.Keycode(kb.KeyV), keyboard.Keycode(kb.KeyB),
			0, 0, 0, keyboard.Keycode(kb.KeyModifierGUI), kc.KeyMod1, keyboard.Keycode(kb.KeySpace),
		}}

	mk := d.AddMatrixKeyboard(colPins, rowPins, matrixKeyCodes)

	// device information
	// model number string r 0x2A24
	// manufacture name string r  2A29
	// pnp id r 2A50
	//

	// d.AddBleKeyboard(24, "Corne-Left", [][]keyboard.Keycode{
	splitKb := [][]keyboard.Keycode{
		{
			keyboard.Keycode(kb.KeyY), keyboard.Keycode(kb.KeyU), keyboard.Keycode(kb.KeyI), keyboard.Keycode(kb.KeyO), keyboard.Keycode(kb.KeyP), keyboard.Keycode(kb.KeyBackspace),
			keyboard.Keycode(kb.KeyH), keyboard.Keycode(kb.KeyJ), keyboard.Keycode(kb.KeyK), keyboard.Keycode(kb.KeyL), keyboard.Keycode(kb.KeySemicolon), jp.KeyColon,
			keyboard.Keycode(kb.KeyN), keyboard.Keycode(kb.KeyM), keyboard.Keycode(kb.KeyComma), keyboard.Keycode(kb.KeyPeriod), keyboard.Keycode(kb.KeyBackslash), keyboard.Keycode(kb.KeyEsc),
			keyboard.Keycode(kb.KeyEnter), kc.KeyMod2, keyboard.Keycode(kb.KeyModifierRightAlt), 0, 0, 0,
		},
		{
			keyboard.Keycode(kb.KeyY), keyboard.Keycode(kb.KeyU), keyboard.Keycode(kb.KeyI), keyboard.Keycode(kb.KeyO), keyboard.Keycode(kb.KeyP), keyboard.Keycode(kb.KeyBackspace),
			keyboard.Keycode(kb.KeyH), keyboard.Keycode(kb.KeyJ), keyboard.Keycode(kb.KeyK), keyboard.Keycode(kb.KeyL), keyboard.Keycode(kb.KeySemicolon), jp.KeyColon,
			keyboard.Keycode(kb.KeyN), keyboard.Keycode(kb.KeyM), keyboard.Keycode(kb.KeyComma), keyboard.Keycode(kb.KeyPeriod), keyboard.Keycode(kb.KeyBackslash), keyboard.Keycode(kb.KeyEsc),
			keyboard.Keycode(kb.KeyEnter), kc.KeyMod2, keyboard.Keycode(kb.KeyModifierRightAlt), 0, 0, 0,
		},
		{
			keyboard.Keycode(kb.Key6), keyboard.Keycode(kb.Key7), keyboard.Keycode(kb.Key8), keyboard.Keycode(kb.Key9), keyboard.Keycode(kb.Key0), keyboard.Keycode(kb.KeyBackspace),
			keyboard.Keycode(kb.KeyH), keyboard.Keycode(kb.KeyJ), keyboard.Keycode(kb.KeyK), keyboard.Keycode(kb.KeyL), keyboard.Keycode(kb.KeySemicolon), jp.KeyColon,
			keyboard.Keycode(kb.KeyN), keyboard.Keycode(kb.KeyM), keyboard.Keycode(kb.KeyComma), keyboard.Keycode(kb.KeyPeriod), keyboard.Keycode(kb.KeyBackslash), keyboard.Keycode(kb.KeyEsc),
			keyboard.Keycode(kb.KeyEnter), kc.KeyMod2, keyboard.Keycode(kb.KeyModifierRightAlt), 0, 0, 0,
		},
	}
	//)

	keyChan := make(chan keyEvent, 10)

	mk.SetCallback(func(layer, index int, state keyboard.State) {
		println("keyboard event")
		keyChan <- keyEvent{layer: layer, indx: index, state: state}
	})
	go func() {
		for {
			var report []byte
			select {
			case x := <-keyChan:
				println("sent key:", x.indx)
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mods bits in array
					key := matrixKeyCodes[x.layer][x.indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
				println("sent key")
			default:

				if !notified {
					time.Sleep(1 * time.Millisecond)
					continue
				}
				println("got key")
				layer := int(rxEvent[0])
				indx := int(rxEvent[1])
				state := keyboard.State(rxEvent[2])
				notified = false

				println("key recieved:", indx)
				if state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mods bits in array
					key := splitKb[layer][indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
				println("sent key")
			}
		}
		//time.Sleep(10 * time.Millisecond)
	}()
	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func connectToSplit() (bluetooth.DeviceCharacteristic, error) {
	var tx bluetooth.DeviceCharacteristic
	var foundDevice bluetooth.ScanResult
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		fmt.Printf("%q:%#v\n", result.LocalName(), result.Address.String())
		if result.LocalName() != "corne-right" {
			return
		}
		foundDevice = result

		// Stop the scan.
		err := adapter.StopScan()
		if err != nil {
			return
		}
	})
	if err != nil {
		return tx, err
	}
	device, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return tx, err
	}
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDNordicUART})
	if err != nil {
		return tx, err
	}
	service := services[0]
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDUARTTX})
	if err != nil {
		return tx, err
	}
	tx = chars[0]
	return tx, nil

}
