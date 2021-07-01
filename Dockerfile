FROM docker.io/golang:1.16.5 AS mqrunner
USER 0
RUN mkdir /go/src/mqrunner
WORKDIR /go/src/mqrunner
COPY mqrunner .
RUN go install ./cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4
USER 0

# install additional packages
RUN microdnf --disableplugin=subscription-manager install shadow-utils bc tar procps-ng gzip findutils pam passwd wget util-linux sudo libselinux-utils \
    bash ca-certificates file gawk glibc-common grep ncurses-compat-libs sed which openldap-clients openssl \
    libgcc libstdc++

# mqm user and group
RUN groupadd -g 1001 mqm && useradd -u 1001 -g 1001 -d /var/mqm -M -e "" -K PASS_MAX_DAYS=-1  mqm
RUN mkdir -p /etc/security && echo "mqm hard nofile 10240" >> /etc/security/limits.conf && echo "mqm soft nofile 10240" >> /etc/security/limits.conf

# copy rpm files
ARG RPMDIR=MQServer
RUN mkdir -p /tmp/MQServer
COPY $RPMDIR /tmp/MQServer

# install rpm's
RUN rpm -Uvh /tmp/MQServer/MQSeriesRuntime*.rpm \
    /tmp/MQServer/MQSeriesServer-*.rpm \
    /tmp/MQServer/MQSeriesGSKit*.rpm \
    /tmp/MQServer/MQSeriesJRE*.rpm \
    /tmp/MQServer/MQSeriesWeb*.rpm \
    /tmp/MQServer/MQSeriesJava*.rpm 

RUN microdnf clean all && rm -fr /tmp/MQServer 

# copy scripts
COPY scripts/*.sh .

# from this point common installation rpm/noi

# create directories && accept license
RUN ./create-directories.sh 1001 && ./accept-license.sh

# copy mqrunner from go builder
COPY --from=mqrunner /go/bin/cmd /opt/mqm/bin/mqrunner
RUN chown 1001:root /opt/mqm/bin/mqrunner & chmod u+x,g+x,o+x /opt/mqm/bin/mqrunner

EXPOSE 1414 9157 9443

USER 1001

#ENV MQ_EPHEMERAL_PREFIX=/etc/mqm
ENV MQ_OVERRIDE_DATA_PATH=/mnt/mqm/data 
ENV MQ_OVERRIDE_INSTALLATION_NAME=Installation1
ENV MQ_USER_NAME="mqm" 
ENV PATH="${PATH}:/opt/mqm/bin"

ENTRYPOINT ["/opt/mqm/bin/mqrunner"]
