package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type ScoreServerClient struct {
	client *http.Client
	url    string
}

func NewScoreServerClient(url string) (*ScoreServerClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Jar: jar}
	return &ScoreServerClient{
		client: client,
		url:    url,
	}, nil
}

func (c *ScoreServerClient) Authenticate(name, password string) error {
	req := ScoreServerAuthRequest{
		Name:     name,
		Password: password,
	}
	j, _ := json.Marshal(req)
	reqBuf := bytes.NewBuffer(j)

	res, err := c.client.Post(c.url+"/api/sessions", "application/json", reqBuf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (c *ScoreServerClient) GetPenalties(after time.Time) (*GetPenaltiesResponse, error) {
	req := QueryRequest{
		Query: fmt.Sprintf("{  penalties(after: \"%s\"){ id problem { code body { title } } team { number } createdAt } }",
			after.Format(time.RFC3339Nano),
		),
	}
	j, _ := json.Marshal(req)
	res, err := c.client.Post(c.url+"/api/graphql?=", "application/json", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to query `%s`", req.Query)
	}

	aBody := new(bytes.Buffer)
	eBody := io.TeeReader(res.Body, aBody)
	var er ErrorsResponse
	if err := json.NewDecoder(eBody).Decode(&er); err != nil {
		return nil, err
	}

	if len(er.Errors) > 0 {
		for _, e := range er.Errors {
			if e.Extensions.Code == "unauthorized" {
				return nil, fmt.Errorf("unauthorized")
			}
		}
		return nil, fmt.Errorf("unknown errors")
	}

	var r GetPenaltiesResponse
	if err := json.NewDecoder(aBody).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
