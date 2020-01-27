FROM golang:1.13-alpine as build
COPY . /antenna
WORKDIR /antenna
RUN go build -o ./bin/antenna ./cmd/antenna
RUN go test ./...

FROM alpine
COPY --from=build /antenna/bin/antenna /usr/local/bin/antenna
COPY --from=build /antenna/configs/config.yml config.yml
CMD [ "antenna" ]
