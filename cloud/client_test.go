package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAccessToken = "DUMMY_ACCESS_TOKEN"
)

func TestGetMe(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ah := fmt.Sprintf("Bearer %s", testAccessToken)
		if req.Header.Get("Authorization") != ah {
			t.Errorf("Unmatched Authorization header: %s", ah)
		}

		res.Header().Add("X-Rate-Limit-Limit", "10")
		res.Header().Add("X-Rate-Limit-Reset", "1577804400")
		res.Header().Add("X-Rate-Limit-Remaining", "10")

		body := map[string]string{
			"id":       "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"nickname": "string",
		}
		p, _ := json.Marshal(body)
		fmt.Fprint(res, string(p))
	}))
	defer sv.Close()

	u := User{}
	cli := NewClient(testAccessToken)
	cli.BaseURL = sv.URL

	err := cli.Get(context.Background(), "users/me", nil, &u)
	assert.NoError(t, err)

	assert.Equal(t, "3fa85f64-5717-4562-b3fc-2c963f66afa6", u.ID)
	assert.Equal(t, "string", u.Nickname)
}

func TestUpdateMe(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ah := fmt.Sprintf("Bearer %s", testAccessToken)
		if req.Header.Get("Authorization") != ah {
			t.Errorf("Unmatched Authorization header: %s", ah)
		}

		res.Header().Add("X-Rate-Limit-Limit", "10")
		res.Header().Add("X-Rate-Limit-Reset", "1577804400")
		res.Header().Add("X-Rate-Limit-Remaining", "10")

		body := map[string]string{
			"nickname": "foobar",
		}
		p, _ := json.Marshal(body)
		fmt.Fprint(res, string(p))
	}))

	defer sv.Close()

	u := User{}
	cli := NewClient(testAccessToken)
	cli.BaseURL = sv.URL

	p := url.Values{}
	p.Add("nickname", "foobar")

	err := cli.Post(context.Background(), "users/me", p, &u)
	assert.NoError(t, err)

	assert.Equal(t, p.Get("nickname"), u.Nickname)
}
