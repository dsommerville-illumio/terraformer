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

type LabelGenerator struct {
	IllumioService
}

func (g LabelGenerator) createResources(svc *illumioapi.PCE, labels []illumioapi.Label) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, label := range labels {
		resourceName := fmt.Sprintf("%s_%s", label.Key, strings.ToLower(label.Value))
		resources = append(resources, terraformutils.NewResource(
			label.Href,
			resourceName,
			"illumio-core_label",
			"illumio-core",
			map[string]string{
				"key":                     label.Key,
				"value":                   label.Value,
				"external_data_set":       label.ExternalDataSet,
				"external_data_reference": label.ExternalDataReference,
			},
			[]string{"external_data_set", "external_data_reference"},
			map[string]interface{}{},
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
	labels, _, err := svc.GetLabels(map[string]string{})
	if err != nil {
		return err
	}
	g.Resources = g.createResources(svc, labels)
	return nil
}
