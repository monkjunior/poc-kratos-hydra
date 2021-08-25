# Builder layer
FROM golang:1.16-buster as builder
LABEL org.opencontainers.image.authors="Monk Junior"

WORKDIR /builder

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Runtime layer
FROM centos:8
LABEL org.opencontainers.image.authors="Monk Junior"
WORKDIR /builder

COPY --from=builder /builder/bin/ory-poc /usr/bin/
COPY . .

EXPOSE 4435

CMD ["ory-poc"]
