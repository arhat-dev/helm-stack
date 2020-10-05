ARG ARCH=mips64le

# there is no officially pre-built mips64le kubectl/helm binary available
FROM arhatdev/builder-go:alpine as builder
FROM arhatdev/go:debian-${ARCH}
ARG APP=helm-stack

ENTRYPOINT [ "sh", "-c" ]
CMD [ "/helm-stack" ]
