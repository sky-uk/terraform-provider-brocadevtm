# go-brocadev-tm - Go Bindings for the Brocade Virtual Traffic Manager (vTM)

## Overview

This is the GoLang API wrapper for Brocade Virtual Traffic Manager (vTM).
This wrapper uses the REST API interface provided by the vTM, currently version 3.8.


## Handled Resources

| Resource                | Create | Read  | Update  | Delete |
|-------------------------|--------|-------|---------|--------|
| Monitor [1]             |   Y    |   Y   |    Y    |   Y    |
| Pool                    |   Y    |   Y   |    Y    |   Y    |
| SSL Server Key          |   Y    |   Y   |    Y    |   Y    |
| Traffic IP Group        |   Y    |   Y   |    Y    |   Y    |
| Virtual Server          |   Y    |   Y   |    Y    |   Y    |
| Rule                    |   Y    |   Y   |    Y    |   Y    |

### Notes
[1] : Currently only HTTP monitoring is supported


## Usage

### Import library

```
import(
    "github.com/sky-uk/go-brocade-vtm"
)
```

### Get a client object

In order to get a client object authentication credentials needs to be
provided. The REST API transport procotol is by default HTTPS and hence
SSL encryption is used by default. Certificate handling behaviour can be 
controlled via the ignoreSSL flag.


```
    client := brocadevtm.NewVTMClient(
        serverUrl,  // https://<server>:<port>
        username,
        password,
        ignoreSSL,  // if server should ignore SSL certificate checking
        debug,      // if enhanced trace log is needed...
        header,     // a set of header key/value pairs to use, not mandatory
    )
```

### Working with Monitors

#### Importing proper packages

In order to work with monitors, import the monitor package section:

```
import (
    "github.com/sky-uk/go-brocade-vtm/api/monitor"
)
```


#### Getting the list of active monitors

```
    api := monitor.NewGetAll()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    monitors := api.GetResponse()
```


#### Retrieving a single monitor

```
    api := monitor.NewGetSingleMonitor( name )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    monitor := api.GetResponse()
```

#### Creating a monitor


```
    newMonitor := monitor.Monitor{
        ...
    }

    api := monitor.NewCreate( name, newMonitor )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    monCreated := api.GetResponse()
```

#### Updating a monitor


```
    updateMonitor := monitor.Monitor{
        ...
    }

    api := monitor.NewUpdate( name, updateMonitor )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    monUpdated := api.GetResponse()
```

#### Deleting a monitor

```
    api := monitor.NewDelete( name )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Monitor %s deleted", name)
    }
```

### Working with Pools

#### Importing proper packages

In order to work with pools, import the pool package section:

```
import (
    "github.com/sky-uk/go-brocade-vtm/api/pool"
)
```


#### Getting the list of pools

```
    api := monitor.NewGetAll()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    // returns a pool.LBPoolList object
    pools := api.GetResponse()
```


#### Retrieving a single pool

```
    api := monitor.NewGetSingle(name)
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    pool := api.GetResponse()
```

#### Creating a pool


```
    newPool := pool.Pool{
        ...
    }

    api := pool.NewCreate( name, newPool )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    poolCreated := api.GetResponse()
```

#### Updating a pool


```
    updatePool := pool.Pool{
        ...
    }

    api := pool.NewUpdate( name, updatePool )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    poolUpdated := api.GetResponse()
```

#### Deleting a pool

```
    api := pool.NewDelete( name )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Pool %s deleted", name)
    }
```

### Working with Server keys

#### Importing proper packages

In order to work with server keys, import the pool package section:

```
import (
    "github.com/sky-uk/go-brocade-vtm/api/sslServerKey"
)
```


#### Getting the list of server keys

```
    api := monitor.NewGetAll()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    
    // returns a sslServerKey.SSLServerKeysList object
    keys := api.GetResponse()
```


#### Retrieving a single server key

```
    api := sslServerKey.NewGetSingle(name)
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    key := api.GetResponse()
```

#### Creating a server key


```
    newKey := sslServerKey.SSLServerKey{
        ...
    }

    api := sslServerKey.NewCreate( name, newKey )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    keyCreated := api.GetResponse()
```

#### Updating a server key


```
    updateKey := sslServerKey.SSLServerKey{
        ...
    }

    api := sslServerKey.NewUpdate( name, updateKey )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    keyUpdated := api.GetResponse()
```

#### Deleting a server key

```
    api := sslServerKey.NewDelete( name )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Server key %s deleted", name)
    }
```

### Working with traffic IP groups

#### Importing proper packages

In order to work with traffic IP groups, import the traffic IP package section:

```
import (
    "github.com/sky-uk/go-brocade-vtm/api/trafficIpGroups"
)
```


#### Getting the list of traffic IP groups

```
    // to get the list of all traffic managers
    api := trafficIpGroups.NewGetAllTrafficManagers()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    // returns TrafficManagerChildren object
    managers := api.GetResponse()


    api := trafficIpGroups.NewGetAll()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    
    // returns a trafficIpGroups.TrafficManagerChildren object
    keys := api.GetResponse()
```


#### Retrieving a single traffic IP group

```
    api := trafficIpGroups.NewGetSingle(name)
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    tipg := api.GetResponse()
```

#### Creating a traffic IP group


```
    newTipg := trafficIpGroups.TrafficIPGroup{
        ...
    }

    api := trafficIpGroups.NewCreate( name, newTipg )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    tipgCreated := api.GetResponse()
```

#### Updating a traffic IP group


```
    updateTipg := trafficIpGroups.TrafficIPGroup{
        ...
    }

    api := trafficIpGroups.NewUpdate( name, updateTipg )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    tipgUpdated := api.GetResponse()
```

#### Deleting a traffic IP group

```
    api := trafficIpGroups.NewDelete( name )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Traffic IP group %s deleted", name)
    }
```

### Working with virtual servers

#### Importing proper packages

In order to work with virtual servers, import the virtual server package section:

```
import (
    "github.com/sky-uk/go-brocade-vtm/api/virtualserver"
)
```


#### Getting the list of virtual servers

```
    api := virtualserver.NewGetAll()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    
    // returns a virtualserver.VirtualServersList object
    keys := api.GetResponse()
```


#### Retrieving a single virtual server

```
    api := virtualserver.NewGetSingle(name)
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    vs := api.GetResponse()
```

#### Creating a virtual server


```
    newVS := virtualserver.VirtualServer{
        ...
    }

    api := virtualserver.NewCreate( name, newVS )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    vsCreated := api.GetResponse()
```

#### Updating a virtual server


```
    updateVS := virtualserver.VirtualServer{
        ...
    }

    api := virtualserver.NewUpdate( name, updateVS )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    vsUpdated := api.GetResponse()
```

#### Deleting a virtual server

```
    api := virtualserver.NewDelete( name )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Virtual server %s deleted", name)
    }
```

### Working with rules

#### Importing proper packages
     
     In order to work with rules, import the rule package section:
     
     ```
     import (
         "github.com/sky-uk/go-brocade-vtm/api/rule"
     )
     ```
     
#### Getting the list of rules

```
    api := rule.NewGetAll()
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    
    // returns a rule.Rules object which contains the name and href of each rule.
    rules := api.GetResponse()
```


#### Retrieving a single rule

```
    api := rule.NewGetRule(name)
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

    // Returns the traffic script as a string
    trafficScript := api.GetResponse()
```

#### Creating a rule

```
    ruleName := "my-rule"
    trafficScript := []byte(`
            if( string.ipmaskmatch( request.getremoteip(), "192.168.123.10" ) ){
                    connection.discard();
            }`)
    
    api := rule.NewCreate( ruleName, trafficScript )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }
    
```

#### Updating a rule

```
    ruleName := "my-rule"
    updatedTrafficScript := []byte(`
                                     if( string.ipmaskmatch( request.getremoteip(), "192.168.123.10" ) ){
                                             connection.discard();
                                     }`)
                                     
    api := virtualserver.NewUpdate( ruleName, updatedTrafficScript )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    }

```

#### Deleting a rule

```
    nameRuleToDelete := "my-rule"

    api := rule.NewDelete( nameRuleToDelete )
    err := client.Do(api)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Rule %s deleted", nameRuleToDelete)
    }
```     

### Using the CLI

Run ```make``` in the root of the repository to build the ```go-brocade-vtm-cli``` binary. Once built the binary may be found in the root of the repository.

Currently only rules may be managed by the cli.


#### Authentication

Set the following shell environment variables to allow the cli to make changes to a Brocade vTM. Note that the cli automatically prefixes the BROCADEVTM_SERVER with https://.

```
$ export BROCADEVTM_USERNAME=admin
$ export BROCADEVTM_PASSWORD=myPassword
$ export BROCADEVTM_SERVER=my-brocade-vtm.example.com:9070
```

#### Getting help

```
$ ./go-brocade-vtm-cli -h
```

#### Working with rules

##### Create a rule

To create a rule ensure the traffic script is in the file /path/to/traffic_script.

```
$ ./go-brocade-vtm-cli rule-create -name my-rule -script /path/to/traffic_script
```

##### Showing rules

To list all rules.
```
$ ./go-brocade-vtm-cli rule-show-all
```

To show a rule called 'my-rule'.
```
$ ./go-brocade-vtm-cli rule-show -name my-rule
```

##### Updating a rule

To update a rule called 'my-rule' with the traffic script (file) /path/to/updated_traffic_script.
```
$ ./go-brocade-vtm-cli rule-update -name my-rule -script /path/to/updated_traffic_script
```

##### Deleting a rule

To delete a rule
```
$ ./go-brocade-vtm-cli rule-delete -name my-rule
```