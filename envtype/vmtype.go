//go:build windows

package envtype

func VirtualizationSystem() string {

	switch {
	case detectVMWware():
		return "VMware Inc."
	case detectVBox():
		return "Oracle VM VirtualBox"
	case detectKVM():
		return "KVM/QEMU"
	case detectXen():
		return "Xen Hypervisor"
	case detectVpc():
		return "Virtual PC"
	default:
		return "Non VM"
	}
}
