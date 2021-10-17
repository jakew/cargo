package renderer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	apiv1 "github.com/jakew/cargo/api/v1"
	"gopkg.in/yaml.v2"
)

func WriteScaffolding(rootPath, cargofile, tmplfile, configfile, outputfile string) error {
	if !isDir(rootPath) {
		if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
			return err
		}
	}

	cargopath := filepath.Join(rootPath, cargofile)
	tmplpath := filepath.Join(rootPath, tmplfile)
	configpath := filepath.Join(rootPath, configfile)

	for name, path := range map[string]string{
		"cargo file":    cargopath,
		"template file": tmplpath,
		"config file":   configpath,
	} {
		if hasFile(path) {
			return fmt.Errorf("%s %s already exists", name, cargopath)
		}
	}

	crg := apiv1.Cargo{
		CargoVersion: "v0",
		Config: apiv1.Config{
			"baseContainer": "alpine",
		},
		Manifests: []apiv1.Manifest{
			{
				Dockerfile: tmplfile,
				OutputFile: outputfile,
				ConfigFiles: []string{
					configfile,
				},
				Config: apiv1.Config{
					"message": "Hello World!",
				},
			},
		},
	}

	cfg := map[string]interface{}{
		"alpineTag": "latest",
	}

	if err := writeFile(cargopath, func(w io.Writer) error {
		return yaml.NewEncoder(w).Encode(crg)
	}); err != nil {
		return err
	}

	if err := writeFile(configpath, func(w io.Writer) error {
		return yaml.NewEncoder(w).Encode(cfg)
	}); err != nil {
		return err
	}

	return writeFile(tmplpath, func(w io.Writer) error {
		_, err := w.Write([]byte("FROM alpine:{{ .alpineTag }}\nRUN echo \"{{ .message }}\""))
		return err
	})
}

func writeFile(filename string, fn func(io.Writer) error) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(f)
}

func isDir(dir string) bool {
	fileInfo, err := os.Stat(dir)
	if !os.IsNotExist(err) && fileInfo.IsDir() {
		return true
	}
	return false
}
