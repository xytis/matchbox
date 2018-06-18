package http

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/divideandconquer/go-merge/merge"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

// unexported key prevents collisions
type key int

const (
	profileKey key = iota
	groupKey
	labelsKey
	requestIDKey
)

var (
	errNoProfileFromContext   = errors.New("api: Context missing a Profile")
	errNoGroupFromContext     = errors.New("api: Context missing a Group")
	errNoLabelsFromContext    = errors.New("api: Context missing parsed Labels")
	errNoRequestIDFromContext = errors.New("api: Context missing Request ID")
)

type unwrappedContext struct {
	context.Context
	Labels    map[string]string
	Group     *storagepb.Group
	Profile   *storagepb.Profile
	Metadata  map[string]interface{}
	RequestID string
}

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

// withRequestId returns a copy of ctx that stores current request id
func withRequestID(ctx context.Context, rID string) context.Context {
	return context.WithValue(ctx, requestIDKey, rID)
}

// requestIDFromContext returns the RequestId from the ctx.
func requestIDFromContext(ctx context.Context) (string, error) {
	rID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return "", errNoRequestIDFromContext
	}
	return rID, nil
}

func (s *Server) mergeMetadata(ctx context.Context) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})
	if gm, err := s.globalMetadata(); err == nil {
		metadata["matchbox"] = gm
	} else {
		return nil, err
	}

	if group, err := groupFromContext(ctx); err == nil {
		if rg, err := group.ToRichGroup(); err == nil {
			metadata = merge.Merge(metadata, rg.Metadata).(map[string]interface{})
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	if profile, err := profileFromContext(ctx); err == nil {
		if rp, err := profile.ToRichProfile(); err == nil {
			metadata = merge.Merge(metadata, rp.Metadata).(map[string]interface{})
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
	if labels, err := labelsFromContext(ctx); err == nil {
		if submap, found := metadata["label"]; found {
			metadata["label"] = merge.Merge(submap, labels)
		} else {
			metadata["label"] = merge.Merge(map[string]string{}, labels)
		}
	} else {
		return nil, err
	}
	return metadata, nil
}

// wrapContext fills context with parsed query information
func (s *Server) wrapContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		labels := labelsFromRequest(s.logger, req)
		ctx = withLabels(ctx, labels)
		if group, err := s.core.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: labels}); err == nil {
			ctx = withGroup(ctx, group)
			if profile, err := s.core.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile}); err == nil {
				ctx = withProfile(ctx, profile)
			}
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// unwrapContext performs reverse operation to wrapContext
func (s *Server) unwrapContext(ctx context.Context) (unwrappedContext, error) {
	var err error
	result := unwrappedContext{Context: ctx}
	result.RequestID, err = requestIDFromContext(ctx)
	if err != nil {
		result.RequestID = "undefined"
	}

	result.Labels, err = labelsFromContext(ctx)
	if err != nil {
		return result, err
	}

	result.Group, err = groupFromContext(ctx)
	if err != nil {
		return result, err
	}

	result.Profile, err = profileFromContext(ctx)
	if err != nil {
		return result, err
	}

	result.Metadata, err = s.mergeMetadata(ctx)
	if err != nil {
		s.logger.Warn("metadata not merged",
			zap.Error(err),
			zap.String("labels", fmt.Sprintf("%v", result.Labels)),
			zap.String("group", result.Group.Id),
			zap.String("profile", result.Profile.Id),
		)
		return result, err
	}

	return result, nil
}
