package storagepb

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
)

// ParseGroup parses bytes into a Group.
func ParseGroup(data []byte) (*Group, error) {
	richGroup := new(RichGroup)
	err := json.Unmarshal(data, richGroup)
	if err != nil {
		return nil, err
	}
	group, err := richGroup.ToGroup()
	if err != nil {
		return nil, err
	}
	if err := group.Normalize(); err != nil {
		return nil, err
	}
	return group, err
}

// Copy creates a copy of Group
func (g *Group) Copy() *Group {
	selectors := make(map[string]string)
	for k, v := range g.Selector {
		selectors[k] = v
	}
	return &Group{
		Id:       g.Id,
		Name:     g.Name,
		Profile:  g.Profile,
		Selector: selectors,
		Metadata: g.Metadata,
	}
}

// Matches returns true if the given labels satisfy all the selector
// requirements, false otherwise.
func (g *Group) Matches(labels map[string]string) bool {
	for key, val := range g.Selector {
		if labels == nil || labels[key] != val {
			return false
		}
	}
	return true
}

// Normalize normalizes Group selectors according to reserved selector rules
// which require "mac" addresses to be valid, normalized MAC addresses.
func (g *Group) Normalize() error {
	for key, val := range g.Selector {
		switch strings.ToLower(key) {
		case "mac":
			macAddr, err := net.ParseMAC(val)
			if err != nil {
				return err
			}
			// range iteration copy with mutable map
			g.Selector[key] = macAddr.String()
		}
	}
	return nil
}

// AssertValid validates a Group. Returns nil if there are no validation
// errors.
func (g *Group) AssertValid() error {
	if g.Id == "" {
		return ErrIdRequired
	}
	if g.Profile == "" {
		return ErrProfileRequired
	}
	return nil
}

// SelectorString returns Group selectors as a string of sorted key value
// pairs for comparisons and output.
func (g *Group) SelectorString() string {
	reqs := make([]string, 0, len(g.Selector))
	for key, value := range g.Selector {
		reqs = append(reqs, key+"="+value)
	}
	// sort by "key=value" pairs for a deterministic ordering
	sort.StringSlice(reqs).Sort()
	return strings.Join(reqs, ",")
}

// MetadataPrettyString returns metadata in pretty string format
func (g *Group) MetadataPrettyString() string {
	if g.Metadata == nil {
		return "{}"
	}
	metadata := make(map[string]interface{})
	err := json.Unmarshal(g.Metadata, &metadata)
	if err != nil {
		return fmt.Sprintf("unable to unmarshal metadata: %v\n%s\n", err, g.Metadata)
	}
	pretty, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Sprintf("unable to marshal metadata: %v\n%s\n", err, g.Metadata)
	}
	return string(pretty)
}

// ToRichGroup converts a Group into a RichGroup suitable for writing and
// user manipulation.
func (g *Group) ToRichGroup() (*RichGroup, error) {
	metadata := make(map[string]interface{})
	if g.Metadata != nil {
		err := json.Unmarshal(g.Metadata, &metadata)
		if err != nil {
			return nil, err
		}
	}
	//TODO: Not a full copy?!
	return &RichGroup{
		Id:       g.Id,
		Name:     g.Name,
		Profile:  g.Profile,
		Selector: g.Selector,
		Metadata: metadata,
	}, nil
}

// ByReqs defines a collection of Group structs which have a deterministic
// sorted order by decreasing number of Requirements, then by sorted key/value
// strings. For example, a Group with Requirements {a:b, c:d} should be ordered
// after one with {a:b} and before one with {a:d, c:d}.
// (legacy) If two groups share selectors, Id is used to tiebreak.
type ByReqs []*Group

func (groups ByReqs) Len() int {
	return len(groups)
}

func (groups ByReqs) Swap(i, j int) {
	groups[i], groups[j] = groups[j], groups[i]
}

func (groups ByReqs) Less(i, j int) bool {
	if len(groups[i].Selector) == len(groups[j].Selector) {
		if groups[i].SelectorString() == groups[j].SelectorString() {
			return groups[i].Id < groups[j].Id
		}
		return groups[i].SelectorString() < groups[j].SelectorString()
	}
	return len(groups[i].Selector) > len(groups[j].Selector)
}

// RichGroup is a user provided Group definition.
type RichGroup struct {
	// machine readable Id
	Id string `json:"id,omitempty"`
	// Human readable name
	Name string `json:"name,omitempty"`
	// Profile id
	Profile string `json:"profile,omitempty"`
	// Selectors to match machines
	Selector map[string]string `json:"selector,omitempty"`
	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToGroup converts a user provided RichGroup into a Group which can be
// serialized as a protocol buffer.
func (rg *RichGroup) ToGroup() (*Group, error) {
	var metadata []byte
	if rg.Metadata != nil {
		var err error
		metadata, err = json.Marshal(rg.Metadata)
		if err != nil {
			return nil, err
		}
	}
	//TODO: Not a full copy?!
	return &Group{
		Id:       rg.Id,
		Name:     rg.Name,
		Profile:  rg.Profile,
		Selector: rg.Selector,
		Metadata: metadata,
	}, nil
}
