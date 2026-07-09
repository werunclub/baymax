package client

import (
	"context"
	"testing"
	"time"

	rpcxClient "github.com/werunclub/rpcx/v4/client"
)

// stubXClient 只实现 Call，用于捕获 CallWithContext 实际下发的 context。
type stubXClient struct {
	rpcxClient.XClient
	gotDeadline time.Time
	hasDeadline bool
}

func (s *stubXClient) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error {
	s.gotDeadline, s.hasDeadline = ctx.Deadline()
	return nil
}

func newTestClient(opts ...Option) (*Client, *stubXClient) {
	stub := &stubXClient{}
	return &Client{
		opts:      newOptions(opts...),
		rpcClient: stub,
	}, stub
}

func TestCallWithContextAppliesDefaultRequestTimeout(t *testing.T) {
	c, stub := newTestClient()

	before := time.Now()
	if err := c.CallWithContext(context.Background(), "Svc.Method", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !stub.hasDeadline {
		t.Fatal("expected default RequestTimeout deadline, got none")
	}
	remaining := stub.gotDeadline.Sub(before)
	if remaining < 29*time.Second || remaining > 31*time.Second {
		t.Fatalf("expected deadline ~30s from now, got %v", remaining)
	}
}

func TestCallWithContextKeepsCallerDeadline(t *testing.T) {
	c, stub := newTestClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	callerDeadline, _ := ctx.Deadline()

	if err := c.CallWithContext(ctx, "Svc.Method", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !stub.hasDeadline || !stub.gotDeadline.Equal(callerDeadline) {
		t.Fatalf("caller deadline should be kept, got %v (has=%v), want %v",
			stub.gotDeadline, stub.hasDeadline, callerDeadline)
	}
}

func TestCallWithContextCustomRequestTimeout(t *testing.T) {
	c, stub := newTestClient(RequestTimeout(3 * time.Second))

	before := time.Now()
	if err := c.CallWithContext(context.Background(), "Svc.Method", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !stub.hasDeadline {
		t.Fatal("expected custom RequestTimeout deadline, got none")
	}
	remaining := stub.gotDeadline.Sub(before)
	if remaining < 2*time.Second || remaining > 4*time.Second {
		t.Fatalf("expected deadline ~3s from now, got %v", remaining)
	}
}

func TestCallWithContextNegativeRequestTimeoutDisables(t *testing.T) {
	c, stub := newTestClient(RequestTimeout(-1))

	if err := c.CallWithContext(context.Background(), "Svc.Method", nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stub.hasDeadline {
		t.Fatalf("negative RequestTimeout should disable deadline, got %v", stub.gotDeadline)
	}
}

func TestDefaultOptionsRequestTimeout(t *testing.T) {
	opts := newOptions()
	if opts.RequestTimeout != 30*time.Second {
		t.Fatalf("expected default RequestTimeout 30s, got %v", opts.RequestTimeout)
	}
}
