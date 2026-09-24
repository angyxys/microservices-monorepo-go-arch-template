# pkg/xerrors

A centralized, reusable library for error management, translation, and formatting within the microservices architecture. It intercepts native gRPC errors and Buf `protovalidate` CEL rule violations, transforming them into clean, secure, and standardized JSON responses for front-end consumption through the API Gateway.

## Features

* **Infrastructure Masking:** Removes technical internal fields such as `@type` that expose the backend's internal architecture.
* **Standardized HTTP Status Codes:** Automatically translates native gRPC enums into real HTTP status codes (e.g., `InvalidArgument` maps to `400 Bad Request`).
* **Field Violation Extraction:** Dynamically processes structured `*validate.Violations` objects to return failed property paths as plain strings (e.g., `name`, `password`), supporting nested struct validation.

## Standard JSON Error Response Format

When a microservice returns a controlled error, the front-end will always receive the following structured payload:

```json
{
  "code": 400,
  "message": "field validations",
  "reason": "InvalidArgument",
  "localized_message": "validate, body request",
  "details": [
    {
      "field": "password",
      "description": "value length must be at least 6 runes"
    }
  ]
}
```

## Integration into an API Gateway

To use this error handler in any new microservice, import the package via Go Workspaces and inject it into the gRPC-Gateway multiplexor options during initialization:

```go
import (
	"pkg/xerrors" // Native local import via Go Workspaces
	"://github.com"
)

func main() {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(xerrors.CustomErrorHandler), // Single-line registration
	)
	// ... rest of server initialization
}
```
