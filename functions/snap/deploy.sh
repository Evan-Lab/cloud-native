gcloud run deploy snap-cmd \
  --source . \
  --function SnapCmd \
  --base-image go125 \
  --region europe-west1 \
  --service-account=snap-cmd@serverless-epitech-dev-476110.iam.gserviceaccount.com \
  --no-allow-unauthenticated \
  --set-env-vars='GOOGLE_CLOUD_PROJECT=serverless-epitech-dev-476110,SNAPSHOT_BUCKET=dev-rplace-bucket,FIRESTORE_DB=dev-rplace-database,SECRET_MANAGER_ID=458258130383'
