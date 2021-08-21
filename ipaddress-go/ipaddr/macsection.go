package ipaddr

import (
	"math/big"
	"unsafe"
)

func createMACSection(segments []*AddressDivision) *MACAddressSection {
	return &MACAddressSection{
		addressSectionInternal{
			addressDivisionGroupingInternal{
				addressDivisionGroupingBase: addressDivisionGroupingBase{
					divisions: standardDivArray{segments},
					addrType:  macType,
					cache: &valueCache{
						stringCache: stringCache{
							macStringCache: &macStringCache{},
						},
					},
				},
			},
		},
	}
}

// error returned for invalid segment count, nil sements, segments with invalid bit size, or inconsistent prefixes
func newMACSection(segments []*AddressDivision) (res *MACAddressSection, err AddressValueError) {
	segsLen := len(segments)
	if segsLen > ExtendedUniqueIdentifier64SegmentCount {
		err = &addressValueError{val: segsLen, addressError: addressError{key: "ipaddress.error.exceeds.size"}}
		return
	}
	res = createMACSection(segments)
	if err = res.initMultAndImplicitPrefLen(MACBitsPerSegment, true); err != nil {
		res = nil
		return
	}
	return
}

func NewMACSection(segments []*MACAddressSegment) (res *MACAddressSection, err AddressValueError) {
	res, err = newMACSection(cloneMACSegsToDivs(segments))
	return
}

func newMACSectionParsed(segments []*AddressDivision) (res *MACAddressSection) {
	res = createMACSection(segments)
	_ = res.initMultAndImplicitPrefLen(MACBitsPerSegment, false)
	return
}

func NewMACSectionFromBytes(bytes []byte, segmentCount int) (res *MACAddressSection, err AddressValueError) {
	if segmentCount < 0 {
		segmentCount = len(bytes)
	}
	expectedByteCount := segmentCount
	segments, err := toSegments(
		bytes,
		segmentCount,
		MACBytesPerSegment,
		MACBitsPerSegment,
		expectedByteCount,
		DefaultMACNetwork.getAddressCreator(),
		nil)
	if err == nil {
		// note prefix len is nil
		res = createMACSection(segments)
		if expectedByteCount == len(bytes) {
			bytes = cloneBytes(bytes)
			res.cache.bytesCache = &bytesCache{lowerBytes: bytes}
			if !res.isMultiple { // not a prefix block
				res.cache.bytesCache.upperBytes = bytes
			}
		}
	}
	return
}

func NewMACSectionFromUint64(bytes uint64, segmentCount int) (res *MACAddressSection) {
	if segmentCount < 0 {
		segmentCount = MediaAccessControlSegmentCount
	}
	segments := createSegmentsUint64(
		segmentCount,
		0,
		uint64(bytes),
		MACBytesPerSegment,
		MACBitsPerSegment,
		DefaultMACNetwork.getAddressCreator(),
		nil)
	// note prefix len is nil
	res = createMACSection(segments)
	return
}

func NewMACSectionFromVals(vals SegmentValueProvider, segmentCount int) (res *MACAddressSection) {
	res = NewMACSectionFromRange(vals, nil, segmentCount)
	return
}

func NewMACSectionFromRange(vals, upperVals SegmentValueProvider, segmentCount int) (res *MACAddressSection) {
	if segmentCount < 0 {
		segmentCount = 0
	}
	segments, isMultiple := createSegments(
		vals,
		upperVals,
		segmentCount,
		MACBitsPerSegment,
		DefaultMACNetwork.getAddressCreator(),
		nil)
	res = createMACSection(segments)
	if isMultiple {
		res.initImplicitPrefLen(MACBitsPerSegment)
		res.isMultiple = true
	}
	return
}

type MACAddressSection struct {
	addressSectionInternal
}

//func (section *MACAddressSection) IsExtended() bool {
//	return section.isExtended
//}

func (section *MACAddressSection) GetBitsPerSegment() BitCount {
	return MACBitsPerSegment
}

func (section *MACAddressSection) GetBytesPerSegment() int {
	return MACBytesPerSegment
}

func (section *MACAddressSection) GetCount() *big.Int {
	return section.cacheCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		}, section.GetSegmentCount(), 6, 0x7fffffffffffff)
	})
}

func (section *MACAddressSection) GetPrefixCount() *big.Int {
	return section.cachePrefixCount(func() *big.Int {
		return section.GetPrefixCountLen(*section.GetPrefixLen())
	})
}

func (section *MACAddressSection) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if prefixLen <= 0 {
		return bigOne()
	} else if bc := section.GetBitCount(); prefixLen >= bc {
		return section.GetCount()
	}
	networkSegmentIndex := getNetworkSegmentIndex(prefixLen, section.GetBytesPerSegment(), section.GetBitsPerSegment())
	hostSegmentIndex := getHostSegmentIndex(prefixLen, section.GetBytesPerSegment(), section.GetBitsPerSegment())
	return section.calcCount(func() *big.Int {
		return count(func(index int) uint64 {
			if (networkSegmentIndex == hostSegmentIndex) && index == networkSegmentIndex {
				segmentPrefixLength := getPrefixedSegmentPrefixLength(section.GetBitsPerSegment(), prefixLen, index)
				return getPrefixValueCount(section.GetSegment(index).ToAddressSegment(), *segmentPrefixLength)
			}
			return section.GetSegment(index).GetValueCount()
		}, networkSegmentIndex+1, 6, 0x7fffffffffffff)
	})
}

func (section *MACAddressSection) WithoutPrefixLen() *MACAddressSection {
	return section.withoutPrefixLen().ToMACAddressSection()
}

func (section *MACAddressSection) SetPrefixLen(prefixLen BitCount) *MACAddressSection {
	return section.setPrefixLen(prefixLen).ToMACAddressSection()
}

func (section *MACAddressSection) SetPrefixLenZeroed(prefixLen BitCount) (*MACAddressSection, IncompatibleAddressError) {
	res, err := section.setPrefixLenZeroed(prefixLen)
	return res.ToMACAddressSection(), err
}

func (section *MACAddressSection) AssignPrefixForSingleBlock() *MACAddressSection {
	return section.assignPrefixForSingleBlock().ToMACAddressSection()
}

func (section *MACAddressSection) AssignMinPrefixForBlock() *MACAddressSection {
	return section.assignMinPrefixForBlock().ToMACAddressSection()
}

func (section *MACAddressSection) GetSegment(index int) *MACAddressSegment {
	return section.getDivision(index).ToMACAddressSegment()
}

func (section *MACAddressSection) ToAddressSection() *AddressSection {
	return (*AddressSection)(unsafe.Pointer(section))
}

// Gets the subsection from the series starting from the given index
// The first segment is at index 0.
func (section *MACAddressSection) GetTrailingSection(index int) *MACAddressSection {
	return section.GetSubSection(index, section.GetSegmentCount())
}

//// Gets the subsection from the series starting from the given index and ending just before the give endIndex
//// The first segment is at index 0.
func (section *MACAddressSection) GetSubSection(index, endIndex int) *MACAddressSection {
	return section.getSubSection(index, endIndex).ToMACAddressSection()
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (section *MACAddressSection) CopySubSegments(start, end int, segs []*MACAddressSegment) (count int) {
	return section.visitSubSegments(start, end, func(index int, div *AddressDivision) bool { segs[index] = div.ToMACAddressSegment(); return false }, len(segs))
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (section *MACAddressSection) CopySegments(segs []*MACAddressSegment) (count int) {
	return section.visitSegments(func(index int, div *AddressDivision) bool { segs[index] = div.ToMACAddressSegment(); return false }, len(segs))
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (section *MACAddressSection) GetSegments() (res []*MACAddressSegment) {
	res = make([]*MACAddressSegment, section.GetSegmentCount())
	section.CopySegments(res)
	return
}

func (section *MACAddressSection) GetLower() *MACAddressSection {
	return section.getLower().ToMACAddressSection()
}

func (section *MACAddressSection) GetUpper() *MACAddressSection {
	return section.getUpper().ToMACAddressSection()
}

func (section *MACAddressSection) Uint64Value() uint64 {
	return section.getLongValue(true)
}

func (section *MACAddressSection) UpperUint64Value() uint64 {
	return section.getLongValue(false)
}

func (section *MACAddressSection) getLongValue(lower bool) (result uint64) {
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return
	}
	seg := section.GetSegment(0)
	if lower {
		result = uint64(seg.GetSegmentValue())
	} else {
		result = uint64(seg.GetUpperSegmentValue())
	}
	bitsPerSegment := section.GetBitsPerSegment()
	for i := 1; i < segCount; i++ {
		result = (result << uint(bitsPerSegment))
		seg = section.GetSegment(i)
		if lower {
			result |= uint64(seg.GetSegmentValue())
		} else {
			result |= uint64(seg.GetUpperSegmentValue())
		}
	}
	return
}

func (section *MACAddressSection) ToPrefixBlock() *MACAddressSection {
	return section.toPrefixBlock().ToMACAddressSection()
}

func (section *MACAddressSection) ToPrefixBlockLen(prefLen BitCount) *MACAddressSection {
	return section.toPrefixBlockLen(prefLen).ToMACAddressSection()
}

func (section *MACAddressSection) ToBlock(segmentIndex int, lower, upper SegInt) *MACAddressSection {
	return section.toBlock(segmentIndex, lower, upper).ToMACAddressSection()
}

func (section *MACAddressSection) Iterator() MACSectionIterator {
	return macSectionIterator{section.sectionIterator(nil)}
}

func (section *MACAddressSection) PrefixIterator() MACSectionIterator {
	return macSectionIterator{section.prefixIterator(false)}
}

func (section *MACAddressSection) PrefixBlockIterator() MACSectionIterator {
	return macSectionIterator{section.prefixIterator(true)}
}

func (section *MACAddressSection) IncrementBoundary(increment int64) *MACAddressSection {
	return section.incrementBoundary(increment).ToMACAddressSection()

}

func getMacMaxValueLong(segmentCount int) uint64 {
	return macMaxValues[segmentCount]
}

var macMaxValues = []uint64{
	0,
	MACMaxValuePerSegment,
	0xffff,
	0xffffff,
	0xffffffff,
	0xffffffffff,
	0xffffffffffff,
	0xffffffffffffff,
	0xffffffffffffffff}

func (section *MACAddressSection) Increment(incrementVal int64) *MACAddressSection {
	if incrementVal == 0 && !section.IsMultiple() {
		return section
	}
	segCount := section.GetSegmentCount()
	lowerValue := section.Uint64Value()
	upperValue := section.UpperUint64Value()
	count := section.GetCount()
	countMinus1 := count.Sub(count, bigOneConst()).Uint64()
	isOverflow := checkOverflow(incrementVal, lowerValue, upperValue, countMinus1, getMacMaxValueLong(segCount))
	if isOverflow {
		return nil
	}
	return increment(
		section.ToAddressSection(),
		incrementVal,
		DefaultMACNetwork.getAddressCreator(),
		countMinus1,
		section.Uint64Value(),
		section.UpperUint64Value(),
		section.addressSectionInternal.getLower,
		section.addressSectionInternal.getUpper,
		section.GetPrefixLen()).ToMACAddressSection()
	//			}
	//			BigInteger lowerValue = getValue();
	//			BigInteger upperValue = getUpperValue();
	//			BigInteger count = getCount();
	//			BigInteger bigIncrement = BigInteger.valueOf(increment);
	//			checkOverflow(increment, bigIncrement, lowerValue, upperValue, count, () -> getMaxValue(getSegmentCount()));
	//			Integer prefixLength = getNetwork().getPrefixConfiguration().allPrefixedAddressesAreSubnets() ? null : getPrefixLength();
	//			MACAddressSection result = fastIncrement(
	//					this,
	//					increment,
	//					getIPAddressCreator(),
	//					this::getLower,
	//					this::getUpper,
	//					prefixLength);
	//			if(result != null) {
	//				return result;
	//			}
	//			return increment(
	//					this,
	//					increment,
	//					bigIncrement,
	//					getIPAddressCreator(),
	//					this::getLower,
	//					this::getUpper,
	//					getNetwork().getPrefixConfiguration().allPrefixedAddressesAreSubnets() ? null : getPrefixLength());
	//
	//*/
	//return nil
}

func (section *MACAddressSection) ReverseBits(perByte bool) (*MACAddressSection, IncompatibleAddressError) {
	res, err := section.reverseBits(perByte)
	return res.ToMACAddressSection(), err
}

func (section *MACAddressSection) ReverseBytes() *MACAddressSection {
	return section.ReverseSegments()
}

//func (section *MACAddressSection) ReverseBytesPerSegment() *MACAddressSection {
//	if !section.IsPrefixed() {
//		return section
//	}
//	return section.WithoutPrefixLen()
//}

func (section *MACAddressSection) ReverseSegments() *MACAddressSection {
	if section.GetSegmentCount() <= 1 {
		if section.IsPrefixed() {
			return section.WithoutPrefixLen()
		}
		return section
	}
	res, _ := section.reverseSegments(
		func(i int) (*AddressSegment, IncompatibleAddressError) {
			return section.GetSegment(i).ToAddressSegment(), nil
		},
	)
	return res.ToMACAddressSection()
}

func (section *MACAddressSection) Append(other *MACAddressSection) *MACAddressSection {
	count := section.GetSegmentCount()
	return section.ReplaceLen(count, count, other, 0, other.GetSegmentCount())
}

func (section *MACAddressSection) Insert(index int, other *MACAddressSection) *MACAddressSection {
	return section.ReplaceLen(index, index, other, 0, other.GetSegmentCount())
}

// Replace replaces the segments of this section starting at the given index with the given replacement segments
func (section *MACAddressSection) Replace(index int, replacement *MACAddressSection) *MACAddressSection {
	return section.ReplaceLen(index, index+replacement.GetSegmentCount(), replacement, 0, replacement.GetSegmentCount())
}

// ReplaceLen replaces segments starting from startIndex and ending before endIndex with the segments starting at replacementStartIndex and
// ending before replacementEndIndex from the replacement section
func (section *MACAddressSection) ReplaceLen(startIndex, endIndex int, replacement *MACAddressSection, replacementStartIndex, replacementEndIndex int) *MACAddressSection {
	return section.replaceLen(startIndex, endIndex, replacement.ToAddressSection(), replacementStartIndex, replacementEndIndex, macBitsToSegmentBitshift).ToMACAddressSection()
}

var (
	canonicalWildcards = new(WildcardsBuilder).SetRangeSeparator(MacDashedSegmentRangeSeparatorStr).SetWildcard(SegmentWildcardStr).ToWildcards()

	//macHexParams         = NewMACStringOptionsBuilder().SetHasSeparator(false).SetExpandedSegments(true).ToOptions()
	//macHexPrefixedParams = NewMACStringOptionsBuilder().SetHasSeparator(false).SetExpandedSegments(true).SetAddressLabel(HexPrefix).ToOptions()
	macNormalizedParams  = NewMACStringOptionsBuilder().SetExpandedSegments(true).ToOptions()
	macCanonicalParams   = NewMACStringOptionsBuilder().SetSeparator(MACDashSegmentSeparator).SetExpandedSegments(true).SetWildcards(canonicalWildcards).ToOptions()
	macCompressedParams  = NewMACStringOptionsBuilder().ToOptions()
	dottedParams         = NewMACStringOptionsBuilder().SetSeparator(MacDottedSegmentSeparator).SetExpandedSegments(true).ToOptions()
	spaceDelimitedParams = NewMACStringOptionsBuilder().SetSeparator(MacSpaceSegmentSeparator).SetExpandedSegments(true).ToOptions()
)

// ToCanonicalString produces a canonical string.
//
//If this section has a prefix length, it will be included in the string.
func (section *MACAddressSection) ToCanonicalString() string {
	return cacheStr(&section.getStringCache().canonicalString,
		func() string {
			return section.ToCustomString(macCanonicalParams)
		})
}

func (section *MACAddressSection) ToNormalizedString() string {
	cch := section.getStringCache()
	strp := &cch.normalizedMACString
	return cacheStr(strp,
		func() string {
			return section.ToCustomString(macNormalizedParams)
		})
}

func (section *MACAddressSection) ToCompressedString() string {
	return cacheStr(&section.getStringCache().compressedMACString,
		func() string {
			return section.ToCustomString(macCompressedParams)
		})
}

// ToDottedString produces the dotted hexadecimal format aaaa.bbbb.cccc
func (section *MACAddressSection) ToDottedString() (string, IncompatibleAddressError) {
	return cacheStrErr(&section.getStringCache().dottedString,
		func() (string, IncompatibleAddressError) {
			dottedGrouping, err := section.GetDottedGrouping()
			if err != nil {
				return "", err
			}
			//getStringCache().dottedString = result = toNormalizedString(MACStringCache.dottedParams, dottedGrouping);
			//return section.toNormalizedOptsString(dottedParams)
			return toNormalizedString(dottedParams, dottedGrouping), nil

			//return ""
		})
}

func (section *MACAddressSection) GetDottedGrouping() (AddressDivisionSeries, IncompatibleAddressError) {
	//start := section.addressSegmentIndex
	segmentCount := section.GetSegmentCount()
	var newSegs []*AddressDivision
	newSegmentBitCount := section.GetBitsPerSegment() << 1
	var segIndex, newSegIndex int
	//if (start & 1) == 0 {
	newSegmentCount := (segmentCount + 1) >> 1
	newSegs = make([]*AddressDivision, newSegmentCount)
	//newSegIndex = segIndex = 0;
	//} else {
	//	newSegmentCount := (segmentCount >> 1) + 1
	//	newSegs = make([]*AddressDivision, newSegmentCount)
	//	segment := section.GetSegment(0)
	//
	//	//func NewDivision(val DivInt, bitCount BitCount, defaultRadix int) *AddressDivision {
	//	//	return NewRangePrefixDivision(val, val, nil, bitCount, defaultRadix)
	//	//}
	//	//
	//	//func NewRangeDivision(val, upperVal DivInt, bitCount BitCount, defaultRadix int) *AddressDivision {
	//	//	return NewRangePrefixDivision(val, upperVal, nil, bitCount, defaultRadix)
	//	//}
	//	//
	//	//func NewPrefixDivision(val DivInt, prefixLen PrefixLen, bitCount BitCount, defaultRadix int) *AddressDivision {
	//	//	return NewRangePrefixDivision(val, val, prefixLen, bitCount, defaultRadix)
	//	//}
	//	vals := NewRangeDivision(segment.getDivisionValue(), segment.getUpperDivisionValue(), newSegmentBitCount, MACDefaultTextualRadix)
	//
	//	//vals := &bitsDivisionVals{
	//	//	value:      segment.getDivisionValue(),
	//	//	upperValue: segment.getUpperDivisionValue(),
	//	//	bitCount:   newSegmentBitCount,
	//	//	radix:      MACDefaultTextualRadix,
	//	//	//joinedCount: joinCount,
	//	//	//prefixLen:   nil,
	//	//}
	//	newSegs[0] = createAddressDivision(vals)
	//
	//	//newSegs[0] = new AddressBitsDivision(segment.getSegmentValue(),
	//	//	segment.getUpperSegmentValue(),
	//	//	newSegmentBitCount, MACDefaultTextualRadix);
	//	newSegIndex = 1
	//	segIndex = 1
	//}
	bitsPerSeg := section.GetBitsPerSegment()
	for segIndex+1 < segmentCount {
		segment1 := section.GetSegment(segIndex)
		segIndex++
		segment2 := section.GetSegment(segIndex)
		segIndex++
		if segment1.isMultiple() && !segment2.IsFullRange() {
			return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.invalid.joined.ranges"}}
			//throw new IncompatibleAddressError(segment1, segIndex - 2, segment2, segIndex - 1, "ipaddress.error.invalid.joined.ranges");
		}
		val := (segment1.GetSegmentValue() << uint(bitsPerSeg)) | segment2.GetSegmentValue()
		upperVal := (segment1.GetUpperSegmentValue() << uint(bitsPerSeg)) | segment2.GetUpperSegmentValue()
		vals := NewRangeDivision(DivInt(val), DivInt(upperVal), newSegmentBitCount, MACDefaultTextualRadix)

		//vals := &bitsDivisionVals{
		//	value:      DivInt((segment1.GetSegmentValue() << bitsPerSeg) | segment2.GetSegmentValue()),
		//	upperValue: DivInt((segment1.GetUpperSegmentValue() << bitsPerSeg) | segment2.GetUpperSegmentValue()),
		//	bitCount:   newSegmentBitCount,
		//	radix:      MACDefaultTextualRadix,
		//	//joinedCount: joinCount,
		//	//prefixLen:   nil,
		//}
		newSegs[newSegIndex] = createAddressDivision(vals)
		newSegIndex++
		//AddressDivision newSeg = new AddressBitsDivision(
		//		(segment1.GetSegmentValue() << getBitsPerSegment()) | segment2.GetSegmentValue(),
		//		(segment1.GetUpperSegmentValue() << getBitsPerSegment()) | segment2.GetUpperSegmentValue(),
		//		newSegmentBitCount,
		//		MACDefaultTextualRadix);
		//newSegs[newSegIndex++] = newSeg;
	}
	if segIndex < segmentCount {
		segment := section.GetSegment(segIndex)
		val := segment.GetSegmentValue() << uint(bitsPerSeg)
		upperVal := segment.GetUpperSegmentValue() << uint(bitsPerSeg)
		vals := NewRangeDivision(DivInt(val), DivInt(upperVal), newSegmentBitCount, MACDefaultTextualRadix)
		//vals := &bitsDivisionVals{
		//	value:      DivInt(segment.GetSegmentValue() << bitsPerSeg),
		//	upperValue: DivInt(segment.GetUpperSegmentValue() << bitsPerSeg),
		//	bitCount:   newSegmentBitCount,
		//	radix:      MACDefaultTextualRadix,
		//	//joinedCount: joinCount,
		//	//prefixLen:   nil,
		//}
		newSegs[newSegIndex] = createAddressDivision(vals)
		//			newSegs[newSegIndex] = new AddressBitsDivision(
		//					segment.getSegmentValue() << bitsPerSeg,
		//					segment.getUpperSegmentValue() << bitsPerSeg,
		//					newSegmentBitCount,
		//MACDefaultTextualRadix);
	}
	grouping := createInitializedGrouping(newSegs, section.GetPrefixLen(), zeroType)
	return grouping, nil
	//AddressDivisionGrouping dottedGrouping;
	//if(cachedPrefixLength == null) {
	//	dottedGrouping = new AddressDivisionGrouping(newSegs);
	//} else {
	//	Integer prefLength = cachedPrefixLength;
	//	dottedGrouping = new AddressDivisionGrouping(newSegs) {{
	//		cachedPrefixLength = prefLength;
	//	}};
	//}
	//return dottedGrouping;
}

// ToSpaceDelimitedString produces a string delimited by spaces: aa bb cc dd ee ff
func (section *MACAddressSection) ToSpaceDelimitedString() string {
	return cacheStr(&section.getStringCache().spaceDelimitedString,
		func() string {
			return section.ToCustomString(spaceDelimitedParams)
		})
}

func (section *MACAddressSection) ToDashedString() string {
	return section.ToCanonicalString()
}

func (section *MACAddressSection) ToColonDelimitedString() string {
	return section.ToNormalizedString()
}
