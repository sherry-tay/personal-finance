resource "google_cloud_run_service" "personal_finance" {
  name = "personal-finance"
  location = "us-west1"

  template {
    spec {
      containers {
        image = var.image
        env {
          name  = "TELEGRAM_BOT_TOKEN"
          value = var.telegram_bot_token
        }
        env {
          name  = "AUTHORIZED_USER"
          value = var.telegram_bot_authorised_username
        }
        env {
          name  = "IS_GCP"
          value = true
        }
        env {
          name  = "PROJECT_ID"
          value = var.project
        }
      }
      container_concurrency = 1
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role    = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.personal_finance.location
  project  = google_cloud_run_service.personal_finance.project
  service  = google_cloud_run_service.personal_finance.name

  policy_data = data.google_iam_policy.noauth.policy_data
}