package main

import (
	"errors"
	"fmt"
	"github.com/d2r2/go-hd44780"
	i2c2 "github.com/d2r2/go-i2c"
)

func main() {
	//create client and make request
	c := NewClient(2)
	go c.get()

	lcd, err := getLCD()
	if err != nil {
		c.errChan <- err
	}
	for {
		select {
		case resp := <-c.respChan: //Handle response here
			write(lcd, resp.Bpi.USD.Rate)
		case err := <-c.errChan: //Handle error here
			write(lcd, err.Error())
		}
	}
}

func write(lcd *hd44780.Lcd, message string) {
	lcd.Home()
	lcd.SetPosition(0, 0)
	fmt.Println(message)
}
func getLCD() (*hd44780.Lcd, error) {
	i2c, err := i2c2.NewI2C(0x27, 2)
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
