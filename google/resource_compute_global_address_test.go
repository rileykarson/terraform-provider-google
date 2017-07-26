package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func TestAccComputeGlobalAddress_basic(t *testing.T) {
	var addr compute.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeGlobalAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeGlobalAddressExists(
						"google_compute_global_address.foobar", &addr),
				),
			},
		},
	})
}

func TestAccComputeGlobalAddress_ipv6(t *testing.T) {
	var addr computeBeta.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeGlobalAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeGlobalAddress_ipv6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeBetaGlobalAddressExists(
						"google_compute_global_address.foobar", &addr),
				),
			},
		},
	})

	if addr.IpVersion != "IPV6" {
		t.Errorf("Expected IP version to be IPV6, got %s", addr.IpVersion)
	}
}

func testAccCheckComputeGlobalAddressDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_global_address" {
			continue
		}

		_, err := config.clientCompute.GlobalAddresses.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Address still exists")
		}
	}

	return nil
}

func testAccCheckComputeGlobalAddressExists(n string, addr *compute.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.GlobalAddresses.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Addr not found")
		}

		*addr = *found

		return nil
	}
}

func testAccCheckComputeBetaGlobalAddressExists(n string, addr *computeBeta.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientComputeBeta.GlobalAddresses.Get(config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Addr not found")
		}

		*addr = *found

		return nil
	}
}

var testAccComputeGlobalAddress_basic = fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
	name = "address-test-%s"
}`, acctest.RandString(10))

var testAccComputeGlobalAddress_ipv6 = fmt.Sprintf(`
resource "google_compute_global_address" "foobar" {
	name = "address-test-%s"
	ip_version = "IPV6"
}`, acctest.RandString(10))
