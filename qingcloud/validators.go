package qingcloud

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

var ColorRegex = regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")
var PortRegex = regexp.MustCompile("^0*(?:6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{1,3}|[0-9])$")

func validateRouterVxnetsCIDR(v interface{}, k string) (ws []string, errors []error) {
	vxnets := v.(map[string]interface{})
	for vxnet, IPNetwork := range vxnets {
		if strings.HasPrefix(vxnet, "vxnet-") && vxnet != "vxnet-0" {
			if _, _, err := net.ParseCIDR(IPNetwork.(string)); err != nil {
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

func validateNetworkCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if _, _, err := net.ParseCIDR(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must a valid CIDR, got error parsing: %s", value, err))
		return
	}
	return
}

func withinArrayString(limits ...string) schema.SchemaValidateFunc {
	var limitsMap = make(map[string]bool)
	for _, v := range limits {
		limitsMap[v] = true
	}

	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if limitsMap[value] {
			return
		}
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match", k, value))
		return
	}
}

func withinArrayInt(limits ...int) schema.SchemaValidateFunc {
	var limitsMap = make(map[int]bool)
	for _, v := range limits {
		limitsMap[v] = true
	}

	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if limitsMap[value] {
			return
		}
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match", k, value))
		return
	}
}

func withinArrayIntRange(begin, end int) schema.SchemaValidateFunc {

	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value >= begin && value <= end {
			return
		}
		errors = append(errors, fmt.Errorf("%q (%q) should > %d  && < %d ", k, value, begin, end))
		return
	}
}
func validateColorString(v interface{}, k string) (ws []string, errors []error) {
	colorstring := v.(string)
	if !ColorRegex.MatchString(colorstring) {
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match", k, colorstring))
		return
	}
	return
}
func validatePortString(v interface{}, k string) (ws []string, errors []error) {
	portstring := v.(string)
	if !PortRegex.MatchString(portstring) {
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match", k, portstring))
		return
	}
	return
}
