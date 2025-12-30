package cidr

import (
	"net"
	"net/http"
	"strings"
)

// WithCheckCIDR - middleware для проверки CIDR.
func WithCheckCIDR(trustedSubnet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
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

			_, network, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if network.Contains(ip) {
				next.ServeHTTP(w, r)
			}
			w.WriteHeader(http.StatusForbidden)
		})
	}
}
