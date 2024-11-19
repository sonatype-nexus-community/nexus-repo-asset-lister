/**
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"strings"
)

type NxrmServer struct {
	baseUrl  string
	username string
	password string
}

func NewNxrmServer(url string, username string, password string) *NxrmServer {
	url = strings.TrimRight(url, "/")
	server := &NxrmServer{
		baseUrl:  url,
		username: username,
		password: password,
	}
	return server
}

func (s *NxrmServer) GetApiUrl(api_path string) string {
	return fmt.Sprintf("%s/service/rest%s", s.baseUrl, api_path)
}

type ComponentHash string

type ApiRepository struct {
	Name       *string                `json:"name"`
	Format     *string                `json:"format"`
	Type       *string                `json:"type"`
	Url        *string                `json:"url"`
	Size       *int                   `json:"size"`
	Attributes map[string]interface{} `json:"attributes"`
}

type ApiComponentAssetChecksums struct {
	Md5    *string `json:"md5"`
	Sha1   *string `json:"sha1"`
	Sha256 *string `json:"sha256"`
}

type ApiComponentAsset struct {
	Id         string                      `json:"Id"`
	Repository string                      `json:"repository"`
	Format     string                      `json:"format"`
	Checksums  *ApiComponentAssetChecksums `json:"checksum"`
}

type ApiComponent struct {
	Id         string              `json:"id"`
	Repository string              `json:"repository"`
	Format     string              `json:"format"`
	Group      *string             `json:"group"`
	Name       string              `json:"name"`
	Version    *string             `json:"version"`
	Assets     []ApiComponentAsset `json:"assets"`
}

type ApiComponentList struct {
	Items             []ApiComponent `json:"items"`
	ContinuationToken *string        `json:"continuationToken"`
}
