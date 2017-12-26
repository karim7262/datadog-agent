package network

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/pkg/collector/check"
)

func TestString(t *testing.T) {
	n := newNetDevCheck()
	assert.Equal(t, "net/dev", n.String())
}

func TestConfigure(t *testing.T) {
	testCases := []struct {
		testName           string
		netdevCheck        check.Check
		cfg                []byte
		initCfg            []byte
		expectedNetDevPath string
	}{
		{
			testName:    "without Prefix",
			netdevCheck: newNetDevCheck(),
			cfg: []byte(`
proc_prefix: ""
`),
			initCfg:            []byte(``),
			expectedNetDevPath: "/proc/1/net/dev",
		},
		{
			testName:    "with Prefix",
			netdevCheck: newNetDevCheck(),
			cfg: []byte(`
proc_prefix: "/host"
`),
			initCfg:            []byte(``),
			expectedNetDevPath: "/host/proc/1/net/dev",
		},
		{
			testName:    "with Prefix with extra slash",
			netdevCheck: newNetDevCheck(),
			// Add a / at the end
			cfg: []byte(`
proc_prefix: "/host/"
`),
			initCfg:            []byte(``),
			expectedNetDevPath: "/host/proc/1/net/dev",
		},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("case %d %s", i, test.testName), func(t *testing.T) {
			err := test.netdevCheck.Configure(test.cfg, test.initCfg)
			assert.Nil(t, err)
			assert.Equal(t, 24, len(test.netdevCheck.ID()))

			// type to NetDevCheck to access to dedicated methods
			netdevCheck := test.netdevCheck.(*NetDevCheck)
			assert.Equal(t, test.expectedNetDevPath, netdevCheck.getNetDevPath())
		})
	}
}
