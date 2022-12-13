// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package illumio

import (
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi"
)

type UnmanagedWorkloadGenerator struct {
	IllumioService
}

func (g UnmanagedWorkloadGenerator) createResources(svc *illumioapi.PCE, workloads []illumioapi.Workload) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, workload := range workloads {
		resourceName := stripIdFromHref(workload.Href)
		resources = append(resources, terraformutils.NewResource(
			workload.Href,
			resourceName,
			"illumio-core_unmanaged_workload",
			"illumio-core",
			map[string]string{
				"name":                    workload.Name,
				"hostname":                workload.Hostname,
				"description":             workload.Description,
				"service_principal_name":  workload.ServicePrincipalName,
				"public_ip":               workload.PublicIP,
				"service_provider":        workload.ServiceProvider,
				"data_center":             workload.DataCenter,
				"data_center_zone":        workload.DataCenterZone,
				"os_id":                   workload.OsID,
				"os_detail":               workload.OsDetail,
				"enforcement_mode":        workload.EnforcementMode,
				"distinguished_name":      workload.DistinguishedName,
				"external_data_set":       workload.ExternalDataSet,
				"external_data_reference": workload.ExternalDataReference,
			},
			[]string{
				"name",
				"hostname",
				"description",
				"service_principal_name",
				"interfaces",
				"public_ip",
				"service_provider",
				"data_center",
				"data_center_zone",
				"os_id",
				"os_detail",
				"online",
				"labels",
				"enforcement_mode",
				"distinguished_name",
				"external_data_set",
				"external_data_reference",
			},
			map[string]interface{}{
				"interfaces": workload.Interfaces,
				"online":     workload.Online,
				"labels":     workload.Labels,
			},
		))
	}
	return resources
}

func (g *UnmanagedWorkloadGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	// only get unmanaged workloads from the PCE
	workloads, _, err := svc.GetWklds(map[string]string{"managed": "false"})
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc, workloads)
	return nil
}
