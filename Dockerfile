FROM golang:1.24.3-alpine AS build

RUN apk add --no-cache git tzdata

ENV TZ=America/La_Paz

RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
COPY .env go.mod go.sum .
RUN go mod download

COPY . .

RUN go build -o turismo_backend .

FROM alpine:latest

ENV TZ=America/La_Paz

RUN apk add --no-cache ca-certificates tzdata

RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

COPY --from=build /app/turismo_backend .

EXPOSE 5750

CMD ["/app/turismo_backend"]