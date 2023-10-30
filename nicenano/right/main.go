package main

import (
	"context"
	_ "embed"
	"log"
	"machine"
	"machine/usb"

	"github.com/funkycode/tinygo-corne/nicenano/niceview"
	keyboard "github.com/sago35/tinygo-keyboard"
)

var display = niceview.New()

func main() {
	usb.Product = "wireless-corne-right"
	// niceview.ClearScreen()

	// nicenano.GetBatReadings()
	// var err error
	// tx, err = connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := run()
	if err != nil {
		log.Fatal(err)
	}
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

	d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		{
			0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007, 0x0008, 0x0009,
			0x000A, 0x000B, 0x000C, 0x000D, 0x000E, 0x000F, 0x0010, 0x0011, 0x0012, 0x0013,
			0x0014, 0x0015, 0x0016, 0x0017,
		},
	})
	bleKeyboard := keyboard.BleTxKeyboard{
		RxBleName: "Corne-Left",
	}
	d.Keyboard = &bleKeyboard
	err := bleKeyboard.Connect()
	if err != nil {
		return err
	}
	d.Loop(context.Background())
	return nil
}
