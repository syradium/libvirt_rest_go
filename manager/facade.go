package manager

import (
	"fmt"
	"strconv"

	"github.com/rgbkrk/libvirt-go"
)

// VPS is a simplified version of actual VPS stored in libvirt
type VPS struct {
	Name     string
	RAM      uint64
	DiskSize uint64
}

// LibvirtFacade is a facade over libvirt for easier VPS management
type LibvirtFacade struct {
	conn            *libvirt.VirConnection
	defaultPoolName string
}

// NewLibvirtFacade is a simple constructor for LibvirtFacade
func NewLibvirtFacade(uri string, defaultPoolName string) (*LibvirtFacade, error) {
	conn, err := libvirt.NewVirConnection(uri)
	if err != nil {
		return nil, err
	}
	facade := LibvirtFacade{conn: &conn, defaultPoolName: defaultPoolName}
	return &facade, nil
}

// Close closes connection to libvirt daemon
func (l *LibvirtFacade) Close() (int, error) {
	return l.conn.CloseConnection()
}

// CreateVPSDisk givens size in megabytes creates a disk for VPS given it's name
func (l *LibvirtFacade) CreateVPSDisk(vpsName string, sizeM uint64) error {
	_, err := l.createVolume(l.defaultPoolName, l.generateVPSDiskName(vpsName), sizeM)
	return err
}

// CreateVPS defines VPS given a unique name, RAM size and creates a disk of given size
func (l *LibvirtFacade) CreateVPS(name string, ramSize uint64, diskSize uint64) (VPS, error) {
	// TODO: Create VPS and Disk for real
	vps := VPS{Name: name, RAM: ramSize, DiskSize: diskSize}
	return vps, nil
}

func (l *LibvirtFacade) generateVPSDiskName(vpsName string) string {
	// TODO: Some intellectual VPS disk name generation
	return fmt.Sprintf("%s-disk-0", vpsName)
}

func (l *LibvirtFacade) createVolume(poolName string, volumeName string, sizeM uint64) (*libvirt.VirStorageVol, error) {
	xmlConfig := `<volume> <name>` + volumeName + `</name>
    <allocation>0</allocation>
    <capacity unit="M">` + strconv.FormatUint(sizeM, 10) + `</capacity>
    <target>
      <path>` + volumeName + `</path>
      <permissions>
        <owner>107</owner>
        <group>107</group>
        <mode>0744</mode>
        <label>virt_image_t</label>
      </permissions>
    </target>
  </volume>`

	pool, err := l.conn.LookupStoragePoolByName(poolName)
	if err != nil {
		return nil, err
	}

	volume, err := pool.StorageVolCreateXML(xmlConfig, 0)
	return &volume, err
}

func (l *LibvirtFacade) createPool(name string, path string) (*libvirt.VirStoragePool, error) {
	poolDefinition := `
  <pool type='dir'>
  	<name>` + name + `</name>
  	<source></source>
  	<target>
  		<path>` + path + `</path>
  			<permissions>
  			<mode>0755</mode>
  			<owner>-1</owner>
  			<group>-1</group>
  		</permissions>
  	</target>
  </pool>`

	pool, err := l.conn.StoragePoolDefineXML(poolDefinition, 0)
	pool.Create(0)

	if err != nil {
		return nil, err
	}

	return &pool, err
}
