package main

import (
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getDeployment(client *kubernetes.Clientset, ns, name string) (*v1.Deployment, error) {
	return client.AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
}

// isDeploymentReady returns true if a deployment is completely ready and available.
func isDeploymentReady(deploy *v1.Deployment) bool {
	r := deploy.Status.Replicas
	return deploy.Status.AvailableReplicas == r &&
		deploy.Status.ReadyReplicas == r
}
