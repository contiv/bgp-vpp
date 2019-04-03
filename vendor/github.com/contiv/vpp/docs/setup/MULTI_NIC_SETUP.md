### Setting up a node with multiple NICs


#### Determining network adapter PCI addresses
First, you need to find out the PCI address of the network interface that
you want VPP to use for multi-node interconnect. On Debian-based
distributions, you can use `lshw`(*):

```
$ sudo lshw -class network -businfo
Bus info          Device      Class      Description
====================================================
pci@0000:00:03.0  ens3        network    Virtio network device
pci@0000:00:04.0  ens4        network    Virtio network device
```
\* On CentOS/RedHat/Fedora distributions, `lshw` may not be available by default, install it by
    ```
    yum -y install lshw
    ```

#### Creating VPP configuration
Next, configure hardware interfaces in the VPP startup config, as
described [here](VPP_CONFIG.md#multi-nic-configuration).


#### Configuring network for node interconnect
Next, you need to tell Contiv which IP subnet / IP addresses should be used for interconnecting
VPP interfaces of the nodes in the cluster.

This can be either done by modification of the deployment YAML file (we use this approach in this guide), 
or using [helm and corresponding helm options](../../k8s/contiv-vpp/README.md).


##### Global configuration:
Global configuration is used in homogeneous environments where all nodes in
a given cluster have the same hardware configuration. 

If you want to use DHCP to configure node IPs, put the following stanza into
the [`contiv.conf`](../../k8s/contiv-vpp.yaml) deployment file:
```
data:
  contiv.conf: |-
    ...
    ipamConfig:
       nodeInterconnectDHCP: true
    ...
```

Alternatively to DHCP, you can also use automatic IP assignment from the `nodeInterconnectCIDR` pool 
(pre-defined and customizable in [`contiv.conf`](../../k8s/contiv-vpp.yaml) deployment file).
Also configure the `defaultGateway` if the PODs should be able to reach the internet:
```
data:
  contiv.conf: |-
    ...
    ipamConfig:
       nodeInterconnectCIDR: 192.168.16.0/24
       defaultGateway: 192.168.16.100
    ...
```


##### Individual configuration:
Individual configuration is used in heterogeneous environments where each node
in a given cluster may be configured differently. The configuration of nodes can
be defined in two ways:

- **contiv deployment file** This approach expects configuration for all cluster nodes
to be part of contiv deployment yaml file. Put the following stanza for each node into
the [`contiv.conf`][1] deployment file, for example:
```
...
    nodeConfig:
    - nodeName: "k8s-master"
      mainVPPInterface:
        interfaceName: "GigabitEthernet0/8/0"
        ip: "192.168.16.1/24"
      gateway: "192.168.16.100"
    - nodeName: "k8s-worker1"
      mainVPPInterface:
        interfaceName: "GigabitEthernet0/8/0"
        ip: "192.168.16.2/24"
      gateway: "192.168.16.100"
...
``` 


- **CRD** The node configuration is applied using custom defined k8s resource. The same behavior
as above can be achieved by the following CRD config. Usage of CRD allows to dynamically 
add new configuration for nodes joining a cluster.

```
$ cat master.yaml

# Configuration for node in the cluster
apiVersion: nodeconfig.contiv.vpp/v1
kind: NodeConfig
metadata:
  name: k8s-master
spec:
  mainVPPInterface:
     interfaceName: "GigabitEthernet0/8/0"
     ip: "192.168.16.1/24"
  gateway: "192.168.16.100"

$ cat worker1.yaml

# Configuration for node in the cluster
apiVersion: nodeconfig.contiv.vpp/v1
kind: NodeConfig
metadata:
  name: k8s-worker1
spec:
  mainVPPInterface:
     interfaceName: "GigabitEthernet0/8/0"
     ip: "192.168.16.2/24"
  gateway: "192.168.16.100"


kubectl apply -f master.yaml
kubectl apply -f worker1.yaml
```

These two ways of configuration are mutually exclusive. You select the one you want to use by
`crdNodeConfigurationDisabled` option in `contiv.conf`.
