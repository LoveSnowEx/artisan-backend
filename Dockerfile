FROM golang:1.22-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.45.0

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
# CMD ["sh", "-c", "while true; do sleep 30; done;"]
