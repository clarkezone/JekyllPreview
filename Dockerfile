FROM golang:alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go build

FROM jekyll/jekyll:ARM
USER root
RUN mkdir /app
COPY --from=builder /build/JekyllBlogPreview /app/.
WORKDIR /app
ADD startjek.sh .
RUN chmod +x startjek.sh
ENV JEKPREV_LOCALDIR=/srv/jekyll/source
ENV JEKPREV_monitorCmd=/app/startjek.sh

#Use these if you want to fork or customize or not use docker compose
#env JEKPREV_REPO=<YOUR REPO>
#env JEKPREV_SECRET=<YOUR SECRET>

ENTRYPOINT [ "/app/JekyllBlogPreview" ]

#to run manually from 
#docker run --rm -it -e JEKPREV_REPO=<repo> -e JEKPREV_SECRET=<secret> -p=4000:4000 -p=80:8080 clarkezone/jekpreview
