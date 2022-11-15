package revdns

import (
	"net"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IPSubnet struct {
	cidr string
	zone string
	netmask uint
	octets uint
}

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type: schema.TypeString,
				Required: true,
			},
			"zone_name": {
				Type: schema.TypeString,
				Computed: true,
			},
			"netmask": {
				Type: schema.TypeInt,
				Computed: true,
			},
			"octets": {
				Type: schema.TypeInt,
				Computed: true,
			},
		},
	}
}


func dataSourceZoneRead(d *schema.ResourceData, meta interface{}) error {
	cidr := d.Get("cidr").(string)
	subnet, err := ParseCIDR(cidr)

	if err != nil  {
		return err
	}

	d.SetId(cidr)
	d.Set("netmask", subnet.netmask)
	d.Set("octets", subnet.octets)
	d.Set("zone_name", subnet.zone)

	return nil
}

func ParseCIDR(cidr string) (*IPSubnet, error) {
	ipv4Addr, ipv4Net, err := net.ParseCIDR(cidr)

	if err != nil {
		return nil, err
	}
	ones, bits := ipv4Net.Mask.Size()

	if ones == 0 && bits == 0 {
		return nil, fmt.Errorf("Unable to parse subnet")
	}

	reversed := [...]byte{ipv4Addr.To4()[3], ipv4Addr.To4()[2], ipv4Addr.To4()[1], ipv4Addr.To4()[0]}

	var zname string
	var octets uint

	if ones < 8 {
		return nil, fmt.Errorf("Netmasks less than 8 not supported")
	} else if ones < 16 {
		octets = 1
		// 8 to 16
		zname = fmt.Sprintf("%v.in-addr.arpa.", reversed[3])
	} else if ones < 24 {
		octets = 2
		zname = fmt.Sprintf("%v.%v.in-addr.arpa.", reversed[2], reversed[3])
	} else if ones < 32 {
		octets = 3
		zname = fmt.Sprintf("%v.%v.%v.in-addr.arpa.", reversed[1], reversed[2], reversed[3])
	} else {
		return nil, fmt.Errorf("Netmask of 32 is not supported")
	}

	return &IPSubnet{cidr, zname, uint(ones), octets}, nil
}
