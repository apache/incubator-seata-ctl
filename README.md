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
      --ip string         Seata Server IP (constants "127.0.0.1")
      --password string   Password (constants "seata")
      --port int          Seata Server Admin Port (constants 7091)
      --username string   Username (constants "seata")

Use "seata-ctl [command] --help" for more information about a command.
```

# Build

Dependency: Go 1.19+

```shell
go build .
```

# How to use

# Seata

Some Seata commands are used to connect to the Seata server for configuration browsing, transaction simulation, and other functions. To use these features, you must first log in using the `Login` command.

## Login

Use this command to log in to the seata section. If the login is successful, the login information will be displayed at the beginning of the command line.

```shell
$ seata-ctl login --ip=127.0.0.1 --port=7091 --username=seata --password=seata
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
      --data string     Configuration map (constants "{}")
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

# Config

Seata-ctl uses a YAML configuration file to set up various cluster addresses and connection details. Since manually creating this file can be complex for users, Seata-ctl provides a template for the configuration file.

## Config

### Describe

Used to initialize the configuration file with no parameters required.

After execution, a file named `config.yaml` will be created in the same directory as the terminal. If the file already exists, a message will indicate this. Currently, selecting a different location is not supported.

### Usage

```
seata-ctl
config
```

After use, a file named 'config. yaml' will be created in the same directory as the terminal. If the file already exists, it will prompt that it already exists. We currently do not support selecting other locations.

Each cluster can fill in multiple instances, using a context mechanism similar to kubeconfig to select which cluster and service to use.

The format of the file and the meanings of each field are as follows:(The parameters of the log section will be explained in detail later in the log section)

```
kubernetes:
    clusters:
        - name: ""	//kubernetes cluster name
          kubeconfigpath: ""	//kubeconfig path
          ymlpath: ""		//yml config path,not currently supported
prometheus:
    servers:
        - name: ""	//prometheus name
          address: ""	//prometheus address
          auth: ""	//prometheus auth info,not currently supported
log:
    clusters:
        - name: ""	//log name
          types: ""	//log type:ElasticSearch,Loki,Local
          address: ""	//log server address
          source: ""	//index name
          username: ""	//auth username
          password: ""	//auth password
          index: ""			//index name
context:
    kubernetes: ""  //choose which clusters to use,depend on name
    prometheus: "" 	//choose which clusters to use,depend on name
    log: ""					//choose which clusters to use,depend on name
```

### Example

If everything is normal, it will output:

```
seata-ctl
config

Config created successfully!
```

If the file already exists, it will be output as follows:

```
seata-ctl
config

Config file already exists!
```

The example of writing a configuration file is as follows：

```
kubernetes:
  clusters:
    - name: "cluster1"
      kubeconfigpath: "/Users/wpy/.kube/config"
      ymlpath: ""
    - name: "seata"
      kubeconfigpath: "/Users/wpy/Documents/Kubernetes/remotekube.txt"
      ymlpath: ""
prometheus:
  servers:
    - name: "prometheus"
      address: "http://localhost:9092"
      auth: ""
log:
  clusters:
    - name: "es"
      types: "ElasticSearch"
      address: "https://localhost:9200"
      source: "logstash-2024.10.24"
      username: "elastic"
      password: "bu4AC50REtt_7rUqddMe"
      index: "log"
    - name: "loki"
      types: "Loki"
      address: "http://localhost:3100"
      source: ""
      username: ""
      password: ""
      index: ""
    - name: "local"
      types: "Local"
      address: "http://localhost:8080"
      source: "seata"
      username: ""
      password: ""
      index: ""
context:
  kubernetes: "cluster1"
  prometheus: "prometheus"
  log: "es"
```

# Kubernetes

Seata ctl provides various commands to operate Seata servers on Kubernetes clusters.

## Install

### Describe

Use this command to apply CRD (Custom Resource Definition) and deploy controllers on a Kubernetes cluster` Seata tl 'will use these definitions to complete the deployment of Seata.

### Usage

```
seata-ctl
install --namespace=default --image=apache/seata-controller:latest --name=seata-k8s-controller-manager
```

#### Arguments

- `namespace`: Specifies the namespace to use (optional). If not specified, it defaults to `default`.
- `image`: Specifies the name of the image to use (optional). If not specified, it defaults to `apache/seata-controller:latest`.
- `name`: Specifies the name of the deployed controller (optional). If not specified, it defaults to `seata-k8s-controller-manager`.

### Example

If the CRD is deployed successfully, a message will display: `create seata crd success`, indicating that a CRD named `seataservers.operator.seata.apache.org` has been deployed in the specified namespace.

If the Deployment is successful, a message will display: `Deployment created successfully`, indicating that a Deployment named `seata-k8s-controller-manager` has been created in the specified namespace, containing a single Pod.

```
install --namespace=default --image=apache/seata-controller:latest --name=seata-k8s-controller-manager

INFO[0007] Create seata crd success                     
INFO[0007] Deployment created successfully  
```

If the CRD already exists, a message will display: `seata crd already exists`.

If the Deployment already exists, a message will display: `Deployment 'seata-k8s-controller-manager' already exists in the 'default' namespace`.

```
install --namespace=default --image=apache/seata-controller:latest --name=seata-k8s-controller-manager

ERRO[0001] install CRD err: failed to send post request: seata crd already exists 
ERRO[0001] install Controller err: Deployment 'seata-k8s-controller-manager' already exists in the 'default' namespace 
```

## Uninstall

### Describe

Use this command to uninstall the CRD and the controller.

### Usage

```
seata-ctl
uninstall --namespace=default --name=seata-k8s-controller-manager
```

#### Arguments

- `namespace`: Specifies the namespace to use (optional). If not specified, it defaults to `default`.
- `name`: Specifies the name of the deployed controller (optional). If not specified, it defaults to `seata-k8s-controller-manager`.

### Example

If the CRD is deleted successfully, a message will display: `CRD seataservers.operator.seata.apache.org deleted successfully`. If the Deployment is deleted successfully, a message will display: `Deployment 'seata-k8s-controller-manager' deleted successfully from namespace 'default'`.

```
uninstall --namespace=default --name=seata-k8s-controller-manager

INFO[0017] delete CRD seataservers.operator.seata.apache.org successfully. 
INFO[0017] deleted Controller seata-k8s-controller-manager successfully  
```

If the CRD does not exist, a message will display: `CRD seataservers.operator.seata.apache.org does not exist, no action taken`. If the Deployment does not exist, a message will display: `Deployment 'seata-k8s-controller-manager' does not exist in namespace 'default', no action taken`.

```
uninstall --namespace=default --name=seata-k8s-controller-manager

ERRO[0005] CRD seataservers.operator.seata.apache.org does not exist, no action taken. 
ERRO[0005] Deployment 'seata-k8s-controller-manager' does not exist in namespace 'default', no action taken. 
```

## Deploy

### Describe

Use this command to deploy the Seata server on the Kubernetes cluster.

### Usage

```
seata-ctl
deploy --name=test --replicas=1 --namespace=default --image=apache/seata-server:latest
```

After the CR is successfully deployed, the controller will automatically deploy the Seata server instance.

> **Note**: Since the specific namespace where the CRD and controller are installed is uncertain, the Deployment command will not check whether the CRD and controller have been successfully deployed. Before using this command, please ensure that the Seata CRD and controller have been successfully deployed in your cluster.

#### Arguments

- `namespace`: Specifies the namespace. If not specified, it defaults to `default`.
- `replicas`: Number of deployment replicas. If not specified, it defaults to `1`.
- `name`: Name of the CR. If not specified, it defaults to `example-seataserver`.
- `image`: Specifies the name of the image to use (optional). If not specified, it defaults to `apache/seata-server:latest`.

### Example

If the deployment is successful, a message will display: `CR install success, name: example-seataserver`.

```
deploy --name=test --replicas=1 --namespace=default --image=apache/seata-server:latest

INFO[0013] CR install success，name: test
```

If the deployment already exists, a message will display: `This seata server already exists!`.

```
deploy --name=test --replicas=1 --namespace=default --image=apache/seata-server:latest

ERRO[0007] deploy err:seata server already exist! name:test 
```

If the deployment fails, a message displaying the relevant error information will appear.

```
deploy --name=test --replicas=1 --namespace=default --image=apache/seata-server:latest

ERRO[0065] deploy err:the server could not find the requested resource 
```

If this error occurs, it indicates that the CRD and controller were not successfully deployed. Please check your deployment status.

## UnDeploy

### Describe

Use this command to uninstall the deployed Seata server.

### Usage

```
seata-ctl
undeploy --name=test --namespace=default
```

After deletion, the controller will automatically remove the Seata server instance.

#### Arguments

- `namespace`: Specifies the namespace. If not specified, it defaults to `default`.
- `name`: Name of the CR. If not specified, it defaults to `example-seataserver`.

### Example

If the uninstallation is successful, a message will display: `CR delete success, name: example-seataserver`, indicating that the CR has been successfully uninstalled, and the controller will automatically delete the Seata server from the cluster.

```
undeploy --name=test --namespace=default

INFO[0308] CR delete success，name: test 
```

If the CR is not found, a message will display: `seataservers.operator.seata.apache.org "example-seataserver" not found`.

```
undeploy --name=test --namespace=default

ERRO[0290] undeploy error: seataservers.operator.seata.apache.org "test" not found 
```

If an exception occurs, the relevant error information will be displayed.

```
undeploy --name=test --namespace=default

ERRO[0284] undeploy error: the server could not find the requested resource 
```

If this error occurs, it indicates that the CRD and controller were not successfully deployed. Please check your deployment status.

## Scale

### Describe

Use this command to modify the number of replicas for the deployed Seata server.

### Usage

```
seata-ctl

scale --name=test --replicas=2 --namespace=default
```

After the operation, the controller will automatically adjust the replica count, which may take some time.

#### Arguments

- `namespace`: Specifies the namespace. If not specified, it defaults to `default`.
- `name`: Name of the CR. If not specified, it defaults to `example-seataserver`.
- `replicas`: Number of deployment replicas. If not specified, it defaults to `1`.

### Example

If the modification is successful, a message will display: `CR scale success, name: example-seataserver`, indicating that the CR has been successfully updated, and the controller will automatically adjust the number of Seata server replicas in the cluster.

```
scale --name=test --replicas=2 --namespace=default

INFO[0070] CR scale success，name: test  
```

If the CR is not found, a message will display: `This seata server does not exist!`.

```
scale --name=test --replicas=2 --namespace=default

ERRO[0058] scale err:This seata server does not exits！test 
```

## Status

### Describe

Use this command to view the status of the deployed Seata server on the cluster.

### Usage

```
seata-ctl
status --name=example-seataserver --namespace=default
```

> **Note**: This command essentially filters and displays Pods based on their labels. If Seata server Pods were created using other methods (for example, outside of the controller), the relevant Pods may not be found.

#### Arguments

`namespace`: Specifies the namespace. If not specified, it defaults to `default`.

`name`: The name of the CR created when deploying the Seata server. If not specified, it defaults to `example-seataserver`.

### Example

If everything is functioning correctly, a list of `seata-pods` will be displayed.

```
status --name=example-seataserver --namespace=default

INFO[0319] status: POD NAME                  STATUS     
INFO[0319] status: ------------------------------------------- 
INFO[0319] status: example-seataserver-0     Running    
```

If no pods are found, an error message will be displayed.

```
status --name=example-seataserver1 --namespace=default

ERRO[0002] get k8s status error: no matching pods found 
```

# Prometheus

Use this command to display data retrieved from Prometheus.

Before using this command, you need to ensure that Prometheus is properly set up and is correctly collecting metrics from Seata.

## Metrics

### Describe

Use this command to show the data retrieved from Prometheus.

### Usage

```
seata-ctl
metrics --target=seata_transaction_summary  
```

#### Arguments

- `target`: Specify the metric entry to query. This may vary depending on the system and version, so make sure it aligns with the metrics in Prometheus.

### Example

If everything is functioning correctly, the relevant Prometheus metrics will be displayed in the form of a line chart.

```
metrics --target=seata_transaction_summary

INFO[0095]  10.00 ┼─────────╮
  9.01 ┤         │
  8.02 ┤         ╰╮
  7.03 ┤          │
  6.04 ┤          ╰╮
  5.05 ┤           │
  4.06 ┤           ╰╮
  3.07 ┤            ╰╮
  2.08 ┤             ╰─────────╮
  1.09 ┤                       ╰───────────╮
  0.10 ┤                                   ╰─────────────
                    seata_transaction_summary 
```

If an error occurs, an error message will be displayed.

```
metrics --target=seata_transaction_summary

ERRO[0093] Failed to show metrics: no data found for metric: seata_transaction_summary 
```

This error may indicate that the metric was not found. Please check your Prometheus metrics.

# Log

Use this command to view logs on Kubernetes, currently supporting ElasticSearch, Loki, and a custom-developed logging platform. Configuration files and query parameters vary between platforms.

## Log

### Describe

Use this command to retrieve Seata logs from different logging platforms.

Before using this command, please ensure that your logging system is fully deployed and functioning properly.

### Usage

#### ElasticSearch

```
seata-ctl

log --label={stream=stderr}  --number=3
```

Before using ES, you need to make the following configurations in the global configuration file. The meanings of the related parameters are as follows,An example of the configuration is as follows:

```
    - name: "es"  //name
      types: "ElasticSearch"  //es type
      address: "https://localhost:9200"   //es address
      source: "logstash-2024.10.24"    //es index name
      username: "elastic"   // es auth username
      password: "bu4AC50REtt_7rUqddMe"  //es auth password
      index: "log"	//doc type
```

> Note: Currently, only login with a username and password is supported for ES. If your ES has TLS encryption enabled, it is not supported at this time.

Note: Currently, only login with a username and password is supported for ES. If your ES has TLS encryption enabled, it is not supported at this time.

n the chart, the `index` field represents the output field for the final log display, which may vary across different index structures in ES.The source contains our index structure and the returned document content, but most of the fields are unnecessary for us. We use the index entry to select among these, and only the matching logs will be outputted.

For example, if the returned structure of our query is this:

```
{
  "took": 2,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 38,
      "relation": "eq"
    },
    "max_score": 1.0,
    "hits": [
      {
        "_index": "logstash-2024.10.24",
        "_id": "KyvwyJIB5EKHPP6lI9Os",
        "_score": 1.0,
        "_ignored": ["log.keyword"],
        "_source": {
          "@timestamp": "2024-10-24T17:39:45.771Z",
          "log": "2024-10-24T17:39:45Z\tINFO\tCreating a new SeataServer Service {default:seata-server-cluster}\t{\"controller\": \"seataserver\", \"controllerGroup\": \"operator.seata.apache.org\", \"controllerKind\": \"SeataServer\", \"SeataServer\": {\"name\":\"example-seataserver\",\"namespace\":\"default\"}, \"namespace\": \"default\", \"name\": \"example-seataserver\", \"reconcileID\": \"5912d1c4-6afd-4d66-8f11-eea6c54a1e4c\"}\n",
          "stream": "stderr",
          "time": "2024-10-24T17:39:45.771071Z",
          "kubernetes": {
            "pod_name": "seata-k8s-controller-manager-6448b86796-6l46n",
            "namespace_name": "default",
            "pod_id": "966edd0e-bc15-4276-9e6a-e255e84859ce",
            "labels": {
              "app": "seata-k8s-controller-manager",
              "pod-template-hash": "6448b86796"
            },
            "host": "minikube",
            "container_name": "seata-k8s-controller-manager",
            "docker_id": "b87df4627de30b0aab402aef4621fa46f45b3a1976dd6f466ab4a42bb8265455",
            "container_hash": "bearslyricattack/seata-controller@sha256:84ddedd9fa63c4a3ab8aa3a5017065014826dc3a9cac8c8fff328fbf9c600f11",
            "container_image": "bearslyricattack/seata-controller:latest"
          }
        }
      }
    ]
  }
}

```

We only need the logs related to the `log` field, so you can set `index` to `log`.

##### Arguments

- `label`: Specify the query entries in the format `{key1=value1, key2=value2, ...}`, where `key` is the name of the index in ES and `value` is the index's value. Neither the key nor the value requires quotation marks in the query.The default value is `{}`.
- `number`: The number of log entries returned, with a default value of 10.

> Note: To accommodate different ES index structures, a preliminary index structure check will be conducted during log queries. If the index does not contain this structure, the query will not proceed and an error will be returned.

#### Loki

```
seata-ctl

log --label={job="seata-transaction"} --start=2025-10-18-14:21:21 --end=2025-10-18-22:16:14 --number=3
```

Before using Loki, you also need to configure it in the configuration file.

```
    - name: "loki"  //loki name
      types: "Loki"	//type name
      address: "http://localhost:3100"  //loki address
      source: ""
      username: ""
      password: ""
```

##### Arguments

- **query**: The label for the query, following the native format, e.g., `{job="seata-transaction"}`, depending on the overall log format created.
- **start**: Start time in UTC timezone, formatted as `2024-10-18-12:32:54`.
- **end**: End time in UTC timezone, formatted as `2024-10-18-12:32:54`.
- **number**: The number of log entries returned.

#### Local

```
seata-ctl

log --level=INFO --number=5
```

Similarly, to meet more diverse needs, we have developed a custom logging system for Seata applications on Kubernetes. This system will continue to evolve along with the Seata-K8s project.

At first,configure it in the configuration file.

```
    - name: "local" //name
      types: "Local"  //type
      address: "http://localhost:8080"  //address
      source: "seata" //application
      username: ""
      password: ""
      index: ""
```

### Example

#### ElasticSearch

If everything is functioning correctly, the specified log entries will be returned.

```
log --label={stream=stderr}  --number=3

INFO[0015] 2024-10-24T17:39:45Z INFO    Creating a new SeataServer Service {default:seata-server-cluster}       {"controller": "seataserver", "controllerGroup": "operator.seata.apache.org", "controllerKind": "SeataServer", "SeataServer": {"name":"example-seataserver","namespace":"default"}, "namespace": "default", "name": "example-seataserver", "reconcileID": "5912d1c4-6afd-4d66-8f11-eea6c54a1e4c"} 
INFO[0015] 2024-10-24T17:39:45Z INFO    Creating a new SeataServer StatefulSet {default:example-seataserver}    {"controller": "seataserver", "controllerGroup": "operator.seata.apache.org", "controllerKind": "SeataServer", "SeataServer": {"name":"example-seataserver","namespace":"default"}, "namespace": "default", "name": "example-seataserver", "reconcileID": "5912d1c4-6afd-4d66-8f11-eea6c54a1e4c"} 
INFO[0015] 2024-10-24T17:39:45Z INFO    SeataServer(default/example-seataserver) has not been synchronized yet, requeue in 10 seconds   {"controller": "seataserver", "controllerGroup": "operator.seata.apache.org", "controllerKind": "SeataServer", "SeataServer": {"name":"example-seataserver","namespace":"default"}, "namespace": "default", "name": "example-seataserver", "reconcileID": "5912d1c4-6afd-4d66-8f11-eea6c54a1e4c"} 
```

If the index does not exist, an error message will be displayed.

```
log --label={stream1111=stderr}  --number=3

ERRO[0001] get log error: invalid index key: stream1111 
```

If no documents are found, an error indicating "no results found" will be returned.

```
log --label={stream=stderr11}  --number=3

ERRO[0060] get log error: no documents found     
```

If any other errors occur, an error message will be displayed.

#### Loki

If everything is functioning correctly, the specified log entries will be returned.

```
log --label={job="seata-transaction"} --start=2024-10-18-14:21:21 --end=2024-10-18-22:16:14 --number=2

INFO[0098] 2023-10-17 14:15:23.363 INFO  [TransactionManager] --- Transaction [XID: 172.16.20.1:8091:123456789] ends with SUCCESS. 
INFO[0098] 2023-10-17 14:15:23.360 INFO  [RMHandler] --- Transaction resource committed successfully: ResourceId: jdbc:mysql://localhost:3306/seata_test, XID: 172.16.20.1:8091:123456789 
```

If no documents are found, an error indicating "no results found" will be returned.

```
log --label={job="seata-transaction"} --start=2025-10-18-14:21:21 --end=2025-10-18-22:16:14 --number=3

ERRO[0058] get log error: loki query returned no results 
```

If any other errors occur, an error message will be displayed.

#### Local

If everything is functioning correctly, the specified log entries will be returned.

```
log --level=INFO --number=5

Peer example-seataserver-2.seata-server-cluster:9091 is connected
Peer example-seataserver-1.seata-server-cluster:9091 is connected
Node <default/example-seataserver-0.seata-server-cluster:9091> term 4 start preVote.
Peer example-seataserver-2.seata-server-cluster:9091 is connected
Peer example-seataserver-1.seata-server-cluster:9091 is connected
```

If any other errors occur, an error message will be displayed.