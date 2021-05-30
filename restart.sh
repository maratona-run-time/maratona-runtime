pushd k8s
kubectl delete -f=deploy.yml,service.yml,pod.yml
kubectl create -f .
popd
