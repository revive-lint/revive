FROM --platform=$BUILDPLATFORM golang:1.27.3@sha256:9048d5e7751706200eb1719c7dae0436950d362f29a1da91fa05e5c63867c60f AS build

ARG VERSION
ARG REVISION
ARG BUILDTIME
ARG BUILDER

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /src
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags "-X 'github.com/mgechev/revive/cli.version=${VERSION}' -X 'github.com/mgechev/revive/cli.commit=${REVISION}' -X 'github.com/mgechev/revive/cli.date=${BUILDTIME}' -X 'github.com/mgechev/revive/cli.builtBy=${BUILDER}'"

FROM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb

ENV GOCACHE=/tmp/.cache/go-build GOMODCACHE=/tmp/.cache/go-mod

COPY --from=build /src/revive /revive

ENTRYPOINT ["/revive"]
