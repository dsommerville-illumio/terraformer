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

type ManagedWorkloadGenerator struct {
	IllumioService
}

func (g ManagedWorkloadGenerator) createResources(svc *illumioapi.PCE, workloads []illumioapi.Workload) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, workload := range workloads {
		resourceName := stripIdFromHref(workload.Href)
		resources = append(resources, terraformutils.NewResource(
			workload.Href,
			resourceName,
			"illumio-core_managed_workload",
			"illumio-core",
			map[string]string{
				"name":                    workload.Name,
				"description":             workload.Description,
				"service_principal_name":  workload.ServicePrincipalName,
				"service_provider":        workload.ServiceProvider,
				"data_center":             workload.DataCenter,
				"data_center_zone":        workload.DataCenterZone,
				"enforcement_mode":        workload.EnforcementMode,
				"external_data_set":       workload.ExternalDataSet,
				"external_data_reference": workload.ExternalDataReference,
			},
			[]string{
				"name",
				"description",
				"service_principal_name",
				"service_provider",
				"data_center",
				"data_center_zone",
				"enforcement_mode",
				"labels",
				"ignored_interface_names",
				"external_data_set",
				"external_data_reference",
			},
			map[string]interface{}{
				"labels":                  workload.Labels,
				"ignored_interface_names": workload.IgnoredInterfaceNames,
			},
		))
	}
	return resources
}

func (g *ManagedWorkloadGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	// only get managed workloads from the PCE
	workloads, _, err := svc.GetWklds(map[string]string{"managed": "true"})
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc, workloads)
	return nil
}
