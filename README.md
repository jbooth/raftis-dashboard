edit main.go to do different stuff on port 8080

also runs etcd on port 4001

./publish.sh to push new image under raftis/dashboard

to deploy to kubernetes, first see http://go/tess to set up your kubectl client, it's pretty quick

cd kubernetes
kubectl create -f raftis-dashboard.yaml # creates pod and replicationController of size 1
kubectl get pods # check pods
kubectl get rc  # check replicationController
kubectl create -f raftis-dashboard-service.yaml # creates service exposing ports 8080 and 4001 outside of kubernetes network
kubectl get svc # check services


after a couple minutes, a second IP will appear under the service definition for raftis-dashboard under `kubectl get svc` -- that's the public IP, use that IP on port 4001 or port 8080 and it should be accessible 
