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

type WorkloadGenerator struct {
	IllumioService
}

func (g *WorkloadGenerator) loadWorkloads(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, workload := range svc.WorkloadsSlice {
		if workload.VEN == nil {
			resources = append(resources, g.loadWorkload(workload, "illumio-core_unmanaged_workload"))
		} else {
			resources = append(resources, g.loadWorkload(workload, "illumio-core_managed_workload"))
		}
	}
	return resources
}

func (g *WorkloadGenerator) loadWorkload(workload illumioapi.Workload, resType string) terraformutils.Resource {
	resourceName := fmt.Sprintf("%s__%s", strings.ToLower(dereference(workload.Hostname)), stripIdFromHref(workload.Href))
	return terraformutils.NewSimpleResource(
		workload.Href,
		resourceName,
		resType,
		"illumio-core",
		[]string{},
	)
}

func (g *WorkloadGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	_, err = svc.GetWklds(map[string]string{})
	if err != nil {
		return err
	}
	g.Resources = append(g.Resources, g.loadWorkloads(svc)...)
	return nil
}
