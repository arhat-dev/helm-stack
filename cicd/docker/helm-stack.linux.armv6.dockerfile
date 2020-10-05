ARG ARCH=armv6

# there is no officially pre-built armv6 kubectl/helm binary available
FROM arhatdev/builder-go:alpine as builder
FROM arhatdev/go:debian-${ARCH}
ARG APP=helm-stack

ENTRYPOINT [ "sh", "-c" ]
CMD [ "/helm-stack" ]
