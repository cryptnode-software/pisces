# syntax=docker/dockerfile:1.2
# Line above allows us to use new buildkit features like cache mounts.

ARG ALPINE_VERSION=3.15
ARG GOLANG_VERSION=1.17.6
ARG DYN_SQL_MIGRATE_VERSION=0.0.16
ARG GOLANGCI_LINT_VERSION=v1.41.0-alpine

# Use a "container-info" container that holds some build-time information. Place it early in the
# file in order to help with docker caching because these commands typically don't change. Note that
# the last line of this file will update the build date which is a very quick process and will prevent
# caching if it is done early in the build.
FROM alpine:${ALPINE_VERSION} as container-info
RUN mkdir -p /container_info/
# Allow builds to pass in so that we can echo it on container startup.
ARG CONTAINER_BUILD_HOSTNAME=unspecified
ARG CONTAINER_BUILD_USER=unspecified
ARG CONTAINER_SRC_VERSION=unspecified
RUN echo "$CONTAINER_SRC_VERSION" > /container_info/src_version && \
    echo "$CONTAINER_BUILD_HOSTNAME" > /container_info/build_hostname && \
    echo "$CONTAINER_BUILD_USER" > /container_info/build_user


# build go code
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS build-base
RUN apk --no-cache add git
# RUN apk add --update-cache git ca-certificates curl openssl
# tell go that these modules are private
ARG GOPRIVATE="github.com/dynata/*,github.com/researchnow/*"
WORKDIR /src
ENV CGO_ENABLED=0
# Run go mod download as separate, early step so that we can cache downloaded packages and only invalidate
# if go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

FROM build-base as build
# We do not need to COPY code since we use a target mount
# COPY . /src/
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/pisces cmd/pisces/*.go


FROM build-base as unit-test
# If this parameter is changed at build time (to some non-cached value), it will force tests to run again.
ARG FORCE_UNIT_TEST_RUN=N
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && \
    go test -v "$([ "$FORCE_UNIT_TEST_RUN" != "N" ] && echo "-count=1")" ./... 2>&1 | tee /out/unit_test_out.txt | go tool test2json > /out/unit_test_out.json; \
    echo ${?} > /out/unit_test_rc.txt

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as mercury-convert
ARG GITHUB_TOKEN
RUN if [ -z "$GITHUB_TOKEN" ]; then echo "missing GITHUB_TOKEN env var which is needed for asset downloads"; false; fi
# update git config to use our token for pulling assets
RUN printf "[url \"https://%s@github.com/\"]\n\tinsteadOf = https://github.com/\n" ${GITHUB_TOKEN} >> /root/.gitconfig
RUN apk --no-cache add git

ARG GOPRIVATE="github.com/dynata/*,github.com/researchnow/*"
WORKDIR /src
ENV CGO_ENABLED=0

# This stage/target copies results from unit tests so that we can easily get them without volume mounts.
FROM scratch as unit-test-results
COPY --from=unit-test /out /

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} as lint-base
FROM build-base as lint
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    mkdir -p /out && \
    golangci-lint run --timeout 10m ./... | tee /out/lint_out.txt; \
    echo ${?} > /out/lint_rc.txt

# This stage/target copies results from lint so that we can easily get them without volume mounts.
FROM scratch as lint-results
COPY --from=lint /out /


# main container
FROM alpine:${ALPINE_VERSION}
LABEL maintainer="ars@dynata.com"

EXPOSE 4081


COPY --from=container-info /container_info /container_info
COPY scripts/. /
COPY --from=dyn-sql-migrate /bin/dyn-sql-migrate /usr/bin/dyn-sql-migrate
COPY db/ /db/
COPY --from=build /out/pisces /

ENTRYPOINT ["/startup.sh"]

# This command will run if anything above was not cached. Leave it as last line in file.
RUN date +"%Y-%m-%d %H:%M:%S UTC" > /container_info/build_date