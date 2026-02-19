const charts = {};

function createChart(motorId) {

    const container = document.getElementById("charts");

    const title = document.createElement("h3");
    title.innerText = "Motor " + motorId;

    const canvas = document.createElement("canvas");
    canvas.id = "chart-" + motorId;
    canvas.width = 800;
    canvas.height = 300;

    container.appendChild(title);
    container.appendChild(canvas);

    const ctx = canvas.getContext("2d");

    charts[motorId] = new Chart(ctx, {
        type: "line",
        data: {
            labels: [],
            datasets: [{
                label: "Temperature °C",
                data: [],
                borderColor: "red",
                fill: false
            }]
        },
        options: {
            responsive: false,
            animation: false
        }
    });
}

async function updateCharts() {

    const response = await fetch("http://localhost:8080/status");
    const motors = await response.json();

    for (const motorId in motors) {

        const motor = motors[motorId];

        if (!charts[motorId]) {
            createChart(motorId);
        }

        const labels = motor.tempHistory.map((_, i) => i);

        charts[motorId].data.labels = labels;
        charts[motorId].data.datasets[0].data = motor.tempHistory;

        charts[motorId].update();
    }
}

setInterval(updateCharts, 1000);
