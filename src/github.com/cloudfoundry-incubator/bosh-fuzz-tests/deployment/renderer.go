package deployment

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	yaml "gopkg.in/yaml.v2"
)

type Renderer interface {
	Render(input bftinput.Input, manifestPath string, cloudConfigPath string) error
}

type renderer struct {
	fs boshsys.FileSystem
}

func NewRenderer(fs boshsys.FileSystem) Renderer {
	return &renderer{
		fs: fs,
	}
}

func (g *renderer) Render(input bftinput.Input, manifestPath string, cloudConfigPath string) error {
	manifest, err := yaml.Marshal(&input)
	if err != nil {
		return bosherr.WrapErrorf(err, "Generating deployment manifest")
	}

	err = g.fs.WriteFile(manifestPath, manifest)
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving generated manifest")
	}

	cloudConfig, err := yaml.Marshal(input.CloudConfig)
	if err != nil {
		return bosherr.WrapErrorf(err, "Generating cloud config")
	}

	err = g.fs.WriteFile(cloudConfigPath, cloudConfig)
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving generated cloud config")
	}

	return nil
}
