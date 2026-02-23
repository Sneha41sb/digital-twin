const charts = {};


// Health logic

function getHealthStatus(temp) {

    if (temp < 70)
        return { text: "Normal", class: "normal" };

    else if (temp < 75)
        return { text: "Warning", class: "warning" };

    else
        return { text: "Critical", class: "critical" };
}


// Risk logic

function getFailureRisk(temp, vibration) {

    if (temp > 75 || vibration > 2)
        return { text: "HIGH RISK", class: "risk-high" };

    else if (temp > 70 || vibration > 1.5)
        return { text: "MEDIUM RISK", class: "risk-medium" };

    else
        return { text: "LOW RISK", class: "risk-low" };
}


// Create card

function createMotorCard(motorId) {

    const dashboard =
        document.getElementById("dashboard");


    const card =
        document.createElement("div");

    card.className =
        "motor-card";


    const header =
        document.createElement("div");

    header.className =
        "motor-header";


    const title =
        document.createElement("h3");

    title.innerText =
        "Motor " + motorId;


    const statusBadge =
        document.createElement("span");

    statusBadge.id =
        "status-" + motorId;

    statusBadge.className =
        "status-badge";


    header.appendChild(title);
    header.appendChild(statusBadge);


    const riskText =
        document.createElement("p");

    riskText.id =
        "risk-" + motorId;


    const canvas =
        document.createElement("canvas");

    canvas.width = 240;
    canvas.height = 120;


    card.appendChild(header);
    card.appendChild(riskText);
    card.appendChild(canvas);

    dashboard.appendChild(card);


    const ctx =
        canvas.getContext("2d");


    charts[motorId] =
        new Chart(ctx, {

        type: "line",

        data: {

            labels: [],

            datasets: [

                {
                    label: "Temperature",
                    borderColor: "red",
                    data: [],
                    yAxisID: "yTemp"
                },

                {
                    label: "Vibration",
                    borderColor: "orange",
                    data: [],
                    yAxisID: "yVib"
                },

                {
                    label: "RPM",
                    borderColor: "blue",
                    data: [],
                    yAxisID: "yRPM"
                }

            ]
        },

        options: {

            responsive: false,

            animation: false,

            plugins: {

                legend: {

                    display: false
                }
            },

            scales: {

                yTemp: {

                    position: "left",

                    min: 40,

                    max: 90
                },

                yVib: {

                    position: "right",

                    min: 0,

                    max: 5
                },

                yRPM: {

                    position: "right",

                    min: 800,

                    max: 2000
                }
            }
        }
    });
}


// Update

async function updateDashboard() {

    const response =
        await fetch("http://localhost:8080/status");

    const motors =
        await response.json();


    for (const motorId in motors) {

        const motor =
            motors[motorId];


        if (!charts[motorId])
            createMotorCard(motorId);


        const labels =
            motor.tempHistory.map((_, i) => i);


        charts[motorId].data.labels =
            labels;


        charts[motorId].data.datasets[0].data =
            motor.tempHistory;


        charts[motorId].data.datasets[1].data =
            motor.tempHistory.map(() => motor.vibration);


        charts[motorId].data.datasets[2].data =
            motor.rpmHistory;


        charts[motorId].update();


        const health =
            getHealthStatus(motor.temperature);


        const badge =
            document.getElementById("status-" + motorId);

        badge.innerText =
            health.text;

        badge.className =
            "status-badge " + health.class;


        const risk =
            getFailureRisk(motor.temperature, motor.vibration);


        const riskText =
            document.getElementById("risk-" + motorId);

        riskText.innerText =
            "Failure Risk: " + risk.text;

        riskText.className =
            risk.class;
    }
}


updateDashboard();

setInterval(updateDashboard, 1000);