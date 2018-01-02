package input

import (
	"reflect"
)

type Input struct {
	DirectorUUID    string           `yaml:"director_uuid"`
	InstanceGroups  []InstanceGroup  `yaml:"jobs"`
	Update          UpdateConfig     `yaml:"update"`
	CloudConfig     CloudConfig      `yaml:"-"`
	Stemcells       []StemcellConfig `yaml:"stemcells,omitempty"`
	Variables       []Variable       `yaml:"variables,omitempty"`
	AvailableErrand []string         `yaml:"-"`
	IsDryRun        bool             `yaml:"-"`
}

type output struct {
	Name          string `yaml:"name"`
	Input         `yaml:"-,inline"`
	Releases      []ReleaseConfig      `yaml:"releases"`
	DiskPools     []DiskConfig         `yaml:"disk_pools,omitempty"`
	ResourcePools []ResourcePoolConfig `yaml:"resource_pools,omitempty"`
	Networks      []NetworkConfig      `yaml:"networks,omitempty"`
	Compilation   CompilationConfig    `yaml:"compilation,omitempty"`
}

func (i *Input) MarshalYAML() (interface{}, error) {
	output := output{
		Input: *i,
		Name:  "foo-deployment",
		Releases: []ReleaseConfig{
			{
				Name:    "foo-release",
				Version: "latest",
			},
		},
	}

	output.Input.Update.CanaryWatchTime = 4000
	output.Input.Update.UpdateWatchTime = 20

	if len(i.CloudConfig.PersistentDiskPools) > 0 {
		output.DiskPools = i.CloudConfig.PersistentDiskPools
	}

	if len(i.CloudConfig.ResourcePools) > 0 {
		output.ResourcePools = i.CloudConfig.ResourcePools
	}

	if !i.IsV2() && len(i.CloudConfig.Networks) > 0 {
		output.Networks = i.CloudConfig.Networks
		i.CloudConfig.Networks = []NetworkConfig{}
	}

	if !i.IsV2() && i.CloudConfig.Compilation.Network != "" {
		output.Compilation = i.CloudConfig.Compilation
		i.CloudConfig.Compilation = CompilationConfig{}
	}

	for index := range i.InstanceGroups {
		g := &i.InstanceGroups[index]

		for index := range g.Jobs {
			j := &g.Jobs[index]
			j.Release = "foo-release"
		}

		for index := range g.Networks {
			c := &g.Networks[index]
			if c.DefaultDNSnGW {
				c.Default = []string{"dns", "gateway"}
			}
		}
	}

	return output, nil
}

func (i *Input) IsV2() bool {
	return len(i.CloudConfig.AvailabilityZones) > 0
}

func (i Input) HasMigratedInstances() bool {
	for _, m := range i.InstanceGroups {
		if len(m.MigratedFrom) > 0 {
			return true
		}
	}
	return false
}

func (i Input) FindInstanceGroupByName(instanceGroupName string) (InstanceGroup, bool) {
	for _, instanceGroup := range i.InstanceGroups {
		if instanceGroup.Name == instanceGroupName {
			return instanceGroup, true
		}
	}
	return InstanceGroup{}, false
}

func (i Input) FindAzByName(azName string) (AvailabilityZone, bool) {
	for _, az := range i.CloudConfig.AvailabilityZones {
		if az.Name == azName {
			return az, true
		}
	}
	return AvailabilityZone{}, false
}

func (i Input) FindDiskPoolByName(diskName string) (DiskConfig, bool) {
	for _, disk := range i.CloudConfig.PersistentDiskPools {
		if disk.Name == diskName {
			return disk, true
		}
	}
	return DiskConfig{}, false
}

func (i Input) FindDiskTypeByName(diskName string) (DiskConfig, bool) {
	for _, disk := range i.CloudConfig.PersistentDiskTypes {
		if disk.Name == diskName {
			return disk, true
		}
	}
	return DiskConfig{}, false
}

func (i Input) FindNetworkByName(networkName string) (NetworkConfig, bool) {
	for _, network := range i.CloudConfig.Networks {
		if network.Name == networkName {
			return network, true
		}
	}
	return NetworkConfig{}, false
}

func (i Input) FindResourcePoolByName(resourcePoolName string) (ResourcePoolConfig, bool) {
	for _, resourcePool := range i.CloudConfig.ResourcePools {
		if resourcePool.Name == resourcePoolName {
			return resourcePool, true
		}
	}
	return ResourcePoolConfig{}, false
}

func (i Input) FindVmTypeByName(vmTypeName string) (VmTypeConfig, bool) {
	for _, vmType := range i.CloudConfig.VmTypes {
		if vmType.Name == vmTypeName {
			return vmType, true
		}
	}
	return VmTypeConfig{}, false
}

func (i Input) FindSubnetByIpRange(ipRange string) (SubnetConfig, bool) {
	for _, network := range i.CloudConfig.Networks {
		for _, subnet := range network.Subnets {
			if subnet.IpPool.IpRange == ipRange {
				return subnet, true
			}
		}
	}

	return SubnetConfig{}, false
}

func (i Input) FindStemcellByName(stemcellName string) (StemcellConfig, bool) {
	for _, stemcell := range i.Stemcells {
		if stemcell.Name == stemcellName {
			return stemcell, true
		}
	}
	return StemcellConfig{}, false
}

type ReleaseConfig struct {
	Name    string
	Version string
}

type CloudConfig struct {
	AvailabilityZones   []AvailabilityZone   `yaml:"azs,omitempty"`
	PersistentDiskPools []DiskConfig         `yaml:"-"`
	PersistentDiskTypes []DiskConfig         `yaml:"persistent_disk_types,omitempty"`
	Networks            []NetworkConfig      `yaml:"networks,omitempty"`
	Compilation         CompilationConfig    `yaml:"compilation,omitempty"`
	VmTypes             []VmTypeConfig       `yaml:"vm_types,omitempty"`
	ResourcePools       []ResourcePoolConfig `yaml:"-"`
}

type DiskConfig struct {
	Name            string            `yaml:"name"`
	Size            int               `yaml:"disk_size"`
	CloudProperties map[string]string `yaml:"cloud_properties"`
}

func (d DiskConfig) IsEqual(other DiskConfig) bool {
	return reflect.DeepEqual(d, other)
}

type CompilationConfig struct {
	Network          string            `yaml:"network,omitempty"`
	AvailabilityZone string            `yaml:"az,omitempty"`
	NumberOfWorkers  int               `yaml:"workers,omitempty"`
	CloudProperties  map[string]string `yaml:"cloud_properties"`
}

type AvailabilityZone struct {
	Name            string            `yaml:"name"`
	CloudProperties map[string]string `yaml:"cloud_properties"`
}

func (a AvailabilityZone) IsEqual(other AvailabilityZone) bool {
	return reflect.DeepEqual(a, other)
}

type VmTypeConfig struct {
	Name            string            `yaml:"name"`
	CloudProperties map[string]string `yaml:"cloud_properties,omitempty"`
}

func (v VmTypeConfig) IsEqual(other VmTypeConfig) bool {
	return reflect.DeepEqual(v, other)
}

type ResourcePoolConfig struct {
	Name            string            `yaml:"name"`
	Stemcell        StemcellConfig    `yaml:"stemcell"`
	CloudProperties map[string]string `yaml:"cloud_properties,omitempty"`
}

func (r ResourcePoolConfig) IsEqual(other ResourcePoolConfig) bool {
	return reflect.DeepEqual(r, other)
}

type UpdateConfig struct {
	Canaries        int   `yaml:"canaries"`
	CanaryWatchTime int   `yaml:"canary_watch_time"`
	MaxInFlight     int   `yaml:"max_in_flight"`
	UpdateWatchTime int   `yaml:"update_watch_time"`
	Serial          *bool `yaml:"serial,omitempty"`
}

type StemcellConfig struct {
	Name    string `yaml:"name,omitempty"`
	OS      string `yaml:"os"`
	Version string `yaml:"version"`
	Alias   string `yaml:"alias"`
}

func (s StemcellConfig) IsEqual(other StemcellConfig) bool {
	return s.Version == other.Version
}

type NetworkConfig struct {
	Name            string            `yaml:"name"`
	Type            string            `yaml:"type"`
	Subnets         []SubnetConfig    `yaml:"subnets"`
	CloudProperties map[string]string `yaml:"cloud_properties,omitempty"`
}

func (n NetworkConfig) IsEqual(other NetworkConfig) bool {
	return reflect.DeepEqual(n, other)
}

type SubnetConfig struct {
	AvailabilityZones []string          `yaml:"azs,omitempty"`
	IpPool            IpPool            `yaml:"-,inline"`
	DNS               []string          `yaml:"dns,omitempty"`
	CloudProperties   map[string]string `yaml:"cloud_properties"`
}
