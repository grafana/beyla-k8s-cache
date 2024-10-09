package meta

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/grafana/beyla-k8s-cache/pkg/meta/cni"
)

const (
	kubeConfigEnvVariable  = "KUBECONFIG"
	defaultResyncTime      = 30 * time.Minute
	defaultSyncTimeout     = 30 * time.Second
	indexPodByContainerIDs = "idx_pod_by_container"
	IndexIP                = "idx_ip"
	typeNode               = "Node"
	typePod                = "Pod"
	typeService            = "Service"
)

var podIndexers = cache.Indexers{
	indexPodByContainerIDs: func(obj interface{}) ([]string, error) {
		pi := obj.(*PodInfo)
		return pi.ContainerIDs, nil
	},
	IndexIP: func(obj interface{}) ([]string, error) {
		pi := obj.(*PodInfo)
		return pi.IPInfo.IPs, nil
	},
}

var serviceAndNodeIndexers = cache.Indexers{
	IndexIP: func(obj interface{}) ([]string, error) {
		return obj.(*IPInfo).IPs, nil
	},
}

func InitInformers(ctx context.Context, kubeconfigPath string, syncTimeout, resyncPeriod time.Duration) (*Informers, error) {
	log := slog.With("component", "kube.Informers")
	k := &Informers{log: log}

	kubeCfg, err := loadKubeconfig(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig can't be loaded: %w", err)
	}
	kubeClient, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		return nil, fmt.Errorf("kubernetes client can't be initialized: %w", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncPeriod)

	if err := k.initPodInformer(informerFactory); err != nil {
		return nil, err
	}
	if err := k.initNodeIPInformer(informerFactory); err != nil {
		return nil, err
	}
	if err := k.initServiceIPInformer(informerFactory); err != nil {
		return nil, err
	}

	k.log.Debug("starting kubernetes informers, waiting for syncronization")
	informerFactory.Start(ctx.Done())
	informerFactory.WaitForCacheSync(ctx.Done())
	k.log.Debug("kubernetes informers started")
	return k, nil
}

func loadKubeconfig(kubeConfigPath string) (*rest.Config, error) {
	// if no config path is provided, load it from the env variable
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv(kubeConfigEnvVariable)
	}
	// otherwise, load it from the $HOME/.kube/config file
	if kubeConfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("can't get user home dir: %w", err)
		}
		kubeConfigPath = path.Join(homeDir, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err == nil {
		return config, nil
	}
	// fallback: use in-cluster config
	config, err = rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("can't access kubenetes. Tried using config from: "+
			"config parameter, %s env, homedir and InClusterConfig. Got: %w",
			kubeConfigEnvVariable, err)
	}
	return config, nil
}

func (k *Informers) initPodInformer(informerFactory informers.SharedInformerFactory) error {
	pods := informerFactory.Core().V1().Pods().Informer()

	k.initContainerListeners(pods)

	// Transform any *v1.Pod instance into a *PodInfo instance to save space
	// in the informer's cache
	if err := pods.SetTransform(func(i interface{}) (interface{}, error) {
		pod, ok := i.(*v1.Pod)
		if !ok {
			// it's Ok. The K8s library just informed from an entity
			// that has been previously transformed/stored
			if pi, ok := i.(*PodInfo); ok {
				return pi, nil
			}
			return nil, fmt.Errorf("was expecting a PodInfo. Got: %T", i)
		}
		containerIDs := make([]string, 0,
			len(pod.Status.ContainerStatuses)+
				len(pod.Status.InitContainerStatuses)+
				len(pod.Status.EphemeralContainerStatuses))
		for i := range pod.Status.ContainerStatuses {
			containerIDs = append(containerIDs,
				rmContainerIDSchema(pod.Status.ContainerStatuses[i].ContainerID))
		}
		for i := range pod.Status.InitContainerStatuses {
			containerIDs = append(containerIDs,
				rmContainerIDSchema(pod.Status.InitContainerStatuses[i].ContainerID))
		}
		for i := range pod.Status.EphemeralContainerStatuses {
			containerIDs = append(containerIDs,
				rmContainerIDSchema(pod.Status.EphemeralContainerStatuses[i].ContainerID))
		}

		ips := make([]string, 0, len(pod.Status.PodIPs))
		for _, ip := range pod.Status.PodIPs {
			// ignoring host-networked Pod IPs
			if ip.IP != pod.Status.HostIP {
				ips = append(ips, ip.IP)
			}
		}

		owner := OwnerFrom(pod.OwnerReferences)
		startTime := pod.GetCreationTimestamp().String()
		if k.log.Enabled(context.TODO(), slog.LevelDebug) {
			k.log.Debug("inserting pod", "name", pod.Name, "namespace", pod.Namespace,
				"uid", pod.UID, "owner", owner,
				"node", pod.Spec.NodeName, "startTime", startTime,
				"containerIDs", containerIDs)
		}
		objectMeta := metav1.ObjectMeta{
			Name:            pod.Name,
			Namespace:       pod.Namespace,
			UID:             pod.UID,
			Labels:          pod.Labels,
			OwnerReferences: pod.OwnerReferences,
		}
		return &PodInfo{
			ObjectMeta:   objectMeta,
			Owner:        owner,
			NodeName:     pod.Spec.NodeName,
			StartTimeStr: startTime,
			ContainerIDs: containerIDs,
			// TODO: is this really needed?
			IPInfo: IPInfo{
				ObjectMeta: objectMeta,
				Kind:       typePod,
				HostIP:     pod.Status.HostIP,
				IPs:        ips,
				Owner:      owner,
			},
		}, nil
	}); err != nil {
		return fmt.Errorf("can't set pods transform: %w", err)
	}
	if err := pods.AddIndexers(podIndexers); err != nil {
		return fmt.Errorf("can't add indexers to Pods informer: %w", err)
	}

	k.pods = pods
	return nil
}

// initContainerListeners listens for deletions of pods, to forward them to the ContainerEventHandler subscribers.
func (k *Informers) initContainerListeners(pods cache.SharedIndexInformer) {
	if _, err := pods.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*PodInfo)
			k.log.Debug("deleting containers for pod", "pod", pod.Name, "containers", pod.ContainerIDs)
			for _, listener := range k.containerEventHandlers {
				listener.OnDeletion(pod.ContainerIDs)
			}
		},
	}); err != nil {
		k.log.Warn("can't attach container listener to the Kubernetes informer."+
			" Your kubernetes metadata might be outdated in the long term", "error", err)
	}
}

// rmContainerIDSchema extracts the hex ID of a container ID that is provided in the form:
// containerd://40c03570b6f4c30bc8d69923d37ee698f5cfcced92c7b7df1c47f6f7887378a9
func rmContainerIDSchema(containerID string) string {
	if parts := strings.Split(containerID, "://"); len(parts) > 1 {
		return parts[1]
	}
	return containerID
}

func (k *Informers) initNodeIPInformer(informerFactory informers.SharedInformerFactory) error {
	nodes := informerFactory.Core().V1().Nodes().Informer()
	// Transform any *v1.Node instance into a *IPInfo instance to save space
	// in the informer's cache
	if err := nodes.SetTransform(func(i interface{}) (interface{}, error) {
		node, ok := i.(*v1.Node)
		if !ok {
			// it's Ok. The K8s library just informed from an entity
			// that has been previously transformed/stored
			if pi, ok := i.(*IPInfo); ok {
				return pi, nil
			}
			return nil, fmt.Errorf("was expecting a IPInfo. Got: %T", i)
		}
		ips := make([]string, 0, len(node.Status.Addresses))
		for _, address := range node.Status.Addresses {
			ip := net.ParseIP(address.Address)
			if ip != nil {
				ips = append(ips, ip.String())
			}
		}
		// CNI-dependent logic (must work regardless of whether the CNI is installed)
		ips = cni.AddOvnIPs(ips, node)

		return &IPInfo{
			ObjectMeta: metav1.ObjectMeta{
				Name:      node.Name,
				Namespace: node.Namespace,
				Labels:    node.Labels,
			},
			IPs:  ips,
			Kind: typeNode,
		}, nil
	}); err != nil {
		return fmt.Errorf("can't set nodesIP transform: %w", err)
	}
	if err := nodes.AddIndexers(serviceAndNodeIndexers); err != nil {
		return fmt.Errorf("can't add %s indexer to Nodes informer: %w", IndexIP, err)
	}
	k.nodesIP = nodes
	return nil
}

func (k *Informers) initServiceIPInformer(informerFactory informers.SharedInformerFactory) error {
	services := informerFactory.Core().V1().Services().Informer()
	// Transform any *v1.Service instance into a *IPInfo instance to save space
	// in the informer's cache
	if err := services.SetTransform(func(i interface{}) (interface{}, error) {
		svc, ok := i.(*v1.Service)
		if !ok {
			// it's Ok. The K8s library just informed from an entity
			// that has been previously transformed/stored
			if pi, ok := i.(*IPInfo); ok {
				return pi, nil
			}
			return nil, fmt.Errorf("was expecting a IPInfo. Got: %T", i)
		}
		if svc.Spec.ClusterIP == v1.ClusterIPNone {
			// this will be normal for headless services
			k.log.Debug("Service doesn't have any ClusterIP. Beyla won't decorate their flows",
				"namespace", svc.Namespace, "name", svc.Name)
		}
		return &IPInfo{
			ObjectMeta: metav1.ObjectMeta{
				Name:            svc.Name,
				Namespace:       svc.Namespace,
				Labels:          svc.Labels,
				OwnerReferences: svc.OwnerReferences,
			},
			Kind: typeService,
			IPs:  svc.Spec.ClusterIPs,
		}, nil
	}); err != nil {
		return fmt.Errorf("can't set servicesIP transform: %w", err)
	}
	if err := services.AddIndexers(serviceAndNodeIndexers); err != nil {
		return fmt.Errorf("can't add %s indexer to Pods informer: %w", IndexIP, err)
	}

	k.servicesIP = services
	return nil
}
