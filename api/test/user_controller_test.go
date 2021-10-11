package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MISW/birdol-server/controller"
	"github.com/MISW/birdol-server/controller/jsonmodel"
	"github.com/MISW/birdol-server/database"
	"github.com/gin-gonic/gin"
)

func TestSignUp(t *testing.T) {
	req_bodies := []jsonmodel.SignupUserRequest {
		{
			Name: "キリト",
			PublicKey: "PFJTQUtleVZhbHVlPjxNb2R1bHVzPjRYQ2NuclBCSENqNlBtU3ExbHY3MGcvUkUwL0ZaWGhpQVI0akQ2Wk1qeGh4bE9mMkxaV0JrRUJUWmRHNnZ0T0lyWHRaMDFsZkl1YmVsNGFWcDEzU2JhTUJ1UmtaRE0vRGhPUkRVY3lXMTZVTko3UjNxaEk1SHNFZlNldUFMa1g1Njg2UXdneDYwUGpVR1Z2bFJSYVI4RXVqR20zY1VoUXBEOEdYVW1sbUU1az08L01vZHVsdXM+PEV4cG9uZW50PkFRQUI8L0V4cG9uZW50PjwvUlNBS2V5VmFsdWU+",
			DeviceID: "Elucidator-810-893",
		},
		{
			Name: "キバオウ",
			PublicKey: "PFJTQUtleVZhbHVlPjxNb2R1bHVzPjAyNkVxVDJjcm5nYzF0QWFpQ0dMUjNwcUJvbFB1UFVhZk5iQUYxMDhzWlpQeVhoV05QRWRhRFBHVG1sTktwekVMcGlwd3NtZHlyTDk3UWNEZnlDS0xFMWVQdXhFMGtCZ2t6NVhKK0pDNG9xWTB1YkhBZFFUSm9rZWFoMGpIb2F1c1ZQYzJIbHllTjFidEd0SDFOU1A2dXRHTnlwNXYvWWphUnFZVXV4ZjBWVT08L01vZHVsdXM+PEV4cG9uZW50PkFRQUI8L0V4cG9uZW50PjwvUlNBS2V5VmFsdWU+",
			DeviceID: "Beeter-9302-293484",
		},
	}

	database.TestingDatabase()

	for i, t_case := range req_bodies {
		t.Logf("Testing: case %d of %d\n", i+1, len(req_bodies))
		json_body, _ := json.Marshal(t_case)
		req_reader := bytes.NewBuffer(json_body)
		response := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(response)
		ctx.Request, _ = http.NewRequest(
			http.MethodPut,
			"api/v2/user",
			req_reader,
		)
		ctx.Request.Header.Set("Content-Type", gin.MIMEJSON)
		
		controller.HandleSignUp()(ctx)

		var response_body map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &response_body)

		const process = 6
		passed := 0
		if response.Code != 200 {
			t.Errorf("[%d/%d] Not expected status code: %d\n", passed, process, response.Code)
		} else {
			passed++
			t.Logf("[%d/%d] Correct status code: %d\n", passed, process, response.Code)
		}

		content_type := response.Header().Get("Content-Type")
		if content_type != gin.MIMEJSON + "; charset=utf-8" {
			t.Errorf("[%d/%d] Not expected Content-Type: %s\n", passed, process, content_type)
		} else {
			passed++
			t.Logf("[%d/%d] Correct Content-Type: %s\n", passed, process, content_type)
		}

		expected_keys := []string {
			"user_id",
			"access_token",
			"refresh_token",
			"account_id",
		}
		
		for _, key := range expected_keys {
			if _, exist := response_body[key]; exist {
				passed++
				t.Logf("[%d/%d] %s exists\n", passed, process, key)
			} else {
				t.Errorf("[%d/%d] %s does not exist\n", passed, process, key)
			}
		}

		if passed == process {
			t.Logf("All the test process passed.\n")
		} else {
			t.Errorf("Some process recoded error.\n")
		}
	}
}

func TestLinkNormal(t *testing.T) {

}

func TestLinkExpire(t *testing.T) {

}

func TestLinkNotExist(t *testing.T) {

}