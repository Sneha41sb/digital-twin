package main

import (
	"fmt"
	"math/rand"
	"time"
)
type Motor struct{
	ID int
	Temperature float64
	Vibration float64
	RPM int
}
func main(){
	rand.Seed(time.Now().UnixNano())
	motor := Motor{
		ID: 1,
		Temperature: 60.0,
		Vibration: 1.0,
		RPM: 1200,
	}
	for {
		tempNoise := rand.Float64()*2-1
		motor.Temperature+=tempNoise
		if(motor.Temperature<40){
			motor.Temperature=40
		}
		if(motor.Temperature>80){
			motor.Temperature=80
		}
		motor.Vibration+=rand.Float64()*0.2-0.1
		if(motor.Temperature>75){
			motor.Vibration +=0.5
		}
		if(motor.Vibration<0){
			motor.Vibration=0
		}
		motor.RPM=900+rand.Intn(600)
	
	fmt.Printf(
		"Motor ID: %d | Temp: %f°C | Vibration: %.2f mm/s | RPM: %d \n",
		motor.ID,
		motor.Temperature,
		motor.Vibration,
		motor.RPM,
	)
	time.Sleep(1* time.Second)
	}
}
