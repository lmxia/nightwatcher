package k8s

import (
	"errors"
	"github.com/lmxia/nightwatcher/utils"
	"log"
	"sync"
	"time"

	gaiaclientset "github.com/lmxia/gaia/pkg/generated/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type ClientManager struct {
	K8sClient  *kubernetes.Clientset
	Gaiaclient *gaiaclientset.Clientset
	RestConfig *restclient.Config
}

var once sync.Once
var k8sClient *ClientManager

func GetClientWithPanic() (*ClientManager, error) {
	once.Do(func() {
		var err error
		k8sClient, err = GetClient()
		if err != nil {
			log.Println("we can't get k8s client")
		}
	})
	if k8sClients == nil {
		return nil, errors.New("can't get k8s client")
	}
	return k8sClient, nil
}

var k8sClients = &sync.Map{} //并发map

func GetClient() (*ClientManager, error) {
	// By default we get in cluster config.
	localKubeConfig, err := utils.LoadsKubeConfig("/Users/xialingming/.kube/config", 1)
	if err != nil {
		return nil, err
	}
	localKubeClientSet := kubernetes.NewForConfigOrDie(localKubeConfig)
	localGaiaClientSet := gaiaclientset.NewForConfigOrDie(localKubeConfig)
	return &ClientManager{
		K8sClient:  localKubeClientSet,
		Gaiaclient: localGaiaClientSet,
	}, nil
}

const (
	// High enough QPS to fit all expected use cases.
	defaultQPS = 1e6
	// High enough Burst to fit all expected use cases.
	defaultBurst = 1e6
	// full resyc cache resource time
	defaultResyncPeriod = 30 * time.Second
)
