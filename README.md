# Cargo

### TL;DR

TL;DR: Cargo builds Dockerfiles from a template and a config YAML file.
Optionally, you can write a manifest for multiple Dockerfiles. Templates can use
Sprig functions.

### Install

Only current install method is installing using Go.
`go get -u github.com/jakew/cargo`


### Usage

Include [Go template](https://pkg.go.dev/text/template#pkg-overview) components
in your Dockerfile template. Example:

```shell
cat <<EOF > template.Dockerfile
FROM {{ .baseImage }}
RUN echo "{{ .message }}"
EOF
```

If you want to specify template variables based off of a file, you can create a
YAML file with an object at the root. Example:
```shell
cat <<EOF > cargoconfig.yaml
baseImage: alpine:latest
message: "Hello, World!"
EOF
```

To see the rendered Dockerfile:
```shell
cargo template.Dockerfile -c cargoconfig.yaml
```

To save it as a file:
```shell
cargo build template.Dockerfile -c cargoconfig.yaml > Dockerfile
```

To build it directly:
```shell
cargo build template.Dockerfile -c cargoconfig.yaml | docker build -
```

### Using a cargo.yaml

You can alternatively create a file called `cargo.yaml` like the following:

```yaml
cargoVersion: v0
manifests:
- dockerfile: template.Dockerfile
  config:
    message: Package 1
  configFiles:
    - pkg1/config.yaml
  outputFile: pkg1/Dockerfile
- dockerfile: template.Dockerfile
  config:
    message: Package 2
  configFile:
    - pkg2/config.yaml
  outputFile: pkg2/Dockerfile
config:
  baseContainer: alpine
```

In this case, the same template is used, but different Dockerfiles are
generated, each getting their own unique configs. To build this, you just run:
```shell
cargo
```

The root `config` value is applied to all manifests. Each manifest also has its
own local `config` set, and `configFile` value indicating if a config YAML file
should be loaded in.