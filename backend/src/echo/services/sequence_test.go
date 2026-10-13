package services

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"beo-echo/backend/src/database"
)

func flakyResponses() []database.MockResponse {
	// Priority decides the order: 500 first, then 200.
	return []database.MockResponse{
		{StatusCode: 200, Priority: 1},
		{StatusCode: 500, Priority: 2},
	}
}

func statuses(endpointID string, n int, onEnd database.SequenceOnEnd) []int {
	var got []int
	for i := 0; i < n; i++ {
		got = append(got, getNextSequenceResponse(endpointID, flakyResponses(), onEnd).StatusCode)
	}
	return got
}

func TestSequence_StayLastIsDefault(t *testing.T) {
	endpointStates = sync.Map{}
	assert.Equal(t, []int{500, 200, 200, 200}, statuses("ep-default", 4, ""))
}

func TestSequence_StayLast(t *testing.T) {
	endpointStates = sync.Map{}
	assert.Equal(t, []int{500, 200, 200, 200}, statuses("ep-stay", 4, database.SequenceStayLast))
}

func TestSequence_Loop(t *testing.T) {
	endpointStates = sync.Map{}
	assert.Equal(t, []int{500, 200, 500, 200, 500}, statuses("ep-loop", 5, database.SequenceLoop))
}

func TestSequence_ResetStartsOver(t *testing.T) {
	endpointStates = sync.Map{}
	assert.Equal(t, []int{500, 200, 200}, statuses("ep-reset", 3, ""))

	ResetEndpointCounter("ep-reset")

	assert.Equal(t, []int{500, 200}, statuses("ep-reset", 2, ""))
}

func TestSequence_ResetOnlyAffectsThatEndpoint(t *testing.T) {
	endpointStates = sync.Map{}
	statuses("ep-a", 2, "")
	statuses("ep-b", 2, "")

	ResetEndpointCounter("ep-a")

	assert.Equal(t, 500, getNextSequenceResponse("ep-a", flakyResponses(), "").StatusCode)
	assert.Equal(t, 200, getNextSequenceResponse("ep-b", flakyResponses(), "").StatusCode, "ep-b keeps its position")
}

func TestSequence_ResetAlsoResetsRoundRobin(t *testing.T) {
	endpointStates = sync.Map{}
	responses := flakyResponses()
	assert.Equal(t, 500, getNextRoundRobinResponse("ep-rr", responses).StatusCode)
	assert.Equal(t, 200, getNextRoundRobinResponse("ep-rr", responses).StatusCode)

	ResetEndpointCounter("ep-rr")

	assert.Equal(t, 500, getNextRoundRobinResponse("ep-rr", responses).StatusCode)
}

func TestSequence_SelectedByResponseMode(t *testing.T) {
	endpointStates = sync.Map{}
	var got []int
	for i := 0; i < 3; i++ {
		resp, err := selectResponse("ep-mode", flakyResponses(), "sequence", database.SequenceStayLast, nil)
		assert.NoError(t, err)
		got = append(got, resp.StatusCode)
	}
	assert.Equal(t, []int{500, 200, 200}, got)
}

func TestSequence_ParallelRequestsNeverSkipOrRepeatSlots(t *testing.T) {
	endpointStates = sync.Map{}
	const total = 200
	responses := make([]database.MockResponse, total)
	for i := range responses {
		responses[i] = database.MockResponse{StatusCode: 1000 + i, Priority: total - i}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	seen := map[int]int{}
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := getNextSequenceResponse("ep-parallel", responses, database.SequenceStayLast)
			mu.Lock()
			seen[r.StatusCode]++
			mu.Unlock()
		}()
	}
	wg.Wait()

	assert.Len(t, seen, total, "every slot is served exactly once")
	for code, n := range seen {
		assert.Equal(t, 1, n, "status %d served %d times", code, n)
	}
}

func TestSequenceOnEnd_ReadFromEndpointConfig(t *testing.T) {
	assert.Equal(t, database.SequenceOnEnd(""), sequenceOnEnd(nil))
	assert.Equal(t, database.SequenceOnEnd(""), sequenceOnEnd(&database.MockEndpoint{}))
	assert.Equal(t, database.SequenceLoop, sequenceOnEnd(&database.MockEndpoint{AdvanceConfig: `{"sequence":{"onEnd":"loop"}}`}))
	assert.Equal(t, database.SequenceOnEnd(""), sequenceOnEnd(&database.MockEndpoint{AdvanceConfig: `{"sequence":{"onEnd":"bogus"}}`}), "invalid config falls back to the default")
}
