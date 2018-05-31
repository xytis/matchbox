package server

import (
	"sort"

	"context"

	"github.com/coreos/matchbox/matchbox/server/config"
	pb "github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/pkg/errors"
)

// Possible service errors
var (
	ErrNoMatchingGroup   = errors.New("matchbox: No matching Group")
	ErrNoMatchingProfile = errors.New("matchbox: No matching Profile")
)

// Server defines the matchbox server interface.
type Server interface {
	// SelectGroup returns the Group matching the given labels.
	SelectGroup(context.Context, *pb.SelectGroupRequest) (*storagepb.Group, error)
	// SelectProfile returns the Profile matching the given labels.
	SelectProfile(context.Context, *pb.SelectProfileRequest) (*storagepb.Profile, error)

	// Create or update a Group.
	GroupPut(context.Context, *pb.GroupPutRequest) (*storagepb.Group, error)
	// Get a machine Group by id.
	GroupGet(context.Context, *pb.GroupGetRequest) (*storagepb.Group, error)
	// Delete a machine Group by id.
	GroupDelete(context.Context, *pb.GroupDeleteRequest) error
	// List all machine Groups.
	GroupList(context.Context, *pb.GroupListRequest) ([]*storagepb.Group, error)

	// Create or update a Profile.
	ProfilePut(context.Context, *pb.ProfilePutRequest) (*storagepb.Profile, error)
	// Get a Profile by id.
	ProfileGet(context.Context, *pb.ProfileGetRequest) (*storagepb.Profile, error)
	// Delete a Profile by id.
	ProfileDelete(context.Context, *pb.ProfileDeleteRequest) error
	// List all Profiles.
	ProfileList(context.Context, *pb.ProfileListRequest) ([]*storagepb.Profile, error)

	// Create or update a template.
	TemplatePut(context.Context, *pb.TemplatePutRequest) (*storagepb.Template, error)
	// Get a template by name.
	TemplateGet(context.Context, *pb.TemplateGetRequest) (*storagepb.Template, error)
	// Delete a template by name.
	TemplateDelete(context.Context, *pb.TemplateDeleteRequest) error
	// List all Templates
	TemplateList(context.Context, *pb.TemplateListRequest) ([]*storagepb.Template, error)
}

// server implements the Server interface.
type server struct {
	store storage.Store
}

// NewServer returns a new Server.
func NewServer(config *config.Config) Server {
	store := createStore(config)
	storage.AssertDefaultTemplates(store)
	return &server{
		store: store,
	}
}

func createStore(config *config.Config) storage.Store {
	switch config.StoreBackend {
	case "filesystem":
		store, err := storage.NewFileStore(config.FileStoreConfig)
		if err != nil {
			panic(errors.Wrap(err, "failure creating filesystem store"))
		}
		return store
	case "etcd":
		store, err := storage.NewEtcdStore(config.EtcdStoreConfig)
		if err != nil {
			panic(errors.Wrap(err, "failure creating etcd store"))
		}
		return store
	default:
		panic("unsuported storage engine")
	}

}

// SelectGroup selects the Group whose selector matches the given labels.
// Groups are evaluated in sorted order from most selectors to least, using
// alphabetical order as a deterministic tie-breaker.
func (s *server) SelectGroup(ctx context.Context, req *pb.SelectGroupRequest) (*storagepb.Group, error) {
	groups, err := s.store.GroupList()
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(storagepb.ByReqs(groups)))
	for _, group := range groups {
		if group.Matches(req.Labels) {
			return group, nil
		}
	}
	return nil, ErrNoMatchingGroup
}

func (s *server) SelectProfile(ctx context.Context, req *pb.SelectProfileRequest) (*storagepb.Profile, error) {
	group, err := s.SelectGroup(ctx, &pb.SelectGroupRequest{Labels: req.Labels})
	if err == nil {
		// lookup the Profile by id
		profile, err := s.ProfileGet(ctx, &pb.ProfileGetRequest{Id: group.Profile})
		if err == nil {
			return profile, nil
		}
		return nil, ErrNoMatchingProfile
	}
	return nil, ErrNoMatchingGroup
}

func (s *server) GroupPut(ctx context.Context, req *pb.GroupPutRequest) (*storagepb.Group, error) {
	if err := req.Group.AssertValid(); err != nil {
		return nil, err
	}
	err := s.store.GroupPut(req.Group)
	if err != nil {
		return nil, err
	}
	return req.Group, nil
}

func (s *server) GroupGet(ctx context.Context, req *pb.GroupGetRequest) (*storagepb.Group, error) {
	group, err := s.store.GroupGet(req.Id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *server) GroupDelete(ctx context.Context, req *pb.GroupDeleteRequest) error {
	return s.store.GroupDelete(req.Id)
}

func (s *server) GroupList(ctx context.Context, req *pb.GroupListRequest) ([]*storagepb.Group, error) {
	groups, err := s.store.GroupList()
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *server) ProfilePut(ctx context.Context, req *pb.ProfilePutRequest) (*storagepb.Profile, error) {
	if err := req.Profile.AssertValid(); err != nil {
		return nil, err
	}
	err := s.store.ProfilePut(req.Profile)
	if err != nil {
		return nil, err
	}
	return req.Profile, nil
}

func (s *server) ProfileGet(ctx context.Context, req *pb.ProfileGetRequest) (*storagepb.Profile, error) {
	profile, err := s.store.ProfileGet(req.Id)
	if err != nil {
		return nil, err
	}
	if err := profile.AssertValid(); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *server) ProfileDelete(ctx context.Context, req *pb.ProfileDeleteRequest) error {
	return s.store.ProfileDelete(req.Id)
}

func (s *server) ProfileList(ctx context.Context, req *pb.ProfileListRequest) ([]*storagepb.Profile, error) {
	profiles, err := s.store.ProfileList()
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (s *server) TemplatePut(ctx context.Context, req *pb.TemplatePutRequest) (*storagepb.Template, error) {
	err := s.store.TemplatePut(req.Template)
	if err != nil {
		return nil, err
	}
	return req.Template, err
}

func (s *server) TemplateGet(ctx context.Context, req *pb.TemplateGetRequest) (*storagepb.Template, error) {
	return s.store.TemplateGet(req.Id)
}

func (s *server) TemplateDelete(ctx context.Context, req *pb.TemplateDeleteRequest) error {
	return s.store.TemplateDelete(req.Id)
}

func (s *server) TemplateList(ctx context.Context, req *pb.TemplateListRequest) ([]*storagepb.Template, error) {
	templates, err := s.store.TemplateList()
	if err != nil {
		return nil, err
	}
	return templates, nil
}
