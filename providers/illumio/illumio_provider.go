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
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type IllumioProvider struct { //nolint
	terraformutils.Provider
	pceHost     string
	pceOrgId    int
	apiUsername string
	apiSecret   string
}

func (p *IllumioProvider) Init(args []string) error {
	p.pceHost = args[0]
	if args[1] != "" {
		orgId, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid org ID value %v, must be integer value", args[2]))
		}
		p.pceOrgId = orgId
	}
	p.apiUsername = args[2]
	p.apiSecret = args[3]

	if p.pceHost == "" {
		if os.Getenv("ILLUMIO_PCE_HOST") == "" {
			return errors.New("Missing PCE host, set --host or ILLUMIO_PCE_HOST env var")
		}
		p.pceHost = os.Getenv("ILLUMIO_PCE_HOST")
	}

	if p.pceOrgId == 0 {
		if os.Getenv("ILLUMIO_PCE_ORG_ID") == "" {
			log.Println("No value set for ILLUMIO_PCE_ORG_ID, using default org ID 1")
			p.pceOrgId = 1
		} else {
			orgId, err := strconv.Atoi(os.Getenv("ILLUMIO_PCE_ORG_ID"))
			if err != nil {
				return err
			}
			p.pceOrgId = orgId
		}
	}

	if p.apiUsername == "" {
		if os.Getenv("ILLUMIO_API_KEY_USERNAME") == "" {
			return errors.New("Missing API key, set --api-key or ILLUMIO_API_KEY_USERNAME env var")
		}
		p.apiUsername = os.Getenv("ILLUMIO_API_KEY_USERNAME")
	}

	if p.apiSecret == "" {
		if os.Getenv("ILLUMIO_API_KEY_SECRET") == "" {
			return errors.New("Missing API secret, set --api-secret or ILLUMIO_API_KEY_SECRET env var")
		}
		p.apiSecret = os.Getenv("ILLUMIO_API_KEY_SECRET")
	}

	return nil
}

func (p *IllumioProvider) GetName() string {
	return "illumio-core"
}

func (p *IllumioProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			p.GetName(): map[string]interface{}{},
		},
	}
}

func (IllumioProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{
		"container_cluster_workload_profile": {
			"container_cluster": []string{"container_cluster_href", "id"},
			"label":             []string{"labels.href", "id", "assign_labels.href", "id"},
		},
		"enforcement_boundary": {
			"service":     []string{"ingress_services.href", "id"},
			"ip_list":     []string{"consumers.ip_list.href", "id", "providers.ip_list.href", "id"},
			"label":       []string{"consumers.label.href", "id", "providers.label.href", "id"},
			"label_group": []string{"consumers.label_group.href", "id", "providers.label_group.href", "id"},
		},
		"label_group": {
			"label":       []string{"labels.href", "id"},
			"label_group": []string{"sub_groups.href", "id"},
		},
		"managed_workload": {
			"label": []string{"labels.href", "id"},
		},
		"pairing_profile": {
			"label": []string{"labels.href", "id"},
		},
		"rule_set": {
			"label":       []string{"scopes.label.href", "id"},
			"label_group": []string{"scopes.label_group.href", "id"},
		},
		"security_rule": {
			"rule_set":           []string{"rule_set_href", "id"},
			"service":            []string{"ingress_services.href", "id"},
			"ip_list":            []string{"consumers.ip_list.href", "id", "providers.ip_list.href", "id"},
			"label":              []string{"consumers.label.href", "id", "providers.label.href", "id"},
			"label_group":        []string{"consumers.label_group.href", "id", "providers.label_group.href", "id"},
			"virtual_server":     []string{"providers.virtual_server.href", "id"},
			"virtual_service":    []string{"consumers.virtual_service.href", "id", "providers.virtual_service.href", "id"},
			"managed_workload":   []string{"consumers.workload.href", "id", "providers.workload.href", "id"},
			"unmanaged_workload": []string{"consumers.workload.href", "id", "providers.workload.href", "id"},
		},
		"service_binding": {
			"virtual_service":    []string{"virtual_service.href", "id"},
			"managed_workload":   []string{"workload.href", "id"},
			"unmanaged_workload": []string{"workload.href", "id"},
		},
		"unmanaged_workload": {
			"label": []string{"labels.href", "id"},
		},
		"virtual_service": {
			"label":   []string{"labels.href", "id"},
			"service": []string{"service.href", "id"},
		},
	}
}

func (p *IllumioProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {
	return map[string]terraformutils.ServiceGenerator{
		"container_cluster":    &ContainerClusterGenerator{},
		"enforcement_boundary": &EnforcementBoundaryGenerator{},
		"ip_list":              &IPListGenerator{},
		"label":                &LabelGenerator{},
		"label_group":          &LabelGroupGenerator{},
		"label_type":           &LabelTypeGenerator{},
		"managed_workload":     &ManagedWorkloadGenerator{},
		"pairing_profile":      &PairingProfileGenerator{},
		"rule_set":             &RuleSetGenerator{},
		"service":              &ServiceGenerator{},
		"unmanaged_workload":   &UnmanagedWorkloadGenerator{},
		"virtual_services":     &VirtualServiceGenerator{},
	}
}

func (p *IllumioProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New("illumio: " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"pce_host":     p.pceHost,
		"org_id":       p.pceOrgId,
		"api_username": p.apiUsername,
		"api_secret":   p.apiSecret,
	})
	return nil
}
