package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {

	type input struct {
		fixture string
	}

	type expect struct {
		dataInit        *Data
		dataVMs         *Data
		tagsInit        *ProtectedTagMap
		tagsVMs         *ProtectedTagMap
		srcTagSetsInit  *ProtectedTagSetMap
		srcTagSetsFinal *ProtectedTagSetMap
		responses       []string
	}

	expect_0 := func() expect {
		// Initialize tags
		dataInit := &Data{VMs: []*VM{{"vm-a211de", "jira_server", []string{"ci", "dev"}, nil},
			{"vm-c7bac01a07", "bastion", []string{"ssh", "dev"}, nil},
		}, FWRules: []*FWRule{{"fw-82af742", "ssh", "dev"}}}

		devTagInit := &Tag{Name: "dev", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		sshTagInit := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		devTagInit.from.m["ssh"] = sshTagInit
		sshTagInit.to.m["dev"] = devTagInit
		tagsInit := &ProtectedTagMap{m: map[string]*Tag{
			"dev": devTagInit,
			"ssh": sshTagInit,
		}}

		// tags state after VMs collection
		dataVMs := &Data{VMs: []*VM{{"vm-a211de", "jira_server", []string{"ci", "dev"}, nil},
			{"vm-c7bac01a07", "bastion", []string{"ssh", "dev"}, nil},
		}, FWRules: []*FWRule{{"fw-82af742", "ssh", "dev"}}}

		devTagVMs := &Tag{Name: "dev", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: map[string]*VM{"vm-a211de": dataVMs.VMs[0], "vm-c7bac01a07": dataVMs.VMs[1]}}
		sshTagVMs := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: map[string]*VM{"vm-c7bac01a07": dataVMs.VMs[1]}}
		devTagVMs.from.m["ssh"] = sshTagVMs
		sshTagVMs.to.m["dev"] = devTagVMs
		tagsVMs := &ProtectedTagMap{m: map[string]*Tag{
			"dev": devTagVMs,
			"ssh": sshTagVMs,
		}}

		dstTagSetsInit := &ProtectedTagSetMap{m: map[string]*ProtectedTagMap{
			"dev": {m: map[string]*Tag{"dev": devTagVMs}}}}
		srcTagSetsInit := &ProtectedTagSetMap{m: map[string]*ProtectedTagMap{
			"ssh": {m: map[string]*Tag{"ssh": sshTagVMs}}}}
		dstTagSetsInit.m["dev"].Link = srcTagSetsInit.m["ssh"]
		dataVMs.VMs[0].DstTagMap = dstTagSetsInit.m["dev"]
		dataVMs.VMs[1].DstTagMap = dstTagSetsInit.m["dev"]

		// after processing the tag sets
		dataFinal := &Data{VMs: []*VM{{"vm-a211de", "jira_server", []string{"ci", "dev"}, nil},
			{"vm-c7bac01a07", "bastion", []string{"ssh", "dev"}, nil},
		}, FWRules: []*FWRule{{"fw-82af742", "ssh", "dev"}}}

		devTagFinal := &Tag{Name: "dev", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: map[string]*VM{"vm-a211de": dataFinal.VMs[0], "vm-c7bac01a07": dataFinal.VMs[1]}}
		sshTagFinal := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: map[string]*VM{"vm-c7bac01a07": dataFinal.VMs[1]}}
		devTagFinal.from.m["ssh"] = sshTagFinal
		sshTagFinal.to.m["dev"] = devTagFinal

		slice := &Slice{StringSlice: []string{"vm-c7bac01a07"},
			pos:   map[string][2]int{"vm-c7bac01a07": {1, 16}},
			bytes: []byte(`["vm-c7bac01a07"]`)}
		dstTagSetsFinal := &ProtectedTagSetMap{m: map[string]*ProtectedTagMap{
			"dev": {m: map[string]*Tag{"dev": devTagFinal}}}}
		srcTagSetsFinal := &ProtectedTagSetMap{m: map[string]*ProtectedTagMap{
			"ssh": {m: map[string]*Tag{"ssh": sshTagFinal}, Slice: slice}}}
		dstTagSetsFinal.m["dev"].Link = srcTagSetsFinal.m["ssh"]
		dataFinal.VMs[0].DstTagMap = dstTagSetsFinal.m["dev"]
		dataFinal.VMs[1].DstTagMap = dstTagSetsFinal.m["dev"]

		responses := []string{`["vm-c7bac01a07"]`, `[]`}

		exp := expect{dataInit, dataVMs, tagsInit, tagsVMs, srcTagSetsInit, srcTagSetsFinal, responses}
		return exp
	}

	expect_1 := func() expect {
		// Initialize tags
		dataInit := &Data{VMs: []*VM{{"vm-b8e6c350", "rabbitmq", []string{"windows-dc"}, nil},
			{"vm-c1e6285f", "k8s node", []string{"http", "ci"}, nil},
			{"vm-cf1f8621", "k8s node", []string{"windows-dc"}, nil},
			{"vm-b462c04", "jira server", []string{"windows-dc", "storage"}, nil},
			{"vm-8d2d12765", "kafka", []string{}, nil},
			{"vm-9cbedf7c66", "etcd node", []string{}, nil},
			{"vm-ae24e37f8a", "frontend server", []string{"api", "dev"}, nil},
			{"vm-e30d5fa49a", "etcd node", []string{"dev", "api"}, nil},
			{"vm-1b1cc9cd", "billing service", []string{}, nil},
			{"vm-f270036588", "kafka", []string{}, nil},
		}, FWRules: []*FWRule{{"fw-c4a11ac", "k8s", "loadbalancer"},
			{"fw-1cb4c1", "django", "django"},
			{"fw-0d91970ee", "corp", "django"},
			{"fw-778beb64", "django", "nat"},
			{"fw-1008d7", "ssh", "ssh"},
			{"fw-1c8ebac1f", "loadbalancer", "corp"},
			{"fw-06bf6a628", "nat", "nat"},
			{"fw-9d030bb4bb", "corp", "nat"},
			{"fw-fbcfcc16e1", "antivirus", "nat"},
			{"fw-c74204", "antivirus", "antivirus"},
		}}

		k8sTagInit := &Tag{Name: "k8s", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		loadbalancerTagInit := &Tag{Name: "loadbalancer", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		djangoTagInit := &Tag{Name: "django", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		corpTagInit := &Tag{Name: "corp", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		natTagInit := &Tag{Name: "nat", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		sshTagInit := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		antivirusTagInit := &Tag{Name: "antivirus", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}

		loadbalancerTagInit.from.m["k8s"] = k8sTagInit
		k8sTagInit.to.m["loadbalancer"] = loadbalancerTagInit
		djangoTagInit.from.m["django"] = djangoTagInit
		djangoTagInit.to.m["django"] = djangoTagInit
		djangoTagInit.from.m["corp"] = corpTagInit
		corpTagInit.to.m["django"] = djangoTagInit
		natTagInit.from.m["django"] = djangoTagInit
		djangoTagInit.to.m["nat"] = natTagInit
		sshTagInit.from.m["ssh"] = sshTagInit
		sshTagInit.to.m["ssh"] = sshTagInit
		corpTagInit.from.m["loadbalancer"] = loadbalancerTagInit
		loadbalancerTagInit.to.m["corp"] = corpTagInit
		natTagInit.from.m["nat"] = natTagInit
		natTagInit.to.m["nat"] = natTagInit
		natTagInit.from.m["corp"] = corpTagInit
		corpTagInit.to.m["nat"] = natTagInit
		natTagInit.from.m["antivirus"] = antivirusTagInit
		antivirusTagInit.to.m["nat"] = natTagInit
		antivirusTagInit.from.m["antivirus"] = antivirusTagInit
		antivirusTagInit.to.m["antivirus"] = antivirusTagInit

		tagsInit := &ProtectedTagMap{m: map[string]*Tag{
			"k8s":          k8sTagInit,
			"loadbalancer": loadbalancerTagInit,
			"django":       djangoTagInit,
			"corp":         corpTagInit,
			"nat":          natTagInit,
			"ssh":          sshTagInit,
			"antivirus":    antivirusTagInit,
		}}

		// tags state after VMs collection
		dataVMs := &Data{VMs: []*VM{{"vm-b8e6c350", "rabbitmq", []string{"windows-dc"}, nil},
			{"vm-c1e6285f", "k8s node", []string{"http", "ci"}, nil},
			{"vm-cf1f8621", "k8s node", []string{"windows-dc"}, nil},
			{"vm-b462c04", "jira server", []string{"windows-dc", "storage"}, nil},
			{"vm-8d2d12765", "kafka", []string{}, nil},
			{"vm-9cbedf7c66", "etcd node", []string{}, nil},
			{"vm-ae24e37f8a", "frontend server", []string{"api", "dev"}, nil},
			{"vm-e30d5fa49a", "etcd node", []string{"dev", "api"}, nil},
			{"vm-1b1cc9cd", "billing service", []string{}, nil},
			{"vm-f270036588", "kafka", []string{}, nil},
		}, FWRules: []*FWRule{{"fw-c4a11ac", "k8s", "loadbalancer"},
			{"fw-1cb4c1", "django", "django"},
			{"fw-0d91970ee", "corp", "django"},
			{"fw-778beb64", "django", "nat"},
			{"fw-1008d7", "ssh", "ssh"},
			{"fw-1c8ebac1f", "loadbalancer", "corp"},
			{"fw-06bf6a628", "nat", "nat"},
			{"fw-9d030bb4bb", "corp", "nat"},
			{"fw-fbcfcc16e1", "antivirus", "nat"},
			{"fw-c74204", "antivirus", "antivirus"},
		}}

		k8sTagVMs := &Tag{Name: "k8s", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		loadbalancerTagVMs := &Tag{Name: "loadbalancer", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		djangoTagVMs := &Tag{Name: "django", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		corpTagVMs := &Tag{Name: "corp", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		natTagVMs := &Tag{Name: "nat", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		sshTagVMs := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		antivirusTagVMs := &Tag{Name: "antivirus", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}

		loadbalancerTagVMs.from.m["k8s"] = k8sTagVMs
		k8sTagVMs.to.m["loadbalancer"] = loadbalancerTagVMs
		djangoTagVMs.from.m["django"] = djangoTagVMs
		djangoTagVMs.to.m["django"] = djangoTagVMs
		djangoTagVMs.from.m["corp"] = corpTagVMs
		corpTagVMs.to.m["django"] = djangoTagVMs
		natTagVMs.from.m["django"] = djangoTagVMs
		djangoTagVMs.to.m["nat"] = natTagVMs
		sshTagVMs.from.m["ssh"] = sshTagVMs
		sshTagVMs.to.m["ssh"] = sshTagVMs
		corpTagVMs.from.m["loadbalancer"] = loadbalancerTagVMs
		loadbalancerTagVMs.to.m["corp"] = corpTagVMs
		natTagVMs.from.m["nat"] = natTagVMs
		natTagVMs.to.m["nat"] = natTagVMs
		natTagVMs.from.m["corp"] = corpTagVMs
		corpTagVMs.to.m["nat"] = natTagVMs
		natTagVMs.from.m["antivirus"] = antivirusTagVMs
		antivirusTagVMs.to.m["nat"] = natTagVMs
		antivirusTagVMs.from.m["antivirus"] = antivirusTagVMs
		antivirusTagVMs.to.m["antivirus"] = antivirusTagVMs

		tagsVMs := &ProtectedTagMap{m: map[string]*Tag{
			"k8s":          k8sTagVMs,
			"loadbalancer": loadbalancerTagVMs,
			"django":       djangoTagVMs,
			"corp":         corpTagVMs,
			"nat":          natTagVMs,
			"ssh":          sshTagVMs,
			"antivirus":    antivirusTagVMs,
		}}

		srcTagSetsInit := &ProtectedTagSetMap{m: map[string]*ProtectedTagMap{}}

		// after processing the tag sets
		k8sTagFinal := &Tag{Name: "k8s", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		loadbalancerTagFinal := &Tag{Name: "loadbalancer", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		djangoTagFinal := &Tag{Name: "django", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		corpTagFinal := &Tag{Name: "corp", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		natTagFinal := &Tag{Name: "nat", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		sshTagFinal := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		antivirusTagFinal := &Tag{Name: "antivirus", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}

		loadbalancerTagFinal.from.m["k8s"] = k8sTagFinal
		k8sTagFinal.to.m["loadbalancer"] = loadbalancerTagFinal
		djangoTagFinal.from.m["django"] = djangoTagFinal
		djangoTagFinal.to.m["django"] = djangoTagFinal
		djangoTagFinal.from.m["corp"] = corpTagFinal
		corpTagFinal.to.m["django"] = djangoTagFinal
		natTagFinal.from.m["django"] = djangoTagFinal
		djangoTagFinal.to.m["nat"] = natTagFinal
		sshTagFinal.from.m["ssh"] = sshTagFinal
		sshTagFinal.to.m["ssh"] = sshTagFinal
		corpTagFinal.from.m["loadbalancer"] = loadbalancerTagFinal
		loadbalancerTagFinal.to.m["corp"] = corpTagFinal
		natTagFinal.from.m["nat"] = natTagFinal
		natTagFinal.to.m["nat"] = natTagFinal
		natTagFinal.from.m["corp"] = corpTagFinal
		corpTagFinal.to.m["nat"] = natTagFinal
		natTagFinal.from.m["antivirus"] = antivirusTagFinal
		antivirusTagFinal.to.m["nat"] = natTagFinal
		antivirusTagFinal.from.m["antivirus"] = antivirusTagFinal
		antivirusTagFinal.to.m["antivirus"] = antivirusTagFinal

		srcTagSetsFinal := &ProtectedTagSetMap{m: map[string]*ProtectedTagMap{}}
		responses := []string{`[]`, `[]`, `[]`, `[]`, `[]`, `[]`, `[]`, `[]`, `[]`, `[]`}

		exp := expect{dataInit, dataVMs, tagsInit, tagsVMs, srcTagSetsInit, srcTagSetsFinal, responses}
		return exp
	}

	expect_2 := func() expect {
		// Initialize tags
		dataInit := &Data{VMs: []*VM{{"vm-ec02d5c153", "kafka", []string{"http"}, nil},
			{"vm-a3ed2eed23", "rabbitmq", []string{"https", "http"}, nil},
			{"vm-2ba4d2f87", "ssh bastion", []string{"http", "windows-dc", "nat", "https", "storage"}, nil},
			{"vm-b35b501", "dev-srv-5", []string{"ssh", "nat", "http", "loadbalancer", "storage"}, nil},
			{"vm-7d1ff7af47", "billing service", []string{"http"}, nil},
		}, FWRules: []*FWRule{{"fw-c8706961d", "loadbalancer", "windows-dc"},
			{"fw-76f36a3", "ssh", "ci"},
			{"fw-487b076a6", "storage", "reverse_proxy"},
			{"fw-dd16d0", "nat", "ssh"},
			{"fw-36719127", "https", "loadbalancer"},
			{"fw-1f8b1e8d8", "loadbalancer", "storage"},
			{"fw-e602b7a05", "nat", "nat"},
			{"fw-4e337463", "reverse_proxy", "storage"},
			{"fw-a646f8da6", "http", "http"},
			{"fw-28c3124", "ssh", "https"},
			{"fw-1310da", "ssh", "nat"},
			{"fw-64ae2f2be7", "corp", "nat"},
			{"fw-488809fc3", "corp", "windows-dc"},
			{"fw-4878f98212", "ssh", "ssh"},
			{"fw-1a0642c", "nat", "corp"},
			{"fw-e6b9108", "windows-dc", "corp"},
		}}

		loadbalancerTagInit := &Tag{Name: "loadbalancer", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		windowsdcTagInit := &Tag{Name: "windows-dc", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		sshTagInit := &Tag{Name: "ssh", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		ciTagInit := &Tag{Name: "ci", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		storageTagInit := &Tag{Name: "storage", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		reverseProxyTagInit := &Tag{Name: "reverse_proxy", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		natTagInit := &Tag{Name: "nat", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		httpsTagInit := &Tag{Name: "https", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		httpTagInit := &Tag{Name: "http", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}
		corpTagInit := &Tag{Name: "corp", from: NewProtectedTagMap(), to: NewProtectedTagMap(),
			vms: make(map[string]*VM)}

		windowsdcTagInit.from.m["loadbalancer"] = loadbalancerTagInit
		loadbalancerTagInit.to.m["windows-dc"] = windowsdcTagInit
		ciTagInit.from.m["ssh"] = sshTagInit
		sshTagInit.to.m["ci"] = ciTagInit
		reverseProxyTagInit.from.m["storage"] = storageTagInit
		storageTagInit.to.m["reverse_proxy"] = reverseProxyTagInit
		sshTagInit.from.m["nat"] = natTagInit
		natTagInit.to.m["ssh"] = sshTagInit
		loadbalancerTagInit.from.m["https"] = httpsTagInit
		httpsTagInit.to.m["loadbalancer"] = loadbalancerTagInit
		storageTagInit.from.m["loadbalancer"] = loadbalancerTagInit
		loadbalancerTagInit.to.m["storage"] = storageTagInit
		natTagInit.from.m["nat"] = natTagInit
		natTagInit.to.m["nat"] = natTagInit
		storageTagInit.from.m["reverse_proxy"] = reverseProxyTagInit
		reverseProxyTagInit.to.m["storage"] = storageTagInit
		httpTagInit.from.m["http"] = httpTagInit
		httpTagInit.to.m["http"] = httpTagInit
		httpsTagInit.from.m["ssh"] = sshTagInit
		sshTagInit.to.m["https"] = httpsTagInit
		natTagInit.from.m["ssh"] = sshTagInit
		sshTagInit.to.m["nat"] = natTagInit
		natTagInit.from.m["corp"] = corpTagInit
		corpTagInit.to.m["nat"] = natTagInit
		windowsdcTagInit.from.m["corp"] = corpTagInit
		corpTagInit.to.m["windows-dc"] = windowsdcTagInit
		sshTagInit.from.m["ssh"] = sshTagInit
		sshTagInit.to.m["ssh"] = sshTagInit
		corpTagInit.from.m["nat"] = natTagInit
		natTagInit.to.m["corp"] = corpTagInit
		corpTagInit.from.m["windows-dc"] = windowsdcTagInit
		windowsdcTagInit.to.m["corp"] = corpTagInit

		tagsInit := &ProtectedTagMap{m: map[string]*Tag{
			"loadbalancer":  loadbalancerTagInit,
			"windows-dc":    windowsdcTagInit,
			"ssh":           sshTagInit,
			"ci":            ciTagInit,
			"storage":       storageTagInit,
			"reverse_proxy": reverseProxyTagInit,
			"nat":           natTagInit,
			"https":         httpsTagInit,
			"http":          httpTagInit,
			"corp":          corpTagInit,
		}}

		responses := []string{`["vm-2ba4d2f87","vm-7d1ff7af47","vm-a3ed2eed23","vm-b35b501"]`,
			`["vm-2ba4d2f87","vm-7d1ff7af47","vm-b35b501","vm-ec02d5c153"]`,
			`["vm-7d1ff7af47","vm-a3ed2eed23","vm-b35b501","vm-ec02d5c153"]`,
			`["vm-2ba4d2f87","vm-7d1ff7af47","vm-a3ed2eed23","vm-ec02d5c153"]`,
			`["vm-2ba4d2f87","vm-a3ed2eed23","vm-b35b501","vm-ec02d5c153"]`,
		}

		// Skipping some structures for validation since the assert fails to validate
		// recursive dependencies. Validating the results based on the responses
		exp := expect{dataInit, nil, tagsInit, nil, nil, nil, responses}
		return exp
	}

	expect_3 := func() expect {
		// Initialize tags
		dataInit := &Data{VMs: []*VM{{"vm-9ea3998", "frontend server", []string{"antivirus"}, nil},
			{"vm-5f3ad2b", "frontend server", []string{"http"}, nil},
			{"vm-d9e0825", "etcd node", []string{"ssh"}, nil},
			{"vm-59574582", "k8s node", []string{"antivirus", "ssh", "api", "windows-dc"}, nil},
			{"vm-f00923", "billing service", []string{"http", "dev", "k8s"}, nil},
			{"vm-575c4a", "rabbitmq", []string{"dev", "k8s"}, nil},
			{"vm-0c1791", "php app", []string{"http", "ci", "reverse_proxy", "dev"}, nil},
			{"vm-2987241", "dev-srv-5", []string{"k8s", "api", "nat", "reverse_proxy"}, nil},
			{"vm-ab51cba10", "ssh bastion", []string{"https", "storage", "loadbalancer", "corp", "django"}, nil},
			{"vm-a3660c", "frontend server", []string{"k8s", "ssh"}, nil},
			{"vm-864a94f", "kafka", []string{"dev"}, nil},
		}, FWRules: []*FWRule{{"fw-dd3c1e", "k8s", "django"},
			{"fw-1688373be0", "antivirus", "corp"},
			{"fw-f1fcfa", "http", "loadbalancer"},
			{"fw-93b1338c12", "api", "corp"},
			{"fw-e1d2dcbf3", "ssh", "storage"},
			{"fw-8e836298", "ssh", "django"},
			{"fw-742ac1", "k8s", "django"},
			{"fw-2bf982", "ssh", "django"},
			{"fw-9c95744ef", "reverse_proxy", "django"},
			{"fw-36b8c424e", "nat", "corp"},
			{"fw-87b059", "api", "django"},
			{"fw-2d389ea", "ci", "https"},
			{"fw-3a7c1e8", "ssh", "storage"},
			{"fw-3279918d6", "k8s", "storage"},
			{"fw-dc2742a339", "http", "django"},
			{"fw-9546840161", "reverse_proxy", "corp"},
			{"fw-966d9bd4", "api", "https"},
			{"fw-8ad7dc", "k8s", "corp"},
			{"fw-6c4649bc", "ci", "loadbalancer"},
			{"fw-c7ee8f3", "http", "loadbalancer"},
			{"fw-b3557d17", "windows-dc", "https"},
			{"fw-fc619d1", "ssh", "storage"},
			{"fw-2a83fc853", "ci", "loadbalancer"},
			{"fw-3069e1653", "api", "loadbalancer"},
			{"fw-2ab931510", "dev", "https"},
			{"fw-eb38b5d", "ssh", "loadbalancer"},
			{"fw-7f0d52f2", "k8s", "loadbalancer"},
			{"fw-e8cd165c", "ci", "loadbalancer"},
			{"fw-e85265e1", "api", "https"},
			{"fw-41c837edbb", "ci", "https"},
		}}

		responses := []string{`[]`,
			`[]`,
			`[]`,
			`[]`,
			`[]`,
			`[]`,
			`[]`,
			`[]`,
			`["vm-0c1791","vm-2987241","vm-575c4a","vm-59574582","vm-5f3ad2b","vm-864a94f","vm-9ea3998","vm-a3660c","vm-d9e0825","vm-f00923"]`,
			`[]`,
			`[]`,
		}

		// Skipping some structures for validation since the assert fails to validate
		// cross dependencies, it stucks. Validating the results based on the responses
		exp := expect{dataInit, nil, nil, nil, nil, nil, responses}
		return exp
	}

	var tests = []struct {
		name   string
		input  input
		expect expect
	}{
		{
			"input-0",
			input{"fixtures/input-0.json"},
			expect_0(),
		},
		{
			"input-1",
			input{"fixtures/input-1.json"},
			expect_1(),
		},
		{
			"input-2",
			input{"fixtures/input-2.json"},
			expect_2(),
		},
		{
			"input-3",
			input{"fixtures/input-3.json"},
			expect_3(),
		},
	}

	for _, test := range tests {
		content, err := ioutil.ReadFile(test.input.fixture)
		assert.NoError(t, err, "reading fixture "+test.input.fixture)

		data, err := loadData(content)
		assert.NoError(t, err, "loading fixture "+test.input.fixture)
		assert.EqualValues(t, test.expect.dataInit, data, "Loaded data")

		storage := &Storage{
			Map:        make(map[string]*VM, len(data.VMs)),
			data:       data,
			tags:       NewProtectedTagMap(),
			srcTagSets: NewProtectedTagSetMap(),
		}

		err = storage.processFWRules()
		assert.NoError(t, err, "processing FW Rules "+test.input.fixture)
		assert.Equal(t, test.expect.dataInit, data, "Loaded data")
		if test.expect.tagsInit != nil {
			assert.Equal(t, test.expect.tagsInit, storage.tags, "Tags after processing FW Rules")
		}

		err = storage.processVMs()
		assert.NoError(t, err, "processing VMs "+test.input.fixture)
		if test.expect.dataVMs != nil {
			assert.Equal(t, test.expect.dataVMs, storage.data, "Loaded data after processing VMs")
		}
		if test.expect.tagsVMs != nil {
			assert.Equal(t, test.expect.tagsVMs, storage.tags, "Tags after processing VMs")
		}
		if test.expect.srcTagSetsInit != nil {
			assert.Equal(t, test.expect.srcTagSetsInit, storage.srcTagSets, "Source tag sets after procesging VMs")
		}

		storage.processTagSets()
		if test.expect.srcTagSetsFinal != nil {
			assert.Equal(t, test.expect.srcTagSetsFinal, storage.srcTagSets, "Source tag sets after procesging tag sets")
		}

		for idx, resp := range test.expect.responses {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "http://localhost/api/v1/attack?vm_id="+storage.data.VMs[idx].VMId, nil)
			storage.ServeHTTP(w, req)
			body, _ := io.ReadAll(w.Result().Body)
			assert.Equal(t, resp, string(body))
		}
	}
}
