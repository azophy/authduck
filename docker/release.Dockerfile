FROM alpine
COPY authduck /usr/bin/authduck
ENTRYPOINT ["/usr/bin/authduck"]
