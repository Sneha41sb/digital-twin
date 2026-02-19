package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Motor defines the digital twin structure
type Motor struct {
	ID          int     `json:"id"`
	Temperature float64 `json:"temperature"`
	Vibration   float64 `json:"vibration"`
	RPM         int     `json:"rpm"`
	TempHistory []float64 `json:"tempHistory"`
	RPMHistory  []int     `json:"rpmHistory"`
}

// Global state
var (
	motors = make(map[int]*Motor)
	mutex  = sync.RWMutex{} // Using RWMutex for better performance on reads
)

func simulateMotor(motor *Motor) {
	for {

		mutex.Lock()

		// Simulate natural temperature drift
		tempNoise := rand.Float64()*2 - 1
		motor.Temperature += tempNoise

		if motor.Temperature < 40 {
			motor.Temperature = 40
		}
		if motor.Temperature > 80 {
			motor.Temperature = 80
		}

		// Vibration simulation
		motor.Vibration += rand.Float64()*0.2 - 0.1
		if motor.Temperature > 75 {
			motor.Vibration += 0.5
		}
		if motor.Vibration < 0 {
			motor.Vibration = 0
		}

		// RPM simulation
		motor.RPM = 1100 + rand.Intn(400)

		// Store history AFTER updating state
		motor.TempHistory = append(motor.TempHistory, motor.Temperature)
		motor.RPMHistory = append(motor.RPMHistory, motor.RPM)

		// Keep only last 20 points
		if len(motor.TempHistory) > 20 {
			motor.TempHistory = motor.TempHistory[1:]
		}

		if len(motor.RPMHistory) > 20 {
			motor.RPMHistory = motor.RPMHistory[1:]
		}

		mutex.Unlock()

		fmt.Printf("[Motor %d] Temp: %.2f°C | Vib: %.2f | RPM: %d\n",
			motor.ID, motor.Temperature, motor.Vibration, motor.RPM)

		time.Sleep(1 * time.Second)
	}
}


// --- API Handlers ---

// Root handler to confirm server is up
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Digital Twin API is running</h1><p>Check <a href='/status'>/status</a> for data.</p>")
}

// Status handler to return JSON of all motors
func statusHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	mutex.RLock()
	defer mutex.RUnlock()

	json.NewEncoder(w).Encode(motors)
}


func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Initialize and start 5 motor twins
	for i := 1; i <= 10; i++ {
		m := &Motor{
			ID:          i,
			Temperature: 60.0,
			Vibration:   1.0,
			RPM:         1200,
		}
		motors[i] = m
		go simulateMotor(m) // Run each motor in its own goroutine
	}

	// Route definitions
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/status", statusHandler)

	fmt.Println("----------------------------------------------")
	fmt.Println("Digital Twin Server starting on http://localhost:8080")
	fmt.Println("Access http://localhost:8080/status to see data")
	fmt.Println("----------------------------------------------")

	// Start the server (this is a blocking call)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Could not start server: %s\n", err)
	}
}