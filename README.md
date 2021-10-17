# Cargo

[![CI](https://github.com/jakew/cargo/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/jakew/cargo/actions?query=branch%3Amain)
[![Go Reference](https://pkg.go.dev/badge/jakew/cargo.svg)](https://pkg.go.dev/jakew/cargo)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakew/cargo)](https://goreportcard.com/report/github.com/jakew/cargo)
[![License](https://img.shields.io/github/license/jakew/cargo)](https://github.com/jakew/cargo/blob/main/LICENSE)

### TL;DR

TL;DR: Cargo renders Dockerfiles from a template and a config YAML file.
Optionally, you can write a manifest for multiple Dockerfiles. Templates can use
Sprig functions.

Reason it exists: The Build args feature in Docker is great but it can't be used
for some specific pieces such as the FROM or ENTRYPOINT expressions. By using a
template, you can customize those. A lot of existing projects have setup bash
scripts using awk. This is an attempt to skip that setup.

### Install

Only current install method is installing using Go.
`go get -u github.com/jakew/cargo`


### Usage

Include [Go template](https://pkg.go.dev/text/template#pkg-overview) components
in your Dockerfile template. Example:

```shell
cat <<EOF > template.Dockerfile
FROM {{ env "BASE_IMAGE" }}
RUN echo "Hello, World!"
EOF

BASE_IMAGE="alpine:latest" cargo render template.Dockerfile
FROM alpine:latest
RUN echo "Hello, World!"
```

You can use a YAML File to specifiy values:
To specify variables using a file, you can create a YAML file with an object at
the root. Example:
```shell
cat <<EOF > template.Dockerfile
FROM {{ .baseImage }}
RUN echo {{ .message }}
EOF

cat <<EOF > cargoconfig.yaml
baseImage: alpine:latest
message: "Hello, World!"
EOF

cargo render template.Dockerfile -c cargoconfig.yaml
FROM alpine:latest
RUN echo Hello, World!
```

To save it as a file:
```shell
cargo render template.Dockerfile -c cargoconfig.yaml > Dockerfile
```

To build it directly:
```shell
cargo render template.Dockerfile -c cargoconfig.yaml | docker build -
```

### Using a Cargo file.

Cargo also allows you to declare your usage in a `cargo.yaml`. The Cargo file
has multiple "manifests", each one results in a seperate Dockerfile.

Config values and/or config files may be added to the manifest directly, or in
the Cargo root object. Values added at the root are shared with all manifests.

Here's an example of a two-package manifest:

```yaml
cargoVersion: v0
manifests:
- name: pkg1 
  dockerfile: template.Dockerfile
  config:
    message: Package 1
  configFiles:
    - pkg1/config.yaml
  outputFile: pkg1/Dockerfile
- name: pkg2
  dockerfile: template.Dockerfile
  config:
    message: Package 2
  configFile:
    - pkg2/config.yaml
  outputFile: pkg2/Dockerfile
config:
  baseContainer: alpine
```

In this case, the same template is used, but different Dockerfiles are
generated, each getting their own unique configs.

The root `config` value is applied to all manifests. Each manifest also has its
own local `config` set, and `configFile` value indicating if a config YAML file
should be loaded in.

Only one set of config values is provided to the rendering template. This is the
result of all of the YAML config values being merged together, with subsequent
values overwriting the prior values for any specific key. The order they are
merged in is:

- root config object.
- root configFiles, in order.
- the specific manifest's config object.
- the specific manifest's configFiles, in order.

Run `cargo` to render these. If you want to render a specific one, provide the
name using `--manifest`: `cargo --manifest pkg2`

### Init
You can setup a basic scaffolding using:
```shell
cargo init
ls
cargo.yaml          config.yaml         template.Dockerfile
```

### Example Repo
You can see an example [here](https://github.com/jakew/cargo/tree/main/example).