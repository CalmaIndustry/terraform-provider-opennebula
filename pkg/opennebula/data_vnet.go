package opennebula

import (
	"context"
	"strconv"

	vnetSc "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/virtualnetwork"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataOpenNebulaVNET() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataresourceOpenNebulaVNETRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "Id of the virtual network",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Virtual Network",
			},
		},
	}
}

func filterVirtualNetworks(vnets []vnetSc.VirtualNetwork, id int, name string, nameOk bool) []*vnetSc.VirtualNetwork {

	match := make([]*vnetSc.VirtualNetwork, 0, 1)
	for i, vnet := range vnets {
		if id != -1 && vnet.ID != id {
			continue
		}
		if nameOk && vnet.Name != name {
			continue
		}
		match = append(match, &vnets[i])
	}
	return match
}

func dataresourceOpenNebulaVNETRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*Configuration)
	controller := config.Controller
	vnets, _ := controller.VirtualNetworks().Info()

	id := d.Get("id")
	name, nameOk := d.GetOk("name")

	match := filterVirtualNetworks(vnets.VirtualNetworks, id.(int), name.(string), nameOk)

	vnet := match[0]

	d.SetId(strconv.FormatInt(int64(vnet.ID), 10))
	d.Set("name", vnet.Name)

	return nil
}
