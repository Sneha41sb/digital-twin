import grpc
from concurrent import futures
import motor_pb2
import motor_pb2_grpc

class TwinService(motor_pb2_grpc.TwinServiceServicer):
    def StreamTelemetry(self, request_iterator, context):
        for data in request_iterator:
            print(f"Received: RPM={data.rpm}, Temp={data.temperature}")

            # Simple AI logic
            expected_temp = data.temperature * 0.95
            status = "Optimal" if data.temperature < 80 else "Critical"

            yield motor_pb2.Prediction(
                time_to_failure=100.5,
                health_status=status,
                expected_temp=expected_temp
            )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    motor_pb2_grpc.add_TwinServiceServicer_to_server(TwinService(), server)
    server.add_insecure_port('[::]:50051')

    print("AI Server running on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()