package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getService(client *kubernetes.Clientset, ns, name string) (*corev1.Service, error) {
	return client.CoreV1().Services(ns).Get(name, metav1.GetOptions{})
}

func isServiceReady(service *corev1.Service) bool {
	// if service type is LoadBalancer, wait until an IP is assigned
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) == 0 {
			return false
		}
		return true
	}

	// for ClusterIP, NodePort and ExternalName return true
	return true
}
