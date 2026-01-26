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
		Temperature: 0.0,
		Vibration: 0.0,
		RPM: 0,
	}
	for {
		motor.Temperature=60+ rand.ExpFloat64()*40
		motor.Vibration=rand.Float64()*5
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