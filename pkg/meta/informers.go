package meta

import (
	"fmt"
	"log/slog"
	"sync"

	"k8s.io/client-go/tools/cache"
)

type Informers struct {
	log *slog.Logger
	// pods and replicaSets cache the different K8s types to custom, smaller object types
	pods     cache.SharedIndexInformer
	nodes    cache.SharedIndexInformer
	services cache.SharedIndexInformer

	// notifier implementation
	mutex     sync.RWMutex
	observers map[string]Observer
}

func (i *Informers) Subscribe(observer Observer) {
	i.mutex.Lock()
	i.observers[observer.ID()] = observer
	i.mutex.Unlock()
	// as a "welcome" message, we send the whole kube metadata to the new observer
	go func() {
		for _, pod := range i.pods.GetStore().List() {
			observer.Notify(Event{Type: Create, Pod: pod.(*PodInfo)})
		}
		for _, node := range i.nodes.GetStore().List() {
			observer.Notify(Event{Type: Create, IP: node.(*IPInfo)})
		}
		for _, service := range i.services.GetStore().List() {
			observer.Notify(Event{Type: Create, IP: service.(*IPInfo)})
		}
	}()
}

func (i *Informers) Unsubscribe(observer Observer) {
	i.mutex.Lock()
	delete(i.observers, observer.ID())
	i.mutex.Unlock()
}

func (i *Informers) Notify(event Event) {
	fmt.Println("Notifying observers")
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	for _, observer := range i.observers {
		fmt.Println("Notifying observer", observer.ID())
		observer.Notify(event)
	}
}
