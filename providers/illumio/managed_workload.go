// Copyright 2023 The Terraformer Authors.
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
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi/v2"
)

type ManagedWorkloadGenerator struct {
	IllumioService
}

func (g ManagedWorkloadGenerator) createResources(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, workload := range svc.WorkloadsSlice {
		resourceName := fmt.Sprintf("%s__%s", strings.ToLower(*workload.Hostname), stripIdFromHref(workload.Href))
		resources = append(resources, terraformutils.NewResource(
			workload.Href,
			resourceName,
			"illumio-core_managed_workload",
			"illumio-core",
			map[string]string{},
			[]string{},
			map[string]interface{}{
				"ignored_interface_names": *workload.IgnoredInterfaceNames,
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
	_, err = svc.GetWklds(map[string]string{"managed": "true"})
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc)
	return nil
}
