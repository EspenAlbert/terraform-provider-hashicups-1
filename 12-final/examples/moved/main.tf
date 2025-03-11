terraform {
  required_providers {
    hashicups = {
      #   source = "hashicorp.com/edu/hashicups"
    }
  }
  required_version = ">1.10"
}

provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

# create this legacy resource first
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


output "edu_order_legacy" {
  value = hashicups_order.edu
}
