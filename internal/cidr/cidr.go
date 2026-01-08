package cidr

import (
	"net"
	"net/http"
	"strings"
)

// WithCheckCIDR - middleware для проверки CIDR.
func WithCheckCIDR(trustedSubnet string, useHeader bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_, network, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			var ip net.IP
			if !useHeader {
				addr := r.RemoteAddr
				ipStr, _, err := net.SplitHostPort(addr)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				ip = net.ParseIP(ipStr)
				if ip == nil {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			} else {
				ipStr := r.Header.Get("X-Real-IP")
				ip = net.ParseIP(ipStr)
				if ip == nil {
					ips := r.Header.Get("X-Forwarded-For")
					ipStrs := strings.Split(ips, ",")
					ipStr = ipStrs[0]
					ip = net.ParseIP(ipStr)
				}

				if ip == nil {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
			if network.Contains(ip) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		})
	}
}
