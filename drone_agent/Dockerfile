FROM golang:1.22-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/drone-brain ./cmd/brain

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/drone-brain /usr/local/bin/drone-brain
ENTRYPOINT ["/usr/local/bin/drone-brain"]
