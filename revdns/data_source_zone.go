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
	isv6 bool
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
	ipaddr, ipnet, err := net.ParseCIDR(cidr)

	if err != nil {
		return nil, err
	}

	addrv4 := ipaddr.To4()
	addrv6 := ipaddr.To16()
	isv6 := false
	if addrv4 != nil {
		isv6 = false
	} else if addrv4 == nil && addrv6 != nil {
		isv6 = true
	} else {
		return nil, fmt.Errorf("Neither v4 or v6")
	}

	ones, bits := ipnet.Mask.Size()
	if ones == 0 && bits == 0 {
		return nil, fmt.Errorf("Unable to parse subnet")
	}

	var zname string
	var octets uint

	if !isv6 {
		reversed := [...]byte{ipaddr.To4()[3], ipaddr.To4()[2], ipaddr.To4()[1], ipaddr.To4()[0]}

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
	} else {
		octets = uint(ones)/4
		strnet := ""
		for i:=uint(0); i<octets; i++ {
			strnet += fmt.Sprintf("%x", addrv6[i])
		}
		strbytes := []byte(strnet)
		zname = ""
		for i := len(strbytes)-1; i>=0; i-- {
			zname += fmt.Sprintf("%v.", string(strbytes[i]))
		}
		zname += "ip6.arpa."
	}

	return &IPSubnet{cidr, zname, uint(ones), octets, isv6}, nil
}
