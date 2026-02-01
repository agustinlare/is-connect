FROM golang:1.24-alpine as base

WORKDIR $GOPATH/src/smallest-golang/app/

COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 go build -o /is-connect .

FROM gcr.io/distroless/static

COPY --from=base /is-connect .

CMD ["./is-connect"]
