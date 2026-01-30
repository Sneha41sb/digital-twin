package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)
type Motor struct{
	ID int
	Temperature float64
	Vibration float64
	RPM int
}
func simulateMotor(motor Motor, wg *sync.WaitGroup){
	defer wg.Done();
	for{
		tempNoise := rand.Float64()*2-1
		motor.Temperature += tempNoise
		if (motor.Temperature <40){
			motor.Temperature =40
		}
		if(motor.Temperature>80){
			motor.Temperature=80
		}
		motor.Vibration+=rand.Float64()*0.2-0.1
		if(motor.Temperature >75){
			motor.Vibration+=0.5
		}
		if(motor.Vibration<0){
			motor.Vibration=0;
		}
		motor.RPM=900+ rand.Intn(600)
		fmt.Printf(
			"[Motor %d] Temp:%.2f°C | Vib:%.2f |RPM:%d\n",
			motor.ID,
			motor.Temperature,
			motor.Vibration,
			motor.RPM,
		)
		time.Sleep(1*time.Second)
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	numberOfMotors := 10

	for i := 1; i <= numberOfMotors; i++ {
		motor := Motor{
			ID:          i,
			Temperature: 60.0,
			Vibration:   1.0,
			RPM:         1200,
		}

		wg.Add(1)
		go simulateMotor(motor, &wg)
	}

	wg.Wait()
}
