pushd k8s
kubectl delete -f .
kubectl delete pod --all
kubectl create -f .
popd
