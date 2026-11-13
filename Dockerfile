FROM --platform=$BUILDPLATFORM golang:1.27.4-alpine3.24@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /src
COPY go.mod ./
RUN go mod download && go mod verify
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
  GOARM=${TARGETVARIANT#v} \
  go build -trimpath \
  -ldflags="-s -w" \
  -o /out/udp_mirror .

RUN apk --no-cache add libcap && \
  setcap cap_net_raw+ep /out/udp_mirror


FROM gcr.io/distroless/static-debian13@sha256:1a64f46364620465ef9e9a8f5413757dfdf22695e050bbfd7303eeaf17a59383
ENV CONFIG_FILE="/config/config.yml"
ENV CMD_FLAGS=""

WORKDIR /app
COPY --from=builder /out/udp_mirror /app/udp_mirror

ENTRYPOINT [ "udp-mirror", "-receivers" ]
EXPOSE 9999
