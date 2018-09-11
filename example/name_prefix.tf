locals {
  name_prefix = "${random_pet.prefix.id}"
}

resource "random_pet" "prefix" {
  length = 2
}
