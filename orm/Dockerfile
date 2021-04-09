FROM golang:alpine AS build
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY utils/ utils/
COPY model/ model/
COPY queue/ queue/
COPY orm/ orm/
RUN go build orm/main.go

FROM alpine
WORKDIR /go/src/app
COPY --from=build /go/src/app/main .
CMD ["./main"]
