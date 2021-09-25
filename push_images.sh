docker-compose build
docker tag compiler mruntime/compiler
docker push mruntime/compiler
docker tag executor mruntime/executor
docker push mruntime/executor
docker tag orchestrator mruntime/orchestrator
docker push mruntime/orchestrator
docker tag orm mruntime/orm
docker push mruntime/orm
docker tag verdict mruntime/verdict
docker push mruntime/verdict
