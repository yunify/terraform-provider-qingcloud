/**
 * Copyright (c) 2016 Magicshui
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
/**
 * Copyright (c) 2017 yunify
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package qingcloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
	"log"
)

func TestAccQingcloudAPP_Tomcat(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPTomcat,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
					//resource.TestCheckResourceAttr(
					//	"qingcloud_app.foo", "name", "Tomcat Cluster 003"),
				),
			},
		},
	})
}

func TestAccQingcloudAPP_Zookeeper(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPZookeeper,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
					//resource.TestCheckResourceAttr(
					//	"qingcloud_app.foo", "name", "Tomcat Cluster 003"),
				),
			},
		},
	})
}

func TestAccQingcloudAPP_KafKa(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.kafka",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPKafka,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
					//resource.TestCheckResourceAttr(
					//	"qingcloud_app.foo", "name", "Tomcat Cluster 003"),
				),
			},
		},
	})
}

func TestAccQingcloudAPP_Etcd(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPEtcd,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
				),
			},
		},
	})
}

func TestAccQingcloudAPP_ELK(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPELK,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
				),
			},
		},
	})
}

func TestAccQingcloudAPP_K8S(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.k8s",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPKubernetes,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
				),
			},
		},
	})
}

func TestAccQingcloudAPP_RadonDB(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPRadonDB,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
				),
			},
		},
	})
}


func TestAccQingcloudAPP_PostgreSql(t *testing.T) {
	var cluster qc.DescribeClustersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_app.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAPPPostgresql,
				Check: resource.ComposeTestCheckFunc(
					testCheckClusterExists(
						"qingcloud_cluster.foo", &cluster),
				),
			},
		},
	})
}

func testCheckClusterExists(n string, cluster *qc.DescribeClustersOutput) resource.TestCheckFunc{
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			fmt.Println("------------",rs)
			if rs.Type != "qingcloud_app" {
				continue
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("No cluster ID is set")
			}
			client := testAccProvider.Meta().(*QingCloudClient)
			input := new(qc.DescribeClustersInput)
			input.Clusters = []*string{qc.String(rs.Primary.ID)}
			input.Verbose = qc.Int(1)
			d, err := client.cluster.DescribeClusters(input)
			log.Printf("[WARN] cluster id %#v", rs.Primary.ID)
			if err != nil {
				return err
			}
			if d == nil || qc.StringValue(d.ClusterSet[0].ClusterID) == "" {
				return fmt.Errorf("cluster not found")
			}
			*cluster = *d
		}
		return nil
	}
}

func testCheckClusterDestroy(s *terraform.State) error {
	return testAccCheckClusterDestroy(s, testAccProvider)
}

func testAccCheckClusterDestroy(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_app" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeClustersInput)
		input.Clusters = []*string{qc.String(rs.Primary.ID)}
		output, err := client.cluster.DescribeClusters(input)
		if err == nil {
			if len(output.ClusterSet) != 0 && qc.StringValue(output.ClusterSet[0].Status) != "deleted" && qc.StringValue(output.ClusterSet[0].Status) != "ceased"{
				return fmt.Errorf("Found  cluster: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAPPConfig = `
resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-jwq1fzqo"
	version_id = "appv-f91dxrsf"
    name = "Tomcat Cluster 003"
    conf = "{\"cluster\":{\"name\":\"Tomcat Cluster 003\",\"description\":\"\",\"tomcat_nodes\":{\"loadbalancer\":[{\"listener\":\"lbl-lp3l5klf\",\"port\":8080,\"policy\":\"lbp-ehq5mruj\"}],\"cpu\":1,\"memory\":2048,\"instance_class\":0,\"count\":2,\"volume_size\":10},\"log_node\":{\"cpu\":1,\"memory\":2048,\"instance_class\":0,\"volume_size\":10},\"vxnet\":\"vxnet-1ezfj3y\",\"global_uuid\":\"09122219576047188\"},\"version\":\"appv-f91dxrsf\",\"env\":{\"java_version\":\"8\",\"tomcat_version\":\"7\",\"tomcat_user\":\"qingAdmin\",\"tomcat_pwd\":\"qing0pwd\",\"war_source\":\"tomcat_manager\",\"tomcat_encoding\":\"UTF-8\",\"tomcat_log_level\":\"INFO\",\"tomcat_log_packages\":\"\",\"threadpool_maxThreads\":\"200\",\"threadpool_minSpareThreads\":\"25\",\"threadpool_maxIdleTime\":\"60000\",\"java_opts\":\"\",\"redis_db_num\":\"0\",\"access_key_id\":\"\",\"zone\":\"pek3a\",\"bucket\":\"\",\"war_name\":\"\",\"mysql_db_name\":\"testdb\",\"jdbc_dsname\":\"TestDB\",\"jdbc_maxActive\":\"100\",\"jdbc_maxIdle\":\"30\",\"jdbc_maxWait\":\"30000\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPConfig1 = `
resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-jwq1fzqo"
	version_id = "appv-f91dxrsf"
    conf = "{\"cluster\":{\"name\":\"Tomcat Cluster 003\",\"description\":\"\",\"tomcat_nodes\":{\"loadbalancer\":[{\"listener\":\"lbl-lp3l5klf\",\"port\":8080,\"policy\":\"lbp-ehq5mruj\"}],\"cpu\":1,\"memory\":2048,\"instance_class\":0,\"count\":2,\"volume_size\":10},\"log_node\":{\"cpu\":1,\"memory\":2048,\"instance_class\":0,\"volume_size\":10},\"vxnet\":\"vxnet-1ezfj3y\",\"global_uuid\":\"09122219576047188\"},\"version\":\"appv-f91dxrsf\",\"env\":{\"java_version\":\"8\",\"tomcat_version\":\"7\",\"tomcat_user\":\"qingAdmin\",\"tomcat_pwd\":\"qing0pwd\",\"war_source\":\"tomcat_manager\",\"tomcat_encoding\":\"UTF-8\",\"tomcat_log_level\":\"INFO\",\"tomcat_log_packages\":\"\",\"threadpool_maxThreads\":\"200\",\"threadpool_minSpareThreads\":\"25\",\"threadpool_maxIdleTime\":\"60000\",\"java_opts\":\"\",\"redis_db_num\":\"0\",\"access_key_id\":\"\",\"zone\":\"pek3a\",\"bucket\":\"\",\"war_name\":\"\",\"mysql_db_name\":\"testdb\",\"jdbc_dsname\":\"TestDB\",\"jdbc_maxActive\":\"100\",\"jdbc_maxIdle\":\"30\",\"jdbc_maxWait\":\"30000\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPTomcat = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_eip" "foo" {
    bandwidth = 2
	name = "terrafromtest_eip"
}

resource "qingcloud_loadbalancer" "foo" {
	name = "terraformtest_lb"
	eip_ids =["${qingcloud_eip.foo.id}"]
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  name = "terraformtest_lbpolicy"
  listener_port = "80"
  listener_protocol = "http"
}

resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-jwq1fzqo"
	version_id = "appv-f91dxrsf"
	name = "Tomcat Cluster 003"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
	lb_listener_id = "${qingcloud_loadbalancer_listener.foo.id}"
    conf = "{\"cluster\":{\"name\":\"1234\",\"description\":\"\",\"tomcat_nodes\":{\"loadbalancer\":1234,\"cpu\":1,\"memory\":2048,\"instance_class\":0,\"count\":2,\"volume_size\":10},\"log_node\":{\"cpu\":1,\"memory\":2048,\"instance_class\":0,\"volume_size\":10},\"global_uuid\":\"09122219576047188\"},\"version\":\"appv-f91dxrsf\",\"env\":{\"java_version\":\"8\",\"tomcat_version\":\"7\",\"tomcat_user\":\"qingAdmin\",\"tomcat_pwd\":\"qing0pwd\",\"war_source\":\"tomcat_manager\",\"tomcat_encoding\":\"UTF-8\",\"tomcat_log_level\":\"INFO\",\"tomcat_log_packages\":\"\",\"threadpool_maxThreads\":\"200\",\"threadpool_minSpareThreads\":\"25\",\"threadpool_maxIdleTime\":\"60000\",\"java_opts\":\"\",\"redis_db_num\":\"0\",\"access_key_id\":\"\",\"zone\":\"pek3a\",\"bucket\":\"\",\"war_name\":\"\",\"mysql_db_name\":\"testdb\",\"jdbc_dsname\":\"TestDB\",\"jdbc_maxActive\":\"100\",\"jdbc_maxIdle\":\"30\",\"jdbc_maxWait\":\"30000\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPZookeeper = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-tg3lbp0a"
	version_id = "appv-9b7na511"
	name = "TerraformTestZooKeeper"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"xxx\",\"description\":\"\",\"zk_node\":{\"cpu\":1,\"memory\":1024,\"instance_class\":0,\"count\":3,\"volume_size\":10},\"global_uuid\":\"45292219573932653\"},\"version\":\"appv-9b7na511\",\"env\":{\"admin_enabled\":false,\"admin_username\":\"super\",\"admin_password\":\"Super12345\",\"tickTime\":2000,\"initLimit\":10,\"syncLimit\":5,\"maxClientCnxns\":1000,\"autopurge_snapRetainCount\":3,\"autopurge_purgeInterval\":0}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPKafka = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "zk" {
	app_type = "cluster"
	app_id = "app-tg3lbp0a"
	version_id = "appv-9b7na511"
	name = "TerraformTestZooKeeper"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"xxx\",\"description\":\"\",\"zk_node\":{\"cpu\":1,\"memory\":1024,\"instance_class\":0,\"count\":3,\"volume_size\":10},\"global_uuid\":\"45292219573932653\"},\"version\":\"appv-9b7na511\",\"env\":{\"admin_enabled\":false,\"admin_username\":\"super\",\"admin_password\":\"Super12345\",\"tickTime\":2000,\"initLimit\":10,\"syncLimit\":5,\"maxClientCnxns\":1000,\"autopurge_snapRetainCount\":3,\"autopurge_purgeInterval\":0}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}

resource "qingcloud_app" "kafka" {
	app_type = "cluster"
	app_id = "app-n9ro0xcp"
	version_id = "appv-vx6yl2x5"
	name = "TerraformTestKafka"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    zk_service = "${qingcloud_app.zk.id}"
    conf = "{\"cluster\":{\"name\":\"xxxx\",\"description\":\"\",\"kafka\":{\"cpu\":1,\"memory\":2048,\"count\":3,\"instance_class\":0,\"volume_size\":30},\"client\":{\"cpu\":1,\"memory\":1024,\"count\":1,\"instance_class\":0,\"volume_size\":10},\"vxnet\":\"xxxxx\",\"zk_service\":\"xxxxxx\",\"global_uuid\":\"90302219564228453\"},\"version\":\"appv-vx6yl2x5\",\"env\":{\"offsets_topic_replication_factor\":3,\"kafka-manager_basicAuthentication_enabled\":false,\"kafka-manager_basicAuthentication_username\":\"admin\",\"kafka-manager_basicAuthentication_password\":\"password\",\"log_retention_bytes\":9663676416,\"log_retention_hours\":168,\"log_segment_bytes\":1073741824,\"log_segment_delete_delay_ms\":60000,\"log_roll_hours\":168,\"auto_create_topics_enable\":true,\"default_replication_factor\":1,\"delete_topic_enable\":true,\"log_cleanup_policy\":\"delete\",\"log_cleaner_enable\":false,\"compression_type\":\"producer\",\"message_max_bytes\":1000000,\"num_io_threads\":8,\"num_partitions\":3,\"num_replica_fetchers\":1,\"queued_max_requests\":500,\"socket_receive_buffer_bytes\":102400,\"socket_send_buffer_bytes\":102400,\"unclean_leader_election_enable\":true,\"advertised_host_name\":\"\",\"advertised_port\":9092,\"kafka-manager_port\":9000}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPEtcd = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-fdyvu2wk"
	version_id = "appv-yycc2fun"
	name = "TerraformTestEtcd"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"xxxx\",\"description\":\"\",\"etcd_node\":{\"cpu\":2,\"memory\":8192,\"count\":3,\"instance_class\":1,\"volume_class\":3,\"volume_size\":10},\"etcd_proxy\":{\"cpu\":1,\"memory\":2048,\"count\":0,\"instance_class\":1},\"vxnet\":\"xxxxx\",\"global_uuid\":\"72002219564220947\"},\"version\":\"appv-yycc2fun\",\"env\":{\"etcd_node.autocompact\":0}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPELK = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-p6au3oyq"
	version_id = "appv-dd93di9d"
	name = "TerraformTestElk"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"ELK\",\"description\":\"\",\"es_node\":{\"cpu\":2,\"memory\":4096,\"count\":3,\"instance_class\":0,\"volume_size\":10},\"kbn_node\":{\"cpu\":2,\"memory\":4096,\"count\":1,\"instance_class\":0},\"lst_node\":{\"cpu\":2,\"memory\":4096,\"count\":1,\"instance_class\":0,\"volume_size\":10},\"vxnet\":\"xxxxx\",\"global_uuid\":\"53002219564294789\"},\"version\":\"appv-dd93di9d\",\"env\":{\"es_node.indices.fielddata.cache.size\":\"90%\",\"es_node.gateway.recover_after_time\":\"5m\",\"es_node.http.cors.allow-origin\":\"*\",\"es_node.indices.queries.cache.size\":\"10%\",\"es_node.indices.memory.index_buffer_size\":\"10%\",\"es_node.indices.requests.cache.size\":\"2%\",\"lst_node.input_conf_content\":\"http { host => \\\"0.0.0.0\\\"  port => 9700 }\",\"es_node.action.destructive_requires_name\":\"true\",\"es_node.logstash_node_ip\":\"\",\"es_node.discovery.zen.no_master_block\":\"write\",\"es_node.http.cors.enabled\":\"true\",\"es_node.script.inline\":\"true\",\"es_node.script.stored\":\"true\",\"es_node.script.file\":\"false\",\"es_node.script.aggs\":\"true\",\"es_node.script.search\":\"true\",\"es_node.script.update\":\"true\",\"es_node.remote_ext_dict\":\"\",\"es_node.remote_ext_stopwords\":\"\",\"es_node.path.repo\":\"[]\",\"es_node.repositories.url.allowed_urls\":\"[]\",\"es_node.es_additional_line1\":\"\",\"es_node.es_additional_line2\":\"\",\"es_node.es_additional_line3\":\"\",\"es_node.logger.action.level\":\"info\",\"es_node.rootLogger.level\":\"info\",\"es_node.logger.deprecation.level\":\"warn\",\"es_node.logger.index_search_slowlog_rolling.level\":\"trace\",\"es_node.logger.index_indexing_slowlog.level\":\"trace\",\"es_node.enable_heap_dump\":false,\"es_node.heap_dump_path\":\"/data/elasticsearch/dump\",\"es_node.clean_logs_older_than_n_days\":7,\"kbn_node.enable_elasticsearch_head\":true,\"kbn_node.enable_elastichd\":true,\"kbn_node.enable_cerebro\":true,\"kbn_node.enable_elasticsearch_sql\":true,\"kbn_node.nginx_client_max_body_size\":\"20m\",\"lst_node.filter_conf_content\":\"\",\"lst_node.output_conf_content\":\"\",\"lst_node.output_es_content\":\"\",\"lst_node.gemfile_append_content\":\"\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPKubernetes = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "elk" {
	app_type = "cluster"
	app_id = "app-p6au3oyq"
	version_id = "appv-dd93di9d"
	name = "TerraformTestElk"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"ELK\",\"description\":\"\",\"es_node\":{\"cpu\":2,\"memory\":4096,\"count\":3,\"instance_class\":0,\"volume_size\":10},\"kbn_node\":{\"cpu\":2,\"memory\":4096,\"count\":1,\"instance_class\":0},\"lst_node\":{\"cpu\":2,\"memory\":4096,\"count\":1,\"instance_class\":0,\"volume_size\":10},\"vxnet\":\"xxxxx\",\"global_uuid\":\"53002219564294789\"},\"version\":\"appv-dd93di9d\",\"env\":{\"es_node.indices.fielddata.cache.size\":\"90%\",\"es_node.gateway.recover_after_time\":\"5m\",\"es_node.http.cors.allow-origin\":\"*\",\"es_node.indices.queries.cache.size\":\"10%\",\"es_node.indices.memory.index_buffer_size\":\"10%\",\"es_node.indices.requests.cache.size\":\"2%\",\"lst_node.input_conf_content\":\"http { host => \\\"0.0.0.0\\\"  port => 9700 }\",\"es_node.action.destructive_requires_name\":\"true\",\"es_node.logstash_node_ip\":\"\",\"es_node.discovery.zen.no_master_block\":\"write\",\"es_node.http.cors.enabled\":\"true\",\"es_node.script.inline\":\"true\",\"es_node.script.stored\":\"true\",\"es_node.script.file\":\"false\",\"es_node.script.aggs\":\"true\",\"es_node.script.search\":\"true\",\"es_node.script.update\":\"true\",\"es_node.remote_ext_dict\":\"\",\"es_node.remote_ext_stopwords\":\"\",\"es_node.path.repo\":\"[]\",\"es_node.repositories.url.allowed_urls\":\"[]\",\"es_node.es_additional_line1\":\"\",\"es_node.es_additional_line2\":\"\",\"es_node.es_additional_line3\":\"\",\"es_node.logger.action.level\":\"info\",\"es_node.rootLogger.level\":\"info\",\"es_node.logger.deprecation.level\":\"warn\",\"es_node.logger.index_search_slowlog_rolling.level\":\"trace\",\"es_node.logger.index_indexing_slowlog.level\":\"trace\",\"es_node.enable_heap_dump\":false,\"es_node.heap_dump_path\":\"/data/elasticsearch/dump\",\"es_node.clean_logs_older_than_n_days\":7,\"kbn_node.enable_elasticsearch_head\":true,\"kbn_node.enable_elastichd\":true,\"kbn_node.enable_cerebro\":true,\"kbn_node.enable_elasticsearch_sql\":true,\"kbn_node.nginx_client_max_body_size\":\"20m\",\"lst_node.filter_conf_content\":\"\",\"lst_node.output_conf_content\":\"\",\"lst_node.output_es_content\":\"\",\"lst_node.gemfile_append_content\":\"\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}

resource "qingcloud_app" "etcd" {
	app_type = "cluster"
	app_id = "app-fdyvu2wk"
	version_id = "appv-yycc2fun"
	name = "TerraformTestEtcd"
	elk_service = "${qingcloud_app.elk.id}"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"xxxx\",\"description\":\"\",\"etcd_node\":{\"cpu\":2,\"memory\":8192,\"count\":3,\"instance_class\":1,\"volume_class\":3,\"volume_size\":10},\"etcd_proxy\":{\"cpu\":1,\"memory\":2048,\"count\":0,\"instance_class\":1},\"vxnet\":\"xxxxx\",\"global_uuid\":\"72002219564220947\"},\"version\":\"appv-yycc2fun\",\"env\":{\"etcd_node.autocompact\":0}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}

resource "qingcloud_app" "k8s" {
	app_type = "cluster"
	app_id = "app-u0llx5j8"
	version_id = "appv-vtpbh5fo"
	name = "TerraformTestK8s"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
	elk_service = "${qingcloud_app.elk.id}"
	etcd_service = "${qingcloud_app.etcd.id}"
    conf = "{\"cluster\":{\"name\":\"my k8s cluster\",\"description\":\"my k8s cluster\",\"master\":{\"cpu\":4,\"memory\":4096,\"count\":1,\"instance_class\":0,\"volume_size\":20},\"node\":{\"cpu\":4,\"memory\":8192,\"count\":2,\"volume_size\":50},\"ssd_node\":{\"cpu\":4,\"memory\":4096,\"count\":0,\"volume_size\":50},\"log\":{\"cpu\":4,\"memory\":4096,\"count\":2,\"instance_class\":0,\"volume_size\":100},\"client\":{\"cpu\":1,\"memory\":1024,\"instance_class\":0},\"vxnet\":\"xxxxxxx\",\"elk_service\":\"xxxxxxx\",\"etcd_service\":\"xxxxxxxx\",\"global_uuid\":\"30002219564297805\"},\"version\":\"appv-vtpbh5fo\",\"resource_group\":\"Basic\",\"env\":{\"access_key_id\":\"TRUJJXSHJOXOSTVEJDYH\",\"secret_access_key\":\"HiFu0hvWaGJIkqYjX9orTkgroOCMtTPZR20ucAx1\",\"pod_vxnets\":\"vxnet-1ezfj3y\",\"api_external_domain\":\"k8s.cluster.local\",\"cluster_port_range\":\"30000-32767\",\"max_pods\":60,\"enable_hostnic\":\"calico\",\"docker_bip\":\"172.30.0.1/16\",\"registry-mirrors\":\"https://registry.docker-cn.com\",\"insecure-registries\":\"\",\"private-registry\":\"\",\"dockerhub_username\":\"\",\"keep_log_days\":3,\"kube_log_level\":0,\"fluent_forward_server\":\"\",\"es_server\":\"\",\"enable_istio\":\"no\"},\"toggle_passwd\":\"on\"}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPRadonDB = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-x87yr57g"
	version_id = "appv-bu2m1f9t"
	name = "TerraformTestRadonDB"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"QingCloud RadonDB\",\"description\":\"\",\"auto_backup_time\":\"-1\",\"radon\":{\"cpu\":4,\"memory\":4096,\"instance_class\":1,\"replica\":1,\"volume_size\":50},\"xenon\":{\"cpu\":4,\"memory\":4096,\"count\":2,\"instance_class\":1,\"volume_size\":50},\"xenon-ap\":{\"cpu\":8,\"memory\":8192,\"instance_class\":1,\"volume_size\":100},\"vxnet\":\"xxxxxxxxx\",\"global_uuid\":\"14912219564210838\"},\"version\":\"appv-bu2m1f9t\",\"env\":{\"user\":\"test\",\"password\":\"ZHu88jie\",\"port\":3306,\"twopc-enable\":\"False\",\"allowip\":\"\",\"max-result-size\":1073741824,\"ddl-timeout\":36000000,\"query-timeout\":300000,\"audit-mode\":\"None\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`

const testAPPPostgresql = `
resource "qingcloud_security_group" "foo" {
    name = "terraformtest_sg"
}

resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
    name = "terraformtest_app"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    name = "terraformtest_vxnet"
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_app" "foo" {
	app_type = "cluster"
	app_id = "app-gtusp816"
	version_id = "appv-bxyz959m"
	name = "TerraformTestPostgreSql"
	vxnet_id = "${qingcloud_vxnet.foo.id}"
    conf = "{\"cluster\":{\"name\":\"PostgreSQL10 Cluster\",\"description\":\"\",\"auto_backup_time\":\"-1\",\"pg\":{\"cpu\":2,\"memory\":4096,\"instance_class\":1,\"volume_size\":20},\"vxnet\":\"xxxxxxx\",\"global_uuid\":\"65912219564211189\"},\"version\":\"appv-bxyz959m\",\"env\":{\"SyncStreamRepl\":\"Yes\",\"AutoFailover\":\"No\",\"DBname\":\"qingcloud\",\"DBusername\":\"qingcloud\",\"DBpassword\":\"qingcloud1234\",\"max_connections\":\"auto-optimized-conns\",\"wal_buffers\":\"8MB\",\"work_mem\":\"4MB\",\"maintenance_work_mem\":\"64MB\",\"effective_cache_size\":\"4GB\",\"wal_keep_segments\":256,\"checkpoint_timeout\":\"5min\",\"autovacuum\":\"on\",\"vacuum_cost_delay\":0,\"autovacuum_naptime\":\"1min\",\"vacuum_cost_limit\":200,\"bgwriter_delay\":200,\"bgwriter_lru_multiplier\":2,\"wal_writer_delay\":200,\"fsync\":\"on\",\"commit_delay\":0,\"commit_siblings\":5,\"enable_bitmapscan\":\"on\",\"enable_seqscan\":\"on\",\"full_page_writes\":\"on\",\"log_min_messages\":\"warning\",\"deadlock_timeout\":1,\"log_lock_waits\":\"off\",\"log_min_duration_statement\":-1,\"temp_buffers\":\"8MB\",\"max_prepared_transactions\":0,\"max_wal_senders\":10,\"bgwriter_lru_maxpages\":100,\"log_statement\":\"none\",\"wal_level\":\"replica\",\"shared_buffers\":\"auto-optimized-sharedbuffers\"}}"
    debug = 0
    owner = "usr-UdsBKvWf"
}
`