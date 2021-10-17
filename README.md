# Cargo

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
cargo render template.Dockerfile -c cargoconfig.yaml > Dockerfile
```

To build it directly:
```shell
cargo render template.Dockerfile -c cargoconfig.yaml | docker build -
```

### Using a cargo.yaml

You can alternatively create a file called `cargo.yaml` like the following:

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

To render these, run:
```shell
cargo
```

If you want to render a specific one, provide the name using `--manifest`:
```shell
cargo --manifest pkg2
```

### Init
You can setup a basic scaffolding using:
```shell
cargo init
```
