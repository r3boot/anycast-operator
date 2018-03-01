package main

import (
	"flag"

	"os"

	"time"

	"github.com/prometheus/common/log"
	"github.com/r3boot/anycast-operator/pkg/k8s"
	"github.com/r3boot/anycast-operator/pkg/loopback"
	"github.com/r3boot/anycast-operator/pkg/utils"
)

var (
	cfgKubeConfig    = flag.String("kubeconfig", "", "Path to kubeconfig")
	cfgAllNamespaces = flag.Bool("all-namespaces", true, "Work on all namespaces")
	cfgNamespace     = flag.String("n", "", "Namespace to work on")
	cfgInterface     = flag.String("i", "dummy0", "Interface to configure IPs on")

	kubeClient  *k8s.KubeClient
	lbInterface *loopback.LoopbackInterface
)

func init() {
	var err error

	flag.Parse()

	kcConfig := &k8s.KubeClientConfig{}

	// Get path to kubeconfig
	if *cfgKubeConfig != "" {
		kcConfig.KubeConfigPath, err = utils.ExpandTilde(*cfgKubeConfig)
		if err != nil {
			log.Fatalf("init: %v", err)
		}
	} else if tmpPath := os.Getenv("KUBECONFIG"); tmpPath != "" {
		kcConfig.KubeConfigPath, err = utils.ExpandTilde(tmpPath)
		if err != nil {
			log.Fatalf("init: %v", err)
		}
	} else {
		kcConfig.KubeConfigPath, err = utils.ExpandTilde("~/.kube/config")
		if err != nil {
			log.Fatalf("init: %v", err)
		}
	}

	kubeClient, err = k8s.NewKubeClient(kcConfig)
	if err != nil {
		log.Fatalf("init: %v", err)
	}

	lbInterface, err = loopback.NewLoopback(&loopback.LoopbackInterfaceConfig{
		Interface: *cfgInterface,
	})
	if err != nil {
		log.Fatalf("init: %v", err)
	}
}

func main() {
	var err error

	namespaces := []string{}

	if *cfgAllNamespaces {
		namespaces, err = kubeClient.GetNamespaces()
		if err != nil {
			log.Fatalf("main: %v", err)
		}
	} else {
		ns := "kube-system"
		if *cfgNamespace != "" {
			ns = *cfgNamespace
		}
		hasNamespace, err := kubeClient.HasNamespace(ns)
		if err != nil {
			log.Fatalf("main: %v", err)
		}
		if !hasNamespace {
			log.Fatalf("Namespace %s does not exist", ns)
		}
		namespaces = append(namespaces, ns)
	}

	for {
		// Fetch ExternalIPs from configured Kubernetes Services
		externalIPs := []string{}
		for _, ns := range namespaces {
			ips, err := kubeClient.GetServiceExternalIPs(ns)
			if err != nil {
				log.Fatalf("main: %v", err)
			}

			for _, ip := range ips {
				externalIPs = append(externalIPs, ip)
			}
		}

		// Fetch the currently configured Anycast IPs
		currentIPs, err := lbInterface.GetAnycastIPs()
		if err != nil {
			log.Fatalf("main: %v", err)
		}

		// Pass 1, check for IP's to remove
		ipsToRemove := []string{}
		for _, currentIP := range currentIPs {
			isExternalIP := false
			for _, externalIP := range externalIPs {
				if currentIP == externalIP {
					isExternalIP = true
					break
				}
			}
			if !isExternalIP {
				ipsToRemove = append(ipsToRemove, currentIP)
			}
		}

		// Pass 2, check for IP's to add
		ipsToAdd := []string{}
		for _, externalIP := range externalIPs {
			isCurrentIP := false
			for _, currentIP := range currentIPs {
				if currentIP == externalIP {
					isCurrentIP = true
					break
				}
			}

			if !isCurrentIP {
				ipsToAdd = append(ipsToAdd, externalIP)
			}
		}

		if len(ipsToAdd) > 0 {
			for _, ip := range ipsToAdd {
				log.Infof("ip addr add %s/32 dev %s", ip, *cfgInterface)
			}

			log.Infof("Added %d ip addresses to %s", len(ipsToAdd), *cfgInterface)
		}

		if len(ipsToRemove) > 0 {
			for _, ip := range ipsToRemove {
				log.Infof("ip addr remove %s/32 dev %s", ip, *cfgInterface)
			}

			log.Infof("Removed %d ip addresses from %s", len(ipsToAdd), *cfgInterface)
		}

		time.Sleep(1 * time.Second)
	}
}
