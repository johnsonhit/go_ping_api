package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sparrc/go-ping"
	"net"
	"net/http"
	"strconv"
	"time"
)


func main() {
	gin.SetMode("release")
	r := gin.Default()
	r.GET("/", rhelp)
	r.GET("/ping", rping)
	r.Run(":61024")
}


func rhelp(c *gin.Context) {
	c.String(http.StatusOK, "example: /ping?host=baidu.com&port=80")
}

func rping(c *gin.Context) {
	ahost := c.DefaultQuery("host", "127.0.0.1")
	aport,err := strconv.Atoi(c.DefaultQuery("port",  "80"))
	if err != nil {
		aport = 80
	}
	aips,_ := net.LookupIP(ahost)
	icmpt,isuccess := icmpPing(ahost)
	tcpt,tsuccess := tcpPing(ahost, aport)
	c.JSON(http.StatusOK, gin.H{
		"host": ahost,
		"port": aport,
		"icmp": gin.H{
			"success": isuccess,
			"delay": icmpt },
		"tcp":  gin.H{
			"success": tsuccess,
			"delay": tcpt },
		"ip": aips,
	})
}



func icmpPing(icmpaddr string) (float32, bool) {
	pinger, err := ping.NewPinger(icmpaddr)
	if err != nil {
		return -1, false
	}
	pinger.Timeout = time.Second * 3
	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Run()
	stats := pinger.Statistics()
	return float32(float32(stats.MinRtt)/float32(time.Millisecond)), true
}


func tcpPing(tcpingaddr string, tcpport int) (float32, bool) {
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", tcpingaddr+":"+strconv.Itoa(tcpport), time.Second*3)
	endTime := time.Now()
	if err != nil {
		return -1, false
	} else {
		defer conn.Close()
		return float32(float32(endTime.Sub(startTime))/float32(time.Millisecond)), true
	}
}