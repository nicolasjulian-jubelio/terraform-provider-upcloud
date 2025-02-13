---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "upcloud_loadbalancer_static_backend_member Resource - terraform-provider-upcloud"
subcategory: ""
description: |-
  This resource represents load balancer's static backend member
---

# upcloud_loadbalancer_static_backend_member (Resource)

This resource represents load balancer's static backend member

## Example Usage

```terraform
variable "lb_zone" {
  type    = string
  default = "fi-hel2"
}

resource "upcloud_network" "lb_network" {
  name = "lb-test-net"
  zone = var.lb_zone
  ip_network {
    address = "10.0.0.0/24"
    dhcp    = true
    family  = "IPv4"
  }
}

resource "upcloud_loadbalancer" "lb" {
  configured_status = "started"
  name              = "lb-test"
  plan              = "development"
  zone              = var.lb_zone
  network           = resource.upcloud_network.lb_network.id
}

resource "upcloud_loadbalancer_backend" "lb_be_1" {
  loadbalancer = resource.upcloud_loadbalancer.lb.id
  name         = "lb-be-1-test"
}

resource "upcloud_loadbalancer_static_backend_member" "lb_be_1_sm_1" {
  backend      = resource.upcloud_loadbalancer_backend.lb_be_1.id
  name         = "lb-be-1-sm-1-test"
  ip           = "10.0.0.10"
  port         = 8000
  weight       = 0
  max_sessions = 0
  enabled      = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backend` (String) ID of the load balancer backend to which the member is connected.
- `ip` (String) Server IP address in the customer private network.
- `max_sessions` (Number) Maximum number of sessions before queueing.
- `name` (String) The name of the member must be unique within the load balancer backend service.
- `port` (Number) Server port.
- `weight` (Number) Used to adjust the server's weight relative to other servers. 
				All servers will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. 
				A value of 0 means the server will not participate in load balancing but will still accept persistent connections.

### Optional

- `enabled` (Boolean) Indicates if the member is enabled. Disabled members are excluded from load balancing.

### Read-Only

- `id` (String) The ID of this resource.


