package services

import (
	"beo-echo/backend/src/database"
	"sync"
	"testing"
)

func TestRoundRobinSelection(t *testing.T) {
	// Clear any existing state before testing
	endpointStates = sync.Map{}

	// Test Case 1: Endpoint dengan 5 responses
	t.Run("Endpoint with 5 responses", func(t *testing.T) {
		endpointID := "endpoint-5-responses"
		responses := []database.MockResponse{
			{Body: "1", Priority: 3}, // Priority diacak, bukan urut
			{Body: "2", Priority: 5}, // Setelah sort: 2(5) -> 4(4) -> 1(3) -> 5(2) -> 3(1)
			{Body: "3", Priority: 1},
			{Body: "4", Priority: 4},
			{Body: "5", Priority: 2},
		}

		// Test 3 full cycles (15 calls total) untuk memastikan round-robin bekerja
		// Expected sequence setelah sort by priority: 2(5) -> 4(4) -> 1(3) -> 5(2) -> 3(1)
		expectedSequence := []string{"2", "4", "1", "5", "3", "2", "4", "1", "5", "3", "2", "4", "1", "5", "3"}

		for i, expected := range expectedSequence {
			response := getNextRoundRobinResponse(endpointID, responses)
			if response.Body != expected {
				t.Errorf("Call %d: expected body '%s', got '%s'", i+1, expected, response.Body)
			}
		}
	})

	// Test Case 2: Endpoint dengan 2 responses
	t.Run("Endpoint with 2 responses", func(t *testing.T) {
		endpointID := "endpoint-2-responses"
		responses := []database.MockResponse{
			{Body: "1", Priority: 3}, // Priority diacak
			{Body: "2", Priority: 1}, // Setelah sort: 1(3) -> 2(1)
		}

		// Test 3 full cycles (6 calls total)
		// Expected sequence setelah sort by priority: 1(3) -> 2(1)
		expectedSequence := []string{"1", "2", "1", "2", "1", "2"}

		for i, expected := range expectedSequence {
			response := getNextRoundRobinResponse(endpointID, responses)
			if response.Body != expected {
				t.Errorf("Call %d: expected body '%s', got '%s'", i+1, expected, response.Body)
			}
		}
	})

	// Test Case 3: Endpoint dengan 1 response
	t.Run("Endpoint with 1 response", func(t *testing.T) {
		endpointID := "endpoint-1-response"
		responses := []database.MockResponse{
			{Body: "1", Priority: 7}, // Priority bebas, cuma satu response
		}

		// Test 5 calls - semuanya harus return response yang sama
		for i := 0; i < 5; i++ {
			response := getNextRoundRobinResponse(endpointID, responses)
			if response.Body != "1" {
				t.Errorf("Call %d: expected body '1', got '%s'", i+1, response.Body)
			}
		}
	})
}

func TestMultipleEndpointsIndependentState(t *testing.T) {
	// Clear any existing state before testing
	endpointStates = sync.Map{}

	// Setup 3 endpoints dengan responses yang berbeda
	endpoint1 := "endpoint-1"
	responses1 := []database.MockResponse{
		{Body: "A1", Priority: 2}, // Priority diacak
		{Body: "A2", Priority: 3}, // Setelah sort: A2(3) -> A1(2) -> A3(1)
		{Body: "A3", Priority: 1},
	}

	endpoint2 := "endpoint-2"
	responses2 := []database.MockResponse{
		{Body: "B1", Priority: 4}, // Priority diacak
		{Body: "B2", Priority: 1}, // Setelah sort: B1(4) -> B2(1)
	}

	endpoint3 := "endpoint-3"
	responses3 := []database.MockResponse{
		{Body: "C1", Priority: 2}, // Priority diacak
		{Body: "C2", Priority: 4}, // Setelah sort: C2(4) -> C4(3) -> C1(2) -> C3(1)
		{Body: "C3", Priority: 1},
		{Body: "C4", Priority: 3},
	}

	// Test interleaved calls untuk memastikan setiap endpoint memiliki state independen
	testCases := []struct {
		endpointID string
		responses  []database.MockResponse
		expected   string
		callNum    int
	}{
		// Round 1 - expected setelah sort by priority
		{endpoint1, responses1, "A2", 1}, // endpoint1 call 1 (A2 has highest priority 3)
		{endpoint2, responses2, "B1", 1}, // endpoint2 call 1 (B1 has highest priority 4)
		{endpoint3, responses3, "C2", 1}, // endpoint3 call 1 (C2 has highest priority 4)

		// Round 2
		{endpoint1, responses1, "A1", 2}, // endpoint1 call 2 (A1 has priority 2)
		{endpoint2, responses2, "B2", 2}, // endpoint2 call 2 (B2 has priority 1)
		{endpoint3, responses3, "C4", 2}, // endpoint3 call 2 (C4 has priority 3)

		// Round 3
		{endpoint1, responses1, "A3", 3}, // endpoint1 call 3 (A3 has priority 1)
		{endpoint2, responses2, "B1", 3}, // endpoint2 call 3 (wrap around to B1)
		{endpoint3, responses3, "C1", 3}, // endpoint3 call 3 (C1 has priority 2)

		// Round 4
		{endpoint1, responses1, "A2", 4}, // endpoint1 call 4 (wrap around to A2)
		{endpoint2, responses2, "B2", 4}, // endpoint2 call 4
		{endpoint3, responses3, "C3", 4}, // endpoint3 call 4 (C3 has priority 1)

		// Round 5
		{endpoint1, responses1, "A1", 5}, // endpoint1 call 5
		{endpoint2, responses2, "B1", 5}, // endpoint2 call 5 (wrap around)
		{endpoint3, responses3, "C2", 5}, // endpoint3 call 5 (wrap around to C2)
	}

	for i, tc := range testCases {
		response := getNextRoundRobinResponse(tc.endpointID, tc.responses)
		if response.Body != tc.expected {
			t.Errorf("Test %d (%s call %d): expected body '%s', got '%s'",
				i+1, tc.endpointID, tc.callNum, tc.expected, response.Body)
		}
	}
}

func TestRoundRobinWithEmptyResponses(t *testing.T) {
	// Clear any existing state before testing
	endpointStates = sync.Map{}

	endpointID := "empty-endpoint"
	responses := []database.MockResponse{}

	response := getNextRoundRobinResponse(endpointID, responses)

	// Harus return empty MockResponse
	if response.Body != "" {
		t.Errorf("Expected empty body for empty responses, got '%s'", response.Body)
	}
}

func TestRoundRobinSequentialCalls(t *testing.T) {
	// Clear any existing state before testing
	endpointStates = sync.Map{}

	t.Run("Test endpoint dengan 5 responses - 20 calls berturut-turut", func(t *testing.T) {
		endpointID := "sequential-test"
		responses := []database.MockResponse{
			{Body: "1", Priority: 3}, // Priority diacak untuk test sorting
			{Body: "2", Priority: 5}, // Setelah sort: 2(5) -> 4(4) -> 1(3) -> 5(2) -> 3(1)
			{Body: "3", Priority: 1},
			{Body: "4", Priority: 4},
			{Body: "5", Priority: 2},
		}

		// Expected pattern setelah sort: 2,4,1,5,3,2,4,1,5,3,...
		expectedBodies := []string{"2", "4", "1", "5", "3"}

		for i := 0; i < 20; i++ {
			response := getNextRoundRobinResponse(endpointID, responses)
			expectedBody := expectedBodies[i%5]

			if response.Body != expectedBody {
				t.Errorf("Call %d: expected body '%s', got '%s'", i+1, expectedBody, response.Body)
			}
		}
	})

	t.Run("Test endpoint dengan 2 responses - 10 calls berturut-turut", func(t *testing.T) {
		endpointID := "sequential-test-2"
		responses := []database.MockResponse{
			{Body: "1", Priority: 8}, // Priority diacak
			{Body: "2", Priority: 3}, // Setelah sort: 1(8) -> 2(3)
		}

		// Expected pattern setelah sort: 1,2,1,2,1,2,1,2,1,2
		expectedBodies := []string{"1", "2"}

		for i := 0; i < 10; i++ {
			response := getNextRoundRobinResponse(endpointID, responses)
			expectedBody := expectedBodies[i%2]

			if response.Body != expectedBody {
				t.Errorf("Call %d: expected body '%s', got '%s'", i+1, expectedBody, response.Body)
			}
		}
	})

	t.Run("Test endpoint dengan 1 response - 5 calls berturut-turut", func(t *testing.T) {
		endpointID := "sequential-test-1"
		responses := []database.MockResponse{
			{Body: "1", Priority: 9}, // Priority bebas, cuma satu response
		}

		// Expected: selalu return "1"
		for i := 0; i < 5; i++ {
			response := getNextRoundRobinResponse(endpointID, responses)

			if response.Body != "1" {
				t.Errorf("Call %d: expected body '1', got '%s'", i+1, response.Body)
			}
		}
	})
}
