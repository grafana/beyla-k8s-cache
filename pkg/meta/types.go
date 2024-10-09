package meta

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Provider interface {
	ContainerPod(containerID string) (*PodInfo, bool)
	IPInfo(ip string) (*IPInfo, bool)
}

// PodInfo contains precollected metadata for Pods.
type PodInfo struct {
	// Informers need that internal object is an ObjectMeta instance
	metav1.ObjectMeta
	NodeName string

	Owner *Owner

	// StartTimeStr caches value of ObjectMeta.StartTimestamp.String()
	StartTimeStr string
	ContainerIDs []string
	IPInfo       IPInfo
}

// IPInfo contains precollected metadata for Pods, Nodes and Services.
// Not all the fields are populated for all the above types. To save
// memory, we just keep in memory the necessary data for each Type.
// For more information about which fields are set for each type, please
// refer to the instantiation function of the respective informers.
type IPInfo struct {
	metav1.ObjectMeta
	Kind  string
	Owner Owner
	IPs   []string
	// Hostname and HostIP would be required only for Pod's IPInfo
	HostName string
	HostIP   string
}

type Owner struct {
	Kind string
	Name string
	// Owner of the owner. For example, a ReplicaSet might be owned by a Deployment
	Owner *Owner
}

// OwnerFrom returns the most plausible Owner reference. It might be
// null if the entity does not have any owner
func OwnerFrom(orefs []metav1.OwnerReference) *Owner {
	// fallback will store any found owner that is not part of the bundled
	// K8s owner types (e.g. argocd rollouts).
	// It will be returned if any of the standard K8s owners are found
	var fallback *Owner
	for i := range orefs {
		or := &orefs[i]
		if or.APIVersion != "apps/v1" {
			fallback = &Owner{Name: or.Name, Kind: or.Kind}
			continue
		}
		return &Owner{Name: or.Name, Kind: or.Kind}
	}
	return fallback
}

// TopOwner returns the top Owner in the owner chain.
// For example, if the owner is a ReplicaSet, it will return the Deployment name.
func (o *Owner) TopOwner() *Owner {
	// we have two levels of ownership at most
	if o != nil && o.Kind == "ReplicaSet" && o.Owner == nil {
		// we heuristically extract the Deployment name from the replicaset name
		if idx := strings.LastIndexByte(o.Name, '-'); idx > 0 {
			o.Owner = &Owner{
				Name: o.Name[:idx],
				Kind: "Deployment",
			}
		} else {
			// just caching the own replicaset as owner, in order to cache the result
			o.Owner = o
		}
		return o.Owner
	}

	// just return the highest existing owner (two levels of ownership maximum)
	if o == nil || o.Owner == nil {
		return o
	}
	return o.Owner
}

func (o *Owner) String() string {
	sb := strings.Builder{}
	o.string(&sb)
	return sb.String()
}

func (o *Owner) string(sb *strings.Builder) {
	if o.Owner != nil {
		o.Owner.string(sb)
		sb.WriteString("->")
	}
	sb.WriteString(o.Kind)
	sb.WriteByte(':')
	sb.WriteString(o.Name)
}
