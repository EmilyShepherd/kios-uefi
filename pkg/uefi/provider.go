package uefi

import (
	"os"
	"strings"
	"syscall"

	"github.com/EmilyShepherd/kios-go-sdk/pkg/bootstrap"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
	"k8s.io/klog/v2"
	kubelet "k8s.io/kubelet/config/v1beta1"
)

const KiosGuid = "3a3ae00f-cf4d-4977-8766-8ea575f6bf4f"
const CACert = "Kube-CA-Cert"
const CAKey = "Kube-CA-Key"
const Hostname = "Hostname"
const ApiEndpoint = "Kube-Api-Endpoint"

const EfiVarsDirectory = "/sys/firmware/efi/efivars"
const BootCurrent = "BootCurrent-8be4df61-93ca-11d2-aa0d-00e098032b8c"
const EfiVarFs = "efivarfs"

type Provider struct {
}

func (p *Provider) Init() error {
	if _, err := os.Stat(EfiVarsDirectory + "/" + BootCurrent); err == nil {
		klog.Info("Efivars is already mounted")
		return nil
	}

	syscall.Mount(EfiVarFs, EfiVarsDirectory, EfiVarFs, 0, "")

	return nil
}

func readVar(name string) []byte {
	data, err := os.ReadFile(EfiVarsDirectory + "/" + name + "-" + KiosGuid)
	if err != nil {
		return []byte{}
	}

	// Skip the first 4 bytes as these are UEFI Variable Attributes
	return data[4:]
}

func (p *Provider) GetClusterCA() bootstrap.Cert {
	ca := bootstrap.Cert{}

	ca.Cert = readVar(CACert)
	ca.Key = readVar(CAKey)

	return ca
}

// TODO
func (p *Provider) GetClusterAuthInfo() kubeconfig.AuthInfo {
	return kubeconfig.AuthInfo{}
}

// TODO
func (p *Provider) GetCredentialProviders() []kubelet.CredentialProvider {
	return []kubelet.CredentialProvider{kubelet.CredentialProvider{}}
}

func (p *Provider) GetHostname() string {
	return string(readVar(Hostname))
}

func (p *Provider) GetNodeLabels() map[string]string {
	labels := make(map[string]string)

	prefix := "Node-Label-"
	guidLength := len(KiosGuid) + 1
	prefixLength := len(prefix)

	entries, _ := os.ReadDir(EfiVarsDirectory)
	for _, file := range entries {
		name := file.Name()
		contentLength := len(name) - guidLength

		if len(name) <= guidLength+prefixLength {
			continue
		}
		if name[:prefixLength] == prefix && name[contentLength:] == "-"+KiosGuid {
			name = name[:contentLength]
			key := strings.ReplaceAll(name[prefixLength:], "_", "/")
			labels[key] = string(readVar(name))
		}
	}

	return labels
}

func (p *Provider) GetClusterEndpoint() string {
	return string(readVar(ApiEndpoint))
}

func (p *Provider) GetContainerRuntimeConfiguration() bootstrap.ContainerRuntimeConfiguration {
	return bootstrap.ContainerRuntimeConfiguration{}
}

func (p *Provider) GetKubeletConfiguration(kubeConfig kubelet.KubeletConfiguration) kubelet.KubeletConfiguration {
	return kubeConfig
}
