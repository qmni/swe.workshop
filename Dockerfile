FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /bin/swe-workshop ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /bin/swe-workshop /app/swe-workshop
EXPOSE 8080
ENTRYPOINT ["/app/swe-workshop"]
