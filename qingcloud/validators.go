package qingcloud

import (
	"fmt"
	"net"
	"strings"
)

func validateRouterVxnetsCIDR(v interface{}, k string) (ws []string, errors []error) {
	vxnets := v.(map[string]interface{})
	for vxnet, IPNetwork := range vxnets {
		if strings.HasPrefix(vxnet, "vxnet-") && vxnet != "vxnet-0" {
			_, ipnet, err = net.ParseCIDR(IPNetwork.(string))
			if err != nil {
				errors = append(errors, fmt.Errorf(
					"%q:%q must contain a valid CIDR, got error parsing: %s", vxnet, IPNetwork, err))
				return
			}
		} else {
			errors = append(errors, fmt.Errorf(
				"%q:%q must contain a valid vxnet id", vxnet, IPNetwork))
			return
		}
	}
	return
}

func validateVxnetsIPNetworkCIDR(v interface{}, k string) (ws []string, errors []error) {
	vaule := v.(string)
	_, ipnet, err := net.ParseCIDR(vaule)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must a valid CIDR, got error parsing: %s", value, err))
		return
	}
	return
}
