package main

import (
	"context"
	"fmt"
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func createPod(submissionID string) {
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

	for i := 0; i < len(pod.Spec.Containers); i++ {
		pod.Spec.Containers[i].Env = append(pod.Spec.Containers[i].Env, apiv1.EnvVar{Name: "SUBMISSION_ID", Value: submissionID})
	}

	fmt.Printf("%#v\n", pod)

	pod, err = clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pod created successfully...")
}
