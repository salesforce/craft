# Namespace.yaml

The `namespace.yaml` file contains information about the namespace of the cluster in which you want to deploy the operator. The namespace.yaml for the Wordpress operator looks like this:


```
apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
  name: craft
```

