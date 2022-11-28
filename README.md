# terraform-provider-revdns

Provides data sources for reverse DNS zone configurations, for both IPv4 and IPv6 zones. The module calculates the revese DNS zone names and record names as a data source, providing an easy way to reference the data in the actual DNS provider's configurations.

At the time of the initial release, the goal is to simplify reverse DNS setup in a demo/PoC terraform project of mine. The implementation might not be complete, there might be bugs or incorrect behaviours at places. If you happen to find any, please kindly file either a bug report or feel free to open a PR on it.

# Configuration

First, as usual, include the providers:

  terraform {
    required_providers {
      revdns = {
        source = "gczuczy/revdns"
      }
    }
  }

Next, a zone has to be declared:

    data "revdns_zone" "revzone4-1" {
      cidr = "172.16.0.0/12"
    }
    
    data "revdns_zone" "revzone6-1" {
      cidr = "2001:db8::/32"
    }

After that, add records to the zones:

    data "revdns_record" "record4-1" {
      zoneid = "172.16.0.0/12"
      address = "172.17.42.2"
      domain = "example.com"
      hostname = "test3"
    }
    
    data "revdns_record" "record6-1" {
      zoneid = "2001:db8::/32"
      address = "2001:db8::ae1f:6bff:feb1:de80"
      domain = "example.com"
      hostname = "test1"
    }

The zone and record providers are publishing the following fields, to be referenced in the DNS provider's configuration:

    rev3 = {
      "cidr" = "172.16.0.0/12"
      "id" = "172.16.0.0/12"
      "netmask" = 12
      "parts" = 1
      "zone_name" = "172.in-addr.arpa."
    }
    rev6 = {
      "cidr" = "2001:db8::/32"
      "id" = "2001:db8::/32"
      "netmask" = 32
      "parts" = 0
      "zone_name" = "8.b.d.0.1.0.0.2.ip6.arpa."
    }
    record3 = {
      "address" = "172.17.42.2"
      "domain" = "example.com"
      "fqdn" = "test3.example.com."
      "hostname" = "test3"
      "id" = "172.17.42.2"
      "record_fqdn" = "2.42.17.172.in-addr.arpa."
      "record_short" = "2.42.17"
      "zoneid" = "172.16.0.0/12"
    }
    record6-1 = {
      "address" = "2001:db8::ae1f:6bff:feb1:de80"
      "domain" = "example.com"
      "fqdn" = "test1.example.com."
      "hostname" = "test1"
      "id" = "2001:db8::ae1f:6bff:feb1:de80"
      "record_fqdn" = "0.8.e.d.1.b.e.f.f.f.b.6.f.1.e.a.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa."
      "record_short" = "0.8.e.d.1.b.e.f.f.f.b.6.f.1.e.a.0.0.0.0.0.0.0.0"
      "zoneid" = "2001:db8::/32"
    }

# Data source descriptions

The provider has two data sources, one for the zones, one for the records.

## Data source revdns_zone

| name | type | Required | Computed | description |
| CIDR | string | yes | no | The CIDR address of the zone |
| zone_name | string | no | yes | The generated name of the reverse DNS zone |
| netmask | int | no | yes | The number of bits in the netmask |
| parts | int | no | yes | The number of octets(v4)/quadlets(v6) present in the reverse name |

The reverse zone's name will be the `revdns_zone.name.zone_name` field.

## Data source revdns_record

| name | type | Required | Computed | description |
| zoneid | string | yes | no | The id field of the corresponding revdns_zone data source. Technically the CIDR. |
| address | string | yes | no | The IP address for the record |
| domain | string | yes | no | The domain name for the hostname |
| hostname | string | yes | no | The short hostname belonging to the IP  |
| fqdn | string | no | yes | The generated fully qualified domain name for the host (hostname.domain.) |
| record_short | string | no | yes | The short reverse record for the IP |
| record_fqdn | string | no | yes | The fully qualified reverse record for the IP |

All fully qualified fields are generated with an ending dot to avoid any confusion or mishaps.

The reverse records will generally be like the following pattern:

    revdns_record.name.record_short IN PTR revdns_record.name.fqdn

So, the `fqdn` field is the PTR record's target, while the `record_fqdn` is the DNS record's name.
