package main

import (
	"sort"
	"strings"
	"sync"
)

func NewProtectedTagMap() *ProtectedTagMap {
	return &ProtectedTagMap{m: make(map[string]*Tag)}
}

type ProtectedTagMap struct {
	sync.RWMutex
	m     map[string]*Tag
	Link  *ProtectedTagMap
	Slice *Slice
}

func (tm *ProtectedTagMap) SetLink(link *ProtectedTagMap) {
	tm.Link = link
}

func (tm *ProtectedTagMap) SetSlice(s *Slice) {
	tm.Slice = s
}

func (tm *ProtectedTagMap) Add(tags ...*Tag) {
	tm.Lock()
	defer tm.Unlock()

	for _, tag := range tags {
		if _, ok := tm.m[tag.Name]; !ok {
			tm.m[tag.Name] = tag
		}
	}
}

func (tm *ProtectedTagMap) CreateIfNotExists(key string) (*Tag, bool) {
	tm.Lock()
	defer tm.Unlock()

	var created bool
	tag, ok := tm.m[key]
	if !ok {
		tag = NewTag(key)
		tm.m[key] = tag
		created = true
	}
	return tag, created
}

func (tm *ProtectedTagMap) Get(key string) (*Tag, bool) {
	tm.RLock()
	defer tm.RUnlock()
	tag, ok := tm.m[key]
	return tag, ok
}

func (tm *ProtectedTagMap) Len() int {
	tm.RLock()
	defer tm.RUnlock()
	return len(tm.m)
}

func (tm *ProtectedTagMap) Keys() []string {
	tm.RLock()
	defer tm.RUnlock()
	keys := make([]string, 0, len(tm.m))
	for key := range tm.m {
		keys = append(keys, key)
	}
	return keys
}

func (tm *ProtectedTagMap) Items() []*Tag {
	tm.RLock()
	defer tm.RUnlock()
	items := make([]*Tag, 0, len(tm.m))
	for _, tag := range tm.m {
		items = append(items, tag)
	}
	return items
}

func NewProtectedTagSetMap() *ProtectedTagSetMap {
	return &ProtectedTagSetMap{m: make(map[string]*ProtectedTagMap)}
}

type ProtectedTagSetMap struct {
	sync.RWMutex
	m map[string]*ProtectedTagMap
}

func (tsm *ProtectedTagSetMap) CreateIfNotExists(tags ...*Tag) (*ProtectedTagMap, bool) {
	var created bool

	var keys []string
	seen := make(map[string]struct{})
	for _, tag := range tags {
		if _, ok := seen[tag.Name]; ok {
			continue
		}
		seen[tag.Name] = exists
		keys = append(keys, tag.Name)
	}

	sort.Strings(keys)
	joinedKey := strings.Join(keys, "|")

	tsm.RLock()
	ptm, ok := tsm.m[joinedKey]
	tsm.RUnlock()
	if ok {
		return ptm, created
	}

	tsm.Lock()
	defer tsm.Unlock()
	ptm, ok = tsm.m[joinedKey]
	if !ok {
		ptm = NewProtectedTagMap()
		ptm.Add(tags...)
		tsm.m[joinedKey] = ptm
		created = true
	}

	return ptm, created
}

func (tsm *ProtectedTagSetMap) Len() int {
	tsm.RLock()
	defer tsm.RUnlock()
	return len(tsm.m)
}

func (tsm *ProtectedTagSetMap) Items() []*ProtectedTagMap {
	tsm.RLock()
	defer tsm.RUnlock()
	items := make([]*ProtectedTagMap, 0, len(tsm.m))
	for _, tagMap := range tsm.m {
		items = append(items, tagMap)
	}
	return items
}

type Tag struct {
	Name string
	from *ProtectedTagMap
	to   *ProtectedTagMap

	vms   map[string]*VM
	mutex sync.RWMutex
}

func NewTag(name string) *Tag {
	return &Tag{Name: name,
		from: NewProtectedTagMap(),
		to:   NewProtectedTagMap(),
		vms:  make(map[string]*VM),
	}
}

func (t *Tag) AddFrom(from *Tag) {
	t.from.Add(from)
}

func (t *Tag) HasFrom() bool {
	return t.from.Len() > 0
}

func (t *Tag) FromKeys() []string {
	return t.from.Keys()
}

func (t *Tag) FromItems() []*Tag {
	return t.from.Items()
}

func (t *Tag) AddTo(to *Tag) {
	t.to.Add(to)
}

func (t *Tag) HasTo() bool {
	return t.to.Len() > 0
}

func (t *Tag) ToKeys() []string {
	return t.to.Keys()
}

func (t *Tag) ToItems() []*Tag {
	return t.to.Items()
}

func (t *Tag) AddVM(vm *VM) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	var added bool
	if _, ok := t.vms[vm.VMId]; !ok {
		t.vms[vm.VMId] = vm
		added = true
	}
	return added
}

func (t *Tag) VMs() []*VM {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var vms []*VM
	for _, vm := range t.vms {
		vms = append(vms, vm)
	}
	return vms
}

func (t *Tag) VMIds() []string {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var vmIds []string
	for vmId := range t.vms {
		vmIds = append(vmIds, vmId)
	}
	return vmIds
}
