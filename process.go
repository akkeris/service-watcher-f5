package main

import (
	k8sconfig "github.com/akkeris/service-watcher-f5/k8sconfig"
	services "github.com/akkeris/service-watcher-f5/services"
	utils "github.com/akkeris/service-watcher-f5/utils"
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	api "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"time"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/stackimpact/stackimpact-go"
	"os"
)

func main() {

	utils.Variableinit()
	utils.Startclient()
	utils.InitBlacklist()
	if os.Getenv("PROFILE") == "true" {
		fmt.Println("Starting profiler...")
		_ = stackimpact.Start(stackimpact.Options{
			AgentKey:       os.Getenv("STACKIMPACT"),
			AppName:        "Service Watcher",
			AppEnvironment: os.Getenv("CLUSTER"),
		})

	}

	k8sconfig.CreateConfig()
	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	utils.Client = clientset.CoreV1().RESTClient()
	listWatch := cache.NewListWatchFromClient(
		utils.Client, "services", "",
		fields.Everything())

	listWatch.ListFunc = func(options api.ListOptions) (runtime.Object, error) {
		return utils.Client.Get().Namespace("none").Resource("services").Do().Get()
	}
	listWatch.WatchFunc = func(options api.ListOptions) (watch.Interface, error) {
		return clientset.CoreV1().Services(api.NamespaceAll).Watch(v1.ListOptions{})
	}

	_, controller := cache.NewInformer(
		listWatch, &corev1.Service{},
		time.Second*0, cache.ResourceEventHandlerFuncs{
			AddFunc:    printEventAdd,
			DeleteFunc: printEventDelete,
		},
	)
	fmt.Println("Watching for changes in Services....")
	controller.Run(wait.NeverStop)
}

func printEventAdd(obj interface{}) {
	_, isService := obj.(*corev1.Service)
	if isService {

		created := obj.(*corev1.Service).ObjectMeta.CreationTimestamp.Unix()
		now := v1.Now().Unix()

		diff := now - created
		if diff < 300 && !Blacklisted(obj.(*corev1.Service).ObjectMeta.Namespace) {
			fmt.Println("ADD")
			var jsn []byte
			var err error

			jsn, err = json.Marshal(obj)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(jsn))
			services.ProcessServiceAdd(obj)
		}
	}
}

func Blacklisted(namespace string) bool {

	return utils.Blacklist[namespace]

}

func printEventDelete(obj interface{}) {
	fmt.Println("DELETE")
	_, isService := obj.(*corev1.Service)
	if isService {
		if !Blacklisted(obj.(*corev1.Service).ObjectMeta.Namespace) {
			var jsn []byte
			var err error
			jsn, err = json.Marshal(obj)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(jsn))
			services.ProcessServiceDelete(obj)
		}
	}

}
