# Custom Resource Abstraction Fabrication Tool

__NOTE__: For users of CRAFT, a detailed documentation can be found [here](). This README is primarily aimed for developers. 

## Contribution
Please refer [Contribution.md](Contribution.md) before pushing the code. If you wish to make a contribution, create a branch, push your code into the branch and create a PR. For more details, check [this article](https://opensource.com/article/19/7/create-pull-request-github). 

## Installing craft and its dependencies
Dependencies for CRAFT are `kustomize` and `kubebuilder`. 

Latest CRAFT release can be found [here](). Download the craft.tar.gz file and run the following command:

```
sudo tar -C /usr/local/ -xzf ~/Downloads/craft.tar.gz
export PATH=$PATH:/usr/local/craft/bin
```

In case the file is downloaded somewhere other than Downloads, replace ~/Downloads/craft.tar.gz by the path where you downloaded the file. 

## Commands of CRAFT
### craft version
Usage : 
```
craft version
```
Displays the information about craft, namely version, revision, build user, build date & time, go version. 

### craft init
Usage :
```
craft init
``` 
Initialises a new project with sample controller.json and resource.json

### craft create
Usage :
```
craft create -c "controller.json" -r "resource.json --podDockerFile "dockerFile" -p
```
Creates operator source code in $GOPATH/src, builds operator.yaml, builds and pushes operator and resource docker images. 

#### craft build
Has 3 sub commands, code, deploy and image. 

#### build code
Usage:
```
craft build code -c "controller.json" -r "resource.json
```
Creates code in $GOPATH/src/operator. 

#### build deploy
Usage:
```
craft build deploy -c "controller.json" -r "resource.json
```
Builds operator.yaml for deployment onto cluster.

#### build image
Usage:
```
craft build image -b -c "controller.json" --podDockerFile "dockerFile"
```
Builds operator and resource docker images. 

#### validate
Usage:
```
craft validate -v "operator.yaml"
```
Validates operator.yaml to see if everything is in shape
