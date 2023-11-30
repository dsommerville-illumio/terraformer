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

type EnforcementBoundaryGenerator struct {
	IllumioService
}

func (g *EnforcementBoundaryGenerator) createResources(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, eb := range svc.EnforcementBoundariesSlice {
		resources = append(resources, terraformutils.NewSimpleResource(
			eb.Href,
			strings.ToLower(eb.Name),
			"illumio-core_enforcement_boundary",
			"illumio-core",
			[]string{},
		))
	}
	return resources
}

func (g *EnforcementBoundaryGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	_, err = svc.GetEnforcementBoundaries(map[string]string{}, DRAFT)
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc)
	return nil
}
