# Install Docker

https://docs.docker.com/desktop/setup/install/mac-install/

# Install Kind

https://kind.sigs.k8s.io/docs/user/quick-start/#installing-with-a-package-manager

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
