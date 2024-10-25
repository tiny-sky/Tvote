package test

import (
	"fmt"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/tiny-sky/Tvote/cmd/http_request"
)

func TestVote(t *testing.T) {
	t.Parallel()
	var err error
	var respBody []byte

	respBody, err = http_request.GraphqlRequest("http://localhost:8088/graphql", "query{cas}")
	if err != nil {
		t.Error(err)
	}
	respJson, err := simplejson.NewJson(respBody)
	if err != nil {
		t.Error(err)
	}
	ticket, _ := respJson.Get("data").Get("cas").String()
	respBody, err = http_request.GraphqlRequest("http://localhost:8088/graphql", fmt.Sprintf(`mutation{vote(names:["Alice", "Bob"], ticket:"%v"){name votes}}`, ticket))
	t.Logf("Vote result: %v", string(respBody))
	respJson, err = simplejson.NewJson(respBody)
	if err != nil {
		t.Error(err)
	}
}
