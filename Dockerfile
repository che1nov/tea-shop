FROM golang:1.25-alpine AS builder

ARG SERVICE
WORKDIR /src

COPY shared ./shared
COPY ${SERVICE} ./${SERVICE}

RUN cd ${SERVICE} && go mod download
RUN cd ${SERVICE} && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /out/app ./app

CMD ["./app"]
