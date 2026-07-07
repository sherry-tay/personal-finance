terraform {
  backend "gcs" {
    bucket = "personal-finance-tf-state-prod"
    prefix = "terraform/state"
  }
}