package controllers

import (
	"context"
	"reflect"
	"testing"

	stroomv1 "github.com/gradata-systems/stroom-k8s-operator/api/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestCreateIngressesPathType(t *testing.T) {
	tests := []struct {
		name     string
		override netv1.PathType
		expected map[string]netv1.PathType
	}{
		{
			name: "preserves route defaults when unset",
			expected: map[string]netv1.PathType{
				"/":                       netv1.PathTypePrefix,
				"/stroom/noauth/datafeed": netv1.PathTypeExact,
				"/stroom/datafeeddirect":  netv1.PathTypeExact,
				"/web-socket/":            netv1.PathTypePrefix,
			},
		},
		{
			name:     "overrides every route when configured",
			override: netv1.PathTypeImplementationSpecific,
			expected: map[string]netv1.PathType{
				"/":                       netv1.PathTypeImplementationSpecific,
				"/stroom/noauth/datafeed": netv1.PathTypeImplementationSpecific,
				"/stroom/datafeeddirect":  netv1.PathTypeImplementationSpecific,
				"/web-socket/":            netv1.PathTypeImplementationSpecific,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			if err := stroomv1.AddToScheme(scheme); err != nil {
				t.Fatal(err)
			}
			reconciler := StroomClusterReconciler{Scheme: scheme}
			cluster := &stroomv1.StroomCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "example"},
				Spec: stroomv1.StroomClusterSpec{
					Ingress: stroomv1.IngressSettings{
						HostName: "stroom.example.com",
						PathType: tt.override,
					},
					NodeSets: []stroomv1.NodeSet{{Name: "data", IngressEnabled: true}},
				},
			}

			got := map[string]netv1.PathType{}
			for _, ingress := range reconciler.createIngresses(context.Background(), cluster) {
				for _, rule := range ingress.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						got[path.Path] = *path.PathType
					}
				}
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("unexpected path types: got %v, want %v", got, tt.expected)
			}
		})
	}
}
