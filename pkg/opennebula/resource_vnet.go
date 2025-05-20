package opennebula

import (
	"context"
	"fmt"

	vn "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/virtualnetwork"
	vnk "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/virtualnetwork/keys"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenNebulaVNET() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenNebulaVNETCreate,
		ReadContext:   resourceOpenNebulaVNETRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: false,
			},
		},
	}
}

func resourceOpenNebulaVNETCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := meta.(*Configuration)
	controller := config.Controller

	vnname := d.Get("name").(string)
	vlanid := d.Get("vlan_id")
	tpl := vn.NewTemplate()
	tpl.Add(vnk.Name, vnname)
	tpl.Add(vnk.VlanID, vlanid.(string))

	vnetID, err := controller.VirtualNetworks().Create(tpl.String(), 100)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create the virtual network",
			Detail:   err.Error(),
		})
		return diags
	}
	fmt.Print(vnetID)

	return resourceOpenNebulaVNETRead(ctx, d, meta)
}

func resourceOpenNebulaVNETRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Use client to get VM info, update d.Set("name", ...)

	return nil
}
