package base

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/internal"
)

// PrefixNLRI defines Prefix NLRI onject
// https://tools.ietf.org/html/rfc7752#section-3.2
type PrefixNLRI struct {
	ProtocolID uint8
	Identifier uint64
	LocalNode  *NodeDescriptor
	Prefix     *PrefixDescriptor
}

func (p *PrefixNLRI) String() string {
	var s string
	s += fmt.Sprintf("Protocol ID: %s\n", internal.ProtocolIDString(p.ProtocolID))
	s += fmt.Sprintf("Identifier: %d\n", p.Identifier)
	s += p.LocalNode.String()
	s += p.Prefix.String()

	return s
}

// MarshalJSON defines a method to Marshal Prefix NLRI object into JSON format
func (p *PrefixNLRI) MarshalJSON() ([]byte, error) {
	var jsonData []byte

	jsonData = append(jsonData, '{')
	jsonData = append(jsonData, []byte("\"protocolID\":")...)
	jsonData = append(jsonData, []byte(fmt.Sprintf("%d,", p.ProtocolID))...)
	jsonData = append(jsonData, []byte("\"identifier\":")...)
	jsonData = append(jsonData, []byte(fmt.Sprintf("%d,", p.Identifier))...)
	jsonData = append(jsonData, []byte("\"localNode\":")...)
	if p.LocalNode != nil {
		b, err := json.Marshal(p.LocalNode)
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, b...)
		jsonData = append(jsonData, ',')
	} else {
		jsonData = append(jsonData, "{},"...)
	}
	jsonData = append(jsonData, []byte("\"prefix\":")...)
	if p.Prefix != nil {
		b, err := json.Marshal(p.Prefix)
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, b...)
	} else {
		jsonData = append(jsonData, "{}"...)
	}
	jsonData = append(jsonData, '}')

	return jsonData, nil
}

// UnmarshalPrefixNLRI builds Prefix NLRI object
func UnmarshalPrefixNLRI(b []byte) (*PrefixNLRI, error) {
	glog.V(6).Infof("PrefixNLRI Raw: %s", internal.MessageHex(b))
	pr := PrefixNLRI{}
	p := 0
	pr.ProtocolID = b[p]
	p++
	pr.Identifier = binary.BigEndian.Uint64(b[p : p+8])
	p += 8
	// Get Node Descriptor's length, skip Node Descriptor Type
	ndl := binary.BigEndian.Uint16(b[p+2 : p+4])
	ln, err := UnmarshalNodeDescriptor(b[p : p+int(ndl)])
	if err != nil {
		return nil, err
	}
	pr.LocalNode = ln
	// Skip Node Descriptor Type and Length 4 bytes
	p += 4
	p += int(ndl)
	pn, err := UnmarshalPrefixDescriptor(b[p:len(b)])
	if err != nil {
		return nil, err
	}
	pr.Prefix = pn

	return &pr, nil
}
