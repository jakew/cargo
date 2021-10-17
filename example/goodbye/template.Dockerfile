# Build the manager binary
FROM golang:{{ .goVersion }} as builder

WORKDIR /workspace
COPY main.go main.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o goodbye main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/goodbye .
USER 65532:65532

CMD ["{{ .who }}"]

ENTRYPOINT ["/goodbye"]