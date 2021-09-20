FROM ubuntu:20.10 as builder
ARG DEBIAN_FRONTEND=noninteractive
ENV GOROOT=/usr/local/go
ENV GOPATH=/uadmin
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get -y update \
    && apt-get -y install wget build-essential git-core golang npm libxml2-dev protobuf-compiler libprotobuf-dev \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/sergeyglazyrindev/uadmin
RUN mkdir /uadmin
RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz
RUN tar -xvf go1.16.4.linux-amd64.tar.gz
RUN mv go /usr/local
COPY . .
ARG GOPATH=/go
RUN make build

FROM ubuntu:20.10 as uadmin
ARG DEBIAN_FRONTEND=noninteractive
ENV GOROOT=/usr/local/go
ENV GOPATH=/uadmin
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get -y update \
    && apt-get -y install wget golang npm \
    && rm -rf /var/lib/apt/lists/*
#RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz
#RUN tar -xvf go1.16.4.linux-amd64.tar.gz
#RUN mv go /usr/local
COPY --from=builder /uadmin/uadmin /uadmin/uadmin
COPY configs/sqlite.yml /uadmin/configs/uadmin.yml
COPY configs/demo.yml /uadmin/configs/demo.yml
ENTRYPOINT ["/uadmin/uadmin", "admin", "serve"]