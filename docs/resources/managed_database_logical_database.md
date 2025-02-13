---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "upcloud_managed_database_logical_database Resource - terraform-provider-upcloud"
subcategory: ""
description: |-
  This resource represents a logical database in managed database
---

# upcloud_managed_database_logical_database (Resource)

This resource represents a logical database in managed database

## Example Usage

```terraform
# PostgreSQL managed database with additional logical database: example_db 
resource "upcloud_managed_database_postgresql" "example" {
  name  = "postgres"
  plan  = "1x1xCPU-2GB-25GB"
  title = "postgres"
  zone  = "fi-hel1"
}

resource "upcloud_managed_database_logical_database" "example_db" {
  service = upcloud_managed_database_postgresql.example.id
  name    = "example_db"
}

# MySQL managed database with additional logical database: example2_db 
resource "upcloud_managed_database_mysql" "example" {
  name = "mymysql"
  plan = "1x1xCPU-2GB-25GB"
  zone = "fi-hel1"
}

resource "upcloud_managed_database_logical_database" "example2_db" {
  service = upcloud_managed_database_mysql.example.id
  name    = "example2_db"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the logical database
- `service` (String) Service's UUID for which this user belongs to

### Optional

- `character_set` (String) Default character set for the database (LC_CTYPE)
- `collation` (String) Default collation for the database (LC_COLLATE)

### Read-Only

- `id` (String) The ID of this resource.


