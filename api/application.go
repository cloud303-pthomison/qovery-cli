package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Applications struct {
	Results []Application `json:"results"`
}

type Application struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Status         Status `json:"status"`
	ConnectionURI  string `json:"connection_uri"`
	TotalDatabases *int   `json:"total_databases"`
	TotalBrokers   *int   `json:"total_brokers"`
	TotalStorage   *int   `json:"total_storage"`
}

func GetApplicationByName(projectId string, branchName string, name string) Application {
	for _, a := range ListApplications(projectId, branchName).Results {
		if a.Name == name {
			return a
		}
	}

	return Application{}
}

func ListApplications(projectId string, branchName string) Applications {
	apps := Applications{}

	if projectId == "" || branchName == "" {
		return apps
	}

	CheckAuthenticationOrQuitWithMessage()

	req, _ := http.NewRequest(http.MethodGet, RootURL+"/project/"+projectId+"/branch/"+branchName+"/application", nil)
	req.Header.Set(headerAuthorization, headerValueBearer+GetAuthorizationToken())

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return apps
	}

	err = CheckHTTPResponse(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &apps)

	return apps
}
