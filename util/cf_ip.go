package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

//https://api.cloudflare.com/client/v4/ips
//https://api.cloudflare.com/client/v4/ips?networks=jdcloud

type CfIP struct {
	Ipv4Cidrs    []string `json:"ipv4_cidrs"`
	Ipv6Cidrs    []string `json:"ipv6_cidrs"`
	ReachableIPs chan string
}

func NewCfIP() (*CfIP, error) {
	resp, err := http.Get("https://api.cloudflare.com/client/v4/ips")
	if err != nil {
		return nil, fmt.Errorf("get cf ip failed %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get cf ip failed %s", resp.Status)
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read cf ip failed %w", err)
	}
	result := struct {
		Result struct {
			Ipv4Cidrs []string `json:"ipv4_cidrs"`
			Ipv6Cidrs []string `json:"ipv6_cidrs"`
			Etag      string   `json:"etag"`
		} `json:"result"`
		Success  bool          `json:"success"`
		Errors   []interface{} `json:"errors"`
		Messages []interface{} `json:"messages"`
	}{}
	err = json.Unmarshal(all, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cf ip failed %w", err)
	}

	return &CfIP{
		Ipv4Cidrs:    result.Result.Ipv4Cidrs,
		Ipv6Cidrs:    result.Result.Ipv4Cidrs,
		ReachableIPs: make(chan string),
	}, nil
}

func (ci CfIP) AllIps(fn func(ip, cidr string)) {
	for _, cidr := range append(ci.Ipv4Cidrs, ci.Ipv6Cidrs...) {
		ips, err := getIPsFromCIDR(cidr)
		if err != nil {
			fmt.Printf("Error parsing Cidr: %v\n", err)
			continue
		}
		for _, ip := range ips {
			fn(ip, cidr)
		}
	}
}

//func (ci CfIP) CheckReachableIps() {
//	maxWorkers := runtime.GOMAXPROCS(0) * 64
//	ips := ci.ips()
//	jobs := make(chan string, len(ips))
//	resultChan := make(chan string, len(ips))
//	var wg sync.WaitGroup
//	for i := 0; i < maxWorkers; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			for ip := range jobs {
//				if isReachable(ip, 443) {
//					resultChan <- ip
//				}
//			}
//		}()
//	}
//	for _, ip := range ci.ips() {
//		jobs <- ip
//	}
//	close(jobs)
//	wg.Wait()
//	var reachable []string
//	for ip := range resultChan {
//		reachable = append(reachable, ip)
//	}
//	fd, err := os.Create("cf_reachable_ips.txt")
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer fd.Close()
//	fd.Write([]byte(strings.Join(reachable, "\n")))
//}

func getIPsFromCIDR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// Remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func isReachable(ip string, port int) bool {
	timeout := time.Millisecond * 70
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		log.Printf("Error connecting to %s:%d: %v\n", ip, port, err)
		return false
	}
	conn.Close()
	return true
}
