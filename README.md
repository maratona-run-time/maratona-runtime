# Maratona-Runtime

Maratona Runtime implements the core service of judging solutions for ICPC and CTF-related coding competitions.

## Architecture

We are using a microservice architecture.
The system, for now, is separated in the following parts:

- **Verdict**: Receives information regarding problem and submission details.
- **Compiler**: Tasked with compiling a received source code and responding with the generated binary code (except for interpreted languages, such as Python).
- **Executor**: Executes a binary code or a interpreted language program against a set of inputs and responds with generated outputs.

Follows a visual representation of the communication between these services:

![Representation of the flow of information in the system communication.](assets/architecture.png)

For now, the system can be started using the [docker-compose](docker-compose.yml).
All communication between the services is done via HTTP, using a single docker network "maratona-net".

## Testing

All tests were developed using Go's "testing" package.

To run all tests on the project:

```go
go test ./...
```

For service-specific tests, search for "*_test.go" files, such as: "verdict_test.go".
One can execute those by running `go test` on the folder where they are located.

## Travis

All configuration is on the project's [travis yaml](.travis.yml) file.

Travis currently tests:

- If the project files are styled respecting official golang style guidelines (by using `gofmt`).
- If the project tests are running successfully
