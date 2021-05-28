FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk --no-cache add ca-certificates

WORKDIR /build

COPY . .

RUN cd user_service
WORKDIR /build/user_service
RUN go build -o main .

WORKDIR /dist

RUN cp /build/user_service/main .
RUN cp /build/config.yml .


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist /


ENTRYPOINT ["/main"]