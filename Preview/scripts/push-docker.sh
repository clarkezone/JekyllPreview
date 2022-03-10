#set -x
export IMG=previewd
export VERSION=$(git describe --abbrev=0)

echo ${IMG}
echo ${VERSION}

podman manifest push ${IMG}:${VERSION} docker://registry.hub.docker.com/clarkezone/${IMG}:${VERSION}

echo podman search registry.hub.docker.com/clarkezone/${IMG} --list-tags
