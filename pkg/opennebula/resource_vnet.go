package opennebula

import (
	"context"
	"fmt"
	"strconv"

	"github.com/OpenNebula/one/src/oca/go/src/goca/parameters"
	vn "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/virtualnetwork"
	vnk "github.com/OpenNebula/one/src/oca/go/src/goca/schemas/virtualnetwork/keys"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenNebulaVNET() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenNebulaVNETCreate,
		ReadContext:   resourceOpenNebulaVNETRead,
		UpdateContext: resourceOpennebulaVirtualNetworkUpdate,
		DeleteContext: resourceOpenNebulaVNETDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnmad": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"physical_device": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceOpenNebulaVNETCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := meta.(*Configuration)
	controller := config.Controller

	vnname := d.Get("name").(string)
	vlanid := d.Get("vlan_id").(string)
	vn_mad := d.Get("vnmad").(string)
	physicaldevice := d.Get("physical_device").(string)

	tpl := vn.NewTemplate()

	tpl.Add(vnk.Name, vnname)
	tpl.Add(vnk.VlanID, vlanid)
	tpl.Add(vnk.VNMad, vn_mad)
	tpl.Add(vnk.PhyDev, physicaldevice)

	vnetID, err := controller.VirtualNetworks().Create(tpl.String(), 100)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create the virtual network",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(fmt.Sprintf("%v", vnetID))

	return resourceOpenNebulaVNETRead(ctx, d, meta)
}

func resourceOpenNebulaVNETRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := meta.(*Configuration)
	controller := config.Controller

	vnetID := d.Id()

	vnet, err := controller.VirtualNetworks().Info()
	if err != nil {
		// If not found, mark resource as removed
		d.SetId("")
		return nil
	}
	fmt.Print(vnet)
	d.SetId(fmt.Sprintf("%v", vnetID))

	return diags
}

func resourceOpennebulaVirtualNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Configuration)
	controller := config.Controller

	imgID, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		fmt.Print(err)
	}

	vnc := controller.VirtualNetwork(int(imgID))

	vnInfos, _ := vnc.Info(false)

	tpl := vnInfos.Template

	if d.HasChange("physical_device") {
		tpl.Del(string(vnk.PhyDev))
		physicaldevice := d.Get("physical_device").(string)
		if len(physicaldevice) > 0 {
			tpl.Add(vnk.PhyDev, physicaldevice)
		}
		vnc.Update(tpl.String(), parameters.Replace)
	}
	return nil
}

func resourceOpenNebulaVNETDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*Configuration)
	controller := config.Controller

	imgID, err := strconv.ParseUint(d.Id(), 10, 0)
	if err != nil {
		fmt.Print(err)
	}

	vnc := controller.VirtualNetwork(int(imgID))
	vnc.Delete()

	return nil
}
