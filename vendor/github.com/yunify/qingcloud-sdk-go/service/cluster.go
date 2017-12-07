// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package service

import (
	"fmt"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

var _ fmt.State
var _ time.Time

type ClusterService struct {
	Config     *config.Config
	Properties *ClusterServiceProperties
}

type ClusterServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Cluster(zone string) (*ClusterService, error) {
	properties := &ClusterServiceProperties{
		Zone: &zone,
	}

	return &ClusterService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/add_cluster_nodes.html
func (s *ClusterService) AddClusterNodes(i *AddClusterNodesInput) (*AddClusterNodesOutput, error) {
	if i == nil {
		i = &AddClusterNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddClusterNodes",
		RequestMethod: "GET",
	}

	x := &AddClusterNodesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type AddClusterNodesInput struct {
	Cluster    *string   `json:"cluster" name:"cluster" location:"params"`
	NodeCount  *int      `json:"node_count" name:"node_count" location:"params"`
	NodeName   *string   `json:"node_name" name:"node_name" location:"params"`
	NodeRole   *string   `json:"node_role" name:"node_role" location:"params"`
	PrivateIPs []*string `json:"private_ips" name:"private_ips" location:"params"`
}

func (v *AddClusterNodesInput) Validate() error {

	return nil
}

type AddClusterNodesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/change_cluster_vxnet.html
func (s *ClusterService) ChangeClusterVxNet(i *ChangeClusterVxNetInput) (*ChangeClusterVxNetOutput, error) {
	if i == nil {
		i = &ChangeClusterVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeClusterVxnet",
		RequestMethod: "GET",
	}

	x := &ChangeClusterVxNetOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ChangeClusterVxNetInput struct {
	Cluster    *string     `json:"cluster" name:"cluster" location:"params"`
	PrivateIPs interface{} `json:"private_ips" name:"private_ips" location:"params"`
	Roles      []*string   `json:"roles" name:"roles" location:"params"`
	VxNet      *string     `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *ChangeClusterVxNetInput) Validate() error {

	return nil
}

type ChangeClusterVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/create_clusters.html
func (s *ClusterService) CreateCluster(i *CreateClusterInput) (*CreateClusterOutput, error) {
	if i == nil {
		i = &CreateClusterInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateCluster",
		RequestMethod: "GET",
	}

	x := &CreateClusterOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type CreateClusterInput struct {
	Conf *string `json:"conf" name:"conf" location:"params"` // Required
}

func (v *CreateClusterInput) Validate() error {

	if v.Conf == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Conf",
			ParentName:    "CreateClusterInput",
		}
	}

	return nil
}

type CreateClusterOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/delete_cluster_nodes.html
func (s *ClusterService) DeleteClusterNodes(i *DeleteClusterNodesInput) (*DeleteClusterNodesOutput, error) {
	if i == nil {
		i = &DeleteClusterNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteClusterNodes",
		RequestMethod: "GET",
	}

	x := &DeleteClusterNodesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DeleteClusterNodesInput struct {
	Cluster *string   `json:"cluster" name:"cluster" location:"params"`
	Force   *int      `json:"force" name:"force" location:"params"`
	Nodes   []*string `json:"nodes" name:"nodes" location:"params"`
}

func (v *DeleteClusterNodesInput) Validate() error {

	return nil
}

type DeleteClusterNodesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/delete_clusters.html
func (s *ClusterService) DeleteClusters(i *DeleteClustersInput) (*DeleteClustersOutput, error) {
	if i == nil {
		i = &DeleteClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteClusters",
		RequestMethod: "GET",
	}

	x := &DeleteClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DeleteClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"`
}

func (v *DeleteClustersInput) Validate() error {

	return nil
}

type DeleteClustersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_cluster_nodes.html
func (s *ClusterService) DescribeClusterNodes(i *DescribeClusterNodesInput) (*DescribeClusterNodesOutput, error) {
	if i == nil {
		i = &DescribeClusterNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusterNodes",
		RequestMethod: "GET",
	}

	x := &DescribeClusterNodesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClusterNodesInput struct {
	Cluster      *string   `json:"cluster" name:"cluster" location:"params"`
	ClusterNodes []*string `json:"cluster_nodes" name:"cluster_nodes" location:"params"`
	Role         *string   `json:"role" name:"role" location:"params"`
}

func (v *DescribeClusterNodesInput) Validate() error {

	return nil
}

type DescribeClusterNodesOutput struct {
	Message    *string        `json:"message" name:"message"`
	Action     *string        `json:"action" name:"action" location:"elements"`
	NodeSet    []*ClusterNode `json:"node_set" name:"node_set" location:"elements"`
	RetCode    *int           `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int           `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_cluster_users.html
func (s *ClusterService) DescribeClusterUsers(i *DescribeClusterUsersInput) (*DescribeClusterUsersOutput, error) {
	if i == nil {
		i = &DescribeClusterUsersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusterUsers",
		RequestMethod: "GET",
	}

	x := &DescribeClusterUsersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClusterUsersInput struct {
	AppVersions   []*string `json:"app_versions" name:"app_versions" location:"params"`
	Apps          []*string `json:"apps" name:"apps" location:"params"`
	ClusterStatus []*string `json:"cluster_status" name:"cluster_status" location:"params"`
	Users         []*string `json:"users" name:"users" location:"params"`
	Zones         []*string `json:"zones" name:"zones" location:"params"`
}

func (v *DescribeClusterUsersInput) Validate() error {

	return nil
}

type DescribeClusterUsersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_clusters.html
func (s *ClusterService) DescribeClusters(i *DescribeClustersInput) (*DescribeClustersOutput, error) {
	if i == nil {
		i = &DescribeClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusters",
		RequestMethod: "GET",
	}

	x := &DescribeClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClustersInput struct {
	AppID      []*string `json:"app_id" name:"app_id" location:"params"`
	AppVersion []*string `json:"app_version" name:"app_version" location:"params"`
	Clusters   []*string `json:"clusters" name:"clusters" location:"params"`
	Role       *string   `json:"role" name:"role" location:"params"`
	// Scope's available values: all, cfgmgmt
	Scope *string   `json:"scope" name:"scope" location:"params"`
	Users []*string `json:"users" name:"users" location:"params"`
}

func (v *DescribeClustersInput) Validate() error {

	if v.Scope != nil {
		scopeValidValues := []string{"all", "cfgmgmt"}
		scopeParameterValue := fmt.Sprint(*v.Scope)

		scopeIsValid := false
		for _, value := range scopeValidValues {
			if value == scopeParameterValue {
				scopeIsValid = true
			}
		}

		if !scopeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Scope",
				ParameterValue: scopeParameterValue,
				AllowedValues:  scopeValidValues,
			}
		}
	}

	return nil
}

type DescribeClustersOutput struct {
	Message    *string    `json:"message" name:"message"`
	Action     *string    `json:"action" name:"action" location:"elements"`
	ClusterSet []*Cluster `json:"cluster_set" name:"cluster_set" location:"elements"`
	RetCode    *int       `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int       `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/get_clusters_stats.html
func (s *ClusterService) GetClustersStats(i *GetClustersStatsInput) (*GetClustersStatsOutput, error) {
	if i == nil {
		i = &GetClustersStatsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetClustersStats",
		RequestMethod: "GET",
	}

	x := &GetClustersStatsOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type GetClustersStatsInput struct {
	Zones []*string `json:"zones" name:"zones" location:"params"`
}

func (v *GetClustersStatsInput) Validate() error {

	return nil
}

type GetClustersStatsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/modify_cluster_attributes.html
func (s *ClusterService) ModifyClusterAttributes(i *ModifyClusterAttributesInput) (*ModifyClusterAttributesOutput, error) {
	if i == nil {
		i = &ModifyClusterAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyClusterAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyClusterAttributesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ModifyClusterAttributesInput struct {
	Cluster     *string `json:"cluster" name:"cluster" location:"params"`
	Description *string `json:"description" name:"description" location:"params"`
	Name        *string `json:"name" name:"name" location:"params"`
}

func (v *ModifyClusterAttributesInput) Validate() error {

	return nil
}

type ModifyClusterAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/modify_cluster_node_attributes.html
func (s *ClusterService) ModifyClusterNodeAttributes(i *ModifyClusterNodeAttributesInput) (*ModifyClusterNodeAttributesOutput, error) {
	if i == nil {
		i = &ModifyClusterNodeAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyClusterNodeAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyClusterNodeAttributesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ModifyClusterNodeAttributesInput struct {
	Cluster     *string `json:"cluster" name:"cluster" location:"params"`
	ClusterNode *string `json:"cluster_node" name:"cluster_node" location:"params"`
	Name        *string `json:"name" name:"name" location:"params"`
}

func (v *ModifyClusterNodeAttributesInput) Validate() error {

	return nil
}

type ModifyClusterNodeAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/resize_cluster.html
func (s *ClusterService) ResizeCluster(i *ResizeClusterInput) (*ResizeClusterOutput, error) {
	if i == nil {
		i = &ResizeClusterInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeCluster",
		RequestMethod: "GET",
	}

	x := &ResizeClusterOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ResizeClusterInput struct {
	Cluster     *string `json:"cluster" name:"cluster" location:"params"`
	CPU         *int    `json:"cpu" name:"cpu" location:"params"`
	Memory      *int    `json:"memory" name:"memory" location:"params"`
	NodeRole    *string `json:"node_role" name:"node_role" location:"params"`
	StorageSize *int    `json:"storage_size" name:"storage_size" location:"params"`
}

func (v *ResizeClusterInput) Validate() error {

	return nil
}

type ResizeClusterOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/restart_cluster_service.html
func (s *ClusterService) RestartClusterService(i *RestartClusterServiceInput) (*RestartClusterServiceOutput, error) {
	if i == nil {
		i = &RestartClusterServiceInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RestartClusterService",
		RequestMethod: "GET",
	}

	x := &RestartClusterServiceOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type RestartClusterServiceInput struct {
	Cluster *string `json:"cluster" name:"cluster" location:"params"`
	Role    *string `json:"role" name:"role" location:"params"`
}

func (v *RestartClusterServiceInput) Validate() error {

	return nil
}

type RestartClusterServiceOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/start_clusters.html
func (s *ClusterService) StartClusters(i *StartClustersInput) (*StartClustersOutput, error) {
	if i == nil {
		i = &StartClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartClusters",
		RequestMethod: "GET",
	}

	x := &StartClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type StartClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"`
}

func (v *StartClustersInput) Validate() error {

	return nil
}

type StartClustersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	JobIDs  map[string]*string `json:"job_ids" name:"job_ids" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/stop_clusters.html
func (s *ClusterService) StopClusters(i *StopClustersInput) (*StopClustersOutput, error) {
	if i == nil {
		i = &StopClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopClusters",
		RequestMethod: "GET",
	}

	x := &StopClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type StopClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"`
	Force    *int      `json:"force" name:"force" location:"params"`
}

func (v *StopClustersInput) Validate() error {

	return nil
}

type StopClustersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	JobIDs  map[string]*string `json:"job_ids" name:"job_ids" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/suspend_clusters.html
func (s *ClusterService) SuspendClusters(i *SuspendClustersInput) (*SuspendClustersOutput, error) {
	if i == nil {
		i = &SuspendClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "SuspendClusters",
		RequestMethod: "GET",
	}

	x := &SuspendClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type SuspendClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"`
}

func (v *SuspendClustersInput) Validate() error {

	return nil
}

type SuspendClustersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/update_cluster_environment.html
func (s *ClusterService) UpdateClusterEnvironment(i *UpdateClusterEnvironmentInput) (*UpdateClusterEnvironmentOutput, error) {
	if i == nil {
		i = &UpdateClusterEnvironmentInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateClusterEnvironment",
		RequestMethod: "GET",
	}

	x := &UpdateClusterEnvironmentOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type UpdateClusterEnvironmentInput struct {
	Cluster *string     `json:"cluster" name:"cluster" location:"params"`
	Env     interface{} `json:"env" name:"env" location:"params"`
	Roles   []*string   `json:"roles" name:"roles" location:"params"`
}

func (v *UpdateClusterEnvironmentInput) Validate() error {

	return nil
}

type UpdateClusterEnvironmentOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/upgrade_clusters.html
func (s *ClusterService) UpgradeClusters(i *UpgradeClustersInput) (*UpgradeClustersOutput, error) {
	if i == nil {
		i = &UpgradeClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpgradeClusters",
		RequestMethod: "GET",
	}

	x := &UpgradeClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type UpgradeClustersInput struct {
	AppVersion *string   `json:"app_version" name:"app_version" location:"params"`
	Clusters   []*string `json:"clusters" name:"clusters" location:"params"`
}

func (v *UpgradeClustersInput) Validate() error {

	return nil
}

type UpgradeClustersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
