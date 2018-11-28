package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanKubernetesNodePool_Basic(t *testing.T) {
	rName := acctest.RandString(10)
	var k8s godo.KubernetesCluster
	var k8sPool godo.KubernetesNodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDigitalOceanKubernetesClusterExists("digitalocean_kubernetes_cluster.foobar", &k8s),
					testAccCheckDigitalOceanKubernetesNodePoolExists("digitalocean_kubernetes_node_pool.barfoo", &k8s, &k8sPool),
					resource.TestCheckResourceAttr("digitalocean_kubernetes_cluster.foobar", "name", rName),
				),
			},
		},
	})
}

func testAccDigitalOceanKubernetesConfigBasicWithNodePool(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_kubernetes_cluster" "foobar" {
	name    = "%s"
	region  = "lon1"
	version = "1.12.1-do.2"
	tags    = ["foo","bar"]

	node_pool {
		size  = "s-1vcpu-2gb"
		count = 1
		tags  = ["one","two"]
	}
}

resource digitalocean_kubernetes_node_pool "barfoo" {
  cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"

	name    = "%s"
	size  = "s-1vcpu-2gb"
	count = 1
	tags  = ["three","four"]
}
`, rName, rName)
}

func testAccCheckDigitalOceanKubernetesNodePoolExists(n string, cluster *godo.KubernetesCluster, pool *godo.KubernetesNodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*godo.Client)

		foundPool, _, err := client.Kubernetes.GetNodePool(context.Background(), cluster.ID, rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundPool.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*pool = *foundPool

		return nil
	}
}
