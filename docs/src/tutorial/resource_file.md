# Resource.json

Resource.json contains the schema information for validating inputs while creating the Custom Resource Definition (CRD). The resource.json has a list of properties and required attributes for a certain operator.

    "type": "object"
    "properties": {},
    "required": [],

The properties field contains the field name, data type and the data patterns. The required field contains the data names which are mandatory for the operator to be created.
 
Note: Resource.json remains the same for an operator regardless of the developer. controller.json for an operator may look different for each developer. 
 
Populate resource.json by identifying required fields and their properties from the resource code. Here’s an example resource.json file for the Wordpress operator:

```
  "type": "object",
  "properties": {
    "bootstrap_email": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "bootstrap_password": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "bootstrap_title": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "bootstrap_url": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "bootstrap_user": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "db_password": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "dbVolumeMount": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "host": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "instance": {
      "enum": [
          "prod",
          "dev"
      ],
      "type": "string"
    },
    "name": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "replicas": {
      "format": "int64",
      "type": "integer",
      "minimum": 1,
      "maximum": 5
    },
    "user": {
      "pattern": "^(.*)$",
      "type": "string"
    },
    "wordpressVolumeMount": {
      "pattern": "^(.*)$",
      "type": "string"
    }
  },
  "required": [
    "bootstrap_email",
    "bootstrap_password",
    "bootstrap_title",
    "bootstrap_url",
    "bootstrap_user",
    "db_password",
    "dbVolumeMount",
    "host",
    "instance",
    "name",
    "replicas",
    "user",
    "wordpressVolumeMount"
  ]
```
