kubectl port-forward deployment/orm 8084:8084 &
kubectl port-forward deployment/postgres 5432:5432 &
kubectl port-forward deployment/orchestrator 8080:8080 &
kubectl port-forward pods/mart 8083:8083 8082:8082 8081:8081 &
kubectl port-forward deployment/rabbitmq 5672:5672 15672:15672 &
