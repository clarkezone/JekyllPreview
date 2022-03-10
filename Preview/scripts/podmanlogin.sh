# make sure CR_PAT is exported in profile with a valid PAT
echo $CR_PAT | podman login ghcr.io -u clarkezone --password-stdin

echo $DOCKER_PW | podman login registry.hub.docker.com/library -u clarkezone --password-stdin
