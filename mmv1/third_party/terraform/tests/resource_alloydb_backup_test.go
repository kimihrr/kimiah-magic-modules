package google

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlloydbBackup_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  BootstrapSharedTestNetwork(t, "alloydb-update"),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbBackup_alloydbBackupFullExample(context),
			},
			{
				ResourceName:            "google_alloydb_backup.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_id", "location", "reconciling", "update_time"},
			},
			{
				Config: testAccAlloydbBackup_update(context),
			},
			{
				ResourceName:            "google_alloydb_backup.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_id", "location", "reconciling", "update_time"},
			},
		},
	})
}

// Updates "label" field from testAccAlloydbBackup_alloydbBackupFullExample
func testAccAlloydbBackup_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_alloydb_backup" "default" {
  location     = "us-central1"
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.default.name

  description = "example description"
  labels = {
    "label" = "updated_key"
    "label2" = "updated_key2"
  }
  depends_on = [google_alloydb_instance.default]
}

resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = data.google_compute_network.default.id
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = data.google_compute_network.default.id
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = data.google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

data "google_compute_network" "default" {
  name = "%{network_name}"
}
`, context)
}

// We expect an error when creating an on-demand backup without location.
// Location is a `required` field.
func TestAccAlloydbBackup_missingLocation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccAlloydbBackup_missingLocation(context),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testAccAlloydbBackup_missingLocation(context map[string]interface{}) string {
	return Nprintf(`
resource "google_alloydb_backup" "default" {
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  cluster_name = google_alloydb_cluster.default.name
  depends_on = [google_alloydb_instance.default]
}

resource "google_alloydb_cluster" "default" {
  location = "us-central1"
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}
  
data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}
  
resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

// Test to create on-demand backup with mandatory fields.
func TestAccAlloydbBackup_createBackupWithMandatoryFields(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbBackup_createBackupWithMandatoryFields(context),
			},
		},
	})
}

func testAccAlloydbBackup_createBackupWithMandatoryFields(context map[string]interface{}) string {
	return Nprintf(`
resource "google_alloydb_backup" "default" {
  backup_id    = "tf-test-alloydb-backup%{random_suffix}"
  location = "us-central1"
  cluster_name = google_alloydb_cluster.default.name
  depends_on = [google_alloydb_instance.default]
}

resource "google_alloydb_cluster" "default" {
  location = "us-central1"
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  network    = google_compute_network.default.id
}

data "google_project" "project" { }

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}

resource "google_alloydb_instance" "default" {
  cluster       = google_alloydb_cluster.default.name
  instance_id   = "tf-test-alloydb-instance%{random_suffix}"
  instance_type = "PRIMARY"

  depends_on = [google_service_networking_connection.vpc_connection]
}

resource "google_compute_global_address" "private_ip_alloc" {
  name          =  "tf-test-alloydb-cluster%{random_suffix}"
  address_type  = "INTERNAL"
  purpose       = "VPC_PEERING"
  prefix_length = 16
  network       = google_compute_network.default.id
  lifecycle {
	ignore_changes = [
		address,
		creation_timestamp,
		id,
		network,
		project,
		self_link
	]
  }
}

resource "google_service_networking_connection" "vpc_connection" {
  network                 = google_compute_network.default.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}
