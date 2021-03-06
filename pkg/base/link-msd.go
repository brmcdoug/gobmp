package base

import (
	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/internal"
)

// LinkMSD defines Link MSD object
// No RFC yet
type LinkMSD struct {
	MSD []MSDTV
}

func (msd *LinkMSD) String() string {
	var s string
	s += "Link MSD TVs:" + "\n"
	for _, tv := range msd.MSD {
		s += tv.String()
	}

	return s
}

// UnmarshalLinkMSD build Link MSD object
func UnmarshalLinkMSD(b []byte) (*LinkMSD, error) {
	glog.V(6).Infof("Link MSD Raw: %s", internal.MessageHex(b))
	msd := LinkMSD{}
	p := 0
	tvsv, err := UnmarshalMSDTV(b[p : p+len(b)])
	if err != nil {
		return nil, err
	}
	msd.MSD = tvsv

	return &msd, nil
}
