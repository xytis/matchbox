package http

import (
	"context"
	"errors"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/divideandconquer/go-merge/merge"
	"github.com/hashicorp/go-multierror"
)

// unexported key prevents collisions
type key int

const (
	profileKey key = iota
	groupKey
	labelsKey
)

var (
	errNoProfileFromContext = errors.New("api: Context missing a Profile")
	errNoGroupFromContext   = errors.New("api: Context missing a Group")
	errNoLabelsFromContext  = errors.New("api: Context missing parsed Labels")
)

// withProfile returns a copy of ctx that stores the given Profile.
func withProfile(ctx context.Context, profile *storagepb.Profile) context.Context {
	return context.WithValue(ctx, profileKey, profile)
}

// profileFromContext returns the Profile from the ctx.
func profileFromContext(ctx context.Context) (*storagepb.Profile, error) {
	profile, ok := ctx.Value(profileKey).(*storagepb.Profile)
	if !ok {
		return nil, errNoProfileFromContext
	}
	return profile, nil
}

// withGroup returns a copy of ctx that stores the given Group.
func withGroup(ctx context.Context, group *storagepb.Group) context.Context {
	return context.WithValue(ctx, groupKey, group)
}

// groupFromContext returns the Group from the ctx.
func groupFromContext(ctx context.Context) (*storagepb.Group, error) {
	group, ok := ctx.Value(groupKey).(*storagepb.Group)
	if !ok {
		return nil, errNoGroupFromContext
	}
	return group, nil
}

func withLabels(ctx context.Context, labels map[string]string) context.Context {
	return context.WithValue(ctx, labelsKey, labels)
}

func labelsFromContext(ctx context.Context) (map[string]string, error) {
	labels, ok := ctx.Value(labelsKey).(map[string]string)
	if !ok {
		return nil, errNoLabelsFromContext
	}
	return labels, nil
}

func mergeMetadata(ctx context.Context) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})
	var errors error
	if group, err := groupFromContext(ctx); err == nil {
		if rg, err := group.ToRichGroup(); err == nil {
			metadata = merge.Merge(metadata, rg.Selector).(map[string]interface{})
			metadata = merge.Merge(metadata, rg.Metadata).(map[string]interface{})
		} else {
			errors = multierror.Append(errors, err)
		}
	}
	if profile, err := profileFromContext(ctx); err == nil {
		if rp, err := profile.ToRichProfile(); err == nil {
			metadata = merge.Merge(metadata, rp.Metadata).(map[string]interface{})
		} else {
			errors = multierror.Append(errors, err)
		}
	}
	if labels, err := labelsFromContext(ctx); err == nil {
		if submap, found := metadata["labels"]; found {
			metadata["labels"] = merge.Merge(submap, labels)
		} else {
			metadata["labels"] = merge.Merge(map[string]string{}, labels)
		}
	}
	return metadata, errors
}
