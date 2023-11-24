// Copyright 2022 The Terraformer Authors.
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
	"log"
	"net/url"
	"strconv"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/brian1917/illumioapi/v2"
)

type IllumioService struct { //nolint
	terraformutils.Service
}

func (s *IllumioService) generateService() (*illumioapi.PCE, error) {
	u, err := url.Parse(s.Args["pce_host"].(string))
	if err != nil {
		return nil, err
	}

	portStr := u.Port()
	if portStr == "" {
		log.Println("No port provided in PCE host string, using default port 443")
		portStr = "443"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &illumioapi.PCE{
		FQDN: u.Hostname(),
		Port: port,
		Org:  s.Args["org_id"].(int),
		User: s.Args["api_username"].(string),
		Key:  s.Args["api_secret"].(string),
	}, nil
}
