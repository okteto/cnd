# syntax = docker/dockerfile:experimental

FROM golang:1.16-buster as build

ENV KUBECTL_VERSION 1.21.0
ENV HELM_VERSION 3.6.0

WORKDIR /okteto

# installing kubectl
RUN curl -sLf --retry 3 -o kubectl https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    cp kubectl /usr/local/bin/kubectl && \
    chmod +x /usr/local/bin/kubectl && \
    /usr/local/bin/kubectl version --client=true

# installing helm
RUN curl -sLf --retry 3 -o helm.tar.gz https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    mkdir -p helm && tar -C helm -xf helm.tar.gz && \
    cp helm/linux-amd64/helm /usr/local/bin/helm && \
    chmod +x /usr/local/bin/helm && \
    /usr/local/bin/helm version

ENV CGO_ENABLED=0
ARG VERSION_STRING=docker
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/root/.cache go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache make build
RUN chmod +x /okteto/bin/okteto

# Test
RUN /okteto/bin/okteto version

FROM alpine:3

RUN apk add --no-cache bash ca-certificates
COPY --from=build /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=build /usr/local/bin/helm /usr/local/bin/helm
COPY --from=build /okteto/bin/okteto /usr/local/bin/okteto

ENV PS1="\[\e[36m\]\${OKTETO_NAMESPACE:-okteto}:\e[32m\]\${OKTETO_NAME:-dev} \[\e[m\]\W> "
