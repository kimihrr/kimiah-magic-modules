package google

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBigtableTable_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable(instanceName, tableName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_splitKeys(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_splitKeys(instanceName, tableName),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"split_keys"},
			},
		},
	})
}

func TestAccBigtableTable_family(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_family(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_deletion_protection_enabled(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with deletion protection equals to protected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED"),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete the table because of deletion protection equals to protected
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
			// adding a column family to the table
			{
				Config: testAccBigtableTable_family(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete column families in the table with deletion protection equals to protected
			{
				Config:      testAccBigtableTable(instanceName, tableName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"column_family"},
			},
			// changing deletion protection field to unprotected without changing the column families
			{
				Config: testAccBigtableTable_deletion_protection_family(instanceName, tableName, "UNPROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// destroying the table is possible when deletion protection is equals to unprotected
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_deletion_protection_disabled(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with deletion protection equals to unprotected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED"),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// adding a column family to the table
			{
				Config: testAccBigtableTable_family(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// removing the column family is possible because the deletion protection field is unprotected
			{
				Config: testAccBigtableTable(instanceName, tableName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// changing the deletion protection field to protected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED"),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// changing the deletion protection field to unprotected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED"),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// destroying the table is possible when deletion protection is equals to unprotected
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_familyMany(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_familyUpdate(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableTable_familyUpdate(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigtableTableDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_table" {
				continue
			}

			config := googleProviderConfig(t)
			c, err := config.BigTableClientFactory(config.userAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
			if err != nil {
				// The instance is already gone
				return nil
			}

			_, err = c.TableInfo(ctx, rs.Primary.Attributes["name"])
			if err == nil {
				return fmt.Errorf("Table still present. Found %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
			}

			c.Close()
		}

		return nil
	}
}

func testAccBigtableTable(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_splitKeys(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  split_keys    = ["a", "b", "c"]
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_family(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, family)
}

func testAccBigtableTable_deletion_protection(instanceName, tableName, deletionProtection string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name
  deletion_protection = "%s"
}
`, instanceName, instanceName, tableName, deletionProtection)
}

func testAccBigtableTable_familyMany(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s-first"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableTable_familyUpdate(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s-third"
  }

  column_family {
    family = "%s-fourth"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family, family)
}

func testAccBigtableTable_deletion_protection_family(instanceName, tableName, deletionProtection, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name
  deletion_protection = "%s"

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, deletionProtection, family)
}

func testAccBigtableTable_destroyTable(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}
`, instanceName, instanceName)
}
