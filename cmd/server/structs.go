package main

type Data struct {
	VMs     []*VM     `json:"vms"`
	FWRules []*FWRule `json:"fw_rules"`
}

type VM struct {
	VMId      string   `json:"vm_id"`
	Name      string   `json:"name"`
	Tags      []string `json:"tags"`
	DstTagMap *ProtectedTagMap
}

type FWRule struct {
	FWId      string `json:"fw_id"`
	SourceTag string `json:"source_tag"`
	DestTag   string `json:"dest_tag"`
}
