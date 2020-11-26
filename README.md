# Provide The Kubernetes + Istio Apis By ZMC

## 功能描述
1. K8s基础资源API接口
2. Istio ServiceMesh控制器(操作vs、ds)

## Test Run ApiServer
```
go run cmd/apiserver/apiserver.go --kubeconfig k8s-mysql
```

## Test Run Controller Manger
```
go run cmd/controller-manager/main.go --kubeconfig k8s-mysql
```

## If Change The Client Typed, Then do this
```
git clone https://github.com/kubernetes/code-generator.git
cd code-generator
go install ./cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen}

cd oasis
bash ./hack/update-codegen.sh
go mod vendor
```