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
    -ldflags "-X github.com/mgechev/revive/cli.version=${VERSION} -X github.com/mgechev/revive/cli.commit=${REVISION} -X github.com/mgechev/revive/cli.date=${BUILDTIME} -X github.com/mgechev/revive/cli.builtBy=${BUILDER}"

FROM scratch

COPY --from=build /src/revive /revive

ENTRYPOINT ["/revive"]
