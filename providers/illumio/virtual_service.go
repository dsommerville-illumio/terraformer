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
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi/v2"
)

type VirtualServiceGenerator struct {
	IllumioService
}

func (g *VirtualServiceGenerator) loadServiceBindings(svc *illumioapi.PCE, bindingsMap map[string][]illumioapi.ServiceBinding) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for virtualServiceName, bindings := range bindingsMap {
		for _, binding := range bindings {
			resourceName := fmt.Sprintf("%s__%s", virtualServiceName, stripIdFromHref(binding.Href))
			resources = append(resources, terraformutils.NewSimpleResource(
				binding.Href,
				strings.ToLower(resourceName),
				"illumio-core_service_binding",
				"illumio-core",
				[]string{},
			))
		}
	}
	return resources
}

func (g *VirtualServiceGenerator) loadVirtualServices(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, vs := range svc.VirtualServicesSlice {
		resources = append(resources, terraformutils.NewSimpleResource(
			vs.Href,
			strings.ToLower(vs.Name),
			"illumio-core_virtual_service",
			"illumio-core",
			[]string{},
		))
	}
	return resources
}

func (g *VirtualServiceGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	_, err = svc.GetVirtualServices(map[string]string{}, DRAFT)
	if err != nil {
		return err
	}

	g.Resources = append(g.Resources, g.loadVirtualServices(svc)...)

	bindingsMap := map[string][]illumioapi.ServiceBinding{}
	for _, virtualService := range svc.VirtualServicesSlice {
		bindings, _, err := svc.GetServiceBindings(
			map[string]string{"virtual_service": virtualService.Href},
		)
		if err != nil {
			log.Printf("Failed to fetch service bindings for virtual service %q", virtualService.Href)
		}
		bindingsMap[virtualService.Name] = bindings
	}

	g.Resources = append(g.Resources, g.loadServiceBindings(svc, bindingsMap)...)

	return nil
}
