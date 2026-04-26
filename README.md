OmniTwin 
OmniTwin is a high-performance Predictive Maintenance dashboard. It leverages a physics-informed Go backend and a gRPC-powered AI engine to monitor, visualize, and predict the health of industrial motors in real-time.

 System Architecture
The project is built on an Event-Driven Architecture designed for low latency and high scalability:
Simulation Layer: A Go-based digital skeleton simulating motor telemetry (RPM, Temperature, Vibration).
Inference Layer: (Optional/Future) Python-based gRPC service providing physics-informed predictions.
Communication: High-speed WebSockets for real-time data streaming between the backend and HMI.
Visualization: A modern Glassmorphism HMI built with HTML5/CSS3 and Chart.js for sub-100ms telemetry rendering.

 Key Features
Physics-Informed Monitoring: Tracks "Actual vs. Predicted" temperatures to identify anomalies using residual analysis.
Dynamic Multi-Motor Support: The HMI automatically instantiates new telemetry cards as motors connect to the stream.
Predictive Maintenance: Calculates Remaining Useful Life (RUL) in real-time based on simulated wear-and-tear patterns.
High-Frequency Visualization: Optimized Chart.js implementation handling 5Hz+ data updates without UI lag.

 Technical Deep Dive (Why this matters)
Polyglot Design: Demonstrates seamless integration between Go (for concurrency/speed) and the potential for Python (for ML).
Resilient UI: Implemented local asset hosting and "Safety Mappings" to ensure HMI availability in air-gapped industrial environments.
Cross-Platform Mastery: Developed and debugged across WSL (Linux) and Windows environments.

 Project Structure
Bash
├── client.go       # Go Backend & WebSocket Server
├── app.js          # Telemetry Logic & Chart Management
├── index.html      # Glassmorphism HMI Layout
├── style.css       # Industrial Dashboard Styling
├── motor.proto     # gRPC definitions for sensor data
└── chart.min.js    # Localized Charting Library
