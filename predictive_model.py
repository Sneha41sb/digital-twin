import numpy as np
from sklearn.linear_model import LogisticRegression

# Generate synthetic data
np.random.seed(42)

data_size = 1000

temperature = np.random.uniform(40, 90, data_size)
vibration = np.random.uniform(0, 3, data_size)
rpm = np.random.uniform(800, 2000, data_size)

# Label: 1 = failure risk
labels = []

for t, v, r in zip(temperature, vibration, rpm):
    if t > 75 or v > 2:
        labels.append(1)
    else:
        labels.append(0)

X = np.column_stack((temperature, vibration, rpm))
y = np.array(labels)

# Train model
model = LogisticRegression()
model.fit(X, y)

print("Model trained!")

# Test prediction
test = [[80, 2.5, 1500]]
prediction = model.predict(test)
prob = model.predict_proba(test)

print("Prediction:", prediction)
print("Probability:", prob)