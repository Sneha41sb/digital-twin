const charts = {};

/**
 * Logic to determine health based on the "Residual" 
 * (Difference between Actual Sensor and AI Prediction)
 */
function getHealthStatus(temp, expected) {
    const residual = Math.abs(temp - expected);
    if (residual < 10) return { text: "Optimal", class: "normal" };
    if (residual < 20) return { text: "Warning", class: "warning" };
    return { text: "Critical Anomaly", class: "critical" };
}

/**
 * Dynamically creates the HTML card and Chart.js instance
 */
function createMotorCard(motorId) {
    const dashboard = document.getElementById("dashboard");
    if (!dashboard) {
        console.error("❌Error: Element with id='dashboard' not found!");
        return;
    }

    const card = document.createElement("div");
    card.className = "motor-card";
    card.innerHTML = `
        <div class="motor-header">
            <h3>MOTOR ${motorId}</h3>
            <span id="status-${motorId}" class="status-badge">Init</span>
        </div>
        <p id="risk-${motorId}" style="font-size: 0.8rem; opacity: 0.8;">Calculating RUL...</p>
        <canvas id="chart-${motorId}"></canvas>
    `;
    dashboard.appendChild(card);

    const ctx = document.getElementById(`chart-${motorId}`).getContext('2d');
    
    charts[motorId] = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                { 
                    label: 'Actual', 
                    borderColor: '#ef4444', 
                    backgroundColor: 'rgba(239, 68, 68, 0.1)',
                    data: [], 
                    tension: 0.3, 
                    pointRadius: 0,
                    fill: true
                },
                { 
                    label: 'Predicted', 
                    borderColor: '#38bdf8', 
                    data: [], 
                    borderDash: [5, 5], 
                    tension: 0.3, 
                    pointRadius: 0 
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: false, 
            scales: { 
                x: { display: false }, 
                y: { grid: { color: 'rgba(255,255,255,0.05)' } } 
            },
            plugins: { legend: { display: false } }
        }
    });
}

// --- The WebSocket Bridge ---
const socket = new WebSocket('ws://localhost:8080/ws');

socket.onopen = () => {
    console.log("WebSocket Stream Active");
    if (!charts["01"]) {
        createMotorCard("01");
    }
};

socket.onmessage = function(event) {
    
    const data = JSON.parse(event.data);
    const motorId = "01"; 

    if (!charts[motorId]) createMotorCard(motorId);

    const currentChart = charts[motorId];
    const now = new Date().toLocaleTimeString();

    /**
     *  DATA MAPPING BRIDGE
     * This ensures the graph works even if Go sends Capitalized keys.
     */
    const actual = data.actual_temp || data.Actual_temp || data.temperature || 0;
    const expected = data.expected_temp || data.Expected_temp || 0;
    const rul = data.time_to_fail || data.Time_to_fail || 0;
    const speed = data.rpm || data.Rpm || 0;

    
    if (currentChart.data.labels.length > 30) {
        currentChart.data.labels.shift();
        currentChart.data.datasets[0].data.shift();
        currentChart.data.datasets[1].data.shift();
    }
    
    currentChart.data.labels.push(now);
    currentChart.data.datasets[0].data.push(actual);
    currentChart.data.datasets[1].data.push(expected);
    
    currentChart.update('none'); 

    
    const health = getHealthStatus(actual, expected);
    const badge = document.getElementById(`status-${motorId}`);
    if (badge) {
        badge.innerText = health.text;
        badge.className = "status-badge " + health.class;
    }

    const riskEl = document.getElementById(`risk-${motorId}`);
    if (riskEl) {
        riskEl.innerText = `RUL: ${Math.round(rul)} Hours | RPM: ${Math.round(speed)}`;
    }
};

socket.onerror = (err) => {
    console.error("WebSocket Error:", err);
};