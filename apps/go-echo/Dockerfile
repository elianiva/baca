# syntax=docker/dockerfile:1

FROM golang:1.22 AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./
COPY views ./views
COPY assets ./assets

# Build
RUN GOOS=linux go build -o /app/baca

FROM gcr.io/distroless/base-debian12
COPY --from=build /app/baca /baca
COPY --from=build /app/views /views
COPY --from=build /app/assets /assets

EXPOSE 8080

# Run
ENTRYPOINT ["/baca", "server"]