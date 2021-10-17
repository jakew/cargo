# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
COPY main.go main.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o hello main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/hello .
USER 65532:65532

CMD ["You!"]

ENTRYPOINT ["/hello"]