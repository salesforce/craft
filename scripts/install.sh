#!/bin/sh

set -eu

OKGREEN='\033[92m'
FAIL='\033[91m'
WARN='\033[93m'
INFO='\033[94m'
ENDC='\033[0m'
function PassPrint() {
    echo "$OKGREEN $1 $ENDC"
}
function FailPrint() {
    echo "$FAIL $1 $ENDC"
}
function WarnPrint() {
    echo "$WARN $1 $ENDC"
}
function InfoPrint() {
    echo "$INFO $1 $ENDC"
}


OS=$(uname -s)
ARCH=$(uname -m)
OS=$(echo $OS | tr '[:upper:]' '[:lower:]')
VERSION="1.13.1"
NEWPATH=""

function installKB() {
    version=2.2.0 # latest stable version
    arch=amd64

    # download the release
    curl -L -O "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${version}/kubebuilder_${version}_${OS}_${arch}.tar.gz"

    # extract the archive
    tar -zxvf kubebuilder_${version}_${OS}_${arch}.tar.gz
    mv kubebuilder_${version}_${OS}_${arch} kubebuilder && sudo mv kubebuilder /usr/local/

    # update your PATH to include /usr/local/kubebuilder/bin
    NEWPATH+=":/usr/local/kubebuilder/bin"
    export PATH=$PATH:/usr/local/kubebuilder/bin
}

function installKustomize() {
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
}

function installGo(){
        arch="amd64"
        curl -L -O https://dl.google.com/go/go$VERSION.$OS-$arch.tar.gz
        sudo tar -C /usr/local -xzf go$VERSION.$OS-$arch.tar.gz
        
        NEWPATH+=":/usr/local/go/bin"
        export PATH=$PATH:/usr/local/go/bin
}
function installDependency(){
    cmd=`command -v curl` || {
        "curl is missing, it is required for downloading dependencies"
        exit 1
    }
    CUR=`pwd`
    cd /tmp
    cmd1=`command -v kubebuilder` || {
        PassPrint "Installing kubebuilder"
        installKB
    }
    cmd2=`command -v go` || {
        PassPrint "Installing go"
        installGo
    }
    cmd3=`command -v kustomize` || {
        PassPrint "Installing kustomize"
        installKustomize
    }
    cmd4=`command -v schema-generate` || {
        if [ -z '$GOPATH' ] ; then 
            WarnPrint "set GOPATH and install schema-generate by: go get -u github.com/a-h/generate/..."
        else 
            PassPrint "Installing schema-generate"
            go get -u github.com/a-h/generate/...
        fi
    }

    cd $CUR

    if [[ -n $cmd1 && -n $cmd2 && cmd3 ]] ; then
        InfoPrint "Dependencies already exist"
    fi
}

function installCraft() {
    VERSION=$1
    PassPrint "Installing craft@$VERSION in /usr/local"
    sudo rm -rf /usr/local/craft
    CUR=`pwd`
    cd /tmp
    curl -L -O https://github.com/salesforce/craft/releases/download/$VERSION/craft.tar.gz
    TYPE=`file craft.tar.gz`
    NEWPATH+=":/usr/local/craft/bin"
    if [[ "$TYPE" != *"gzip compressed data"* ]]; then
        FailPrint "Downloaded craft.tar.gz is not of correct format. Maybe SSO is required."
        echo """
        Try downloading https://github.com/salesforce/craft/releases/download/$VERSION/craft.tar.gz from browser.
        Then follow:
            sudo tar -C /usr/local -xzf craft.tar.gz
            export PATH=\$PATH$NEWPATH
        """
        exit 1
    fi
    tar -xf craft.tar.gz
    sudo mv craft /usr/local
    cd $CUR
}

function install(){
    installDependency
    installCraft $1
    if [[ -n $NEWPATH ]] ; then
        WarnPrint "export PATH=\$PATH$NEWPATH"
    fi
}
case $OS in
  darwin | linux)
    case $ARCH in
      x86_64)
        install ${1:-"v0.1.0-alpha"}
        ;;
      *)
        echo "There is no linkerd $OS support for $arch. Please open an issue with your platform details."
        exit 1
        ;;
    esac
    ;;
  *)
    echo "There is no linkerd support for $OS/$arch. Please open an issue with your platform details."
    exit 1
    ;;
esac


