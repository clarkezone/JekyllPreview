podman machine ssh < EOF
	sudo -i
	rpm-ostree install qemu-user-static
	systemctl reboot
EOF
