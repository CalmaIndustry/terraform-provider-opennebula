package opennebula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENNEBULA_ENDPOINT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENNEBULA_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENNEBULA_PASSWORD", nil),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable TLS validation",
				DefaultFunc: schema.EnvDefaultFunc("OPENNEBULA_INSECURE", true),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"opennebula_vnet": resourceOpenNebulaVNET(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
