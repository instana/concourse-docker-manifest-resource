FROM golang:1.13 AS build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY ./ /src
WORKDIR /src
RUN go build -o /tmp/check ./cmd/check
RUN go build -o /tmp/in ./cmd/in
RUN go build -o /tmp/out ./cmd/out

FROM docker:stable
COPY config.json /root/.docker/config.json
COPY --from=build /tmp/check /opt/resource/check
COPY --from=build /tmp/in /opt/resource/in
COPY --from=build /tmp/out /opt/resource/out
