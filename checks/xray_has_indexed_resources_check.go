package checks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/xray"
	"github.com/rdar-lab/JCheck/common"
	"net/http"
)

func GetXrayHasIndexedResourcesCheck() *common.CheckDef {
	return &common.CheckDef{
		Name:        "XrayHasIndexedResourcesCheck",
		Group:       "Xray",
		Description: "Performs a check that validates that XRAY has configured indexed resources",
		IsReadOnly:  true,
		CheckFunc: func(c context.Context) (string, error) {

			serverConf, err := config.GetDefaultServerConf()
			if err != nil {
				return "", err
			}

			if serverConf.XrayUrl == "" {
				return "", errors.New("xray service is not configured")
			}

			xrayServiceMgr, err := commands.CreateXrayServiceManager(serverConf)
			if err != nil {
				return "", err
			}

			httpClientsDetails := xrayServiceMgr.Config().GetServiceDetails().CreateHttpClientDetails()

			msg, err := checkIndexedRepos(serverConf, err, xrayServiceMgr, httpClientsDetails)
			if msg != "" || err != nil {
				return msg, err
			}

			return checkIndexedBuilds(serverConf, err, xrayServiceMgr, httpClientsDetails)
		},
	}
}

func checkIndexedBuilds(serverConf *config.ServerDetails, err error, xrayServiceMgr *xray.XrayServicesManager, httpClientsDetails httputils.HttpClientDetails) (string, error) {
	indexedBuildsUrl := clientutils.AddTrailingSlashIfNeeded(serverConf.XrayUrl) + "api/v1/binMgr/default/builds"
	resp, body, _, err := xrayServiceMgr.Client().SendGet(indexedBuildsUrl, true, &httpClientsDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}

	indexedBuildsResponse :=
		struct {
			IndexedBuilds    []struct{} `json:"indexed_builds,omitempty"`
			NonIndexedBuilds []struct{} `json:"non_indexed_builds,omitempty"`
		}{
			IndexedBuilds:    make([]struct{}, 0),
			NonIndexedBuilds: make([]struct{}, 0),
		}

	err = json.Unmarshal(body, &indexedBuildsResponse)

	if err != nil {
		return "", errors.New("failed unmarshalling indexed build response")
	}

	if len(indexedBuildsResponse.IndexedBuilds) == 0 {
		return "", errors.New("detected no indexed resources")
	} else {
		return fmt.Sprintf("detected %d indexed builds", len(indexedBuildsResponse.IndexedBuilds)), nil
	}
}

func checkIndexedRepos(serverConf *config.ServerDetails, err error, xrayServiceMgr *xray.XrayServicesManager, httpClientsDetails httputils.HttpClientDetails) (string, error) {
	indexedReposUrl := clientutils.AddTrailingSlashIfNeeded(serverConf.XrayUrl) + "api/v1/binMgr/default/repos"

	resp, body, _, err := xrayServiceMgr.Client().SendGet(indexedReposUrl, true, &httpClientsDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return "", err
	}

	indexedReposResponse :=
		struct {
			IndexedRepos    []struct{} `json:"indexed_repos,omitempty"`
			NonIndexedRepos []struct{} `json:"non_indexed_repos,omitempty"`
		}{
			IndexedRepos:    make([]struct{}, 0),
			NonIndexedRepos: make([]struct{}, 0),
		}

	err = json.Unmarshal(body, &indexedReposResponse)

	if err != nil {
		return "", errors.New("failed unmarshalling indexed repos response")
	}

	if len(indexedReposResponse.NonIndexedRepos) == 0 && len(indexedReposResponse.IndexedRepos) == 0 {
		return "", errors.New("detected no repositories")
	} else if len(indexedReposResponse.IndexedRepos) > 0 {
		return fmt.Sprintf("detected %d indexed repositories", len(indexedReposResponse.IndexedRepos)), nil
	}
	return "", nil
}
