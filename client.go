package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Sneha41sb/digital-twin/proto"
)

// WebSocket setup
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var dataChan = make(chan []byte)

// Handle dashboard connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}
	defer conn.Close()

	log.Println("Dashboard connected")

	for data := range dataChan {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func main() {

	
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", handleConnections)

	go func() {
		log.Println("Dashboard running at http://localhost:8080")
		http.ListenAndServe(":8080", nil)
	}()

	
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	client := pb.NewTwinServiceClient(conn)

	stream, err := client.StreamTelemetry(context.Background())
	if err != nil {
		log.Fatal("Stream error:", err)
	}

	
	rpm := 1500.0

	for {
		// Simulate physics
		temp := 30.0 + (rpm * 0.015) + rand.Float64()*3
		vibration := (rpm * 0.001) + rand.Float64()

		data := &pb.SensorData{
			Rpm:         float32(rpm),
			Temperature: float32(temp),
			Vibration:   float32(vibration),
			Timestamp:   uint64(time.Now().Unix()),
		}

		
		if err := stream.Send(data); err != nil {
			log.Println("Send error:", err)
			continue
		}

		
		resp, err := stream.Recv()
		if err != nil {
			log.Println("Receive error:", err)
			continue
		}

		
		payload, _ := json.Marshal(map[string]interface{}{
			"rpm":           data.Rpm,
			"actual_temp":   data.Temperature,
			"expected_temp": resp.ExpectedTemp,
			"status":        resp.HealthStatus,
			"time_to_fail":  resp.TimeToFailure,
		})

		select {
		case dataChan <- payload:
		default:
		}

		log.Printf("RPM: %.1f | Temp: %.1f | Status: %s",
			rpm, data.Temperature, resp.HealthStatus)

		time.Sleep(200 * time.Millisecond)

		// Cycle RPM
		rpm += 20
		if rpm > 6000 {
			rpm = 1500
		}
	}
}