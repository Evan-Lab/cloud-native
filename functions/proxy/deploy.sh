gcloud run deploy discord-proxy \
  --source . \
  --function DiscordProxy \
  --base-image go125 \
  --region europe-west1 \
  --service-account=discord-hello@serverless-epitech-dev-476110.iam.gserviceaccount.com \
  --set-env-vars='SECRET_MANAGER_ID=458258130383,GOOGLE_CLOUD_PROJECT=serverless-epitech-dev-476110'