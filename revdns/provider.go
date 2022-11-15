package revdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"revdns_zone": dataSourceZone(),
			"revdns_record": dataSourceRecord(),
		},
	}
}
