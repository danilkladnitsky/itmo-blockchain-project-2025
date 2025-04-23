FROM golang:1.22.3-alpine AS build

WORKDIR /app

RUN apk update && apk add --no-cache curl git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CONFIG_PATH ./config/local.yaml
ENV GOBIN="/go/bin"
ENV PATH="${GOBIN}:${PATH}"

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN ls -la /go/bin && /go/bin/swag --version

RUN /go/bin/swag init -d cmd/memory-book,internal/handlers/check_alive,internal/handlers/register,internal/handlers/authorize,internal/handlers/ELK_auth,internal/lib/api/error,internal/lib/api/response

RUN go build -o main cmd/memory-book/main.go

EXPOSE 8083

CMD ["./main"]
