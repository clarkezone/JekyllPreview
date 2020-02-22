FROM golang:alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go build

FROM jekyll/jekyll
USER root
RUN mkdir /app
COPY --from=builder /build/JekyllBlogPreview /app/.
WORKDIR /app
ADD startjek.sh .
RUN chmod +x startjek.sh
ENV JEKPREV_LOCALDIR=/srv/jekyll/source
ENV JEKPREV_monitorCmd=/app/startjek.sh

#env JEKPREV_REPO=https://github.com/clarkezone/clarkezone.github.io.git
#env JEKPREV_SECRET=ONETWOTHREE

ENTRYPOINT [ "/app/JekyllBlogPreview" ]

#docker run --rm -it -p=4000:4000 -p=8080:8080 clarkezone/jekpreview
