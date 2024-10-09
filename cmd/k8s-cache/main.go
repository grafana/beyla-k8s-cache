package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grafana/beyla-k8s-cache/pkg/meta"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))

	infors, err := meta.InitInformers(context.Background(), "", time.Minute, time.Minute)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.URL.Path)
		path := strings.Split(r.URL.Path, "/")
		switch path[1] {
		case "ip":
			if ipinfo, ok := infors.IPInfo(path[2]); ok {
				w.Write(toJSON(ipinfo))
				return
			}
		case "cnt":
			if podinfo, ok := infors.ContainerPod(path[2]); ok {
				w.Write(toJSON(podinfo))
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	panic(err)
}

func toJSON(obj any) []byte {
	content, _ := json.Marshal(obj)
	return content
}
