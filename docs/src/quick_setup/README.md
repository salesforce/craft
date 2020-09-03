# Quick Setup
## Prerequisites

* [Go](https://golang.org/dl/) version v1.13+
* [Docker](https://docs.docker.com/get-docker/) version 17.03+
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.11.3+
* [Kustomize](https://kubernetes-sigs.github.io/kustomize/installation/) version v3.1.0+
* [Kubebuilder](https://book.kubebuilder.io/quick-start.html#installation) version v2.0.0+

## Installation

The latest version of CRAFT can be found [here](https://github.com/salesforce/craft/releases/). Extract it into **/usr/local** directory.

```
sudo tar -C /usr/local -xzf ~/Downloads/craft.tar.gz 
```

Note: Replace **/Downloads/craft.tar.gz** with the directory where you downloaded the latest version of CRAFT. 

Next, add the /usr/local/craft/bin directory to your PATH environment variable.

```
export PATH=$PATH:/usr/local/craft/bin
```

You can also add CRAFT to your PATH environment variable permanently.

```
sudo vim /etc/paths 
```

Add the line **/usr/local/craft/bin** at the end of the file and save the file.

## Create a CRAFT Application

From the command line, **cd** into a directory where you&#39;d like to store your CRAFT application and run this command:

```
craft init
```

This will initiate a CRAFT application in your current directory and create the following skeleton files:

- controller.json: This file holds Custom Resource Definition (CRD) information like group, domain, operator image, and reconciliation frequency.
- resource.json: This file contains the schema information for validating inputs while creating the CRD.

##Next Steps

Follow the Wordpress operator tutorial to understand how to use CRAFT to create and deploy an operator into a cluster. This deep-dive tutorial demonstrates the entire scope and scale of a CRAFT application.