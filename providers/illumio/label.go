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

type LabelGenerator struct {
	IllumioService
}

func (g LabelGenerator) createResources(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, label := range svc.LabelsSlice {
		resourceName := fmt.Sprintf("%s__%s", label.Key, label.Value)
		resources = append(resources, terraformutils.NewSimpleResource(
			label.Href,
			strings.ToLower(resourceName),
			"illumio-core_label",
			"illumio-core",
			[]string{},
		))
	}
	return resources
}

func (g *LabelGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}
	// pass empty params to get all labels from the PCE
	_, err = svc.GetLabels(map[string]string{})
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc)
	return nil
}
