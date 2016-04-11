# Bootkube

Bootkube provides an entire kubernetes control plane in a single binary, and includes commands to render and consume assets for bootstrapping a self-hosted kubernetes cluster. Bootkube is expected to be used simply for bootstrapping purposes.

## Usage

Bootkube has two modes of operation. 

### Render assets

First, you can use bootkube to render out all of the assets (including kubernetes object manifests, TLS assets and kubeconfig) that you need to run a self-hosted kubernetes cluster. This feature is still experimental and changing rapidly.

To use this feature, run:

```
bootkube render <options>
```

You can customize the generated manifests by passing flags to the command. For more information on the supported commands, run `bootkube help render`.

### Start bootkube

To start bootkube use the `start` subcommand:

```
bootkube start <options>
```

Bootkube expects a directory containing the manifests to be provided as a command line flag, as well as other TLS assets (all of which can be taken from the `render` command). To see the available flags, run `bootkube help start`. 

When you start bootkube, you must also give it the addresses of your etcd servers, and enough information for bootkube to create an ssh tunnel to the node that will become a member of the master control plane. Upon startup, bootkube will create a reverse proxy using an ssh connection, which will allow a bootstrap kubelet to contact the apiserver running as part of bootkube.

## Build

First, clone the repo into the proper location in your $GOPATH:

```
git clone git@github.com:coreos/bootkube.git && cd bootkube
```

Then, to build:

```
make all
```

And optionally, to install into $GOPATH/bin:

```
make install
```

## License

bootkube is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
