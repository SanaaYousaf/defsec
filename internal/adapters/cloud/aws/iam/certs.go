package iam

import (
	"fmt"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/concurrency"
	"github.com/aquasecurity/defsec/pkg/providers/aws/iam"
	"github.com/aquasecurity/defsec/pkg/state"
	iamapi "github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (a *adapter) adaptServerCertificates(state *state.State) error {
	a.Tracker().SetServiceLabel("Discovering server certificates...")

	var certs []iamtypes.ServerCertificateMetadata

	input := &iamapi.ListServerCertificatesInput{}
	for {
		certsOutput, err := a.api.ListServerCertificates(a.Context(), input)
		if err != nil {
			return err
		}
		certs = append(certs, certsOutput.ServerCertificateMetadataList...)
		a.Tracker().SetTotalResources(len(certs))
		if !certsOutput.IsTruncated {
			break
		}
		input.Marker = certsOutput.Marker
	}

	a.Tracker().SetServiceLabel("Adapting server certificates...")

	state.AWS.IAM.ServerCertificates = concurrency.Adapt(certs, a.RootAdapter, a.adaptServerCertificate)
	return nil
}

func (a *adapter) adaptServerCertificate(certInfo iamtypes.ServerCertificateMetadata) (*iam.ServerCertificate, error) {
	cert, err := a.api.GetServerCertificate(a.Context(), &iamapi.GetServerCertificateInput{
		ServerCertificateName: certInfo.ServerCertificateName,
	})
	if err != nil {
		return nil, err
	}

	if cert.ServerCertificate.ServerCertificateMetadata == nil || cert.ServerCertificate.ServerCertificateMetadata.Arn == nil {
		return nil, fmt.Errorf("server certificate metadata is nil")
	}

	metadata := a.CreateMetadataFromARN(*cert.ServerCertificate.ServerCertificateMetadata.Arn)

	expiration := defsecTypes.TimeUnresolvable(metadata)
	if cert.ServerCertificate.ServerCertificateMetadata.Expiration != nil {
		expiration = defsecTypes.Time(*cert.ServerCertificate.ServerCertificateMetadata.Expiration, metadata)
	}

	return &iam.ServerCertificate{
		Metadata:   metadata,
		Name:       defsecTypes.String(*certInfo.ServerCertificateName, metadata),
		Expiration: expiration,
	}, nil
}

func (a *adapter) adaptVirtualMfadevices(state *state.State) error {
	a.Tracker().SetServiceLabel("Discovering Virtual Mfa Devices...")

	var mfadevices []iamtypes.VirtualMFADevice

	input := &iamapi.ListVirtualMFADevicesInput{}
	for {
		certsOutput, err := a.api.ListVirtualMFADevices(a.Context(), input)
		if err != nil {
			return err
		}
		mfadevices = append(mfadevices, certsOutput.VirtualMFADevices...)
		a.Tracker().SetTotalResources(len(mfadevices))
		if !certsOutput.IsTruncated {
			break
		}
		input.Marker = certsOutput.Marker
	}

	a.Tracker().SetServiceLabel("Adapting virtual mfa devices...")

	state.AWS.IAM.VirtualMfaDevices = concurrency.Adapt(mfadevices, a.RootAdapter, a.adaptVirtualMfaDevice)
	return nil
}

func (a *adapter) adaptVirtualMfaDevice(mfaapi iamtypes.VirtualMFADevice) (*iam.VirtualMfaDevice, error) {
	metadata := a.CreateMetadataFromARN(*mfaapi.User.Arn)

	return &iam.VirtualMfaDevice{
		Metadata:     metadata,
		SerialNumber: defsecTypes.String(*mfaapi.SerialNumber, metadata),
	}, nil
}
