package main

import (
	"net/http"
	"os"

	"errors"
	"net"

	"github.com/gin-gonic/gin"
)

// Get the local network interface IP (private IP)
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP address: IPv4 or IPv6
		isIpNet bool
	)
	// Get all network interfaces
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// Get the first non-loopback network interface IP
	for _, addr = range addrs {
		// Check if the network address is an IP address: IPv4 or IPv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet {
			if !ipNet.IP.IsLoopback() {
				if ipNet.IP.IsPrivate() { // Private network address
					// Skip IPv6
					if ipNet.IP.To4() != nil {
						ipv4 = ipNet.IP.String()
						return
					}
				}
			}
		}
	}

	err = errors.New("ERR_NO_LOCAL_IP_FOUND")
	return
}

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")
	engine.GET("/", func(ctx *gin.Context) {
		name := ctx.Query("name")       // Get value from request query parameter
		project := os.Getenv("project") // Get value from environment variable
		ip, _ := GetLocalIP()
		ctx.HTML(http.StatusOK, "home.html", gin.H{"name": name, "project": project, "ip": ip})
	})
	port := os.Getenv("port") // Get value from environment variable
	engine.Run("0.0.0.0:" + port)
}
