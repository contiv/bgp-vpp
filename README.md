# Referenced https://www.brianlinkletter.com/how-to-build-a-network-of-linux-routers-using-quagga/

# bgp-vpp
A BGP Speaker implementation for Contiv-VPP

# Quagga Setup
First, you need to enable the zebra feature in the Global configuration using:
//This will be done in the Gateway Configuration File

```
[zebra]
    [zebra.config]
        enabled = true
        url = "unix:/var/run/quagga/zserv.api"
        redistribute-route-type-list = ["connect"]
        version = 2
```
Second, start by entering the commands:
```
  $ sudo su //ensure that sudo is being used if you decide to not use this command to start
  # apt-get update
  # apt-get install quagga quagga-doc
```
Third, configure Quagga daemons by editing the file /etc/quagga/daemons to start the zebra daemon.
```
# nano /etc/quagga/daemons
```
File should appear like:
```
zebra=yes
bgpd=no
ospfd=no
ospf6d=no
ripd=no
ripngd=no
isisd=no
babeld=no
Save the file and quit the editor.
```
Fourth, create config files for the zebra daemon.
```
# cp /usr/share/doc/quagga/examples/zebra.conf.sample /etc/quagga/zebra.conf
# chown quagga.quaggavty /etc/quagga/*.conf
# chmod 640 /etc/quagga/*.conf
```
Fifth, start Quagga:
```
# /etc/init.d/quagga start
```
Now you will need to run the gateway and any worker bgps by using the command:
```
$ sudo -E ./gobgpd -f ./gobgpd.conf
```
You will also need to run the executable GoBGP that you create in GoLand to run bgp on the master. In GoLand, for example:
```
$ go run main.go
```
If you have worker nodes running, make sure to add the subnet IP and nexthop IP. For example:
```
$ ./gobgp global rib add -a ipv4 10.1.2.0/24 nexthop 192.168.16.2
```
Finally, check if your network is set up by viewing the routing table in the gateway using:
```
$ netstat -rn
```
You will also be able to ping the connection.
