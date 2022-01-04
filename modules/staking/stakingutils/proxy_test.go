package stakingutils

import (
	"testing"

	"github.com/stretchr/testify/suite"

	testutils "github.com/HarleyAppleChoi/junomum/modules/utils"
	_ "github.com/proullon/ramsql/driver"
)

func TestStakingProxyTestSuite(t *testing.T) {
	suite.Run(t, new(StakingProxyTestSuite))
}

type StakingProxyTestSuite struct {
	testutils.ProxyTestSuite
}
