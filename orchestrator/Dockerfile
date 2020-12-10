FROM golang:alpine
EXPOSE 8080
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY errors/ errors/
COPY orchestrator/ orchestrator/
CMD ["go", "run", "orchestrator/main.go"]