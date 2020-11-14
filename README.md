# pubsub_consumer
Payments entries are posted in a topic called `payment-audit-log`.
This Cloud Function with push subscription validates the messages, identity poison and invalid messages, and store them in GCS buckets for further processing with analytics tools.

## Deploy
```shell
$ gcloud functions deploy PaymentAuditLog \
--entry-point ProcessLog --runtime go113 \
--trigger-topic payment-audit-log --set-env-vars BUCKET=<BUCKET_NAME>
```

## Environment Variables
- BUCKET: The bucket where valid messages must be stored for further processing.
