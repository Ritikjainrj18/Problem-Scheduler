FROM golang:1.23.4-alpine

WORKDIR /app

COPY ./cmd ./cmd
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download
COPY . .

CMD ["/bin/sh", "-c", "go run ./cmd/migrate/main.go up && go run ./cmd/main.go"]

# CMD [ "go","run","./cmd/main.go" ] did shell as two commands were needed 



