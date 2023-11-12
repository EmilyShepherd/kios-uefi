package uefi

import (
	"os"
	"syscall"

	"k8s.io/klog/v2"
)

const KiosGuid = "3a3ae00f-cf4d-4977-8766-8ea575f6bf4f"

const DefaultEfiVarsDirectory = "/sys/firmware/efi/efivars"
const BootCurrent = "BootCurrent-8be4df61-93ca-11d2-aa0d-00e098032b8c"
const EfiVarFs = "efivarfs"

var EfiVarsDirectory string = GetEfiDirectory()

func GetEfiDirectory() string {
	if efiDir := os.Getenv("EFI_DIRECTORY"); efiDir != "" {
		return efiDir
	}

	return DefaultEfiVarsDirectory
}

func EnsureEfiVarFs() error {
	if _, err := os.Stat(EfiVarsDirectory + "/" + BootCurrent); err == nil {
		klog.Info("Efivars is already mounted")
		return nil
	}

	syscall.Mount(EfiVarFs, EfiVarsDirectory, EfiVarFs, 0, "")

	return nil
}

func GetValue(name string) []byte {
	data, err := os.ReadFile(EfiVarsDirectory + "/" + name + "-" + KiosGuid)
	if err != nil {
		return []byte{}
	}

	// Skip the first 4 bytes as these are UEFI Variable Attributes
	return data[4:]
}

func GetMultiValue(prefix string) map[string][]byte {
	values := make(map[string][]byte)

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
			key := name[prefixLength:]
			values[key] = GetValue(name)
		}
	}

	return values
}
