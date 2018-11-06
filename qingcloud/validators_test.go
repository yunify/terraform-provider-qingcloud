package qingcloud

import (
	"testing"
)

func TestValidateNetworkCIDR(t *testing.T) {
	validCIDRNetworkAddress := []string{"192.168.10.0/24", "0.0.0.0/0", "10.121.10.0/24"}
	for _, v := range validCIDRNetworkAddress {
		_, errors := validateNetworkCIDR(v, "cidr_network_address")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid cidr network address: %q", v, errors)
		}
	}
	invalidCIDRNetworkAddress := []string{"1.2.3.4", "0x38732/21"}
	for _, v := range invalidCIDRNetworkAddress {
		_, errors := validateNetworkCIDR(v, "cidr_network_address")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid cidr network address", v)
		}
	}
}
func TestWithinArrayString(t *testing.T) {
	exceptValues := []string{"a", "bc", "def", "test"}
	validValues := []string{"a", "bc", "def", "test"}
	for _, v := range validValues {
		_, errors := withinArrayString(exceptValues...)(v, "array string")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid value in %#v: %q", v, exceptValues, errors)
		}
	}
	invalidValues := []string{"lal", "hhh", "gg"}
	for _, v := range invalidValues {
		_, errors := withinArrayString(exceptValues...)(v, "array string")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid value", v)
		}
	}
}

func TestWithinArrayInt(t *testing.T) {
	exceptValues := []int{1, 2, 3, 4}
	validValues := []int{1, 2, 3, 4}
	for _, v := range validValues {
		_, errors := withinArrayInt(exceptValues...)(v, "array int")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid value in %#v: %q", v, exceptValues, errors)
		}
	}
	invalidValues := []int{6, 100, -1, 20}
	for _, v := range invalidValues {
		_, errors := withinArrayInt(exceptValues...)(v, "array int")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid value", v)
		}
	}
}
func TestWithinArrayIntRange(t *testing.T) {
	validIntegers := []int{-259, 0, 1, 5, 999}
	min := -259
	max := 999
	for _, v := range validIntegers {
		_, errors := withinArrayIntRange(min, max)(v, "int range")
		if len(errors) != 0 {
			t.Fatalf("%q should be an integer in range (%d, %d): %q", v, min, max, errors)
		}
	}
	invalidIntegers := []int{-260, -99999, 1000, 25678}
	for _, v := range invalidIntegers {
		_, errors := withinArrayIntRange(min, max)(v, "int range")
		if len(errors) == 0 {
			t.Fatalf("%q should be an integer outside range (%d, %d)", v, min, max)
		}
	}
}
func TestValidateColorString(t *testing.T) {
	validColor := []string{"#1f1f1F", "#AFAFAF", "#1AFFa1", "#222fff", "#F00"}
	for _, v := range validColor {
		_, errors := validateColorString(v, "color string")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid color string: %q", v, errors)
		}
	}
	invalidColor := []string{"123456", "#afafah", "#123abce", "aFaE3f",
		"F00", "#afaf", "#F0h"}
	for _, v := range invalidColor {
		_, errors := validateColorString(v, "color string")
		if len(errors) == 0 {
			t.Fatalf("%q should be a invalid color string: %q", v, errors)
		}
	}
}

func TestValidatePortString(t *testing.T) {
	validPort := []string{"65535", "0", "2333"}
	for _, v := range validPort {
		_, errors := validatePortString(v, "port string")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid port string: %q", v, errors)
		}
	}
	invalidPort := []string{"65536", "-1", "hhhh"}
	for _, v := range invalidPort {
		_, errors := validatePortString(v, "port string")
		if len(errors) == 0 {
			t.Fatalf("%q should be a invalid port string: %q", v, errors)
		}
	}

}

func TestValidateVolumeSize(t *testing.T) {
	validSize := []int{10, 500, 1000}
	for _, v := range validSize {
		_, errors := validateVolumeSize(v, "volume size")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid volume size: %q", v, errors)
		}
	}
	invalidSize := []int{9, 5001, 5500}
	for _, v := range invalidSize {
		_, errors := validateVolumeSize(v, "volume size")
		if len(errors) == 0 {
			t.Fatalf("%q should be a invalid volume size: %q", v, errors)
		}
	}
}
