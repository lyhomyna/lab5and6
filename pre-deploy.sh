#!/bin/bash
set -e

PROJECT_ID="cloudlab5and6"
IMAGE_NAME="jdm-app"

echo "Enabling Google Cloud APIs"
gcloud services enable \
    compute.googleapis.com \
    sqladmin.googleapis.com \
    run.googleapis.com \
    cloudbuild.googleapis.com \
    iam.googleapis.com

echo "Configuring Docker Auth"
gcloud auth configure-docker --quiet

echo "Authenticating for Terraform"
gcloud auth application-default login --no-launch-browser

echo "Building and Pushing Image"
gcloud builds submit --tag gcr.io/$PROJECT_ID/$IMAGE_NAME:latest .

echo "All systems ready for Terraform"

