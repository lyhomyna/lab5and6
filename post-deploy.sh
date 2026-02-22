#!/bin/bash
set -e

PROJECT_ID="cloudlab5and6"
IMAGE_NAME="jdm-app"

echo "Running terraform destroy"
cd terraform
terraform destroy -auto-aprove

echo "Removing docker image from GCR"
gcloud container images delete gcr.io/${PROJECT_ID}/${IMAGE_NAME}:latest --quiet

echo "Enabling Google Cloud APIs"
gcloud services disable \
    compute.googleapis.com \
    sqladmin.googleapis.com \
    run.googleapis.com \
    cloudbuild.googleapis.com \
    iam.googleapis.com \
    container.googleapis.com \
    storage.googleapis.com

echo "All resources were deleted"

