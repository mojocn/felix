package shadowos

import "testing"

func TestCfIP_CheckReachableIps(t *testing.T) {
	cf, err := NewCfIP()
	if err != nil {
		t.Errorf("NewCfIP() error = %v", err)
		return
	}
	cf.CheckReachableIps()
}
