# Linescrape
![CodeQL](https://github.com/fischersean/linescrape/workflows/CodeQL/badge.svg)

AWS Lambda functions and related utilities for collecting betting line and win projection data

## Deloying
Install AWS SAM and make sure it is functioning correctly with you AWS account. 

Next, create an S3 bucket and place the API config (LineScrapeAPIv1.yaml) in the bucket. Replace the 'Bucket' key in the template.yaml with your bucket name.

Once SAM is installed and the S3 bucket is configured:

`make build && make deploy`

NOTE: You will have to manually change the route in the API Gateway configuration to point to your own URL.
