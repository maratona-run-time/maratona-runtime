FROM golang:alpine AS build
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY utils/ utils/
COPY queue/ queue/
COPY orchestrator/ orchestrator/
RUN go build -o main orchestrator/*.go

FROM alpine
WORKDIR /go/src/app
COPY orchestrator/pod.yml .
COPY --from=build /go/src/app/main .
CMD ["./main"]
