# Install Docker

https://docs.docker.com/desktop/setup/install/mac-install/

# Install Kind

```sh
brew install kind
```

https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager

Kind uses a single docker container to simulate a node. Inside the docker container it runs systemd, and uses systemd to manage kubelet and containerd. Then, through the kubelet inside the container, it starts up the other K8s components such as kube-apiserver, etcd, etc. Finally, once the CNI is deployed on top, the whole cluster is complete.
Internally, Kind also uses kubeadm for deployment.

# Pull Golang Image From docker

```sh
docker pull golang:latest
docker image ls
```

# Build and run the blog image

```sh
docker build -t blog:v1.0.0 ./blog/
docker image ls
docker run -p 80:5678 blog:v1.0.0
curl http://localhost/?name=fulim
```

# Kubernetes Concept

In Kubernetes cluster

1. 1 of the node is control plan to receive request from client and manage compute machines
2. Other nodes is compute machines where it can deploy multiple pods
3. Every Pod can deploy multiple container
   ![alt text](<CleanShot 2026-09-15 at 11.32.23 PM.png>)

# Create a k8s cluster

Create a k8s cluster with 1 control-plane node and 3 worker nodes. In practice, this creates <bold>and starts</bold> 4 containers in docker, to simulate 4 physical nodes. If the name is already specified in the yaml file, there's no need to specify it again in the create command.

```sh
kind create cluster --name dqq --config deploy/cluster-dqq.yaml
```

View clusters

```sh
kind get clusters
```

View the nodes in a cluster. Each node corresponds to a container in Docker.

```sh
kind get nodes --name dqq
```

Export kind logs for troubleshooting

```sh
kind export logs --name dqq ./kind_log
```

[Delete a cluster]

```sh
kind delete cluster --name dqq
```

# Build images

Manually pull the golang image

```sh
docker pull golang
```

Write the Dockerfile.
Build the image. -t specifies the tag

```Shell
docker build -t blog:v1.0.0 ./blog
docker build -t lottery:v1.0.0 ./lottery
docker build -t search:v1.0.0 ./search
```

View images. You should see golang:latest, blog:v1.0.0, lottery:v1.0.0 and search:v1.0.0

```sh
docker image ls
```

Run the image

```sh
docker run -p 80:5678 blog:v1.0.0
```

Load the image into the k8s cluster

```Shell
// kind load docker-image my-custom-image --name my-cluster-name
kind load docker-image blog:v1.0.0  --name dqq
kind load docker-image lottery:v1.0.0  --name dqq
kind load docker-image search:v1.0.0  --name dqq
```

Check which images exist on a node

```sh
docker exec -it dqq-worker crictl images
```

[Remove images from docker]

```Shell
docker image rm -f blog:v1.0.0
docker image rm -f lottery:v1.0.0
docker image rm -f search:v1.0.0
```

# Container Orchestration

Check the current nodes and get their names

```sh
kubectl get node
```

Label the nodes

```Shell
kubectl label nodes dqq-worker hp=true ls=true
kubectl label nodes dqq-worker ingress-ready=true
kubectl label nodes dqq-worker2 ls=true
kubectl label nodes dqq-worker2 ingress-ready=true
kubectl label nodes dqq-worker3 hp=true ls=true
```

View nodes with a given label

```Shell
kubectl get node -l ingress-nginx
kubectl get node -l ls=true
```

View all labels on a node

```sh
kubectl label node dqq-worker --list
```

Create a Deployment to deploy Pods

```Shell
kubectl apply -f deploy/dep-blog.yaml
kubectl apply -f deploy/dep-lottery.yaml
kubectl apply -f deploy/dep-search.yaml
```

View pods

```sh
kubectl get pod -o wide
```

View pod details

```sh
kubectl describe pod ${pod_name}
```

Delete a pod. After deleting a pod, it will be immediately restarted on another node

```sh
kubectl delete pod dep-blog-fbf9ccd46-r8c59
```

View Deployments

```sh
kubectl get deployment -o wide
```

[To delete a pod, you must delete the corresponding deployment]

```Shell
kubectl delete deployment dep-blog
kubectl delete deployment dep-lottery
kubectl delete deployment dep-search
```

Create a service

```Shell
kubectl apply -f deploy/svc-blog.yaml
kubectl apply -f deploy/svc-lottery.yaml
kubectl apply -f deploy/svc-search.yaml
```

View services

```sh
kubectl get svc
```

Install Ingress Controller:
Download the file

```sh
wget -O deploy/ingress-controller.yaml https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

Label the node

```sh
kubectl label nodes dqq-control-plane ingress-ready=true
```

Install the Ingress Nginx Controller

```sh
kubectl apply -f deploy/ingress-controller.yaml
```

Check whether the corresponding Pod has been created. -n specifies the namespace; any kubectl operation on a resource needs to specify the namespace, unless it's the default namespace: default

```sh
kubectl get pod -n ingress-nginx
```

View pods across all namespaces

```sh
kubectl get pods --all-namespaces
```

Delete an entire namespace and all resources under it

```sh
kubectl delete namespace ingress-nginx
```

View pod details

```sh
kubectl describe pod ${pod_name} -n ingress-nginx
```

Create Ingress

```sh
kubectl apply -f deploy/ingress-nginx.yaml
```

# HPA, Pod Horizontal Autoscaling

Install metrics-server

```sh
wget -O deploy/metrics-server.yaml https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Modify the downloaded file, adding `- --kubelet-insecure-tls` under containers' args

```sh
kubectl apply -f deploy/metrics-server.yaml
```

Check whether metrics-server installed successfully

```Shell
kubectl get pod -o wide -n kube-system | grep metrics
kubectl describe pod metrics-server-587b667b55-whqcw -n kube-system
```

Edit the pod spec file to limit CPU usage. Edit deploy/dep-blog.yaml:

```yaml
containers:   # which containers the pod contains
    - name: blog   # container name
        image: blog:v1.0.0  # which image the container runs
        imagePullPolicy: IfNotPresent   # image pull policy10m
        resources:
        requests:    # controls the pod's hardware resource usage
            cpu: 10m  # 1 CPU core is 1000m, 10m is equivalent to 0.01 CPU
```

Create HPA

```sh
kubectl apply -f deploy/hpa-blog.yaml
```

View HPA

```sh
kubectl get hpa
```

Under TARGETS, the CPU usage will initially show unknown. Run the command again after a while and you'll see CPU usage at 0%. If it's still unknown, check the logs:

```sh
kubectl describe hpa hpa-blog
```

Delete HPA

```sh
kubectl delete hpa hpa-blog
```

Run the stress test

```sh
bash deploy/stress_test.sh
```

After a while, run `kubectl get hpa` again and you'll see CPU usage has spiked. Check the number of blog pods, which should have reached 5:

```sh
kubectl get pod | grep blog
```

Stop the stress test, then run `kubectl get hpa` and `kubectl get pod | grep blog` again to check.
