# Controllers and Reconciler

A controller is the core of Kubernetes and operators. A controller ensures that, for any given object, the actual state of the world (both the cluster state, and potentially external state like running containers for Kubelet or loadbalancers for a cloud provider) matches the desired state in the object. Each controller focuses on one root API type but may interact with other API types.

A reconciler tracks changes to the root API type, checks and updates changes in the operator image at controller-runtime. It runs an operation and returns an exit code, and through this process it checks if reconciliation is needed and determines the frequency for reconciliation. Based on the exit code, the next operation is added to the controller queue.

These are the 14 exit codes that a reconciler can return:

```
## ExitCode to state mapping
  201: "Succeeded", // create or update
  202: "AwaitingVerification", // create or update
  203: "Error", // create or update
  211: "Ready", // verify
  212: "InProgress", // verify
  213: "Error", // verify
  214: "Missing", // verify
  215: "UpdateRequired", // verify
  216: "RecreateRequired", // verify
  217: "Deleting", // verify
  221: "Succeeded", // delete
  222: "InProgress", // delete
  223: "Error", // delete
  224: "Missing", // delete
```

