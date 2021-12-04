# Maratona-Runtime

[![Go Reference](https://pkg.go.dev/badge/github.com/maratona-run-time/Maratona-Runtime?readme=expanded.svg)](https://pkg.go.dev/github.com/maratona-run-time/Maratona-Runtime?readme=expanded)
[![Go Report Card](https://goreportcard.com/badge/github.com/maratona-run-time/Maratona-Runtime)](https://goreportcard.com/report/github.com/maratona-run-time/Maratona-Runtime)
[![codecov](https://codecov.io/gh/maratona-run-time/Maratona-Runtime/branch/master/graph/badge.svg?token=G1GDE2TBXU)](https://codecov.io/gh/maratona-run-time/Maratona-Runtime)
[![Maintainability](https://api.codeclimate.com/v1/badges/b4d5cb940795135cca45/maintainability)](https://codeclimate.com/github/maratona-run-time/Maratona-Runtime/maintainability)

Maratona Runtime implements the core service of judging solutions for ICPC and CTF-related coding competitions.

## Architecture

We are using a microservice architecture.
The system, for now, is separated in the following parts:

- **Verdict**: Receives information regarding problem and submission details.
- **Compiler**: Tasked with compiling a received source code and responding with the generated binary code (except for interpreted languages, such as Python).
- **Executor**: Executes a binary code or a interpreted language program against a set of inputs and responds with generated outputs.
- **Orchestrator**: Receives challenge submissions and is responsible for triggering their evaluation.
- **ORM**: Responsible for administering a PostgreSQL server and saving submissions and their verdicts.

Follows a visual representation of the communication between these services:

![Representation of the flow of information in the system communication.](assets/architecture.png)

For now, the system can be started using the [docker-compose](docker-compose.yml) and Kubernetes.

### Using docker-compose

The communication between the services is done via HTTP, using two docker networks:

- "maratona-net" supports general purpose communication between the Verdict, Compiler, Executor, Orchestrator and ORM.
- "database-net", on the other hand, only supports database-related communication, between the ORM and the Postgres database.

### Using Kubernetes

If you want to use your local docker images, you might have to configure minikube to use the local Docker environment. To do that run:

```
eval $(minikube -p minikube docker-env)
```

After this change the image names used on k8s/ files, removing the `mruntime/` prefix.

You'll need to rebuild any pre-existing images to make them accessible from the cluster after this step.
To build the images, run `docker-compose build`.

Another possibility is pulling the docker images from our [DockerHub registry](https://hub.docker.com/orgs/mruntime).

Then, to deploy the project run:

```bash
minikube start --cpus 4 --memory=8192 --vm=true
minikube addons enable ingress
kubectl apply -f k8s/
kubectl get ingress 
```

Add on your `/etc/hosts` the ip provided on the last command above for the route mart-route.

### Troubleshooting

```bash
Error from server (InternalError): error when creating "k8s/ingress.yml": Internal error occurred: failed calling webhook "validate.nginx.ingress.kubernetes.io": an error on the server ("") has prevented the request from succeeding
```

Run:

```bash
kubectl delete validatingwebhookconfigurations ingress-nginx-admission
```

## Submissions Database

We are using a PostgreSQL database to store challenges and submissions details.

![Relational scheme of our database.](assets/db.png)

## Testing

All tests were developed using Go's "testing" package.

To run all tests on the project:

```go
go test ./...
```

For service-specific tests, search for "*_test.go" files, such as: "verdict_test.go".
One can execute those by running `go test` on the folder where they are located.

### Coverage

To keep track of coverage information this project uses [codecov.io](https://codecov.io/github/maratona-run-time/Maratona-Runtime).

To check the test coverage locally one can use:

```bash
export CGO_ENABLED=0
go test ./... -coverprofile=coverage.txt -covermode=atomic
bash <(curl -s https://codecov.io/bash) -t $CODECOV_TOKEN
```

These commands should output a codecov URL with the coverage report.

## Travis

All configuration is on the project's [travis yaml](.travis.yml) file.

Travis currently tests:

- If the project files are styled respecting official golang style guidelines (by using `gofmt`).
- If the project tests are running successfully

Travis also uploads the updated code coverage report to [codecov.io](codecov.io).

## Tasks and Organization

The developers organized the development tasks on a [Trello Kanban](https://trello.com/b/tZnrTevw/kanban) (in portuguese).
