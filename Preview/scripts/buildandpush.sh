#podman build . -t ghcr.io/clarkezone/jekyllblogpreview:latest
#podman push ghcr.io/clarkezone/jekyllblogpreview:latest
podman build . -t registry.dev.clarkezone.dev/jekyllblogpreview:latest
podman push registry.dev.clarkezone.dev/jekyllblogpreview:latest
