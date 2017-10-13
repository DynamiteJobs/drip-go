package drip_test

import (
	"os"
	"testing"

	"github.com/atishpatel/drip-go"
)

var (
	// TODO: setup for test
	apiKey    = os.Getenv("DRIP_API_KEY")
	accountID = os.Getenv("DRIP_ACCOUNT_ID")
	testEmail = "test@test.com"
)

type mockSubscribersResp struct {
	desc         string
	minSubs      int
	hasError     bool
	minCodeError int
}

type mockResp struct {
	desc         string
	hasError     bool
	minCodeError int
}

func TestNew(t *testing.T) {
	t.Logf("APIKey(%s) AccountID(%s", apiKey, accountID)
	var err error
	_, err = drip.New("", "123")
	if err != drip.ErrBadAPIKey {
		t.Errorf("Failed to get ErrBadAPIKey")
	}
	_, err = drip.New("acb123", "")
	if err != drip.ErrBadAccountID {
		t.Errorf("Failed to get ErrBadAccountID")
	}
	_, err = drip.New(apiKey, accountID)
	if err != nil {
		t.Errorf("Failed because got error: %s", err)
	}
}

func TestListSubscribers(t *testing.T) {
	tables := []struct {
		req  *drip.ListSubscribersReq
		resp *mockSubscribersResp
	}{
		{
			req: &drip.ListSubscribersReq{},
			resp: &mockSubscribersResp{
				desc:         "failed to get min 1 sub with no parms",
				minSubs:      1,
				hasError:     false,
				minCodeError: 0,
			},
		},
	}

	dripClient, err := drip.New(apiKey, accountID)
	if err != nil {
		t.Fatalf("Failed to get drip client: %s", err)
	}
	for _, table := range tables {
		resp, err := dripClient.ListSubscribers(table.req)
		if err != nil && table.resp.hasError != true {
			t.Fatalf("hasError %s: %s", table.resp.desc, err)
		}
		if len(resp.Errors) < table.resp.minCodeError {
			t.Fatalf("minCodeError %s", table.resp.desc)
		}
		if len(resp.Subscribers) < table.resp.minSubs {
			t.Fatalf("minSubs %s", table.resp.desc)
		}
	}
}

func TestUpdateSubscriber(t *testing.T) {
	tables := []struct {
		req  *drip.UpdateSubscribersReq
		resp *mockSubscribersResp
	}{
		{
			req: &drip.UpdateSubscribersReq{
				Subscribers: []drip.UpdateSubscriber{
					drip.UpdateSubscriber{
						Email:    "test@test.com",
						NewEmail: "test@test.com",
						Tags:     []string{"dev", "test"},
					},
				},
			},
			resp: &mockSubscribersResp{
				desc:         "failed to get min 1 sub with no parms",
				minSubs:      1,
				hasError:     false,
				minCodeError: 0,
			},
		},
		{
			req: &drip.UpdateSubscribersReq{},
			resp: &mockSubscribersResp{
				desc:         "failed to get error with no id and email",
				minSubs:      0,
				hasError:     true,
				minCodeError: 0,
			},
		},
	}

	dripClient, err := drip.New(apiKey, accountID)
	if err != nil {
		t.Fatalf("Failed to get drip client: %s", err)
	}
	for _, table := range tables {
		resp, err := dripClient.UpdateSubscriber(table.req)
		if err != nil && table.resp.hasError != true {
			t.Fatalf("hasError %s: %s", table.resp.desc, err)
		}
		if len(resp.Errors) < table.resp.minCodeError {
			t.Fatalf("minCodeError %s", table.resp.desc)
		}
		if len(resp.Subscribers) < table.resp.minSubs {
			t.Fatalf("minSubs %s", table.resp.desc)
		}
	}
}

func TestDeleteSubscriber(t *testing.T) {
	tables := []struct {
		idOrEmail string
		resp      *mockResp
	}{
		{
			idOrEmail: "test@test.com",
			resp: &mockResp{
				desc:         "failed to delete email",
				hasError:     false,
				minCodeError: 0,
			},
		},
	}

	dripClient, err := drip.New(apiKey, accountID)
	if err != nil {
		t.Fatalf("Failed to get drip client: %s", err)
	}
	for _, table := range tables {
		resp, err := dripClient.DeleteSubscriber(table.idOrEmail)
		if err != nil && table.resp.hasError != true {
			t.Fatalf("hasError %s: %s", table.resp.desc, err)
		}
		if len(resp.Errors) < table.resp.minCodeError {
			t.Fatalf("minCodeError %s", table.resp.desc)
		}
	}
}
