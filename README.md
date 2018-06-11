# kubectl plugin wait
Example kubectl plugin that waits for a kubernetes resource to be ready.

```
kubectl plugin wait po/mypod
```

**Supported Resource Types**
* pods
* deployments

## Install the plugin
```console
$ curl -sLO https://github.com/carolynvs/kubectl-wait-plugin/releases/download/latest/wait.zip

$ unzip wait.zip -d ~/.kube/plugins/
Archive:  wait.zip
   creating: ~/.kube/plugins/wait/
  inflating: ~/.kube/plugins/wait/plugin.yaml
  inflating: ~/.kube/plugins/wait/wait
```

## Build the plugin
```console
$ make deploy
mkdir -p ~/.kube/plugins/wait
go build -o ~/.kube/plugins/wait/wait
cp plugin.yaml ~/.kube/plugins/wait/
```
