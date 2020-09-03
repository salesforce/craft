# Custom Resource Abstraction Fabrication Tool

For new users of CRAFT, a detailed documentation can be found [here](https://salesforce.github.io/craft). 

## Contribution
Please refer [Contribution.md](Contribution.md) before pushing the code. If you wish to make a contribution, create a branch, push your code into the branch and create a PR. For more details, check [this article](https://opensource.com/article/19/7/create-pull-request-github). 

## Installing craft and its dependencies
Dependencies for CRAFT are `kustomize` and `kubebuilder`. 


```
# dowload latest craft binary from releases and extract 
curl -L https://github.com/salesforce/craft/releases/download/0.1.7/craft.tar.gz | tar -xz -C /tmp/

# move to a path that you can use for long term
sudo mv /tmp/craft /usr/local/craft
export PATH=$PATH:/usr/local/craft/bin
```
## CRAFT Usage
To know more about how to use craft cli you can refer to [here](https://salesforce.github.io/craft/craft_cli.html)
