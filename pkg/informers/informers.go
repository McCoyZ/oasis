package informers

import (
	"time"

	istioinformers "istio.io/client-go/pkg/informers/externalversions"
	k8sinformers "k8s.io/client-go/informers"
	msinformers "zmc.io/oasis/pkg/client/informers/externalversions"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
	k8sclient "k8s.io/client-go/kubernetes"
	meshclient "zmc.io/oasis/pkg/client/clientset/versioned"
)

// default re-sync period for all informer factories
const defaultResync = 600 * time.Second

// InformerFactory is a group all shared informer factories which kubesphere needed
// callers should check if the return value is nil
type InformerFactory interface {
	KubernetesSharedInformerFactory() k8sinformers.SharedInformerFactory
	IstioSharedInformerFactory() istioinformers.SharedInformerFactory
	MeshSharedInformerFactory() msinformers.SharedInformerFactory

	// Start shared informer factory one by one if they are not nil
	Start(stopCh <-chan struct{})
}

type informerFactories struct {
	informerFactory      k8sinformers.SharedInformerFactory
	istioInformerFactory istioinformers.SharedInformerFactory
	msInformerFactory    msinformers.SharedInformerFactory
}

func NewInformerFactories(client k8sclient.Interface, msClient meshclient.Interface, istioClient istioclient.Interface) InformerFactory {
	factory := &informerFactories{}

	if client != nil {
		factory.informerFactory = k8sinformers.NewSharedInformerFactory(client, defaultResync)
	}

	if msClient != nil {
		factory.msInformerFactory = msinformers.NewSharedInformerFactory(msClient, defaultResync)
	}

	if istioClient != nil {
		factory.istioInformerFactory = istioinformers.NewSharedInformerFactory(istioClient, defaultResync)
	}

	return factory
}

func (f *informerFactories) KubernetesSharedInformerFactory() k8sinformers.SharedInformerFactory {
	return f.informerFactory
}

func (f *informerFactories) MeshSharedInformerFactory() msinformers.SharedInformerFactory {
	return f.msInformerFactory
}

func (f *informerFactories) IstioSharedInformerFactory() istioinformers.SharedInformerFactory {
	return f.istioInformerFactory
}

func (f *informerFactories) Start(stopCh <-chan struct{}) {
	if f.informerFactory != nil {
		f.informerFactory.Start(stopCh)
	}

	if f.msInformerFactory != nil {
		f.msInformerFactory.Start(stopCh)
	}

	if f.istioInformerFactory != nil {
		f.istioInformerFactory.Start(stopCh)
	}
}
