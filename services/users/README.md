# services/users

The core microservice responsible for user management, onboarding, and domain business logic. It utilizes high-performance internal communication via gRPC over HTTP/2 and exposes a public, developer-friendly RESTful interface using an embedded reverse API Gateway.

## Technologies Used

* **Go + Native gRPC:** High-throughput binary communication and strict API type safety.
* **gRPC-Gateway:** Automatic reverse-proxy for translating HTTP/JSON requests into gRPC payloads.
* **Buf (protovalidate):** Declarative, high-performance request validation powered by Common Expression Language (CEL) syntax embedded in the .proto definition files.
* **pkg/xerrors:** Shared Monorepo component used to maintain standardized JSON error schemas.

## Service Architecture

The service runs in a hybrid fashion inside a single process utilizing separate network ports:
* **Port :50051:** Internal gRPC server for secure, low-latency inter-service communication.
* **Port :8080:** Public-facing HTTP API Gateway exposed to Web, Mobile, or third-party client applications.

## Development Commands

Ensure you are located in the monorepo root or within this specific service directory to run the Buf toolchain:

### 1. Synchronize Buf Dependencies
Download official Google API descriptors and Buf validation signature schemas from the registry:
```bash
buf dep update
```

### 2. Compile Protocol Buffer Definitions
Generate Go structs (.pb.go), gRPC server boilerplate (_grpc.pb.go), and reverse proxy configurations (.pb.gw.go) inside the `pb/` directory:
```bash
buf generate
```

### 3. Start the Service Locally
Launch the application containing both the microservice runtime and the HTTP API Gateway:
```bash
go run main.go
```

## Integration Testing (HTTP/REST)

You can verify the user registration process by sending an HTTP POST request to the API Gateway using curl:

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Carlos", "username": "cortega", "password": "supersecretpassword"}'
```

### Graceful Shutdown
The application handles standard operating system signals (SIGINT, SIGTERM). Upon receiving a termination event, the Gateway rejects new incoming traffic, allows active RPC operations a 5-second window to complete, and tears down open socket listeners cleanly without aborting data transactions midway.
