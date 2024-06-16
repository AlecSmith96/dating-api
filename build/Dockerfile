FROM golang:1.22-alpine as build

WORKDIR /app

COPY . .

COPY go.mod go.sum ./
RUN go mod download

RUN go mod vendor

# Build binary
RUN  GO111MODULE="on" CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=build /app/main .

# migration files are required to connect to db
COPY --from=build /app/db/ ./db/

# Expose port 8080 and run binary
EXPOSE 8080
CMD ["./main"]