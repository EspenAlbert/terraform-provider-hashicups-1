terraform {
  required_providers {
    hashicups = {}
  }
  required_version = ">1.10"
}

provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

# LEGACY RESOURCE
# resource "hashicups_order_legacy" "edu" {
#   items = [{
#     coffee = {
#       id = 3
#     }
#     quantity = 2
#     },
#     {
#       coffee = {
#         id = 2
#       }
#       quantity = 3
#   }]
# }

# NEW RESOURCE
moved {
  from = hashicups_order_legacy.edu
  to   = hashicups_order.edu
}

resource "hashicups_order" "edu" {
  items = [{
    coffee = {
      id = 3
    }
    quantity = 2
    },
    {
      coffee = {
        id = 2
      }
      quantity = 3
  }]
}


# LEGACY OUTPUT VAR
# output "edu_order_legacy" {
#   value = hashicups_order_legacy.edu
# }

# NEW OUTPUT VAR
output "edu_order" {
  value = hashicups_order.edu
}
