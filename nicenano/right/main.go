package main

import (
	_ "embed"
	"log"
	"machine"
	"machine/usb"
	"time"

	kb "machine/usb/hid/keyboard"

	keyboard "github.com/sago35/tinygo-keyboard"
	kc "github.com/sago35/tinygo-keyboard/keycodes"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/bluetooth"
)

var tx = &bluetooth.Characteristic{}
var adapter = bluetooth.DefaultAdapter

type keyEvent struct {
	layer, indx int
	state       keyboard.State
}

func main() {
	usb.Product = "wireless-corne-right"
	err := advertise()
	if err != nil {
		log.Fatal(err)
	}
	err = run()
	if err != nil {
		log.Fatal(err)
	}
}

func advertise() error {

	err := adapter.Enable()
	if err != nil {
		return err
	}
	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "corne-right",
	})
	err = adv.Start()
	if err != nil {
		return err
	}

	var buf = make([]byte, 0, 3)

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDNordicUART,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: tx,
				UUID:   bluetooth.CharacteristicUUIDUARTTX,
				Value:  buf[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	return nil
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D111,
		machine.D113,
		machine.D115,
		machine.D002,
		machine.D029,
		machine.D031,
	}

	rowPins := []machine.Pin{
		machine.D022,
		machine.D024,
		machine.D100,
		machine.D011,
	}

	mk := d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
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

	keyChan := make(chan keyEvent, 10)

	mk.SetCallback(func(layer, index int, state keyboard.State) {
		println("keyboard event")
		keyChan <- keyEvent{layer: layer, indx: index, state: state}
	})
	go func() {
		for x := range keyChan {
			println("sent key:", x.indx)

			_, err := tx.Write([]byte{
				uint8(x.layer), uint8(x.indx), uint8(x.state),
			})

			if err != nil {
				println("failed to send key:", err.Error())
			}
			println("sent key")
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
