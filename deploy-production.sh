#!/bin/bash

set -e

docker build -t gcr.io/supermarket-05052018/jstorer/gannett:$TRAVIS_COMMIT .

echo $GCLOUD_SERVICE_KEY_PRD | base64 --decode -i > ${HOME}/gcloud-service-key.json
gcloud auth activate-service-account --key-file ${HOME}/gcloud-service-key.json

gcloud --quiet config set project supermarket
gcloud --quiet config set container/cluster supermarket-05052018
gcloud --quiet config set compute/zone us-east1-b
gcloud --quiet container clusters get-credentials supermarket-cluster

gcloud docker push gcr.io/supermarket-05052018/jstorer/gannett

yes | gcloud beta container images add-tag gcr.io/supermarket-05052018/jstorer/gannett:$TRAVIS_COMMIT gcr.io/supermarket-05052018/jstorer/gannett:latest

kubectl config view
kubectl config current-context

kubectl set image deployment/jstorer/gannett jstorer/gannett=gcr.io/supermarket-05052018/jstorer/gannett:$TRAVIS_COMMIT

#
# sleep 30
# npm run e2e_test