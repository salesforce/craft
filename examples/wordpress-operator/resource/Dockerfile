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
