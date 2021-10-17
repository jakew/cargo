package renderer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"

	apiv1 "github.com/jakew/cargo/api/v1"
)

// PrintDockerfile will write tmplfile to w with the cfgfile values provided.
func PrintDockerfile(w io.Writer, tmplfile string, cfgfile []string) error {
	config := apiv1.Config{}

	for _, f := range cfgfile {
		cfg, err := parseConfig(f)
		if err != nil {
			return fmt.Errorf("unable to load config file %s: %w", f, err)
		}
		config = mergeConfig(config, cfg)
	}

	tmpl, err := parseTemplate(tmplfile)
	if err != nil {
		return fmt.Errorf("unable to load template: %w", err)
	}

	if err := executeTemplate(w, tmpl, config); err != nil {
		return fmt.Errorf("unable to execute template: %w", err)
	}

	return nil
}

func contains(strs []string, str string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

func toLower(strs []string) []string {
	lc := make([]string, len(strs))
	copy(lc, strs)
	for k, v := range lc {
		lc[k] = strings.ToLower(v)
	}
	return lc
}

// RenderCargo reads and parses path and then builds and writes each manifest.
func RenderCargo(path string, manifestNames []string) error {
	dir := filepath.Dir(path)
	manifestNames = toLower(manifestNames)

	cargo := &apiv1.Cargo{}
	if err := loadYamlFromFile(path, cargo); err != nil {
		return fmt.Errorf("failed to open manifest: %w", err)
	}

	cfg := cargo.Config
	if cfg == nil {
		cfg = apiv1.Config{}
	}

	for _, f := range cargo.ConfigFiles {
		cfgFromFile, err := parseConfig(filepath.Join(dir, f))
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		cfg = mergeConfig(cfg, cfgFromFile)
	}

	for i, manifest := range cargo.Manifests {
		if len(manifestNames) != 0 {
			if !contains(manifestNames, manifest.Name) {
				continue
			}
		}
		if err := RenderManifest(manifest, cfg, dir); err != nil {
			return fmt.Errorf("unable to build manifest %d: %w", i, err)
		}
	}

	return nil
}

// RenderManifest appends cfg and config values and files provided and then
// writes the template.
func RenderManifest(man apiv1.Manifest, cfg apiv1.Config, dir string) error {
	if !validManifest(man) {
		return errors.New("must have Dockerfile template and output")
	}

	if man.Config != nil {
		cfg = mergeConfig(cfg, man.Config)
	}

	if man.ConfigFiles != nil {
		for _, f := range man.ConfigFiles {
			cfgFromFile, err := parseConfig(filepath.Join(dir, f))
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cfg = mergeConfig(cfg, cfgFromFile)
		}
	}

	tmpl, err := parseTemplate(filepath.Join(dir, man.Dockerfile))
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	outputpath := filepath.Join(dir, man.OutputFile)
	outdir := filepath.Dir(outputpath)
	if !isDir(outdir) {
		if err := os.MkdirAll(outdir, os.ModePerm); err != nil {
			return err
		}
	}
	w, err := os.Create(outputpath)
	if err != nil {
		return fmt.Errorf("failed to open Dockerfile for writing: %w", err)
	}
	defer w.Close()

	if err := executeTemplate(w, tmpl, cfg); err != nil {
		return fmt.Errorf("failed to print Dockerfile: %w", err)
	}

	return nil
}

// mergeConfig merges the provided Config values. If a key was already seen and
// the previous and current values are also Config, they are merged, otherwise
// the second value overwrites the first.
func mergeConfig(cfgs ...apiv1.Config) apiv1.Config {
	n := apiv1.Config{}
	for _, cfg := range cfgs {
		for k, v2 := range cfg {
			if v1, ok := n[k]; ok {
				if nv1, ok1 := v1.(apiv1.Config); ok1 {
					if nv2, ok2 := v2.(apiv1.Config); ok2 {
						n[k] = mergeConfig(nv1, nv2)
						continue
					}
				}
			}
			n[k] = v2
		}
	}

	return n
}

func hasFile(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func loadYamlFromFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	return yaml.NewDecoder(f).Decode(v)
}

func parseConfig(path string) (apiv1.Config, error) {
	cfg := apiv1.Config{}
	if err := loadYamlFromFile(path, cfg); err != nil {
		return nil, fmt.Errorf("unable to load config file %s: %w", path, err)
	}
	return cfg, nil
}

func validManifest(manifest apiv1.Manifest) bool {
	return manifest.OutputFile != "" && manifest.Dockerfile != ""
}

func parseTemplate(path string) (*template.Template, error) {
	tmpl, err := template.
		New(filepath.Base(path)).
		Funcs(sprig.TxtFuncMap()).
		ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load template: %w", err)
	}
	return tmpl, nil
}

func executeTemplate(w io.Writer, tmpl *template.Template, cfg apiv1.Config) error {
	if err := tmpl.Execute(w, cfg); err != nil {
		return fmt.Errorf("unable to execute template: %w", err)
	}

	return nil
}
