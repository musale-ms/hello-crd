package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/musale-ms/hello-crd/pkg/hello"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubecfg string

	// check the home directory exists and fill the path to the .kube/config
	if home := homedir.HomeDir(); home != "" {
		kubecfg = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		fmt.Println("Using the config in the cluster")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	helloResource := schema.GroupVersionResource{Group: "localusr.io", Version: "v1", Resource: "hello"}

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.Resource(helloResource).Namespace("").List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.Resource(helloResource).Namespace("").Watch(context.TODO(), options)
			},
		},
		&unstructured.Unstructured{},
		0,
		cache.Indexers{},
	)
	helloEventHandler := hello.HelloEventHandler{}

	informer.AddEventHandler(helloEventHandler.NewEvent())

	stop := make(chan struct{})
	defer close(stop)

	go informer.Run(stop)

	if !cache.WaitForCacheSync(stop, informer.HasSynced) {
		panic("Timeout waiting for the cache to sync")
	}

	fmt.Println("Custom Resource Controller started successfully")

	<-stop
}
