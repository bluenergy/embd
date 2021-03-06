// +build ignore

package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	"github.com/kidoman/embd/motion/servo"
)

func main() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	pwm := pca9685.New(bus, 0x41)
	pwm.Freq = 50
	pwm.Debug = true
	defer pwm.Close()

	servo := servo.New(pwm, 0)
	servo.Debug = true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	turnTimer := time.Tick(500 * time.Millisecond)
	left := true

	servo.SetAngle(90)
	defer func() {
		servo.SetAngle(90)
	}()

	for {
		select {
		case <-turnTimer:
			left = !left
			switch left {
			case true:
				servo.SetAngle(70)
			case false:
				servo.SetAngle(110)
			}
		case <-c:
			return
		}
	}
}
