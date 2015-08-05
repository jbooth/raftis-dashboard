package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kubeclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"net/http"
)

var ()

func main() {
	//client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
	http.HandleFunc("/listPods", listPods)
	http.HandleFunc("/getEtcdNode", getEtcdNode)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listPods(w http.ResponseWriter, r *http.Request) {
	config := &kubeclient.Config{
		//Host:     "http://10.65.224.102:8080",
		Host:     "http://172.20.2.3:8080",
		Username: "jabooth",
	}
	client, err := kubeclient.New(config)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	pods, err := client.Pods(api.NamespaceDefault).List(labels.Everything(), fields.Everything())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	json, err := json.Marshal(pods)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(json)

	return

}

func getEtcdNode(w http.ResponseWriter, r *http.Request) {
	etcd := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	nodePath := r.FormValue("node")
	resp, err := etcd.Get(nodePath, false, false)
	w.Write([]byte(fmt.Sprintf("resp: %+v err: %s", resp, err)))
}
