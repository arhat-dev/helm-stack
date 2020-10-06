ARG ARCH=amd64

FROM arhatdev/builder-go:alpine as builder

ARG ARCH=amd64

# Add kubctl v1.16,v1.17,v1.18,v1.19
RUN KUBE_ARCH="$(basename ${ARCH} v7)"; if [ "${KUBE_ARCH}" = "x86" ]; then KUBE_ARCH="386"; fi ;\
    export KUBE_ARCH ;\
    mkdir -p /opt/kube/v1.16 &&\
    curl -o /opt/kube/v1.16/kubectl -SsL "https://storage.googleapis.com/kubernetes-release/release/v1.16.15/bin/linux/${KUBE_ARCH}/kubectl" &&\
    chmod +x /opt/kube/v1.16/kubectl &&\
    mkdir -p /opt/kube/v1.17 &&\
    curl -o /opt/kube/v1.17/kubectl -SsL "https://storage.googleapis.com/kubernetes-release/release/v1.17.12/bin/linux/${KUBE_ARCH}/kubectl" &&\
    chmod +x /opt/kube/v1.17/kubectl &&\
    mkdir -p /opt/kube/v1.18 &&\
    curl -o /opt/kube/v1.18/kubectl -SsL "https://storage.googleapis.com/kubernetes-release/release/v1.18.9/bin/linux/${KUBE_ARCH}/kubectl" &&\
    chmod +x /opt/kube/v1.18/kubectl &&\
    mkdir -p /opt/kube/v1.19 &&\
    curl -o /opt/kube/v1.19/kubectl -SsL "https://storage.googleapis.com/kubernetes-release/release/v1.19.2/bin/linux/${KUBE_ARCH}/kubectl" &&\
    chmod +x /opt/kube/v1.19/kubectl

# Add helm v2,v3
RUN KUBE_ARCH="$(basename ${ARCH} v7)"; if [ "${KUBE_ARCH}" = "x86" ]; then KUBE_ARCH="386"; fi ;\
    export KUBE_ARCH ;\
    mkdir -p /tmp/helm /opt/helm/v2 &&\
    curl -SsL "https://get.helm.sh/helm-v2.16.12-linux-${KUBE_ARCH}.tar.gz" | tar -C /tmp/helm -zxf - &&\
    mv "/tmp/helm/linux-${KUBE_ARCH}/helm" /opt/helm/v2/helm &&\
    rm -rf /tmp/helm &&\
    mkdir -p /tmp/helm /opt/helm/v3 &&\
    curl -SsL "https://get.helm.sh/helm-v3.3.4-linux-${KUBE_ARCH}.tar.gz" | tar -C /tmp/helm -zxf - &&\
    mv "/tmp/helm/linux-${KUBE_ARCH}/helm" /opt/helm/v3/helm &&\
    rm -rf /tmp/helm

FROM arhatdev/go:alpine-${ARCH}
ARG APP=helm-stack

COPY --from=builder /opt /opt

ENTRYPOINT [ "sh", "-c" ]
CMD [ "/helm-stack" ]
