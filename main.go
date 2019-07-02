package main

import (
	"Bitcoin-Tracker/client"
	"errors"
	"fmt"
	"github.com/d2r2/go-hd44780"
	i2c2 "github.com/d2r2/go-i2c"
)

func main() {
	//create client and make request
	c := client.NewClient(30)
	go c.Get()

	lcd, err := getLCD()
	if err != nil {
		c.ErrChan <- err
	}
	for {
		select {
		case resp := <-c.RespChan: //Handle response here
			write(lcd, resp.Bpi.USD.Rate)
		case err := <-c.ErrChan: //Handle error here
			write(lcd, err.Error())
		}
	}
}

func write(lcd *hd44780.Lcd, message string) {
	lcd.Home()
	lcd.SetPosition(0, 0)
	fmt.Fprint(lcd, message)
}
func getLCD() (*hd44780.Lcd, error) {
	i2c, err := i2c2.NewI2C(0x27, 1)
	if err != nil {
		return nil, errors.New("Unable to locate i2c")
	}
	defer i2c.Close()

	lcd, err := hd44780.NewLcd(i2c, hd44780.LCD_16x2)
	if err != nil {
		return nil, errors.New("Error with LCD")
	}

	lcd.BacklightOn()
	lcd.Clear()

	return lcd, nil
}
