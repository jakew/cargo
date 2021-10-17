package v1

// Cargo describes the set of container images to render.
type Cargo struct {

	// CargoVersion is the version of Cargo to use.
	CargoVersion string `yaml:"cargoVersion"`

	// Manifests describe each container image.
	Manifests []Manifest `yaml:"manifests"`

	// Config is the shared starting config for container image.
	Config map[string]interface{} `yaml:"config"`

	// ConfigFiles lists the config files to be loaded and appended with the
	// Config before the container image is rendered.
	ConfigFiles []string `yaml:"configFiles"`
}

// Manifest describes a single container image.
type Manifest struct {

	// Name to refer to this manifest as. If blank, you can't render it by name.
	Name string `yaml:"name"`

	// Dockerfile template render.
	Dockerfile string `yaml:"dockerfile"`

	// Config values to be used when rendering.
	Config Config `yaml:"config"`

	// ConfigFiles are loaded and appended to Config before being used when
	// rendering.
	ConfigFiles []string `yaml:"configFiles"`

	// OutputFile describes where the result is written.
	OutputFile string `yaml:"outputFile"`
}
