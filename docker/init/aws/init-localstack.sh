#!/bin/bash

# Set default region (optional, but good practice)
export AWS_DEFAULT_REGION=eu-west-1

# Create Bucket

awslocal s3 mb s3://ecommerce-uploads

# Create SQS queue
awslocal sqs create-queue \
  --queue-name ecommerce-events \
  --region eu-west-1

# Verify queue exists in eu-west-1
echo "Verifying queue creation in eu-west-1..."
QUEUE_URL=$(awslocal sqs get-queue-url \
  --queue-name ecommerce-events \
  --region eu-west-1 \
  --query 'QueueUrl' \
  --output text 2>/dev/null)

if [ $? -eq 0 ]; then
  echo "✅ Queue created successfully: $QUEUE_URL"
else
  echo "❌ Failed to create queue in eu-west-1"
  exit 1
fi

# List all queues (will show across all regions)
echo "All queues across all regions:"
awslocal sqs list-queues --region eu-west-1

# List S3 buckets in eu-west-1
echo "S3 buckets in eu-west-1:"
awslocal s3 ls --region eu-west-1

# Optional: Create the same resources in other regions if needed
# echo "Creating resources in us-east-1 for cross-region testing..."
# awslocal sqs create-queue --queue-name ecommerce-events-us --region us-east-1

echo "✅ LocalStack initialization complete in eu-west-1"