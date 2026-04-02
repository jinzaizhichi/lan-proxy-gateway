package rules

import (
	"slices"
	"testing"
)

func TestNormalProbeRulesUseOrdinaryProxyPath(t *testing.T) {
	rules := NormalProbeRules("Proxy")

	for _, want := range []string{
		"DOMAIN,ipwho.is,Proxy",
		"DOMAIN,api.ip.sb,Proxy",
	} {
		if !slices.Contains(rules, want) {
			t.Fatalf("NormalProbeRules() missing %q", want)
		}
	}
}

func TestAIProxyRulesKeepResidentialProbe(t *testing.T) {
	rules := AIProxyRules("AI Only")

	if !slices.Contains(rules, "DOMAIN,checkip.amazonaws.com,AI Only") {
		t.Fatal("AIProxyRules() should keep checkip.amazonaws.com on AI Only")
	}
}
