package broker

import (
	"bytes"
	"crypto/tls"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type CertAuthHook struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *CertAuthHook) ID() string {
	return "cert-auth"
}

// Provides indicates which hook methods this hook provides.
func (h *CertAuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
	}, []byte{b})
}

// OnConnectAuthenticate returns true/allowed for all requests.
func (h *CertAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	tlsConn, ok := cl.Net.Conn.(*tls.Conn)
	if !ok {
		return false
	}

	state := tlsConn.ConnectionState()

	for _, cert := range state.PeerCertificates {
		if cert.Issuer.CommonName == "Mochi MQTT Root CA" && strings.HasPrefix(cert.Subject.CommonName, "cluster") {
			return true
		}
	}
	return false
}

// OnACLCheck returns true/allowed for all checks.
func (h *CertAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	return true
}
