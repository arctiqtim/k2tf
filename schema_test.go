package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
	"testing"
)

func TestSchemaSupportsAttribute(t *testing.T) {
	type args struct {
		attrName string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"kubernetes_deployment.metadata",
			args{
				"kubernetes_deployment.metadata",
			},
			true,
			false,
		},
		{
			"kubernetes_pod.metadata.labels",
			args{
				"kubernetes_pod.metadata.labels",
			},
			true,
			false,
		},
		{
			"kubernetes_pod.metadata.name",
			args{
				"kubernetes_pod.metadata.name",
			},
			true,
			false,
		},
		{
			"kubernetes_deployment.spec.template.spec.container.name",
			args{
				"kubernetes_deployment.spec.template.spec.container.name",
			},
			true,
			false,
		},
		{
			"kubernetes_deployment.spec.toleration",
			args{
				"kubernetes_deployment.spec.toleration",
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsAttributeSupported(tt.args.attrName)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAttributeSupported() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsAttributeSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsKubernetesKindSupported(t *testing.T) {
	tests := []struct {
		name string
		g string
		v string
		k string
		want bool
	}{
		{"ConfigMap", "core", "v1", "ConfigMap", true},
		{"ClusterRole", "rbac.authorization.k8s.io", "v1", "ClusterRole", true},
		{"ClusterRoleBinding", "rbac.authorization.k8s.io", "v1", "ClusterRoleBinding", true},
		{"DaemonSet", "apps", "v1", "DaemonSet", true},
		{"Namespace", "core", "v1", "Namespace", true},
		{"Pod", "", "v1", "pod", true},
		{"PodDisruptionBudget_false", "policy", "v1beta1", "PodDisruptionBudget", false},
		{"Deployment", "apps", "v1", "deployment", true},
		{"Ingress_false", "extensions", "v1beta1", "ingress", false},
		{"ReplicaSet_false", "apps", "v1", "ReplicaSet", false},
		{"Secret", "core", "v1", "Secret", true},
		{"Service", "core", "v1", "Service", true},
		{"Endpoints", "core", "v1", "endpoints", true},
		{"ValidatingWebhookConfiguration_false", "admissionregistration.k8s.io", "v1beta1", "ValidatingWebhookConfiguration", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := testCreateRuntimeObject(t, tt.g, tt.v, tt.k)
			if got := IsKubernetesKindSupported(obj); got != tt.want {
				t.Errorf("IsResourceSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testCreateRuntimeObject(t *testing.T, g, v, k string) runtime.Object {
	apiVersion := v
	if g != "" && g != "core" {
		apiVersion = g + "/" + v
	}

	yaml := fmt.Sprintf(`
apiVersion: %s
kind: %s
metadata:
  name: test-%s
`, apiVersion, strcase.ToCamel(k), strings.ToLower(k))

	obj := testParseK8SYAML(t, yaml)
	if obj == nil {
		t.Fatal("test setup error, runtime.Object is nil")
	}
	return obj
}