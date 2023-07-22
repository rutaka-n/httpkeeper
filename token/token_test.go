package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Parallel()
	secret := []byte("test_secret")
	service := "test service"
	client := "test client"
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0IHNlcnZpY2UiLCJhdWQiOlsidGVzdCBjbGllbnQiXSwiZXhwIjoxNjcyNTc1MTMyfQ.2AAIMPtT3d11iIeRIkqZkTq_yF7fbMos6lgDlFUvv88"

	actualToken, err := Generate(secret, service, client, time.Date(2023, 01, 01, 12, 12, 12, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, expectedToken, actualToken)
}

func TestValidate(t *testing.T) {
	t.Parallel()
	var secret []byte = []byte("test_secret")
	service := "test service"
	cases := []struct {
		desc      string
		client    string
		secret    []byte
		service   string
		expiresAt time.Time
		isValid   bool
	}{
		{
			desc:      "valid token",
			client:    "test client",
			secret:    secret,
			service:   service,
			expiresAt: time.Now().Add(1 * time.Minute),
			isValid:   true,
		},
		{
			desc:      "token generated with different secret",
			client:    "test client",
			secret:    []byte("another secret"),
			service:   service,
			expiresAt: time.Now().Add(1 * time.Minute),
			isValid:   false,
		},
		{
			desc:      "expired token",
			client:    "test client",
			secret:    secret,
			service:   service,
			expiresAt: time.Now().Add(-1 * time.Minute),
			isValid:   false,
		},
	}

	for i := range cases {
		tc := cases[i]
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			token, err := Generate(tc.secret, tc.service, tc.client, tc.expiresAt)
			require.NoError(t, err)

			actualPayload, err := Validate(secret, token)
			if tc.isValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				return
			}
			require.Equal(t, tc.service, actualPayload.Issuer)
			require.NotEmpty(t, actualPayload.Audience)
			require.Equal(t, tc.client, actualPayload.Audience[0])
		})
	}
}
