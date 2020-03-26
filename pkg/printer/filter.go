package printer

import (
	"fmt"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

// Filter defines which results are to be shown within a printer
type Filter func([]*types.Result) []*types.Result

var mappings = map[string]Filter{
	"only-issues": OnlyIssues,
}

type FlagFilters struct {
	Filters []Filter
}

func NewFlagFilter() *FlagFilters {
	return &FlagFilters{
		Filters: make([]Filter, 0),
	}
}

func (flag *FlagFilters) Set(name string) error {
	filt, exist := mappings[name]
	if !exist {
		return fmt.Errorf("Unknown filter name %s", name)
	}
	flag.Filters = append(flag.Filters, filt)
	return nil
}

func (flag *FlagFilters) String() string {
	return ""
}

// OnlyIssues will remove all results that are OK
func OnlyIssues(res []*types.Result) []*types.Result {
	ret := make([]*types.Result, 0, cap(res))
	for _, r := range res {
		if r.IssueType != types.OK {
			ret = append(ret, r)
		}
	}
	return ret
}
