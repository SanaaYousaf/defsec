package efs

import (
	"github.com/aquasecurity/defsec/pkg/providers/aws/efs"
	"github.com/aquasecurity/defsec/pkg/terraform"
)

func Adapt(modules terraform.Modules) efs.EFS {
	return efs.EFS{
		FileSystems: adaptFileSystems(modules),
	}
}

func adaptFileSystems(modules terraform.Modules) []efs.FileSystem {
	var filesystems []efs.FileSystem
	for _, module := range modules {
		for _, resource := range module.GetResourcesByType("aws_efs_file_system") {
			filesystems = append(filesystems, adaptFileSystem(resource))
		}
	}
	return filesystems
}

func adaptFileSystem(resource *terraform.Block) efs.FileSystem {
	encryptedAttr := resource.GetAttribute("encrypted")
	encryptedVal := encryptedAttr.AsBoolValueOrDefault(false, resource)

	kmskeyidAttr := resource.GetAttribute("kms_key_id")
	kmskeyidVal := kmskeyidAttr.AsStringValueOrDefault("", resource)

	return efs.FileSystem{
		Metadata:  resource.GetMetadata(),
		Encrypted: encryptedVal,
		KmsKeyId:  kmskeyidVal,
	}
}
