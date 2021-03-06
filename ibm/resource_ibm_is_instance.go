package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcclassicv1"
	"github.ibm.com/ibmcloud/vpc-go-sdk/vpcv1"
)

const (
	isInstanceName                    = "name"
	isInstanceKeys                    = "keys"
	isInstanceTags                    = "tags"
	isInstanceNetworkInterfaces       = "network_interfaces"
	isInstancePrimaryNetworkInterface = "primary_network_interface"
	isInstanceNicName                 = "name"
	isInstanceProfile                 = "profile"
	isInstanceNicPortSpeed            = "port_speed"
	isInstanceNicPrimaryIpv4Address   = "primary_ipv4_address"
	isInstanceNicPrimaryIpv6Address   = "primary_ipv6_address"
	isInstanceNicSecondaryAddress     = "secondary_addresses"
	isInstanceNicSecurityGroups       = "security_groups"
	isInstanceNicSubnet               = "subnet"
	isInstanceNicFloatingIPs          = "floating_ips"
	isInstanceUserData                = "user_data"
	isInstanceVolumes                 = "volumes"
	isInstanceVPC                     = "vpc"
	isInstanceZone                    = "zone"
	isInstanceBootVolume              = "boot_volume"
	isInstanceVolAttName              = "name"
	isInstanceVolAttVolume            = "volume"
	isInstanceVolAttVolAutoDelete     = "auto_delete"
	isInstanceVolAttVolCapacity       = "capacity"
	isInstanceVolAttVolIops           = "iops"
	isInstanceVolAttVolName           = "name"
	isInstanceVolAttVolBillingTerm    = "billing_term"
	isInstanceVolAttVolEncryptionKey  = "encryption_key"
	isInstanceVolAttVolType           = "type"
	isInstanceVolAttVolProfile        = "profile"
	isInstanceImage                   = "image"
	isInstanceCPU                     = "vcpu"
	isInstanceCPUArch                 = "architecture"
	isInstanceCPUCores                = "cores"
	isInstanceCPUCount                = "count"
	isInstanceGpu                     = "gpu"
	isInstanceGpuCores                = "cores"
	isInstanceGpuCount                = "count"
	isInstanceGpuManufacturer         = "manufacturer"
	isInstanceGpuMemory               = "memory"
	isInstanceGpuModel                = "model"
	isInstanceMemory                  = "memory"
	isInstanceStatus                  = "status"
	isInstanceGeneration              = "generation"

	isInstanceProvisioning     = "provisioning"
	isInstanceProvisioningDone = "done"
	isInstanceAvailable        = "available"
	isInstanceDeleting         = "deleting"
	isInstanceDeleteDone       = "done"
	isInstanceFailed           = "failed"

	isInstanceActionStatusStopping = "stopping"
	isInstanceActionStatusStopped  = "stopped"
	isInstanceStatusPending        = "pending"
	isInstanceStatusRunning        = "running"
	isInstanceStatusFailed         = "failed"

	isInstanceBootName       = "name"
	isInstanceBootSize       = "size"
	isInstanceBootIOPS       = "iops"
	isInstanceBootEncryption = "encryption"
	isInstanceBootProfile    = "profile"

	isInstanceVolumeAttachments = "volume_attachments"
	isInstanceVolumeAttaching   = "attaching"
	isInstanceVolumeAttached    = "attached"
	isInstanceVolumeDetaching   = "detaching"
	isInstanceResourceGroup     = "resource_group"
)

func resourceIBMISInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMisInstanceCreate,
		Read:     resourceIBMisInstanceRead,
		Update:   resourceIBMisInstanceUpdate,
		Delete:   resourceIBMisInstanceDelete,
		Exists:   resourceIBMisInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateISName,
				Description:  "Instance name",
			},

			isInstanceVPC: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "VPC id",
			},

			isInstanceZone: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Zone name",
			},

			isInstanceProfile: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Profile info",
			},

			isInstanceKeys: {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: applyOnce,
				Description:      "SSH key Ids for the instance",
			},

			isInstanceTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         resourceIBMVPCHash,
				Description: "list of tags for the instance",
			},

			isInstanceVolumeAttachments: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_crn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			isInstancePrimaryNetworkInterface: {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Primary Network interface info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceNicPortSpeed: {
							Type:             schema.TypeInt,
							Optional:         true,
							DiffSuppressFunc: applyOnce,
							Deprecated:       "This field is deprected",
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceNetworkInterfaces: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceNicName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceNicPrimaryIpv4Address: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceNicSecurityGroups: {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						isInstanceNicSubnet: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			isInstanceGeneration: {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				ValidateFunc:     validateGeneration,
				Removed:          "This field is removed",
			},

			isInstanceUserData: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "User data given for the instance",
			},

			isInstanceImage: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "image name",
			},

			isInstanceBootVolume: {
				Type:             schema.TypeList,
				DiffSuppressFunc: applyOnce,
				Optional:         true,
				Computed:         true,
				MaxItems:         1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceBootName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceBootEncryption: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						isInstanceBootSize: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceBootIOPS: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceBootProfile: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			isInstanceVolumes: {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of volumes",
			},

			isInstanceResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Instance resource group",
			},

			isInstanceCPU: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceCPUArch: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceCPUCount: {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			isInstanceGpu: {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "This field is deprecated",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceGpuCores: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceGpuCount: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceGpuMemory: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						isInstanceGpuManufacturer: {
							Type:     schema.TypeString,
							Computed: true,
						},
						isInstanceGpuModel: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			isInstanceMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance memory",
			},

			isInstanceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance status",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func classicInstanceCreate(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcclassicv1.InstancePrototype{
		Image: &vpcclassicv1.ImageIdentity{
			ID: &image,
		},
		Zone: &vpcclassicv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcclassicv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		Vpc: &vpcclassicv1.VPCIdentity{
			ID: &vpcID,
		},
	}

	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcclassicv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceBootName]
		namestr := name.(string)
		if ok {
			volTemplate.Name = &namestr
		}
		enc, ok := bootvol[isInstanceBootEncryption]
		encstr := enc.(string)
		if ok && encstr != "" {
			volTemplate.EncryptionKey = &vpcclassicv1.EncryptionKeyIdentity{
				Crn: &encstr,
			}
		}
		volcap := 100
		volcapint64 := int64(volcap)
		volprof := "general-purpose"
		volTemplate.Capacity = &volcapint64
		volTemplate.Profile = &vpcclassicv1.VolumeProfileIdentity{
			Name: &volprof,
		}

		deletebool := true
		instanceproto.BootVolumeAttachment = &vpcclassicv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}
	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcclassicv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcclassicv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, ok := primnic[isInstanceNicName]
		namestr := name.(string)
		if ok {
			primnicobj.Name = &namestr
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
		if ok {
			secgrpSet := secgrpintf.(*schema.Set)
			if secgrpSet.Len() != 0 {
				var secgrpobjs = make([]vpcclassicv1.SecurityGroupIdentityIntf, secgrpSet.Len())
				for i, secgrpIntf := range secgrpSet.List() {
					secgrpIntfstr := secgrpIntf.(string)
					secgrpobjs[i] = &vpcclassicv1.SecurityGroupIdentity{
						ID: &secgrpIntfstr,
					}
				}
				primnicobj.SecurityGroups = secgrpobjs
			}
		}

		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcclassicv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcclassicv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcclassicv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
			if ok {
				secgrpSet := secgrpintf.(*schema.Set)
				if secgrpSet.Len() != 0 {
					var secgrpobjs = make([]vpcclassicv1.SecurityGroupIdentityIntf, secgrpSet.Len())
					for i, secgrpIntf := range secgrpSet.List() {
						secgrpIntfstr := secgrpIntf.(string)
						secgrpobjs[i] = &vpcclassicv1.SecurityGroupIdentity{
							ID: &secgrpIntfstr,
						}
					}
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcclassicv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcclassicv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcclassicv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	options := &vpcclassicv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}
	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForClassicInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc instance (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func instanceCreate(d *schema.ResourceData, meta interface{}, profile, name, vpcID, zone, image string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	instanceproto := &vpcv1.InstancePrototype{
		Image: &vpcv1.ImageIdentity{
			ID: &image,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &zone,
		},
		Profile: &vpcv1.InstanceProfileIdentity{
			Name: &profile,
		},
		Name: &name,
		Vpc: &vpcv1.VPCIdentity{
			ID: &vpcID,
		},
	}
	if boot, ok := d.GetOk(isInstanceBootVolume); ok {
		bootvol := boot.([]interface{})[0].(map[string]interface{})
		var volTemplate = &vpcv1.VolumePrototypeInstanceByImageContext{}
		name, ok := bootvol[isInstanceBootName]
		namestr := name.(string)
		if ok {
			volTemplate.Name = &namestr
		}
		// enc, ok := bootvol[isInstanceBootEncryption]
		// encstr := enc.(string)
		// if ok && encstr != "" {
		// 	volTemplate.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
		// 		Crn: &encstr,
		// 	}
		// }
		volcap := 100
		volcapint64 := int64(volcap)
		volprof := "general-purpose"
		volTemplate.Capacity = &volcapint64
		volTemplate.Profile = &vpcv1.VolumeProfileIdentity{
			Name: &volprof,
		}
		deletebool := true
		instanceproto.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
			DeleteVolumeOnInstanceDelete: &deletebool,
			Volume:                       volTemplate,
		}
	}

	if primnicintf, ok := d.GetOk(isInstancePrimaryNetworkInterface); ok {
		primnic := primnicintf.([]interface{})[0].(map[string]interface{})
		subnetintf, _ := primnic[isInstanceNicSubnet]
		subnetintfstr := subnetintf.(string)
		var primnicobj = &vpcv1.NetworkInterfacePrototype{}
		primnicobj.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnetintfstr,
		}
		name, _ := primnic[isInstanceNicName]
		namestr := name.(string)
		if namestr != "" {
			primnicobj.Name = &namestr
		}
		secgrpintf, ok := primnic[isInstanceNicSecurityGroups]
		if ok {
			secgrpSet := secgrpintf.(*schema.Set)
			if secgrpSet.Len() != 0 {
				var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
				for i, secgrpIntf := range secgrpSet.List() {
					secgrpIntfstr := secgrpIntf.(string)
					secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
						ID: &secgrpIntfstr,
					}
				}
				primnicobj.SecurityGroups = secgrpobjs
			}
		}
		instanceproto.PrimaryNetworkInterface = primnicobj
	}

	if nicsintf, ok := d.GetOk(isInstanceNetworkInterfaces); ok {
		nics := nicsintf.([]interface{})
		var intfs []vpcv1.NetworkInterfacePrototype
		for _, resource := range nics {
			nic := resource.(map[string]interface{})
			nwInterface := &vpcv1.NetworkInterfacePrototype{}
			subnetintf, _ := nic[isInstanceNicSubnet]
			subnetintfstr := subnetintf.(string)
			nwInterface.Subnet = &vpcv1.SubnetIdentity{
				ID: &subnetintfstr,
			}
			name, ok := nic[isInstanceNicName]
			namestr := name.(string)
			if ok && namestr != "" {
				nwInterface.Name = &namestr
			}
			secgrpintf, ok := nic[isInstanceNicSecurityGroups]
			if ok {
				secgrpSet := secgrpintf.(*schema.Set)
				if secgrpSet.Len() != 0 {
					var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
					for i, secgrpIntf := range secgrpSet.List() {
						secgrpIntfstr := secgrpIntf.(string)
						secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
							ID: &secgrpIntfstr,
						}
					}
					nwInterface.SecurityGroups = secgrpobjs
				}
			}
			intfs = append(intfs, *nwInterface)
		}
		instanceproto.NetworkInterfaces = intfs
	}

	keySet := d.Get(isInstanceKeys).(*schema.Set)
	if keySet.Len() != 0 {
		keyobjs := make([]vpcv1.KeyIdentityIntf, keySet.Len())
		for i, key := range keySet.List() {
			keystr := key.(string)
			keyobjs[i] = &vpcv1.KeyIdentity{
				ID: &keystr,
			}
		}
		instanceproto.Keys = keyobjs
	}

	if userdata, ok := d.GetOk(isInstanceUserData); ok {
		userdatastr := userdata.(string)
		instanceproto.UserData = &userdatastr
	}

	if grp, ok := d.GetOk(isInstanceResourceGroup); ok {
		grpstr := grp.(string)
		instanceproto.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &grpstr,
		}

	}

	options := &vpcv1.CreateInstanceOptions{
		InstancePrototype: instanceproto,
	}

	instance, response, err := sess.CreateInstance(options)
	if err != nil {
		log.Printf("[DEBUG] Instance err %s\n%s", err, response)
		return err
	}
	d.SetId(*instance.ID)

	log.Printf("[INFO] Instance : %s", *instance.ID)
	d.Set(isInstanceStatus, instance.Status)

	_, err = isWaitForInstanceAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isInstanceTags); ok || v != "" {
		oldList, newList := d.GetChange(isInstanceTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc instance (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	profile := d.Get(isInstanceProfile).(string)
	name := d.Get(isInstanceName).(string)
	vpcID := d.Get(isInstanceVPC).(string)
	zone := d.Get(isInstanceZone).(string)
	image := d.Get(isInstanceImage).(string)

	if userDetails.generation == 1 {
		err := classicInstanceCreate(d, meta, profile, name, vpcID, zone, image)
		if err != nil {
			return err
		}
	} else {
		err := instanceCreate(d, meta, profile, name, vpcID, zone, image)
		if err != nil {
			return err
		}
	}

	return resourceIBMisInstanceUpdate(d, meta)
}

func isWaitForClassicInstanceAvailable(instanceC *vpcclassicv1.VpcClassicV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for instance (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isInstanceProvisioning},
		Target:     []string{isInstanceStatusRunning, "available", "failed", ""},
		Refresh:    isClassicInstanceRefreshFunc(instanceC, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForInstanceAvailable(instanceC *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for instance (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isInstanceProvisioning},
		Target:     []string{isInstanceStatusRunning, "available", "failed", ""},
		Refresh:    isInstanceRefreshFunc(instanceC, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicInstanceRefreshFunc(instanceC *vpcclassicv1.VpcClassicV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcclassicv1.GetInstanceOptions{
			ID: &id,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting instance: %s\n%s", err, response)
		}

		d.Set(isInstanceStatus, *instance.Status)

		if *instance.Status == "available" || *instance.Status == "failed" || *instance.Status == "running" {
			return instance, *instance.Status, nil
		}

		return instance, isInstanceProvisioning, nil
	}
}

func isInstanceRefreshFunc(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcv1.GetInstanceOptions{
			ID: &id,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
		}
		d.Set(isInstanceStatus, *instance.Status)

		if *instance.Status == "available" || *instance.Status == "failed" || *instance.Status == "running" {
			return instance, *instance.Status, nil
		}

		return instance, isInstanceProvisioning, nil
	}
}

func resourceIBMisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	ID := d.Id()
	if userDetails.generation == 1 {
		err := classicInstanceGet(d, meta, ID)
		if err != nil {
			return err
		}
	} else {
		err := instanceGet(d, meta, ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicInstanceGet(d *schema.ResourceData, meta interface{}, id string) error {
	instanceC, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	getinsOptions := &vpcclassicv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
	}
	d.Set(isInstanceName, *instance.Name)
	if instance.Profile != nil {
		d.Set(isInstanceProfile, *instance.Profile.Name)
	}
	cpuList := make([]map[string]interface{}, 0)
	if instance.Vcpu != nil {
		currentCPU := map[string]interface{}{}
		currentCPU[isInstanceCPUArch] = *instance.Vcpu.Architecture
		currentCPU[isInstanceCPUCount] = *instance.Vcpu.Count
		cpuList = append(cpuList, currentCPU)
	}
	d.Set(isInstanceCPU, cpuList)

	d.Set(isInstanceMemory, *instance.Memory)
	gpuList := make([]map[string]interface{}, 0)
	// if instance.Gpu != nil {
	// 	currentGpu := map[string]interface{}{}
	// 	currentGpu[isInstanceGpuManufacturer] = instance.Gpu.Manufacturer
	// 	currentGpu[isInstanceGpuModel] = instance.Gpu.Model
	// 	currentGpu[isInstanceGpuCores] = instance.Gpu.Cores
	// 	currentGpu[isInstanceGpuCount] = instance.Gpu.Count
	// 	currentGpu[isInstanceGpuMemory] = instance.Gpu.Memory
	// 	gpuList = append(gpuList, currentGpu)

	// }
	d.Set(isInstanceGpu, gpuList)

	if instance.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
		currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name
		currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
		getnicoptions := &vpcclassicv1.GetNetworkInterfaceOptions{
			InstanceID: &id,
			ID:         instance.PrimaryNetworkInterface.ID,
		}
		insnic, response, err := instanceC.GetNetworkInterface(getnicoptions)
		if err != nil {
			return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
		}
		currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
		if len(insnic.SecurityGroups) != 0 {
			secgrpList := []string{}
			for i := 0; i < len(insnic.SecurityGroups); i++ {
				secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
			}
			currentPrimNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
		}

		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isInstancePrimaryNetworkInterface, primaryNicList)
	}

	if instance.NetworkInterfaces != nil {
		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range instance.NetworkInterfaces {
			if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
				currentNic := map[string]interface{}{}
				currentNic["id"] = *intfc.ID
				currentNic[isInstanceNicName] = *intfc.Name
				currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
				getnicoptions := &vpcclassicv1.GetNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         intfc.ID,
				}
				insnic, response, err := instanceC.GetNetworkInterface(getnicoptions)
				if err != nil {
					return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
				}
				currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
				if len(insnic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(insnic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
					}
					currentNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
				}
				interfacesList = append(interfacesList, currentNic)

			}
		}

		d.Set(isInstanceNetworkInterfaces, interfacesList)
	}

	if instance.Image != nil {
		d.Set(isInstanceImage, *instance.Image.ID)
	}

	d.Set(isInstanceStatus, *instance.Status)
	d.Set(isInstanceVPC, *instance.Vpc.ID)
	d.Set(isInstanceZone, *instance.Zone.Name)

	var volumes []string
	volumes = make([]string, 0)
	if instance.VolumeAttachments != nil {
		for _, volume := range instance.VolumeAttachments {
			if volume.Volume != nil && *volume.Volume.ID != *instance.BootVolumeAttachment.Volume.ID {
				volumes = append(volumes, *volume.Volume.ID)
			}
		}
	}
	d.Set(isInstanceVolumes, newStringSet(schema.HashString, volumes))
	if instance.VolumeAttachments != nil {
		volList := make([]map[string]interface{}, 0)
		for _, volume := range instance.VolumeAttachments {
			vol := map[string]interface{}{}
			if volume.Volume != nil {
				vol["id"] = *volume.ID
				vol["volume_id"] = *volume.Volume.ID
				vol["name"] = *volume.Name
				vol["volume_name"] = *volume.Volume.Name
				vol["volume_crn"] = *volume.Volume.Crn
				volList = append(volList, vol)
			}
		}
		d.Set(isInstanceVolumeAttachments, volList)
	}
	if instance.BootVolumeAttachment != nil {
		bootVolList := make([]map[string]interface{}, 0)
		bootVol := map[string]interface{}{}
		bootVol[isInstanceBootName] = *instance.BootVolumeAttachment.Name
		// getvolattoptions := &vpcclassicv1.GetVolumeAttachmentOptions{
		// 	InstanceID: &ID,
		// 	ID:         instance.BootVolumeAttachment.Volume.ID,
		// }
		// vol, _, err := instanceC.GetVolumeAttachment(getvolattoptions)
		// if err != nil {
		// 	return fmt.Errorf("Error while retrieving boot volume %s for instance %s: %v", getvolattoptions.ID, d.Id(), err)
		// }
		if instance.BootVolumeAttachment.Volume.Crn != nil {
			bootVol[isInstanceBootEncryption] = *instance.BootVolumeAttachment.Volume.Crn
		}
		// bootVol[isInstanceBootSize] = instance.BootVolumeAttachment.Capacity
		// bootVol[isInstanceBootIOPS] = instance.BootVolumeAttachment.Iops
		// bootVol[isInstanceBootProfile] = instance.BootVolumeAttachment.Name
		bootVolList = append(bootVolList, bootVol)

		d.Set(isInstanceBootVolume, bootVolList)
	}
	tags, err := GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Instance (%s) tags: %s", d.Id(), err)
	}
	d.Set(isInstanceTags, tags)

	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc/compute/vs")
	d.Set(ResourceName, instance.Name)
	d.Set(ResourceCRN, instance.Crn)
	d.Set(ResourceStatus, instance.Status)
	if instance.ResourceGroup != nil {
		d.Set(isInstanceResourceGroup, instance.ResourceGroup.ID)
		d.Set(ResourceGroupName, instance.ResourceGroup.ID)
	}
	return nil
}

func instanceGet(d *schema.ResourceData, meta interface{}, id string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
	}
	d.Set(isInstanceName, *instance.Name)
	if instance.Profile != nil {
		d.Set(isInstanceProfile, *instance.Profile.Name)
	}
	cpuList := make([]map[string]interface{}, 0)
	if instance.Vcpu != nil {
		currentCPU := map[string]interface{}{}
		currentCPU[isInstanceCPUArch] = *instance.Vcpu.Architecture
		currentCPU[isInstanceCPUCount] = *instance.Vcpu.Count
		cpuList = append(cpuList, currentCPU)
	}
	d.Set(isInstanceCPU, cpuList)

	d.Set(isInstanceMemory, *instance.Memory)
	gpuList := make([]map[string]interface{}, 0)
	// if instance.Gpu != nil {
	// 	currentGpu := map[string]interface{}{}
	// 	currentGpu[isInstanceGpuManufacturer] = instance.Gpu.Manufacturer
	// 	currentGpu[isInstanceGpuModel] = instance.Gpu.Model
	// 	currentGpu[isInstanceGpuCores] = instance.Gpu.Cores
	// 	currentGpu[isInstanceGpuCount] = instance.Gpu.Count
	// 	currentGpu[isInstanceGpuMemory] = instance.Gpu.Memory
	// 	gpuList = append(gpuList, currentGpu)

	// }
	d.Set(isInstanceGpu, gpuList)

	if instance.PrimaryNetworkInterface != nil {
		primaryNicList := make([]map[string]interface{}, 0)
		currentPrimNic := map[string]interface{}{}
		currentPrimNic["id"] = *instance.PrimaryNetworkInterface.ID
		currentPrimNic[isInstanceNicName] = *instance.PrimaryNetworkInterface.Name
		currentPrimNic[isInstanceNicPrimaryIpv4Address] = *instance.PrimaryNetworkInterface.PrimaryIpv4Address
		getnicoptions := &vpcv1.GetNetworkInterfaceOptions{
			InstanceID: &id,
			ID:         instance.PrimaryNetworkInterface.ID,
		}
		insnic, response, err := instanceC.GetNetworkInterface(getnicoptions)
		if err != nil {
			return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
		}
		currentPrimNic[isInstanceNicSubnet] = *insnic.Subnet.ID
		if len(insnic.SecurityGroups) != 0 {
			secgrpList := []string{}
			for i := 0; i < len(insnic.SecurityGroups); i++ {
				secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
			}
			currentPrimNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
		}

		primaryNicList = append(primaryNicList, currentPrimNic)
		d.Set(isInstancePrimaryNetworkInterface, primaryNicList)
	}

	if instance.NetworkInterfaces != nil {
		interfacesList := make([]map[string]interface{}, 0)
		for _, intfc := range instance.NetworkInterfaces {
			if *intfc.ID != *instance.PrimaryNetworkInterface.ID {
				currentNic := map[string]interface{}{}
				currentNic["id"] = *intfc.ID
				currentNic[isInstanceNicName] = *intfc.Name
				currentNic[isInstanceNicPrimaryIpv4Address] = *intfc.PrimaryIpv4Address
				getnicoptions := &vpcv1.GetNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         intfc.ID,
				}
				insnic, response, err := instanceC.GetNetworkInterface(getnicoptions)
				if err != nil {
					return fmt.Errorf("Error getting network interfaces attached to the instance %s\n%s", err, response)
				}
				currentNic[isInstanceNicSubnet] = *insnic.Subnet.ID
				if len(insnic.SecurityGroups) != 0 {
					secgrpList := []string{}
					for i := 0; i < len(insnic.SecurityGroups); i++ {
						secgrpList = append(secgrpList, string(*(insnic.SecurityGroups[i].ID)))
					}
					currentNic[isInstanceNicSecurityGroups] = newStringSet(schema.HashString, secgrpList)
				}
				interfacesList = append(interfacesList, currentNic)

			}
		}

		d.Set(isInstanceNetworkInterfaces, interfacesList)
	}

	if instance.Image != nil {
		d.Set(isInstanceImage, *instance.Image.ID)
	}

	d.Set(isInstanceStatus, *instance.Status)
	d.Set(isInstanceVPC, *instance.Vpc.ID)
	d.Set(isInstanceZone, *instance.Zone.Name)

	var volumes []string
	volumes = make([]string, 0)
	if instance.VolumeAttachments != nil {
		for _, volume := range instance.VolumeAttachments {
			if volume.Volume != nil && *volume.Volume.ID != *instance.BootVolumeAttachment.Volume.ID {
				volumes = append(volumes, *volume.Volume.ID)
			}
		}
	}
	d.Set(isInstanceVolumes, newStringSet(schema.HashString, volumes))
	if instance.VolumeAttachments != nil {
		volList := make([]map[string]interface{}, 0)
		for _, volume := range instance.VolumeAttachments {
			vol := map[string]interface{}{}
			if volume.Volume != nil {
				vol["id"] = *volume.ID
				vol["volume_id"] = *volume.Volume.ID
				vol["name"] = *volume.Name
				vol["volume_name"] = *volume.Volume.Name
				vol["volume_crn"] = *volume.Volume.Crn
				volList = append(volList, vol)
			}
		}
		d.Set(isInstanceVolumeAttachments, volList)
	}
	if instance.BootVolumeAttachment != nil {
		bootVolList := make([]map[string]interface{}, 0)
		bootVol := map[string]interface{}{}
		bootVol[isInstanceBootName] = *instance.BootVolumeAttachment.Name
		// getvolattoptions := &vpcclassicv1.GetVolumeAttachmentOptions{
		// 	InstanceID: &ID,
		// 	ID:         instance.BootVolumeAttachment.Volume.ID,
		// }
		// vol, _, err := instanceC.GetVolumeAttachment(getvolattoptions)
		// if err != nil {
		// 	return fmt.Errorf("Error while retrieving boot volume %s for instance %s: %v", getvolattoptions.ID, d.Id(), err)
		// }
		if instance.BootVolumeAttachment.Volume.Crn != nil {
			bootVol[isInstanceBootEncryption] = *instance.BootVolumeAttachment.Volume.Crn
		}
		// bootVol[isInstanceBootSize] = instance.BootVolumeAttachment.Capacity
		// bootVol[isInstanceBootIOPS] = instance.BootVolumeAttachment.Iops
		// bootVol[isInstanceBootProfile] = instance.BootVolumeAttachment.Name
		bootVolList = append(bootVolList, bootVol)

		d.Set(isInstanceBootVolume, bootVolList)
	}
	tags, err := GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Instance (%s) tags: %s", d.Id(), err)
	}
	d.Set(isInstanceTags, tags)

	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/vpc-ext/compute/vs")
	d.Set(ResourceName, *instance.Name)
	d.Set(ResourceCRN, *instance.Crn)
	d.Set(ResourceStatus, *instance.Status)
	if instance.ResourceGroup != nil {
		d.Set(isInstanceResourceGroup, *instance.ResourceGroup.ID)
		d.Set(ResourceGroupName, *instance.ResourceGroup.Name)
	}
	return nil
}

func classicInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	id := d.Id()
	if d.HasChange(isInstanceVolumes) {
		ovs, nvs := d.GetChange(isInstanceVolumes)
		ov := ovs.(*schema.Set)
		nv := nvs.(*schema.Set)

		remove := expandStringList(ov.Difference(nv).List())
		add := expandStringList(nv.Difference(ov).List())

		if len(add) > 0 {
			for i := range add {
				createvolattoptions := &vpcclassicv1.CreateVolumeAttachmentOptions{
					InstanceID: &id,
					Volume: &vpcclassicv1.VolumeIdentity{
						ID: &add[i],
					},
				}
				vol, response, err := instanceC.CreateVolumeAttachment(createvolattoptions)
				if err != nil {
					return fmt.Errorf("Error while attaching volume %q for instance %s\n%s: %q", add[i], d.Id(), err, response)
				}
				_, err = isWaitForClassicInstanceVolumeAttached(instanceC, d, id, *vol.ID)
				if err != nil {
					return err
				}
			}

		}
		if len(remove) > 0 {
			for i := range remove {
				listvolattoptions := &vpcclassicv1.ListVolumeAttachmentsOptions{
					InstanceID: &id,
				}
				vols, _, err := instanceC.ListVolumeAttachments(listvolattoptions)
				if err != nil {
					return err
				}
				for _, vol := range vols.VolumeAttachments {
					if *vol.Volume.ID == remove[i] {
						delvolattoptions := &vpcclassicv1.DeleteVolumeAttachmentOptions{
							InstanceID: &id,
							ID:         vol.ID,
						}
						response, err := instanceC.DeleteVolumeAttachment(delvolattoptions)
						if err != nil {
							return fmt.Errorf("Error while removing volume %q for instance %s\n%s: %q", remove[i], d.Id(), err, response)
						}
						_, err = isWaitForClassicInstanceVolumeDetached(instanceC, d, d.Id(), *vol.ID)
						if err != nil {
							return err
						}
						break
					}
				}
			}
		}
	}

	if d.HasChange("primary_network_interface.0.security_groups") && !d.IsNewResource() {
		ovs, nvs := d.GetChange("primary_network_interface.0.security_groups")
		ov := ovs.(*schema.Set)
		nv := nvs.(*schema.Set)
		remove := expandStringList(ov.Difference(nv).List())
		add := expandStringList(nv.Difference(ov).List())
		if len(add) > 0 {
			networkID := d.Get("primary_network_interface.0.id").(string)
			for i := range add {
				createsgnicoptions := &vpcclassicv1.CreateSecurityGroupNetworkInterfaceBindingOptions{
					SecurityGroupID: &add[i],
					ID:              &networkID,
				}
				_, response, err := instanceC.CreateSecurityGroupNetworkInterfaceBinding(createsgnicoptions)
				if err != nil {
					return fmt.Errorf("Error while creating security group %q for primary network interface of instance %s\n%s: %q", add[i], d.Id(), err, response)
				}
				_, err = isWaitForClassicInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return err
				}
			}

		}
		if len(remove) > 0 {
			networkID := d.Get("primary_network_interface.0.id").(string)
			for i := range remove {
				deletesgnicoptions := &vpcclassicv1.DeleteSecurityGroupNetworkInterfaceBindingOptions{
					SecurityGroupID: &remove[i],
					ID:              &networkID,
				}
				response, err := instanceC.DeleteSecurityGroupNetworkInterfaceBinding(deletesgnicoptions)
				if err != nil {
					return fmt.Errorf("Error while removing security group %q for primary network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
				}
				_, err = isWaitForClassicInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return err
				}
			}
		}
	}

	// if d.HasChange("primary_network_interface.0.name") && !d.IsNewResource() {
	// 	newName := d.Get("primary_network_interface.0.name").(string)
	// 	networkID := d.Get("primary_network_interface.0.id").(string)
	// 	_, err := instanceC.UpdateInterface(d.Id(), networkID, newName, 0)
	// 	if err != nil {
	// 		return fmt.Errorf("Error while updating name %s for primary network interface of instance %s: %q", newName, d.Id(), err)
	// 	}
	// 	_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	if d.HasChange(isInstanceNetworkInterfaces) && !d.IsNewResource() {
		nics := d.Get(isInstanceNetworkInterfaces).([]interface{})
		for i := range nics {
			securitygrpKey := fmt.Sprintf("network_interfaces.%d.security_groups", i)
			// networkNameKey := fmt.Sprintf("network_interfaces.%d.name", i)
			if d.HasChange(securitygrpKey) {
				ovs, nvs := d.GetChange(securitygrpKey)
				ov := ovs.(*schema.Set)
				nv := nvs.(*schema.Set)
				remove := expandStringList(ov.Difference(nv).List())
				add := expandStringList(nv.Difference(ov).List())
				if len(add) > 0 {
					networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
					networkID := d.Get(networkIDKey).(string)
					for i := range add {
						createsgnicoptions := &vpcclassicv1.CreateSecurityGroupNetworkInterfaceBindingOptions{
							SecurityGroupID: &add[i],
							ID:              &networkID,
						}
						_, response, err := instanceC.CreateSecurityGroupNetworkInterfaceBinding(createsgnicoptions)
						if err != nil {
							return fmt.Errorf("Error while creating security group %q for network interface of instance %s\n%s: %q", add[i], d.Id(), err, response)
						}
						_, err = isWaitForClassicInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
						if err != nil {
							return err
						}
					}

				}
				if len(remove) > 0 {
					networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
					networkID := d.Get(networkIDKey).(string)
					for i := range remove {
						deletesgnicoptions := &vpcclassicv1.DeleteSecurityGroupNetworkInterfaceBindingOptions{
							SecurityGroupID: &remove[i],
							ID:              &networkID,
						}
						response, err := instanceC.DeleteSecurityGroupNetworkInterfaceBinding(deletesgnicoptions)
						if err != nil {
							return fmt.Errorf("Error while removing security group %q for network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
						}
						_, err = isWaitForClassicInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
						if err != nil {
							return err
						}
					}
				}

			}

			// if d.HasChange(networkNameKey) {
			// 	newName := d.Get(networkNameKey).(string)
			// 	networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
			// 	networkID := d.Get(networkIDKey).(string)
			// 	_, err := instanceC.UpdateInterface(d.Id(), networkID, newName, 0)
			// 	if err != nil {
			// 		return fmt.Errorf("Error while updating name %s for network interface %s of instance %s: %q", newName, networkID, d.Id(), err)
			// 	}
			// 	_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		}

	}

	if d.HasChange(isInstanceName) {
		name := d.Get(isInstanceName).(string)
		updnetoptions := &vpcclassicv1.UpdateInstanceOptions{
			ID:   &id,
			Name: &name,
		}
		_, _, err = instanceC.UpdateInstance(updnetoptions)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceTags) {
		getinsOptions := &vpcclassicv1.GetInstanceOptions{
			ID: &id,
		}
		instance, response, err := instanceC.GetInstance(getinsOptions)
		if err != nil {
			log.Printf("Error Getting Instance: %s\n%s", err, response)
		}
		oldList, newList := d.GetChange(isInstanceTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Instance (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func instanceUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}
	id := d.Id()
	if d.HasChange(isInstanceVolumes) {
		ovs, nvs := d.GetChange(isInstanceVolumes)
		ov := ovs.(*schema.Set)
		nv := nvs.(*schema.Set)

		remove := expandStringList(ov.Difference(nv).List())
		add := expandStringList(nv.Difference(ov).List())

		if len(add) > 0 {
			for i := range add {
				createvolattoptions := &vpcv1.CreateVolumeAttachmentOptions{
					InstanceID: &id,
					Volume: &vpcv1.VolumeIdentity{
						ID: &add[i],
					},
				}
				vol, _, err := instanceC.CreateVolumeAttachment(createvolattoptions)
				if err != nil {
					return fmt.Errorf("Error while attaching volume %q for instance %s: %q", add[i], d.Id(), err)
				}
				_, err = isWaitForInstanceVolumeAttached(instanceC, d, id, *vol.ID)
				if err != nil {
					return err
				}
			}

		}
		if len(remove) > 0 {
			for i := range remove {
				listvolattoptions := &vpcv1.ListVolumeAttachmentsOptions{
					InstanceID: &id,
				}
				vols, _, err := instanceC.ListVolumeAttachments(listvolattoptions)
				if err != nil {
					return err
				}
				for _, vol := range vols.VolumeAttachments {
					if *vol.Volume.ID == remove[i] {
						delvolattoptions := &vpcv1.DeleteVolumeAttachmentOptions{
							InstanceID: &id,
							ID:         vol.ID,
						}
						_, err := instanceC.DeleteVolumeAttachment(delvolattoptions)
						if err != nil {
							return fmt.Errorf("Error while removing volume %q for instance %s: %q", remove[i], d.Id(), err)
						}
						_, err = isWaitForInstanceVolumeDetached(instanceC, d, d.Id(), *vol.ID)
						if err != nil {
							return err
						}
						break
					}
				}
			}
		}
	}

	if d.HasChange("primary_network_interface.0.security_groups") && !d.IsNewResource() {
		ovs, nvs := d.GetChange("primary_network_interface.0.security_groups")
		ov := ovs.(*schema.Set)
		nv := nvs.(*schema.Set)
		remove := expandStringList(ov.Difference(nv).List())
		add := expandStringList(nv.Difference(ov).List())
		if len(add) > 0 {
			networkID := d.Get("primary_network_interface.0.id").(string)
			for i := range add {
				createsgnicoptions := &vpcv1.CreateSecurityGroupNetworkInterfaceBindingOptions{
					SecurityGroupID: &add[i],
					ID:              &networkID,
				}
				_, response, err := instanceC.CreateSecurityGroupNetworkInterfaceBinding(createsgnicoptions)
				if err != nil {
					return fmt.Errorf("Error while creating security group %q for primary network interface of instance %s\n%s: %q", add[i], d.Id(), err, response)
				}
				_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return err
				}
			}

		}
		if len(remove) > 0 {
			networkID := d.Get("primary_network_interface.0.id").(string)
			for i := range remove {
				deletesgnicoptions := &vpcv1.DeleteSecurityGroupNetworkInterfaceBindingOptions{
					SecurityGroupID: &remove[i],
					ID:              &networkID,
				}
				response, err := instanceC.DeleteSecurityGroupNetworkInterfaceBinding(deletesgnicoptions)
				if err != nil {
					return fmt.Errorf("Error while removing security group %q for primary network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
				}
				_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("primary_network_interface.0.name") && !d.IsNewResource() {
		newName := d.Get("primary_network_interface.0.name").(string)
		networkID := d.Get("primary_network_interface.0.id").(string)
		updatepnicfoptions := &vpcv1.UpdateNetworkInterfaceOptions{
			InstanceID: &id,
			ID:         &networkID,
			Name:       &newName,
		}
		_, response, err := instanceC.UpdateNetworkInterface(updatepnicfoptions)
		if err != nil {
			return fmt.Errorf("Error while updating name %s for primary network interface of instance %s\n%s: %q", newName, d.Id(), err, response)
		}
		_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return err
		}
	}

	if d.HasChange(isInstanceNetworkInterfaces) && !d.IsNewResource() {
		nics := d.Get(isInstanceNetworkInterfaces).([]interface{})
		for i := range nics {
			securitygrpKey := fmt.Sprintf("network_interfaces.%d.security_groups", i)
			networkNameKey := fmt.Sprintf("network_interfaces.%d.name", i)
			if d.HasChange(securitygrpKey) {
				ovs, nvs := d.GetChange(securitygrpKey)
				ov := ovs.(*schema.Set)
				nv := nvs.(*schema.Set)
				remove := expandStringList(ov.Difference(nv).List())
				add := expandStringList(nv.Difference(ov).List())
				if len(add) > 0 {
					networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
					networkID := d.Get(networkIDKey).(string)
					for i := range add {
						createsgnicoptions := &vpcv1.CreateSecurityGroupNetworkInterfaceBindingOptions{
							SecurityGroupID: &add[i],
							ID:              &networkID,
						}
						_, response, err := instanceC.CreateSecurityGroupNetworkInterfaceBinding(createsgnicoptions)
						if err != nil {
							return fmt.Errorf("Error while creating security group %q for network interface of instance %s\n%s: %q", add[i], d.Id(), err, response)
						}
						_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
						if err != nil {
							return err
						}
					}

				}
				if len(remove) > 0 {
					networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
					networkID := d.Get(networkIDKey).(string)
					for i := range remove {
						deletesgnicoptions := &vpcv1.DeleteSecurityGroupNetworkInterfaceBindingOptions{
							SecurityGroupID: &remove[i],
							ID:              &networkID,
						}
						response, err := instanceC.DeleteSecurityGroupNetworkInterfaceBinding(deletesgnicoptions)
						if err != nil {
							return fmt.Errorf("Error while removing security group %q for network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response)
						}
						_, err = isWaitForInstanceAvailable(instanceC, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
						if err != nil {
							return err
						}
					}
				}

			}

			if d.HasChange(networkNameKey) {
				newName := d.Get(networkNameKey).(string)
				networkIDKey := fmt.Sprintf("network_interfaces.%d.id", i)
				networkID := d.Get(networkIDKey).(string)
				updatepnicfoptions := &vpcv1.UpdateNetworkInterfaceOptions{
					InstanceID: &id,
					ID:         &networkID,
					Name:       &newName,
				}
				_, response, err := instanceC.UpdateNetworkInterface(updatepnicfoptions)
				if err != nil {
					return fmt.Errorf("Error while updating name %s for network interface of instance %s\n%s: %q", newName, d.Id(), err, response)
				}
				if err != nil {
					return err
				}
			}
		}

	}

	if d.HasChange(isInstanceName) {
		name := d.Get(isInstanceName).(string)
		updnetoptions := &vpcv1.UpdateInstanceOptions{
			ID:   &id,
			Name: &name,
		}
		_, _, err = instanceC.UpdateInstance(updnetoptions)
		if err != nil {
			return err
		}
	}

	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
	}
	if d.HasChange(isInstanceTags) {
		oldList, newList := d.GetChange(isInstanceTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on update of resource vpc Instance (%s) tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicInstanceUpdate(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := instanceUpdate(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceIBMisInstanceRead(d, meta)
}

func classicInstanceDelete(d *schema.ResourceData, meta interface{}, id string) error {
	instanceC, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	getinsOptions := &vpcclassicv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Instance (%s): %s\n%s", id, err, response)
	}
	actiontype := "stop"
	createinsactoptions := &vpcclassicv1.CreateInstanceActionOptions{
		InstanceID: &id,
		Type:       &actiontype,
	}
	_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Creating Instance Action: %s\n%s", err, response)
	}
	_, err = isWaitForClassicInstanceActionStop(instanceC, d, meta, id)
	if err != nil {
		return err
	}
	listvolattoptions := &vpcclassicv1.ListVolumeAttachmentsOptions{
		InstanceID: &id,
	}
	vols, response, err := instanceC.ListVolumeAttachments(listvolattoptions)
	if err != nil {
		return fmt.Errorf("Error Listing volume attachments to the instance: %s\n%s", err, response)
	}
	bootvolid := ""
	for _, vol := range vols.VolumeAttachments {
		if *vol.Type == "data" {
			delvolattoptions := &vpcclassicv1.DeleteVolumeAttachmentOptions{
				InstanceID: &id,
				ID:         vol.ID,
			}
			_, err := instanceC.DeleteVolumeAttachment(delvolattoptions)
			if err != nil {
				return fmt.Errorf("Error while removing volume attachment %q for instance %s: %q", *vol.ID, d.Id(), err)
			}
			_, err = isWaitForClassicInstanceVolumeDetached(instanceC, d, d.Id(), *vol.ID)
			if err != nil {
				return err
			}
			break
		}
		if *vol.Type == "boot" {
			bootvolid = *vol.Volume.ID
		}
	}
	deleteinstanceOptions := &vpcclassicv1.DeleteInstanceOptions{
		ID: &id,
	}
	_, err = instanceC.DeleteInstance(deleteinstanceOptions)
	if err != nil {
		return err
	}
	_, err = isWaitForClassicInstanceDelete(instanceC, d, d.Id())
	if err != nil {
		return err
	}
	if _, ok := d.GetOk(isInstanceBootVolume); ok {
		_, err = isWaitForClassicVolumeDeleted(instanceC, bootvolid, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return err
		}
	}
	return nil
}

func instanceDelete(d *schema.ResourceData, meta interface{}, id string) error {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Instance (%s): %s\n%s", id, err, response)
	}
	actiontype := "stop"
	createinsactoptions := &vpcv1.CreateInstanceActionOptions{
		InstanceID: &id,
		Type:       &actiontype,
	}
	_, response, err = instanceC.CreateInstanceAction(createinsactoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error Creating Instance Action: %s\n%s", err, response)
	}
	_, err = isWaitForInstanceActionStop(instanceC, d, meta, id)
	if err != nil {
		return err
	}
	listvolattoptions := &vpcv1.ListVolumeAttachmentsOptions{
		InstanceID: &id,
	}
	vols, response, err := instanceC.ListVolumeAttachments(listvolattoptions)
	if err != nil {
		return fmt.Errorf("Error Listing volume attachments to the instance: %s\n%s", err, response)
	}
	bootvolid := ""
	for _, vol := range vols.VolumeAttachments {
		if *vol.Type == "data" {
			delvolattoptions := &vpcv1.DeleteVolumeAttachmentOptions{
				InstanceID: &id,
				ID:         vol.ID,
			}
			_, err := instanceC.DeleteVolumeAttachment(delvolattoptions)
			if err != nil {
				return fmt.Errorf("Error while removing volume Attachment %q for instance %s: %q", *vol.ID, d.Id(), err)
			}
			_, err = isWaitForInstanceVolumeDetached(instanceC, d, d.Id(), *vol.ID)
			if err != nil {
				return err
			}
			break
		}
		if *vol.Type == "boot" {
			bootvolid = *vol.Volume.ID
		}
	}
	deleteinstanceOptions := &vpcv1.DeleteInstanceOptions{
		ID: &id,
	}
	_, err = instanceC.DeleteInstance(deleteinstanceOptions)
	if err != nil {
		return err
	}
	_, err = isWaitForInstanceDelete(instanceC, d, d.Id())
	if err != nil {
		return err
	}
	if _, ok := d.GetOk(isInstanceBootVolume); ok {
		_, err = isWaitForVolumeDeleted(instanceC, bootvolid, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceIBMisInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		err := classicInstanceDelete(d, meta, id)
		if err != nil {
			return err
		}
	} else {
		err := instanceDelete(d, meta, id)
		if err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}

func classicInstanceExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	instanceC, err := classicVpcClient(meta)
	if err != nil {
		return false, err
	}
	getinsOptions := &vpcclassicv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
	}
	return true, nil
}

func instanceExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	instanceC, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getinsOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := instanceC.GetInstance(getinsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
	}
	return true, nil
}

func resourceIBMisInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return false, err
	}
	id := d.Id()
	if userDetails.generation == 1 {
		exists, err := classicInstanceExists(d, meta, id)
		return exists, err
	} else {
		exists, err := instanceExists(d, meta, id)
		return exists, err
	}
}

func isWaitForClassicInstanceDelete(instanceC *vpcclassicv1.VpcClassicV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceDeleting, isInstanceAvailable},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcclassicv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return instance, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
			}
			if *instance.Status == isInstanceFailed {
				return instance, *instance.Status, fmt.Errorf("The  instance %s failed to delete: %v", d.Id(), err)
			}
			return instance, isInstanceDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForInstanceDelete(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceDeleting, isInstanceAvailable},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return instance, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
			}
			if *instance.Status == isInstanceFailed {
				return instance, *instance.Status, fmt.Errorf("The  instance %s failed to delete: %v", d.Id(), err)
			}
			return instance, isInstanceDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isWaitForClassicInstanceActionStop(instanceC *vpcclassicv1.VpcClassicV1, d *schema.ResourceData, meta interface{}, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceStatusRunning, isInstanceStatusPending, isInstanceActionStatusStopping},
		Target:  []string{isInstanceActionStatusStopped, isInstanceStatusFailed, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcclassicv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				return nil, "", fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
			}
			if *instance.Status == isInstanceStatusFailed {
				return instance, *instance.Status, fmt.Errorf("The  instance %s failed to stop: %v", d.Id(), err)
			}
			return instance, *instance.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func isWaitForInstanceActionStop(instanceC *vpcv1.VpcV1, d *schema.ResourceData, meta interface{}, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceStatusRunning, isInstanceStatusPending, isInstanceActionStatusStopping},
		Target:  []string{isInstanceActionStatusStopped, isInstanceStatusFailed, ""},
		Refresh: func() (interface{}, string, error) {
			getinsoptions := &vpcv1.GetInstanceOptions{
				ID: &id,
			}
			instance, response, err := instanceC.GetInstance(getinsoptions)
			if err != nil {
				return nil, "", fmt.Errorf("Error Getting Instance: %s\n%s", err, response)
			}
			if *instance.Status == isInstanceStatusFailed {
				return instance, *instance.Status, fmt.Errorf("The  instance %s failed to stop: %v", d.Id(), err)
			}
			return instance, *instance.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForClassicInstanceVolumeAttached(instanceC *vpcclassicv1.VpcClassicV1, d *schema.ResourceData, id, volID string) (interface{}, error) {
	log.Printf("Waiting for instance volume (%s) to be attched.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceVolumeAttaching},
		Target:     []string{isInstanceVolumeAttached, ""},
		Refresh:    isClassicInstanceVolumeRefreshFunc(instanceC, id, volID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isClassicInstanceVolumeRefreshFunc(instanceC *vpcclassicv1.VpcClassicV1, id, volID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getvolattoptions := &vpcclassicv1.GetVolumeAttachmentOptions{
			InstanceID: &id,
			ID:         &volID,
		}
		vol, response, err := instanceC.GetVolumeAttachment(getvolattoptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Attaching volume: %s\n%s", err, response)
		}

		if *vol.Status == isInstanceVolumeAttached {
			return vol, isInstanceVolumeAttached, nil
		}

		return vol, isInstanceVolumeAttaching, nil
	}
}

func isWaitForInstanceVolumeAttached(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id, volID string) (interface{}, error) {
	log.Printf("Waiting for instance volume (%s) to be attched.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceVolumeAttaching},
		Target:     []string{isInstanceVolumeAttached, ""},
		Refresh:    isInstanceVolumeRefreshFunc(instanceC, id, volID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isInstanceVolumeRefreshFunc(instanceC *vpcv1.VpcV1, id, volID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getvolattoptions := &vpcv1.GetVolumeAttachmentOptions{
			InstanceID: &id,
			ID:         &volID,
		}
		vol, response, err := instanceC.GetVolumeAttachment(getvolattoptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Attaching volume: %s\n%s", err, response)
		}

		if *vol.Status == isInstanceVolumeAttached {
			return vol, isInstanceVolumeAttached, nil
		}

		return vol, isInstanceVolumeAttaching, nil
	}
}

func isWaitForClassicInstanceVolumeDetached(instanceC *vpcclassicv1.VpcClassicV1, d *schema.ResourceData, id, volID string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceVolumeAttached, isInstanceVolumeDetaching},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getvolattoptions := &vpcclassicv1.GetVolumeAttachmentOptions{
				InstanceID: &id,
				ID:         &volID,
			}
			vol, response, err := instanceC.GetVolumeAttachment(getvolattoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return vol, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error Detaching volume: %s\n%s", err, response)
			}
			if *vol.Status == isInstanceFailed {
				return vol, *vol.Status, fmt.Errorf("The instance %s failed to detach volume %s: %v", d.Id(), volID, err)
			}
			return vol, isInstanceVolumeDetaching, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForInstanceVolumeDetached(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id, volID string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isInstanceVolumeAttached, isInstanceVolumeDetaching},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getvolattoptions := &vpcv1.GetVolumeAttachmentOptions{
				InstanceID: &id,
				ID:         &volID,
			}
			vol, response, err := instanceC.GetVolumeAttachment(getvolattoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return vol, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error Detaching: %s\n%s", err, response)
			}
			if *vol.Status == isInstanceFailed {
				return vol, *vol.Status, fmt.Errorf("The instance %s failed to detach volume %s: %v", d.Id(), volID, err)
			}
			return vol, isInstanceVolumeDetaching, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
