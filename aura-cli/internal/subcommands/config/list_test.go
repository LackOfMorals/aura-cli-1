package config_test

import (
	"fmt"
	"testing"

	"github.com/neo4j/cli/aura-cli/internal/test/testutils"
	"github.com/neo4j/cli/common/clicfg"
)

func TestListConfig(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.OverwriteConfig("{}")

	helper.ExecuteCommand("config list")

	helper.AssertOutJson(fmt.Sprintf(`{"auth-url": "%s","base-url": "%s","beta-enabled": false,"output": "default"}`, clicfg.DefaultAuraAuthUrl, clicfg.DefaultAuraBaseUrl))
}
