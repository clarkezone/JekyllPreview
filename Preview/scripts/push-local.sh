#set -x
export IMG=previewd
export VERSION=$(git describe --abbrev=0)

echo ${IMG}
echo ${VERSION}

podman manifest push ${IMG}:${VERSION} docker://registry.dev.clarkezone.dev/${IMG}:${VERSION}

echo podman search registry.dev.clarkezone.dev/${IMG} --list-tags
