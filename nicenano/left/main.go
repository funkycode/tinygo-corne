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
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

const (
	KeyPersentage keyboard.Keycode = 0xF000 | 0xC4
)

func main() {
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

	d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
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
		}})

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

	d.Loop(context.Background())
	return nil
}
