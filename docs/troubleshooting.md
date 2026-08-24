<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Troubleshooting

`seata-ctl` provides read-only diagnostics for Seata Server.

## Diagnose

Run a quick health sweep:

```bash
seata-ctl
login --ip 127.0.0.1 --port 7091 --username seata --password seata
diagnose run
```

The command checks:

- server configuration
- TCP connectivity
- login token
- status endpoint
- global transaction query
- global lock query
- optional database connectivity

Use `--output table|json|yaml` to switch formats.

## Common failures

- `server address is not configured`: log in again with the right IP and port.
- `please login`: log in through the REPL first.
- `tcp connectivity failed`: check the Seata Server port and network path.
