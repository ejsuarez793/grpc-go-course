FROM golang:1.17.8-alpine3.14
WORKDIR /build

# Fetch Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build -o ./docker-golang-grpc ./blog/server/server.go

EXPOSE 50051
CMD ["./docker-golang-grpc"]

# Create final image
# FROM alpine
# WORKDIR /
# COPY --from=builder ./build/docker-golang-grpc .
# EXPOSE 8080
#CMD ["./docker-golang-grpc"]