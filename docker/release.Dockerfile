FROM ubuntu:22.04
COPY authduck /usr/bin/authduck
ENTRYPOINT ["/usr/bin/authduck"]
