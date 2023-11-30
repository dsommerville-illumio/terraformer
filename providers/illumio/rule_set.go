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

type RuleSetGenerator struct {
	IllumioService
}

func (g *RuleSetGenerator) loadRules(svc *illumioapi.PCE, rulesMap map[string][]illumioapi.Rule) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for ruleSetName, rules := range rulesMap {
		for _, rule := range rules {
			resourceName := fmt.Sprintf("%s__%s", ruleSetName, stripIdFromHref(rule.Href))
			resources = append(resources, terraformutils.NewSimpleResource(
				rule.Href,
				strings.ToLower(resourceName),
				"illumio-core_security_rule",
				"illumio-core",
				[]string{},
			))
		}
	}
	return resources
}

func (g *RuleSetGenerator) loadRuleSets(svc *illumioapi.PCE) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, ruleSet := range svc.RuleSetsSlice {
		resources = append(resources, terraformutils.NewSimpleResource(
			ruleSet.Href,
			strings.ToLower(ruleSet.Name),
			"illumio-core_rule_set",
			"illumio-core",
			[]string{},
		))
	}
	return resources
}

func (g *RuleSetGenerator) InitResources() error {
	svc, err := g.generateService()
	if err != nil {
		return err
	}

	_, err = svc.GetRulesets(map[string]string{}, DRAFT)
	if err != nil {
		return err
	}

	g.Resources = append(g.Resources, g.loadRuleSets(svc)...)

	rules := map[string][]illumioapi.Rule{}
	for _, ruleSet := range svc.RuleSetsSlice {
		rules[ruleSet.Name] = *ruleSet.Rules
	}

	g.Resources = append(g.Resources, g.loadRules(svc, rules)...)

	return nil
}
