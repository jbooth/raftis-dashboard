package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/kubernetes/pkg/api"
	kubeclient "k8s.io/kubernetes/pkg/client"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"net/http"
	"strconv"
)

var ()

func main() {
	//client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		body := `
			<H1>Hello world!!!</H1>
			<a href="listServices">list services</a>
			<a href="listPods">list pods</a>
			<a href="listRCs">list rcs</a>
			</br>
			<a href="launchRaftis">launch raftis (9 hosts, cluster0)</a>
		`
		w.Write([]byte(body))
	})
	http.HandleFunc("/listPods", listPods)
	http.HandleFunc("/listServices", listServices)
	http.HandleFunc("/listRCs", listRCs)
	http.HandleFunc("/getEtcdNode", getEtcdNode)
	http.HandleFunc("/launchRaftis", launchRaftis)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Kubeclient() (* kubeclient.Client, error) {
	config := &kubeclient.Config{
		//Host:     "http://10.65.224.102:8080",
		Host:     "http://172.20.2.3:8080",
		Username: "jabooth",
	}
	return kubeclient.New(config)
}

// There probably is something like that, can't find it
func paramWithDefault(r *http.Request, name string, defValue string) string {
	param := r.FormValue(name)
	if param == "" {
		param = defValue
	}
	return param
}

func doList(w http.ResponseWriter, r *http.Request, listf func(* kubeclient.Client, string) (interface{}, error)) {
	ns := paramWithDefault(r, "ns", "raftis")
	client, err := Kubeclient()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	list, err := listf(client, ns)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	json, err := json.Marshal(list)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(json)

	return
}

func listServices(w http.ResponseWriter, r *http.Request) {
	svs := func(client *kubeclient.Client, ns string) (interface{}, error) {
		return client.Services(ns).List(labels.Everything())
	}
	doList(w, r, svs)
	return
}

func listPods(w http.ResponseWriter, r *http.Request) {
	pods := func(client *kubeclient.Client, ns string) (interface{}, error) {
		return client.Pods(ns).List(labels.Everything(), fields.Everything())
	}
	doList(w, r, pods)
}

func listRCs(w http.ResponseWriter, r *http.Request) {
	rcs := func(client *kubeclient.Client, ns string) (interface{}, error) {
		return client.ReplicationControllers(ns).List(labels.Everything())
	}
	doList(w, r, rcs)
}

func launchRaftis(w http.ResponseWriter, r *http.Request) {
	ns := "eplatonov"
	base := paramWithDefault(r, "base", "cluster0")
	replicasStr := paramWithDefault(r, "replicas", "9")
	replicas, err := strconv.Atoi(replicasStr)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	mountPath := paramWithDefault(r, "mountPath", "/var/raftis/" + base)
	client, err := Kubeclient()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	etcdBase := "/raftis/" + base
	specName := "raftis-" + base
	requestController := &api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name: specName,
			Labels: map[string]string{
				"name": specName,
			},
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: replicas,
			Selector: map[string]string{
				"name": specName,
			},
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{
						"name": specName,
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container {
						api.Container {
							Name: "raftis",
							Image: "raftis/raftis:latest",
							Env: []api.EnvVar {
								api.EnvVar{
									Name: "ETCDURL",
									Value: "http://raftis-dashboard:4001",
								},
								api.EnvVar{
									Name: "ETCDBASE",
									Value: etcdBase,
								},
								api.EnvVar{
									Name: "NUMHOSTS",
									Value: replicasStr,
								},
							},
							Ports: []api.ContainerPort {
								api.ContainerPort{
									ContainerPort: 1103,
								},
								api.ContainerPort{
									ContainerPort: 6379,
								},
							},
							VolumeMounts: []api.VolumeMount {
								api.VolumeMount{
									MountPath: mountPath,
									Name: "data",
								},
							},
						},
					},
					Volumes: []api.Volume {
						api.Volume{
							Name: "data",
						},
					},
				},
			},
		},
	}
	_, err = client.ReplicationControllers(ns).Create(requestController)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	raftisSvc := &api.Service{
		ObjectMeta: api.ObjectMeta{
			Name: specName,
			Labels: map[string]string{
				"name": specName,
			},
		},
		Spec: api.ServiceSpec{
			Type: api.ServiceTypeLoadBalancer,
			Ports: []api.ServicePort{
				api.ServicePort{
					Name: "raftis",
					Port: 6379,
					TargetPort: util.NewIntOrStringFromInt(6379),
				},
			},
			Selector: map[string]string{
				"name": specName,
			},
		},
	}

	svc, err := client.Services(ns).Create(raftisSvc)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(svc)
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
	recStr := paramWithDefault(r, "recursive", "false")
	rec, err := strconv.ParseBool(recStr)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := etcd.Get(nodePath, false, rec)
	w.Write([]byte(fmt.Sprintf("resp: %+v err: %s", resp, err)))
}
