package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "log"
	"net/http"
	"net/url"
)

const token = "wBTvk4rH5auTh4up8nOaVCcJBYWT3jr2Wk7QnlcOc4Qa7dvkgaDBGK6pP3hUaneP_aw0vGveE3XqDEfXSBIsQy7slH24kQ_SZVlojNYkNrg"

//EngEnv is the Engage environment.
type EngEnv struct {
	Host  string
	Token string
}

//Metrics contains the measurable stsuff in Engage.
type Metrics struct {
	RateLimit                      int32  `json:"rateLimit"`
	MaxBatchSize                   int32  `json:"maxBatchSize"`
	CurrentRateLimit               int32  `json:"currentRateLimit"`
	SupporterRead                  int32  `json:"supporterRead"`
	TotalAPICalls                  int32  `json:"totalAPICalls"`
	LastAPICall                    string `json:"lastAPICall"`
	TotalAPICallFailures           int32  `json:"totalAPICallFailures"`
	LastAPICallFailure             string `json:"lastAPICallFailure"`
	SupporterReads                 int32  `json:"supporterRead"`
	SupporterAdd                   int32  `json:"supporterAdd"`
	SupporterUpdate                int32  `json:"supporterUpdate"`
	SupporterDelete                int32  `json:"supporterDelete"`
	ActivityEvent                  int32  `json:"activityEvent"`
	ActivitySubscribe              int32  `json:"activitySubscribe"`
	ActivityFundraise              int32  `json:"activityFundraise"`
	ActivityTargetedLetter         int32  `json:"activityTargetedLetter"`
	ActivityPetition               int32  `json:"activityPetition"`
	ActivitySubscriptionManagement int32  `json:"activitySubscriptionManagement"`
}

//MetricReturn is returned by Engage when asking for metrics.
type MetricReturn struct {
	ID        string
	Timestamp string
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `jsin:"serverId"`
	}
	Payload Metrics
}

//Measure reads metrics and returns them.
func (e EngEnv) Measure() (*Metrics, error) {
	u, _ := url.Parse("/api/integration/ext/v1/metrics")
	x := fmt.Sprintf("https://%v", e.Host)
	b, _ := url.Parse(x)
	t := b.ResolveReference(u)
	fmt.Printf("Meterics URL is %v", t)
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, t.String(), nil)
	req.Header.Set("authToken", e.Token)
	var body []byte
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Printf("body: %v\n", string(body))
	var m MetricReturn
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MetricReturn is: %+v\n", m)
	return &m.Payload, err
}

func main() {
	e := EngEnv{"hq.uat.igniteaction.net", token}
	fmt.Printf("EngEnv is %+v\n", e)
	m, err := e.Measure()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Metrics: %+v\n", m)
}
