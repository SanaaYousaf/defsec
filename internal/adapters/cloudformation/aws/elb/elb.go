package elb

import (
	"github.com/aquasecurity/defsec/pkg/providers/aws/elb"
	"github.com/aquasecurity/defsec/pkg/scanners/cloudformation/parser"
)

// Adapt ...
func Adapt(cfFile parser.FileContext) elb.ELB {
	return elb.ELB{
		LoadBalancers:        getLoadBalancers(cfFile),
		TargetGroups:         getTargetGroups(cfFile),
		LoadBalancersV1:      getLoadBalancersV1(cfFile),
		LoadBalancerPolicies: getLoadBalancersPolicy(cfFile),
	}
}
