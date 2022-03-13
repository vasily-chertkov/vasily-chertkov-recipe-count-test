package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Storage struct {
	Map        map[string]*VM
	data       *Data
	tags       *ProtectedTagMap
	srcTagSets *ProtectedTagSetMap
}

func NewStorage(content []byte) (*Storage, error) {
	start := time.Now()
	defer func(now time.Time) {
		log.Infoln("Preprocessing time:", time.Since(now))
	}(start)

	data, err := loadData(content)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		Map:        make(map[string]*VM, len(data.VMs)),
		data:       data,
		tags:       NewProtectedTagMap(),
		srcTagSets: NewProtectedTagSetMap(),
	}

	err = storage.processFWRules()
	if err != nil {
		return nil, err
	}

	err = storage.processVMs()
	if err != nil {
		return nil, err
	}

	storage.processTagSets()
	return storage, nil
}

func (s Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["vm_id"]
	if !ok {
		msg := "missing mandatory argument 'vm_id'"
		log.Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}
	if len(keys) > 1 {
		msg := "expected single argument 'vm_id', got multiple"
		log.Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	vmId := keys[0]
	vm, ok := s.Map[vmId]
	if !ok {
		msg := fmt.Sprintf("vm_id not found: %s", vmId)
		log.Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// All this stuff has to be done with interfaces,
	// but for now accessing the fields directly...
	if vm.DstTagMap == nil ||
		vm.DstTagMap.Link == nil ||
		vm.DstTagMap.Link.Slice == nil {
		w.Write([]byte("[]"))
		return
	}

	b := s.Map[vmId].DstTagMap.Link.Slice.Bytes(vmId)
	w.Write(b)
}

func loadData(content []byte) (*Data, error) {
	start := time.Now()
	defer func(now time.Time) {
		log.Infoln("Loading data time:", time.Since(now))
		PrintMemUsage()
	}(start)

	data := Data{}
	err := jsonfast.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("json malformed: %s", err)
	}
	return &data, nil
}

func (s *Storage) processFWRules() error {
	start := time.Now()
	defer func(now time.Time) {
		log.Infoln("Processing FW Rules time:", time.Since(now))
		PrintMemUsage()
	}(start)

	seen := make(map[string]struct{}, len(s.data.FWRules))
	for _, fwRule := range s.data.FWRules {
		if _, ok := seen[fwRule.FWId]; ok {
			msg := fmt.Sprintf("duplicates fw_id found: %s", fwRule.FWId)
			return fmt.Errorf(msg)
		}
		seen[fwRule.FWId] = exists

		dst, _ := s.tags.CreateIfNotExists(fwRule.DestTag)
		src, _ := s.tags.CreateIfNotExists(fwRule.SourceTag)
		dst.AddFrom(src)
		src.AddTo(dst)
	}
	return nil
}

func (s *Storage) processVMs() error {
	start := time.Now()
	defer func(now time.Time) {
		log.Infoln("Processing VMs time:", time.Since(now))
		PrintMemUsage()
	}(start)

	dstTagSets := NewProtectedTagSetMap()

	worker := func(id int, wg *sync.WaitGroup, vmChan <-chan *VM) {
		defer wg.Done()
		for vm := range vmChan {
			var dsttags []*Tag
			for _, t := range vm.Tags {
				if tag, ok := s.tags.Get(t); ok {
					tag.AddVM(vm)

					// If the tag has "From" directions, then this tag
					// was seen as a Destination in FW Rules
					if tag.HasFrom() {
						dsttags = append(dsttags, tag)
					}
				}
			}

			// No Destination Tags found for the VM, it's not reachable
			if len(dsttags) == 0 {
				continue
			}

			tagMap, created := dstTagSets.CreateIfNotExists(dsttags...)
			if created {
				var srctags []*Tag
				for _, dtag := range dsttags {
					srctags = append(srctags, dtag.FromItems()...)
				}
				tm, _ := s.srcTagSets.CreateIfNotExists(srctags...)
				tagMap.SetLink(tm)
			}

			vm.DstTagMap = tagMap
		}
	}

	workersCount := runtime.NumCPU()
	wg := new(sync.WaitGroup)
	wg.Add(workersCount)

	vmChan := make(chan *VM, len(s.data.VMs))
	for w := 1; w <= workersCount; w++ {
		go worker(w, wg, vmChan)
	}

	for _, vm := range s.data.VMs {
		if _, ok := s.Map[vm.VMId]; ok {
			msg := fmt.Sprintf("duplicates vm_id found: %s", vm.VMId)
			return fmt.Errorf(msg)
		}
		s.Map[vm.VMId] = vm
		vmChan <- vm
	}
	close(vmChan)

	wg.Wait()

	return nil
}

func (s *Storage) processTagSets() {
	start := time.Now()
	defer func(now time.Time) {
		log.Infoln("Processing Tag Sets time:", time.Since(now))
		PrintMemUsage()
	}(start)

	worker := func(id int, wg *sync.WaitGroup, tagMapChan <-chan *ProtectedTagMap) {
		defer wg.Done()
		for tagmap := range tagMapChan {
			seen := make(map[string]struct{})
			var vmIds []string
			for _, tag := range tagmap.Items() {
				for _, vm := range tag.VMs() {
					if _, ok := seen[vm.VMId]; ok {
						continue
					}
					seen[vm.VMId] = exists
					vmIds = append(vmIds, vm.VMId)
				}
			}
			s := NewStringSlice(vmIds...)
			tagmap.SetSlice(s)
		}
	}

	workersCount := runtime.NumCPU()
	wg := new(sync.WaitGroup)
	wg.Add(workersCount)

	tagMapChan := make(chan *ProtectedTagMap, s.srcTagSets.Len())
	for w := 1; w <= workersCount; w++ {
		go worker(w, wg, tagMapChan)
	}

	for _, tagMap := range s.srcTagSets.Items() {
		tagMapChan <- tagMap
	}
	close(tagMapChan)

	wg.Wait()

}
