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
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi"
)

type LabelGroupGenerator struct {
	IllumioService
}

func (g LabelGroupGenerator) createResources(svc *illumioapi.PCE, labelGroups []illumioapi.LabelGroup) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, labelGroup := range labelGroups {
		resourceName := fmt.Sprintf("%s__%s", labelGroup.Key, labelGroup.Name)
		resources = append(resources, terraformutils.NewSimpleResource(
			labelGroup.Href,
			strings.ToLower(resourceName),
			"illumio-core_label_group",
			"illumio-core",
			[]string{},
		))
	}
	return resources
}

func (g *LabelGroupGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	// pass empty params to get all label groups from the PCE
	labelGroups, _, err := svc.GetLabelGroups(map[string]string{}, DRAFT)
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc, labelGroups)
	return nil
}
