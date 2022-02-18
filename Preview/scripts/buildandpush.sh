#podman build . -t ghcr.io/clarkezone/jekyllblogpreview:latest
#podman push ghcr.io/clarkezone/jekyllblogpreview:latest
export IMG=previewd
export VERSION=0.0.1.0

podman build --arch=amd64 -t ${IMG}:${VERSION}.amd64 -f Dockerfile
podman build --arch=arm64 -t ${IMG}:${VERSION}.arm64 -f Dockerfile

podman manifest create ${IMG}:${VERSION}
podman manifest add ${IMG}:${VERSION} ${IMG}:${VERSION}.amd64
podman manifest add ${IMG}:${VERSION} ${IMG}:${VERSION}.arm64


#podman build . -t registry.dev.clarkezone.dev/jekyllblogpreview:latest
#podman push registry.dev.clarkezone.dev/jekyllblogpreview:latest
