FROM golang:1.15.6-alpine

RUN apk add --no-cache git

WORKDIR /app/magic-cards

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/magic-cards .


# This container exposes port 8080 to the outside world
EXPOSE 8000

# Run the binary program produced by `go install`
CMD ["./out/magic-cards"]