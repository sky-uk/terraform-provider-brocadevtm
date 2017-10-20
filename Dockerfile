# Use phusion/baseimage as base image. To make your builds
# reproducible, make sure you lock down to a specific version, not
# to `latest`! See
# https://github.com/phusion/baseimage-docker/blob/master/Changelog.md
# for a list of version numbers.
FROM phusion/baseimage:0.9.22

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]

# Install some packages.
RUN apt-get update; apt-get -y install aptitude wget unzip git

# Install Golang binaries
RUN wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
RUN tar -xvf go1.9.1.linux-amd64.tar.gz; mv go /usr/local/

# Install Terraform.
RUN mkdir /terraform;
WORKDIR /terraform
RUN wget https://releases.hashicorp.com/terraform/0.10.7/terraform_0.10.7_linux_amd64.zip
RUN unzip terraform_0.10.7_linux_amd64.zip
WORKDIR /
RUN mv /terraform /usr/local/

ENV PATH $PATH:/usr/local/terraform:/usr/local/go/bin

# Setup golang deps
RUN mkdir -p /gows
ENV GOPATH /gows
RUN go get -u github.com/tools/godep
RUN go get -u github.com/golang/lint
RUN go get -u github.com/axw/gocov
RUN go get -u github.com/AlekSi/gocov-xml
RUN go get -u github.com/matm/gocov-html
RUN go get -u github.com/go-playground/overalls

RUN apt-get -y install make binutils

# Build the Brocade vTM provider
#ADD . /gows/src/github.com/sky-uk/terraform-provider-brocadevtm
#RUN cd /gows/src/github.com/sky-uk/terraform-provider-brocadevtm; make fmt; make ; cp /gows/bin/terraform-provider-brocadevtm /usr/local/terraform/

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.ssh/id_rsa*
