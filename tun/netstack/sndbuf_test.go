package netstack

import "testing"

func TestTCPSendBuf(t *testing.T) {
	for _, c := range []struct {
		env  string
		set  bool
		want int
	}{{"", false, 128 << 10}, {"0", true, 0}, {"262144", true, 262144}, {"100", true, 4096}, {"x", true, 128 << 10}, {"-5", true, 128 << 10}} {
		if c.set {
			t.Setenv("WG_TCP_SNDBUF", c.env)
		}
		if got := tcpSendBuf(); got != c.want {
			t.Errorf("WG_TCP_SNDBUF=%q set=%v: got %d, want %d", c.env, c.set, got, c.want)
		}
	}
}
