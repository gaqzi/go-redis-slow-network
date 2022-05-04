go-redis on a slow network
==========================

We had a problem recently where our AWS ElastiCache instances hit 
the `NetworkBandwidthOutAllowanceExceeded` [limit] and at the same time our 
connections to Redis increased by a lot. Looking at the docs for go-redis 
for ["Large number of open connections"][large-conn] we saw:

> Under high load, some commands will time out and go-redis will close such 
> connections because they can still receive some data later and can’t be 
> reused. Closed connections are first put into TIME_WAIT state and remain 
> there for double maximum segment life time which is usually 1 minute   

This seemed to explain why our connection count went up. But I wanted to 
verify that this is what happened, so this is my attempt at replicating 
this by running two different test scenarios: 

1. [test-normal.sh]: there is no network obstruction, should create 
   `n*workers` connections
2. [test-slow.sh]: add an [outgoing network delay of 100ms][delay], this 
   should create a new connection for every timed out call.

For both scenarios I have configured 10 concurrent workers and a read 
timeout of 50ms. To verify the number of connections we're checking Redis' 
[INFO] which counts `total_connections_received`.

[limit]: https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/TroubleshootingConnections.html#Network-limits
[large-conn]: https://redis.uptrace.dev/guide/go-redis-debugging.html#large-number-of-open-connections=
[delay]: https://medium.com/@kazushi/simulate-high-latency-network-using-docker-containerand-tc-commands-a3e503ea4307
[test-normal.sh]: ./test-normal.sh
[test-slow.sh]: ./test-slow.sh
[INFO]: https://kapeli.com/dash_share?docset_file=Redis&docset_name=Redis&path=commands/info.html&platform=redis&repo=Main&source=redis.io/commands/info&version=6.2.3

## Results

### test-normal.sh

<details>
<summary>`total_connections_received:11`</summary>

```shell
$ ./test-normal.sh
[+] Running 2/2
 ⠿ Network go-redis-slow-network_default         Created                                                            0.0s
 ⠿ Container go-redis-slow-network-slow_redis-1  Started                                                            0.4s
[+] Running 2/0
 ⠿ Container go-redis-slow-network-slow_redis-1  Running                                                            0.0s
 ⠿ Container go-redis-slow-network-client-1      Created                                                            0.0s
Attaching to go-redis-slow-network-client-1
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Connecting to redis: slow_redis:6379
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 0
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 1
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 2
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 3
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 4
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 5
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 6
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 7
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 8
go-redis-slow-network-client-1  | 2022/05/04 03:48:34 Starting worker 9
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 1
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 2
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 4
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 6
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 7
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 0
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 8
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 9
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 3
go-redis-slow-network-client-1  | 2022/05/04 03:48:39 caught ctx.Done in 5
go-redis-slow-network-client-1 exited with code 0
# Server
redis_version:5.0.7
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:66bd629f924ac924
redis_mode:standalone
os:Linux 5.10.104-linuxkit x86_64
arch_bits:64
multiplexing_api:epoll
atomicvar_api:atomic-builtin
gcc_version:9.3.0
process_id:1
run_id:35100fe5036857b197a728fe684638376e8a44e3
tcp_port:6379
uptime_in_seconds:6
uptime_in_days:0
hz:10
configured_hz:10
lru_clock:7468952
executable:/work/redis-server
config_file:

# Clients
connected_clients:1
client_recent_max_input_buffer:2
client_recent_max_output_buffer:0
blocked_clients:0

# Memory
used_memory:859064
used_memory_human:838.93K
used_memory_rss:6631424
used_memory_rss_human:6.32M
used_memory_peak:1416296
used_memory_peak_human:1.35M
used_memory_peak_perc:60.66%
used_memory_overhead:845774
used_memory_startup:796080
used_memory_dataset:13290
used_memory_dataset_perc:21.10%
allocator_allocated:1170896
allocator_active:1474560
allocator_resident:4882432
total_system_memory:7290019840
total_system_memory_human:6.79G
used_memory_lua:41984
used_memory_lua_human:41.00K
used_memory_scripts:0
used_memory_scripts_human:0B
number_of_cached_scripts:0
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:1.26
allocator_frag_bytes:303664
allocator_rss_ratio:3.31
allocator_rss_bytes:3407872
rss_overhead_ratio:1.36
rss_overhead_bytes:1748992
mem_fragmentation_ratio:8.33
mem_fragmentation_bytes:5835256
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:49694
mem_aof_buffer:0
mem_allocator:jemalloc-5.2.1
active_defrag_running:0
lazyfree_pending_objects:0

# Persistence
loading:0
rdb_changes_since_last_save:10
rdb_bgsave_in_progress:0
rdb_last_save_time:1651636114
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:-1
rdb_current_bgsave_time_sec:-1
rdb_last_cow_size:0
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok
aof_last_cow_size:0

# Stats
total_connections_received:11
total_commands_processed:356021
instantaneous_ops_per_sec:59833
total_net_input_bytes:8900799
total_net_output_bytes:3916171
instantaneous_input_kbps:1460.78
instantaneous_output_kbps:642.78
rejected_connections:0
sync_full:0
sync_partial_ok:0
sync_partial_err:0
expired_keys:10
expired_stale_perc:4.51
expired_time_cap_reached_count:0
evicted_keys:0
keyspace_hits:356011
keyspace_misses:0
pubsub_channels:0
pubsub_patterns:0
latest_fork_usec:0
migrate_cached_sockets:0
slave_expires_tracked_keys:0
active_defrag_hits:0
active_defrag_misses:0
active_defrag_key_hits:0
active_defrag_key_misses:0

# Replication
role:master
connected_slaves:0
master_replid:784ecf7071754be601c8e149f5039be0de3f8437
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

# CPU
used_cpu_sys:3.342943
used_cpu_user:0.635310
used_cpu_sys_children:0.000000
used_cpu_user_children:0.000000

# Cluster
cluster_enabled:0

# Keyspace
[+] Running 3/3
 ⠿ Container go-redis-slow-network-client-1      Removed                                                                                                                                                                                         0.0s
 ⠿ Container go-redis-slow-network-slow_redis-1  Removed                                                                                                                                                                                         0.2s
 ⠿ Network go-redis-slow-network_default         Removed                                                                                                                                                                                         0.1s
Look at total_connections_received:!
DONE
```
</details>

### test-slow.sh

<details>
<summary>`total_connections_received:321`</summary>

```shell
$ ./test-slow.sh
[+] Running 2/2
 ⠿ Network go-redis-slow-network_default         Created                                                                                                                                                                                         0.0s
 ⠿ Container go-redis-slow-network-slow_redis-1  Started                                                                                                                                                                                         0.4s
[+] Running 2/0
 ⠿ Container go-redis-slow-network-slow_redis-1  Running                                                                                                                                                                                         0.0s
 ⠿ Container go-redis-slow-network-client-1      Created                                                                                                                                                                                         0.0s
Attaching to go-redis-slow-network-client-1
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Connecting to redis: slow_redis:6379
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 0
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 1
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 2
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 3
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 4
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 5
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 6
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 7
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 8
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 Starting worker 9
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 1: read tcp 172.31.0.3:52196->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 3: read tcp 172.31.0.3:52212->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 4: read tcp 172.31.0.3:52200->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 6: read tcp 172.31.0.3:52202->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 9: read tcp 172.31.0.3:52214->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 0: read tcp 172.31.0.3:52210->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 7: read tcp 172.31.0.3:52206->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 8: read tcp 172.31.0.3:52204->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 2: read tcp 172.31.0.3:52208->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 5: read tcp 172.31.0.3:52198->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 5: read tcp 172.31.0.3:52230->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 4: read tcp 172.31.0.3:52224->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 8: read tcp 172.31.0.3:52220->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 2: read tcp 172.31.0.3:52232->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 6: read tcp 172.31.0.3:52226->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 0: read tcp 172.31.0.3:52216->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 7: read tcp 172.31.0.3:52234->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 3: read tcp 172.31.0.3:52222->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 1: read tcp 172.31.0.3:52228->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 9: read tcp 172.31.0.3:52218->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 1: read tcp 172.31.0.3:52246->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 9: read tcp 172.31.0.3:52236->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 5: read tcp 172.31.0.3:52250->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 7: read tcp 172.31.0.3:52238->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 3: read tcp 172.31.0.3:52240->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 0: read tcp 172.31.0.3:52244->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 6: read tcp 172.31.0.3:52242->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 8: read tcp 172.31.0.3:52254->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 4: read tcp 172.31.0.3:52248->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:32 caught error in routine 2: read tcp 172.31.0.3:52252->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 5: read tcp 172.31.0.3:52262->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 3: read tcp 172.31.0.3:52258->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 9: read tcp 172.31.0.3:52260->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 2: read tcp 172.31.0.3:52256->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 4: read tcp 172.31.0.3:52274->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 7: read tcp 172.31.0.3:52266->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 1: read tcp 172.31.0.3:52264->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 0: read tcp 172.31.0.3:52270->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 6: read tcp 172.31.0.3:52268->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 8: read tcp 172.31.0.3:52272->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 6: read tcp 172.31.0.3:52286->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 1: read tcp 172.31.0.3:52282->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 9: read tcp 172.31.0.3:52292->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 3: read tcp 172.31.0.3:52290->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 4: read tcp 172.31.0.3:52278->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 5: read tcp 172.31.0.3:52288->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 0: read tcp 172.31.0.3:52280->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 8: read tcp 172.31.0.3:52276->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 7: read tcp 172.31.0.3:52284->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 2: read tcp 172.31.0.3:52294->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 2: read tcp 172.31.0.3:52296->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 6: read tcp 172.31.0.3:52304->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 9: read tcp 172.31.0.3:52308->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 1: read tcp 172.31.0.3:52300->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 3: read tcp 172.31.0.3:52312->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 5: read tcp 172.31.0.3:52298->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 4: read tcp 172.31.0.3:52302->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 0: read tcp 172.31.0.3:52314->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 7: read tcp 172.31.0.3:52306->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 8: read tcp 172.31.0.3:52310->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 0: read tcp 172.31.0.3:52334->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 5: read tcp 172.31.0.3:52320->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 4: read tcp 172.31.0.3:52322->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 1: read tcp 172.31.0.3:52332->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 8: read tcp 172.31.0.3:52316->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 7: read tcp 172.31.0.3:52324->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 3: read tcp 172.31.0.3:52318->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 2: read tcp 172.31.0.3:52326->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 6: read tcp 172.31.0.3:52330->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 9: read tcp 172.31.0.3:52328->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 9: read tcp 172.31.0.3:52354->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 4: read tcp 172.31.0.3:52348->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 2: read tcp 172.31.0.3:52342->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 6: read tcp 172.31.0.3:52336->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 8: read tcp 172.31.0.3:52340->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 0: read tcp 172.31.0.3:52350->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 5: read tcp 172.31.0.3:52346->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 1: read tcp 172.31.0.3:52352->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 3: read tcp 172.31.0.3:52338->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 7: read tcp 172.31.0.3:52344->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 8: read tcp 172.31.0.3:52358->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 7: read tcp 172.31.0.3:52356->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 2: read tcp 172.31.0.3:52370->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 1: read tcp 172.31.0.3:52364->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 5: read tcp 172.31.0.3:52366->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 0: read tcp 172.31.0.3:52362->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 6: read tcp 172.31.0.3:52368->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 9: read tcp 172.31.0.3:52374->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 3: read tcp 172.31.0.3:52372->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:33 caught error in routine 4: read tcp 172.31.0.3:52360->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52384->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52390->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52392->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52388->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52386->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52376->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52394->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52378->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52380->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52396->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52408->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52414->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52410->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52402->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52406->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52412->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52400->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52404->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52398->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52416->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52424->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52436->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52432->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52428->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52418->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52426->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52430->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52420->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52434->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52422->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52444->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52446->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52440->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52442->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52448->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52452->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52438->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52450->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52454->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52456->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52462->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52476->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52470->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52466->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52472->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52474->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52460->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52468->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52464->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52458->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52478->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52482->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52480->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52484->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52486->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52494->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52496->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52488->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52490->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52492->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 4: read tcp 172.31.0.3:52502->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 9: read tcp 172.31.0.3:52512->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 6: read tcp 172.31.0.3:52498->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 5: read tcp 172.31.0.3:52514->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 1: read tcp 172.31.0.3:52508->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 8: read tcp 172.31.0.3:52506->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 0: read tcp 172.31.0.3:52500->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 3: read tcp 172.31.0.3:52504->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 2: read tcp 172.31.0.3:52516->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:34 caught error in routine 7: read tcp 172.31.0.3:52510->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 9: read tcp 172.31.0.3:52528->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 5: read tcp 172.31.0.3:52532->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 0: read tcp 172.31.0.3:52524->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 1: read tcp 172.31.0.3:52520->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 3: read tcp 172.31.0.3:52518->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 8: read tcp 172.31.0.3:52522->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 4: read tcp 172.31.0.3:52526->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 2: read tcp 172.31.0.3:52536->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 6: read tcp 172.31.0.3:52530->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 7: read tcp 172.31.0.3:52534->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 0: read tcp 172.31.0.3:52548->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 5: read tcp 172.31.0.3:52546->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 2: read tcp 172.31.0.3:52554->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 6: read tcp 172.31.0.3:52556->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 8: read tcp 172.31.0.3:52552->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 3: read tcp 172.31.0.3:52542->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 9: read tcp 172.31.0.3:52538->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 7: read tcp 172.31.0.3:52540->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 4: read tcp 172.31.0.3:52544->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 1: read tcp 172.31.0.3:52550->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 3: read tcp 172.31.0.3:52560->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 0: read tcp 172.31.0.3:52568->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 4: read tcp 172.31.0.3:52566->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 5: read tcp 172.31.0.3:52572->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 1: read tcp 172.31.0.3:52576->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 8: read tcp 172.31.0.3:52562->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 9: read tcp 172.31.0.3:52564->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 2: read tcp 172.31.0.3:52570->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 7: read tcp 172.31.0.3:52558->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 6: read tcp 172.31.0.3:52574->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 0: read tcp 172.31.0.3:52582->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 4: read tcp 172.31.0.3:52586->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 8: read tcp 172.31.0.3:52578->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 3: read tcp 172.31.0.3:52580->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 6: read tcp 172.31.0.3:52590->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 9: read tcp 172.31.0.3:52594->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 5: read tcp 172.31.0.3:52584->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 1: read tcp 172.31.0.3:52588->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 7: read tcp 172.31.0.3:52596->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 2: read tcp 172.31.0.3:52592->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 4: read tcp 172.31.0.3:52604->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 0: read tcp 172.31.0.3:52598->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 8: read tcp 172.31.0.3:52600->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 9: read tcp 172.31.0.3:52612->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 5: read tcp 172.31.0.3:52608->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 6: read tcp 172.31.0.3:52618->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 2: read tcp 172.31.0.3:52606->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 1: read tcp 172.31.0.3:52610->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 3: read tcp 172.31.0.3:52616->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 7: read tcp 172.31.0.3:52614->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 8: read tcp 172.31.0.3:52620->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 0: read tcp 172.31.0.3:52624->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 4: read tcp 172.31.0.3:52622->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 2: read tcp 172.31.0.3:52634->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 6: read tcp 172.31.0.3:52638->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 1: read tcp 172.31.0.3:52632->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 3: read tcp 172.31.0.3:52636->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 7: read tcp 172.31.0.3:52626->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 9: read tcp 172.31.0.3:52628->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:35 caught error in routine 5: read tcp 172.31.0.3:52630->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 8: read tcp 172.31.0.3:52648->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 6: read tcp 172.31.0.3:52646->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 2: read tcp 172.31.0.3:52644->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 4: read tcp 172.31.0.3:52642->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 1: read tcp 172.31.0.3:52640->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 0: read tcp 172.31.0.3:52650->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 3: read tcp 172.31.0.3:52652->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 9: read tcp 172.31.0.3:52654->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 5: read tcp 172.31.0.3:52658->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 7: read tcp 172.31.0.3:52656->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 2: read tcp 172.31.0.3:52664->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 4: read tcp 172.31.0.3:52666->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 9: read tcp 172.31.0.3:52674->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 5: read tcp 172.31.0.3:52676->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 8: read tcp 172.31.0.3:52660->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 0: read tcp 172.31.0.3:52670->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 1: read tcp 172.31.0.3:52668->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 3: read tcp 172.31.0.3:52678->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 6: read tcp 172.31.0.3:52662->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 7: read tcp 172.31.0.3:52672->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 4: read tcp 172.31.0.3:52684->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 2: read tcp 172.31.0.3:52696->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 8: read tcp 172.31.0.3:52680->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 9: read tcp 172.31.0.3:52688->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 3: read tcp 172.31.0.3:52692->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 5: read tcp 172.31.0.3:52694->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 0: read tcp 172.31.0.3:52686->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 1: read tcp 172.31.0.3:52690->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 7: read tcp 172.31.0.3:52682->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 6: read tcp 172.31.0.3:52698->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 2: read tcp 172.31.0.3:52708->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 1: read tcp 172.31.0.3:52710->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 6: read tcp 172.31.0.3:52700->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 7: read tcp 172.31.0.3:52712->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 5: read tcp 172.31.0.3:52702->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 8: read tcp 172.31.0.3:52714->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 4: read tcp 172.31.0.3:52716->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 9: read tcp 172.31.0.3:52718->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 0: read tcp 172.31.0.3:52706->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 3: read tcp 172.31.0.3:52704->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 9: read tcp 172.31.0.3:52734->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 8: read tcp 172.31.0.3:52728->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 5: read tcp 172.31.0.3:52726->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 2: read tcp 172.31.0.3:52722->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 1: read tcp 172.31.0.3:52724->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 6: read tcp 172.31.0.3:52730->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 0: read tcp 172.31.0.3:52720->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 4: read tcp 172.31.0.3:52732->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 3: read tcp 172.31.0.3:52738->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 7: read tcp 172.31.0.3:52736->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 8: read tcp 172.31.0.3:52752->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 3: read tcp 172.31.0.3:52750->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 4: read tcp 172.31.0.3:52748->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 0: read tcp 172.31.0.3:52746->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 7: read tcp 172.31.0.3:52740->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 1: read tcp 172.31.0.3:52742->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 2: read tcp 172.31.0.3:52758->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 9: read tcp 172.31.0.3:52754->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 5: read tcp 172.31.0.3:52756->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:36 caught error in routine 6: read tcp 172.31.0.3:52744->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 4: read tcp 172.31.0.3:52762->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 3: read tcp 172.31.0.3:52764->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 8: read tcp 172.31.0.3:52760->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 7: read tcp 172.31.0.3:52768->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 9: read tcp 172.31.0.3:52766->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 0: read tcp 172.31.0.3:52772->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 1: read tcp 172.31.0.3:52774->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 2: read tcp 172.31.0.3:52770->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 6: read tcp 172.31.0.3:52776->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 5: read tcp 172.31.0.3:52778->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 4: read tcp 172.31.0.3:52784->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 8: read tcp 172.31.0.3:52780->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 3: read tcp 172.31.0.3:52782->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 1: read tcp 172.31.0.3:52786->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 7: read tcp 172.31.0.3:52792->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 0: read tcp 172.31.0.3:52788->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 5: read tcp 172.31.0.3:52794->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 2: read tcp 172.31.0.3:52796->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 6: read tcp 172.31.0.3:52798->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 9: read tcp 172.31.0.3:52790->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 4: read tcp 172.31.0.3:52802->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 4
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 7: read tcp 172.31.0.3:52800->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 7
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 3: read tcp 172.31.0.3:52804->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 1: read tcp 172.31.0.3:52808->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 3
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 1
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 8: read tcp 172.31.0.3:52806->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 8
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 6: read tcp 172.31.0.3:52818->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 6
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 5: read tcp 172.31.0.3:52810->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 5
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 9: read tcp 172.31.0.3:52814->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 9
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 0: read tcp 172.31.0.3:52812->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 0
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught error in routine 2: read tcp 172.31.0.3:52816->172.31.0.2:6379: i/o timeout
go-redis-slow-network-client-1  | 2022/05/04 03:51:37 caught ctx.Done in 2
go-redis-slow-network-client-1 exited with code 0
# Server
redis_version:5.0.7
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:66bd629f924ac924
redis_mode:standalone
os:Linux 5.10.104-linuxkit x86_64
arch_bits:64
multiplexing_api:epoll
atomicvar_api:atomic-builtin
gcc_version:9.3.0
process_id:1
run_id:0e7b0693335a036ae105d42f991f491115ae74e5
tcp_port:6379
uptime_in_seconds:6
uptime_in_days:0
hz:10
configured_hz:10
lru_clock:7469129
executable:/work/redis-server
config_file:

# Clients
connected_clients:1
client_recent_max_input_buffer:4
client_recent_max_output_buffer:0
blocked_clients:0

# Memory
used_memory:859064
used_memory_human:838.93K
used_memory_rss:6426624
used_memory_rss_human:6.13M
used_memory_peak:1415976
used_memory_peak_human:1.35M
used_memory_peak_perc:60.67%
used_memory_overhead:845774
used_memory_startup:796080
used_memory_dataset:13290
used_memory_dataset_perc:21.10%
allocator_allocated:1311216
allocator_active:1548288
allocator_resident:4902912
total_system_memory:7290019840
total_system_memory_human:6.79G
used_memory_lua:41984
used_memory_lua_human:41.00K
used_memory_scripts:0
used_memory_scripts_human:0B
number_of_cached_scripts:0
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:1.18
allocator_frag_bytes:237072
allocator_rss_ratio:3.17
allocator_rss_bytes:3354624
rss_overhead_ratio:1.31
rss_overhead_bytes:1523712
mem_fragmentation_ratio:8.07
mem_fragmentation_bytes:5630456
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_clients_slaves:0
mem_clients_normal:49694
mem_aof_buffer:0
mem_allocator:jemalloc-5.2.1
active_defrag_running:0
lazyfree_pending_objects:0

# Persistence
loading:0
rdb_changes_since_last_save:10
rdb_bgsave_in_progress:0
rdb_last_save_time:1651636291
rdb_last_bgsave_status:ok
rdb_last_bgsave_time_sec:-1
rdb_current_bgsave_time_sec:-1
rdb_last_cow_size:0
aof_enabled:0
aof_rewrite_in_progress:0
aof_rewrite_scheduled:0
aof_last_rewrite_time_sec:-1
aof_current_rewrite_time_sec:-1
aof_last_bgrewrite_status:ok
aof_last_write_status:ok
aof_last_cow_size:0

# Stats
total_connections_received:321
total_commands_processed:320
instantaneous_ops_per_sec:55
total_net_input_bytes:8274
total_net_output_bytes:3460
instantaneous_input_kbps:1.36
instantaneous_output_kbps:0.60
rejected_connections:0
sync_full:0
sync_partial_ok:0
sync_partial_err:0
expired_keys:10
expired_stale_perc:4.51
expired_time_cap_reached_count:0
evicted_keys:0
keyspace_hits:310
keyspace_misses:0
pubsub_channels:0
pubsub_patterns:0
latest_fork_usec:0
migrate_cached_sockets:0
slave_expires_tracked_keys:0
active_defrag_hits:0
active_defrag_misses:0
active_defrag_key_hits:0
active_defrag_key_misses:0

# Replication
role:master
connected_slaves:0
master_replid:6f10be4ec5dac095263b948ef7edbfb00a736713
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0

# CPU
used_cpu_sys:0.075459
used_cpu_user:0.044867
used_cpu_sys_children:0.000000
used_cpu_user_children:0.000000

# Cluster
cluster_enabled:0

# Keyspace
[+] Running 3/3
 ⠿ Container go-redis-slow-network-client-1      Removed                                                                                                                                                                                         0.0s
 ⠿ Container go-redis-slow-network-slow_redis-1  Removed                                                                                                                                                                                         0.3s
 ⠿ Network go-redis-slow-network_default         Removed                                                                                                                                                                                         0.1s
Look at total_connections_received:!
DONE
```

</details>
