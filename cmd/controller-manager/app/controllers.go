package app

import (
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"zmc.io/oasis/pkg/controller/destinationrule"
	"zmc.io/oasis/pkg/controller/virtualservice"
	"zmc.io/oasis/pkg/informers"
	"zmc.io/oasis/pkg/simple/client/k8s"
)

func addControllers(
	mgr manager.Manager,
	client k8s.Client,
	informerFactory informers.InformerFactory,
	serviceMeshEnabled bool,
	stopCh <-chan struct{}) error {

	kubernetesInformer := informerFactory.KubernetesSharedInformerFactory()
	istioInformer := informerFactory.IstioSharedInformerFactory()
	msInformer := informerFactory.MeshSharedInformerFactory()

	var vsController, drController manager.Runnable
	if serviceMeshEnabled {
		vsController = virtualservice.NewVirtualServiceController(kubernetesInformer.Core().V1().Services(),
			istioInformer.Networking().V1beta1().VirtualServices(),
			istioInformer.Networking().V1beta1().DestinationRules(),
			msInformer.Servicemesh().V1alpha1().Strategies(),
			client.Kubernetes(),
			client.Istio(),
			client.Mesh())

		drController = destinationrule.NewDestinationRuleController(kubernetesInformer.Apps().V1().Deployments(),
			istioInformer.Networking().V1beta1().DestinationRules(),
			kubernetesInformer.Core().V1().Services(),
			msInformer.Servicemesh().V1alpha1().ServicePolicies(),
			client.Kubernetes(),
			client.Istio(),
			client.Mesh())
	}

	controllers := map[string]manager.Runnable{
		"virtualservice-controller":  vsController,
		"destinationrule-controller": drController,
	}

	for name, ctrl := range controllers {
		if ctrl == nil {
			klog.V(4).Infof("%s is not going to run due to dependent component disabled.", name)
			continue
		}

		if err := mgr.Add(ctrl); err != nil {
			klog.Error(err, "add controller to manager failed", "name", name)
			return err
		}
	}

	return nil
}
