package scanner

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseScanReturnsAccessPoints(t *testing.T) {
	// Output with two valid access points and one invalid access point (no SSID). There
	// may never be information missing for an AP, but test just in case
	output := `
BSS mac1(on testwlan0)
        last seen: 10s [boottime]
        signal: 3.30 dBm
        last seen: 20 ms ago
        SSID: ssid1
        Supported rates: test rates

BSS mac2(on testwlan0)
        last seen: 10s [boottime]
        signal: -2.00 dBm
        last seen: 10 ms ago
        SSID: ssid2
        Supported rates: test rates

BSS mac3(on testwlan0)
        last seen: 10s [boottime]
        signal: 5.00 dBm
        last seen: 100 ms ago
        Supported rates: test rates`

	output = strings.TrimSpace(output)

	iw := &IWCommand{}
	accessPoints := iw.parseScan(output)

	require.Len(t, accessPoints, 2)

	// Access point 1
	accessPoint1 := accessPoints[0]
	assert.Equal(t, accessPoint1.Address, "mac1")
	assert.Equal(t, accessPoint1.SSID, "ssid1")
	assert.Equal(t, accessPoint1.Strength, 3.3)

	// Access point 2
	accessPoint2 := accessPoints[1]
	assert.Equal(t, accessPoint2.Address, "mac2")
	assert.Equal(t, accessPoint2.SSID, "ssid2")
	assert.Equal(t, accessPoint2.Strength, -2.0)
}
