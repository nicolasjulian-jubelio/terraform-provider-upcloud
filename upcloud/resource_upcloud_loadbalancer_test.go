package upcloud

import (
	"regexp"
	"testing"

	"github.com/UpCloudLtd/terraform-provider-upcloud/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccUpcloudLoadBalancer(t *testing.T) {
	testDataS1 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_s1.tf")
	testDataS2 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_s2.tf")
	testDataS3 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_s3.tf")
	testDataS4 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_s4.tf")

	var providers []*schema.Provider
	lbName := "upcloud_loadbalancer.lb"
	dnsName := "upcloud_loadbalancer_resolver.lb_dns_1"
	be1Name := "upcloud_loadbalancer_backend.lb_be_1"
	be2Name := "upcloud_loadbalancer_backend.lb_be_2"
	be1sm1Name := "upcloud_loadbalancer_static_backend_member.lb_be_1_sm_1"
	be1dm1Name := "upcloud_loadbalancer_dynamic_backend_member.lb_be_1_dm_1"
	be2sm1Name := "upcloud_loadbalancer_static_backend_member.lb_be_2_sm_1"
	fe1Name := "upcloud_loadbalancer_frontend.lb_fe_1"
	fe1Rule1Name := "upcloud_loadbalancer_frontend_rule.lb_fe_1_r1"
	fe1TLS1Name := "upcloud_loadbalancer_frontend_tls_config.lb_fe_1_tls1"
	cbd1Name := "upcloud_loadbalancer_dynamic_certificate_bundle.lb_cb_d1"
	cbm1Name := "upcloud_loadbalancer_manual_certificate_bundle.lb_cb_m1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		Steps: []resource.TestStep{
			{
				Config: testDataS1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lbName, "plan", "development"),
					resource.TestCheckResourceAttr(lbName, "zone", "fi-hel2"),
					resource.TestCheckResourceAttrSet(lbName, "dns_name"),
					resource.TestCheckResourceAttr(lbName, "labels.%", "2"),
					resource.TestCheckResourceAttr(lbName, "labels.key", "value"),
					resource.TestCheckResourceAttr(lbName, "labels.test-step", "1"),
					resource.TestCheckResourceAttr(dnsName, "name", "lb-resolver-1-test"),
					resource.TestCheckResourceAttr(dnsName, "nameservers.#", "2"),
					resource.TestCheckResourceAttr(be1Name, "name", "lb-be-1-test"),
					resource.TestCheckResourceAttr(be1Name, "resolver_name", "lb-resolver-1-test"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.timeout_server", "10"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.timeout_tunnel", "3600"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_type", "tcp"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_interval", "10"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_fall", "3"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_rise", "3"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_url", "https://10.0.0.10/healthz"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_tls_verify", "true"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_expected_status", "200"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.sticky_session_cookie_name", ""),
					resource.TestCheckResourceAttr(be1Name, "properties.0.outbound_proxy_protocol", ""),
					resource.TestCheckResourceAttr(be2Name, "name", "lb-be-2-test"),
					resource.TestCheckResourceAttr(be2Name, "resolver_name", ""),
					resource.TestCheckResourceAttr(be1sm1Name, "name", "lb-be-1-sm-1-test"),
					resource.TestCheckResourceAttr(be1sm1Name, "ip", "10.0.0.10"),
					resource.TestCheckResourceAttr(be1sm1Name, "port", "8000"),
					resource.TestCheckResourceAttr(be1sm1Name, "weight", "100"),
					resource.TestCheckResourceAttr(be1sm1Name, "max_sessions", "1000"),
					resource.TestCheckResourceAttr(be1sm1Name, "enabled", "true"),
					resource.TestCheckResourceAttr(be1dm1Name, "name", "lb-be-1-dm-1-test"),
					resource.TestCheckResourceAttr(be1dm1Name, "weight", "10"),
					resource.TestCheckResourceAttr(be1dm1Name, "max_sessions", "10"),
					resource.TestCheckResourceAttr(be1dm1Name, "enabled", "false"),
					resource.TestCheckResourceAttr(be2sm1Name, "name", "lb-be-2-sm-1-test"),
					resource.TestCheckResourceAttr(be2sm1Name, "enabled", "true"),
					resource.TestCheckResourceAttr(be2sm1Name, "weight", "1"),
					resource.TestCheckResourceAttr(fe1Name, "name", "lb-fe-1-test"),
					resource.TestCheckResourceAttr(fe1Name, "port", "8080"),
					resource.TestCheckResourceAttr(fe1Name, "mode", "http"),
					resource.TestCheckResourceAttr(fe1Name, "default_backend_name", "lb-be-1-test"),
					resource.TestCheckResourceAttr(fe1Name, "properties.0.timeout_client", "10"),
					resource.TestCheckResourceAttr(fe1Name, "properties.0.inbound_proxy_protocol", "false"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "name", "lb-fe-1-r1-test"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "priority", "10"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_port.0.method", "equal"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_port.0.value", "80"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_port_range.0.range_start", "100"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_port_range.0.range_end", "1000"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_ip.0.value", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.body_size.0.method", "equal"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.body_size.0.value", "8000"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.body_size_range.0.range_start", "1000"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.body_size_range.0.range_end", "1001"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.path.0.value", "/application"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.path.0.method", "starts"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.path.0.ignore_case", "true"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url.0.value", "/application"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url.0.method", "starts"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url.0.ignore_case", "true"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url_query.0.method", "starts"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url_query.0.value", "type=app"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.host.0.value", "example.com"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.http_method.0.value", "PATCH"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.cookie.0.value", "123456"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.cookie.0.name", "x-session-id"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.cookie.0.method", "exact"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.header.0.method", "exact"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.header.0.name", "status"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.header.0.value", "active"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.header.0.ignore_case", "true"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url_param.0.method", "exact"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url_param.0.name", "status"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url_param.0.value", "active"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.url_param.0.ignore_case", "true"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.num_members_up.0.backend_name", "lb-be-1-test"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.num_members_up.0.value", "1"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.num_members_up.0.method", "less"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.num_members_up.0.method", "less"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.use_backend.0.backend_name", "lb-be-1-test"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.http_redirect.0.location", "/app"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.http_redirect.1.scheme", "https"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.http_return.0.status", "404"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.http_return.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.http_return.0.payload", "UmVzb3VyY2Ugbm90IGZvdW5kIQ=="),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.tcp_reject.0.active", "true"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.set_forwarded_headers.0.active", "true"),
					resource.TestCheckResourceAttr(cbd1Name, "name", "lb-cb-d1-test"),
					resource.TestCheckResourceAttr(cbd1Name, "key_type", "rsa"),
					resource.TestCheckResourceAttr(cbd1Name, "hostnames.0", "example.com"),
					resource.TestCheckResourceAttr(cbm1Name, "name", "lb-cb-m1-test"),
					resource.TestCheckResourceAttr(cbm1Name, "certificate", "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZhekNDQTFPZ0F3SUJBZ0lVRzR1KzRDZmlHQ3pQSDk4dDA4QXh5VkE0QzVnd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TWpBek1UY3hNekUxTURoYUZ3MHlNekF6Ck1UY3hNekUxTURoYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdnZ0lpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQ0R3QXdnZ0lLQW9JQ0FRRHZaeG4vK3pUeVc0RVJ2S2t3V29tVXBpOG8ydEp6MWR2ZXIrREpySzNnCkNObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZDFxNjlGSStXY1BFMWFhYjk5MDlJQWsvR0dpSzIKelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpURThDUjlaaTJPTEV5eEdRMzZwQTYxOTBueFZnMgpTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm9qSC9uWGs1Vy9Yd3paYm9JSHNTRlZZaksyemxnCm9xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMTdFVG9kSHNLSlpCb1h4RHNaVVRHQ0RESkNpbXoKVzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyTW5YZWJrT2pnY3ZVWTROSXd6cFlWMFcyR0dHRwp3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzaldsdjBUZUd4RStCWlorZzN4UHFIVEd6MndIekM1CjVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTlh6aW5MMGtlRnh1cExYWThiS2hFelh6U0xLeisKQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMUzd6YlhGb1FycDF4L2RXbTlTOERWRVhWb1VBMQpXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGJKOG0wVGZ2Zmt1aFZpdVZBTi9wK01wOVduUStICjExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdjFycE5QelVVUHd0QXF4OTIzOXd3SVI5WTE1Y2wKT1FJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dId1lEVlIwagpCQmd3Rm9BVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dEd1lEVlIwVEFRSC9CQVV3QXdFQi96QU5CZ2txCmhraUc5dzBCQVFzRkFBT0NBZ0VBQ00reGJiOW9rVUdPcWtRWmtndHh3eVF0em90WFRXNnlIYmdmeGhCd3d1disKc0ltT1QyOUw5WFVZblM5cTgrQktBSHVGa3JTRUlwR0ZVVWhDVWZTa2xlNmZuR0c1b05PWm1DdEszN3RoZlFVOQp2NEJvOHFCWkhqREh3azlWVHRhWm1BazBLYnhmaHVneVdWQ1ZsbURacm9TQ09pV0drVFZoc1hhS0RrYnc0RWwwCjJzY3lnYkFDdFZ4bkU1WjlmU0F3MU9QWXJZYUcySW5HTDQvMHVSZXo4aXl1UE9lNUNiL0RkNDl1eHFzR1FkM1IKQzdKNC9vWnB2b0V6UVJtakxib1FzQzkwU2ZqaFNpcGhHQlNiYUpCZGRsMDBrNVZzVXJxS1haU004cVFxVWZWLwpubEJtYjJOblVsa2RlOEtIczBQamhCaG8rLzdmaitMN21GYTJsNWpmdWlsdHVxdmgyWnladFJjd2didmJlaUxPCmZQSWlMQ2dTbnMwaitZMkVrS1drRUp6RXJQVm5sOTdaQktZclBaYmRYMFY5b2dvTC9qeEV5NzlsbzlKczI5djYKUkY2NmdvSlUwMkVKZTUwMmk3WHJzMzFZQ0tuSGd2ejUwTDZha0JpYWRSNmtrTXVXdkJ1d1l6MElaS1RMcXhqZAowOEdlUkJVeWFsUFZodGZKbzNNdXRuYUllL1pWVTdLQUl3S1Znb20zS09EY1RpWllQV3RWKzFnL0UvN3A1aGh2CkJERzFqcklRc1ZrZG4yNWZhNXNkNU9Qa1AvbDBRdXY1em16UEk3S1MrS2ZlWS92NHFBOTBtNGk2dkZORlRtbTAKSFNXV0JZTlR4blIxYjk2UElUcnRzOE15am9YTFg2QnUxVkZOSlByMkpnMDJMVlZvcTZSSWJlMVVvNjE5b2pBPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="),
					resource.TestCheckResourceAttr(cbm1Name, "private_key", "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUpSUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQ1M4d2dna3JBZ0VBQW9JQ0FRRHZaeG4vK3pUeVc0RVIKdktrd1dvbVVwaThvMnRKejFkdmVyK0RKckszZ0NObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZAoxcTY5RkkrV2NQRTFhYWI5OTA5SUFrL0dHaUsyelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpUCkU4Q1I5WmkyT0xFeXhHUTM2cEE2MTkwbnhWZzJTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm8KakgvblhrNVcvWHd6WmJvSUhzU0ZWWWpLMnpsZ29xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMQo3RVRvZEhzS0paQm9YeERzWlVUR0NEREpDaW16VzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyCk1uWGVia09qZ2N2VVk0Tkl3enBZVjBXMkdHR0d3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzalcKbHYwVGVHeEUrQlpaK2czeFBxSFRHejJ3SHpDNTVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTgpYemluTDBrZUZ4dXBMWFk4YktoRXpYelNMS3orQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMClM3emJYRm9RcnAxeC9kV205UzhEVkVYVm9VQTFXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGIKSjhtMFRmdmZrdWhWaXVWQU4vcCtNcDlXblErSDExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdgoxcnBOUHpVVVB3dEFxeDkyMzl3d0lSOVkxNWNsT1FJREFRQUJBb0lDQVFDcUNtd2dNbmcvOEJoejFiSENRM3hYCkZkYUhTUzJUMHdHUllSRGpqZ0FPRVpyMERxN3IzQnFDLy9Jd1RMZlRaZ2dKQmpPaWpPd0I4TE01cGVPRkwxWngKZjVVRDRDQVpZUkJ4MEJxRFZjcjBWajM2R3B6MjlLUnZFV3JDTWptaitlZUtHZ3NVVEp3TmpnRGk1N091dUdlWQpmaG4yT2lJSXlWVmFSanF4NWV5cTJlcTFSOVMvd3BlVElSek9zdTlyU29la1V5SDFZZDBTMS9TdXpLU0lYS1orCkNSdXZrZ0NaaGVrRjMyUUMyY1VlUzBTb3FFY1VtUEJXY0dzRk4xTFV1K3ZQNzBBZ0ZZV0lQbHBXZHRQVzIrME0KbnZPNy9sSVI1amY4QkpOS0tDcklWMFVKb3ZTV3h1VGlxYjNpVUFnTUwxQTNnQXJwZUVOaEFRMjZYWXIweXhMRQpiMTRObWdnZzRUSktRMWp4aTZlRys5Nkp0WStteFdwaHJUNTVUT2s5blpEY2JTV29ibTgwcExJZEE0QTN0bjhJCjI5SXZ2dkdhVnh2NXFXaU5sL0sxWEZhK2JRRVY3K1AxT0RWSnV5VXEzdFg4dDlONVJ2c3RiUE5XNU8xQ01STGsKZExESVMwRFFKYytMa1ZGaGpZRXlJOWluQmtlRXV1ekI0L0k3M1JnTXpKMFZvZU1hSXEwdWpxbVF6KzN0ei9JUgp6VTFSN2FndEZmNHUvTS9jeVl2U0E4UEZVSis1Q25ldHFtWC91dzR0em10WmtLd3JnMVArblVpRTV4Z0ZLT1FZCjUyaU81aXRKWHo4aHJGdFpmbVMrcTI3eFRFWWlEdTFFSFg0Y1pLaFBQNVh3RzkyekN0ZUxJenc5QUh2TS9aRTcKNmI0OFMwQWR6T3dFN3VaVVhheEw5UUtDQVFFQTkrOVROV2RuMjVxK3lSWUlsc2x1NFI1N0xSQm1ucXFhaDljSQowZ1JGS0RZUFZTWFl0RlhVci9mK0psSmR4YUxEYjYwV0o4c2hWcHZ1MTlJTE0wNTR6M3ZYa0w0QjQwSWI1Z3JnCnlzbVZKVlpoYVdTNWU5RW02TEh2bXRoVWl0MzVDWnlLcTlkN2FwdzFOQ1VIUkFEWE9wSndBWkc4ckttRmtDeFYKTnFpVThnWi9LOXpPK0gra215bGVKZTkrZHhCR2RCZTczbnFsRTVPQVdJcjVLcFNpMU5sQWYxOG0yYnVUWkxiTQo2Rnd1MEc0SjUzd1ZaTFRhbWxVelNmdjFTRGRrMEZIbnpkdEJPMWI3OTNJeVpOQXc4eGlkMGVPVy90OE1QTm1RCmFXWnNqaWc4WEZVMlFaMXl5YUd6amoxNUpoVzl2WGN1SUpGc21rRzZiRS8wSm9ZZnd3S0NBUUVBOXpDNXNYN1IKQlN2MzlwVnFROXI4bDJXNVVqQk9iczFick05Q0NPQWFFZCtlY2xyZWZZRU5yMUVtL0x3Qk9GK1RuRnVnMjJZOQpKczliNjgrY0wzSEtweVJ5KzVVaGM1NzU0b3pDRklDemJzZHphM005TGVzVUtTVTBYYnlhNm5vc0UvWjNOWGpyCmpLQ1ROZEU2eFY4OTZ1Vm5FTXZpcDh4M1kvR0VBejh4U3lCOUFveGNtanVqZnVVdmlZbmN5SnBqUUoxOEhsZk8KMlloWWdTd3VxSFhISUswUXhGckJiTjIzZXZ5TUkrOVVSSVlnOWxtemFWR2hqczlHd0Y3UUQ0dFoyOXgwQVFOSwpBUFArMnI4ZUlLbjFucEs4VnFRdVNrNVdsZjg3L0xvV1Zha2xtRXkvdEJ0alpjNVBQanpMVmRTLy9QNWlkc2F1CkZuTnc2VDNmZWEwelV3S0NBUUVBMWsySDU2WXNzRFhPYUxOaDB5dmphalJWbGJzU2FGemdXei8wQU12dUZ2YTcKUkFjRmk4S1FwMVU4MlZUaWRyemNIc0JHWVRrRDVQKzliOUMvRzZiZFo4SU1yckI5b3poMk10MytOV29PUDRxdApnbEtzdktnbzhJTTBydXdFRDFBVVBVbVExejNYRUd4YTFHcVpJQjkxNmN1L2dxdThvS1dhcSthVjlUdThHb0toCkU0RzFhRGUwU09WMTJtWnJNbkRmNU9MSzRWK3pKZnVkdVdyT09nN2x2QUxZNi8rTDdqRmpFbStySjhEZU9neVQKQlFKTTM1SXZUYTBOT3dyTWxaSkQwb2lwUzFjVHlEM0VacnJQY2pJOXpUSGU0QmZQWVJmY1ZSQmM4YTIxY1I2NApKYnNGdmF0aEY0VnNWU3N2ZDByZGlWSGxqZ01GRTBSeTVjSXFMODVJendLQ0FRRUEzdE1nZ1QwRkpIbFhKQVBhCmIrS0drZDlUNkIrOWhDcEFPbzMyUTlQb0REYWRPUTVxdzQzRERVZkZNa3d6ZVdMR3lFcmN2UW56ay9tV0xnTFAKRXdHcm9YRzg2TWF0Q2ZIRDVoSG1uZDdLWU5FUVhVcmJXbm92aVV1TllmWXpXNnpYOFFMYXdPd0l3Wks2UU9nago1Mm1NZ2lOYS9nd2NmQkJYaTFOYUlpY2p3MG85QmtBSzljbE8vNE9QajVjajIvMDMvVFk1Zll5LzNONElraUNHCnlycW96dTdUVDMxVUlWUFlJdGhuWjdsRktDUVVzSjE1bWpYSXdkaGRPZW45K2hVdTRuOWVYczlkTlhDOVN1aS8KT3NpYXJlQXVRSmZ0Vm5RNW55c2VJeHFJS1oyNVV3blVRWUh5M3dIVDh4R1FaZ1hMTHo4TStXN3QzVFVoRWxBQgpGRWtxR3dLQ0FRRUFtT2VzMXFINlA0dDRsZ0VjK01Ubzc2VmJvYnhab29HQWdMWmc3N1J4c2hPaGxMdCtYZnIzCjFsOWgycFJ0eXFIdDRuMG9ob1A0VzNSN3VuQ05rNWl3SGNJSDVjNmhGQTYyelVEM1JjeEhJZERhdDBGL3RoWDgKTUpndDlsM0MrMVJwZFZlV0hlbEJOM0JlM0FtWEpXL1ZsL3lGTXVjWWxETnlXUFlPSmRuQ1BRZ0FGVFJnVnJlUQpiUjZCY29neUVRTVEzenRMWnNBRnRaZ25Sak05YkpLN1JjYzg5bGxaa1BuMHVKZjNKVWxMeVFFN0l2bEJsWi9tClZnUUhiRTkwQStZNzFpb1piQWh5TFcwTE9lTmhBS3NRNFJZbDlnb0N4dGp1ZnE2NnFTNGNzdGN6c2J5N083dFAKeXZkSXp2eEZRZmx4Yk8ra1ludDNkcFRIdUNuUkNIMFM0UT09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K"),
					resource.TestCheckResourceAttr(fe1TLS1Name, "name", "lb-fe-1-tls1-test"),
				),
			},
			{
				Config:            testDataS2,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lbName, "plan", "development"),
					resource.TestCheckResourceAttr(lbName, "labels.%", "2"),
					resource.TestCheckResourceAttr(lbName, "labels.key", "value"),
					resource.TestCheckResourceAttr(lbName, "labels.test-step", "2"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_port.0.method", "equal"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "matchers.0.src_port.0.value", "80"),
					resource.TestCheckResourceAttr(fe1Rule1Name, "actions.0.use_backend.0.backend_name", "lb-be-1-test-1"),
					resource.TestCheckResourceAttr(be1Name, "resolver_name", "lb-resolver-1-test-1"),
					resource.TestCheckResourceAttr(fe1Name, "default_backend_name", "lb-be-1-test-1"),
					resource.TestCheckResourceAttr(fe1Name, "properties.0.timeout_client", "20"),
					resource.TestCheckResourceAttr(fe1Name, "properties.0.inbound_proxy_protocol", "true"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.timeout_server", "20"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.timeout_tunnel", "4000"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_type", "http"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_interval", "10"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_fall", "3"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_rise", "3"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_url", "/"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.health_check_expected_status", "200"),
					resource.TestCheckResourceAttr(be1Name, "properties.0.sticky_session_cookie_name", ""),
					resource.TestCheckResourceAttr(be1Name, "properties.0.outbound_proxy_protocol", "v2"),
					resource.TestCheckResourceAttr(be2sm1Name, "weight", "0"),
				),
			},
			{
				Config:            testDataS3,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lbName, "network", ""),
					resource.TestCheckResourceAttr(lbName, "networks.#", "2"),
					resource.TestCheckResourceAttr(lbName, "networks.0.name", "lan-1"),
					resource.TestCheckResourceAttr(lbName, "networks.0.type", "private"),
					resource.TestCheckResourceAttr(lbName, "networks.1.name", "lan-2"),
					resource.TestCheckResourceAttr(lbName, "networks.1.type", "private"),
					resource.TestCheckResourceAttr(lbName, "labels.%", "0"),
				),
			},
			{
				Config:            testDataS4,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lbName, "networks.0.name", "lan-a"),
					resource.TestCheckResourceAttr(lbName, "networks.0.type", "private"),
					resource.TestCheckResourceAttr(lbName, "networks.1.name", "lan-b"),
					resource.TestCheckResourceAttr(lbName, "networks.1.type", "private"),
				),
			},
		},
	})
}

func TestAccUpcloudLoadBalancer_HTTPRedirectValidation(t *testing.T) {
	// These test data files should fail in pre-plan validation. Thus, these tests are run in plan-only mode.
	testDataE1 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_e1.tf")
	testDataE2 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_e2.tf")
	testDataE3 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_e3.tf")
	testDataE4 := utils.ReadTestDataFile(t, "testdata/upcloud_loadbalancer/loadbalancer_e4.tf")

	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(&providers),
		Steps: []resource.TestStep{
			{
				Config:      testDataE1,
				ExpectError: regexp.MustCompile("actions block should contain at least one action"),
				PlanOnly:    true,
			},
			{
				Config:      testDataE2,
				ExpectError: regexp.MustCompile("either location or scheme should be defined for http_redirect"),
				PlanOnly:    true,
			},
			{
				Config:      testDataE3,
				ExpectError: regexp.MustCompile(`Error: expected .*scheme to be one of \["http" "https"\], got invalid`),
				PlanOnly:    true,
			},
			{
				Config:      testDataE4,
				ExpectError: regexp.MustCompile(`only either location or scheme should be defined at a time for http_redirect`),
				PlanOnly:    true,
			},
		},
	})
}
