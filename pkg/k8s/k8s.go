package k8s

import (
	"fmt"
	"os"

	"github.com/prometheus/common/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	KUBERNETES_SERVICE_HOST = "KUBERNETES_SERVICE_HOST"
	KUBERNETES_SERVICE_PORT = "KUBERNETES_SERVICE_PORT"
)

type KubeClientConfig struct {
	KubeConfigPath string
	AllNamespaces  bool
	Namespace      string
}

type KubeClient struct {
	clientSet *kubernetes.Clientset
	kcConfig  *KubeClientConfig
}

func NewKubeClient(kcConfig *KubeClientConfig) (*KubeClient, error) {
	client := &KubeClient{
		kcConfig: kcConfig,
	}

	log.Infof("NewKubeClient: Fetching configuration")
	config, err := client.getConfig(kcConfig.KubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("NewKubeClient: %v", err)
	}

	log.Infof("NewKubeClient: Loading client")
	client.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return client, nil
}

func (c *KubeClient) getConfig(kubeConfigPath string) (*rest.Config, error) {
	var err error

	log.Infof("kubeConfigPath = %s", kubeConfigPath)

	runningInCluster := false
	if os.Getenv(KUBERNETES_SERVICE_HOST) != "" && os.Getenv(KUBERNETES_SERVICE_PORT) != "" {
		runningInCluster = true
	}

	log.Infof("runningInCluster: %v", runningInCluster)
	fmt.Printf("runningInCluster: %v", runningInCluster)

	config := &rest.Config{}
	if runningInCluster {
		log.Infof("KubeClient.getConfig: Running inside cluster")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("KubeClient.getConfig: rest.InClusterConfig: %v", err)
		}
	} else {
		log.Infof("KubeClient.getConfig: Running outside of cluster")
		_, err := os.Stat(kubeConfigPath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("KubeClient.getConfig: no configuration found")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("KubeClient.getConfig: clientcmd.BuildConfigFromFlags: %v", err)
		}
	}

	return config, nil
}

func (c *KubeClient) GetNamespaces() ([]string, error) {
	result, err := c.clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("KubeClient.GetNamespaces: List: %v", err)
	}

	allNamespaces := []string{}
	for _, item := range result.Items {
		allNamespaces = append(allNamespaces, item.Name)
	}

	return allNamespaces, nil
}

func (c *KubeClient) HasNamespace(name string) (bool, error) {
	result, err := c.clientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("Kubeclient.HasNamespace: Get: %v", err)
	}

	return result.Namespace == name, nil
}

func (c *KubeClient) GetServiceExternalIPs(ns string) ([]string, error) {
	result, err := c.clientSet.CoreV1().Services(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("KubeClient.GetServiceExternalIPs: List: %v", err)
	}

	allExternalIPs := []string{}
	for _, item := range result.Items {
		for _, ip := range item.Spec.ExternalIPs {
			allExternalIPs = append(allExternalIPs, ip)
		}
	}

	return allExternalIPs, nil
}
