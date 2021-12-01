package ipaddr

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"unsafe"
)

func createGrouping(divs []*AddressDivision, prefixLength PrefixLen, addrType addrType) *AddressDivisionGrouping {
	grouping := &AddressDivisionGrouping{
		addressDivisionGroupingInternal{
			addressDivisionGroupingBase: addressDivisionGroupingBase{
				divisions:    standardDivArray{divs},
				prefixLength: prefixLength,
				addrType:     addrType,
				cache:        &valueCache{},
			},
		},
	}
	assignStringCache(&grouping.addressDivisionGroupingBase, addrType)
	return grouping
}

func createGroupingMultiple(divs []*AddressDivision, prefixLength PrefixLen, isMultiple bool) *AddressDivisionGrouping {
	result := createGrouping(divs, prefixLength, zeroType)
	result.isMult = isMultiple
	return result
}

func createInitializedGrouping(divs []*AddressDivision, prefixLength PrefixLen) *AddressDivisionGrouping {
	result := createGrouping(divs, prefixLength, zeroType)
	result.initMultiple() // assigns isMult
	return result
}

// Creates an arbitrary grouping of divisions.
// To create address sections or addresses, use the constructors that are specific to the address version or type.
// The AddressDivision instances can be created with the NewDivision, NewRangeDivision, NewPrefixDivision or NewRangePrefixDivision functions.
func NewDivisionGrouping(divs []*AddressDivision, prefixLength PrefixLen) *AddressDivisionGrouping {
	return createInitializedGrouping(divs, prefixLength)
}

var (
	emptyBytes = []byte{}
)

type addressDivisionGroupingInternal struct {
	addressDivisionGroupingBase

	// get rid of addressSegmentIndex and isExtended
	// You just don't need positionality from sections.
	// Being mixed or converting to IPv6 from MACSize are properties of the address.
	// isExtended really only used for IPv6/MACSize conversion.
	// addressSegmentindex really only used for mixed
	// Both of those are really "address-level" concepts.
	//
	// TODO LATER refactor to support infiniband, which will involve multiple types.
	// But that will be a joint effort with Java and will wait to later.

	// The index of the containing address where this section starts, only used by IPv6 where we trach the "IPv4-embedded" part of an address section
	//addressSegmentIndex int8

	//isExtended bool
}

func createSegmentArray(length int) []*AddressDivision {
	return make([]*AddressDivision, length)
}

func (grouping *addressDivisionGroupingInternal) initMultiple() {
	divCount := grouping.getDivisionCount()
	for i := divCount - 1; i >= 0; i-- {
		div := grouping.getDivision(i)
		if div.isMultiple() {
			grouping.isMult = true
			return
		}
	}
	return
}

// getDivision returns the division or panics if the index is negative or too large
func (grouping *addressDivisionGroupingInternal) getDivision(index int) *AddressDivision {
	divsArray := grouping.divisions
	if divsArray != nil {
		return divsArray.(standardDivArray).divisions[index]
	}
	panic("invalid index") // must be consistent with above code which panics with invalid index
}

// getDivisionsInternal returns the divisions slice, only to be used internally
func (grouping *addressDivisionGroupingInternal) getDivisionsInternal() []*AddressDivision {
	divsArray := grouping.divisions
	if divsArray != nil {
		return divsArray.(standardDivArray).getDivisions()
	}
	return nil
}

func (grouping *addressDivisionGroupingInternal) getDivisionCount() int {
	divsArray := grouping.divisions
	if divsArray != nil {
		return divsArray.(standardDivArray).getDivisionCount()
	}
	return 0
}

func adjust1To1Indices(sourceStart, sourceEnd, sourceCount, targetStart, targetCount int) (newSourceStart, newSourceEnd, newTargetStart int) {
	//targetIndex := 0
	if sourceStart < 0 {
		targetStart -= sourceStart
		sourceStart = 0
	}
	// how many to copy?
	if sourceEnd > sourceCount { // end index exceeds available
		sourceEnd = sourceCount
	}
	calcCount := sourceEnd - sourceStart
	if calcCount <= 0 { // end index below start index
		return sourceStart, sourceStart, targetStart
	}
	// if not enough space in target, adjust count and end
	if space := targetCount - targetStart; calcCount > space {
		if space <= 0 {
			return sourceStart, sourceStart, targetStart
		}
		sourceEnd = sourceStart + space
	}
	return sourceStart, sourceEnd, targetStart
}

func adjustIndices(
	startIndex, endIndex, sourceCount,
	replacementStartIndex, replacementEndIndex, replacementSegmentCount int) (int, int, int, int) {
	//segmentCount := section.GetSegmentCount()
	if startIndex < 0 {
		startIndex = 0
	} else if startIndex > sourceCount {
		startIndex = sourceCount
	}
	if endIndex < startIndex {
		endIndex = startIndex
	} else if endIndex > sourceCount {
		endIndex = sourceCount
	}
	if replacementStartIndex < 0 {
		replacementStartIndex = 0
	} else if replacementStartIndex > replacementSegmentCount {
		replacementStartIndex = replacementSegmentCount
	}
	if replacementEndIndex < replacementStartIndex {
		replacementEndIndex = replacementStartIndex
	} else if replacementEndIndex > replacementSegmentCount {
		replacementEndIndex = replacementSegmentCount
	}
	return startIndex, endIndex, replacementStartIndex, replacementEndIndex
}

func (grouping *addressDivisionGroupingInternal) visitDivisions(target func(index int, div *AddressDivision) bool, targetLen int) (count int) {
	if grouping.hasNoDivisions() {
		return
	}
	count = grouping.GetDivisionCount()
	if count > targetLen {
		count = targetLen
	}
	for start := 0; start < count; start++ {
		if target(start, grouping.getDivision(start)) {
			break
		}
	}
	return
}

func (grouping *addressDivisionGroupingInternal) visitSubDivisions(start, end int, target func(index int, div *AddressDivision) (stop bool), targetLen int) (count int) {
	if grouping.hasNoDivisions() {
		return
	}
	targetIndex := 0
	start, end, targetIndex = adjust1To1Indices(start, end, grouping.GetDivisionCount(), targetIndex, targetLen)

	// now iterate start to end
	index := start
	for index < end {
		exitEarly := target(targetIndex, grouping.getDivision(index))
		index++
		if exitEarly {
			break
		}
		targetIndex++
	}
	return index - start
}

// copySubDivisions copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (grouping *addressDivisionGroupingInternal) copySubDivisions(start, end int, divs []*AddressDivision) (count int) {
	//return grouping.visitSubDivisions(start, end, func(index int, div *AddressDivision) bool { divs[index] = div; return false }, len(divs))
	//divsArray := grouping.divisions
	//if divsArray != nil {
	//	return divsArray.(standardDivArray).copySubDivisions(start, end, divs)
	//}
	divsArray := grouping.divisions
	if divsArray != nil {
		targetIndex := 0
		start, end, targetIndex = adjust1To1Indices(start, end, grouping.GetDivisionCount(), targetIndex, len(divs))
		//return copy(grouping.divs,divsArray[start:end])
		return divsArray.(standardDivArray).copySubDivisions(start, end, divs)
		//xxxx
	}
	return
}

// copyDivisions copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (grouping *addressDivisionGroupingInternal) copyDivisions(divs []*AddressDivision) (count int) {
	//return grouping.visitDivisions(func(index int, div *AddressDivision) bool { divs[index] = div; return false }, len(divs))
	divsArray := grouping.divisions
	if divsArray != nil {
		return divsArray.(standardDivArray).copyDivisions(divs)
	}
	return
}

func (grouping *addressDivisionGroupingInternal) getSubDivisions(start, end int) []*AddressDivision {
	divsArray := grouping.divisions
	if divsArray != nil {
		return divsArray.(standardDivArray).getSubDivisions(start, end)
	} else if start != 0 || end != 0 {
		panic("invalid subslice")
	}
	return make([]*AddressDivision, 0)
}

func (grouping *addressDivisionGroupingInternal) isAddressSection() bool {
	return grouping != nil && grouping.matchesAddrSectionType()
}

//func (grouping *addressDivisionGroupingInternal) isAddressSection() bool {
//	if grouping == nil {
//		return false
//	}
//	if grouping.matchesAddrSectionType() {
//		return true
//	}
//	var bitCount BitCount
//	count := grouping.GetDivisionCount()
//	// all divisions must be equal size and have an exact number of bytes
//	for i := 0; i < count; i++ {
//		div := grouping.getDivision(i)
//		if i == 0 {
//			bitCount = div.GetBitCount()
//			if bitCount%8 != 0 || bitCount > SegIntSize {
//				return false
//			}
//		} else if bitCount != div.GetBitCount() {
//			return false
//		}
//	}
//	return true
//}

//func (grouping *addressDivisionGroupingInternal) CompareSize(other AddressDivisionSeries) int { // the getCount() is optimized which is why we do not defer to the method in addressDivisionGroupingBase
func (grouping *addressDivisionGroupingInternal) compareSize(other StandardDivisionGroupingType) int { // the getCount() is optimized which is why we do not defer to the method in addressDivisionGroupingBase
	if other == nil || other.ToAddressDivisionGrouping() == nil {
		// our size is 1 or greater, other 0
		return 1
	}
	if !grouping.isMultiple() {
		if other.IsMultiple() {
			return -1
		}
		return 0
	} else if !other.IsMultiple() {
		return 1
	}
	return grouping.getCount().CmpAbs(other.GetCount())
}

func (grouping *addressDivisionGroupingInternal) getCount() *big.Int {
	if !grouping.isMultiple() {
		return bigOne()
	} else if section := grouping.toAddressSection(); section != nil {
		return section.GetCount()
	}
	return grouping.addressDivisionGroupingBase.getCount()
}

func (grouping *addressDivisionGroupingInternal) GetPrefixCount() *big.Int {
	if section := grouping.toAddressSection(); section != nil {
		return section.GetPrefixCount()
	}
	return grouping.addressDivisionGroupingBase.GetPrefixCount()
}

func (grouping *addressDivisionGroupingInternal) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if section := grouping.toAddressSection(); section != nil {
		return section.GetPrefixCountLen(prefixLen)
	}
	return grouping.addressDivisionGroupingBase.GetPrefixCountLen(prefixLen)
}

func (grouping *addressDivisionGroupingInternal) getDivisionStrings() []string {
	if grouping.hasNoDivisions() {
		return []string{}
	}
	result := make([]string, grouping.GetDivisionCount())
	for i := range result {
		result[i] = grouping.getDivision(i).String()
	}
	return result
}

func (grouping *addressDivisionGroupingInternal) getSegmentStrings() []string {
	if grouping.hasNoDivisions() {
		return []string{}
	}
	result := make([]string, grouping.GetDivisionCount())
	for i := range result {
		result[i] = grouping.getDivision(i).GetWildcardString()
	}
	return result
}

func (grouping *addressDivisionGroupingInternal) toAddressDivisionGrouping() *AddressDivisionGrouping {
	return (*AddressDivisionGrouping)(unsafe.Pointer(grouping))
}

func (grouping *addressDivisionGroupingInternal) toAddressSection() *AddressSection {
	return grouping.toAddressDivisionGrouping().ToAddressSection()
}

func (grouping *addressDivisionGroupingInternal) matchesIPv6AddressType() bool {
	return grouping.getAddrType().isIPv6() // no need to check segment count because addresses cannot be constructed with incorrect segment count
}

func (grouping *addressDivisionGroupingInternal) matchesIPv4AddressType() bool {
	return grouping.getAddrType().isIPv4() // no need to check segment count because addresses cannot be constructed with incorrect segment count
}

func (grouping *addressDivisionGroupingInternal) matchesIPAddressType() bool {
	return grouping.matchesIPSectionType() // no need to check segment count because addresses cannot be constructed with incorrect segment count (note the zero IPAddress has zero segments)
}

func (grouping *addressSectionInternal) matchesMACAddressType() bool {
	return grouping.getAddrType().isMAC()
}

func (grouping *addressDivisionGroupingInternal) matchesAddrSectionType() bool {
	return !grouping.getAddrType().isNil() || grouping.hasNoDivisions()
}

func (grouping *addressDivisionGroupingInternal) matchesIPv6SectionType() bool {
	addrType := grouping.getAddrType()
	return addrType.isIPv6() || (addrType.isNil() && grouping.hasNoDivisions())
}

func (grouping *addressDivisionGroupingInternal) matchesIPv6v4MixedGroupingType() bool {
	addrType := grouping.getAddrType()
	return addrType.isIPv6v4Mixed() || (addrType.isNil() && grouping.hasNoDivisions())
}

func (grouping *addressDivisionGroupingInternal) matchesIPv4SectionType() bool {
	addrType := grouping.getAddrType()
	return addrType.isIPv4() || (addrType.isNil() && grouping.hasNoDivisions())
}

func (grouping *addressDivisionGroupingInternal) matchesIPSectionType() bool {
	addrType := grouping.getAddrType()
	return addrType.isIP() || (addrType.isNil() && grouping.hasNoDivisions())
}

func (grouping *addressDivisionGroupingInternal) matchesMACSectionType() bool {
	addrType := grouping.getAddrType()
	return addrType.isMAC() || (addrType.isNil() && grouping.hasNoDivisions())
}

func (grouping *addressDivisionGroupingInternal) initDivs() *addressDivisionGroupingInternal {
	if grouping.divisions == nil {
		return &zeroSection.addressDivisionGroupingInternal
	}
	return grouping
}

func (grouping *addressDivisionGroupingInternal) toString() string {
	if sect := grouping.toAddressSection(); sect != nil {
		return sect.ToNormalizedString()
	}
	return fmt.Sprintf("%v", grouping.initDivs().divisions) //TODO see how I print with Format, make sure this is consistent
}

func (grouping addressDivisionGroupingInternal) Format(state fmt.State, verb rune) {
	if sect := grouping.toAddressSection(); sect != nil {
		sect.Format(state, verb)
		return
	}
	grouping.defaultFormat(state, verb)
}

func (grouping addressDivisionGroupingInternal) defaultFormat(state fmt.State, verb rune) {
	//state.Write([]byte(fmt.Sprintf("%"+string(verb), grouping.initDivs().divisions)))
	s := flagsFromState(state, verb)
	state.Write([]byte(fmt.Sprintf(s, grouping.initDivs().divisions.(standardDivArray).divisions)))
	//return fmt.Sprintf(s, grouping.initDivs().divisions.(standardDivArray).divisions)
	//return fmt.Sprintf("%v", grouping.initDivs().divisions)
	//	Line 393 grouping.go
	//
	//		see fmtBytes and fmtBytes in fmt printf.go
	//		xxx can I pass on to the slice?  How do I do that? xxxx
	//	xxx most likely you just reconstruct the string as it is described xxxx
	//	xxxx but you know, it really only does cool stuff with []byte
	//	xxxx which is useless when we are a range and they can use GetBytes to do what they want
	//	xxxx but even for stringer, it seems that fmtSx does all kinds of cool shit
	//	xxxx
	//	xxxx
	//	for one thing, we can just put the verb there directly, the verb is most of the formatting
	//	seems bytes uses sharpV with v and d
	//	shit, it seems there is a lot of shit to cover
	//	we need to just reproduce
	//
	//	java says it is:
	//	%[flags][width][.precision]conversion-character
	//	https://docs.oracle.com/javase/7/docs/api/java/util/Formatter.html
	//	these seem to apply: '-' '#' ' ' '0' '+'
	//according to fmt_test, the order is:
	//	'# '
	//	'#+'
	//	'-#20.8X'
	//	' +.68d'
	//	'#-014.6U'
	//	'-020E'
	//	'#-06v'
	//	'#+-6d'
	//	'# 02x'
	//	'# -010X'
	//	'# -010X'
	//	' +07.2f'
	//	'+07.2f'
	//	'-05.1f'
	//	reorderTests seems crazy, skip reorder
	//
	//	I think there is no order for the flags
	//	but go '# +-0' then width then precision
	//
	//	then there is this:
	//
	//func (flagPrinter) Format(f State, c rune) {
	//	s := "%"
	//	for i := 0; i < 128; i++ {
	//		if f.Flag(i) {
	//			s += string(i)
	//		}
	//	}
	//	if w, ok := f.Width(); ok {
	//		s += Sprintf("%d", w)
	//	}
	//	if p, ok := f.Precision(); ok {
	//		s += Sprintf(".%d", p)
	//	}
	//	s += string(c)
	//	io.WriteString(f, "["+s+"]") this seems to be a flag printer thing, just ignore it
	//}

	/*
		// For the formats %+v %#v, we set the plusV/sharpV flags
			// and clear the plus/sharp flags since %+v and %#v are in effect
			// different, flagless formats set at the top level.
			plusV  bool
			sharpV bool
	*/
	/*
		%f     default width, default precision
		%9f    width 9, default precision
		%.2f   default width, precision 2
		%9.2f  width 9, precision 2
		%9.f   width 9, precision 0

		-	pad with spaces on the right rather than the left (left-justify the field)
		#	alternate format: add leading 0b for binary (%#b), 0 for octal (%#o),
			0x or 0X for hex (%#x or %#X); suppress 0x for %p (%#p);
			for %q, print a raw (backquoted) string if strconv.CanBackquote
			returns true;
			always print a decimal point for %e, %E, %f, %F, %g and %G;
			do not remove trailing zeros for %g and %G;
			write e.g. U+0078 'x' if the character is printable for %U (%#U).
		' '	(space) leave a space for elided sign in numbers (% d);
			put spaces between bytes printing strings or slices in hex (% x, % X)
		0	pad with leading zeros rather than spaces;
			for numbers, this moves the padding after the sign
	*/
	/*
		switch verb {
				case 'v', 's', 'x', 'X', 'q':
					// Is it an error or Stringer?
					// The duplication in the bodies is necessary:
					// setting handled and deferring catchPanic
					// must happen before calling the method.
					switch v := p.arg.(type) {
					case error:
						handled = true
						defer p.catchPanic(p.arg, verb, "Error")
						p.fmtString(v.Error(), verb)
						return

					case Stringer:
						handled = true
						defer p.catchPanic(p.arg, verb, "String")
						p.fmtString(v.String(), verb)
						return
					}
				}
	*/

}

func (grouping *addressDivisionGroupingInternal) GetPrefixLen() PrefixLen {
	return grouping.prefixLength
}

func (grouping *addressDivisionGroupingInternal) IsPrefixed() bool {
	return grouping.prefixLength != nil
}

//TODO LATER eventually when supporting large divisions,
//might move containsPrefixBlock(prefixLen BitCount), containsSinglePrefixBlock(prefixLen BitCount),
// GetMinPrefixLenForBlock, and GetPrefixLenForSingleBlock into groupingBase code
// IsPrefixBlock, IsSinglePrefixBlock
// which looks straightforward since none deal with DivInt, instead they all call into divisionValues interface

func (grouping *addressDivisionGroupingInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	if section := grouping.toAddressSection(); section != nil {
		return section.ContainsPrefixBlock(prefixLen)
	}
	prefixLen = checkSubnet(grouping.toAddressDivisionGrouping(), prefixLen)
	divisionCount := grouping.GetDivisionCount()
	var prevBitCount BitCount
	for i := 0; i < divisionCount; i++ {
		division := grouping.getDivision(i)
		bitCount := division.GetBitCount()
		totalBitCount := bitCount + prevBitCount
		if prefixLen < totalBitCount {
			divPrefixLen := prefixLen - prevBitCount
			if !division.containsPrefixBlock(divPrefixLen) {
				return false
			}
			for i++; i < divisionCount; i++ {
				division = grouping.getDivision(i)
				if !division.IsFullRange() {
					return false
				}
			}
			return true
		}
		prevBitCount = totalBitCount
	}
	return true
}

func (grouping *addressDivisionGroupingInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	prefixLen = checkSubnet(grouping.toAddressDivisionGrouping(), prefixLen)
	divisionCount := grouping.GetDivisionCount()
	var prevBitCount BitCount
	for i := 0; i < divisionCount; i++ {
		division := grouping.getDivision(i)
		bitCount := division.getBitCount()
		totalBitCount := bitCount + prevBitCount
		if prefixLen >= totalBitCount {
			if division.isMultiple() {
				return false
			}
		} else {
			divPrefixLen := prefixLen - prevBitCount
			if !division.ContainsSinglePrefixBlock(divPrefixLen) {
				return false
			}
			for i++; i < divisionCount; i++ {
				division = grouping.getDivision(i)
				if !division.IsFullRange() {
					return false
				}
			}
			return true
		}
		prevBitCount = totalBitCount
	}
	return true
}

func (grouping *addressDivisionGroupingInternal) IsPrefixBlock() bool { //Note for any given prefix length you can compare with GetMinPrefixLenForBlock
	prefLen := grouping.GetPrefixLen()
	return prefLen != nil && grouping.ContainsPrefixBlock(*prefLen)
}

func (grouping *addressDivisionGroupingInternal) IsSinglePrefixBlock() bool { //Note for any given prefix length you can compare with GetPrefixLenForSingleBlock
	calc := func() bool {
		prefLen := grouping.GetPrefixLen()
		return prefLen != nil && grouping.ContainsSinglePrefixBlock(*prefLen)
	}
	cache := grouping.cache
	if cache == nil {
		return calc()
	}
	res := cache.isSinglePrefixBlock
	if res == nil {
		if calc() {
			res = &trueVal

			// we can also set related cache fields
			pref := grouping.GetPrefixLen()
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.equivalentPrefix))
			atomic.StorePointer(dataLoc, unsafe.Pointer(pref))

			dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))
			atomic.StorePointer(dataLoc, unsafe.Pointer(pref))
		} else {
			res = &falseVal
		}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))
		atomic.StorePointer(dataLoc, unsafe.Pointer(res))
	}
	return *res
}

func (grouping *addressDivisionGroupingInternal) GetMinPrefixLenForBlock() BitCount {
	calc := func() BitCount {
		count := grouping.GetDivisionCount()
		totalPrefix := grouping.GetBitCount()
		for i := count - 1; i >= 0; i-- {
			div := grouping.getDivision(i)
			segBitCount := div.getBitCount()
			segPrefix := div.GetMinPrefixLenForBlock()
			if segPrefix == segBitCount {
				break
			} else {
				totalPrefix -= segBitCount
				if segPrefix != 0 {
					totalPrefix += segPrefix
					break
				}
			}
		}
		return totalPrefix
	}
	cache := grouping.cache
	if cache == nil {
		return calc()
	}
	res := cache.minPrefix
	if res == nil {
		val := calc()
		res = cacheBitCount(val)
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))
		atomic.StorePointer(dataLoc, unsafe.Pointer(res))
	}
	return *res
}

func (grouping *addressDivisionGroupingInternal) GetPrefixLenForSingleBlock() PrefixLen {
	calc := func() PrefixLen {
		count := grouping.GetDivisionCount()
		var totalPrefix BitCount
		for i := 0; i < count; i++ {
			div := grouping.getDivision(i)
			divPrefix := div.GetPrefixLenForSingleBlock()
			if divPrefix == nil {
				return nil
			}
			divPrefLen := *divPrefix
			totalPrefix += divPrefLen
			if divPrefLen < div.GetBitCount() {
				//remaining segments must be full range or we return nil
				for i++; i < count; i++ {
					laterDiv := grouping.getDivision(i)
					if !laterDiv.IsFullRange() {
						return nil
					}
				}
			}
		}
		return cacheBitCount(totalPrefix)
	}
	cache := grouping.cache
	if cache == nil {
		return calc()
	}
	res := cache.equivalentPrefix
	if res == nil {
		res = calc()
		if res == nil {
			res = noPrefix
			// we can also set related cache fields
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))
			atomic.StorePointer(dataLoc, unsafe.Pointer(&falseVal))
		} else {
			// we can also set related cache fields
			var isSingleBlock *bool
			if grouping.IsPrefixed() && PrefixEquals(res, grouping.GetPrefixLen()) {
				isSingleBlock = &trueVal
			} else {
				isSingleBlock = &falseVal
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.isSinglePrefixBlock))
			atomic.StorePointer(dataLoc, unsafe.Pointer(isSingleBlock))

			dataLoc = (*unsafe.Pointer)(unsafe.Pointer(&cache.minPrefix))
			atomic.StorePointer(dataLoc, unsafe.Pointer(res))
		}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.equivalentPrefix))
		atomic.StorePointer(dataLoc, unsafe.Pointer(res))
	}
	if res == noPrefix {
		return nil
	}
	return res
}

func (grouping *addressDivisionGroupingInternal) GetValue() *big.Int {
	if grouping.hasNoDivisions() {
		return bigZero()
	}
	return bigZero().SetBytes(grouping.getBytes())
}

func (grouping *addressDivisionGroupingInternal) GetUpperValue() *big.Int {
	if grouping.hasNoDivisions() {
		return bigZero()
	}
	return bigZero().SetBytes(grouping.getUpperBytes())
}

//func (grouping *addressDivisionGroupingInternal) Compare(item AddressItem) int {
//	xxx lowercase it xxxx
//	return CountComparator.Compare(grouping.toAddressDivisionGrouping(), item)
//}

//func (grouping *addressDivisionGroupingInternal) Equal(other GenericGroupingType) bool { xxxx need subs to have this xxxx
//	// For an identity comparison need to access the *addressDivisionGroupingBase or something
//	//otherSection := other.to
//	//if section.toAddressSection() == otherSection {
//	//	return true
//	//}
//	if section := grouping.toAddressSection(); section != nil {
//		if otherGrp, ok := other.(StandardDivisionGroupingType); ok {
//			otherSect := otherGrp.ToAddressDivisionGrouping().ToAddressSection()
//			return otherSect != nil && section.EqualsSection(otherSect)
//		}
//		return false
//	}
//	matchesStructure, count := grouping.matchesTypeAndCount(other)
//	if !matchesStructure {
//		return false
//	} else {
//		for i := 0; i < count; i++ {
//			one := grouping.getDivision(i)
//			two := other.GetGenericDivision(i)
//			if !one.Equal(two) { //this checks the division types and also the bit counts
//				return false
//			}
//		}
//	}
//	return true
//}

func (grouping *addressDivisionGroupingInternal) GetBytes() []byte {
	if grouping.hasNoDivisions() {
		return emptyBytes
	}
	cached := grouping.getBytes()
	return cloneBytes(cached)
}

func (grouping *addressDivisionGroupingInternal) GetUpperBytes() []byte {
	if grouping.hasNoDivisions() {
		return emptyBytes
	}
	cached := grouping.getUpperBytes()
	return cloneBytes(cached)
}

// CopyBytes gets the value for the lowest address in the range represented by this address division grouping.
//
// If the value fits in the given slice, the same slice is returned with the value.
// Otherwise, a new slice is allocated and returned with the value.
//
// You can use getBitCount() to determine the required array length for the bytes.
func (grouping *addressDivisionGroupingInternal) CopyBytes(bytes []byte) []byte {
	if grouping.hasNoDivisions() {
		if bytes != nil {
			return bytes
		}
		return emptyBytes
	}
	return getBytesCopy(bytes, grouping.getBytes())
}

func (grouping *addressDivisionGroupingInternal) CopyUpperBytes(bytes []byte) []byte {
	if grouping.hasNoDivisions() {
		if bytes != nil {
			return bytes
		}
		return emptyBytes
	}
	return getBytesCopy(bytes, grouping.getUpperBytes())
}

func (grouping *addressDivisionGroupingInternal) getBytes() (bytes []byte) {
	bytes, _ = grouping.getCachedBytes(grouping.calcBytes)
	return
}

func (grouping *addressDivisionGroupingInternal) getUpperBytes() (bytes []byte) {
	_, bytes = grouping.getCachedBytes(grouping.calcBytes)
	return
}

func (grouping *addressDivisionGroupingInternal) calcBytes() (bytes, upperBytes []byte) {
	addrType := grouping.getAddrType()
	divisionCount := grouping.GetDivisionCount()
	isMultiple := grouping.isMultiple()
	if addrType.isIPv4() || addrType.isMAC() {
		bytes = make([]byte, divisionCount)
		if isMultiple {
			upperBytes = make([]byte, divisionCount)
		} else {
			upperBytes = bytes
		}
		for i := 0; i < divisionCount; i++ {
			seg := grouping.getDivision(i).ToAddressSegment()
			bytes[i] = byte(seg.GetSegmentValue())
			if isMultiple {
				upperBytes[i] = byte(seg.GetUpperSegmentValue())
			}
		}
	} else if addrType.isIPv6() {
		byteCount := divisionCount << 1
		bytes = make([]byte, byteCount)
		if isMultiple {
			upperBytes = make([]byte, byteCount)
		} else {
			upperBytes = bytes
		}
		for i := 0; i < divisionCount; i++ {
			seg := grouping.getDivision(i).ToAddressSegment()
			byteIndex := i << 1
			val := seg.GetSegmentValue()
			bytes[byteIndex] = byte(val >> 8)
			var upperVal SegInt
			if isMultiple {
				upperVal = seg.GetUpperSegmentValue()
				upperBytes[byteIndex] = byte(upperVal >> 8)
			}
			nextByteIndex := byteIndex + 1
			bytes[nextByteIndex] = byte(val)
			if isMultiple {
				upperBytes[nextByteIndex] = byte(upperVal)
			}
		}
	} else {
		byteCount := grouping.GetByteCount()
		bytes = make([]byte, byteCount)
		if isMultiple {
			upperBytes = make([]byte, byteCount)
		} else {
			upperBytes = bytes
		}
		for k, byteIndex, bitIndex := divisionCount-1, byteCount-1, BitCount(8); k >= 0; k-- {
			div := grouping.getDivision(k)
			val := div.GetDivisionValue()
			var upperVal DivInt
			if isMultiple {
				upperVal = div.GetUpperDivisionValue()
			}
			divBits := div.GetBitCount()
			for divBits > 0 {
				rbi := 8 - bitIndex
				bytes[byteIndex] |= byte(val << uint(rbi))
				val >>= uint(bitIndex)
				if isMultiple {
					upperBytes[byteIndex] |= byte(upperVal << uint(rbi))
					upperVal >>= uint(bitIndex)
				}
				if divBits < bitIndex {
					bitIndex -= divBits
					break
				} else {
					divBits -= bitIndex
					bitIndex = 8
					byteIndex--
				}
			}
		}
	}
	return
}

// Returns whether the series represents a range of values that are sequential.
// Generally, this means that any division covering a range of values must be followed by divisions that are full range, covering all values.
func (grouping *addressDivisionGroupingInternal) IsSequential() bool {
	count := grouping.GetDivisionCount()
	if count > 1 {
		for i := 0; i < count; i++ {
			if grouping.getDivision(i).isMultiple() {
				for i++; i < count; i++ {
					if !grouping.getDivision(i).IsFullRange() {
						return false
					}
				}
				return true
			}
		}
	}
	return true
}

//func (grouping *addressDivisionGroupingInternal) Equal(other GenericGroupingType) bool {
//	// For an identity comparison need to access the *addressDivisionGroupingBase or something
//	//otherSection := other.to
//	//if section.toAddressSection() == otherSection {
//	//	return true
//	//}
//	if section := grouping.toAddressSection(); section != nil {
//		if otherGrouping, ok := other.(StandardDivisionGroupingType); ok {
//			if otherSection := otherGrouping.ToAddressDivisionGrouping().ToAddressSection(); otherSection != nil {
//				return section.EqualsSection(otherSection)
//			}
//		}
//		return false
//	}
//	matchesStructure, count := grouping.matchesTypeAndCount(other)
//	if !matchesStructure {
//		return false
//	} else {
//		for i := 0; i < count; i++ {
//			one := grouping.getDivision(i)
//			two := other.GetGenericDivision(i)
//			if !one.Equal(two) { //this checks the division types and also the bit counts
//				return false
//			}
//		}
//	}
//	return true
//}

//protected static interface GroupingCreator<S extends AddressDivisionBase> {
//		S createDivision(long value, long upperValue, int bitCount, int radix);
//	}

func (grouping *addressDivisionGroupingInternal) createNewDivisions(bitsPerDigit BitCount) ([]*AddressDivision, IncompatibleAddressError) {
	return grouping.createNewPrefixedDivisions(bitsPerDigit, nil)
}

//protected static interface PrefixedGroupingCreator<S extends AddressDivisionBase> {
//	S createDivision(long value, long upperValue, int bitCount, int radix, IPAddressNetwork<?, ?, ?, ?, ?> network, Integer prefixLength);
//}

func (grouping *addressDivisionGroupingInternal) createNewPrefixedDivisions(bitsPerDigit BitCount, networkPrefixLength PrefixLen) ([]*AddressDivision, IncompatibleAddressError) {
	//if(bitsPerDigit >= Integer.SIZE) {
	//	//keep in mind once you hit 5 bits per digit, radix 32, you need 32 different digits, and there are only 26 alphabet characters and 10 digit chars, so 36
	//	//so once you get higher than that, you need a new character set.
	//	//AddressLargeDivision allows all the way up to base 85
	//	throw new AddressValueException(bitsPerDigit);
	//}
	bitCount := grouping.GetBitCount()
	//List<Integer> bitDivs = new ArrayList<Integer>(bitsPerDigit);
	var bitDivs []BitCount

	// here we divide into divisions, each with an exact number of digits.
	// Each digit takes 3 bits.  So the division bit-sizes are a multiple of 3 until the last one.

	//ipv6 octal:
	//seg bit counts: 63, 63, 2
	//ipv4 octal:
	//seg bit counts: 30, 2

	largestBitCount := BitCount(64) // uint64, size of DivInt

	//int largestBitCount = Long.SIZE - 1;
	largestBitCount -= largestBitCount % bitsPerDigit // round off to a multiple of 3 bits
	for {
		if bitCount <= largestBitCount {
			mod := bitCount % bitsPerDigit
			secondLast := bitCount - mod
			if secondLast > 0 {
				//bitDivs.add(cacheBits(secondLast));
				bitDivs = append(bitDivs, secondLast)
			}
			if mod > 0 {
				bitDivs = append(bitDivs, mod)
				//bitDivs.add(cacheBits(mod));
			}
			break
		} else {
			bitCount -= largestBitCount
			bitDivs = append(bitDivs, largestBitCount)
			//bitDivs.add(cacheBits(largestBitCount));
		}
	}

	// at this point bitDivs has our division sizes

	divCount := len(bitDivs)
	divs := make([]*AddressDivision, divCount)
	//S divs[] = groupingArrayCreator.apply(divCount);
	currentSegmentIndex := 0
	seg := grouping.getDivision(currentSegmentIndex)
	segLowerVal := seg.GetDivisionValue()
	segUpperVal := seg.GetUpperDivisionValue()
	segBits := seg.GetBitCount()
	bitsSoFar := BitCount(0)

	// 2 to the x is all ones shift left x, then not, then add 1
	// so, for x == 1, 1111111 -> 1111110 -> 0000001 -> 0000010
	radix := ^(^(0) << uint(bitsPerDigit)) + 1
	//int radix = AddressDivision.getRadixPower(BigInteger.valueOf(2), bitsPerDigit).intValue();
	//fill up our new divisions, one by one
	for i := divCount - 1; i >= 0; i-- {

		//int originalDivBitSize, divBitSize;
		divBitSize := bitDivs[i]
		originalDivBitSize := divBitSize
		//long divLowerValue, divUpperValue;
		//divLowerValue = divUpperValue = 0;
		var divLowerValue, divUpperValue uint64
		for {
			if segBits >= divBitSize { // this segment fills the remainder of this division
				diff := uint(segBits - divBitSize)
				segBits = BitCount(diff)
				//udiff := uint(diff);
				segL := segLowerVal >> diff
				segU := segUpperVal >> diff

				// if the division upper bits are multiple, then the lower bits inserted must be full range
				if divLowerValue != divUpperValue {
					if segL != 0 || segU != ^(^uint64(0)<<uint(divBitSize)) {
						return nil, &incompatibleAddressError{addressError: addressError{key: "ipaddress.error.invalid.joined.ranges"}}
					}
				}

				divLowerValue |= segL
				divUpperValue |= segU

				shift := ^(^uint64(0) << diff)
				segLowerVal &= shift
				segUpperVal &= shift

				// if a segment's bits are split into two divisions, and the bits going into the first division are multi-valued,
				// then the bits going into the second division must be full range
				if segL != segU {
					if segLowerVal != 0 || segUpperVal != ^(^uint64(0)<<uint(segBits)) {
						return nil, &incompatibleAddressError{addressError: addressError{key: "ipaddress.error.invalid.joined.ranges"}}
					}
				}

				var segPrefixBits PrefixLen
				if networkPrefixLength != nil {
					segPrefixBits = getDivisionPrefixLength(originalDivBitSize, *networkPrefixLength-bitsSoFar)
				}
				//Integer segPrefixBits = networkPrefixLength == null ? null : getSegmentPrefixLength(originalDivBitSize, networkPrefixLength - bitsSoFar);
				div := NewRangePrefixDivision(divLowerValue, divUpperValue, segPrefixBits, originalDivBitSize, radix)
				//S div = groupingCreator.createDivision(divLowerValue, divUpperValue, originalDivBitSize, radix, network, segPrefixBits);
				divs[divCount-i-1] = div
				if segBits == 0 && i > 0 {
					//get next seg
					currentSegmentIndex++
					seg = grouping.getDivision(currentSegmentIndex)
					segLowerVal = seg.getDivisionValue()
					segUpperVal = seg.getUpperDivisionValue()
					segBits = seg.getBitCount()
				}
				break
			} else {
				// if the division upper bits are multiple, then the lower bits inserted must be full range
				if divLowerValue != divUpperValue {
					if segLowerVal != 0 || segUpperVal != ^(^uint64(0)<<uint(segBits)) {
						return nil, &incompatibleAddressError{addressError: addressError{key: "ipaddress.error.invalid.joined.ranges"}}
					}
				}
				diff := uint(divBitSize - segBits)
				divLowerValue |= segLowerVal << diff
				divUpperValue |= segUpperVal << diff
				divBitSize = BitCount(diff)

				//get next seg
				currentSegmentIndex++
				seg = grouping.getDivision(currentSegmentIndex)
				segLowerVal = seg.getDivisionValue()
				segUpperVal = seg.getUpperDivisionValue()
				segBits = seg.getBitCount()
			}
		}
		bitsSoFar += originalDivBitSize
	}
	return divs, nil
}

type AddressDivisionGrouping struct {
	addressDivisionGroupingInternal
}

func (grouping *AddressDivisionGrouping) Compare(item AddressItem) int {
	return CountComparator.Compare(grouping, item)
}

func (grouping *AddressDivisionGrouping) CompareSize(other StandardDivisionGroupingType) int {
	if grouping == nil {
		if other != nil && other.ToAddressDivisionGrouping() != nil {
			// we have size 0, other has size >= 1
			return -1
		}
		return 0
	}
	return grouping.compareSize(other)
}

func (grouping *AddressDivisionGrouping) GetCount() *big.Int {
	if grouping == nil {
		return bigZero()
	}
	return grouping.getCount()
}

func (grouping *AddressDivisionGrouping) IsMultiple() bool {
	return grouping != nil && grouping.isMultiple()
}

// copySubDivisions copies the existing divisions from the given start index until but not including the division at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (grouping *AddressDivisionGrouping) CopySubDivisions(start, end int, divs []*AddressDivision) (count int) {
	return grouping.copySubDivisions(start, end, divs)
}

// CopyDivisions copies the existing divisions from the given start index until but not including the division at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied
func (grouping *AddressDivisionGrouping) CopyDivisions(divs []*AddressDivision) (count int) {
	return grouping.copyDivisions(divs)
}

func (grouping *AddressDivisionGrouping) GetDivisionStrings() []string {
	return grouping.getDivisionStrings()
}

func (grouping *AddressDivisionGrouping) IsAddressSection() bool {
	return grouping != nil && grouping.isAddressSection()
}

func (grouping *AddressDivisionGrouping) IsIPAddressSection() bool {
	return grouping.ToAddressSection().IsIPAddressSection()
}

func (grouping *AddressDivisionGrouping) IsIPv4AddressSection() bool {
	return grouping.ToAddressSection().IsIPv4AddressSection()
}

func (grouping *AddressDivisionGrouping) IsIPv6AddressSection() bool {
	return grouping.ToAddressSection().IsIPv6AddressSection()
}

func (grouping *AddressDivisionGrouping) IsIPv6v4MixedAddressGrouping() bool {
	return grouping.matchesIPv6v4MixedGroupingType()
}

func (grouping *AddressDivisionGrouping) IsMACAddressSection() bool {
	return grouping.ToAddressSection().IsMACAddressSection()
}

// ToAddressSection converts to an address section.
// If the conversion cannot happen due to division size or count, the result will be the zero value.
func (grouping *AddressDivisionGrouping) ToAddressSection() *AddressSection {
	if grouping == nil || !grouping.isAddressSection() {
		return nil
	}
	return (*AddressSection)(unsafe.Pointer(grouping))
}

func (grouping *AddressDivisionGrouping) ToIPv4v6MixedAddressGrouping() *IPv6v4MixedAddressGrouping {
	if grouping.matchesIPv6v4MixedGroupingType() {
		return (*IPv6v4MixedAddressGrouping)(grouping)
	}
	return nil
}

func (grouping *AddressDivisionGrouping) ToIPAddressSection() *IPAddressSection {
	return grouping.ToAddressSection().ToIPAddressSection()
}

func (grouping *AddressDivisionGrouping) ToIPv6AddressSection() *IPv6AddressSection {
	return grouping.ToAddressSection().ToIPv6AddressSection()
}

func (grouping *AddressDivisionGrouping) ToIPv4AddressSection() *IPv4AddressSection {
	return grouping.ToAddressSection().ToIPv4AddressSection()
}

func (grouping *AddressDivisionGrouping) ToMACAddressSection() *MACAddressSection {
	return grouping.ToAddressSection().ToMACAddressSection()
}

func (grouping *AddressDivisionGrouping) ToAddressDivisionGrouping() *AddressDivisionGrouping {
	return grouping
}

func (grouping *AddressDivisionGrouping) GetDivision(index int) *AddressDivision {
	return grouping.getDivision(index)
}

func (grouping *AddressDivisionGrouping) String() string {
	if grouping == nil {
		return nilString()
	}
	return grouping.toString()
}
