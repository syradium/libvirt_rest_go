package worker

// CreateVPS creates VPS
func createVPS(name string, ram uint64, diskSize uint64) error {
	_, err := libvirtConn.CreateVPS(name, ram, diskSize)
	return err
}
