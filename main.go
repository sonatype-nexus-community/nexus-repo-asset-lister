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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sonatype-nexus-community/nexus-repo-asset-lister/util"
)

const (
	ENV_NXRM_USERNAME = "NXRM_USERNAME"
	ENV_NXRM_PASSWORD = "NXRM_PASSWORD"

	REPO_TYPE_PROXY = "proxy"
)

var (
	debugLogging    bool   = false
	currentRuntime  string = runtime.GOOS
	commit                 = "unknown"
	outputDirectory string
	nxrmUrl         string
	nxrmUsername    string
	nxrmPassword    string
	version         = "dev"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: nexus-repo-asset-lister [OPTIONS]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&nxrmUrl, "url", "http://localhost:8081", "URL including protocol to your Sonatype Nexus Repository Manager")
	flag.StringVar(&nxrmUsername, "username", "", fmt.Sprintf("Username used to authenticate to Sonatype Nexus Repository (can also be set using the environment variable %s)", ENV_NXRM_USERNAME))
	flag.StringVar(&nxrmPassword, "password", "", fmt.Sprintf("Password used to authenticate to Sonatype Nexus Repository (can also be set using the environment variable %s)", ENV_NXRM_PASSWORD))
	flag.StringVar(&outputDirectory, "o", cwd, "Directory to write asset lists to")
	flag.BoolVar(&debugLogging, "X", false, "Enable debug logging")
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&util.LogFormatter{Module: "NX-AL"})

	flag.Usage = usage
	flag.Parse()

	// Disable Debug Logging if not requested
	if !debugLogging {
		log.SetLevel(log.InfoLevel)
	}

	// Load Credentials
	err := loadCredentials()
	if err != nil {
		os.Exit(1)
	}

	if strings.TrimSpace(nxrmUrl) == "" {
		println("URL to Sonatype Nexus Repository must be supplied")
		os.Exit(1)
	}

	// Output Banner
	println(strings.Repeat("⬢⬡", 42))
	println("")
	println("	███████╗ ██████╗ ███╗   ██╗ █████╗ ████████╗██╗   ██╗██████╗ ███████╗  ")
	println(" 	██╔════╝██╔═══██╗████╗  ██║██╔══██╗╚══██╔══╝╚██╗ ██╔╝██╔══██╗██╔════╝  ")
	println("	███████╗██║   ██║██╔██╗ ██║███████║   ██║    ╚████╔╝ ██████╔╝█████╗    ")
	println(" 	╚════██║██║   ██║██║╚██╗██║██╔══██║   ██║     ╚██╔╝  ██╔═══╝ ██╔══╝    ")
	println(" 	███████║╚██████╔╝██║ ╚████║██║  ██║   ██║      ██║   ██║     ███████╗  ")
	println(" 	╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝   ╚═╝      ╚═╝   ╚═╝     ╚══════╝  ")
	println("")
	println(fmt.Sprintf("	Running on:		%s/%s", currentRuntime, runtime.GOARCH))
	println(fmt.Sprintf("	Version: 		%s (%s)", version, commit))
	println("")
	println(strings.Repeat("⬢⬡", 42))
	println("")
	println("Collecting Assets from:", nxrmUrl)
	println("")

	nxrmServer := NewNxrmServer(nxrmUrl, nxrmUsername, nxrmPassword)
	allRepositories, err := getAllProxyRepositories(nxrmServer)
	if err != nil {
		println(fmt.Sprintf("Error: %v", err))
	}

	for i, r := range *allRepositories {
		if r.Type != nil && *r.Type == REPO_TYPE_PROXY {
			println(fmt.Sprintf("%00d: PROXY of type %s named %s", i, *r.Type, *r.Name))
			componentHashes, err := getAssetsInRepository(nxrmServer, &r)
			if err != nil {
				println(fmt.Sprintf("Error: %v", err))
			}

			println(fmt.Sprintf("  		  %d Compnent Hashes", len(*componentHashes)))

			outputFilename := fmt.Sprintf("%s-%s-%s.json", *r.Type, *r.Format, *r.Name)
			jsonData, err := json.Marshal(componentHashes)
			if err != nil {
				println(fmt.Sprintf("Error: %v", err))
			}

			err = os.WriteFile(path.Join(outputDirectory, strings.ToLower(outputFilename)), jsonData, os.ModePerm)
			if err != nil {
				println(fmt.Sprintf("Failed writing component hashes: %v", err))
			}
		}
	}
}

func loadCredentials() error {
	if strings.TrimSpace(nxrmUsername) == "" {
		log.Debug("Username not supplied as argument - checking environment variable")
		envUsername := os.Getenv(ENV_NXRM_USERNAME)
		if strings.TrimSpace(envUsername) == "" {
			log.Error("No username has been supplied either via argument or environment variable. Cannot continue.")
			return fmt.Errorf("No username has been supplied either via argument or environment variable. Cannot continue.")
		} else {
			nxrmUsername = envUsername
		}
	}

	if strings.TrimSpace(nxrmPassword) == "" {
		log.Debug("Password not supplied as argument - checking environment variable")
		envPassword := os.Getenv(ENV_NXRM_PASSWORD)
		if strings.TrimSpace(envPassword) == "" {
			log.Error("No password has been supplied either via argument or environment variable. Cannot continue.")
			return fmt.Errorf("No password has been supplied either via argument or environment variable. Cannot continue.")
		} else {
			nxrmPassword = envPassword
		}
	}

	return nil
}

func getApiClient() *http.Client {
	return http.DefaultClient
}

func getAllProxyRepositories(server *NxrmServer) (*[]ApiRepository, error) {
	request, err := http.NewRequest(http.MethodGet, server.GetApiUrl("/v1/repositories"), nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(server.username, server.password)

	localVarHTTPResponse, err := getApiClient().Do(request)
	if err != nil {
		return nil, err
	}

	var repositories []ApiRepository
	var localVarBody []byte
	localVarBody, err = io.ReadAll(localVarHTTPResponse.Body)
	if err != nil {
		return nil, err
	}
	localVarHTTPResponse.Body.Close()
	err = json.Unmarshal(localVarBody, &repositories)
	if err != nil {
		return nil, err
	}

	return &repositories, nil
}

func getAssetsInRepository(server *NxrmServer, repository *ApiRepository) (*[]ComponentHash, error) {
	allAssetHashes := make([]ComponentHash, 0)

	firstComponentPage, err := getAssetsPageForRepository(server, *repository.Name, nil)
	if err != nil {
		return nil, err
	}

	for _, c := range firstComponentPage.Items {
		for _, a := range c.Assets {
			allAssetHashes = append(allAssetHashes, ComponentHash(*a.Checksums.Sha1))
		}
	}
	log.Debug("Component Hashes after first page:", len(allAssetHashes))

	lastContinuationToken := firstComponentPage.ContinuationToken

	for lastContinuationToken != nil {
		componentPage, err := getAssetsPageForRepository(server, *repository.Name, lastContinuationToken)
		if err != nil {
			return nil, err
		}

		for _, c := range firstComponentPage.Items {
			for _, a := range c.Assets {
				allAssetHashes = append(allAssetHashes, ComponentHash(*a.Checksums.Sha1))
			}
		}
		log.Debug(fmt.Sprintf("Component Hashes after page: %d - cont token: %s ", len(allAssetHashes), *lastContinuationToken))

		lastContinuationToken = componentPage.ContinuationToken
	}

	return &allAssetHashes, nil
}

func getAssetsPageForRepository(server *NxrmServer, repository_name string, continuation_token *string) (*ApiComponentList, error) {
	apiUrl := server.GetApiUrl(fmt.Sprintf("/v1/components?repository=%s", repository_name))

	if continuation_token != nil {
		apiUrl = fmt.Sprintf("%s&continuationToken=%s", apiUrl, *continuation_token)
	}

	request, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(server.username, server.password)

	localVarHTTPResponse, err := getApiClient().Do(request)
	if err != nil {
		return nil, err
	}

	var componentList ApiComponentList
	var localVarBody []byte
	localVarBody, err = io.ReadAll(localVarHTTPResponse.Body)
	if err != nil {
		return nil, err
	}
	localVarHTTPResponse.Body.Close()
	err = json.Unmarshal(localVarBody, &componentList)
	if err != nil {
		return nil, err
	}

	return &componentList, nil
}
