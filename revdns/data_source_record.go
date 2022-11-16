package revdns

import (
	"net"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRecordRead,
		Schema: map[string]*schema.Schema{
			"zoneid": {
				Type: schema.TypeString,
				Required: true,
			},
			"address": {
				Type: schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type: schema.TypeString,
				Required: true,
			},
			"domain": {
				Type: schema.TypeString,
				Required: true,
			},
			"fqdn": {
				Type: schema.TypeString,
				Computed: true,
			},
			"record_short": {
				Type: schema.TypeString,
				Computed: true,
			},
			"record_fqdn": {
				Type: schema.TypeString,
				Computed: true,
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

	var ipslice []string = make([]string, 4-subnet.octets)
	for i:=uint(3); i>=subnet.octets; i-- {
		ipslice[3-i] = fmt.Sprintf("%v", ipaddress.To4()[i])
	}
	revip := strings.Join(ipslice, `.`)

	d.Set("record_short", revip)
	d.Set("record_fqdn", fmt.Sprintf("%v.%v", revip, subnet.zone))

	fqdn := fmt.Sprintf("%v.%v", hostname, domain)
	if fqdn[len(fqdn)-1] != '.' {
		fqdn += "."
	}
	d.Set("fqdn", fqdn)

	return nil
}

func ParseAddress(address string, cidr string) (*net.IP, error) {
	ipaddr := net.ParseIP(address)
	if err != nil {
		return nil, err
	}

	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	if !net.Contains(ipaddr) {
		return nil, fmt.Errorf("CIDR(%v) doesn't contain address(%v)", cidr, address)
	}

	return &ipaddr, nil
}
