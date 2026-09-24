# Microservices Monorepo

A scalable, high-performance microservices architecture powered by Go, gRPC, and Buf. This repository uses Go Workspaces to manage local dependencies efficiently without remote versioning overhead, exposing both strict binary gRPC endpoints and public-facing RESTful JSON APIs via gRPC-Gateway.

## Repository Structure

The project is segregated into reusable shared packages and isolated business domain microservices:

* [pkg/xerrors](pkg/xerrors/README.md) - Centralized error management library that intercepts gRPC errors and converts Buf validation violations into unified, frontend-friendly JSON schemas.
* [services/users](services/users/README.md) - Core user onboarding and authentication microservice handling domain logic, request verification, and integrated HTTP reverse proxy gateway routing.

## Prerequisites

Ensure you have the following toolchains installed locally before working on this repository:

* Go 1.22 or higher (with Workspace support enabled)
* Buf CLI toolchain (for Protocol Buffer management and dependency resolution)

## Workspace Architecture

This monorepo utilizes Go Workspaces (`go.work`), allowing components inside `services/` to import internal shared libraries from `pkg/` natively using direct disk paths. Local changes in shared code register immediately across all utilizing services without requiring independent version tagging or manual code replacement.

The root setup includes an ignored local configurations file to isolate each developer's workspace path routing definitions.

## Global Development Commands

Run the following operations from the root directory to manage the global system stack:

### 1. Initialize and Sync Remote Schemas
Download the required external schemas (such as Google API annotations and Buf validation definitions) from the Buf Schema Registry (BSR):
```bash
buf dep update
```

### 2. Generate Global Multi-Service Code Boilerplate
Compile all active components, messages, gRPC client/servers, and HTTP API reverse-proxies simultaneously across the workspace layout definitions:
```bash
buf generate
```

### 3. Run a Specific Microservice Runtime
Navigate to the desired service directory to execute its compilation matrix and spin up its low-latency listeners:
```bash
cd services/users
go run main.go
```

## Production Security Notes

* All environment keys, configuration constants, database passwords, and asymmetric encryption tokens must reside inside isolated local `.env` files.
* Ensure you duplicate the provided `.env.example` templates to seed your local values. Do not commit actual `.env` configuration runtime payloads to the remote version control history.
