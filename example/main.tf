resource "google_compute_backend_bucket" "image_backend" {
  name        = "${local.name_prefix}-backend-bucket"
  bucket_name = "${google_storage_bucket.image_bucket.name}"
  enable_cdn  = true
}

resource "google_storage_bucket" "image_bucket" {
  name     = "${local.name_prefix}-bucket"
  location = "EU"
}

locals {
  name_prefix = "${random_pet.prefix.id}"
}

resource "random_pet" "prefix" {
  length = 2
}
