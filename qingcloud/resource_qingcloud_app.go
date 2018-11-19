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
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"fmt"
	"encoding/json"
)

const (
//resourceVolumeSize = "size"
//resourceVolumeType = "type"
)

func resourceQingcloudApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudAppCreate,
		Read:   resourceQingcloudAppRead,
		Update: resourceQingcloudAppUpdate,
		Delete: resourceQingcloudAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"app_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"conf": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"debug": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vxnet_id":&schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"lb_listener_id":&schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"zk_service": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"elk_service": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"etcd_service": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudAppCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).app
	input := new(qc.DeployAppVersionInput)
	input.AppType = qc.String(d.Get("app_type").(string))
	input.VersionID = qc.String(d.Get("version_id").(string))
	input.AppID = qc.String(d.Get("app_id").(string))
	input.Conf = qc.String(d.Get("conf").(string))
	input.Owner = qc.String(d.Get("owner").(string))

	var dat map[string]interface{}
	conf := *input.Conf
	if err := json.Unmarshal([]byte(conf), &dat); err == nil {
		fmt.Println("==============json str 转map=======================")
		cluster := dat["cluster"]
		fmt.Println("cluster====",cluster)
		cluster_mp,_ := cluster.(map[string]interface{})
		if v, ok := d.GetOk("name"); ok  {
			cluster_mp["name"] = v.(string)
		}
		if v, ok := d.GetOk("vxnet_id"); ok {
			cluster_mp["vxnet"] = v.(string)
		}
		if v, ok := d.GetOk("lb_listener_id"); ok {
			if *input.AppType == "cluster" {
				tomcat_nodes := cluster_mp["tomcat_nodes"]
				tomcat_nodes_mp,_ := tomcat_nodes.(map[string]interface{})
				lb_policy_id := v.(string);
				lbInfo := GetLbinfo(lb_policy_id);
				tomcat_nodes_mp["loadbalancer"] = lbInfo
			}
		}
		if v, ok := d.GetOk("zk_service"); ok {
			cluster_mp["zk_service"] = v.(string)
		}
		if v, ok := d.GetOk("elk_service"); ok {
			cluster_mp["elk_service"] = v.(string)
		}
		if v, ok := d.GetOk("etcd_service"); ok {
			cluster_mp["etcd_service"] = v.(string)
		}
		dat["cluster"] = cluster_mp
		jsondat, _ := json.Marshal(dat)
		input.Conf = qc.String(string(jsondat))
		fmt.Println("conf====",*input.Conf)
	} else {
		fmt.Println(err)
	}

	fmt.Println("conf====",*input.Conf)
	var output *qc.DeployAppVersionOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DeployAppVersion(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	fmt.Println("output====",&output)
	d.SetId(*output.ClusterID);
	fmt.Println("ClusterID====",d.Id())
	clusterclt := meta.(*QingCloudClient).cluster
	if _, err = ClusterTransitionStateRefresh(clusterclt, d.Id()); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudAppRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cluster
	input := new(qc.DescribeClustersInput)
	input.Clusters = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeClustersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeClusters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	//if isInstanceDeleted(output.ClusterSet) {
	//	d.SetId("")
	//	return nil
	//}
	cluster := output.ClusterSet[0]
	fmt.Println("cluster====",cluster)
	//d.Set("app_id", qc.StringValue(cluster.AppID))
	//d.Set("app_version_info", qc.StringValue(cluster.AppVersionInfo.VersionID))
	//d.Set("cluster_id",qc.StringValue(cluster.ClusterID))
	d.Set("name",qc.StringValue(cluster.Name))
	//d.Set("status",cluster.Status)
	//d.Set("node_count",cluster.NodeCount)
	d.Set("vxnet_id",cluster.VxNet.VxNetID)

	return nil
}

func resourceQingcloudAppUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	//d.SetPartial("app_id")
	//d.SetPartial("app_version_info")
	d.SetPartial("cluster_id")
	//d.SetPartial("name")
	//d.SetPartial("status")
	//d.SetPartial("node_count")
	//d.SetPartial("vxnet_id")
	d.Partial(false)
	return resourceQingcloudAppRead(d, meta)
}

func resourceQingcloudAppDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cluster
	input := new(qc.DeleteClustersInput)
	input.Clusters = []*string{qc.String(d.Id())}

	input.DirectCease = qc.Int(1)
	var output *qc.DeleteClustersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DeleteClusters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	fmt.Println("===output===", &output)
	fmt.Println("job_id=",*output.JobIDs[d.Id()])

	if _, err := ClusterDeleteStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
