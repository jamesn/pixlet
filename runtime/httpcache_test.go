package runtime

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitHTTP(t *testing.T) {
	c := NewInMemoryCache()
	InitHTTP(c)

	b, err := os.ReadFile("testdata/httpcache.star")
	assert.NoError(t, err)

	app, err := NewApplet("httpcache.star", b)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

// TestDetermineTTL tests the DetermineTTL function.
func TestDetermineTTL(t *testing.T) {
	type test struct {
		ttl         int
		retryAfter  int
		resHeader   string
		statusCode  int
		method      string
		expectedTTL time.Duration
	}

	tests := map[string]test{
		"test request cache control headers": {
			ttl:         3600,
			resHeader:   "",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test response cache control headers": {
			ttl:         0,
			resHeader:   "public, max-age=3600, s-maxage=7200, no-transform",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test too long response cache control headers": {
			ttl:         0,
			resHeader:   "max-age=604800",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test max-age of zero": {
			ttl:         0,
			resHeader:   "max-age=0",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
		"test both request and response cache control headers": {
			ttl:         3600,
			resHeader:   "public, max-age=60, s-maxage=7200, no-transform",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 3600 * time.Second,
		},
		"test 500 response code": {
			ttl:         3600,
			resHeader:   "",
			statusCode:  500,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
		"test too low ttl": {
			ttl:         3,
			resHeader:   "",
			statusCode:  200,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
		"test DELETE request": {
			ttl:         0,
			resHeader:   "",
			statusCode:  200,
			method:      "DELETE",
			expectedTTL: 5 * time.Second,
		},
		"test POST request configured with TTL": {
			ttl:         30,
			resHeader:   "",
			statusCode:  200,
			method:      "POST",
			expectedTTL: 30 * time.Second,
		},
		"test POST request without configured TTL": {
			ttl:         0,
			resHeader:   "",
			statusCode:  200,
			method:      "POST",
			expectedTTL: 5 * time.Second,
		},
		"test 429 request": {
			ttl:         30,
			retryAfter:  60,
			resHeader:   "",
			statusCode:  429,
			method:      "GET",
			expectedTTL: 60 * time.Second,
		},
		"test 429 request below minimum": {
			ttl:         30,
			retryAfter:  3,
			resHeader:   "",
			statusCode:  429,
			method:      "GET",
			expectedTTL: 5 * time.Second,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := &http.Request{
				Header: map[string][]string{
					"X-Tidbyt-Cache-Seconds": {fmt.Sprintf("%d", tc.ttl)},
				},
				Method: tc.method,
			}

			res := &http.Response{
				Header: map[string][]string{
					"Cache-Control": {tc.resHeader},
				},
				StatusCode: tc.statusCode,
			}

			if tc.retryAfter > 0 {
				res.Header.Set("Retry-After", fmt.Sprintf("%d", tc.retryAfter))
			}

			ttl := determineTTL(req, res)
			assert.Equal(t, tc.expectedTTL, ttl)
		})
	}
}

func TestDetermineTTLJitter(t *testing.T) {
	req := &http.Request{
		Header: map[string][]string{
			"X-Tidbyt-Cache-Seconds": {"60"},
		},
		Method: "GET",
	}

	res := &http.Response{
		StatusCode: 200,
	}

	// Use a seeded random source for deterministic behavior
	oldRand := rand.New(rand.NewSource(50))
	// Temporarily replace the global random function behavior by calculating expected value
	// With seed 50 and 60 seconds base: jitter = 6, so range is -6 to +6
	// We need to determine what Int63n(13) returns with seed 50
	testRand := rand.New(rand.NewSource(50))
	jitter := int64(float64(60) * 0.1) // = 6
	randomJitter := testRand.Int63n(2*jitter+1) - jitter // Int63n(13) - 6
	expectedTTL := 60 + randomJitter
	
	// Since we can't control the global rand, let's just verify the TTL is within the jitter range
	ttl := DetermineTTL(req, res)
	// TTL should be 60 +/- 6 seconds (10% jitter)
	assert.GreaterOrEqual(t, int(ttl.Seconds()), 54)
	assert.LessOrEqual(t, int(ttl.Seconds()), 66)
	
	_ = oldRand // Suppress unused variable warning
	_ = expectedTTL // Suppress unused variable warning
}

func TestDetermineTTLNoHeaders(t *testing.T) {
	req := &http.Request{
		Method: "GET",
	}

	res := &http.Response{
		StatusCode: 200,
	}

	ttl := DetermineTTL(req, res)
	// DetermineTTL applies jitter, so the result should be MinRequestTTL +/- 10%
	// MinRequestTTL is 5 seconds, so jitter range is -0.5 to +0.5 seconds
	assert.GreaterOrEqual(t, int(ttl.Seconds()), 4)
	assert.LessOrEqual(t, int(ttl.Seconds()), 6)
}
