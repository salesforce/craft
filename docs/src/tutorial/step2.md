# Creating the Operator using CRAFT

Now that we have created controller.json, resource.json and the DockerFile, let's create the Operator using CRAFT.

First, let us check whether CRAFT is working in our machine.
```
$ craft version
```
This should display the version and other info regarding CRAFT.

---
***!TIP***

 If this gives an error saying "$GOPATH is not set", then set GOPATH to the location where you've installed Go.

---

Now that we have verified that CRAFT is working properly, creating the Operator with CRAFT is a fairly straight forward process:
```
craft create -c config/controller.json -r config/resource.json \
--podDockerFile resource/DockerFile -p
```
---
***NOTE***

If the execution in the terminal stops at a certain point, do not assume that it has hanged. The command takes a little while to execute, so give it some time.

---
This will create the Operator template in $GOPATH/src, build operator.yaml for deployment, build and push Docker images for operator and resource. We shall see what these are individually in the next section.
