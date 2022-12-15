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
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi"
)

type ContainerClusterGenerator struct {
	IllumioService
}

func (g ContainerClusterGenerator) loadWorkloadProfiles(
	svc *illumioapi.PCE,
	workloadProfilesMap map[string][]illumioapi.ContainerWorkloadProfile,
) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for clusterName, profiles := range workloadProfilesMap {
		for _, ccwp := range profiles {
			resourceName := fmt.Sprintf("%s__%s", clusterName, ccwp.Namespace)
			resources = append(resources, terraformutils.NewResource(
				ccwp.Href,
				strings.ToLower(resourceName),
				"illumio-core_container_cluster_workload_profile",
				"illumio-core",
				map[string]string{},
				[]string{},
				map[string]interface{}{
					"managed":       ccwp.Managed,
					"assign_labels": ccwp.AssignLabels,
					"labels":        ccwp.Labels,
				},
			))
		}
	}
	return resources
}

func (g ContainerClusterGenerator) loadContainerClusters(svc *illumioapi.PCE, containerClusters []illumioapi.ContainerCluster) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, cc := range containerClusters {
		resourceName := fmt.Sprintf("%s__%s", strings.ToLower(cc.Name), stripIdFromHref(cc.Href))
		resources = append(resources, terraformutils.NewResource(
			cc.Href,
			resourceName,
			"illumio-core_container_cluster",
			"illumio-core",
			map[string]string{},
			[]string{},
			map[string]interface{}{
				"online": cc.Online,
			},
		))
	}
	return resources
}

func (g *ContainerClusterGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}

	containerClusters, _, err := svc.GetContainerClusters(map[string]string{})
	if err != nil {
		return err
	}
	g.Resources = append(g.Resources, g.loadContainerClusters(svc, containerClusters)...)

	workloadProfiles := map[string][]illumioapi.ContainerWorkloadProfile{}
	for _, cc := range containerClusters {
		clusterID := stripIdFromHref(cc.Href)
		profiles, _, err := svc.GetContainerWkldProfiles(map[string]string{}, clusterID)
		if err != nil {
			log.Printf("Failed to fetch workload profiles for container cluster with ID %q", clusterID)
		}
		workloadProfiles[cc.Name] = profiles
	}
	if err != nil {
		return err
	}

	g.Resources = append(g.Resources, g.loadWorkloadProfiles(svc, workloadProfiles)...)

	return nil
}
