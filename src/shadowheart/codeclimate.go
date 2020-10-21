package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// https://stackoverflow.com/questions/21830447/json-cannot-unmarshal-object-into-go-value-of-type
type repoDetails struct {
	Repository []repository `json:"data"`
}

type repository struct {
	Relationship relationship `json:"relationships"`
}

type relationship struct {
	Snapshot snapshot `json:"snapshot"`
}

type snapshot struct {
	SnapshotData snapshotData `json:"data"`
}

type snapshotData struct {
	ID string `json:"id"`
}

// For building the metrics
type metricsData struct {
	Metrics metrics `json:"data"`
}

type metrics struct {
	Attributes attributes `json:"attributes"`
	Meta       meta       `json:"meta"`
}

type attributes struct {
	Ratings []ratings `json:"ratings"`
}

type ratings struct {
	Letter string `json:"letter"`
}

type meta struct {
	Measures measure `json:"measures"`
}

type measure struct {
	Remediation        units `json:"remediation"`
	TechnicalDebtRatio units `json:"technical_debt_ratio"`
}

type units struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

var repos repoDetails
var mets metricsData

var repoID = os.Getenv("CODECLIMATE_REPO_ID")
var codeClimateURL = os.Getenv("CODECLIMATE_URL")

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	// Client HTTPClient
	Client HTTPClient
)

func init() {
	Client = &http.Client{
		Timeout: time.Second * 10, // Timeout after 10 seconds
	}
}

func addParams(repoURL string, paramArgs map[string]string) string {
	baseURL, err := url.Parse(repoURL)
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	for param, arg := range paramArgs {
		params.Set(param, arg)
	}

	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}

func makeRequest(baseURL string) (int, []byte) {
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("CODECLIMATE_TOKEN")
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Authorization", "Token token="+token)

	res, getErr := Client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return res.StatusCode, body
}

func getSnapshotID(data []byte) string {
	jsonErr := json.Unmarshal(data, &repos)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return repos.Repository[0].Relationship.Snapshot.SnapshotData.ID
}

func getRepositoryStats(data []byte) map[string]string {

	jsonErr := json.Unmarshal(data, &mets)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	repoStats := make(map[string]string)
	remediationValue := mets.Metrics.Meta.Measures.Remediation.Value
	remediationUnit := mets.Metrics.Meta.Measures.Remediation.Unit
	repoStats["Remediation"] = strconv.FormatFloat(float64(remediationValue), 'f', 1, 32) +
		" " + remediationUnit

	techDebtRatioValue := mets.Metrics.Meta.Measures.TechnicalDebtRatio.Value
	techDebtRatioUnit := mets.Metrics.Meta.Measures.TechnicalDebtRatio.Unit
	repoStats["Technical Debt Ratio"] = strconv.FormatFloat(float64(techDebtRatioValue), 'f', 1, 32) +
		" " + techDebtRatioUnit
	repoStats["Maintainability"] = mets.Metrics.Attributes.Ratings[0].Letter
	return repoStats

}

// BuildRepositoryStats is a function that gives us the metrics
func BuildRepositoryStats() map[string]string {
	// Get the snapshot ID
	url := codeClimateURL + "/repos/" + repoID + "/ref_points"

	params := make(map[string]string)
	params["filter[branch]"] = "develop"
	params["filter[analyzed]"] = "true"
	params["page[size]"] = "1"

	urlWithParams := addParams(url, params)
	_, body := makeRequest(urlWithParams)
	snapshotID := getSnapshotID(body)

	// Get the repo details for the given snapshot
	url = codeClimateURL + "/repos/" + repoID + "/snapshots/" + snapshotID
	_, body = makeRequest(url)
	stats := getRepositoryStats(body)
	return stats
}
