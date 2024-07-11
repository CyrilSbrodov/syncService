package kubernetes

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// KubernetesDeployer - структура дэплоера
type KubernetesDeployer struct {
	clientset *kubernetes.Clientset
}

// NewKubernetesDeployer - конструктор деплоера
func NewKubernetesDeployer() (*KubernetesDeployer, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("could not find kubeconfig")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubernetesDeployer{clientset: clientset}, nil
}

// CreatePod - создание нового pod'a с проверкой на уже существующий с таким же именем
func (d *KubernetesDeployer) CreatePod(name string) error {
	_, err := d.clientset.CoreV1().Pods("default").Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "algorithm-container",
						Image: "algorithm-image",
					},
				},
			},
		}
		_, err := d.clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
		return err
	}
	return nil
}

// DeletePod - удаление существующего pod'a, если такого нет, то выходит из функции
func (d *KubernetesDeployer) DeletePod(name string) error {
	_, err := d.clientset.CoreV1().Pods("default").Get(context.Background(), name, metav1.GetOptions{})
	if err == nil {
		return d.clientset.CoreV1().Pods("default").Delete(context.Background(), name, metav1.DeleteOptions{})
	}
	return nil
}

// GetPodList - функция получения всех pod
func (d *KubernetesDeployer) GetPodList() ([]string, error) {
	pods, err := d.clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}
