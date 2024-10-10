package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/beyla-k8s-cache/pkg/meta"
)

type observer struct {
	name string
}

func (o *observer) ID() string {
	return o.name
}

func (o *observer) Notify(event meta.Event) {
	fmt.Println(o.name, ":", toJSON(event))
}

func main() {
	//slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))

	infors, err := meta.InitInformers(context.Background(), "", time.Minute)
	if err != nil {
		panic(err)
	}

	infors.Subscribe(&observer{name: "observer1"})

	infors.Subscribe(&observer{name: "observer2"})

	time.Sleep(time.Second * 100000)
}

func toJSON(obj any) string {
	content, _ := json.Marshal(obj)
	return string(content)
}
