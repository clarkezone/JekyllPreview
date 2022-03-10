#set -x
export IMG=previewd
export VERSION=$(git describe --abbrev=0)

echo ${IMG}
echo ${VERSION}

podman manifest push ${IMG}:${VERSION} docker://ghcr.io/clarkezone/${IMG}:${VERSION}

echo podman search ghcr.io/clarkezone/${IMG} --list-tags
