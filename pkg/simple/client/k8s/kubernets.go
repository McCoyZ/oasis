package k8s

import (
	"strings"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
	k8sclient "k8s.io/client-go/kubernetes"
	meshclient "zmc.io/oasis/pkg/client/clientset/versioned"
)

type Client interface {
	Kubernetes() k8sclient.Interface
	Mesh() meshclient.Interface
	Istio() istioclient.Interface
	Discovery() discovery.DiscoveryInterface
	Master() string
	Config() *rest.Config
}

type kubernetesClient struct {
	// kubernetes client interface
	k8s k8sclient.Interface

	// generated clientset
	ms meshclient.Interface

	istio istioclient.Interface

	// discovery client
	discoveryClient *discovery.DiscoveryClient

	master string

	config *rest.Config
}

// NewKubernetesClientOrDie creates KubernetesClient and panic if there is an error
func NewKubernetesClientOrDie(options *KubernetesOptions) *kubernetesClient {
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		panic(err)
	}

	config.QPS = options.QPS
	config.Burst = options.Burst

	k := &kubernetesClient{
		k8s:             k8sclient.NewForConfigOrDie(config),
		ms:              meshclient.NewForConfigOrDie(config),
		istio:           istioclient.NewForConfigOrDie(config),
		discoveryClient: discovery.NewDiscoveryClientForConfigOrDie(config),
		master:          config.Host,
		config:          config,
	}

	if options.Master != "" {
		k.master = options.Master
	}
	// The https prefix is automatically added when using sa.
	// But it will not be set automatically when reading from kubeconfig
	// which may cause some problems in the client of other languages.
	if !strings.HasPrefix(k.master, "http://") && !strings.HasPrefix(k.master, "https://") {
		k.master = "https://" + k.master
	}
	return k
}

// NewKubernetesClient creates a KubernetesClient
func NewKubernetesClient(options *KubernetesOptions) (Client, error) {
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, err
	}

	config.QPS = options.QPS
	config.Burst = options.Burst

	var k kubernetesClient
	k.k8s, err = k8sclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.ms, err = meshclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.istio, err = istioclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	k.master = options.Master
	k.config = config

	return &k, nil
}

func (k *kubernetesClient) Kubernetes() k8sclient.Interface {
	return k.k8s
}

func (k *kubernetesClient) Mesh() meshclient.Interface {
	return k.ms
}

func (k *kubernetesClient) Istio() istioclient.Interface {
	return k.istio
}

func (k *kubernetesClient) Discovery() discovery.DiscoveryInterface {
	return k.discoveryClient
}

// master address used to generate kubeconfig for downloading
func (k *kubernetesClient) Master() string {
	return k.master
}

func (k *kubernetesClient) Config() *rest.Config {
	return k.config
}
