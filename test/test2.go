package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/azure"
)

func TestAzureLinuxVMCreation(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../", 
    }

    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    // Get outputs from Terraform
    nicName := terraform.Output(t, terraformOptions, "nic_name")
    vmName := terraform.Output(t, terraformOptions, "vm_name")
    resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

    // Test 1: Check if NIC exists and is connected
    nic := azure.GetNetworkInterface(t, nicName, resourceGroupName, "")
    assert.Equal(t, vmName, *nic.VirtualMachine.ID, "NIC is not connected to the correct VM")

    // Test 2: Verify Ubuntu version on the VM
    expectedVersion := "22.04"  // Change based on your project
    actualVersion := azure.GetVirtualMachine(t, vmName, resourceGroupName, "").StorageProfile.ImageReference.Sku
    assert.Contains(t, actualVersion, expectedVersion, "VM is not running the expected Ubuntu version")
}
