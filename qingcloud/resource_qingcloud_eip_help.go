package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyEipAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	input := new(qc.ModifyEIPAttributesInput)
	input.EIP = qc.String(d.Id())
	nameUpdate := false
	descriptionUpdate := false
	input.EIPName, nameUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	if nameUpdate || descriptionUpdate {
		var err error
		simpleRetry(func() error {
			_, err = clt.ModifyEIPAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func changeEIPBandwidth(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange(resourceEipBandwidth) && !d.IsNewResource() {
		clt := meta.(*QingCloudClient).eip
		input := new(qc.ChangeEIPsBandwidthInput)
		input.EIPs = []*string{qc.String(d.Id())}
		input.Bandwidth = qc.Int(d.Get(resourceEipBandwidth).(int))
		var err error
		simpleRetry(func() error {
			_, err = clt.ChangeEIPsBandwidth(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := EIPTransitionStateRefresh(clt, d.Id()); err != nil {
			return nil
		}
	}
	return nil
}
func changeEIPBillMode(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	if d.HasChange(resourceEipBillMode) && !d.IsNewResource() {
		input := new(qc.ChangeEIPsBillingModeInput)
		input.EIPs = []*string{qc.String(d.Id())}
		input.BillingMode = qc.String(d.Get(resourceEipBillMode).(string))
		var err error
		simpleRetry(func() error {
			_, err = clt.ChangeEIPsBillingMode(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := EIPTransitionStateRefresh(clt, d.Id()); err != nil {
			return nil
		}
	}
	return nil
}

func getEIPResourceMap(data *qc.EIP) map[string]interface{} {
	var a = make(map[string]interface{}, 3)
	a["resource_id"] = qc.StringValue(data.Resource.ResourceID)
	a["resource_name"] = qc.StringValue(data.Resource.ResourceName)
	a["resource_type"] = qc.StringValue(data.Resource.ResourceType)
	return a
}
func waitEipLease(d *schema.ResourceData, meta interface{}) error {
	if !d.IsNewResource() {
		return nil
	}
	clt := meta.(*QingCloudClient).eip
	input := new(qc.DescribeEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeEIPs(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	//wait for lease info
	WaitForLease(output.EIPSet[0].StatusTime)
	return nil
}

func isEipDeleted(eipSet []*qc.EIP) bool {
	if len(eipSet) == 0 || qc.StringValue(eipSet[0].Status) == "ceased" || qc.StringValue(eipSet[0].Status) == "ceased" {
		return true
	}
	return false
}
