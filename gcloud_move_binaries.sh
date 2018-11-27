#!/bin/bash

FILES=./*

for f in $FILES
do
  echo "Moving $f to remote instance..."
  gcloud compute scp $f samnaser@containerd-image-generator:/home/samnaser/$f --project=dashpole-gke-dev --zone=us-central1-c
done
