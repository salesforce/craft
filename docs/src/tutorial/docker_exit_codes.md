# CRUD operations in CRAFT operators

Define CRUD (Create, Read, Update, and Delete) operations for your operator. This diagram illustrates the flow for the 14 possible outputs of a CRUD operation:

![operator-full-lifecycle](https://www.stephenzoio.com/images/operator-full-lifecycle.png)
*Credits: Stephen Zoio & his project operatify*

With CRAFT, you can account for all the14 cases by using 14 unique docker exit codes. Use these docker exit codes to check the result of the operation and route the path accordingly. The 14 exit codes that CRAFT provides are:

```
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

The 14 exit codes are correspondingly mapped to the 14 possibilities that arise from the operations.

While writing our CRUD operation definitions, we use the exit codes to specify output. For example, in the `wordpress_manager.py`, we can see the CUD operations:

```
def create_wordpress(spec):
    /*
    ..
    */
    if result.returncode == 0:
        init_wordpress(spec)
        sys.exit(201)
    else:
        sys.exit(203)

def delete_wordpress(spec):
    /*
    ..
    */
    if result.returncode == 0:
        sys.exit(221)
    else:
        sys.exit(223)

def update_wordpress(spec):
    /*
    ..
    */
    if result.returncode == 0:
        sys.exit(201)
    else:
        sys.exit(203)

def verify_wordpress(spec):
    /*
    ..
    */
    if result.returncode == 0:
        result = subprocess.run(['kubectl', 'get', 'deployment', 'wordpress-' + spec['instance'],  '-o', 'yaml'], stdout=subprocess.PIPE)
        deployment_out = yaml.safe_load(result.stdout)
        if deployment_out['spec']['replicas'] != spec['replicas']:
            print("Change in replicas.")
            sys.exit(214)
        sys.exit(211)
    else:
        sys.exit(214)
```

We map the corresponding output possibility to the exit code.

Now that we have defined our CRUD operations, let's create our operator.
