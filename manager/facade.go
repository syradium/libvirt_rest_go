package manager

import "github.com/rgbkrk/libvirt-go"

// LibvirtFacade is a facade over libvirt for easier VPS management
type LibvirtFacade struct {
	conn *libvirt.VirConnection
}

// NewLibvirtFacade is a simple constructor for LibvirtFacade
func NewLibvirtFacade(uri string) (*LibvirtFacade, error) {
	conn, err := libvirt.NewVirConnection(uri)
	if err != nil {
		return nil, err
	}
	facade := LibvirtFacade{conn: &conn}
	return &facade, nil
}

// Close closes connection to libvirt daemon
func (l *LibvirtFacade) Close() (int, error) {
	return l.conn.CloseConnection()
}

func (l *LibvirtFacade) createVolume(poolName string, volumeName string, sizeG float64) (*libvirt.VirStorageVol, error) {
	xmlConfig := `<volume> <name>` + volumeName + `</name>
    <allocation>0</allocation>
    <capacity unit="G">2</capacity>
    <target>
      <path>` + volumeName + `.img</path>
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

func (l *LibvirtFacade) createPool(name string) (*libvirt.VirStoragePool, error) {
	poolDefinition :=
		`
  <pool type='dir'>
  <name>` + name + `</name>
  <uuid>8c79f996-cb2a-d24d-9822-ac7547ab2d01</uuid>
  <capacity unit='bytes'>4306780815</capacity>
  <allocation unit='bytes'>237457858</allocation>
  <available unit='bytes'>4069322956</available>
  <source>
  </source>
  <target>
  <path>/root/images</path>
  <permissions>
  <mode>0755</mode>
  <owner>-1</owner>
  <group>-1</group>
  </permissions>
  </target>
  </pool>
  `

	pool, err := l.conn.StoragePoolDefineXML(poolDefinition, 0)
	pool.Create(0)

	if err != nil {
		return nil, err
	}

	return &pool, err
}
