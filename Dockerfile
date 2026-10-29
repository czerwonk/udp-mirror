FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4 AS builder

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


FROM gcr.io/distroless/static-debian13@sha256:034df5c1737e8241f6ed6423a12e720ee0f6d6f088609ee172f98f6f2166ef5d
ENV CONFIG_FILE="/config/config.yml"
ENV CMD_FLAGS=""

WORKDIR /app
COPY --from=builder /out/udp_mirror /app/udp_mirror

ENTRYPOINT [ "udp-mirror", "-receivers" ]
EXPOSE 9999
