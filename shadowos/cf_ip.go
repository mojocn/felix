package shadowos

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
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

func (ci CfIP) ips() []string {
	var allIps []string
	for _, cidr := range ci.Ipv4Cidrs {
		ips, err := getIPsFromCIDR(cidr)
		if err != nil {
			fmt.Printf("Error parsing CIDR: %v\n", err)
			continue
		}
		allIps = append(allIps, ips...)
	}
	//for _, cidr := range ci.Ipv6Cidrs {
	//	ips, err := getIPsFromCIDR(cidr)
	//	if err != nil {
	//		fmt.Printf("Error parsing CIDR: %v\n", err)
	//		continue
	//	}
	//	allIps = append(allIps, ips...)
	//}
	return allIps
}

func (ci CfIP) CheckReachableIps() {
	maxWorkers := runtime.GOMAXPROCS(0)

	jobs := make(chan string)
	wg := sync.WaitGroup{}
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range jobs {
				if checkPort(ip, 443) {
					ci.ReachableIPs <- ip
				}
			}
		}()
	}
	go func() {
		for _, ip := range ci.ips() {
			jobs <- ip
		}
		close(jobs)
	}()

	go func() {
		for ii := range ci.ReachableIPs {
			log.Println(ii)
		}
	}()
	wg.Wait()
}

func getIPsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
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

func checkPort(ip string, port int) bool {
	timeout := time.Millisecond * 40
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
