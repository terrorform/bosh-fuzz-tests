package dummy

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"os"
)

type DirectorOptions struct {
	Port             	  int
	DatabaseName     	  string
	DatabaseServer   	  string
	DatabaseUser     	  string
	DatabasePassword 	  string
	DatabasePort     	  int
	BaseDir          	  string
	DummyCPIPath     	  string
	VerifyMultidigestPath     string
}

type DirectorConfig struct {
	options        DirectorOptions
	numWorkers     int
	baseDir        string
	fs             boshsys.FileSystem
	assetsProvider bltassets.Provider
	dummyCPIPath   string
}

func NewDirectorConfig(
	options DirectorOptions,
	baseDir string,
	fs boshsys.FileSystem,
	assetsProvider bltassets.Provider,
	numWorkers int,
	dummyCPIPath string,
) *DirectorConfig {
	return &DirectorConfig{
		options:        options,
		numWorkers:     numWorkers,
		baseDir:        baseDir,
		fs:             fs,
		assetsProvider: assetsProvider,
		dummyCPIPath:   dummyCPIPath,
	}
}

func (c *DirectorConfig) DirectorConfigPath() string {
	return filepath.Join(c.baseDir, "director.yml")
}

func (c *DirectorConfig) CPIPath() string {
	return filepath.Join(c.baseDir, "cpi")
}

func (c *DirectorConfig) WorkerConfigPath(index int) string {
	return filepath.Join(c.baseDir, fmt.Sprintf("worker-%d.yml", index))
}

func (c *DirectorConfig) DirectorPort() int {
	return c.options.Port
}

func (c *DirectorConfig) Write() error {
	directorTemplatePath, err := c.assetsProvider.FullPath("director.yml")
	if err != nil {
		return err
	}

	t := template.Must(template.ParseFiles(directorTemplatePath))
	err = c.saveConfig(c.options.Port, c.DirectorConfigPath(), t)

	if err != nil {
		return err
	}

	cpiTemplatePath, err := c.assetsProvider.FullPath("cpi.sh")
	if err != nil {
		return err
	}

	cpiTemplate := template.Must(template.ParseFiles(cpiTemplatePath))

	err = c.saveCPIConfig(c.CPIPath(), cpiTemplate)

	if err != nil {
		return err
	}

	for i := 1; i <= c.numWorkers; i++ {
		port := c.options.Port + i
		err = c.saveConfig(port, c.WorkerConfigPath(i), t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DirectorConfig) saveConfig(port int, path string, t *template.Template) error {
	buffer := bytes.NewBuffer([]byte{})
	context := c.options
	context.Port = port
	context.BaseDir = c.baseDir
	err := t.Execute(buffer, context)
	if err != nil {
		return err
	}
	err = c.fs.WriteFile(path, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *DirectorConfig) saveCPIConfig(cpiPath string, t *template.Template) error {
	buffer := bytes.NewBuffer([]byte{})
	context := c.options
	context.DummyCPIPath = c.dummyCPIPath
	context.BaseDir = c.baseDir

	err := t.Execute(buffer, context)
	if err != nil {
		return err
	}
	err = c.fs.WriteFile(cpiPath, buffer.Bytes())
	if err != nil {
		return err
	}

	c.fs.Chmod(cpiPath, os.ModePerm)

	return nil
}
