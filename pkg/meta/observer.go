package meta

type EventType int

const (
	Create = EventType(0)
	Update = EventType(1)
	Delete = EventType(2)
)

type Event struct {
	Type EventType
	Pod  *PodInfo
	IP   *IPInfo
}

type Observer interface {
	ID() string
	Notify(event Event)
}

type Notifier interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	Notify(event Event)
}
