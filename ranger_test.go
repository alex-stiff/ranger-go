package ranger

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetFakeRangerServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		policy := `[
			{
				"id": 1,
				"name": "test policy",
				"service": "kafka"
			},
			{
				"id": 2,
				"name": "another policy",
				"service": "hdfs"
			}
		]`

		kafkaPolicy := `[
			{
				"id": 1,
				"name": "test policy",
				"service": "kafka"
			}
		]`

		hdfsPolicy := `[
			{
				"id": 2,
				"name": "another policy",
				"service": "hdfs"
			}
		]`

		switch r.URL.Path {
		case "/":
			w.WriteHeader(http.StatusOK)
		case "/service/public/v2/api/policy":
			w.WriteHeader(http.StatusOK)
			if r.URL.Query().Get("serviceName") == "kafka" {
				w.Write([]byte(kafkaPolicy))
			} else if r.URL.Query().Get("serviceName") == "hdfs" {
				w.Write([]byte(hdfsPolicy))
			} else {
				w.Write([]byte(policy))
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))

	return server
}

func TestNewClient(t *testing.T) {
	uri := "https://example.ranger.com"
	username := "testuser"
	password := "testpassword"

	c := NewClient(uri, username, password)

	if c.BaseURL != uri {
		t.Errorf("expected BaseURL %s, got %s", uri, c.BaseURL)
	}

	if c.Username != username {
		t.Errorf("expected Username %s, got %s", username, c.Username)
	}

	if c.Password != password {
		t.Errorf("expected Password %s, got %s", password, c.Password)
	}

	if c == nil {
		t.Error("expected client to be created, got nil")
	}
}

func TestDoRequest(t *testing.T) {
	testServer := GetFakeRangerServer()

	c := NewClient(testServer.URL, "testuser", "testpassword")

	req, err := http.NewRequest("GET", testServer.URL, nil)

	if err != nil {
		t.Errorf("expected no error creating request, got %v", err)
	}

	resp, err := c.doRequest(req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Error("expected response, got nil")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetPolicies(t *testing.T) {
	testServer := GetFakeRangerServer()

	defer testServer.Close()

	c := NewClient(testServer.URL, "testuser", "testpassword")

	policies, err := c.GetPolicies()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(policies))
	}

	if policies[0].ID != 1 {
		t.Errorf("expected policy ID 1, got %d", policies[0].ID)
	}

	if policies[0].Name != "test policy" {
		t.Errorf("expected policy name 'test policy', got '%s'", policies[0].Name)
	}

	if policies[0].Service != "kafka" {
		t.Errorf("expected policy service 'kafka', got '%s'", policies[0].Service)
	}

	if policies[1].ID != 2 {
		t.Errorf("expected policy ID 2, got %d", policies[1].ID)
	}

	if policies[1].Name != "another policy" {
		t.Errorf("expected policy name 'another policy', got '%s'", policies[1].Name)
	}

	if policies[1].Service != "hdfs" {
		t.Errorf("expected policy service 'hdfs', got '%s'", policies[1].Service)
	}
}

func TestGetPoliciesWithService(t *testing.T) {
	testServer := GetFakeRangerServer()

	defer testServer.Close()

	c := NewClient(testServer.URL, "testuser", "testpassword")

	policies, err := c.GetPolicies("hdfs")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(policies) != 1 {
		t.Errorf("expected 1 policy for service 'hdfs', got %d", len(policies))
	}

	if policies[0].ID != 2 {
		t.Errorf("expected policy ID 2, got %d", policies[0].ID)
	}

	if policies[0].Name != "another policy" {
		t.Errorf("expected policy name 'another policy', got '%s'", policies[0].Name)
	}

	if policies[0].Service != "hdfs" {
		t.Errorf("expected policy service 'kafka', got '%s'", policies[0].Service)
	}
}
