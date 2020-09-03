# Controller.json

Custom Resource Definition (CRD) information like the domain, group, image, repository, etc. are stored in the controller.json file. This skeleton file for controller.json is created when you [create a CRAFT application in quickstart](quick_setup/README.md): 

    "group": "",
    "resource": "",
    "repo": "",
    "domain": "",
    "namespace": "",
    "version": "",
    "operator_image": "",
    "image": "",
    "imagePullSecrets": "",
    "imagePullPolicy": "",
    "cpu_limit": "",
    "memory_limit": "",
    "vault_addr": "",
    "runOnce": "",
    "reconcileFreq": ""

This table explains the controller.json attributes:

| Attribute        | Description                                                                                                                    |   |   |   |
|------------------|--------------------------------------------------------------------------------------------------------------------------------|---|---|---|
| group            | See the Kubernetes API Concepts page for more information.                                                                     |   |   |   |
| resource         |                                                                                                                                |   |   |   |
| namespace        |                                                                                                                                |   |   |   |
| version          |                                                                                                                                |   |   |   |
| repo             | The repo where you want to store the operator template.                                                                        |   |   |   |
| domain           | The domain web address for this project.                                                                                       |   |   |   |
| operator_image   | The docker registry files used to push operator image into docker.                                                             |   |   |   |
| image            | The docker registry files used to push resource image into docker.                                                             |   |   |   |
| imagePullSecrets | Restricted data to be stored in the operator like access, permissions, etc.                                                    |   |   |   |
| imagePullPolicy  | Method of updating images. Default pull policy is IfNotPresent causes Kubelet to skip pulling an image if one already exists.  |   |   |   |
| cpu_limit        | CPU limit allocated to the operator created.                                                                                   |   |   |   |
| memory_limit     | Memory limit allocated to the operator created.                                                                                |   |   |   |
| vault_addr       | Address of the vault.                                                                                                          |   |   |   |
| runOnce          | If set to 0 reconciliation stops. If set to 1, reconciliation runs according to the specified frequency.                       |   |   |   |
| reconcileFreq    | Frequency interval (in minutes) between two reconciliations.                                                                   |   |   |   |


Here’s an example of a controller.json file for the Wordpress operator: 

```
{
 "group": "wordpress",
  "resource": "WordpressAPI",
  "repo": "wordpress",
  "domain": "salesforce.com",
  "namespace": "default",
  "version": "v1",
  "operator_image": "ops0-artifactrepo1-0-prd.data.sfdc.net/cco/wordpress-operator",
  "image": "ops0-artifactrepo1-0-prd.data.sfdc.net/cco/wordpress:latest",
  "imagePullSecrets": "registrycredential",
  "imagePullPolicy": "IfNotPresent",
  "cpu_limit": "500m",
  "memory_limit": "200Mi",
  "vault_addr": "http://10.215.194.253:8200"
}
```


