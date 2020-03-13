package vsClient

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/types"
)

// VsClient is a client for vSphere.
type Client struct {
	Govmomi *govmomi.Client
	Rest    *rest.Client
}

func New(ctx context.Context, u url.URL, insecure bool) (*Client, error) {
	var clt Client

	gc, err := govmomi.NewClient(ctx, &u, insecure)
	if err != nil {
		return &clt, fmt.Errorf("connecting to govmomi api failed: %w", err)
	}
	clt.Govmomi = gc

	clt.Rest = rest.NewClient(clt.Govmomi.Client)
	clt.Rest.Login(ctx, u.User)

	return &clt, nil
}

// Tag adds an existing tag to a VirtualMachine.
func (client *Client) MoTag(ctx context.Context, vm types.ManagedObjectReference, tagID string) error {
	// Get the tag manager which does the tagging.
	m := tags.NewManager(client.Rest)

	// Attach tag to VM.
	err := m.AttachTag(ctx, tagID, vm)
	if err != nil {
		return fmt.Errorf("attach tag to VM failed: %w", err)
	}

	return nil
}

func (client *Client) Logout(ctx context.Context) error {
	err := client.Govmomi.Logout(ctx)
	if err != nil {
		return fmt.Errorf("govmomi api logout failed: %w", err)
	}

	err = client.Rest.Logout(ctx)
	if err != nil {
		return fmt.Errorf("rest api logout failed: %w", err)
	}

	return nil
}
