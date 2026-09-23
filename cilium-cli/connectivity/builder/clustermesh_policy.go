// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package builder

import (
	"github.com/cilium/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium/cilium-cli/connectivity/tests"
)

type clusterMeshPolicy struct{}

func (clusterMeshPolicy) build(ct *check.ConnectivityTest, _ map[string]string) {
	multiCluster := func() bool { return ct.Params().MultiCluster != "" }

	newTest("clustermesh-policy-cnp-ingress", ct).
		WithCondition(multiCluster).
		WithCiliumPolicy(echoIngressFromOtherClientPolicyYAML).
		WithScenarios(tests.PodToPod(tests.WithCrossClusterOnly())).
		WithExpectations(func(a *check.Action) (egress, ingress check.Result) {
			if a.Destination().HasLabel("kind", "echo") && !a.Source().HasLabel("other", "client") {
				return check.ResultDropCurlTimeout, check.ResultDropCurlTimeout
			}
			return check.ResultOK, check.ResultOK
		})

	newTest("clustermesh-policy-cnp-egress", ct).
		WithCondition(multiCluster).
		WithCiliumPolicy(clientEgressToEchoPolicyYAML).
		WithScenarios(tests.PodToPod(tests.WithCrossClusterOnly()))
}
