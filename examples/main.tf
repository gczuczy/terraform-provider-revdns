
terraform {
  required_providers {
    revdns = {
      source = "gczuczy/revdns"
    }
  }
}

data "revdns_zone" "revzone1" {
  cidr = "10.0.0.0/8"
}

data "revdns_zone" "revzone2" {
  cidr = "192.168.0.0/16"
}

data "revdns_zone" "revzone3" {
  cidr = "172.16.0.0/12"
}

data "revdns_zone" "revzone6-1" {
  cidr = "2001:db8::/32"
}

data "revdns_record" "record1" {
  zoneid = "10.0.0.0/8"
  address = "10.254.3.0"
  domain = "example.com"
  hostname = "test1"
}

data "revdns_record" "record2" {
  zoneid = "192.168.0.0/16"
  address = "192.168.1.2"
  domain = "example.com"
  hostname = "test2"
}


data "revdns_record" "record3" {
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

output "rev1" {
  description = "Rev1 zone data source"
  value = data.revdns_zone.revzone1
}

output "rev2" {
  description = "Rev2 zone data source"
  value = data.revdns_zone.revzone2
}

output "rev3" {
  description = "Rev3 zone data source"
  value = data.revdns_zone.revzone3
}

output "rev6-1" {
  description = "Rev6-1 zone data source"
  value = data.revdns_zone.revzone6-1
}

output "record1" {
  description = "Address test 1"
  value = data.revdns_record.record1
}

output "record2" {
  description = "Address test 1"
  value = data.revdns_record.record2
}

output "record3" {
  description = "Address test 1"
  value = data.revdns_record.record3
}

output "record6-1" {
  description = "Address test 6-1"
  value = data.revdns_record.record6-1
}
