#Stage 1 build and test 

FROM golang:alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go build
RUN apk --no-cache add gcc build-base

# values required by unit / integration tests

ARG ARG_TEST_JEKPREV_REPO_NOAUTH
ARG ARG_TEST_JEKPREV_LOCALDIR
ARG ARG_TEST_JEKPREV_BRANCHSWITCH
ARG ARG_TEST_JEKPREV_SECURE_REPO_NOAUTH
ARG ARG_TEST_JEKPREV_SECURECLONEPW

ENV TEST_JEKPREV_REPO_NOAUTH $ARG_TEST_JEKPREV_REPO_NOAUTH
ENV TEST_JEKPREV_LOCALDIR $ARG_TEST_JEKPREV_LOCALDIR
ENV TEST_JEKPREV_BRANCHSWITCH $ARG_TEST_JEKPREV_BRANCHSWITCH
ENV TEST_JEKPREV_SECURE_REPO_NOAUTH $ARG_TEST_JEKPREV_SECURE_REPO_NOAUTH
ENV TEST_JEKPREV_SECURECLONEPW $ARG_TEST_JEKPREV_SECURECLONEPW

RUN go test -v

#Stage 2 add to jekyll image

#The clarkezone/jekyll:x64 image doesn't currently work.  The ARM one does.
FROM jekyll/jekyll
#FROM clarkezone/jekyll:ARM
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

# to build
#docker build . \
# --build-arg ARG_TEST_JEKPREV_REPO_NOAUTH=value \
# --build-arg ARG_TEST_JEKPREV_LOCALDIR=value \
# --build-arg ARG_TEST_JEKPREV_BRANCHSWITCH=value \
# --build-arg ARG_TEST_JEKPREV_SECURE_REPO_NOAUTH=value \
# --build-arg ARG_TEST_JEKPREV_SECURECLONEPW=value \
