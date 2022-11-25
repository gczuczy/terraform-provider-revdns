package revdns

import (
	"testing"
)

type ZoneTestInputs struct {
	cidr string
	error bool
}

var zoneCases []ZoneTestInputs = []ZoneTestInputs{
	{cidr: "192.168.1.3/24", error: false},
	{cidr: "192.168.1.0/24", error: false},
	{cidr: "192.168.0.0/16", error: false},
	{cidr: "172.16.0.0/12", error: false},
	{cidr: "10.0.0.0/8", error: false},
	{cidr: "192.168.1.3/a24", error: true},
	{cidr: "10.0.0.0/0", error: true},
	{cidr: "10.1.2.0/32", error: true},
	{cidr: "2001:db8::/32", error: false},
	{cidr: "20h1:db8::/33", error: true},
}

func TestParseCIDR(t *testing.T) {

	for _,data := range zoneCases {
		_, err := ParseCIDR(data.cidr)
		if !data.error && err != nil {
			t.Fatalf("ParseCIDR(%v) failed: %v", data.cidr, err)
			return
		}
		if data.error && err == nil {
			t.Fatalf("ParseCIDR(%v) accepted erroneous input", data.cidr)
			return
		}
	}
}
