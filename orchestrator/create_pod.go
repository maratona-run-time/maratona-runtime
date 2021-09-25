package main

import (
	"fmt"
	"os"
	"context"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes/scheme"
)

func createPod() {
	fmt.Println("Trying to create Pod")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode

	yaml, _ := os.ReadFile("./pod.yml")
	fmt.Println(yaml)

	obj, _, err := decode(yaml, nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	pod := obj.(*apiv1.Pod)

	fmt.Printf("%#v\n", pod)

	pod, err = clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pod created successfully...")
}
