package main

import (
	"github.com/EmilyShepherd/kios-go-sdk/pkg/bootstrap"
	"github.com/EmilyShepherd/kios-uefi/pkg/uefi"
)

var Bootstrap = bootstrap.Bootstrap{
	Provider: &uefi.Provider{},
}

func main() {
	Bootstrap.Run()
}
