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
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi/v2"
)

type ServiceGenerator struct {
	IllumioService
}

func (g ServiceGenerator) createResources(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, service := range svc.ServicesSlice {
		resources = append(resources, terraformutils.NewSimpleResource(
			service.Href,
			strings.ToLower(service.Name),
			"illumio-core_service",
			"illumio-core",
			[]string{},
		))
	}
	return resources
}

func (g *ServiceGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	// pass empty params to get all IP lists from the PCE
	_, err = svc.GetServices(map[string]string{}, DRAFT)
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc)
	return nil
}
