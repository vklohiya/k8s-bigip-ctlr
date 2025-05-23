package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sync"
)

var _ = Describe("AS3Handler Tests", func() {

	Describe("Prepare AS3 Declaration", func() {
		var mem1, mem2, mem3, mem4 PoolMember
		var as3Handler AS3Handler
		BeforeEach(func() {
			as3Handler = AS3Handler{
				AS3Parser:   &AS3Parser{},
				PostManager: &PostManager{declUpdate: sync.Mutex{}},
			}
			mem1 = PoolMember{
				Address:         "1.2.3.5",
				Port:            8080,
				ConnectionLimit: 5,
				AdminState:      "disable",
			}
			mem2 = PoolMember{
				Address:         "1.2.3.6",
				Port:            8081,
				ConnectionLimit: 5,
			}
			mem3 = PoolMember{
				Address: "1.2.3.7",
				Port:    8082,
			}
			mem4 = PoolMember{
				Address: "1.2.3.8",
				Port:    8083,
			}
		})
		It("VirtualServer Declaration", func() {

			rsCfg := &ResourceConfig{}
			rsCfg.MetaData.Active = true
			rsCfg.MetaData.ResourceType = VirtualServer
			rsCfg.Virtual.Name = "crd_vs_172.13.14.15"
			rsCfg.Virtual.PoolName = "default_pool_svc1"
			rsCfg.Virtual.Destination = "/test/172.13.14.5:8080"
			rsCfg.Virtual.AllowVLANs = []string{"flannel_vxlan"}
			rsCfg.Virtual.IpIntelligencePolicy = "/Common/ip-intelligence-policy"
			rsCfg.Virtual.HTTPCompressionProfile = "/Common/compressionProfile"
			rsCfg.Virtual.BigIPRouteDomain = 10
			rsCfg.Virtual.AdditionalVirtualAddresses = []string{"172.13.14.17", "172.13.14.18"}
			rsCfg.Virtual.ProfileAdapt = ProfileAdapt{"/Common/example-requestadapt", "/Common/example-responseadapt"}
			rsCfg.Virtual.Policies = []nameRef{
				{
					Name:      "policy1",
					Partition: "test",
				},
				{
					Name:      "policy2",
					Partition: "test",
				},
			}
			rsCfg.Pools = Pools{
				Pool{
					Name:            "pool1",
					MinimumMonitors: intstr.IntOrString{Type: 1, StrVal: "all"},
					Members:         []PoolMember{mem1, mem2},
					MonitorNames: []MonitorName{
						{Name: "/test/http_monitor"},
					},
				},
			}
			rsCfg.Virtual.IRules = []string{"none", "/common/irule1", "/common/irule_2", "/common/http_redirect_irule"}
			rsCfg.IRulesMap = IRulesMap{
				NameRef{"custom_iRule", DEFAULT_PARTITION}: &IRule{
					Name:      "custom_iRule",
					Partition: DEFAULT_PARTITION,
					Code:      "tcl code blocks",
				},
				NameRef{HttpRedirectIRuleName, DEFAULT_PARTITION}: &IRule{
					Name:      HttpRedirectIRuleName,
					Partition: DEFAULT_PARTITION,
					Code:      "tcl code blocks",
				},
			}
			rsCfg.IntDgMap = InternalDataGroupMap{
				NameRef{"static-internal-dg", "test"}: DataGroupNamespaceMap{
					"intDg1": &InternalDataGroup{
						Name:      "static-internal-dg",
						Partition: "test",
						Type:      "string",
						Records: []InternalDataGroupRecord{
							{
								Name: "apiTye",
								Data: AS3,
							},
						},
					},
				},
			}
			rsCfg.Policies = Policies{
				Policy{
					Name:     "policy1",
					Strategy: "first-match",
					Rules: Rules{
						&Rule{
							Conditions: []*condition{
								{
									Values: []string{"test.com"},
									Equals: true,
								},
							},
							Actions: []*action{
								{
									Forward:  true,
									Request:  true,
									Redirect: true,
									HTTPURI:  true,
									HTTPHost: true,
									Pool:     "default_svc_1",
								},
							},
						},
					},
				},
				Policy{
					Name:     "policy2",
					Strategy: "first-match",
					Rules: Rules{
						&Rule{
							Conditions: []*condition{
								{
									Host:     true,
									Values:   []string{"prod.com"},
									Equals:   true,
									HTTPHost: true,
									Request:  true,
								},
								{
									PathSegment: true,
									Index:       1,
									HTTPURI:     true,
									Equals:      true,
									Values:      []string{"/foo"},
									Request:     true,
								},
							},
							Actions: []*action{
								{
									Forward:       true,
									Request:       true,
									Redirect:      true,
									HTTPURI:       true,
									HTTPHost:      true,
									Pool:          "default_svc_2",
									Log:           true,
									Location:      PrimaryBigIP,
									Message:       "log action",
									Replace:       true,
									Value:         "urihost",
									WAF:           true,
									Policy:        "/common/policy3",
									Enabled:       true,
									Drop:          true,
									PersistMethod: SourceAddress,
								},
								{
									PersistMethod: DestinationAddress,
								},
								{
									PersistMethod: CookieHash,
								},
								{
									PersistMethod: CookieInsert,
								},
								{
									PersistMethod: CookieRewrite,
								},
								{
									PersistMethod: CookiePassive,
								},
								{
									PersistMethod: Universal,
								},
								{
									PersistMethod: Carp,
								},
								{
									PersistMethod: Hash,
								},
								{
									PersistMethod: Disable,
								},
								{
									PersistMethod: "Disable",
								},
							},
						},
					},
				},
				Policy{
					Name:     "policy3",
					Strategy: "first-match",
					Rules: Rules{
						&Rule{
							Conditions: []*condition{
								{
									Path:    true,
									Name:    "condition3",
									Values:  []string{"/common/test"},
									HTTPURI: true,
									Equals:  true,
									Index:   3,
								},
								{
									Tcp:     true,
									Address: true,
									Values:  []string{"10.10.10.10"},
									Request: true,
								},
							},
							Actions: []*action{
								{
									Forward:  true,
									Request:  true,
									Redirect: true,
									HTTPURI:  true,
									HTTPHost: true,
									Pool:     "default_svc_2",
								},
							},
						},
					},
				},
			}
			rsCfg.Monitors = Monitors{
				{
					Name:       "http_monitor",
					Interval:   10,
					Type:       "http",
					TargetPort: 8080,
					Timeout:    10,
					Send:       "GET /health",
				},
				{
					Name:       "https_monitor",
					Interval:   10,
					Type:       "https",
					TargetPort: 8443,
					Timeout:    10,
					Send:       "GET /health",
				},
				{
					Name:       "tcp_monitor",
					Interval:   10,
					Type:       "tcp",
					TargetPort: 3600,
					Timeout:    10,
					Send:       "GET /health",
				},
			}

			rsCfg.Virtual.Profiles = ProfileRefs{
				ProfileRef{
					Name:      "serverssl",
					Partition: "Common",
					Context:   "serverside",
				},
				ProfileRef{
					Name:    "serversslnew",
					Context: "serverside",
				},
				ProfileRef{
					Name:      "clientssl",
					Partition: "Common",
					Context:   "clientside",
				},
				ProfileRef{
					Name:    "clientsslnew",
					Context: "clientside",
				},
			}

			rsCfg2 := &ResourceConfig{}
			rsCfg2.MetaData.Active = false
			rsCfg2.MetaData.defaultPoolType = BIGIP
			rsCfg2.MetaData.ResourceType = VirtualServer
			rsCfg2.Virtual.Name = "crd_vs_172.13.14.16"
			rsCfg.Virtual.PoolName = "default_pool_svc2"
			rsCfg2.Pools = Pools{
				Pool{
					Name:            "pool1",
					Members:         []PoolMember{mem3, mem4},
					MinimumMonitors: intstr.IntOrString{Type: 0, IntVal: 1},
				},
			}

			enabled := false
			rsCfg2.customProfiles = make(map[SecretKey]CustomProfile)
			cert := certificate{Cert: "crthash", Key: "keyhash"}
			rsCfg2.customProfiles[SecretKey{
				Name:         "default_svc_test_com_cssl",
				ResourceName: "crd_vs_172.13.14.15",
			}] = CustomProfile{
				Name:          "default_svc_test_com_cssl",
				Partition:     "test",
				Context:       "clientside",
				Certificates:  []certificate{cert},
				SNIDefault:    false,
				TLS1_0Enabled: &enabled,
				TLS1_1Enabled: &enabled,
				TLS1_2Enabled: &enabled,
			}
			certOnly := certificate{Cert: "crthash"}
			rsCfg2.customProfiles[SecretKey{
				Name:         "default_svc_test_com_sssl",
				ResourceName: "crd_vs_172.13.14.15",
			}] = CustomProfile{
				Name:          "default_svc_test_com_sssl",
				Partition:     "test",
				Context:       "serverside",
				Certificates:  []certificate{certOnly},
				ServerName:    "test.com",
				SNIDefault:    false,
				TLS1_0Enabled: &enabled,
				TLS1_1Enabled: &enabled,
				TLS1_2Enabled: &enabled,
			}

			rsCfg3 := &ResourceConfig{}
			rsCfg3.MetaData.Active = false
			rsCfg3.MetaData.defaultPoolType = BIGIP
			rsCfg3.MetaData.ResourceType = VirtualServer
			rsCfg3.Virtual.Name = "crd_vs_172.13.14.17"
			rsCfg3.Virtual.PoolName = "default_pool_svc3"
			rsCfg3.Pools = Pools{
				Pool{
					Name:            "pool1",
					Members:         []PoolMember{mem3, mem4},
					MinimumMonitors: intstr.IntOrString{Type: 0, IntVal: 1},
					MonitorNames: []MonitorName{
						{Name: "/Common/http",
							Reference: "bigip"},
						{Name: "default_pool_svc3_http"},
						{Name: "default_pool_svc3_http_81"},
					},
				},
			}

			rsCfg3.customProfiles = make(map[SecretKey]CustomProfile)
			rsCfg3.customProfiles[SecretKey{
				Name:         "default_svc_test_com_cssl",
				ResourceName: "crd_vs_172.13.14.17",
			}] = CustomProfile{
				Name:         "default_svc_test_com_cssl",
				Partition:    "test",
				Context:      "clientside",
				Certificates: []certificate{cert},
				SNIDefault:   false,
			}
			rsCfg3.customProfiles[SecretKey{
				Name:         "default_svc_test_com_sssl",
				ResourceName: "crd_vs_172.13.14.17",
			}] = CustomProfile{
				Name:         "default_svc_test_com_sssl",
				Partition:    "test",
				Context:      "serverside",
				Certificates: []certificate{certOnly},
				ServerName:   "test.com",
				SNIDefault:   false,
			}

			rsCfg4 := &ResourceConfig{}
			rsCfg4.MetaData.Active = false
			rsCfg4.MetaData.defaultPoolType = BIGIP
			rsCfg4.MetaData.ResourceType = VirtualServer
			rsCfg4.Virtual.Name = "crd_vs_172.13.14.18"
			rsCfg4.Virtual.PoolName = "default_pool_svc3"
			rsCfg4.Pools = Pools{
				Pool{
					Name:            "pool1",
					Members:         []PoolMember{mem3, mem4},
					MinimumMonitors: intstr.IntOrString{Type: 0, IntVal: 1},
				},
			}

			rsCfg4.customProfiles = make(map[SecretKey]CustomProfile)
			rsCfg4.customProfiles[SecretKey{
				Name:         "default_svc_test_com_cssl",
				ResourceName: "crd_vs_172.13.14.18",
			}] = CustomProfile{
				Name:         "default_svc_test_com_cssl",
				Partition:    "test",
				Context:      "clientside",
				Certificates: []certificate{cert},
				SNIDefault:   false,
			}
			rsCfg4.customProfiles[SecretKey{
				Name:         "default_svc_test_com_sssl",
				ResourceName: "crd_vs_172.13.14.18",
			}] = CustomProfile{
				Name:         "default_svc_test_com_sssl",
				Partition:    "test",
				Context:      "serverside",
				Certificates: []certificate{certOnly},
				ServerName:   "test.com",
				SNIDefault:   false,
			}

			config := ResourceConfigRequest{
				ltmConfig:          make(LTMConfig),
				shareNodes:         true,
				gtmConfig:          GTMConfig{},
				defaultRouteDomain: 1,
			}
			zero := 0
			config.ltmConfig["default"] = &PartitionConfig{ResourceMap: make(ResourceMap), Priority: &zero}
			config.ltmConfig["default"].ResourceMap["crd_vs_172.13.14.15"] = rsCfg
			config.ltmConfig["default"].ResourceMap["crd_vs_172.13.14.16"] = rsCfg2
			config.ltmConfig["default"].ResourceMap["crd_vs_172.13.14.17"] = rsCfg3
			config.ltmConfig["default"].ResourceMap["crd_vs_172.13.14.18"] = rsCfg4

			agentPostCfg := as3Handler.createAPIConfig(config, true, "", false)
			decl := agentPostCfg.data
			Expect(decl).ToNot(Equal(""), "Failed to Create AS3 Declaration")
			Expect(strings.Contains(decl, "pool1")).To(BeTrue())
			Expect(strings.Contains(decl, "default_pool_svc2")).To(BeTrue())
			Expect(strings.Contains(decl, "default_pool_svc3")).To(BeTrue())
			Expect(strings.Contains(decl, "/Common/example-requestadapt")).To(BeTrue())
			Expect(strings.Contains(decl, "/Common/example-responseadapt")).To(BeTrue())
			Expect(strings.Contains(decl, "/Common/compressionProfile")).To(BeTrue())

			sharedApp := as3Application{}
			as3Handler.createPoolDecl(rsCfg3, sharedApp, false, "test", Cluster)
			as3Handler.createPoolDecl(rsCfg4, sharedApp, false, "test", Cluster)
			Expect(len(sharedApp)).To(Equal(1))
			Expect(len(sharedApp["pool1"].(*as3Pool).Monitors)).To(Equal(3))
			rsCfg4.Pools[0].MonitorNames = []MonitorName{
				{Name: "/Common/http",
					Reference: "bigip"},
				{Name: "default_pool_svc3_http"},
				{Name: "default_pool_svc3_http_81"},
			}
			sharedApp = as3Application{}
			as3Handler.createPoolDecl(rsCfg3, sharedApp, false, "test", Cluster)
			as3Handler.createPoolDecl(rsCfg4, sharedApp, false, "test", Cluster)
			Expect(len(sharedApp)).To(Equal(1))
			Expect(len(sharedApp["pool1"].(*as3Pool).Monitors)).To(Equal(3))
			rsCfg4.Pools[0].MonitorNames = []MonitorName{
				{Name: "/Common/http",
					Reference: "bigip1"},
				{Name: "default_pool_svc4_http"},
				{Name: "default_pool_svc3_http_81"},
			}
			sharedApp = as3Application{}
			as3Handler.createPoolDecl(rsCfg3, sharedApp, false, "test", Cluster)
			as3Handler.createPoolDecl(rsCfg4, sharedApp, false, "test", Cluster)
			Expect(len(sharedApp)).To(Equal(1))
			Expect(len(sharedApp["pool1"].(*as3Pool).Monitors)).To(Equal(5))

			rsCfg3.Pools[0].MonitorNames = []MonitorName{
				{Name: "/Common/http",
					Reference: "bigip"},
				{Name: "default_pool_svc3_http_81"},
			}
			sharedApp = as3Application{}
			as3Handler.createPoolDecl(rsCfg3, sharedApp, false, "test", Cluster)
			as3Handler.createPoolDecl(rsCfg4, sharedApp, false, "test", Cluster)
			Expect(len(sharedApp)).To(Equal(1))
			Expect(len(sharedApp["pool1"].(*as3Pool).Monitors)).To(Equal(4))

			rsCfg3.Pools[0].MonitorNames = []MonitorName{}
			sharedApp = as3Application{}
			as3Handler.createPoolDecl(rsCfg3, sharedApp, false, "test", Cluster)
			as3Handler.createPoolDecl(rsCfg4, sharedApp, false, "test", Cluster)
			Expect(len(sharedApp)).To(Equal(1))
			Expect(len(sharedApp["pool1"].(*as3Pool).Monitors)).To(Equal(3))

		})
		It("TransportServer Declaration", func() {
			rsCfg := &ResourceConfig{}
			rsCfg.MetaData.Active = true
			rsCfg.MetaData.ResourceType = TransportServer
			rsCfg.Virtual.Name = "crd_vs_172.13.14.16"
			rsCfg.Virtual.Mode = "standard"
			rsCfg.Virtual.IpProtocol = "tcp"
			rsCfg.Virtual.FTPProfile = "/Common/ftpProfile1"
			rsCfg.Virtual.TranslateServerAddress = true
			rsCfg.Virtual.TranslateServerPort = true
			rsCfg.Virtual.AllowVLANs = []string{"flannel_vxlan"}
			rsCfg.Virtual.Destination = "172.13.14.6:1600"
			rsCfg.customProfiles = make(map[SecretKey]CustomProfile)
			rsCfg.Virtual.ProfileL4 = "/Common/profileL4"
			rsCfg.Virtual.ProfileDOS = "/Common/profileDOS"
			rsCfg.Virtual.ProfileBotDefense = "/Common/profileBotDefense"
			rsCfg.Virtual.TCP.Client = "/Common/tcpClient"
			rsCfg.Virtual.TCP.Server = "/Common/tcpServer"
			rsCfg.Virtual.TranslateServerAddress = true
			rsCfg.Virtual.TranslateServerPort = true
			rsCfg.Virtual.Source = "ts"
			rsCfg.Virtual.PoolName = "/Common/customPool"
			rsCfg.Virtual.Destination = "/test/172.13.14.15:8080"

			rsCfg.Pools = Pools{
				Pool{
					Name:            "pool1",
					Members:         []PoolMember{mem1, mem2},
					MinimumMonitors: intstr.IntOrString{Type: 0, IntVal: 1},
				},
			}

			config := ResourceConfigRequest{
				ltmConfig:          make(LTMConfig),
				shareNodes:         true,
				gtmConfig:          GTMConfig{},
				defaultRouteDomain: 1,
			}

			zero := 0
			config.ltmConfig["default"] = &PartitionConfig{ResourceMap: make(ResourceMap), Priority: &zero}
			config.ltmConfig["default"].ResourceMap["crd_vs_172.13.14.15"] = rsCfg

			agentPostCfg := as3Handler.createAPIConfig(config, false, "", false)
			decl := agentPostCfg.data
			Expect(decl).ToNot(Equal(""), "Failed to Create AS3 Declaration")
			Expect(strings.Contains(decl, "adminState")).To(BeTrue())
			Expect(strings.Contains(decl, "connectionLimit")).To(BeTrue())
			Expect(strings.Contains(decl, "profileFTP")).To(BeTrue())

		})
		It("Ingresslink Declaration", func() {
			rsCfg := &ResourceConfig{}
			rsCfg.MetaData.Active = true
			rsCfg.MetaData.ResourceType = IngressLink
			rsCfg.Virtual.Name = "crd_il_172.13.14.16"
			rsCfg.Virtual.IpProtocol = "http"
			rsCfg.Virtual.ProfileL4 = "/Common/profileL4"
			rsCfg.Virtual.ProfileDOS = "/Common/profileDOS"
			rsCfg.Virtual.ProfileBotDefense = "/Common/profileBotDefense"
			rsCfg.Virtual.TCP.Client = "/Common/tcpClient"
			rsCfg.Virtual.TCP.Server = "/Common/tcpServer"
			rsCfg.Virtual.Profiles = ProfileRefs{
				ProfileRef{
					Name:      "serverssl",
					Partition: "Common",
					Context:   "udp",
				},
				ProfileRef{
					Name:         "clientssl",
					Partition:    "Common",
					Context:      "udp",
					BigIPProfile: true,
				},
			}
			rsCfg.Virtual.TranslateServerAddress = true
			rsCfg.Virtual.TranslateServerPort = true
			rsCfg.Virtual.Source = "inglink"
			rsCfg.Virtual.PoolName = "/Common/customPool"
			rsCfg.Virtual.Destination = "/test/172.13.14.5:8080"
			rsCfg.Pools = Pools{
				Pool{
					Name:            "pool1",
					Members:         []PoolMember{mem1, mem2},
					MinimumMonitors: intstr.IntOrString{Type: 0, IntVal: 1},
				},
			}

			config := ResourceConfigRequest{
				ltmConfig:          make(LTMConfig),
				shareNodes:         true,
				gtmConfig:          GTMConfig{},
				defaultRouteDomain: 1,
			}

			zero := 0
			config.ltmConfig["default"] = &PartitionConfig{ResourceMap: make(ResourceMap), Priority: &zero}
			config.ltmConfig["default"].ResourceMap["crd_il_172.13.14.16"] = rsCfg

			agentPostCfg := as3Handler.createAPIConfig(config, false, "", false)
			decl := agentPostCfg.data
			Expect(decl).ToNot(Equal(""), "Failed to Create AS3 Declaration")
			Expect(strings.Contains(decl, "adminState")).To(BeTrue())
			Expect(strings.Contains(decl, "connectionLimit")).To(BeTrue())
			Expect(strings.Contains(decl, "profileL4")).To(BeTrue())

		})
		It("Delete partition", func() {
			config := ResourceConfigRequest{
				ltmConfig:          make(LTMConfig),
				shareNodes:         true,
				gtmConfig:          GTMConfig{},
				defaultRouteDomain: 1,
			}

			zero := 0
			config.ltmConfig["default"] = &PartitionConfig{ResourceMap: make(ResourceMap), Priority: &zero}
			agentPostCfg := as3Handler.createAPIConfig(config, false, "", false)
			decl := agentPostCfg.data
			var as3Config map[string]interface{}
			_ = json.Unmarshal([]byte(decl), &as3Config)
			deletedTenantDecl := as3Tenant{
				"class": "Tenant",
			}
			adc := as3Config["declaration"].(map[string]interface{})

			Expect(agentPostCfg.incomingTenantDeclMap["default"]).To(Equal(deletedTenantDecl), "Failed to Create AS3 Declaration for deleted tenant")
			Expect(adc["default"]).To(Equal(map[string]interface{}(deletedTenantDecl)), "Failed to Create AS3 Declaration for deleted tenant")
		})
		It("Handles Persistence Methods", func() {
			svc := &as3Service{}
			// Default persistence methods
			defaultValues := []string{"cookie", "destination-address", "hash", "msrdp",
				"sip-info", "source-address", "tls-session-id", "universal"}
			for _, defaultValue := range defaultValues {
				svc.addPersistenceMethod(defaultValue)
				Expect(svc.PersistenceMethods).To(Equal(&[]as3MultiTypeParam{as3MultiTypeParam(defaultValue)}))
			}

			// Persistence methods with no value and None
			svc = &as3Service{}
			svc.addPersistenceMethod("")
			Expect(svc.PersistenceMethods).To(BeNil())
			svc.addPersistenceMethod("none")
			Expect(svc.PersistenceMethods).To(Equal(&[]as3MultiTypeParam{}))

			// Custom persistence methods
			svc.addPersistenceMethod("/Common/pm1")
			Expect(svc.PersistenceMethods).To(Equal(&[]as3MultiTypeParam{as3ResourcePointer{BigIP: "/Common/pm1"}}))
			svc.addPersistenceMethod("pm2")
			Expect(svc.PersistenceMethods).To(Equal(&[]as3MultiTypeParam{as3ResourcePointer{BigIP: "pm2"}}))
		})
	})

	Describe("Prepare AS3 Declaration with HAMode", func() {
		var as3Handler AS3Handler
		BeforeEach(func() {
			as3Handler = AS3Handler{
				AS3Parser:   &AS3Parser{},
				PostManager: &PostManager{declUpdate: sync.Mutex{}},
			}
		})
		It("VirtualServer Declaration", func() {
			config := ResourceConfigRequest{
				ltmConfig: make(LTMConfig),
			}
			agentPostCfg := as3Handler.createAPIConfig(config, false, "", false)
			decl := agentPostCfg.data
			Expect(decl).ToNot(Equal(""), "Failed to Create AS3 Declaration")
			Expect(strings.Contains(decl, "\"class\":\"Tenant\"")).To(BeTrue())
			Expect(strings.Contains(decl, "\"class\":\"AS3\"")).To(BeTrue())
			Expect(strings.Contains(decl, "\"class\":\"ADC\"")).To(BeTrue())
			Expect(strings.Contains(decl, "\"class\":\"Application\"")).To(BeTrue())
		})
	})

	Describe("GTM Config", func() {
		var as3Handler AS3Handler
		BeforeEach(func() {
			as3Handler = AS3Handler{
				AS3Parser:   &AS3Parser{},
				PostManager: &PostManager{declUpdate: sync.Mutex{}},
			}
			DEFAULT_PARTITION = "default"
		})
		// Commenting this test case
		// with new GTM partition support we will not delete partition, instead we flush contents
		//It("Empty GTM Config", func() {
		//	adc := as3ADC{}
		//	adc = agent.createAS3GTMConfigADC(ResourceConfigRequest{
		//		gtmConfig: GTMConfig{},
		//	}, adc)
		//
		//	Expect(len(adc)).To(BeZero(), "Invalid GTM Config")
		//})

		It("Empty GTM Partition Config / Delete Case", func() {
			adc := as3Handler.createLTMAndGTMConfigADC(ResourceConfigRequest{
				gtmConfig: GTMConfig{
					DEFAULT_PARTITION: GTMPartitionConfig{},
				},
			}, false, false)
			Expect(len(adc)).To(Equal(1), "Invalid GTM Config")

			Expect(adc).To(HaveKey(DEFAULT_PARTITION))
			tenant := adc[DEFAULT_PARTITION].(as3Tenant)

			Expect(tenant).To(HaveKey(as3SharedApplication))
			sharedApp := tenant[as3SharedApplication].(as3Application)
			Expect(len(sharedApp)).To(Equal(2))
			Expect(sharedApp).To(HaveKeyWithValue("class", "Application"))
			Expect(sharedApp).To(HaveKeyWithValue("template", "shared"))
		})

		It("Valid GTM Config", func() {
			monitors := []Monitor{
				{
					Name:     "pool1_monitor",
					Interval: 10,
					Timeout:  10,
					Type:     "http",
					Send:     "GET /health",
				},
			}
			gtmConfig := GTMConfig{
				DEFAULT_PARTITION: GTMPartitionConfig{
					WideIPs: map[string]WideIP{
						"test.com": {
							DomainName: "test.com",
							RecordType: "A",
							LBMethod:   "round-robin",
							Pools: []GSLBPool{
								{
									Name:       "pool1",
									RecordType: "A",
									LBMethod:   "round-robin",
									Members:    []string{"vs1", "vs2"},
									Monitors:   monitors,
								},
							},
						},
					},
				},
			}
			adc := as3Handler.createLTMAndGTMConfigADC(ResourceConfigRequest{gtmConfig: gtmConfig}, false, false)
			Expect(adc).To(HaveKey(DEFAULT_PARTITION))
			tenant := adc[DEFAULT_PARTITION].(as3Tenant)

			Expect(tenant).To(HaveKey(as3SharedApplication))
			sharedApp := tenant[as3SharedApplication].(as3Application)

			Expect(sharedApp).To(HaveKey("test.com"))
			Expect(sharedApp["test.com"].(as3GLSBDomain).Class).To(Equal("GSLB_Domain"))

			Expect(sharedApp).To(HaveKey("pool1"))
			Expect(sharedApp["pool1"].(as3GSLBPool).Class).To(Equal("GSLB_Pool"))

			Expect(sharedApp).To(HaveKey("pool1_monitor"))
			Expect(sharedApp["pool1_monitor"].(as3GSLBMonitor).Class).To(Equal("GSLB_Monitor"))
		})
	})

	Describe("Misc", func() {
		var as3Handler AS3Handler
		BeforeEach(func() {
			as3Handler = AS3Handler{
				AS3Parser:   &AS3Parser{},
				PostManager: &PostManager{declUpdate: sync.Mutex{}},
			}
		})
		It("Service Address declaration", func() {
			rsCfg := &ResourceConfig{
				ServiceAddress: []ServiceAddress{
					{
						ArpEnabled: true,
					},
				},
			}
			app := as3Application{}
			as3Handler.createServiceAddressDecl(rsCfg, "1.2.3.4", app)

			val, ok := app["crd_service_address_1_2_3_4"]
			Expect(ok).To(BeTrue())
			Expect(val).NotTo(BeNil())
		})
		It("Test Deleted Partition", func() {
			as3Handler.defaultPartition = "test"
			cisLabel := "test"
			deletedPartition := as3Handler.getDeletedTenantDeclaration("test", cisLabel, 0)
			Expect(deletedPartition[as3SharedApplication]).NotTo(BeNil())
			deletedPartition = as3Handler.getDeletedTenantDeclaration("default", cisLabel, 0)
			Expect(deletedPartition[as3SharedApplication]).To(BeNil())
		})
	})

	Describe("JSON comparision of AS3 declaration", func() {
		It("Verify with two empty declarations", func() {
			ok := DeepEqualJSON("", "")
			Expect(ok).To(BeTrue(), "Failed to compare empty declarations")
		})
		It("Verify with empty and non empty declarations", func() {
			cmcfg1 := `{"key": "value"}`
			ok := DeepEqualJSON("", as3Declaration(cmcfg1))
			Expect(ok).To(BeFalse())
			ok = DeepEqualJSON(as3Declaration(cmcfg1), "")
			Expect(ok).To(BeFalse())
		})
		It("Verify two equal JSONs", func() {
			ok := DeepEqualJSON(`{"key": "value"}`, `{"key": "value"}`)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("Poll tenant status", func() {
		var agent *Agent
		var config *agentPostConfig
		var mockBaseAPIHandler *BaseAPIHandler
		BeforeEach(func() {
			agent = &Agent{}
			mockBaseAPIHandler = newMockBaseAPIHandler()
			tenantDeclMap := make(map[string]as3Tenant)
			tenantResponseMap := make(map[string]tenantResponse)
			tenantResponseMap["test"] = tenantResponse{}
			tenantDeclMap["test"] = as3Tenant{
				"class":              "Tenant",
				"defaultRouteDomain": 0,
				as3SharedApplication: "shared",
				"label":              "cis2.x",
			}
			config = &agentPostConfig{
				acceptedTaskId: "123",
				reqMeta: requestMeta{
					id: 1,
				},
				as3APIURL:             "https://127.0.0.1/mgmt/shared/appsvcs/declare",
				data:                  `{"class": "AS3", "declaration": {"class": "ADC", "test": {"class": "Tenant", "testApp": {"class": "Application", "webcert":{"class": "Certificate", "certificate": "abc", "privateKey": "abc", "chainCA": "abc"}}}}}`,
				incomingTenantDeclMap: tenantDeclMap,
				tenantResponseMap:     tenantResponseMap,
			}
		})
		It("Verify tenant status", func() {
			mockBaseAPIHandler.httpClient, _ = getMockHttpClient([]responseCtx{{
				tenant: "test",
				status: http.StatusOK,
				body:   io.NopCloser(strings.NewReader("{\"results\": [{\"code\": 200, \"message\": \"success\", \"tenant\": \"test\"}], \"declaration\": {\"class\": \"ADC\", \"test\": {\"class\": \"Tenant\", \"testApp\": {\"class\": \"Application\", \"webcert\":{\"class\": \"Certificate\", \"certificate\": \"abc\", \"privateKey\": \"abc\", \"chainCA\": \"abc\"}}}}}")),
			}, {
				tenant: "test",
				status: http.StatusUnprocessableEntity,
				body:   io.NopCloser(strings.NewReader("{\"id\": \"123\"}")),
			}}, http.MethodGet)
			agent.APIHandler = &APIHandler{
				LTM: &LTMAPIHandler{
					BaseAPIHandler: mockBaseAPIHandler,
				},
				GTM: &GTMAPIHandler{
					BaseAPIHandler: mockBaseAPIHandler,
				},
			}

			// verify status ok
			agent.LTM.APIHandler.pollTenantStatus(config)
			Expect(config.acceptedTaskId).To(BeEmpty(), "Accepted task id should be empty")
			Expect(config.tenantResponseMap["test"].agentResponseCode).To(Equal(http.StatusOK), "Response code should be 200")

			config.acceptedTaskId = "123"
			agent.LTM.APIHandler.pollTenantStatus(config)
			Expect(config.acceptedTaskId).To(BeEmpty(), "Accepted task id should be empty")
			Expect(config.tenantResponseMap["test"].agentResponseCode).To(Equal(http.StatusUnprocessableEntity), "Response code should be 422")
		})
	})

})
