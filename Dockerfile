FROM golang:1.13-alpine as build
COPY . /antenna
WORKDIR /antenna
RUN go build -o ./bin/antenna ./cmd/antenna
RUN CGO_ENABLED=0 go test ./...

FROM alpine
COPY --from=build /antenna/bin/antenna /usr/local/bin/antenna
RUN adduser -Dg '' antenna
WORKDIR /home/antenna
COPY --from=build /antenna/configs/config.yml config.yml
USER antenna
CMD antenna config.yml
