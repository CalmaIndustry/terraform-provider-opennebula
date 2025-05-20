package opennebula

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	ver "github.com/hashicorp/go-version"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Configuration struct {
	OneVersion *ver.Version
	Controller *goca.Controller
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	username, ok := d.GetOk("username")
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "username should be defined",
		})
		return nil, diags
	}

	password, ok := d.GetOk("password")
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "password should be defined",
		})
		return nil, diags
	}

	endpoint, ok := d.GetOk("endpoint")
	if !ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "endpoint should be defined",
		})
		return nil, diags
	}

	insecure := d.Get("insecure")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure.(bool)},
	}

	oneClient := goca.NewClient(goca.NewConfig(username.(string),
		password.(string),
		endpoint.(string)),
		&http.Client{Transport: tr})

	versionStr, err := goca.NewController(oneClient).SystemVersion()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to get OpenNebula release number",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	version, err := ver.NewVersion(versionStr)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse OpenNebula version",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	log.Printf("[INFO] OpenNebula version: %s", versionStr)

	cfg := &Configuration{
		OneVersion: version,
	}

	cfg.Controller = goca.NewController(oneClient)

	return cfg, nil
}
