package uefi

import (
	"net"
	"sort"
	"strings"

	"github.com/EmilyShepherd/kios-go-sdk/pkg/bootstrap"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
	kubelet "k8s.io/kubelet/config/v1beta1"
)

type Provider struct {
}

func (p *Provider) Init() error {
	return EnsureEfiVarFs()
}

const CACert = "Kube-CA-Cert"
const CAKey = "Kube-CA-Key"

func (p *Provider) GetClusterCA() bootstrap.Cert {
	ca := bootstrap.Cert{}

	ca.Cert = GetValue(CACert)
	ca.Key = GetValue(CAKey)

	return ca
}

// TODO
func (p *Provider) GetClusterAuthInfo() kubeconfig.AuthInfo {
	return kubeconfig.AuthInfo{}
}

// TODO
func (p *Provider) GetCredentialProviders() []kubelet.CredentialProvider {
	return []kubelet.CredentialProvider{}
}

const Hostname = "Hostname"

func (p *Provider) GetHostname() string {
	return string(GetValue(Hostname))
}

// Reads the desired nameservers from the EFI variables
// EFI variables with the name NameserverX, where X is any string, will
// be read and interpreted acording, depending on the value length:
//   - 4 bytes = Binary representation of IPv4 address
//   - 16+ bytes - Binary representation of IPv6 address. Bytes on the end
//     are interpreted as a scope zone.
func (p *Provider) GetNameservers() (nameservers []string) {
	ns := GetMultiValue("Nameserver")
	nameservers = make([]string, len(ns))
	sort.Strings(nameservers)

	i := 0
	for _, value := range ns {
		if len(value) == net.IPv4len || len(value) == net.IPv6len {
			nameservers[i] = net.IP(value).String()
		} else if len(value) > net.IPv6len {
			nameservers[i] = net.IP(value[:net.IPv6len]).String()
			nameservers[i] += "%" + string(value[net.IPv6len:])
		} else {
			nameservers[i] = string(value)
		}

		i += 1
	}

	return
}

func (p *Provider) GetNodeLabels() (labels map[string]string) {
	l := GetMultiValue("Node-Label-")
	labels = make(map[string]string, len(l))

	for key, value := range l {
		key = strings.ReplaceAll(key, "_", "/")
		labels[key] = string(value)
	}

	return
}

const ApiEndpoint = "Kube-Api-Endpoint"

func (p *Provider) GetClusterEndpoint() string {
	return string(GetValue(ApiEndpoint))
}

func (p *Provider) GetContainerRuntimeConfiguration() bootstrap.ContainerRuntimeConfiguration {
	return bootstrap.ContainerRuntimeConfiguration{}
}

func (p *Provider) GetKubeletConfiguration(kubeConfig kubelet.KubeletConfiguration) kubelet.KubeletConfiguration {
	return kubeConfig
}
