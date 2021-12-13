FROM node:16 as assets

WORKDIR /assets

COPY assets/yarn.lock assets/package.json ./

RUN yarn

COPY assets ./

RUN yarn build

FROM golang:1.17-alpine as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY --from=assets /assets/dist /go/src/app/assets/dist

COPY pkg ./pkg

COPY main.go .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o build

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app/build .

CMD ["/app/build"]