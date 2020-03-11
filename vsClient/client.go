package vsClient

import (
	"context"
	"log"
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
		log.Println(err)
		return vsClient, err
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
		log.Fatal(err)
		return err
	}

	return nil
}

func (clt *Client) Logout(ctx) {
	err = clt.Govmomi.Logout(ctx)
	if err != nil {
		log.Printf("govmomi logout failed: %v", err)
	}
	err = clt.Rest.Logout(ctx)
	if err != nil {
		log.Printf("rest logout failed: %v", err)
	}
}
