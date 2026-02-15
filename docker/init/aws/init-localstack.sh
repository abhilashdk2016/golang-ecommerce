#!/bin/bash

awslocal s3 mb s3://ecommerce-uploads
awslocal sqs create-queue --queue-name ecomm-events

echo "Localstack initialization complete"