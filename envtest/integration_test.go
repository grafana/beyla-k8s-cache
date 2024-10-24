package envtest

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

//	. "github.com/onsi/ginkgo/v2"
//	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	_ "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	_ "github.com/stretchr/testify/mock"

	_ "k8s.io/api/autoscaling/v2"

	_ "k8s.io/api/core/v1"

	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	_ "sigs.k8s.io/controller-runtime/pkg/envtest"
	_ "sigs.k8s.io/controller-runtime/pkg/log"
	_ "sigs.k8s.io/controller-runtime/pkg/log/zap"
	_ "sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/grafana/beyla-k8s-cache/pkg/meta"
	"github.com/grafana/beyla-k8s-cache/pkg/service"
)

var (
	ctx        context.Context
	k8sManager manager.Manager
	k8sClient  client.Client
	theClient  *kubernetes.Clientset
	testEnv    *envtest.Environment
	cancel     context.CancelFunc
)

func TestAPIs(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))

	ctx, cancel = context.WithCancel(context.TODO())

	testEnv = &envtest.Environment{}

	cfg, err := testEnv.Start()
	require.NoError(t, err)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	require.NoError(t, err)

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme.Scheme})
	require.NoError(t, err)

	svc := service.InformersCache{
		Port: 50055, // TODO: get it automatically
	}

	config := k8sManager.GetConfig()
	theClient, err = kubernetes.NewForConfig(config)
	require.NoError(t, err)

	err = k8sClient.Create(ctx, &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "test-container", Image: "nginx"},
			},
		},
	})
	require.NoError(t, err)

	go func() {
		require.NoError(t, svc.Run(ctx,
			meta.WithResyncPeriod(30*time.Minute),
			meta.WithKubeClient(theClient),
		))
	}()

	go func() {
		err = k8sManager.Start(ctx)
		require.NoError(t, err)
	}()

	/*
	   var _ = AfterSuite(func() {
	   	cancel()
	   	By("tearing down the test environment")
	   	err := testEnv.Stop()
	   	Expect(err).NotTo(HaveOccurred())
	   })
	*/

	time.Sleep(5 * time.Second)

	err = k8sClient.Create(ctx, &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-podacacao",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "test-container", Image: "nginx"},
			},
		},
	})
	//_, err := theClient.CoreV1().Pods("default").Create(ctx, &corev1.Pod{
	//}, v1.CreateOptions{})
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
}
