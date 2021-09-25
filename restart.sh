pushd k8s
kubectl delete -f=deploy.yml,service.yml
kubectl create -f .
popd
