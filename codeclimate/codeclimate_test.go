package codeclimate

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/eorate/shadowheart/mocks"

	"github.com/stretchr/testify/assert"
)

func init() {
	Client = &mocks.MockClient{}

}

var testRepoID = "7d5f95bdf79e483f8d1f509f"
var repoTestData = "{\"data\":[{\"id\":\"20e576a84fa645a399a29ce2\"," +
	"\"type\":\"ref_points\",\"attributes\":{\"analyzed\":true," +
	"\"branch\":\"develop\"," +
	"\"commit_sha\":\"a0b57436178e48b9a30f62650843b7a9e9cbf870\"," +
	"\"created_at\":\"2020-09-10T14:38:49.333Z\"," +
	"\"ref\":\"refs/heads/develop\"},\"relationships\":{\"snapshot\"" +
	":{\"data\":{\"id\":\"32363380cbea4ac8b1642df9\"," +
	"\"type\":\"snapshots\"}}}}],\"links\":" +
	"{\"self\":\"https://api.codeclimate.com/v1/repos/" +
	"7d5f95bdf79e483f8d1f509f/ref_points?filter%5Banalyzed%5D=true" +
	"\\u0026filter%5Bbranch%5D=develop\\u0026page%5Bnumber%5D=1" +
	"\\u0026page%5Bsize%5D=1\",\"next\":\"https://api.codeclimate.com/" +
	"v1/repos/7d5f95bdf79e483f8d1f509f/ref_points?filter%5Banalyzed%5D=true" +
	"\\u0026filter%5Bbranch%5D=develop\\u0026page%5Bnumber%5D=2\\u0026" +
	"page%5Bsize%5D=1\",\"last\":\"https://api.codeclimate.com/v1/repos/" +
	"7d5f95bdf79e483f8d1f509f/ref_points?filter%5Banalyzed%5D=true" +
	"\\u0026filter%5Bbranch%5D=develop\\u0026page%5Bnumber%5D=65" +
	"\\u0026page%5Bsize%5D=1\"}}"

var metricsTestData = "{\"data\":{\"id\":\"32363380cbea4ac8b1642df9\"," +
	"\"type\":\"snapshots\",\"attributes\":" +
	"{\"commit_sha\":\"a0b57436178e48b9a30f62650843b7a9e9cbf870\"," +
	"\"committed_at\":\"2020-09-10T14:38:49.333Z\"," +
	"\"created_at\":\"2020-10-19T09:52:22.337Z\",\"lines_of_code\":1117," +
	"\"ratings\":[{\"path\":\"/\",\"letter\":\"A\",\"measure\":" +
	"{\"value\":0.0,\"unit\":\"percent\",\"meta\":{\"remediation_time\":" +
	"{\"value\":0.0,\"unit\":\"minute\"},\"implementation_time\":" +
	"{\"value\":15527.071644135749,\"unit\":\"minute\"}}}," +
	"\"pillar\":\"Maintainability\"}],\"gpa\":null," +
	"\"worker_version\":61100},\"meta\":{\"issues_count\":24," +
	"\"measures\":{\"remediation\":{\"value\":8.1,\"unit\":\"minute\"}," +
	"\"technical_debt_ratio\":{\"value\":3.7,\"unit\":\"percent\"," +
	"\"meta\":{\"remediation_time\":{\"value\":0.0,\"unit\":\"minute\"}," +
	"\"implementation_time\":{\"value\":15527.071644135749," +
	"\"unit\":\"minute\"}}}}}}}"

func TestGetSnapshotID(t *testing.T) {

	result := getSnapshotID([]byte(repoTestData))

	assert.Equal(t, "32363380cbea4ac8b1642df9", result)
}

func TestAddParams(t *testing.T) {
	url := codeClimateURL + "/repos/" + testRepoID + "/ref_points"

	params := make(map[string]string)
	params["filter[branch]"] = "develop"
	params["filter[analyzed]"] = "true"
	params["page[size]"] = "1"

	urlWithParams := addParams(url, params)
	expectedURLWithParams := "https://api.codeclimate.com/v1/repos/" +
		"7d5f95bdf79e483f8d1f509f/ref_points?filter%5Banalyzed%5D=true&" +
		"filter%5Bbranch%5D=develop&page%5Bsize%5D=1"

	assert.Equal(t, urlWithParams, expectedURLWithParams)
}

func TestMakeRequest(t *testing.T) {
	url := codeClimateURL + "/repos/" + testRepoID + "/ref_points"

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(metricsTestData)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	statusCode, body := makeRequest(url)

	assert.Equal(t, statusCode, 200)
	assert.Equal(t, string(body), metricsTestData)
}

func TestGetRepositoryStats(t *testing.T) {

	stats := map[string]string{
		"Maintainability":      "A",
		"Remediation":          "8.1 minute",
		"Technical Debt Ratio": "3.7 percent"}
	result := getRepositoryStats([]byte(metricsTestData))
	assert.Equal(t, result, stats)
}
