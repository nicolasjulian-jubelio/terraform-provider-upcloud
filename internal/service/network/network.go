package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/UpCloudLtd/terraform-provider-upcloud/internal/utils"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/go-cty/cty"
)

func ResourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource represents an SDN private network that cloud servers from the same zone can be attached to.",
		ReadContext:   resourceNetworkRead,
		CreateContext: resourceNetworkCreate,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ip_network": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "A list of IP subnets within the network",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Description:  "The CIDR range of the subnet",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Description: "Is DHCP enabled?",
							Required:    true,
						},
						"dhcp_default_route": {
							Type:        schema.TypeBool,
							Description: "Is the gateway the DHCP default route?",
							Computed:    true,
							Optional:    true,
						},
						"dhcp_dns": {
							Type:        schema.TypeSet,
							Description: "The DNS servers given by DHCP",
							Computed:    true,
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.Any(validation.IsIPv4Address, validation.IsIPv6Address),
							},
						},
						"dhcp_routes": {
							Type:        schema.TypeSet,
							Description: "The additional DHCP classless static routes given by DHCP",
							Computed:    true,
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
						},
						"family": {
							Type:        schema.TypeString,
							Description: "IP address family",
							Required:    true,
							ValidateDiagFunc: func(v interface{}, _ cty.Path) diag.Diagnostics {
								switch v.(string) {
								case upcloud.IPAddressFamilyIPv4, upcloud.IPAddressFamilyIPv6:
									return nil
								default:
									return diag.Diagnostics{diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "'family' has incorrect value",
										Detail: fmt.Sprintf("'family' should have value of %s or %s",
											upcloud.IPAddressFamilyIPv4,
											upcloud.IPAddressFamilyIPv6),
									}}
								}
							},
						},
						"gateway": {
							Type:        schema.TypeString,
							Description: "Gateway address given by DHCP",
							Computed:    true,
							Optional:    true,
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Description: "A valid name for the network",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The network type",
				Computed:    true,
			},
			"zone": {
				Type:        schema.TypeString,
				Description: "The zone the network is in, e.g. `de-fra1`. You can list available zones with `upctl zone list`.",
				Required:    true,
				ForceNew:    true,
			},
			"router": {
				Type:        schema.TypeString,
				Description: "The UUID of a router",
				Optional:    true,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*service.Service)

	req := request.CreateNetworkRequest{}
	if v := d.Get("name"); v != nil {
		req.Name = v.(string)
	}

	if v := d.Get("zone"); v != nil {
		req.Zone = v.(string)
	}

	if v := d.Get("router"); v != nil {
		req.Router = v.(string)
	}

	if v, ok := d.GetOk("ip_network"); ok {
		ipn := v.([]interface{})[0]
		ipnConf := ipn.(map[string]interface{})

		uipn := upcloud.IPNetwork{
			Address:          ipnConf["address"].(string),
			DHCP:             upcloud.FromBool(ipnConf["dhcp"].(bool)),
			DHCPDefaultRoute: upcloud.FromBool(ipnConf["dhcp_default_route"].(bool)),
			Family:           ipnConf["family"].(string),
			Gateway:          ipnConf["gateway"].(string),
		}

		for _, dns := range ipnConf["dhcp_dns"].(*schema.Set).List() {
			uipn.DHCPDns = append(uipn.DHCPDns, dns.(string))
		}

		for _, route := range ipnConf["dhcp_routes"].(*schema.Set).List() {
			uipn.DHCPRoutes = append(uipn.DHCPRoutes, route.(string))
		}

		req.IPNetworks = append(req.IPNetworks, uipn)
	}

	network, err := client.CreateNetwork(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(network.UUID)

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*service.Service)

	req := request.GetNetworkDetailsRequest{
		UUID: d.Id(),
	}

	network, err := client.GetNetworkDetails(ctx, &req)
	if err != nil {
		return utils.HandleResourceError(d.Get("name").(string), d, err)
	}

	_ = d.Set("name", network.Name)
	_ = d.Set("type", network.Type)
	_ = d.Set("zone", network.Zone)

	if network.Router != "" {
		_ = d.Set("router", network.Router)
	}

	if len(network.IPNetworks) > 1 {
		return diag.Errorf("too many ip_networks: %d", len(network.IPNetworks))
	}

	if len(network.IPNetworks) == 1 {
		ipn := map[string]interface{}{
			"address":            network.IPNetworks[0].Address,
			"dhcp":               network.IPNetworks[0].DHCP.Bool(),
			"dhcp_default_route": network.IPNetworks[0].DHCPDefaultRoute.Bool(),
			"dhcp_dns":           network.IPNetworks[0].DHCPDns,
			"dhcp_routes":        network.IPNetworks[0].DHCPRoutes,
			"family":             network.IPNetworks[0].Family,
			"gateway":            network.IPNetworks[0].Gateway,
		}

		_ = d.Set("ip_network", []map[string]interface{}{
			ipn,
		})
	}

	d.SetId(network.UUID)

	return nil
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*service.Service)

	req := request.ModifyNetworkRequest{
		UUID: d.Id(),
	}

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		req.Name = v.(string)
	}

	if d.HasChange("ip_network") {
		v := d.Get("ip_network")

		ipn := v.([]interface{})[0]
		ipnConf := ipn.(map[string]interface{})

		uipn := upcloud.IPNetwork{
			Address:          ipnConf["address"].(string),
			DHCP:             upcloud.FromBool(ipnConf["dhcp"].(bool)),
			DHCPDefaultRoute: upcloud.FromBool(ipnConf["dhcp_default_route"].(bool)),
			Family:           ipnConf["family"].(string),
			Gateway:          ipnConf["gateway"].(string),
		}

		for _, dns := range ipnConf["dhcp_dns"].(*schema.Set).List() {
			uipn.DHCPDns = append(uipn.DHCPDns, dns.(string))
		}

		for _, route := range ipnConf["dhcp_routes"].(*schema.Set).List() {
			uipn.DHCPRoutes = append(uipn.DHCPRoutes, route.(string))
		}

		req.IPNetworks = []upcloud.IPNetwork{uipn}
	}

	network, err := client.ModifyNetwork(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("router") {
		_, v := d.GetChange("router")
		if v.(string) == "" {
			err = client.DetachNetworkRouter(ctx, &request.DetachNetworkRouterRequest{NetworkUUID: d.Id()})
		} else {
			err = client.AttachNetworkRouter(ctx, &request.AttachNetworkRouterRequest{NetworkUUID: d.Id(), RouterUUID: v.(string)})
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(network.UUID)

	return nil
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*service.Service)

	req := request.DeleteNetworkRequest{
		UUID: d.Id(),
	}
	err := client.DeleteNetwork(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
