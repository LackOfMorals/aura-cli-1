package credential_test

import (
	"testing"

	"github.com/neo4j/cli/aura-cli/internal/test/testutils"
)

func TestRemoveCredential(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetCredentialsValue("aura.credentials", []map[string]string{{"name": "test", "client-id": "testclientid", "client-secret": "testclientsecret"}})

	helper.ExecuteCommand("credential remove test")

	helper.AssertCredentialsValue("aura.credentials", "[]")
}
