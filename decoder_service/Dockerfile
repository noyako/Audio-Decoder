FROM docker:stable AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache go
RUN apk --no-cache add ca-certificates

WORKDIR /build

COPY . .

RUN cd decoder_service
WORKDIR /build/decoder_service
RUN go build -o main .

WORKDIR /dist

RUN cp /build/decoder_service/main .
RUN cp /build/config.yml .
RUN cp /build/decoder_service/start.sh .


FROM docker:stable

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist /

CMD sh start.sh