package qingcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

func ClusterDeleteStateRefresh(clt *qc.ClusterService, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeClustersInput)
		input.Clusters = []*string{qc.String(id)}
		var output *qc.DescribeClustersOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeClusters(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		if len(output.ClusterSet) != 1 {
			return output, "creating", nil
		}
		cluster := output.ClusterSet[0]
		return cluster, qc.StringValue(cluster.Status), nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "suspended","deleted"},
		Target:     []string{"ceased"},
		Refresh:    refreshFunc,
		Timeout:    30 * time.Minute,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func GetLbinfo(lb_policy_id string) ([]map[string]interface{}) {
	lbMap := []map[string]interface{}{}
	lb := map[string]interface{}{"listener": lb_policy_id, "port": 8080}
	lbMap = append(lbMap, lb)
	return lbMap;
}