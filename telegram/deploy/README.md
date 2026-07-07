## Deploying onto Google Cloud

Prerequisites:
1. Enabled Cloud Run, Cloud Storage, Container Registry in Google Cloud Platform
1. Have a Google Cloud Platform service account key file
1. Installed Terraform CLI SDK
1. Installed gcloud CLI SDK

### Activate Google Cloud Platform service account with key file and configure Docker to push images to Google Container Registry:

```
$ gcloud auth activate-service-account --key-file <key_file_path>
$ gcloud auth configure-docker
```

### Build, tag and push the latest Docker image to Google Container Registry:

```
$ docker build -t personal-finance .
$ docker tag personal-finance <gcr>/<project>/personal-finance
$ docker push <gcr>/<project>/personal-finance
```

Note: To check list of images or image tags/info in Google Container Registry:

```
$ gcloud container images list --repository <gcr>/<project>
$ gcloud container images list-tags <gcr>/<project>/personal-finance
```

### Create a Google Cloud Storage to store Terraform state files to enable Terraform state-locking:

```
$ gsutil mb gs://personal-finance-tf-state-prod
```

Note: To check contents of Google Cloud Storage bucket:

```
$ gsutil cp -r gs://personal-finance-tf-state-prod/ .
```

### Run Terraform configuration for Google Cloud Platform:

```
$ terraform init -backend-config="credentials=<key_file_path>"
$ terraform plan
$ terraform apply
```

Note: To check if new image has been deployed successfully on Cloud Run:

```
$ gcloud run revisions list
```

Note: Use `gcloud config list` to check gcloud configuration and `gcloud config set project <project-id>` to change projects if necessary

### Set the Telegram webhook manually using url from Terraform output:

```
$ curl https://api.telegram.org/bot<token>/setWebhook?url=<terraform_output_url>
```
