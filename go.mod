module arhat.dev/helm-stack

go 1.15

require (
	arhat.dev/pkg v0.5.4-0.20201208233302-107b8822e93b
	github.com/rogpeppe/go-internal v1.6.2
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/multierr v1.6.0
	k8s.io/apimachinery v0.19.4
	k8s.io/kubectl v0.19.4
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/creack/pty => github.com/jeffreystoke/pty v1.1.12-0.20201126201855-c1c1e24408f9

replace (
	k8s.io/api => github.com/kubernetes/api v0.19.4
	k8s.io/apiextensions-apiserver => github.com/kubernetes/apiextensions-apiserver v0.19.4
	k8s.io/apimachinery => github.com/kubernetes/apimachinery v0.19.4
	k8s.io/apiserver => github.com/kubernetes/apiserver v0.19.4
	k8s.io/cli-runtime => github.com/kubernetes/cli-runtime v0.19.4
	k8s.io/client-go => github.com/kubernetes/client-go v0.19.4
	k8s.io/cloud-provider => github.com/kubernetes/cloud-provider v0.19.4
	k8s.io/cluster-bootstrap => github.com/kubernetes/cluster-bootstrap v0.19.4
	k8s.io/code-generator => github.com/kubernetes/code-generator v0.19.4
	k8s.io/component-base => github.com/kubernetes/component-base v0.19.4
	k8s.io/cri-api => github.com/kubernetes/cri-api v0.19.4
	k8s.io/csi-translation-lib => github.com/kubernetes/csi-translation-lib v0.19.4
	k8s.io/klog => github.com/kubernetes/klog v1.0.0
	k8s.io/klog/v2 => github.com/kubernetes/klog/v2 v2.4.0
	k8s.io/kube-aggregator => github.com/kubernetes/kube-aggregator v0.19.4
	k8s.io/kube-controller-manager => github.com/kubernetes/kube-controller-manager v0.19.4
	k8s.io/kube-proxy => github.com/kubernetes/kube-proxy v0.19.4
	k8s.io/kube-scheduler => github.com/kubernetes/kube-scheduler v0.19.4
	k8s.io/kubectl => github.com/kubernetes/kubectl v0.19.4
	k8s.io/kubelet => github.com/kubernetes/kubelet v0.19.4
	k8s.io/kubernetes => github.com/kubernetes/kubernetes v1.19.4
	k8s.io/legacy-cloud-providers => github.com/kubernetes/legacy-cloud-providers v0.19.4
	k8s.io/metrics => github.com/kubernetes/metrics v0.19.4
	k8s.io/sample-apiserver => github.com/kubernetes/sample-apiserver v0.19.4
	k8s.io/utils => github.com/kubernetes/utils v0.0.0-20201027101359-01387209bb0d
	vbom.ml/util => github.com/fvbommel/util v0.0.2
)
