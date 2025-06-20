variable "client_id" {
  description = "Client ID of the connected app"
  type        = string
}

variable "private_key" {
  description = "Private Key associated to the public certificate"
  type        = string
  sensitive   = true
}

variable "username" {
  description = "Salesforce Username of a System Administrator"
  type        = string
}

variable "api_version" {
  description = "API version of the salesforce org"
  type        = string
  default     = "58.0"
}

variable "login_url" {
  description = "Salesforce login URL"
  type        = string
  default     = "https://login.salesforce.com"
} 