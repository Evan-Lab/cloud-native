gcloud run deploy snap-cmd \
  --source . \
  --function SnapCmd \
  --base-image go125 \
  --region europe-west1 \
  --service-account=snap-cmd@serverless-epitech-dev-476110.iam.gserviceaccount.com \
  --no-allow-unauthenticated \
  --set-env-vars='GOOGLE_CLOUD_PROJECT=serverless-epitech-dev-476110,SNAPSHOT_BUCKET=dev-rplace-bucket,FIRESTORE_DB=dev-rplace-database'

# Create OR Update the Trigger
gcloud eventarc triggers create snap-cmd-trigger \
  --location=europe-west1 \
  --destination-run-service=snap-cmd \
  --destination-run-region=europe-west1 \
  --event-filters="type=google.cloud.pubsub.topic.v1.messagePublished" \
  --event-filters="topic=projects/serverless-epitech-dev-476110/topics/command.snap" \
  --service-account=snap-cmd@serverless-epitech-dev-476110.iam.gserviceaccount.com \
  || \
gcloud eventarc triggers update snap-cmd-trigger \
  --location=europe-west1 \
  --destination-run-service=snap-cmd \
  --service-account=snap-cmd@serverless-epitech-dev-476110.iam.gserviceaccount.com