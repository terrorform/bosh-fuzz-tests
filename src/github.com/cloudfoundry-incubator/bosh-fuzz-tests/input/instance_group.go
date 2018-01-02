package input

import "reflect"

type InstanceGroup struct {
	Name               string                       `yaml:"name"`
	Instances          int                          `yaml:"instances"`
	AvailabilityZones  []string                     `yaml:"azs,omitempty"`
	PersistentDiskSize int                          `yaml:"persistent_disk,omitempty"`
	PersistentDiskPool string                       `yaml:"persistent_disk_pool,omitempty"`
	PersistentDiskType string                       `yaml:"persistent_disk_type,omitempty"`
	Networks           []InstanceGroupNetworkConfig `yaml:"networks"`
	MigratedFrom       []MigratedFromConfig         `yaml:"migrated_from,omitempty"`
	VmType             string                       `yaml:"vm_type,omitempty"`
	ResourcePool       string                       `yaml:"resource_pool,omitempty"`
	Stemcell           string                       `yaml:"stemcell,omitempty"`
	Jobs               []Job                        `yaml:"templates,omitempty"`
	Lifecycle          string                       `yaml:"lifecycle,omitempty"`
}

func (j InstanceGroup) IsEqual(other InstanceGroup) bool {
	return reflect.DeepEqual(j, other)
}

func (j InstanceGroup) HasPersistentDisk() bool {
	return j.PersistentDiskSize != 0 || j.PersistentDiskPool != "" || j.PersistentDiskType != ""
}

func (j InstanceGroup) FindNetworkByName(networkName string) (InstanceGroupNetworkConfig, bool) {
	for _, network := range j.Networks {
		if network.Name == networkName {
			return network, true
		}
	}
	return InstanceGroupNetworkConfig{}, false
}

type Job struct {
	Name    string
	Release string
}

type InstanceGroupNetworkConfig struct {
	Name          string   `yaml:"name"`
	DefaultDNSnGW bool     `yaml:"-"`
	Default       []string `yaml:"default,omitempty"`
	StaticIps     []string `yaml:"static_ips,omitempty"`
}

type MigratedFromConfig struct {
	Name             string `yaml:"name"`
	AvailabilityZone string `yaml:"az"`
}
