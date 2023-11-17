# Introduction

This is a CLI tool for [Seata](https://github.com/seata/seata) named `seata-ctl`.

```shell
$ seata-ctl -h
seata-ctl is a CLI tool for Seata

Usage:
  seata-ctl [flags]
  seata-ctl [command]

Available Commands:
  version     Print the version number of seata-ctl

Flags:
  -h, --help              help for seata-ctl
      --ip string         Seata Server IP (default "127.0.0.1")
      --password string   Password (default "seata")
      --port int          Seata Server Admin Port (default 7091)
      --username string   Username (default "seata")

Use "seata-ctl [command] --help" for more information about a command.
```

# Build

Dependency: Go 1.19+

```shell
go build .
```

# How to use

## Login

```shell
$ seata-ctl --ip 127.0.0.1 --port 7091 --username seata --password seata
127.0.0.1:7091 > # input command here
```

## Help

```shell
127.0.0.1:7091 > help
Usage:
  [command] [flag] 

Available Commands:
  get         Get the resource
  help        Help about any command
  quit        Quit the session
  reload      Reload the configuration
  set         Set the resource
  try         Try example transactions
```

## Get

### Usage

```shell
127.0.0.1:7091 > get -h    
Get the resource

Usage:
   get [flags]
   get [command]

Available Commands:
  config        Get the configuration
  status        Get the status

Flags:
  -h, --help   help for get

Use "get [command] --help" for more information about a command.
```

### Example

1. Get the status of the Seata server cluster:

```shell
127.0.0.1:7091 > get status
+-------+--------------------+--------+
| TYPE  | ADDRESS            | STATUS |
+-------+--------------------+--------+
| nacos | 192.168.163.1:7091 | ok     |
+-------+--------------------+--------+
| nacos | 192.168.163.2:7091 | ok     |
+-------+--------------------+--------+
```

2. Get the configuration `server.servicePort`:

```shell
127.0.0.1:7091 > get config --key '["server.servicePort"]'
+--------------------+-------+
| KEY                | VALUE |
+--------------------+-------+
| server.servicePort | 8091  |
+--------------------+-------+
```

3. Get multiple configuration (which could be written in multiple lines by ending with '\\'):

```shell
127.0.0.1:7091 > get config --key '[ \
    "server.servicePort", \
    "server.recovery.timeoutRetryPeriod", \
    "server.undo.logDeletePeriod" \
]'
+------------------------------------+----------+
| KEY                                | VALUE    |
+------------------------------------+----------+
| server.recovery.timeoutRetryPeriod | 1000     |
| server.servicePort                 | 8091     |
| server.undo.logDeletePeriod        | 86400000 |
+------------------------------------+----------+
```

4. Get all configurations

```shell
127.0.0.1:7091 > get config
+------------------------------------+----------+
| KEY                                | VALUE    |
+------------------------------------+----------+
...
```

## Set

### Usage

```shell
127.0.0.1:7091 > set -h                    
Set the resource

Usage:
   set [flags]
   set [command]

Available Commands:
  config        Set the configuration

Flags:
  -h, --help   help for set

Use "set [command] --help" for more information about a command.
```

```shell
127.0.0.1:7091 > set config -h
Set the configuration

Usage:
   set config [flags]

Flags:
      --config-center   If set configuration center conf
      --data string     Configuration map (default "{}")
  -h, --help            help for config
      --registry        If set registry conf
```

### Example

1. Set the registry config, such as setting type to `eureka`:

```shell
127.0.0.1:7091 > set config --registry --data '{"registry.type": "eureka"}'
+---------------+------+--------+
| KEY           | FROM | TO     |
+---------------+------+--------+
| registry.type | file | eureka |
+---------------+------+--------+
```

You can found that the Seata server is registered at `eureka` registry.

2. Set the configuration center config, such as setting type to `nacos`

```shell
127.0.0.1:7091 > set config --config-center --data '{"config.type": "nacos"}'
+-------------+------+-------+
| KEY         | FROM | TO    |
+-------------+------+-------+
| config.type | file | nacos |
+-------------+------+-------+
```

You can found that the configuration in `nacos` is loaded.

3. Set a configuration item which can be dynamically configured (such as `server.undo.logSaveDays`):

```shell
127.0.0.1:7091 > set config --data '{"server.undo.logSaveDays": "5"}'
+-------------------------+------+----+
| KEY                     | FROM | TO |
+-------------------------+------+----+
| server.undo.logSaveDays | 6    | 5  |
+-------------------------+------+----+
```

4. Set multiple configurations at the same time:

```shell
127.0.0.1:7091 > set config --data '{ \
    "server.maxCommitRetryTimeout": "3000", \
    "server.maxRollbackRetryTimeout": "3000", \
    "server.undo.logSaveDays": "14" \
}'
+--------------------------------+------+------+
| KEY                            | FROM | TO   |
+--------------------------------+------+------+
| server.maxCommitRetryTimeout   | -1   | 3000 |
| server.maxRollbackRetryTimeout | -1   | 3000 |
| server.undo.logSaveDays        | 5    | 14   |
+--------------------------------+------+------+
```

## Try

### Usage

```shell
127.0.0.1:7091 > try -h
Try if this node is ready

Usage:
   try [flags]
   try [command]

Available Commands:
  begin       begin a txn
  commit      commit a txn
  rollback    rollback a txn

Flags:
  -h, --help   help for try

Use "try [command] --help" for more information about a command.
```

### Example

1. Try to begin an example transaction:

```shell
127.0.0.1:7091 > try begin --timeout 300000
Try an example txn successfully, xid=192.168.163.1:8091:8755443813836259333
```

2. Commit a transaction by xid:

```shell
127.0.0.1:7091 > try commit --xid 192.168.163.1:8091:8755443813836259333
Commit txn successfully, xid=192.168.163.1:8091:8755443813836259333
```

3. Rollback a transaction by xid:

```shell
127.0.0.1:7091 > try rollback --xid 192.168.163.1:8091:8755443813836259333
Rollback txn successfully, xid=192.168.163.1:8091:8755443813836259333
```

Lifecycle of the transactions could be checked in Web console UI. (exposed at `7091` by default).

## Reload

### Usage

```shell
reload -h
Reload the configuration

Usage:
   reload [flags]

Flags:
  -h, --help   help for reload
```

### Example

```shell
127.0.0.1:7091 > reload
Reload Successful!
```

## Quit

Quit the session:

```shell
127.0.0.1:7091 > quit
Quit the session
```