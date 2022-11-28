package revdns

import (
	"net"
	"fmt"
	"strings"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRecordRead,
		Schema: map[string]*schema.Schema{
			"zoneid": {
				Type: schema.TypeString,
				Required: true,
				Description: "The revdns_zone's ID where this record belongs to",
			},
			"address": {
				Type: schema.TypeString,
				Required: true,
				Description: "IP address to generate the reverse record for",
			},
			"hostname": {
				Type: schema.TypeString,
				Required: true,
				Description: "Short hostname associated with the IP",
			},
			"domain": {
				Type: schema.TypeString,
				Required: true,
				Description: "The domain name to be used to generate the FQDN with",
			},
			"fqdn": {
				Type: schema.TypeString,
				Computed: true,
				Description: "The Fully Qualified hostname",
			},
			"record_short": {
				Type: schema.TypeString,
				Computed: true,
				Description: "Non-qualified reverse record name",
			},
			"record_fqdn": {
				Type: schema.TypeString,
				Computed: true,
				Description: "Fully qualified reverse record name",
			},
		},
	}
}


func dataSourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	address := d.Get("address").(string)
	zoneid := d.Get("zoneid").(string)
	domain := d.Get("domain").(string)
	hostname := d.Get("hostname").(string)

	subnet, err := ParseCIDR(zoneid)
	if err != nil {
		return err
	}

	ipaddress, err := ParseAddress(address, zoneid)
	if err != nil {
		return err
	}

	d.SetId(address)

	revname, err := RevName(*ipaddress, subnet.netmask)
	if err != nil {
		return err
	}

	d.Set("record_short", revname)
	d.Set("record_fqdn", fmt.Sprintf("%v.%v", revname, subnet.zone))

	fqdn := fmt.Sprintf("%v.%v", hostname, domain)
	if fqdn[len(fqdn)-1] != '.' {
		fqdn += "."
	}
	d.Set("fqdn", fqdn)

	return nil
}

func ParseAddress(address string, cidr string) (*net.IP, error) {
	ipaddr := net.ParseIP(address)
	if ipaddr == nil {
		return nil, fmt.Errorf("Unable to parse IP address %v", address)
	}

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	if !ipnet.Contains(ipaddr) {
		return nil, fmt.Errorf("CIDR(%v) doesn't contain address(%v)", cidr, address)
	}

	return &ipaddr, nil
}

func RevName(ipaddress net.IP, netmask uint) (string, error) {
	addrv4 := ipaddress.To4()
	addrv6 := ipaddress.To16()
	isv6 := false
	if addrv4 != nil {
		isv6 = false
	} else if addrv4 == nil && addrv6 != nil {
		isv6 = true
	} else {
		return "", fmt.Errorf("Neither v4 or v6")
	}

	var recname string
	if isv6 {
		mask_quadlets := uint(math.Ceil(float64(netmask)/4))
		addr_quadlets := uint(math.Ceil(float64(128-netmask)/4))
		strnet := ""
		addrv6 := ipaddress.To16()
		fmt.Printf("RevName6(%v, %v):\n - mask quadlets: %v\n - addr quadlets: %v\n",
			ipaddress, netmask, mask_quadlets, addr_quadlets);
		if addrv6 == nil {
			return "", fmt.Errorf("Not a valid v6 address")
		}

		// write the whole v6 address as a string
		// each character is a quadlet (4 bits) here
		for i:=uint(0); i<=uint(len(addrv6))-1; i++ {
			strnet += fmt.Sprintf("%02x", addrv6[i])
		}
		strbytes := []byte(strnet)

		var recbytes []string = make([]string, addr_quadlets)
		for i,j := uint(len(strbytes)), 0; i > mask_quadlets; i,j = i-1, j+1 {
			recbytes[j] = string(strbytes[i-1])
		}
		recname = strings.Join(recbytes, ".")
	} else {
		octets := netmask/8
		var ipslice []string = make([]string, 4-octets)
		for i:=uint(3); i>=octets; i-- {
			ipslice[3-i] = fmt.Sprintf("%v", addrv4[i])
		}
		recname = strings.Join(ipslice, `.`)
	}
	return recname, nil
}
