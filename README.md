## CopyImage Controller 
To copy public registry container images referenced in Deployment/DaemonSets to another private OCI Registry

> Currently only support Docker(index.docker.io) registry

## Demo
[![asciicast](https://asciinema.org/a/rX2K9kK2p5vsli5vGPIYc7oEW.png)](https://asciinema.org/a/rX2K9kK2p5vsli5vGPIYc7oEW)

## Development
This controller is developed using [kubebuilder](https://book.kubebuilder.io) to generate all the bootstrap code as well as kuberetes manifests. As we do not require any CRDs for this controller, API Types creation is skipped.

Below controllers are created for watching/reconciling `deployments` and `daemonsets` respectively. 

- CopyImageDeploymentReconciler
- CopyImageDaemonsetReconciler

Both controllers filters only create/update events on their respective watch objects as shown for CopyImageDeploymentReconciler

```go
// CopyImageController for Deployments
func (r *CopyImageDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Deployment{}).
		WithEventFilter(predicate.Funcs{
			DeleteFunc: func(e event.DeleteEvent) bool {
				// Suppress Delete events to avoid filtering them out in the Reconcile function
				return false
			},
		}).
		Complete(r)
}
```

There a utility, `utils.CloneImage` which intialize the private registry credentials along with the main logic to check/copy images to private registry.
### Local Setup:

- Prerequisite
    Make sure the following are installed.
    - Install Git
        - [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
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

- Update code with your settings
  - Update Makefile with your Docker Repository
    - IMG ?= {your_repo_name}/copyimage-controller:latest
  - Update Secret `dockersecret` defined in `config/manager/manager.yaml` with base64 values applicable for you

- Setup Local Dev cluster
    - kind [kind](https://kind.sigs.k8s.io/)
        - Requires [go(1.14+)](https://golang.org/doc/devel/release#policy) and [docker](https://www.docker.com/)
        - Install [kind]    (https://kind.sigs.k8s.io/docs/user/quick-start/)
        - Create the cluster `kind create cluster`
        - Test the setup `kind get cluster` and `kubectl get nodes`
- Run the controller
    - Change directory in another shell `cd copyimage-controller`
    - Make run with your credentials `PRIV_OCI_REGISTRY=index.docker.io PRIV_OCI_REPOSITORY=your_repo_name PRIV_OCI_REGISTRY_USERNAME=your_username PRIV_OCI_REGISTRY_PASSWORD=your password make run`

- Verify the controller
    - `cd copyimage-controller/test`
    - `kubectl create -f nginx-singlecontainer.yaml`
    - `kubectl get deploy/nginx-single-deployment -o yaml`

- Deploy the controller as a deployment
    - set below variables in your terminal
        ```
        export PRIV_OCI_REGISTRY="index.docker.io"
        export PRIV_OCI_REPOSITORY=""
        export PRIV_OCI_REGISTRY_USERNAME=""
        export PRIV_OCI_REGISTRY_PASSWORD=""```
    - Build Controller Docker Image `make docker-build`
    - Push Controller Docker Image `make docker-push`
    - Deploy Controller and other Manifests `make deploy`

## TODO

- [ ] Do not Remote Write Image if it already exists
- [ ] Initialize an OCI client to reduce repeated authention with `remote.Write`
- [ ] Update Message for `kubectl rollout history`
- [ ] Concurrency Handling of `Utils.CloneImage` Method
- [ ] Handle Optimistic Concurrency effectively if possible
- [ ] Better Handle `utils.CloneImage` errors
- [ ] Externalize the Hardcoded Public Registry Name. It could be a list of public  registries
## Supplementary Resources

For those who are interested, there is documentation on kubebuilder available here:

- [http://book.kubebuilder.io](http://book.kubebuilder.io)