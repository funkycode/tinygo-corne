//go:build tinygo

package main

import (
	"context"
	_ "embed"
	"log"
	"machine"
	kb "machine/usb/hid/keyboard"

	"github.com/funkycode/tinygo-corne/nicenano/niceview"
	keyboard "github.com/sago35/tinygo-keyboard"
	kc "github.com/sago35/tinygo-keyboard/keycodes"
	"tinygo.org/x/bluetooth"
)

const (
	KeyPersentage keyboard.Keycode = 0xF000 | 0xC4
)

type KeyEvent struct {
	layer, indx int
	state       keyboard.State
}

func main() {
	machine.InitSerial()
	niceview.ClearScreen()

	err := run()

	if err != nil {
		log.Fatal(err)
	}
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

	bluetooth.SetSecParamsBonding()
	bluetooth.SetSecCapabilities(bluetooth.NoneGapIOCapability)

	err := adapter.Enable()
	if err != nil {
		log.Fatal(err)
	}
	adv := adapter.DefaultAdvertisement()

	niceview.ClearScreen()

	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinygo-corne",
		ServiceUUIDs: []bluetooth.UUID{

			bluetooth.ServiceUUIDDeviceInformation,
			bluetooth.ServiceUUIDBattery,
			bluetooth.ServiceUUIDHumanInterfaceDevice,
		},
		/*
		   gatt
		   gacc
		   parameters ?
		   battery service
		   device information
		   hid

		*/
	})
	adv.Start()
	// device information
	// model number string r 0x2A24
	// manufacture name string r  2A29
	// pnp id r 2A50
	//

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
				Value: reportMap, //hidReporMap, // make([]byte, 0, len(HID_REPORT_MAP)),
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					print("report map")
				},
			},
			{

				Handle: &reportIn,
				UUID:   bluetooth.CharacteristicUUIDReport,
				Value:  reportValue[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					print("report in")
				},
			},
			{
				// protocl mode
				UUID:  bluetooth.New16BitUUID(0x2A4E),
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicReadPermission,
				Value: []byte{uint8(1)},
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					print("protocol mode")
				},
			},
			{
				UUID:  bluetooth.CharacteristicUUIDHIDControlPoint,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission,
				//	Value: []byte{0x02},
			},
		},
	})
	/*
		d.AddBleKeyboard(24, "Corne-Left", [][]keyboard.Keycode{
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
		})
	*/
	keyChan := make(chan KeyEvent, 5)

	mk.SetCallback(func(layer, index int, state keyboard.State) {
		keyChan <- KeyEvent{layer: layer, indx: index, state: state}
	})
	go func() {
		for {
			select {
			case x := <-keyChan:
				var report []byte
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
			}
		}
	}()

	d.Loop(context.Background())

	return nil
}
