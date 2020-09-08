# Deploy operator onto the cluster

In the previous step, we created an operator. Now, we deploy it to the cluster. This involves deploying the namespace and the operator files to the cluster. Create the namespace where you want to deploy the operator with the namespace.yaml file we created in step 2 using this command:

```
kubectl apply -f config/deploy/namespace.yaml
```
When the command runs successfully, it returns `namespace/craft created`. You can check the namespace created by running this command:

```
kubectl get namespace
```
This should display all the existing namespaces, out of which `craft` is one. Install the operator onto the cluster with this command:

```
kubectl apply -f config/deploy/operator.yaml
```

This will create the required pod in the cluster. We can verify the creation by running:

```
kubectl get pods
```

This returns the wordpress pod along with the other pods running on your machine:
```
NAME                                           READY   STATUS         RESTARTS   AGE
wordpress-controller-manager-8844cf545-gn5rt   1/2     Running           0       11s
```

Great, your pod is running! You are ready to deploy the resource.

---
***NOTE***

If your pod’s status is `ContainerCreating`, run the command again in a few seconds and the status should change to running.

---

Deploy the resource onto the cluster using the wordpress-dev-withoutvault YAML file created by CRAFT.

```
kubectl -n craft apply -f config/deploy/wordpress-dev-withoutvault.yaml
```

This deploys the wordpress resource onto the cluster. To verify, run:
```
kubectl -n craft port-forward  svc/wordpress-dev 9090:80
```

Open `http://localhost:9090` on the browser and you’ll see the Wordpress application. You can check the logs to see that reconciliation is running as configured. To see that, we can use [stern](https://github.com/wercker/stern).

```
stern -n craft .
```
