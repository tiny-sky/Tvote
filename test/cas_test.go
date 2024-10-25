package test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tiny-sky/Tvote/cmd/http_request"
)

func TestCas(t *testing.T) {
	t.Parallel()
	var testCounts int = 100
	var wg sync.WaitGroup
	wg.Add(testCounts)
	for i := 0; i < testCounts; i++ {
		go func() {
			defer wg.Done()
			var err error
			var respBody []byte

			respBody, err = http_request.GraphqlRequest("http://localhost:8088/graphql", "query{cas}")
			if err != nil {
				t.Error(err)
			}
			firstTicket := string(respBody)
			t.Logf("First Ticket: %v", firstTicket)

			time.Sleep(2 * time.Second)

			respBody, err = http_request.GraphqlRequest("http://localhost:8088/graphql", "query{cas}")
			if err != nil {
				t.Error(err)
			}
			secondTicket := string(respBody)
			t.Logf("Second Ticket: %v", secondTicket)

			assert.NotEqual(t, firstTicket, secondTicket)
		}()
	}
	wg.Wait()
}
