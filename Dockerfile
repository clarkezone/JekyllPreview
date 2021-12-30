#Stage 1 build and test 
#docker.io prefix required by podman
FROM docker.io/golang:alpine as builder
RUN mkdir /build
ADD src/. /build
WORKDIR /build
RUN go mod tidy
RUN go build
RUN apk --no-cache add gcc build-base

# values required by unit / integration tests

#ARG ARG_TEST_JEKPREV_REPO_NOAUTH
#ARG ARG_TEST_JEKPREV_LOCALDIR
#ARG ARG_TEST_JEKPREV_BRANCHSWITCH
#ARG ARG_TEST_JEKPREV_SECURE_REPO_NOAUTH
#ARG ARG_TEST_JEKPREV_SECURECLONEPW
#
#ENV TEST_JEKPREV_REPO_NOAUTH $ARG_TEST_JEKPREV_REPO_NOAUTH
#ENV TEST_JEKPREV_LOCALDIR $ARG_TEST_JEKPREV_LOCALDIR
#ENV TEST_JEKPREV_BRANCHSWITCH $ARG_TEST_JEKPREV_BRANCHSWITCH
#ENV TEST_JEKPREV_SECURE_REPO_NOAUTH $ARG_TEST_JEKPREV_SECURE_REPO_NOAUTH
#ENV TEST_JEKPREV_SECURECLONEPW $ARG_TEST_JEKPREV_SECURECLONEPW

#RUN go test -v

# generate clean, final image for end users
FROM alpine:3.11.3
COPY --from=builder /build/JekyllBlogPreview .

# executable
ENTRYPOINT [ "./JekyllBlogPreview" ]
# arguments that can be overridden
#CMD [ "3", "300" ]
