# pubsub_consumer
This Cloud Function receives push messages from the PubSub topic `payment-audit-log`.
It validates the message format, schema, then store valid messages in a GCS bucket for further processing with analytics tools.

## Deploy
```shell
$ gcloud functions deploy PaymentAuditLog \
--entry-point ProcessLog --runtime go113 \
--trigger-topic payment-audit-log --set-env-vars BUCKET=<BUCKET_NAME>
```

## Environment Variables
- BUCKET: The bucket where valid messages must be stored for further processing.
