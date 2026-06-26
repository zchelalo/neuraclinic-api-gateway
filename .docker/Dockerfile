FROM golang:1.25.8-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/neuraclinic-api-gateway ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/neuraclinic-api-gateway /usr/local/bin/neuraclinic-api-gateway

ENV ENVIRONMENT=production

EXPOSE 8000

CMD ["neuraclinic-api-gateway"]
