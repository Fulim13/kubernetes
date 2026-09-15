package main

import (
	"net/http"
	"os"

	"errors"
	"net"

	"github.com/gin-gonic/gin"
)

// Get local NIC IP (private IP)
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP address
		isIpNet bool
	)
	// Get all NICs
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// Take the first non-lo NIC IP
	for _, addr = range addrs {
		// Check whether this network address is an IP address: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet {
			if !ipNet.IP.IsLoopback() {
				if ipNet.IP.IsPrivate() { // private address
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
		name := ctx.Query("name")       // get value from request params
		project := os.Getenv("project") // get value from environment variable
		ip, _ := GetLocalIP()
		ctx.HTML(http.StatusOK, "home.html", gin.H{"name": name, "project": project, "ip": ip})
	})
	port := os.Getenv("port") // get value from environment variable
	engine.Run("0.0.0.0:" + port)
}
