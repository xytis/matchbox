package http

import (
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
)

var (
	validMACStr = "52:da:00:89:d8:10"

	testProfileIgnitionYAML = &storagepb.Profile{
		Id:       "g1h2i3j4",
		Template: map[string]string{"ignition": "fake-template"},
	}

	testProfileGeneric = &storagepb.Profile{
		Id:       "g1h2i3j4",
		Template: map[string]string{"ignition": "fake-template"},
	}

	testGroupWithMAC = &storagepb.Group{
		Id:       "test-group",
		Name:     "test group",
		Profile:  "g1h2i3j4",
		Selector: map[string]string{"mac": validMACStr},
	}
)
