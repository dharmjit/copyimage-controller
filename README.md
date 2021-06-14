## CopyImage Controller 
To copy public registry container images referenced in Deployment/DaemonSets to another private OCI Registry

> Currently only support Docker(index.docker.io) registry

This controller is developed using [kubebuilder](https://book.kubebuilder.io). As we do not require any CRDs for this controller, API Types creation is skipped.

### Build/Run CopyImage Controller

Local Dev:
- Prerequisite
    Make sure the following are installed.

    - Install Go 1.14+
        - [go](https://golang.org/)

    - Install kubebuilder
        - [kubebuilder](https://book.kubebuilder.io/getting_started/installation_and_setup.html)  

    - Install kubectl
        - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl)

    - Install kind
        - [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
    
    - (Optional) Install Make
        - [Make](https://www.gnu.org/software/make/)
- Clone this project (which has already been setup for you)
  - requires [git](https://git-scm.com/downloads)
  - `git clone https://github.com/dharmjit/copyimage-controller`
- Setup Local Dev cluster
    - kind [kind](https://kind.sigs.k8s.io/)
        - Requires [go(1.14+)](https://golang.org/doc/devel/release#policy) and [docker](https://www.docker.com/)
        - Create the cluster `kind create cluster`
        - Test the setup `kind get cluster` and `kubectl get nodes`
- Run the controller
    - Change directory in another shell `cd copyimage-controller`
    - Make run with your credentials `PRIV_OCI_REGISTRY=<<>> PRIV_OCI_REPOSITORY=<<>> PRIV_OCI_REGISTRY_USERNAME=<<>> PRIV_OCI_REGISTRY_PASSWORD=<<>> make run`
- Verify the controller
    - `cd copyimage-controller/test`
    - `kubectl create -f nginx-singlecontainer.yaml`
    - `kubectl get deploy/nginx-single-deployment -o yaml`

Remote Deployment:

- 

## Supplementary Resources

For those who are interested, there is documentation on kubebuilder available here:

- [http://book.kubebuilder.io](http://book.kubebuilder.io)