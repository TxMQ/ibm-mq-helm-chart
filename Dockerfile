FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4
USER 0

RUN microdnf install -y shadow-utils bc tar procps-ng gzip findutils pam passwd wget util-linux sudo libselinux-utils
RUN microdnf install -y libgcc libstdc++

# mqm user and group, outside the range of auto assigned uid/gid
RUN groupadd -g 1001 mqm && useradd -u 1001 -g 1001 -d /var/mqm -M -e "" -K PASS_MAX_DAYS=-1  mqm
RUN mkdir -p /etc/security && echo "mqm hard nofile 10240" >> /etc/security/limits.conf && echo "mqm soft nofile 10240" >> /etc/security/limits.conf

RUN mkdir -p /tmp/MQServer
COPY MQServer /tmp/MQServer

RUN rpm -Uvh /tmp/MQServer/MQSeriesRuntime*.rpm /tmp/MQServer/MQSeriesServer-*.rpm /tmp/MQServer/MQSeriesGSKit*.rpm /tmp/MQServer/MQSeriesJRE*.rpm /tmp/MQServer/MQSeriesWeb*.rpm /tmp/MQServer/MQSeriesJava*.rpm && /opt/mqm/bin/mqlicense -accept

RUN microdnf clean all && rm -fr /tmp/MQServer 

# link /var/mqm to the /mnt/mqm/data; /var/mqm link is owned by root
RUN rm -fr /var/mqm && mkdir -mode 2775 -p /mnt/mqm/data && chown mqm:mqm /mnt/mqm/data && ln -s /mnt/mqm/data /var/mqm

# expose listening ports
# default mq port 1414, metrics 9157, web 9443
EXPOSE 1414 9157 9443

RUN mkdir -p /etc/mqm && chown mqm:mqm /etc/mqm

#ENV MQ_EPHEMERAL_PREFIX=/etc/mqm PATH="${PATH}:/opt/mqm/bin"
ENV PATH="${PATH}:/opt/mqm/bin"

RUN /opt/mqm/bin/crtmqdir -s -f
RUN /opt/mqm/bin/security/amqpamcf

USER 1001

ENTRYPOINT ["/bin/bash", "-c", "/usr/bin/sleep 100d"]

