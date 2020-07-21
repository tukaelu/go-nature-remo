package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDevices(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ah := fmt.Sprintf("Bearer %s", testAccessToken)
		if req.Header.Get("Authorization") != ah {
			t.Errorf("Unmatched Authorization header: %s", ah)
		}

		res.Header().Add("X-Rate-Limit-Limit", "10")
		res.Header().Add("X-Rate-Limit-Reset", "1577804400")
		res.Header().Add("X-Rate-Limit-Remaining", "10")

		var body []map[string]interface{}
		el := map[string]interface{}{
			"id":                 "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"name":               "string",
			"temperature_offset": 0,
			"humidity_offset":    0,
			"created_at":         "2020-07-21T00:12:41.230Z",
			"updated_at":         "2020-07-21T00:12:41.230Z",
			"firmware_version":   "string",
			"mac_address":        "string",
			"serial_number":      "string",
			"newest_events": map[string]interface{}{
				"te": map[string]interface{}{
					"val":        0,
					"created_at": "2020-07-21T00:12:41.230Z",
				},
				"hu": map[string]interface{}{
					"val":        0,
					"created_at": "2020-07-21T00:12:41.230Z",
				},
				"il": map[string]interface{}{
					"val":        0,
					"created_at": "2020-07-21T00:12:41.230Z",
				},
				"mo": map[string]interface{}{
					"val":        0,
					"created_at": "2020-07-21T00:12:41.230Z",
				},
			},
		}
		body = append(body, el)

		p, _ := json.Marshal(body)
		fmt.Fprint(res, string(p))
	}))
	defer sv.Close()

	cli := NewClient(testAccessToken)
	cli.BaseURL = sv.URL

	ctx := context.Background()

	api := &devices{cli: cli}
	d, err := api.GetDevices(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "3fa85f64-5717-4562-b3fc-2c963f66afa6", d[0].ID)
	assert.Equal(t, "string", d[0].SerialNumber)
}
