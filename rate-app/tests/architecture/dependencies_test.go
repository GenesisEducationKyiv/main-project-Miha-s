package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestShouldNotDependOn(t *testing.T) {
	t.Run("Configuration should not depend on anything", func(t *testing.T) {
		archtest.Package(t, "btc-test-task/internal/common/...").
			ShouldNotDependOn("btc-test-task/internal/...")
	})

	t.Run("currencyrate should not depend on anything except common", func(t *testing.T) {
		archtest.Package(t, "btc-test-task/internal/currencyrate/...").
			Ignoring("btc-test-task/internal/common/...").
			ShouldNotDependOn("btc-test-task/internal/...")
	})

	t.Run("email should not depend on anything except common", func(t *testing.T) {
		archtest.Package(t, "btc-test-task/internal/email/...").
			Ignoring("btc-test-task/internal/common/...").
			ShouldNotDependOn("btc-test-task/internal/...")
	})

	t.Run("repository should not depend on anything except common", func(t *testing.T) {
		archtest.Package(t, "btc-test-task/internal/repository/...").
			Ignoring("btc-test-task/internal/common/...").
			ShouldNotDependOn("btc-test-task/internal/...")
	})

	t.Run("server should not depend on lifecycle", func(t *testing.T) {
		archtest.Package(t, "btc-test-task/internal/server/...").
			ShouldNotDependOn("btc-test-task/internal/lifecycle/...")
	})
}
