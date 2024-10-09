package meta

import (
	"log/slog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// containerEventHandler listens for the deletion of containers, as triggered
// by a Pod deletion.
type containerEventHandler interface {
	OnDeletion(containerID []string)
}

type Informers struct {
	log *slog.Logger
	// pods and replicaSets cache the different K8s types to custom, smaller object types
	pods       cache.SharedIndexInformer
	nodesIP    cache.SharedIndexInformer
	servicesIP cache.SharedIndexInformer

	containerEventHandlers []containerEventHandler
}

func (r *Informers) ContainerPod(containerID string) (*PodInfo, bool) {
	objs, err := r.pods.GetIndexer().ByIndex(indexPodByContainerIDs, containerID)
	if err != nil {
		r.log.Debug("error accessing index by container ID. Ignoring", "error", err, "containerID", containerID)
		return nil, false
	}
	if len(objs) == 0 {
		return nil, false
	}
	return objs[0].(*PodInfo), true
}

func (r *Informers) IPInfo(ip string) (*IPInfo, bool) {
	if info, ok := r.fetchInformersByIP(ip); ok {
		// Owner data might be discovered after the owned, so we fetch it
		// at the last moment
		if info.Owner.Name == "" {
			info.Owner = r.getOwner(&info.ObjectMeta, info)
		}
		return info, true
	}

	return nil, false
}

func (r *Informers) fetchInformersByIP(ip string) (*IPInfo, bool) {
	if info, ok := r.infoForIP(r.pods.GetIndexer(), ip); ok {
		info := info.(*PodInfo)
		// it might happen that the Host is discovered after the Pod
		if info.IPInfo.HostName == "" {
			info.IPInfo.HostName = r.getHostName(info.IPInfo.HostIP)
		}
		return &info.IPInfo, true
	}
	if info, ok := r.infoForIP(r.nodesIP.GetIndexer(), ip); ok {
		return info.(*IPInfo), true
	}
	if info, ok := r.infoForIP(r.servicesIP.GetIndexer(), ip); ok {
		return info.(*IPInfo), true
	}
	return nil, false
}

func (r *Informers) infoForIP(idx cache.Indexer, ip string) (any, bool) {
	objs, err := idx.ByIndex(IndexIP, ip)
	if err != nil {
		r.log.Debug("error accessing index. Ignoring", "ip", ip, "error", err)
		return nil, false
	}
	if len(objs) == 0 {
		return nil, false
	}
	return objs[0], true
}

func (r *Informers) getOwner(meta *metav1.ObjectMeta, info *IPInfo) Owner {
	if len(meta.OwnerReferences) > 0 {
		return *OwnerFrom(meta.OwnerReferences)
	}
	// If no owner references found, return itself as owner
	return Owner{
		Name: meta.Name,
		Kind: info.Kind,
	}
}

func (r *Informers) getHostName(hostIP string) string {
	if info, ok := r.infoForIP(r.nodesIP.GetIndexer(), hostIP); ok {
		return info.(*IPInfo).Name
	}
	return ""
}
