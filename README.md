Terraform Provider
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
$ mkdir -p $GOPATH/src/github.com/sky-uk/; cd $GOPATH/src/github.com/sky-uk/
$ git clone https://github.com/sky-uk/terraform-provider-brocadevtm
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sky-uk/terraform-provider-brocadevtm
$ make build
```

Using the provider
----------------------

See the [BrocadeVTM Provider wiki](http://github.com/sky-uk/terraform-provider-brocadevtm/wiki) to get started using the BrocadeVTM provider.

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

