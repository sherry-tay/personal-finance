variable "key_file" {
  type        = string
  description = "Path to service account private key file"
}

variable "project" {
  type        = string
  description = "GCP project ID"
}

variable "image" {
  type        = string
  description = "Docker image to deploy"
}

variable "telegram_bot_token" {
  type        = string
  description = "Telegram bot token to connect to Bot API"
}

variable "telegram_bot_authorised_username" {
  type        = string
  description = "Telegram username allowed to use this Telegram bot"
}