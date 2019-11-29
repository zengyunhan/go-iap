package amazon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestHandle497Error(t *testing.T) {
	t.Parallel()
	var expected, actual error
	client := New("developerSecret")

	server, client := testTools(
		497,
		"{\"message\":\"Purchase token/app user mismatch\",\"status\":false}",
	)
	defer server.Close()

	// status 400
	expected = errors.New("Purchase token/app user mismatch")
	_, actual = client.Verify(
		context.Background(),
		"99FD_DL23EMhrOGDnur9-ulvqomrSg6qyLPSD3CFE=",
		"q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ",
	)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestHandle400Error(t *testing.T) {
	t.Parallel()
	var expected, actual error
	client := New("developerSecret")

	server, client := testTools(
		400,
		"{\"message\":\"Failed to parse receipt Id\",\"status\":false}",
	)
	defer server.Close()

	// status 400
	expected = errors.New("Failed to parse receipt Id")
	_, actual = client.Verify(
		context.Background(),
		"99FD_DL23EMhrOGDnur9-ulvqomrSg6qyLPSD3CFE=",
		"q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ",
	)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNew(t *testing.T) {
	expected := &Client{
		URL:    SandboxURL,
		Secret: "developerSecret",
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	actual := New("developerSecret")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithEnvironment(t *testing.T) {
	expected := &Client{
		URL:    ProductionURL,
		Secret: "developerSecret",
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	os.Setenv("IAP_ENVIRONMENT", "production")
	actual := New("developerSecret")
	os.Clearenv()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithClient(t *testing.T) {
	expected := &Client{
		URL:    ProductionURL,
		Secret: "developerSecret",
		httpCli: &http.Client{
			Timeout: time.Second * 2,
		},
	}
	os.Setenv("IAP_ENVIRONMENT", "production")

	cli := &http.Client{
		Timeout: time.Second * 2,
	}
	actual := NewWithClient("developerSecret", cli)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifySubscription(t *testing.T) {
	t.Parallel()
	server, client := testTools(
		200,
		"{\"purchaseDate\":1558424877035,\"receiptId\":\"q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ\",\"productId\":\"com.amazon.iapsamplev2.expansion_set_3\",\"parentProductId\":null,\"productType\":\"SUBSCRIPTION\",\"renewalDate\":1561103277035,\"quantity\":1,\"betaProduct\":false,\"testTransaction\":true,\"term\":\"1 Week\",\"termSku\":\"sub1-weekly\"}",
	)
	defer server.Close()

	expected := IAPResponse{
		ReceiptID:       "q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ",
		ProductType:     "SUBSCRIPTION",
		ProductID:       "com.amazon.iapsamplev2.expansion_set_3",
		PurchaseDate:    1558424877035,
		RenewalDate:     1561103277035,
		CancelDate:      0,
		TestTransaction: true,
		Quantity:        1,
		Term:            "1 Week",
		TermSku:         "sub1-weekly",
	}

	actual, _ := client.Verify(
		context.Background(),
		"99FD_DL23EMhrOGDnur9-ulvqomrSg6qyLPSD3CFE=",
		"q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ",
	)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyEntitled(t *testing.T) {
	t.Parallel()
	server, client := testTools(
		200,
		"{\"purchaseDate\":1402008634018,\"receiptId\":\"q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ\",\"productId\":\"com.amazon.iapsamplev2.expansion_set_3\",\"parentProductId\":null,\"productType\":\"ENTITLED\",\"cancelDate\":null,\"quantity\":1,\"betaProduct\":false,\"testTransaction\":true}",
	)
	defer server.Close()

	expected := IAPResponse{
		ReceiptID:       "q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ",
		ProductType:     "ENTITLED",
		ProductID:       "com.amazon.iapsamplev2.expansion_set_3",
		PurchaseDate:    1402008634018,
		CancelDate:      0,
		TestTransaction: true,
		Quantity:        1,
	}

	actual, _ := client.Verify(
		context.Background(),
		"99FD_DL23EMhrOGDnur9-ulvqomrSg6qyLPSD3CFE=",
		"q1YqVrJSSs7P1UvMTazKz9PLTCwoTswtyEktM9JLrShIzCvOzM-LL04tiTdW0lFKASo2NDEwMjCwMDM2MTC0AIqVAsUsLd1c4l18jIxdfTOK_N1d8kqLLHVLc8oK83OLgtPNCit9AoJdjJ3dXG2BGkqUrAxrAQ",
	)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyTimeout(t *testing.T) {
	t.Parallel()
	// HTTP 100 is "continue" so it will time out
	server, client := testTools(100, "timeout response")
	defer server.Close()

	ctx := context.Background()
	_, actual := client.Verify(ctx, "timeout", "timeout")

	// Actual should be a "request canceled" *url.Error
	urlErr, ok := actual.(*url.Error)
	if !ok {
		t.Errorf("Expected *url.Error, got %T", actual)
	}

	if !urlErr.Timeout() {
		t.Errorf("got %v\nwant timeout", actual)
	}
}

func testTools(code int, body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	client := &Client{URL: server.URL, Secret: "developerSecret", httpCli: &http.Client{Timeout: 2 * time.Second}}
	return server, client
}
