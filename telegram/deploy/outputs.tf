output "url" {
  value = google_cloud_run_service.personal_finance.status[0].url
}