package controllers

import (
	"context"
	"testing"

	stroomv1 "github.com/gradata-systems/stroom-k8s-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestDatafeedIngressIncludesHAProxyRewrite(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := stroomv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	reconciler := StroomClusterReconciler{Scheme: scheme}
	cluster := &stroomv1.StroomCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default"},
		Spec: stroomv1.StroomClusterSpec{
			Ingress: stroomv1.IngressSettings{
				HostName:   "stroom.example.com",
				SecretName: "tls",
				ClassName:  "haproxy",
			},
			NodeSets: []stroomv1.NodeSet{{
				Name:           "backend",
				Role:           stroomv1.ProcessingNodeRole,
				IngressEnabled: true,
			}},
		},
	}

	for _, ingress := range reconciler.createIngresses(context.Background(), cluster) {
		if ingress.Name == cluster.GetBaseName()+"-datafeed" {
			got := ingress.Annotations["haproxy.org/path-rewrite"]
			want := "/stroom/datafeeddirect /stroom/noauth/datafeed"
			if got != want {
				t.Fatalf("unexpected HAProxy rewrite: got %q, want %q", got, want)
			}
			if ingress.Annotations["nginx.ingress.kubernetes.io/rewrite-target"] != "/stroom/noauth/datafeed" {
				t.Fatal("nginx rewrite annotation changed")
			}
			return
		}
	}
	t.Fatal("datafeed ingress not generated")
}
