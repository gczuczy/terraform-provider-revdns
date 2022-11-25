package revdns

import (
	"testing"
)

type ZoneRecordTestInputs struct {
	address string
	error bool
}

type RecordTestInputs struct {
	cidr string
	addresses []ZoneRecordTestInputs
}

var recordCases []RecordTestInputs = []RecordTestInputs{
	{cidr: "192.168.1.0/24",
		addresses: []ZoneRecordTestInputs{
			{address: "192.168.1.4", error: false},
			{address: "192.168.2.4", error: true},
		},
	},
	{cidr: "192.168.0.0/16",
		addresses: []ZoneRecordTestInputs{
			{address: "192.168.1.4", error: false},
			{address: "192.162.2.4", error: true},
		},
	},
	{cidr: "172.16.0.0/12",
		addresses: []ZoneRecordTestInputs{
			{address: "172.16.9.5", error: false},
			{address: "172.11.2.4", error: true},
		},
	},
	{cidr: "10.0.0.0/8",
		addresses: []ZoneRecordTestInputs{
			{address: "10.53.23.67", error: false},
			{address: "192.168.2.4", error: true},
		},
	},
	{cidr: "2001:db8::/32",
		addresses: []ZoneRecordTestInputs{
			{address: "2001:db8::ae1f:6bff:feb1:de80", error: false},
			{address: "2001:cb8::AE1F:6BFF:FEB1:DE80", error: true},
			{address: "2g01:db8::AE1F:6BFF:FEB1:DE80", error: true},
		},
	},
}

func TestParseAddress(t *testing.T) {

	for _,zone := range recordCases {
		_, err := ParseCIDR(zone.cidr)
		// also checkin on the errors here to be on the safe side
		if err != nil {
			t.Fatalf("ParseCIDR(%v) failed: %v", zone.cidr, err)
			return
		}

		for _, address := range zone.addresses {
			_, err := ParseAddress(address.address, zone.cidr)
			if !address.error && err != nil {
				t.Fatalf("Unable to parse IP %v: %v", address.address, err)
				return
			}
			if address.error && err == nil {
				t.Fatalf("Address(%v) should be invalid for cidr(%v)", address.address, zone.cidr)
				return
			}
		}
	}
}
