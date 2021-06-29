package ipaddr

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type AddressNetwork interface {
	//GetAddressCreator() AddressCreator
}

//TODO I think I probably want to get rid of the address creators from networks (but they are still useful when passing into certain functions), I realize now they make little sense
//But I will still have caching.

// IPAddressNetwork represents the full collection of addresses for a given IP version.
// You can create your own network objects satisfying this interface, allowing you to create your own address types,
// or to provide your own IP address conversion between IPv4 and IPv6.
// When creating your own network, for IP addresses to be associated with it, you must:
// - create each address using the creator methods in the instance creator returned from GetIPAddressCreator(),
//	which will associate each address with said network when creating the address
// - return the network object from the IPAddressStringParameters implementation used for parsing an IPAddressString,
//	which will associate the parsed address with the network
// Addresses deprived from an existing address, using masking, iterating, or any other address manipulation,
// will be associated with the same network as the original address, by using the network's address creator instance.
// Addresses created by instantiation not through the network's creator instance will be associated with the default network.
type IPAddressNetwork interface {
	AddressNetwork

	//GetIPAddressCreator() IPAddressCreator

	GetLoopback() *IPAddress

	//GetNetworkIPAddress(PrefixLen) *IPAddress

	GetNetworkMask(prefixLength BitCount) *IPAddress

	GetPrefixedNetworkMask(prefixLength BitCount) *IPAddress

	GetHostMask(prefixLength BitCount) *IPAddress

	GetPrefixedHostMask(prefixLength BitCount) *IPAddress
}

type ipAddressNetwork struct {
	subnetsMasksWithPrefix, subnetMasks, hostMasksWithPrefix, hostMasks []*IPAddress
}

//
//
//
//
//
type IPv6AddressNetwork struct {
	ipAddressNetwork
	creator IPv6AddressCreator
}

func (network *IPv6AddressNetwork) GetIPv6AddressCreator() *IPv6AddressCreator {
	return &network.creator
}

//func (network *IPv6AddressNetwork) GetIPAddressCreator() IPAddressCreator {
//	return network.GetIPv6AddressCreator()
//}

//func (network *IPv6AddressNetwork) GetAddressCreator() AddressCreator {
//	return network.GetIPv6AddressCreator()
//}

func (network *IPv6AddressNetwork) GetLoopback() *IPAddress {
	//TODO use the creator
	return nil
}

//func (network *IPv6AddressNetwork) GetNetworkIPAddress(prefLen PrefixLen) *IPAddress {
//	return network.GetNetworkIPv6Address(prefLen).ToIPAddress()
//}

func (network *IPv6AddressNetwork) GetNetworkMask(prefLen BitCount) *IPAddress {
	return getMask(IPv6, zeroIPv6Seg.ToAddressDivision(), prefLen, network.subnetMasks, true, false)
	//return network.getNetworkIPv6Mask(prefLen, false).ToIPAddress()
}

func (network *IPv6AddressNetwork) GetPrefixedNetworkMask(prefLen BitCount) *IPAddress {
	return getMask(IPv6, zeroIPv6Seg.ToAddressDivision(), prefLen, network.subnetsMasksWithPrefix, true, true)
	//return network.getNetworkIPv6Mask(prefLen, true).ToIPAddress()
}

func (network *IPv6AddressNetwork) GetHostMask(prefLen BitCount) *IPAddress {
	return getMask(IPv6, zeroIPv6Seg.ToAddressDivision(), prefLen, network.hostMasks, false, false)
	//return network.getNetworkIPv4Mask(prefLen, false).ToIPAddress()
}

func (network *IPv6AddressNetwork) GetPrefixedHostMask(prefLen BitCount) *IPAddress {
	return getMask(IPv6, zeroIPv6Seg.ToAddressDivision(), prefLen, network.hostMasksWithPrefix, false, true)
	//return network.getNetworkIPv4Mask(prefLen, false).ToIPAddress()
}

//func (network *IPv6AddressNetwork) GetNetworkIPv6Address(prefLen PrefixLen) *IPv6Address {
//	//
//	return nil
//}

//func (network *IPv6AddressNetwork) getNetworkIPv6Mask(prefLen BitCount, withPrefixLength bool) *IPv6Address {
//	return getMask(IPv6, zeroIPv6Seg.ToAddressDivision(), prefLen, cache, true, withPrefixLength).ToIPv6Address()
//	//xxx
//	//func getMask(version IPVersion, zeroSeg *AddressDivision, networkPrefixLength BitCount, cache []*IPAddress,  network,  withPrefixLength bool) *IPAddress {
//	//return nil
//}

var _ IPAddressNetwork = &IPv6AddressNetwork{}

var DefaultIPv6Network = IPv6AddressNetwork{
	ipAddressNetwork: ipAddressNetwork{
		make([]*IPAddress, IPv6BitCount+1),
		make([]*IPAddress, IPv6BitCount+1),
		make([]*IPAddress, IPv6BitCount+1),
		make([]*IPAddress, IPv6BitCount+1),
	},
}

//
//
//
//
//

type IPv4AddressNetwork struct {
	ipAddressNetwork
	creator IPv4AddressCreator
}

func (network *IPv4AddressNetwork) GetIPv4AddressCreator() *IPv4AddressCreator {
	return &network.creator
}

//func (network *IPv4AddressNetwork) GetIPAddressCreator() IPAddressCreator {
//	return network.GetIPv4AddressCreator()
//}

//func (network *IPv4AddressNetwork) GetAddressCreator() AddressCreator {
//	return network.GetIPv4AddressCreator()
//}

func (network *IPv4AddressNetwork) GetLoopback() *IPAddress {
	//TODO
	return nil
}

//func (network *IPv4AddressNetwork) GetNetworkIPAddress(prefLen PrefixLen) *IPAddress {
//	return network.GetNetworkIPv4Address(prefLen).ToIPAddress()
//}
//func (network *IPv4AddressNetwork) GetNetworkIPv4Address(prefLen PrefixLen) *IPv4Address {
//	// get the ipv4 network address for a given prefix len, which is the all ones host (but for what address?)
//	return nil
//}

func (network *IPv4AddressNetwork) GetNetworkMask(prefLen BitCount) *IPAddress {
	return getMask(IPv4, zeroIPv4Seg.ToAddressDivision(), prefLen, network.subnetMasks, true, false)
	//return network.getNetworkIPv4Mask(prefLen, false).ToIPAddress()
}

func (network *IPv4AddressNetwork) GetPrefixedNetworkMask(prefLen BitCount) *IPAddress {
	return getMask(IPv4, zeroIPv4Seg.ToAddressDivision(), prefLen, network.subnetsMasksWithPrefix, true, true)
	//return network.getNetworkIPv4Mask(prefLen, true).ToIPAddress()
}

func (network *IPv4AddressNetwork) GetHostMask(prefLen BitCount) *IPAddress {
	return getMask(IPv4, zeroIPv4Seg.ToAddressDivision(), prefLen, network.hostMasks, false, false)
	//return network.getNetworkIPv4Mask(prefLen, false).ToIPAddress()
}

func (network *IPv4AddressNetwork) GetPrefixedHostMask(prefLen BitCount) *IPAddress {
	return getMask(IPv4, zeroIPv4Seg.ToAddressDivision(), prefLen, network.hostMasksWithPrefix, false, true)
	//return network.getNetworkIPv4Mask(prefLen, false).ToIPAddress()
}

//func (network *IPv4AddressNetwork) getNetworkIPv4Mask(prefLen BitCount, withPrefixLength bool) *IPv4Address {
//	return getMask(IPv4, zeroIPv4Seg.ToAddressDivision(), prefLen, cache, true, withPrefixLength).ToIPv4Address()
//	//return nil
//}

var _ IPAddressNetwork = &IPv4AddressNetwork{}

var DefaultIPv4Network = IPv4AddressNetwork{
	ipAddressNetwork: ipAddressNetwork{
		make([]*IPAddress, IPv4BitCount+1),
		make([]*IPAddress, IPv4BitCount+1),
		make([]*IPAddress, IPv4BitCount+1),
		make([]*IPAddress, IPv4BitCount+1),
	},
}

/*
func cacheStr(cachedString **string, stringer func() string) (str string) {
	cachedVal := *cachedString
	if cachedVal == nil {
		str = stringer()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(cachedString))
		atomic.StorePointer(dataLoc, unsafe.Pointer(&str))
	} else {
		str = *cachedVal
	}
	return
}
*/

var maskMutex sync.Mutex

func getMask(version IPVersion, zeroSeg *AddressDivision, networkPrefixLength BitCount, cache []*IPAddress, network, withPrefixLength bool) *IPAddress {
	bits := networkPrefixLength
	//IPVersion version = getIPVersion();
	addressBitLength := GetBitCount(version)
	if bits < 0 {
		bits = 0
	} else if bits > addressBitLength {
		bits = addressBitLength
	}
	cacheIndex := bits
	subnet := cache[cacheIndex]
	if subnet != nil {
		return subnet
	}

	maskMutex.Lock()
	subnet = cache[cacheIndex]
	if subnet != nil {
		maskMutex.Unlock()
		return subnet
	}
	//
	//
	//

	var onesSubnetIndex, zerosSubnetIndex int
	if network {
		onesSubnetIndex = int(addressBitLength)
		zerosSubnetIndex = 0
	} else {
		onesSubnetIndex = 0
		zerosSubnetIndex = int(addressBitLength)
	}
	onesSubnet := cache[onesSubnetIndex]
	zerosSubnet := cache[zerosSubnetIndex]
	segmentCount := GetSegmentCount(version)
	bitsPerSegment := GetBitsPerSegment(version)
	//bytesPerSegment := GetBytesPerSegment(version);
	//if(onesSubnet == nil || zerosSubnet == nil) {
	//synchronized(cache) {
	//onesSubnet = cache[onesSubnetIndex];
	maxSegmentValue := GetMaxSegmentValue(version)
	if onesSubnet == nil {
		//IPAddressCreator<T, ?, ?, S, ?> creator = getAddressCreator();
		newSegments := createSegmentArray(segmentCount)

		//if network && withPrefixLength {
		if withPrefixLength {
			if network {
				segment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, nil))
				//lastSegment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, getDivisionPrefixLength(bitsPerSegment, bitsPerSegment) /* bitsPerSegment */))
				lastSegment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, cacheBitCount(bitsPerSegment) /* bitsPerSegment */))
				lastIndex := len(newSegments) - 1
				fillDivs(newSegments[:lastIndex], segment)
				//S segment = creator.createSegment(maxSegmentValue, IPAddressSection.getSegmentPrefixLength(bitsPerSegment, addressBitLength) /* null */ );
				//Arrays.fill(newSegments, 0, newSegments.length - 1, segment);
				newSegments[lastIndex] = lastSegment
				onesSubnet = createIPAddress(createSection(newSegments, cacheBitCount(addressBitLength), version.toType(), 0), noZone)
				//onesSubnet = creator.createAddressInternal(newSegments, cacheBits(addressBitLength)); /* address creation */

			} else {
				segment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, cacheBits(0)))
				//newSegments[0] = firstSegment
				//fillDivs(newSegments[1:], segment)
				fillDivs(newSegments, segment)
				onesSubnet = createIPAddress(createSection(newSegments, cacheBits(0), version.toType(), 0), noZone)
			}
		} else {
			segment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, nil))
			//S segment = creator.createSegment(maxSegmentValue);
			//Arrays.fill(newSegments, segment);
			fillDivs(newSegments, segment)
			onesSubnet = createIPAddress(createSection(newSegments, nil, version.toType(), 0), noZone) /* address creation */
			//onesSubnet = creator.createAddressInternal(newSegments); /* address creation */
		}
		//initMaskCachedValues(onesSubnet.getSection(), network, withPrefixLength, networkAddress, addressBitLength, onesSubnetIndex, segmentCount, bitsPerSegment, bytesPerSegment);

		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache[onesSubnetIndex]))
		atomic.StorePointer(dataLoc, unsafe.Pointer(onesSubnet))

		//cache[onesSubnetIndex] = onesSubnet;
	}
	//zerosSubnet = cache[zerosSubnetIndex];
	if zerosSubnet == nil {
		//IPAddressCreator<T, ?, ?, S, ?> creator = getAddressCreator();
		newSegments := createSegmentArray(segmentCount)
		//S seg;
		if withPrefixLength {
			prefLen := cacheBits(0)
			if network {
				segment := createAddressDivision(zeroSeg.deriveNewSeg(0, prefLen))
				//seg = creator.createSegment(0, IPAddressSection.getSegmentPrefixLength(bitsPerSegment, 0) /* 0 */);
				fillDivs(newSegments, segment)
				//Arrays.fill(newSegments, seg);
				zerosSubnet = createIPAddress(createSection(newSegments, prefLen, version.toType(), 0), noZone)
				//zerosSubnet = creator.createAddressInternal(newSegments, prefLen); /* address creation */
				//if(getPrefixConfiguration().zeroHostsAreSubnets() && !networkAddress) {
				//	zerosSubnet = (T) zerosSubnet.getLower();
				//}
			} else {

				//segment := createAddressDivision(zeroSeg.deriveNewSeg(0, xxx))
				//lastSegment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, getDivisionPrefixLength(bitsPerSegment, bitsPerSegment) /* bitsPerSegment */))
				lastSegment := createAddressDivision(zeroSeg.deriveNewSeg(0, cacheBitCount(bitsPerSegment) /* bitsPerSegment */))
				lastIndex := len(newSegments) - 1
				fillDivs(newSegments[:lastIndex], zeroSeg)
				//S segment = creator.createSegment(maxSegmentValue, IPAddressSection.getSegmentPrefixLength(bitsPerSegment, addressBitLength) /* null */ );
				//Arrays.fill(newSegments, 0, newSegments.length - 1, segment);
				newSegments[lastIndex] = lastSegment
				zerosSubnet = createIPAddress(createSection(newSegments, cacheBitCount(addressBitLength), version.toType(), 0), noZone)
				//onesSubnet = creator.createAddressInternal(newSegments, cacheBits(addressBitLength)); /* address creation */

			}
		} else {
			segment := createAddressDivision(zeroSeg.deriveNewSeg(0, nil))
			fillDivs(newSegments, segment)
			zerosSubnet = createIPAddress(createSection(newSegments, nil, version.toType(), 0), noZone)
			//seg = creator.createSegment(0);
			//Arrays.fill(newSegments, seg);
			//zerosSubnet = creator.createAddressInternal(newSegments); /* address creation */
		}
		//initMaskCachedValues(zerosSubnet.getSection(), network, withPrefixLength, networkAddress, addressBitLength, zerosSubnetIndex, segmentCount, bitsPerSegment, bytesPerSegment);

		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache[zerosSubnetIndex]))
		atomic.StorePointer(dataLoc, unsafe.Pointer(zerosSubnet))

		//cache[zerosSubnetIndex] = zerosSubnet;
	}
	//}
	//}

	//synchronized(cache) {
	//subnet = cache[cacheIndex];
	//if(subnet == nil) {
	//BiFunction<T, Integer, S> segProducer = getSegmentProducer();
	prefix := bits
	onesSegment := onesSubnet.getDivision(0)
	zerosSegment := zerosSubnet.getDivision(0)
	//onesSegment := segProducer(onesSubnet, 1);
	//zerosSegment := segProducer(zerosSubnet, 1);
	//IPAddressCreator<T, ?, ?, S, ?> creator = getAddressCreator();

	//ArrayList<S> segmentList = new ArrayList<S>(segmentCount);
	newSegments := createSegmentArray(segmentCount)[:0]
	i := 0
	//for ; bits > 0; i++, bits -= bitsPerSegment {
	for ; bits > 0; i, bits = i+1, bits-bitsPerSegment {
		if bits <= bitsPerSegment {
			//S segment = null;
			var segment *AddressDivision

			//first do a check whether we have already created a segment like the one we need
			offset := ((bits - 1) % bitsPerSegment) + 1
			for j, entry := 0, offset; j < segmentCount; j, entry = j+1, entry+bitsPerSegment {
				//for j := 0, entry = offset; j < segmentCount; j++, entry += bitsPerSegment {
				if entry != cacheIndex { //we already know that the entry at cacheIndex is null
					prev := cache[entry]
					if prev != nil {
						segment = prev.getDivision(j)
						//segment = segProducer.apply(prev, j);
						break
					}
				}
			}

			//if none of the other addresses with a similar segment are created yet, we need a new segment.
			if segment == nil {
				//int networkMask = fullMask & (fullMask << (segmentBitSize - i));
				//int mask = getSegmentNetworkMask(bits);
				if network {
					mask := maxSegmentValue & (maxSegmentValue << (bitsPerSegment - bits))
					if withPrefixLength {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, getDivisionPrefixLength(bitsPerSegment, bits)))
						//segment = creator.createSegment(mask, IPAddressSection.getSegmentPrefixLength(bitsPerSegment, bits));
					} else {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, nil))
						//segment = creator.createSegment(mask);
					}
				} else {
					mask := maxSegmentValue & ^(maxSegmentValue << (bitsPerSegment - bits))
					if withPrefixLength {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, getDivisionPrefixLength(bitsPerSegment, bits)))
					} else {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, nil))
					}
					//segment = creator.createSegment(getSegmentHostMask(bits));
				}
			}
			//segmentList.add(segment);
			newSegments = append(newSegments, segment)
		} else {
			if network {
				newSegments = append(newSegments, onesSegment)
			} else {
				newSegments = append(newSegments, zerosSegment)
			}
			//segmentList.add(network ? onesSegment : zerosSegment);
		}
	}
	for ; i < segmentCount; i++ {
		if network {
			newSegments = append(newSegments, zerosSegment)
		} else {
			newSegments = append(newSegments, onesSegment)
		}
		//segmentList.add(network ? zerosSegment : onesSegment);
	}
	//S newSegments[] = creator.createSegmentArray(segmentList.size());
	//segmentList.toArray(newSegments);
	var prefLen PrefixLen
	if withPrefixLength {
		prefLen = cacheBitCount(prefix)
	}
	subnet = createIPAddress(createSection(newSegments, prefLen, version.toType(), 0), noZone)

	//if withPrefixLength {
	//	subnet = createIPAddress(createSection(newSegments, cacheBitCount(prefix), version.toType(), 0), noZone)
	//
	//	//subnet = creator.createAddressInternal(newSegments, cacheBits(prefix)); /* address creation */
	//	//if(getPrefixConfiguration().zeroHostsAreSubnets() && !networkAddress) {
	//	//	subnet = (T) subnet.getLower();
	//	//}
	//} else {
	//	subnet = createIPAddress(createSection(newSegments, nil, version.toType(), 0), noZone)
	//	//subnet = creator.createAddressInternal(newSegments); /* address creation */
	//}
	//initialize the cache fields since we know what they are now - they do not have to be calculated later
	//initMaskCachedValues(subnet.getSection(), network, withPrefixLength, networkAddress, addressBitLength, prefix, segmentCount, bitsPerSegment, bytesPerSegment);
	//cache[cacheIndex] = subnet; //last thing is to put into the cache - don't put it there before we are done with it
	//} // end subnet from cache is null

	dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache[cacheIndex]))
	atomic.StorePointer(dataLoc, unsafe.Pointer(subnet))

	//} //end synchronized

	//
	//
	maskMutex.Unlock()

	//} // end subnet from cache is null
	return subnet
}

type MACAddressNetwork struct {
	creator MACAddressCreator
}

func (network *MACAddressNetwork) GetMACAddressCreator() *MACAddressCreator {
	return &network.creator
}

//
//func (network *MACAddressNetwork) GetAddressCreator() AddressCreator {
//	return network.GetMACAddressCreator()
//}

var _ AddressNetwork = &MACAddressNetwork{}

var DefaultMACNetwork MACAddressNetwork