# CRAFTing Wordpress

### Creating Operator
#### Configure Private Registry
Update `config/controller.json` with your environment variables like docker registry, docker secret name, etc.
```
    "imagePullSecrets": "docker-registry-credentials",
```
#### Generate Wordpress operator
```
craft create -c config/controller.json -r config/resource.json --podDockerFile resource/DockerFile -p
```

### Deploy operator to cluster
#### Install operator
```
kubectl apply -f config/deploy/operator.yaml
```
#### Create docker private registry secret
```
kubectl create secret docker-registry docker-registry-credentials --docker-server=<your-private-docker-registry> \
--docker-username=<some_user> --docker-password=<some_password> --namespace=craft
```

### Deploy wordpress resource
```
kubectl apply -f config/deploy/operator.yaml
```
#### Verfication of the wordpress operator

For testing we are going to access the wordpress service using kubectl proxy
```
kubectl -n craft port-forward  svc/wordpress-dev 9090:80
```
To access the wordpress application, open `http://localhost:9090`
 

### Cleanup: Delete resource, operator and namespace.
```
kubectl delete -f conf/deploy/operator.yaml

```
