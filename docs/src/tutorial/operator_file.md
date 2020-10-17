The `operator.yaml` file contains all the metadata required to deploy the operator into the cluster. It is the backbone of the operator as it contains the schema validations, the specification properties, the API version rules, etc. This file is automatically populated by CRAFT based on the information provided in the controller.json and the resource.json file. The `operator.yaml` that CRAFT generates for the Wordpress operator can be found at `examples/wordpress-operator/config/deploy/operator.yaml`.

---
***Note***  

 Our operator's default user currently holds the minimum *rbac* required to run the operator and only the operator itself. If you need any more control of the *rbac* add those permissions in the `operator.yaml`.

---
