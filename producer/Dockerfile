FROM golang:1.23.4-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD [ "go","run","main.go" ]


