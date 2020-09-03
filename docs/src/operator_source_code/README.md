# Where is the Operator Code?

The source code for the operator is stored in the `$GOPATH/src` path.
```
$ cd $GOPATH/src/wordpress
$ ls

```

The folder contains all the files required to run an operator like the configurations, API files, controllers, reconciliation files, main files, etc. All these files contain information about the operator and it's runtime characteristics, such as the CRUD logic, reconciliation frequency, etc. These files can be classified in four sections:

1.  Build infrastructure.

2.  Launch Configuration.

3.  Entry Point

4.  Controllers and Reconciler.

CRAFT creates these files when you create an operator. This saves you a few weeks of effort to write and connect your operator.

## Build Infrastructure

These files are used to build the operator:
-   go.mod : A Go module for the project that lists all the dependencies.
-   Makefile : File makes targets for building and deploying the controller and reconciler.
-   PROJECT : Kubebuilder metadata for scaffolding new components.
-   DockerFile : File with instructions on running the operator. Specifies the docker entrypoint for the operator.


## Launch Configuration


The launch configurations are in the config/ directory. It holds the CustomResourceDefinitions, RBAC configuration, and WebhookConfigurations. Each folder in config/ contains a refactored part of the launch configuration.
-   `config/default` contains a Kustomize base for launching the controller with standard configurations.
-   `config/manager` can be used to launch controllers as pods in the cluster
-   `config/rbac` contains permissions required to run your controllers under their own service account.


## Entry point

The basic entry point for the operator is in the main.go file. This file can be used to:
-   Set up flags for metrics.
-   Initialise all the controller parameters, including the reconciliation frequency and the parameters we received from controller.json and resource.json
-   Instantiate a manager to keep track of all the running controllers and clients to the API server.
-   Run a manager that runs all of the controllers and keeps track of them until it receives a shutdown signal when it stops all of the controllers.

