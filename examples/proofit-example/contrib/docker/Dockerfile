ARG  BASE=proofit/proofit:release
FROM $BASE
ARG  ARCH=amd64
ENV environment=proofit-example
ENV UADMIN_PATH=/proofit-example
ENV GOROOT=/usr/local/go
ENV GOPATH=/proofit-example
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
RUN apt-get -y update \
    && apt-get -y install libxml2 \
    && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /proofit-example
COPY proofit.$ARCH /proofit-example/proofit-example
RUN chmod +x /proofit-example/proofit-example
RUN mkdir -p /proofit-example/uploads
EXPOSE 8080
COPY proofit-example.yml /proofit-example/configs/proofit-example.yml
ENTRYPOINT ["/proofit-example/proofit-example", "admin", "serve"]
