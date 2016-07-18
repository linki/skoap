FROM registry.opensource.zalan.do/stups/ubuntu:16.04-35

COPY entrypoint.sh /

COPY build/linux/skoap /skoap
COPY scm-source.json /scm-source.json

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/skoap"]
