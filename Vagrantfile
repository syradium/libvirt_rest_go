# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
	config.vm.box = "debian/jessie64"
	config.vm.network "private_network", ip: "192.168.15.10"
	config.vm.provision "shell", privileged: true, inline: <<-EOF
	set -e
	apt-get update
	apt-get install -y qemu-kvm libvirt-bin

        systemctl enable lxc
	systemctl enable libvirtd
	EOF
	config.vm.synced_folder ".", "/vagrant", type: "virtualbox"
end
