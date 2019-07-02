package main

import (
	"fmt"
	"github.com/d2r2/go-hd44780"
	i2c2 "github.com/d2r2/go-i2c"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//create client and make request
	c := NewClient(30)
	go c.Get()

	i2c, err := i2c2.NewI2C(0x27, 1)
	if err != nil {
		panic("Unable to locate i2c")
	}
	defer i2c.Close()

	lcd, err := hd44780.NewLcd(i2c, hd44780.LCD_16x2)
	if err != nil {
		panic("Error with LCD")
	}

	lcd.BacklightOn()
	lcd.Clear()

	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case resp := <-c.RespChan: //Handle response here
			lcd.Home()
			lcd.SetPosition(0, 0)
			fmt.Fprint(lcd, resp.Bpi.USD.Rate)
		case err := <-c.ErrChan: //Handle error here
			lcd.Home()
			lcd.SetPosition(1, 0)
			fmt.Fprint(lcd, err)
		case <-sigs:
			fmt.Println("Got shutdown, exiting")
			// Break out of the outer for statement and end the program
			break loop
		}
	}
}
