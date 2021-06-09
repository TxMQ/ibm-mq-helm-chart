FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4
USER 0

COPY MQServer /tmp
RUN ls /tmp
