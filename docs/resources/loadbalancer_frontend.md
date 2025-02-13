---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "upcloud_loadbalancer_frontend Resource - terraform-provider-upcloud"
subcategory: ""
description: |-
  This resource represents load balancer frontend service
---

# upcloud_loadbalancer_frontend (Resource)

This resource represents load balancer frontend service

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

resource "upcloud_loadbalancer_frontend" "lb_fe_1" {
  loadbalancer         = resource.upcloud_loadbalancer.lb.id
  name                 = "lb-fe-1-test"
  mode                 = "http"
  port                 = 8080
  default_backend_name = resource.upcloud_loadbalancer_backend.lb_be_1.name
  networks {
    name = resource.upcloud_loadbalancer.lb.networks[1].name
  }
}

resource "upcloud_loadbalancer" "lb" {
  configured_status = "started"
  name              = "lb-test"
  plan              = "development"
  zone              = var.lb_zone
  networks {
    name    = "Private-Net"
    type    = "private"
    family  = "IPv4"
    network = resource.upcloud_network.lb_network.id
  }
  networks {
    name   = "Public-Net"
    type   = "public"
    family = "IPv4"
  }
}

resource "upcloud_loadbalancer_backend" "lb_be_1" {
  loadbalancer = resource.upcloud_loadbalancer.lb.id
  name         = "lb-be-1-test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `default_backend_name` (String) The name of the default backend where traffic will be routed. Note, default backend can be overwritten in frontend rules.
- `loadbalancer` (String) ID of the load balancer to which the frontend is connected.
- `mode` (String) When load balancer operating in `tcp` mode it acts as a layer 4 proxy. In `http` mode it acts as a layer 7 proxy.
- `name` (String) The name of the frontend must be unique within the load balancer service.
- `port` (Number) Port to listen incoming requests

### Optional

- `networks` (Block List) Networks that frontend will be listening. Networks are required if load balancer has `networks` defined. This field will be required when deprecated field `network` is removed from load balancer resource. (see [below for nested schema](#nestedblock--networks))
- `properties` (Block List, Max: 1) Frontend properties. Properties can set back to defaults by defining empty `properties {}` block. (see [below for nested schema](#nestedblock--properties))

### Read-Only

- `id` (String) The ID of this resource.
- `rules` (List of String) Set of frontend rule names
- `tls_configs` (List of String) Set of TLS config names

<a id="nestedblock--networks"></a>
### Nested Schema for `networks`

Required:

- `name` (String) Name of the load balancer network


<a id="nestedblock--properties"></a>
### Nested Schema for `properties`

Optional:

- `inbound_proxy_protocol` (Boolean) Enable or disable inbound proxy protocol support.
- `timeout_client` (Number) Client request timeout in seconds.


