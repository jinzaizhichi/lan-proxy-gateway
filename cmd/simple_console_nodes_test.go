package cmd

import "testing"

func TestParseSubscriptionHintItemRecognizesCommonHints(t *testing.T) {
	tests := []struct {
		input     string
		wantKind  string
		wantValue string
	}{
		{input: "剩余流量：41.45 GB", wantKind: "remaining", wantValue: "41.45 GB"},
		{input: "套餐到期: 2027-03-28", wantKind: "expiry", wantValue: "2027-03-28"},
		{input: "最新网址：fliggycloud.cc", wantKind: "website", wantValue: "fliggycloud.cc"},
		{input: "距离下次重置剩余：29 天", wantKind: "reset", wantValue: "29 天"},
	}

	for _, tc := range tests {
		kind, value, ok := parseSubscriptionHintItem(tc.input)
		if !ok {
			t.Fatalf("parseSubscriptionHintItem(%q) should match", tc.input)
		}
		if kind != tc.wantKind || value != tc.wantValue {
			t.Fatalf("parseSubscriptionHintItem(%q) = (%q, %q), want (%q, %q)", tc.input, kind, value, tc.wantKind, tc.wantValue)
		}
	}

	if _, _, ok := parseSubscriptionHintItem("香港 02"); ok {
		t.Fatal("expected regular node name to stay selectable")
	}
}

func TestSimpleSelectableNodesFiltersSubscriptionHints(t *testing.T) {
	nodes := []string{
		"最新网址：fliggycloud.cc",
		"香港 01",
		"剩余流量：41.45 GB",
		"香港 02",
		"香港 02",
		"距离下次重置剩余：29 天",
	}

	got := simpleSelectableNodes(nodes)
	if len(got) != 2 {
		t.Fatalf("len(simpleSelectableNodes()) = %d, want 2 (%v)", len(got), got)
	}
	if got[0] != "香港 01" || got[1] != "香港 02" {
		t.Fatalf("simpleSelectableNodes() = %v, want only real nodes", got)
	}
}

func TestSortSimpleNodeDelayEntriesOrdersByDelay(t *testing.T) {
	entries := []simpleNodeDelayEntry{
		{Name: "香港 03", Delay: -1, DelayLabel: "失败", Reachable: false},
		{Name: "香港 02", Delay: 86, DelayLabel: "86ms", Reachable: true},
		{Name: "香港 01", Delay: 42, DelayLabel: "42ms", Reachable: true, Current: true},
		{Name: "香港 04", Delay: -1, DelayLabel: "超时", Reachable: false, Current: true},
	}

	sortSimpleNodeDelayEntries(entries)

	want := []string{"香港 01", "香港 02", "香港 04", "香港 03"}
	for i, name := range want {
		if entries[i].Name != name {
			t.Fatalf("sorted entry %d = %q, want %q", i, entries[i].Name, name)
		}
	}
}
