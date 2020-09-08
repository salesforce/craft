# Resource DockerFile - How does it help?

For our Wordpress operator, the resource Dockerfile looks like this:

```
FROM centos/python-36-centos7:latest

USER root

RUN pip install --upgrade pip
RUN python3 -m pip install pyyaml
RUN python3 -m pip install jinja2
RUN python3 -m pip install hvac

ADD kubectl kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl

ARG vault_token=dummy
ENV VAULT_TOKEN=${vault_token}

ADD templates templates

ADD initwordpress.sh .
RUN chmod 755 initwordpress.sh
ADD wordpress_manager.py .
RUN chmod 755 wordpress_manager.py

RUN find / -perm /6000 -type f -exec chmod a-s {} \; || true

ENTRYPOINT ["python3", "wordpress_manager.py"]
```
The above resource Dockerfile for the Wordpress operator has Docker run two files/scripts: 
 
1. `initwordpress.sh` : This file contains instructions to initialize the wordpress resource and install the required components. 
2. `wordpress_manager.py` : This file contains CRUD operations defined by the Wordpress resource. These operations don’t return the usual output, but return exit codes. 

---
**FAQ**

**Q. The methods in wordpress_manager.py take input_data/spec as input parameter. Where is it being passed in?**

**A.** Both the input_data/spec and action_type are passed as args at pod spec. Operator deploys a pod for each execution of action_type (create/update/delete/verify)

**Q. The kubectl binary is being packaged within the Operator's container image. How/where is it being given the credentials to access the API server?**

**A.** We can also use the `client-go` library instead of `kubectl` binary in Wordpress example. `kubectl` also uses the default service account available credentials in the pod itself. Just like client-go kubectl default to that service account and we are not specifically configuring any other credentials. [K8s Docs](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/#accessing-the-api-from-a-pod)

---
