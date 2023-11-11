# kiOS UEFI Build

The bootstrap container and sample pod definition to run kiOS in generic
UEFI environments.

## Configuration

`uefi-bootstrap` auto configures kiOS by reading EFI Variables.

kiOS' vendor GUID is `3a3ae00f-cf4d-4977-8766-8ea575f6bf4f` and the
following can be configured:

### Hostname

The hostname of the instance can be set using the `Hostname` variable.
This is saved to `/etc/hostname`, which is used by kiOS to set the
instance hostname.

### Api Endpoint

The `Kube-Api-Endpoint` variable specifies the URL of the api-server for
the cluster. If this is specified, Webhook Authorization and
Authentication will be turned on by default. If this is not specified,
these will be disabled by default, and no KubeConfig will be created for
the kubelet, causing it to run in standalone mode.

### Cluster CA

The `Kube-CA-Cert` and `Kube-CA-Key` variables may contain the public /
private components of the cluster's CA and will be saved to
`/etc/kubernetes/pki/ca.crt` and `/etc/kubernetes/pki/ca.key`
respectively. If an api endpoint is provided, the public certificate is
used for the server ca in the kubelet's KubeConfig file, and as the
trusted ca for kubelet authentication.

### Nameservers

You may specify any number[^1] of nameservers to save into
`/etc/resolv.conf`. All variables starting with `Nameserver` are read in
string sorted order.

[^1]: You may specify as many nameservers as desired, however
kubernetes, will only ever respect the first three.

The format of this variable is normally interpreted as a binary
representation of an IP address. If it is 4 bytes, it is assumed to be
IPv4. If it is 16 or more, it is assumed to be IPv6, with the trailing
bytes interpreted as characters identifying the scope for the address.

#### Examples

`\xC0\x00\x02\x10` results in `192.0.2.16`

`\x20\x01\xDB\x08\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x23` - results in `2001:db8::1123`

`\x20\x01\xDB\x08\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x23\x65\x74\x68\x30` - results in `2001:db8::1123%eth0`

### Node Labels

Labels can be set for the node by setting an EFI variable called
`Node-Label-[label-key]`. Any underscores in the `[label-key]` are
interpreted as forward slashes.

For example, in order to set `example.com/label=foobar` we set the value
of `foobar` to the EFI Variable named `Node-Label-example.com_label`.
