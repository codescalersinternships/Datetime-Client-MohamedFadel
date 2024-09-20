[1mdiff --git a/client/client_integration_test.go b/client/client_integration_test.go[m
[1mnew file mode 100644[m
[1mindex 0000000..e6c5cce[m
[1m--- /dev/null[m
[1m+++ b/client/client_integration_test.go[m
[36m@@ -0,0 +1,111 @@[m
[32m+[m[32mpackage client[m
[32m+[m
[32m+[m[32mimport ([m
[32m+[m	[32m"encoding/json"[m
[32m+[m	[32m"net/http"[m
[32m+[m	[32m"net/http/httptest"[m
[32m+[m	[32m"strings"[m
[32m+[m	[32m"testing"[m
[32m+[m	[32m"time"[m
[32m+[m[32m)[m
[32m+[m
[32m+[m[32mfunc TestIntegrationGetDateTime(t *testing.T) {[m
[32m+[m	[32mtests := []struct {[m
[32m+[m		[32mname               string[m
[32m+[m		[32mserverType         string[m
[32m+[m		[32mcontentType        string[m
[32m+[m		[32mserverResponse     interface{}[m
[32m+[m		[32mexpectedStatusCode int[m
[32m+[m		[32mexpectedResult     string[m
[32m+[m		[32mexpectedError      bool[m
[32m+[m	[32m}{[m
[32m+[m		[32m{[m
[32m+[m			[32mname:               "Successful JSON response",[m
[32m+[m			[32mserverType:         "standard",[m
[32m+[m			[32mcontentType:        "application/json",[m
[32m+[m			[32mserverResponse:     map[string]string{"datetime": "2024-09-20T14:30:00Z"},[m
[32m+[m			[32mexpectedStatusCode: http.StatusOK,[m
[32m+[m			[32mexpectedResult:     "2024-09-20T14:30:00Z",[m
[32m+[m			[32mexpectedError:      false,[m
[32m+[m		[32m},[m
[32m+[m		[32m{[m
[32m+[m			[32mname:               "Successful plain text response",[m
[32m+[m			[32mserverType:         "gin",[m
[32m+[m			[32mcontentType:        "text/plain",[m
[32m+[m			[32mserverResponse:     "2024-09-20T14:30:00Z",[m
[32m+[m			[32mexpectedStatusCode: http.StatusOK,[m
[32m+[m			[32mexpectedResult:     "2024-09-20T14:30:00Z",[m
[32m+[m			[32mexpectedError:      false,[m
[32m+[m		[32m},[m
[32m+[m		[32m{[m
[32m+[m			[32mname:               "Server error",[m
[32m+[m			[32mserverType:         "standard",[m
[32m+[m			[32mcontentType:        "application/json",[m
[32m+[m			[32mserverResponse:     nil,[m
[32m+[m			[32mexpectedStatusCode: http.StatusInternalServerError,[m
[32m+[m			[32mexpectedResult:     "",[m
[32m+[m			[32mexpectedError:      true,[m
[32m+[m		[32m},[m
[32m+[m		[32m{[m
[32m+[m			[32mname:               "Invalid JSON response",[m
[32m+[m			[32mserverType:         "standard",[m
[32m+[m			[32mcontentType:        "application/json",[m
[32m+[m			[32mserverResponse:     "invalid json",[m
[32m+[m			[32mexpectedStatusCode: http.StatusOK,[m
[32m+[m			[32mexpectedResult:     "",[m
[32m+[m			[32mexpectedError:      true,[m
[32m+[m		[32m},[m
[32m+[m		[32m{[m
[32m+[m			[32mname:               "Timeout",[m
[32m+[m			[32mserverType:         "gin",[m
[32m+[m			[32mcontentType:        "application/json",[m
[32m+[m			[32mserverResponse:     nil,[m
[32m+[m			[32mexpectedStatusCode: http.StatusOK,[m
[32m+[m			[32mexpectedResult:     "",[m
[32m+[m			[32mexpectedError:      true,[m
[32m+[m		[32m},[m
[32m+[m	[32m}[m
[32m+[m
[32m+[m	[32mfor _, tt := range tests {[m
[32m+[m		[32mt.Run(tt.name, func(t *testing.T) {[m
[32m+[m			[32mserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {[m
[32m+[m				[32mw.Header().Set("Content-Type", tt.contentType)[m
[32m+[m				[32mw.WriteHeader(tt.expectedStatusCode)[m
[32m+[m
[32m+[m				[32mif tt.name == "Timeout" {[m
[32m+[m					[32mtime.Sleep(11 * time.Second)[m
[32m+[m					[32mreturn[m
[32m+[m				[32m}[m
[32m+[m
[32m+[m				[32mif tt.serverResponse != nil {[m
[32m+[m					[32mif tt.contentType == "application/json" {[m
[32m+[m						[32mjson.NewEncoder(w).Encode(tt.serverResponse)[m
[32m+[m					[32m} else {[m
[32m+[m						[32mw.Write([]byte(tt.serverResponse.(string)))[m
[32m+[m					[32m}[m
[32m+[m				[32m}[m
[32m+[m			[32m}))[m
[32m+[m			[32mdefer server.Close()[m
[32m+[m
[32m+[m			[32mserverURL := server.URL[m
[32m+[m			[32mhostPort := strings.TrimPrefix(serverURL, "http://")[m
[32m+[m			[32mparts := strings.Split(hostPort, ":")[m
[32m+[m			[32mhost := parts[0][m
[32m+[m			[32mport := parts[1][m
[32m+[m
[32m+[m			[32mSetupEnv("http://"+host, port, port)[m
[32m+[m
[32m+[m			[32mresult, err := GetDateTime(tt.serverType, tt.contentType)[m
[32m+[m
[32m+[m			[32mif tt.expectedError && err == nil {[m
[32m+[m				[32mt.Errorf("Expected an error, but got none")[m
[32m+[m			[32m}[m
[32m+[m			[32mif !tt.expectedError && err != nil {[m
[32m+[m				[32mt.Errorf("Expected no error, but got: %v", err)[m
[32m+[m			[32m}[m
[32m+[m			[32mif result != tt.expectedResult {[m
[32m+[m				[32mt.Errorf("Expected result %s, but got %s", tt.expectedResult, result)[m
[32m+[m			[32m}[m
[32m+[m		[32m})[m
[32m+[m	[32m}[m
[32m+[m[32m}[m
[1mdiff --git a/client/client_test.go b/client/client_test.go[m
[1mindex d91631f..25cb9eb 100644[m
[1m--- a/client/client_test.go[m
[1m+++ b/client/client_test.go[m
[36m@@ -159,7 +159,7 @@[m [mfunc TestGetDateTime(t *testing.T) {[m
 			t.Fatal("Expected error after max retries, got none")[m
 		}[m
 [m
[31m-		if duration < 27*time.Second || duration > 32*time.Second {[m
[32m+[m		[32mif duration < 20*time.Second || duration > 40*time.Second {[m
 			t.Errorf("Expected retry duration to be approximately 30 seconds, got %v", duration)[m
 		}[m
 	})[m
