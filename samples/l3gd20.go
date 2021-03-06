// +build ignore

package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/l3gd20"
)

func main() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	gyro := l3gd20.New(bus, l3gd20.R250DPS)
	gyro.Debug = true
	defer gyro.Close()

	gyro.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	orientations, err := gyro.Orientations()
	if err != nil {
		log.Panic(err)
	}

	timer := time.Tick(250 * time.Millisecond)

	for {
		select {
		case <-timer:
			orientation := <-orientations
			log.Printf("x: %v, y: %v, z: %v", orientation.X, orientation.Y, orientation.Z)
		case <-quit:
			return
		}
	}
}
