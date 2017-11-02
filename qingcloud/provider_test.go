package qingcloud

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"qingcloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("QINGCLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("QINGCLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("QINGCLOUD_SECRET_KEY"); v == "" {
		t.Fatal("QINGCLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("QINGCLOUD_ZONE"); v == "" {
		log.Println("[INFO] Test: Using pek3a as test region")
		os.Setenv("QINGCLOUD_ZONE", DEFAULT_ZONE)
	}
}
