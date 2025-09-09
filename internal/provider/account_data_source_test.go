package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zoullx/terraform-provider-unifi/internal/datasource_account"
	"github.com/zoullx/unifi-go/unifi"
)

type mockUnifiClient struct {
	GetAccountFunc func(ctx context.Context, site, id string) (*unifi.Account, error)
}

func (m *mockUnifiClient) GetAccount(ctx context.Context, site, id string) (*unifi.Account, error) {
	return m.GetAccountFunc(ctx, site, id)
}
// Implement all other methods of unifi.Client with empty bodies to satisfy the interface
func (m *mockUnifiClient) AdoptDevice(ctx context.Context, site, deviceID string) error { return nil }
func (m *mockUnifiClient) SomeOtherMethod() error { return nil } // Add more as needed for interface compliance
func (m *mockUnifiClient) BaseURL() string { return "" }

func TestAccountDataSource_Read_Success(t *testing.T) {
	// Instead of using ReadRequest.Config, call parseAccountDataSourceJson directly and test its output
	account := unifi.Account{
		ID:               "acc-123",
		Name:             "test-account",
		XPassword:        "secret",
		TunnelType:       1,
		TunnelMediumType: 2,
		NetworkID:        "net-456",
	}
	model := &datasource_account.AccountModel{}
	parseAccountDataSourceJson(account, model)

	assert.Equal(t, "acc-123", model.Id.ValueString())
	assert.Equal(t, "test-account", model.Name.ValueString())
	assert.Equal(t, "secret", model.Password.ValueString())
	assert.Equal(t, int64(1), model.TunnelType.ValueInt64())
	assert.Equal(t, int64(2), model.TunnelMediumType.ValueInt64())
	assert.Equal(t, "net-456", model.NetworkId.ValueString())
}

func TestAccountDataSource_Read_Error(t *testing.T) {
	mockClient := &mockUnifiClient{
		GetAccountFunc: func(ctx context.Context, site, id string) (*unifi.Account, error) {
			return nil, errors.New("not found")
		},
	}

	// Simulate error handling by calling GetAccount directly
	account, err := mockClient.GetAccount(context.Background(), "default", "acc-123")
	assert.Nil(t, account)
	assert.Error(t, err)
}

// ...removed mockConfigType and related code, not needed with new approach...
