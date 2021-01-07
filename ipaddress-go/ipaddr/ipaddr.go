package ipaddr

import "unsafe"

type IPVersion string

const (
	PrefixLenSeparator = '/'

	INDETERMINATE_VERSION IPVersion = ""
	IPv4                  IPVersion = "IPv4"
	IPv6                  IPVersion = "IPv6"
)

func (version IPVersion) isIPv6() bool {
	return version == IPv6
}

func (version IPVersion) isIPv4() bool {
	return version == IPv4
}

func (version IPVersion) isIndeterminate() bool {
	return version == INDETERMINATE_VERSION
}

// returns an index starting from 0 with INDETERMINATE_VERSION being the highest
func (version IPVersion) index() int {
	if version.isIPv4() {
		return 0
	} else if version.isIPv6() {
		return 1
	}
	return 2
}

func (version IPVersion) String() string {
	return string(version)
}

//
//
//
//type IPAddress struct {
//	addressInternal
//}

// necessary to avoid direct access to IPAddress
type ipAddressInternal struct {
	addressInternal
}

func (addr *ipAddressInternal) ToAddress() *Address {
	return (*Address)(unsafe.Pointer(addr))
}

func (addr *ipAddressInternal) toIPAddress() *IPAddress {
	return (*IPAddress)(unsafe.Pointer(addr))
}

//func (addr *ipAddressInternal) checkIdentity(section *IPAddressSection) *IPAddress {
//	sec := section.toAddressSection()
//	if sec == addr.section {
//		return addr.toIPAddress()
//	}
//	return &IPAddress{ipAddressInternal{addressInternal{
//		section: sec, zone: addr.zone, cache: &addressCache{},
//	}}}
//}

// isIPv4() returns whether this matches an IPv4 address.
// It needs to match an IPv4 section and also have 4 segments.
func (addr *ipAddressInternal) isIPv4() bool {
	if addr.section == nil {
		return false
	}
	return addr.section.matchesIPv4Address()
}

// isIPv6() returns whether this matches an IPv6 address.
// It needs to match an IPv6 section and also have 8 segments.
func (addr *ipAddressInternal) isIPv6() bool {
	if addr.section == nil {
		return false
	}
	return addr.section.matchesIPv6Address()
}

func (addr *ipAddressInternal) getIPVersion() IPVersion {
	if addr.isIPv4() {
		return IPv4
	} else if addr.isIPv6() {
		return IPv6
	}
	return INDETERMINATE_VERSION
}

//func (addr *ipAddressInternal) getSection() *IPAddressSection {
//	return addr.section.ToIPAddressSection()
//}

func (addr *ipAddressInternal) GetNetworkPrefixLength() PrefixLen {
	section := addr.section
	if section == nil {
		return nil
	}
	return section.ToIPAddressSection().GetNetworkPrefixLength()
}

func (addr *ipAddressInternal) GetBlockMaskPrefixLength(network bool) PrefixLen {
	section := addr.section
	if section == nil {
		return nil
	}
	return section.ToIPAddressSection().GetBlockMaskPrefixLength(network)
}

func (addr *ipAddressInternal) GetSegment(index int) *IPAddressSegment {
	return addr.getSegment(index).ToIPAddressSegment()
}

var zeroIPAddr *IPAddress

func init() {
	zeroIPAddr = &IPAddress{
		ipAddressInternal{
			addressInternal{
				section: &AddressSection{},
				cache:   &addressCache{},
			},
		},
	}
}

//
//
// IPAddress represents an IPAddress, either IPv4 or IPv6.
// Only the zero-value IPAddress can be neither IPv4 or IPv6.
// The zero value has no segments, which is not compatible with zero value for ivp4 or ipv6.
type IPAddress struct {
	ipAddressInternal
}

func (addr *IPAddress) init() *IPAddress {
	if addr.section == nil {
		return zeroIPAddr // this has a zero section
	}
	return addr
}

func (addr IPAddress) String() string {
	//address := addr.init()
	//if address.section.cache != nil {
	//	addrType := address.section.cache.addrType
	//	_ = addrType
	//	//TODO a different default string if we know we are IPv4 or IPv6.  But we must do full check, same as when calling ToIPvxAddress() or ToIPvxAddressSection(), so that the result of this is consistent.
	//}
	return addr.init().ipAddressInternal.String()
}

func (addr *IPAddress) GetLower() *IPAddress {
	addr = addr.init()
	return addr.getLower().ToIPAddress()
}

func (addr *IPAddress) GetUpper() *IPAddress {
	addr = addr.init()
	return addr.getUpper().ToIPAddress()
}

func (addr *IPAddress) ToPrefixBlock() *IPAddress {
	addr = addr.init()
	prefixLength := addr.GetNetworkPrefixLength()
	if prefixLength == nil {
		return addr
	}
	return addr.ToPrefixBlockLen(*prefixLength)
}

func (addr *IPAddress) ToPrefixBlockLen(prefLen BitCount) *IPAddress {
	addr = addr.init()
	return addr.checkIdentity(addr.section.toPrefixBlockLen(prefLen)).ToIPAddress()
}

func (addr *IPAddress) IsIPv4() bool {
	if addr == nil {
		return false
	}
	return addr.isIPv4()
}

func (addr *IPAddress) IsIPv6() bool {
	if addr == nil {
		return false
	}
	return addr.isIPv6()
}

func (addr *IPAddress) GetIPVersion() IPVersion {
	addr = addr.init()
	return addr.getIPVersion()
}

// this makes no sense in the golang world, since it cannot be customized since not virtual
// gave a lot of thought to this, you cannot stick a converted object in each and every address,
// and so there really is no way to do this internally
//func (addr *IPAddress) IsIPv6Convertible() bool {
//	return addr.getConverter().IsIPv6Convertible(addr)
//}
//
//func (addr *IPAddress) IsIPv4Convertible() bool {
//	return addr.getConverter().IsIPv4Convertible(addr)
//}

func (addr *IPAddress) ToIPv6Address() *IPv6Address {
	if addr != nil {
		//addr = addr.init()
		if addr.isIPv6() {
			return (*IPv6Address)(unsafe.Pointer(addr))
		}
		//TODO consider allowing IPv4-mapped, see waht golang does, but consider we might not be ipv4 or ipv6 if zero-valued
		//return addr.cache.network.(IPAddressNetwork).GetConverter().ToIPv6(addr)
	}
	return nil
}

func (addr *IPAddress) ToIPv4Address() *IPv4Address {
	if addr != nil {
		//addr = addr.init()
		if addr.isIPv4() {
			return (*IPv4Address)(unsafe.Pointer(addr))
		}
		//return addr.cache.network.(IPAddressNetwork).GetConverter().ToIPv4(addr)
	}
	return nil
}

func (addr *IPAddress) SpanWithRange(other *IPAddress) *IPAddressSeqRange {
	if thisAddr := addr.ToIPv4Address(); thisAddr != nil {
		if oth := other.ToIPv4Address(); oth != nil {
			return thisAddr.SpanWithRange(oth).ToIPAddressSeqRange()
		}
	} else if thisAddr := addr.ToIPv6Address(); thisAddr != nil {
		if oth := other.ToIPv6Address(); oth != nil {
			return thisAddr.SpanWithRange(oth).ToIPAddressSeqRange()
		}
	}
	return nil
}

func (addr *IPAddress) Mask(other *IPAddress) *IPAddress {
	if thisAddr := addr.ToIPv4Address(); thisAddr != nil {
		if oth := other.ToIPv4Address(); oth != nil {
			return thisAddr.Mask(oth).ToIPAddress()
		}
	} else if thisAddr := addr.ToIPv6Address(); thisAddr != nil {
		if oth := other.ToIPv6Address(); oth != nil {
			return thisAddr.Mask(oth).ToIPAddress()
		}
	}
	return nil
}

func (addr *IPAddress) ToSequentialRange() *IPAddressSeqRange {
	if addr != nil {
		if addr.IsIPv4() {
			return addr.ToIPv4Address().ToSequentialRange().ToIPAddressSeqRange()
		} else if addr.IsIPv6() {
			return addr.ToIPv6Address().ToSequentialRange().ToIPAddressSeqRange()
		}
	}
	return nil
}

// TODO make sure everything in IPv4 and IPv6 is "overridden", in the sense all methods will check for no divisions and
// create the default zero-segments if necessary, so we never expose a zero value with 0 segments
// The zero values of everythign else will have sections with no segments
