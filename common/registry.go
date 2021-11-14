package common

import (
	"sort"
	"sync"
)

var once sync.Once

type ChecksRegistry struct {
	checks []*CheckDef
}

var instance *ChecksRegistry

func GetRegistry() *ChecksRegistry {
	once.Do(func() {
		instance = &ChecksRegistry{
			checks: make([]*CheckDef, 0),
		}
	})
	return instance
}

func (reg *ChecksRegistry) Register(def *CheckDef) {
	reg.checks = append(reg.checks, def)

	// Perform sorting
	sort.Slice(reg.checks, func(i, j int) bool {
		if reg.checks[i].Group < reg.checks[j].Group {
			return true
		}
		if reg.checks[i].Group > reg.checks[j].Group {
			return false
		}
		return reg.checks[i].Name < reg.checks[j].Name
	})

}

func (reg *ChecksRegistry) GetAllChecks() []*CheckDef {
	return reg.checks
}
