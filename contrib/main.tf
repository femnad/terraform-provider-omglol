provider "omglol" {
  username = var.username
  api_key  = var.api_key
}

resource "omglol_dns" "new" {
  name = "bar"
  type = "CNAME"
  data = "bar.example.com"
}
