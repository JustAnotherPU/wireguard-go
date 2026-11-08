package netstack

import (
	"net/netip"
	"testing"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
)

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

// The cap reaches the stack that CreateNetTUN builds: WG_TCP_SNDBUF=0 keeps gVisor's own default (upstream
// behaviour), unset gives the 128 KiB cap, a value gives that value as both the default and the maximum.
func TestCreateNetTUNSendBuf(t *testing.T) {
	upstream := tcp.DefaultSendBufferSize
	for _, c := range []struct {
		env              string
		set              bool
		wantDef, wantMax int
	}{{"0", true, upstream, tcp.MaxBufferSize}, {"", false, 128 << 10, 128 << 10}, {"65536", true, 65536, 65536}} {
		t.Run("env="+c.env, func(t *testing.T) { // a subtest: each t.Setenv is undone before the next case
			if c.set {
				t.Setenv("WG_TCP_SNDBUF", c.env)
			}
			dev, tnet, err := CreateNetTUN([]netip.Addr{netip.MustParseAddr("10.0.0.2")}, nil, 1280)
			if err != nil {
				t.Fatal(err)
			}
			defer dev.Close()
			var got tcpip.TCPSendBufferSizeRangeOption
			if e := (*netTun)(tnet).stack.TransportProtocolOption(tcp.ProtocolNumber, &got); e != nil {
				t.Fatal(e)
			}
			if got.Default != c.wantDef || got.Max != c.wantMax {
				t.Errorf("default %d max %d, want %d %d", got.Default, got.Max, c.wantDef, c.wantMax)
			}
		})
	}
}
