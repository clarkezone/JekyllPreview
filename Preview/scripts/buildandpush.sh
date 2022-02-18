#podman build . -t ghcr.io/clarkezone/jekyllblogpreview:latest
#podman push ghcr.io/clarkezone/jekyllblogpreview:latest
export IMG=previewd
#export VERSION=0.0.1.0
export VERSION=$(git describe --abbrev=0)

echo ${IMG}
echo ${VERSION}

podman build --arch=amd64 -t ${IMG}:${VERSION}.amd64 -f Dockerfile
podman build --arch=arm64 -t ${IMG}:${VERSION}.arm64 -f Dockerfile

podman manifest create ${IMG}:${VERSION}
podman manifest add ${IMG}:${VERSION} containers-storage:localhost/${IMG}:${VERSION}.amd64
podman manifest add ${IMG}:${VERSION} containers-storage:localhost/${IMG}:${VERSION}.arm64

podman manifest push ${IMG}:{$VERSION} docker:/registry.hub.docker.com/clarkezone/${IMG}:{$VERSION}
podman manifest push ${IMG}:{$VERSION} docker:/ghcr.io/clarkezone${IMG}:{$VERSION}
