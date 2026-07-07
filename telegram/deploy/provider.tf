provider "google" {
  credentials = var.key_file
  project     = var.project
  region      = "us-west1"
}