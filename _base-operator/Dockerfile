# Copyright (c) 2020, salesforce.com, inc.
# All rights reserved.
# SPDX-License-Identifier: BSD-3-Clause
# For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY controller.json controller.json
COPY main.go main.go
COPY reconciler/ reconciler/
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/controller.json .
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
