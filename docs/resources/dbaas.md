---
page_title: "infomaniak_dbaas"
subcategory: "DBaaS"
description: |-
  The DBaas resource allows the user to manage a DBaas project
---

# infomaniak_dbaas

The DBaas resource allows the user to manage a DBaas project.

To get your `public_cloud_id`:
```sh
account_id=$(curl -s -H "Authorization: Bearer $INFOMANIAK_TOKEN" https://api.infomaniak.com/2/profile | jq '.data.preferences.account.current_account_id')
curl -s -H "Authorization: Bearer $INFOMANIAK_TOKEN" https://api.infomaniak.com/1/public_clouds?account_id=$account_id | jq '.data[] | {"name": .customer_name, "cloud_id": .id}'
```

To get your `public_cloud_project_id`:
```sh
public_cloud_id=1234  # use the ID retrieved from the step above
curl -s -H "Authorization: Bearer $INFOMANIAK_TOKEN" https://api.infomaniak.com/1/public_clouds/$public_cloud_id/projects | jq '.data[] | {"name": .name, "project_id": .public_cloud_project_id}'
```

## Example

```hcl
resource "infomaniak_dbaas" "db-0" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  
  name      = "db-0"
  pack_name = "pro-4"
  type      = "mysql"
  version   = "8.0.42"
  region    = "dc4-a"

  allowed_cidrs = [
    "162.1.15.122/32",
    "1.1.1.1",
    "2345:425:2CA1:0000:0000:567:5673:23b5/64",
  ]

  mysqlconfig = {
    connect_timeout = 10
  }
}
```

Be careful with `allowed_cidrs`:
```hcl
resource "infomaniak_dbaas" "db-0" {
  public_cloud_id = xxxxx
  public_cloud_project_id = yyyyy
  
  name      = "db-0"
  pack_name = "pro-4"
  type      = "mysql"
  version   = "8.0.42"
  region    = "dc4-a"

  allowed_cidrs = [] // If you set an empty list here, it means no one can access it ! Even you !

  mysqlconfig = {
    connect_timeout = 10
  }
}
```

## Schema

### Required

- `public_cloud_id` (Integer) The id of the Public Cloud where DBaaS is installed.
- `public_cloud_project_id` (Integer) The id of the public cloud project where DBaaS is installed.
- `region` (String) Region where the instance live.
- `pack_name` (String) The name of the pack corresponding the DBaaS project.
- `type` (String) The type of the database to use.
- `version` (String) The version of the database to use.
- `name` (String) The name of the DBaaS shown on the manager.
- `allowed_cidrs` (List of String) The list of allowed cidrs to access to the database.
- `mysqlconfig` (Object) MySQL configuration parameters. Please refer to [this section](#mysql-configuration-attributes)

### Read-Only

- `id` (Integer) A computed value representing the unique identifier for the architecture. Mandatory for acceptance testing.
- `kube_identifier` (String) A computed value that gives the kubernetes identifier of the DbaaS
- `host` (String) The host to access the Database.
- `port` (String) The port to access the Database.
- `user` (String) The user to access the Database.
- `password` (String, Sensitive) The password to access the Database.
- `ca` (String) The database CA certificate.

## MySQL Configuration Attributes

The `mysqlconfig` block supports the following attributes:

- `auto_increment_increment` (Integer) AUTO_INCREMENT increment value.
- `auto_increment_offset` (Integer) AUTO_INCREMENT offset value.
- `character_set_server` (String) Default character set for the server.
- `connect_timeout` (Integer) Timeout for establishing a connection.
- `group_concat_max_len` (Integer) Maximum length of GROUP_CONCAT() result.
- `information_schema_stats_expiry` (Integer) Expiration time for information schema statistics.
- `innodb_change_buffer_max_size` (Integer) Maximum size of the InnoDB change buffer.
- `innodb_flush_neighbors` (Integer) Whether to flush neighbor pages when flushing a page.
- `innodb_ft_max_token_size` (Integer) Maximum token size for InnoDB full-text search.
- `innodb_ft_min_token_size` (Integer) Minimum token size for InnoDB full-text search.
- `innodb_ft_server_stopword_table` (String) Server-wide stopword table for InnoDB full-text search.
- `innodb_lock_wait_timeout` (Integer) Timeout for InnoDB lock waits.
- `innodb_log_buffer_size` (Integer) Size of the InnoDB log buffer.
- `innodb_online_alter_log_max_size` (Integer) Maximum size of the online alter log.
- `innodb_print_all_deadlocks` (String) Whether to print all deadlocks to the error log.
- `innodb_read_io_threads` (Integer) Number of InnoDB read I/O threads.
- `innodb_rollback_on_timeout` (String) Whether to rollback transactions on lock wait timeout.
- `innodb_stats_persistent_sample_pages` (Integer) Number of index pages sampled for persistent stats.
- `innodb_thread_concurrency` (Integer) Maximum number of concurrent threads.
- `innodb_write_io_threads` (Integer) Number of InnoDB write I/O threads.
- `interactive_timeout` (Integer) Timeout for interactive connections.
- `lock_wait_timeout` (Integer) Timeout for all lock waits.
- `log_bin_trust_function_creators` (String) Whether to trust function creators for binary logging.
- `long_query_time` (Float) Threshold for slow query logging.
- `max_allowed_packet` (Integer) Maximum packet size allowed.
- `max_connections` (Integer) Maximum number of simultaneous connections.
- `max_digest_length` (Integer) Maximum digest length for statement digest.
- `max_heap_table_size` (Integer) Maximum size of user-created MEMORY tables.
- `max_prepared_stmt_count` (Integer) Maximum number of prepared statements.
- `min_examined_row_limit` (Integer) Minimum examined row limit for query logging.
- `net_buffer_length` (Integer) Buffer size for TCP/IP and socket communication.
- `net_read_timeout` (Integer) Timeout for reading from a connection.
- `net_write_timeout` (Integer) Timeout for writing to a connection.
- `performance_schema_max_digest_length` (Integer) Maximum digest length for Performance Schema.
- `require_secure_transport` (String) Whether secure transport is required.
- `sort_buffer_size` (Integer) Sort buffer size per session.
- `sql_mode` (List of String) SQL mode settings.
- `table_definition_cache` (Integer) Number of table definitions that can be cached.
- `table_open_cache` (Integer) Number of open tables cache instances.
- `table_open_cache_instances` (Integer) Number of table open cache instances.
- `thread_stack` (Integer) Stack size for each thread.
- `transaction_isolation` (String) Default transaction isolation level.
- `wait_timeout` (Integer) Timeout for non-interactive connections.
