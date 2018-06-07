package http

import (
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// labelsFromRequest returns request query parameters.
func labelsFromRequest(logger *zap.Logger, req *http.Request) map[string]string {
	values := req.URL.Query()
	labels := map[string]string{}
	for key := range values {
		switch strings.ToLower(key) {
		case "mac":
			// set mac if and only if it parses
			if hw, err := parseMAC(values.Get(key)); err == nil {
				labels[key] = hw.String()
			} else {
				if logger != nil {
					logger.Warn("ignoring unparseable MAC address",
						zap.Error(err),
						zap.String("mac", values.Get(key)),
					)
				}
			}
		default:
			// matchers don't use multi-value keys, drop later values
			labels[key] = values.Get(key)
		}
	}
	return labels
}

// parseMAC wraps net.ParseMAC with logging.
func parseMAC(s string) (net.HardwareAddr, error) {
	macAddr, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	return macAddr, err
}
