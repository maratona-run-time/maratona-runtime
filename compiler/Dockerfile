FROM golang:alpine AS build
EXPOSE 8081
WORKDIR /go/src/app
COPY go.sum go.mod ./
RUN go mod download
COPY model/ model/
COPY utils/ utils/
COPY compiler/ compiler/
RUN go build compiler/main.go

FROM alpine
WORKDIR /go/src/app
RUN apk add g++ gcc go
COPY --from=build /go/src/app/main .
CMD ["./main"]
