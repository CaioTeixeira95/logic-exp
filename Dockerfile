# syntax=docker/dockerfile:1

FROM golang:1.19 AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /logic-exp ./cmd/

FROM build AS run-test
ENV GIN_MODE=release
RUN go test -v -race ./...

FROM build AS run-app

EXPOSE 8080

# Run
CMD [ "/logic-exp" ]