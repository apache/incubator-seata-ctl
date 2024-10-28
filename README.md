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

# Deploy

## Describe

Use this command to deploy the Seata server on the Kubernetes cluster.

## Usage

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

## Example

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

# UnDeploy

## Describe

Use this command to uninstall the deployed Seata server.

## Usage

```
seata-ctl
undeploy --name=test --namespace=default
```

After deletion, the controller will automatically remove the Seata server instance.

#### Arguments

- `namespace`: Specifies the namespace. If not specified, it defaults to `default`.
- `name`: Name of the CR. If not specified, it defaults to `example-seataserver`.

## Example

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

# Scale

## Describe

Use this command to modify the number of replicas for the deployed Seata server.

## Usage

```
seata-ctl

scale --name=test --replicas=2 --namespace=default
```

After the operation, the controller will automatically adjust the replica count, which may take some time.

#### Arguments

- `namespace`: Specifies the namespace. If not specified, it defaults to `default`.
- `name`: Name of the CR. If not specified, it defaults to `example-seataserver`.
- `replicas`: Number of deployment replicas. If not specified, it defaults to `1`.

## Example

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

# Status

## Describe

Use this command to view the status of the deployed Seata server on the cluster.

## Usage

```
seata-ctl
status --name=example-seataserver --namespace=default
```

> **Note**: This command essentially filters and displays Pods based on their labels. If Seata server Pods were created using other methods (for example, outside of the controller), the relevant Pods may not be found.

#### Arguments

`namespace`: Specifies the namespace. If not specified, it defaults to `default`.

`name`: The name of the CR created when deploying the Seata server. If not specified, it defaults to `example-seataserver`.

## Example

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

