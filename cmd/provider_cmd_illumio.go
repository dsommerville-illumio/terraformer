// Copyright 2022 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	illumio_terraforming "github.com/GoogleCloudPlatform/terraformer/providers/illumio"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/spf13/cobra"
)

func newCmdIllumioImporter(options ImportOptions) *cobra.Command {
	var pceHost string
	var pceOrgId string
	var pceAPIKey string
	var pceAPISecret string

	cmd := &cobra.Command{
		Use:   "illumio",
		Short: "Import current state to Terraform configuration from the Illumio PCE",
		Long:  "Import current state to Terraform configuration from the Illumio PCE",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := newIllumioProvider()
			err := Import(provider, options, []string{pceHost, pceOrgId, pceAPIKey, pceAPISecret})
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.AddCommand(listCmd(newIllumioProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "label,ip_list", "label=href1:href2")
	cmd.PersistentFlags().StringVarP(&pceHost, "host", "H", "", "host=https://my.pce.com:8443")
	cmd.PersistentFlags().StringVarP(&pceOrgId, "org-id", "", "", "org-id=123")
	cmd.PersistentFlags().StringVarP(&pceAPIKey, "api-key", "", "", "api-key=<PCE_API_KEY>")
	cmd.PersistentFlags().StringVarP(&pceAPISecret, "api-secret", "", "", "api-secret=<PCE_API_SECRET>")
	return cmd
}

func newIllumioProvider() terraformutils.ProviderGenerator {
	return &illumio_terraforming.IllumioProvider{}
}
