#sudo docker run --name jbuild --rm -it -v ~/src/synology/clarkezoneblog:/src jekyllbuilder

#Podman with volumes don't appear to work on M1 mac
#If podbuild builds fails with not valid until date, restart machine

#FROM arm64v8/ruby:latest
#FROM debian
FROM docker.io/ruby:latest

RUN apt-get update\
&& apt-get install -y --no-install-recommends apt-utils curl ca-certificates build-essential gcc libssl-dev
RUN curl -Lo cmake-3.15.2.tar.gz https://github.com/Kitware/CMake/releases/download/v3.15.2/cmake-3.15.2.tar.gz
RUN tar xzf cmake-3.15.2.tar.gz && rm cmake-3.15.2.tar.gz

RUN cd cmake-3.15.2 && ls && ./configure --prefix=/opt/cmake  && make
RUN  cd cmake-3.15.2 &&  make install

RUN gem install \
  github-pages \
  jekyll \
  jekyll-redirect-from \
  kramdown \
  rdiscount \
  rouge
VOLUME /src
EXPOSE 4000
WORKDIR /src
RUN cd /bin && ln -s /opt/cmake/bin/cmake cmake

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8
ENV JEKYLL_ENV production
ENTRYPOINT ["bash"]
