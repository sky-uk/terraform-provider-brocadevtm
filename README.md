****Terraform Provider
==================

Terraform provider for BrocadeVTM appliance


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/sky-uk/terraform-provider-brocadevtm`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:sky-uk/terraform-provider-brocadevtm
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sky-uk/terraform-provider-brocadevtm
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-brocadevtm
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```



Example Templates
------------------

To help people understand how resources can be created we have put together a collection of examples that will allow you to do so . 

Pool Resource 
--------------



```
resource "brocadevtm_pool" "pool_demo" {
       name = "pool_demo"
       monitorlist = ["ping"]
       node {
             node="127.0.0.1:80"
             priority=1
             state="active"
             weight=1
      }
      node {
            node="127.0.0.1:81"
            priority=1
            state="active"
            weight=1
     }
     node {
           node="127.0.0.1:82"
           priority=1
           state="active"
           weight=1
    }
    node {
          node="127.0.0.1:83"
          priority=1
          state="active"
          weight=1
   }
      max_connection_attempts = 5
}

```
