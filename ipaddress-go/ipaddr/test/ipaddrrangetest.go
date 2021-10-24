package test

import (
	"github.com/seancfoley/ipaddress/ipaddress-go/ipaddr"
	"math/big"
	"strconv"
	"strings"
)

type ipAddressRangeTester struct {
	ipAddressTester
}

func (t ipAddressRangeTester) run() {
	t.testEquivalentPrefix("*.*.*.*", 0)
	t.testEquivalentPrefix("0-127.*.*.*", 1)
	t.testEquivalentPrefix("128-255.*.*.*", 1)
	t.testEquivalentPrefix("*.*.*.*/1", 0)
	//t.testEquivalentPrefix("0.*.*.*/1", 1)
	t.testEquivalentPrefix("0.*.*.*/1", 8)
	t.testEquivalentPrefix("128-255.*.*.*/1", 1)
	t.testEquivalentPrefix("1.2.*.*", 16)
	t.testEquivalentPrefix("1.2.*.*/24", 16)
	t.testEquivalentMinPrefix("1.2.*.0/24", cacheTestBits(16), 16)
	t.testEquivalentMinPrefix("1.2.0-255.0/24", cacheTestBits(16), 16)
	t.testEquivalentPrefix("1.2.1.0/24", 24)
	t.testEquivalentPrefix("1.2.1.*/24", 24)
	t.testEquivalentPrefix("1.2.1.*", 24)
	t.testEquivalentMinPrefix("1.2.*.4", nil, 32)
	t.testEquivalentPrefix("1.2.252-255.*", 22)
	t.testEquivalentPrefix("1.2.252-255.0-255", 22)
	t.testEquivalentPrefix("1.2.0-3.0-255", 22)
	t.testEquivalentPrefix("1.2.128-131.0-255", 22)
	t.testEquivalentMinPrefix("1.2.253-255.0-255", nil, 24)
	t.testEquivalentMinPrefix("1.2.252-255.0-254", nil, 32)
	t.testEquivalentMinPrefix("1.2.251-255.0-254", nil, 32)
	t.testEquivalentMinPrefix("1.2.251-255.0-255", nil, 24)

	t.testEquivalentMinPrefix("1.2.1-3.*", nil, 24)
	t.testEquivalentPrefix("1.2.0-3.*", 22)

	t.testEquivalentPrefix("*:*", 0)
	t.testEquivalentPrefix("::/0", 0)
	t.testEquivalentMinPrefix("0-1::/0", nil, 128)
	t.testEquivalentPrefix("::/1", 1)
	t.testEquivalentMinPrefix("0-1::/1", nil, 128)
	t.testEquivalentMinPrefix("8000-ffff::/1", nil, 128)
	t.testEquivalentPrefix("8000-ffff:*", 1)
	t.testEquivalentMinPrefix("7fff-ffff:*", nil, 16)
	t.testEquivalentMinPrefix("7fff-ffff:*/1", nil, 16)
	t.testEquivalentPrefix("11:8000-ffff:*/1", 17)
	t.testEquivalentPrefix("11:8000-ffff:*", 17)
	t.testEquivalentPrefix("1:2:*", 32)
	t.testEquivalentMinPrefix("1:2:*:*::/64", cacheTestBits(32), 32)
	t.testEquivalentPrefix("1:2:*:*/64", 32)
	t.testEquivalentPrefix("1:2:3:4:5:*:*/64", 80)
	t.testEquivalentMinPrefix("1:2:*::/64", nil, 64)
	t.testEquivalentMinPrefix("1:2:*::", nil, 128)
	t.testEquivalentPrefix("1:2:8000-ffff:*", 33)
	t.testEquivalentPrefix("1:2:0000-7fff:*", 33)
	t.testEquivalentPrefix("1:2:c000-ffff:*", 34)
	t.testEquivalentPrefix("1:2:0000-3fff:*", 34)
	t.testEquivalentPrefix("1:2:8000-bfff:*", 34)
	t.testEquivalentPrefix("1:2:4000-7fff:*", 34)
	t.testEquivalentPrefix("1:2:fffc-ffff:*", 46)
	t.testEquivalentPrefix("1:2:fffc-ffff:0-ffff:*", 46)
	t.testEquivalentMinPrefix("1:2:fffd-ffff:0-ffff:*", nil, 48)
	t.testEquivalentMinPrefix("1:2:fffc-ffff:0-fffe:*", nil, 64)
	t.testEquivalentMinPrefix("1:2:fffb-ffff:0-fffe:*", nil, 64)
	t.testEquivalentMinPrefix("1:2:fffb-ffff:0-ffff:*", nil, 48)

	t.testTrees()

	t.testStrings()

	t.testReverse("1:2:*:4:5:6:a:b", false, false)
	t.testReverse("1:1:1:1-fffe:2:3:3:3", false, false)                                   // 0x1-0xfffe reverseBitsPerByte throws
	t.testReverse("1-fffe:0-ffff:0-ffff:0-fffe:1-ffff:1-ffff:1-fffe:1-ffff", false, true) // all reversible
	t.testReverse("1-fffe:0-ffff:1-ffff:0-fffe:0-fffe:1-ffff:0-ffff:1-fffe", true, true)  // all reversible
	t.testReverse("1:1:1:0-fffe:1-fffe:*:1:1", false, false)                              // 100-feff or aa01-aafe are byte reversible becoming 100-feff and xx01-xxfe where x is reverse of a
	t.testReverse("ffff:80:*:ff:01:ffff", false, false)
	t.testReverse("ffff:8000:fffe::7fff:0001:ffff", true, false)
	t.testReverse("ffff:8000:*:8000:1:*:01:ffff", true, false)
	t.testReverse("ffff:8118:ffff:*:1-fffe:ffff", false, true)
	t.testReverse("ffff:8181:c3c3::4224:2400:0-fffe", false, true)
	t.testReverse("ffff:1:ff:ff:*:*", false, false)

	t.testPrefixes("255.127.0.0/16",
		16, -5,
		"255.127.0.0/24",
		"255.0.0.0/8",
		"255.96.*.*/11",
		//"255.96.0.0/11",
		"255.127.0.0/16",
		"255.127.0.0/16")

	t.testPrefixes("255.127.0.0/17",
		16, -17,
		"255.127.0.0/24",
		"255.127.0.0/16",
		//"0.0.0.0/0",
		"0.0.0-127.*/0",
		"255.127.0-127.*/16",
		//"255.127.0.0/16",
		"255.127.0.0/16")

	t.testPrefixes("ffff:ffff:1:ffff::/64",
		16, -5,
		"ffff:ffff:1:ffff::/80",
		"ffff:ffff:1::/48",
		"ffff:ffff:1:ffe0:*:*:*:*/59",
		//"ffff:ffff:1:ffe0::/59",
		//"ffff::/16",
		"ffff::*:*:*:*/16",
		"ffff::/16")

	t.testPrefixes("ffff:ffff:1:ffff::/64",
		16, 1,
		"ffff:ffff:1:ffff::/80",
		"ffff:ffff:1::/48",
		"ffff:ffff:1:ffff::/65",
		//"ffff::/16",
		"ffff::*:*:*:*/16",
		"ffff::/16")

	t.testBitwiseOr("1.2.0.0/16", cacheTestBits(8), "0.0.3.248", "1.2.3.248-255")
	t.testBitwiseOr("1.2.0.0/16", cacheTestBits(7), "0.0.2.0", "1.2.2-3.*")
	t.testBitwiseOr("1.2.*.*", cacheTestBits(7), "0.0.3.0", "")
	t.testBitwiseOr("1.2.0-3.*", cacheTestBits(0), "0.0.3.0", "1.2.3.*")
	t.testBitwiseOr("1.2.0.0/16", cacheTestBits(7), "0.0.3.0", "1.2.3.*")
	t.testBitwiseOr("0.0.0.0/0", cacheTestBits(0), "128.192.224.240", "128-255.192-255.224-255.240-255")
	t.testBitwiseOr("*.*", cacheTestBits(0), "128.192.224.240", "128-255.192-255.224-255.240-255")
	t.testBitwiseOr("0.0.0.0/0", cacheTestBits(0), "128.192.224.64", "")
	t.testBitwiseOr("*.*", cacheTestBits(0), "128.192.224.64", "")
	t.testPrefixBitwiseOr("1.3.0.0/15", 24, "0.0.255.1", "1.3.255.0", "1.3.255.1/15")
	t.testPrefixBitwiseOr("1.3.0.1/15", 24, "0.0.255.1", "1.3.255.1/24", "1.3.255.1/15")
	t.testPrefixBitwiseOr("1.3.0.1/15", 24, "0.0.255.0", "1.3.255.1/24", "1.3.255.1/15")
	t.testPrefixBitwiseOr("1.2.0.0/22", 24, "0.0.3.248", "1.2.3.0/24", ("1.2.3.248-255/22"))
	t.testPrefixBitwiseOr("1.2.0.0/24", 24, "0.0.3.248", "1.2.3.0/24", ("1.2.3.248-255/24"))
	t.testPrefixBitwiseOr("1.2.0.0/22", 23, "0.0.3.0", "1.2.2.0/23", "1.2.3.0-255/22")
	t.testPrefixBitwiseOr("1.2.0.0/24", 23, "0.0.3.0", "1.2.2.*", ("1.2.3.0-255/24"))
	t.testPrefixBitwiseOr("1:2::/46", 47, "0:0:3::", "1:2:2::/47", "1:2:3:*:*:*:*:*/46")

	t.testPrefixBitwiseOr("0.0.0.0/16", 18, "0.0.2.8", "0.0.0-192.0/18", "")

	t.testBitwiseOr("1:2::/32", cacheTestBits(16), "0:0:3:fff8::", "1:2:3:fff8-ffff:*")
	t.testBitwiseOr("1:2::/32", cacheTestBits(15), "0:0:2::", "1:2:2-3:*")
	t.testBitwiseOr("1:2:*", cacheTestBits(0), "0:0:8000::", "1:2:8000-ffff:*")
	t.testBitwiseOr("1:2:*", cacheTestBits(0), "0:0:c000::", "1:2:c000-ffff:*")
	t.testBitwiseOr("1:2::/32", cacheTestBits(15), "0:0:3::", "1:2:3:*")
	t.testBitwiseOr("::/0", cacheTestBits(0), "8000:c000:e000:fff0::", "8000-ffff:c000-ffff:e000-ffff:fff0-ffff:*")
	t.testBitwiseOr("*:*", cacheTestBits(0), "8000:c000:e000:fff0::", "8000-ffff:c000-ffff:e000-ffff:fff0-ffff:*")
	t.testBitwiseOr("::/0", cacheTestBits(0), "8000:c000:e000:4000::", "")
	t.testBitwiseOr("1:1::/16", cacheTestBits(32), "0:2:3::ffff", "1:2:3::ffff")       //note the prefix length is dropped to become "1.2.3.*", but equality still holds
	t.testBitwiseOr("1:1:0:*:0/16", nil, "0:2:3:*:ffff", "1:3:3:*:*:*:*:ffff")         //note the prefix length is dropped to become "1.2.3.*", but equality still holds
	t.testBitwiseOr("1:0:0:1::/16", cacheTestBits(32), "0:2:3::ffff", "1:2:3:1::ffff") //note the prefix length is dropped to become "1.2.3.*", but equality still holds
	t.testPrefixBitwiseOr("::/32", 34, "0:0:2:8::", "0:0:0-c000::/34", "")

	t.testDelimitedCount("1,2-3,4:3:4,5:6:7:8:ffff:ffff", 8)
	t.testDelimitedCount("1,2::3,6:7:8:4,5-6:6,8", 16)
	t.testDelimitedCount("1:2:3:*:4::5", 1)
	t.testDelimitedCount("1:2,3,*:3:ffff:ffff:6:4:5,ff,7,8,99", 15)
	t.testDelimitedCount("0,1-2,3,5:3::6:4:5,ffff,7,8,99", 30)

	//if(false) {
	//	testMatches(true, "1.2.3.4/16", "1.2.*.*");
	//	testMatches(true, "1.2.3.4/16", "1.2.*");
	//	testMatches(false, "1.2.3.4/15", "1.2.*.*");
	//	testMatches(false, "1.2.3.4/17", "1.2.*.*");
	//} else {
	t.testMatches(true, "1.2.3.4/16", "1.2.3.4")
	t.testMatches(true, "1.2.3.4/15", "1.2.3.4")
	t.testMatches(true, "1.2.3.4/17", "1.2.3.4")

	t.testMatches(true, "1.2.0.4/16", "1.2.0.4")
	t.testMatches(true, "1.2.3.0/16", "1.2.3.0")

	t.testMatches(true, "1.2.3.4/14", "1.2.3.4")
	t.testMatches(true, "1.2.0.4/14", "1.2.0.4")
	t.testMatches(true, "1.2.0.0/14", "1.2.0.0")
	t.testMatches(true, "1.0.3.0/14", "1.0.3.0")

	//}

	t.testMatches(true, "1.2.0.0/16", "1.2.*.*")
	t.testMatches(true, "1.2.0.0/16", "1.2.*")

	t.testMatches(true, "1.4.0.0/14", "1.4-7.*")
	t.testMatches(true, "1.4.0.0/14", "1.4-7.*.*")

	t.testMatches(false, "1.2.3.4/16", "1.2.*/255.255.0.0")
	t.testMatches(false, "1.2.3.4/15", "1.2.3.*/255.254.0.0")
	t.testMatches(false, "1.2.3.4/17", "1.2.3.*/255.255.128.0")

	t.testMatches(true, "1.2.0.0/16", "1.2.*/255.255.0.0")
	t.testMatches(true, "1.2.3.*/15", "1.2.3.*/255.254.0.0")
	t.testMatches(true, "1.2.3.*/17", "1.2.3.*/255.255.128.0")

	t.testMatches(false, "1.1.3.4/15", "1.2.3.*/255.254.0.0")
	t.testMatches(false, "1.1.3.4/17", "1.2.3.*/255.255.128.0")

	t.testMatches(true, "1:2::/32", "1:2:*:*:*:*:*:*")
	t.testMatches(true, "1:2::/32", "1:2:*:*:*:*:*.*.*.*")
	t.testMatches(true, "1:2::/32", "1:2:*")
	t.testMatches(false, "1:2::/32", "1:2:*:*:*:*:3:*")
	t.testMatches(false, "1:2::/32", "1:2:*:*:*:*:*.*.3.*")
	t.testMatches(false, "1:2::/31", "1:2:*")
	t.testMatches(false, "1:2::/33", "1:2::*")

	t.testMatches(true, "1:2::/32", "1:2:*:*:*:*:*:*/ffff:ffff::")
	t.testMatches(true, "1:2::/31", "1:2-3:*:*:*:*:*:*/ffff:fffe::")
	t.testMatches(true, "1:2::/33", "1:2:0-7fff:*:*:*:*:*/ffff:ffff:8000::")

	t.testMatches(false, "1:2::/24", "1:__:*")
	t.testMatches(false, "1:2::/28", "1:_::/32")
	t.testMatches(false, "1:2::/20", "1:___::/32")
	t.testMatches(false, "1:2::/16", "1:____::/32")
	t.testMatches(false, "1:ffef::/24", "1:ff__::/32")
	t.testMatches(false, "1:ffef::/24", "1:ff__:*:*")

	t.testMatches(true, "1::/24", "1:__:*")
	t.testMatches(true, "1::/28", "1:_::/32")
	t.testMatches(true, "1::/20", "1:___::/32")
	t.testMatches(true, "1::/16", "1:____::/32")
	t.testMatches(true, "1:ff00::/24", "1:ff__::/32")
	t.testMatches(true, "1:ff00::/24", "1:ff__:*:*")

	t.testMatches(true, "250-255.200-255.0-255.20-29", "25_.2__.___.2_")
	t.testMatches(true, "150-159.100-199.0-99.10-19", "15_.1__.__.1_")
	t.testMatches(false, "251-255.200-255.0-255.20-29", "25_.2__.___.2_")
	t.testMatches(false, "150-158.100-199.0-99.10-19", "15_.1__.__.1_")
	t.testMatches(true, "250-25f:200-2ff:0-fff:20-2f::", "25_:2__:___:2_::")
	t.testMatches(true, "150-15f:100-1ff:0-ff:10-1f::", "15_:1__:__:1_::")
	t.testMatches(false, "250-25f:201-2ff:0-fff:20-2f::", "25_:2__:___:2_::")
	t.testMatches(false, "150-15f:100-1ef:0-ff:10-1f::", "15_:1__:__:1_::")
	t.testMatches(true, "::250-25f:200-2ff:0-fff:20-2f", "::25_:2__:___:2_")
	t.testMatches(true, "::150-15f:100-1ff:0-ff:10-1f", "::15_:1__:__:1_")
	t.testMatches(true, "250-25f:200-2ff::0-fff:20-2f", "25_:2__::___:2_")
	t.testMatches(true, "150-15f:100-1ff::0-ff:10-1f", "15_:1__::__:1_")

	t.testMatches(true, "1:2:3:4:5:6:1-2.*.0.4", "1:2:3:4:5:6:100-2ff:4")         // mixed starting with range
	t.testMatches(true, "1:2:3:4:5:6:1.2.0.4-5", "1:2:3:4:5:6:102:4-5")           // mixed ending with range
	t.testMatches(true, "1:2:3:4:5:6:1.2.0.*", "1:2:3:4:5:6:102:0-ff")            // mixed ending with range
	t.testMatches(true, "1:2:3:4:5:6:1.2.0._", "1:2:3:4:5:6:102:0-9")             // mixed ending with range
	t.testMatches(true, "1:2:3:4:5:6:1.2.0.1_", "1:2:3:4:5:6:102:a-13")           // mixed ending with range
	t.testMatches(true, "1:2:3:4:5:6:1.2.0.4-5", "1:2:3:4:5:6:102:5-4")           // mixed ending with range
	t.testMatches(true, "1:2:3:4:5:6:1.2.0.4-5", "1:2:3:4:5:6:1.2.0.5-4")         // mixed ending with range
	t.testMatches(true, "1:2:3:4:5:6:1.2-255.0.4-5", "1:2:3:4:5:6:1.255-2.0.5-4") // mixed ending with range
	t.testMatches(false, "1:2:3:4:5:6:1-3.2.0.4-5", "1:2:3:4:5:6:3-1.2.0.5-4")    // inet.ipaddr.IncompatibleAddressException: 1-3, 2, IP Address error: IPv4 segment ranges cannot be converted to IPv6 segment ranges
	t.testMatches(true, "1:2:3:4:5:6:1-3.*.0.4-5", "1:2:3:4:5:6:3-1.*.0.5-4")
	t.testMatches(true, "1:2:3:4:5:6:1-3.*.0.4-5", "1:2:3:4:5:6:3ff-100:5-4")

	t.testMatches(true, "1.2.2-3.4", "1.2.3-2.4")
	t.testMatches(true, "1.255-2.2-3.4", "1.2-255.3-2.4")
	t.testMatches(true, "1:2:3:4:5:6:7:7-8", "1:2:3:4:5:6:7:8-7")
	t.testMatches(true, "1-ffff:2:3:4:5:6:7:7-8", "ffff-1:2:3:4:5:6:7:8-7")
	t.testMatches(true, "1-ffff:2:3:4:aa-5:6:7:7-8", "ffff-1:2:3:4:5-aa:6:7:8-7")
	t.testMatches(true, "1.2.*.4", "1.2.255-0.4")
	t.testMatches(true, "1:2:3:4:5:*:7:7-8", "1:2:3:4:5:ffff-0:7:8-7")

	t.testMatchesInetAton(true, "1.2.3", "1.2.0.3", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "0x1.0x2.2-0x3.0x4", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "0x1.0x2.0x2-0x3.0x4", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "0x1.0x2.0x2-3.0x4", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "01.02.2-03.04", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "01.02.2-3.04", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "01.02.02-03.04", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "01.02.0x2-03.04", true)
	t.testMatchesInetAton(true, "1.2.2-3.4", "01.02.0x2-0x3.04", true)
	t.testMatchesInetAton(true, "1.2.0200-0277.4", "01.02.02__.04", true)
	t.testMatchesInetAton(true, "1.2.0x20-0x2f.4", "01.02.0x2_.04", true)
	t.testMatchesInetAton(true, "1.2.0x10-0x1f.4", "01.02.0x1_.04", true)
	t.testMatchesInetAton(true, "1.2.*.4", "01.02.0x__.04", true)
	t.testMatchesInetAton(true, "1.2.0-077.4", "01.02.0__.04", true)

	t.testMatchesInetAton(true, "1.2.2-3.4", "01.02.0x2-0x3.04", true)

	t.testMatchesInetAton(true, "0.0.0-1.4", "00.0x0.0x00-0x000001.04", true)
	t.testMatchesInetAton(true, "11.10-11.10-11.10-11", "11.012-0xb.0xa-013.012-0xB", true)
	t.testMatchesInetAton(true, "11.10-11.*.10-11", "11.012-0xb.0x0-0xff.012-0xB", true)
	t.testMatchesInetAton(true, "1.*", "1.*.*.0x0-0xff", true)
	t.testMatchesInetAton(true, "1.*", "1.0-255.0-65535", true)
	t.testMatchesInetAton(true, "1.*", "1.0-0xff.0-0xffff", true)
	t.testMatchesInetAton(true, "1.*", "1.0x0-0xff.00-0xffff", true)

	t.testMatchesInetAton(true, "11.11.0-11.*", "11.11.0-0xbff", true)
	t.testMatchesInetAton(true, "11.0.0.11-11", "11.0x00000000000000000b-0000000000000000000013", true)
	t.testMatchesInetAton(true, "11.1-11.*/16", "11.0x10000-786431/16", true)
	t.testMatchesInetAton(true, "11.1-11.*/16", "11.0x10000-0xbffff/16", true)

	t.testMatches(true, "1:2:3:4:5:6:1.2.3.4/128", "1:2:3:4:5:6:102:304")

	t.testMatches(false, "1:2:3:4:5:6:1.2.3.4/96", "1:2:3:4:5:6:*:*")
	t.testMatches(false, "1:2:3:4:5:6:255.2.3.4/97", "1:2:3:4:5:6:8000-ffff:*")
	t.testMatches(false, "1:2:3:4:5:6:1.2.3.4/112", "1:2:3:4:5:6:102:*")
	t.testMatches(false, "1:2:3:4:5:6:1.2.255.4/115", "1:2:3:4:5:6:102:e000-ffff")

	t.testMatches(false, "1.2.3.4/0", "*.*")
	t.testMatches(false, "1.2.3.4/0", "*.*.*.*")
	t.testMatches(false, "1:2:3:4:5:6:7:8/0", "*:*")
	t.testMatches(false, "1:2:3:4:5:6:7:8/0", "*:*:*:*:*:*:*:*")

	t.testMatches(true, "1:2:3:4:5:6:1.2.3.4/96", "1:2:3:4:5:6:102:304")
	t.testMatches(true, "1:2:3:4:5:6:255.2.3.4/97", "1:2:3:4:5:6:ff02:304")
	t.testMatches(true, "1:2:3:4:5:6:1.2.3.4/112", "1:2:3:4:5:6:102:304")
	t.testMatches(true, "1:2:3:4:5:6:1.2.255.4/115", "1:2:3:4:5:6:102:ff04")

	t.testMatches(true, "1.2.3.4/0", "1.2.3.4")
	t.testMatches(true, "1.2.3.4/0", "1.2.3.4")
	t.testMatches(true, "1:2:3:4:5:6:7:8/0", "1:2:3:4:5:6:7:8")
	t.testMatches(true, "1:2:3:4:5:6:7:8/0", "1:2:3:4:5:6:7:8")

	t.testMatches(true, "1:2:3:4:5:6:0.0.0.0/96", "1:2:3:4:5:6:*:*")
	t.testMatches(false, "1:2:3:4:5:6:255.0.0.0/97", "1:2:3:4:5:6:8000-ffff:*")
	t.testMatches(true, "1:2:3:4:5:6:255.0.0.0/97", "1:2:3:4:5:6:ff00:0")
	t.testMatches(true, "1:2:3:4:5:6:128.0.0.0/97", "1:2:3:4:5:6:8000-ffff:*")
	t.testMatches(true, "1:2:3:4:5:6:1.2.0.0/112", "1:2:3:4:5:6:102:*")
	t.testMatches(false, "1:2:3:4:5:6:1.2.255.0/115", "1:2:3:4:5:6:102:e000-ffff")
	t.testMatches(true, "1:2:3:4:5:6:1.2.255.0/115", "1:2:3:4:5:6:102:FF00")
	t.testMatches(true, "1:2:3:4:5:6:1.2.224.0/115", "1:2:3:4:5:6:102:e000-ffff")

	t.testMatches(true, "0.0.0.0/0", "*.*")
	t.testMatches(true, "0.0.0.0/0", "*.*.*.*")
	t.testMatches(true, "::/0", "*:*")
	t.testMatches(true, "::/0", "*:*:*:*:*:*:*:*")

	t.testMatches(true, "1-02.03-4.05-06.07", "1-2.3-4.5-6.7")
	t.testMatches(true, "1-002.003-4.005-006.007", "1-2.3-4.5-6.7")

	t.testMatches(true, "1-2.0-0.00-00.00-0", "1-2.0.0.0")
	t.testMatches(true, "1-2:0-0:00-00:00-0:0-000:0000-0000:0000-00:0000-0", "1-2:0:0:0:0:0:0:0")
	t.testMatches(true, "00-0.0-0.00-00.00-0", "0.0.0.0")
	t.testMatches(true, "0-00:0-0:00-00:00-0:0-000:0000-0000:0000-00:0000-0", "::")

	t.testMatches(true, "-1.22.33.4", "0-1.22.33.4")
	t.testMatches(true, "-1.22.33.4", "0-1.22.33.4")
	t.testMatches(true, "22.1-.33.4", "22.1-255.33.4")
	t.testMatches(true, "22.33.4.1-", "22.33.4.1-255")
	t.testMatches(true, "aa:-1:cc::d:ee:f", "aa:0-1:cc::d:ee:f")
	t.testMatches(true, "aa:dd-:cc::d:ee:f", "aa:dd-ffff:cc::d:ee:f")
	t.testMatches(true, "aa:dd-:cc::d:ee:f-", "aa:dd-ffff:cc::d:ee:f-ffff")
	t.testMatches(true, "-:0:0:0:0:0:0:0", "0-ffff:0:0:0:0:0:0:0")
	t.testMatches(true, "0-:0:0:0:0:0:0:0", "-ffff:0:0:0:0:0:0:0")
	t.testMatches(true, "ffff:0:0:0:0:0:0:0", "ffff-:0:0:0:0:0:0:0")
	t.testMatches(true, "-0:0:0:0:0:0:0:0", "::")
	t.testMatches(true, "0:0:-:0:0:0:0:0", "0:0:0-ffff:0:0:0:0:0")
	t.testMatches(true, "0:0:0-:0:0:0:0:0", "0:0:-ffff:0:0:0:0:0")
	t.testMatches(true, "0:0:ffff:0:0:0:0:0", "0:0:ffff-:0:0:0:0:0")
	t.testMatches(true, "0:0:-0:0:0:0:0:0", "::")
	t.testMatches(true, "0:-:0:0:0:0:0:0", "0:0-ffff:0:0:0:0:0:0")
	t.testMatches(true, "0:0-:0:0:0:0:0:0", "0:-ffff:0:0:0:0:0:0")
	t.testMatches(true, "0:ffff:0:0:0:0:0:0", "0:ffff-:0:0:0:0:0:0")
	t.testMatches(true, "0:-0:0:0:0:0:0:0", "::")

	t.testMatches(true, "::1:0:0:0.0.0.0", "0:0:0:1::0.0.0.0")

	t.testMatches(true, "1::-1:16", "1::0-1:16")
	//if isNoAutoSubnets {
	//	t.testMatches(true, "1::-1:16/16", "1::0-1:16")
	//	t.testMatches(true, "1::-1:16", "1::0-1:16/16")
	//	t.testMatches(true, "1:-1::16/16", "1:0-1::16")
	//	t.testMatches(true, "1:-1::16", "1:0-1::16/16")
	//} else if false {
	//	t.testMatches(true, "1:-1::16/32", "1:0-1:*")
	//	t.testMatches(true, "1:-1:*", "1:0-1::16/32")
	//} else {
	t.testMatches(true, "1:-1::16/32", "1:0-1::16")
	t.testMatches(true, "1:-1::16", "1:0-1::16/32")
	//}
	t.testMatches(true, "0.0.0.-", "0.0.0.*")           // ok
	t.testMatches(true, "1-.0.0.1-", "1-255.0.0.1-255") // ok // more than one inferred range

	t.testMatches(true, "0b1.0b01.0b101.1-0b11111111", "1.1.5.1-255")
	t.testMatches(true, "0b1.0b01.0b101.0b11110000-0b11111111", "1.1.5.240-255")
	t.testMatches(true, "0b1.0b01.0b101.0b1111____", "1.1.5.240-255")

	t.testMatches(true, "::0b0000111100001111-0b1111000011110000:3", "::f0f-f0f0:3")
	t.testMatches(true, "::0b000011110000____:3", "::f00-f0f:3")
	t.testMatches(true, "::0B000011110000____:3", "::f00-f0f:3")
	t.testMatches(true, "::f0f-0b1111000011110000:3", "::f0f-f0f0:3")
	t.testMatches(true, "::0b0000111100001111-f0f0:3", "::f0f-f0f0:3")
	t.testMatches(true, "::0B0000111100001111-f0f0:3", "::f0f-f0f0:3")

	t.ipv4test(true, "1.2.*.4/1")
	t.ipv4test(false, "1.2.*.4/-1")
	t.ipv4test(false, "1.2.*.4/")
	t.ipv4test(false, "1.2.*.4/x")
	t.ipv4test(true, "1.2.*.4/33") //we are now allowing extra-large prefixes
	t.ipv6test(true, "1:*::1/1")
	t.ipv6test(false, "1:*::1/-1")
	t.ipv6test(false, "1:*::1/")
	t.ipv6test(false, "1:*::1/x")
	t.ipv6test(false, "1:*::1/129") //we are not allowing extra-large prefixes

	//masks that have wildcards in them
	t.ipv4test(false, "1.2.3.4/*")
	t.ipv4test(false, "1.2.*.4/*")
	t.ipv4test(false, "1.2.3.4/1-2.2.3.4")
	t.ipv4test(false, "1.2.*.4/1-2.2.3.4")
	t.ipv4test(false, "1.2.3.4/**")
	t.ipv4test(false, "1.2.*.4/**")
	t.ipv4test(false, "1.2.3.4/*.*")
	t.ipv4test(false, "1.2.*.4/*.*")
	t.ipv4test(false, "1.2.3.4/*:*")
	t.ipv4test(false, "1.2.*.4/*:*")
	t.ipv4test(false, "1.2.3.4/*:*:*:*:*:*:*:*")
	t.ipv4test(false, "1.2.*.4/*:*:*:*:*:*:*:*")
	t.ipv4test(false, "1.2.3.4/1.2.*.4")
	t.ipv4test(false, "1.2.*.4/1.2.*.4")
	t.ipv4test(true, "1.2.*.4/1.2.3.4")
	t.ipv6test(false, "1:2::1/*")
	t.ipv6test(false, "1:*::1/*")
	t.ipv6test(false, "1:2::1/1:1-2:3:4:5:6:7:8")
	t.ipv6test(false, "1:*::1/1:1-2:3:4:5:6:7:8")
	t.ipv6test(false, "1:2::1/**")
	t.ipv6test(false, "1:*::1/**")
	t.ipv6test(false, "1:2::1/*:*")
	t.ipv6test(false, "1:*::1/*:*")
	t.ipv6test(false, "1:2::1/*.*")
	t.ipv6test(false, "1:*::1/*.*")
	t.ipv6test(false, "1:2::1/*.*.*.*")
	t.ipv6test(false, "1:*::1/*.*.*.*")
	t.ipv6test(false, "1:2::1/1:*::2")
	t.ipv6test(false, "1:*::1/1:*::2")
	t.ipv6test(true, "1:*::1/1::2")

	t.ipv4rangetest(true, "1.1.*.100-101", ipaddr.WildcardAndRange)
	t.ipv4rangetest(true, "1.2.*.101-100", ipaddr.WildcardAndRange)   //downwards range
	t.ipv4rangetest(false, "1.2.*.1010-100", ipaddr.WildcardAndRange) //downwards range
	t.ipv4rangetest(true, "1.2.*.101-101", ipaddr.WildcardAndRange)
	t.ipv6rangetest(true, "1:2:f4:a-ff:0-2::1", ipaddr.WildcardAndRange)
	t.ipv6rangetest(true, "1:2:4:ff-a:0-2::1", ipaddr.WildcardAndRange)     //downwards range
	t.ipv6rangetest(false, "1:2:4:ff1ff-a:0-2::1", ipaddr.WildcardAndRange) //downwards range
	t.ipv4rangetest(true, "1.2.*.101-100/24", ipaddr.WildcardAndRange)      //downwards range but covered CIDR

	//these tests create strings that validate ipv4 and ipv6 differently, allowing ranges for one and not the other
	t.ipv4rangestest(true, "1.*.3.4", ipaddr.WildcardAndRange, ipaddr.NoRange)
	t.ipv4rangestest(false, "1.*.3.4", ipaddr.NoRange, ipaddr.WildcardAndRange)
	t.ipv6rangestest(false, "a:*::1.*.3.4", ipaddr.WildcardAndRange, ipaddr.NoRange)
	t.ipv6rangestest(true, "a:*::1.*.3.4", ipaddr.NoRange, ipaddr.WildcardAndRange)
	t.ipv6rangestest(false, "a:*::", ipaddr.WildcardAndRange, ipaddr.NoRange)
	t.ipv6rangestest(true, "a:*::", ipaddr.NoRange, ipaddr.WildcardAndRange)

	//		octal, hex, dec overflow
	//		do it with 1, 2, 3, 4 segments
	t.ipv4_inet_aton_test(true, "0.0.0.1-255")
	t.ipv4_inet_aton_test(false, "0.0.0.1-256")
	t.ipv4_inet_aton_test(true, "0.0.512-65535")
	t.ipv4_inet_aton_test(false, "0.0.512-65536")
	t.ipv4_inet_aton_test(true, "0.65536-16777215")
	t.ipv4_inet_aton_test(false, "0.65536-16777216")
	t.ipv4_inet_aton_test(true, "16777216-4294967295")
	t.ipv4_inet_aton_test(true, "0b00000001000000000000000000000000-4294967295")
	//t.ipv4_inet_aton_test(true, "0b1000000000000000000000000-4294967295");
	t.ipv4_inet_aton_test(false, "16777216-4294967296")
	t.ipv4_inet_aton_test(false, "0.0.0.0x1x")
	t.ipv4_inet_aton_test(false, "0.0.0.1x")
	t.ipv4_inet_aton_test(true, "0.0.0.0x1-0xff")
	t.ipv4_inet_aton_test(false, "0.0.0.0x1-0x100")
	t.ipv4_inet_aton_test(true, "0.0.0xfffe-0xffff")
	t.ipv4_inet_aton_test(false, "0.0.0xfffe-0x10000")
	t.ipv4_inet_aton_test(false, "0.0.0x10000-0x10001")
	t.ipv4_inet_aton_test(true, "0.0-0xffffff")
	t.ipv4_inet_aton_test(false, "0.0-0x1000000")
	t.ipv4_inet_aton_test(true, "0x11000000-0xffffffff")
	t.ipv4_inet_aton_test(false, "0x11000000-0x100000000")
	t.ipv4_inet_aton_test(false, "0x100000000-0x100ffffff")
	t.ipv4_inet_aton_test(true, "0.0.0.00-0377")
	t.ipv4_inet_aton_test(false, "0.0.0.00-0400")
	t.ipv4_inet_aton_test(true, "0.0.0x100-017777")
	t.ipv4_inet_aton_test(false, "0.0.0x100-0200000")
	t.ipv4_inet_aton_test(true, "0.0x10000-077777777")
	//t.ipv4_inet_aton_test(false, "0.0x1-077777777"); the given address throw IncompatibleAddressException as expected, would need to rewrite the test to make that a pass
	t.ipv4_inet_aton_test(false, "0.0x10000-0100000000")
	t.ipv4_inet_aton_test(true, "0x1000000-03777777777")
	t.ipv4_inet_aton_test(true, "0x1000000-037777777777")
	t.ipv4_inet_aton_test(true, "0x1000000-0b11111111111111111111111111111111") //[0-1, 0, 0-255, 0-255]
	t.ipv4_inet_aton_test(false, "0x1000000-040000000000")

	t.ipv4test(true, "*") //toAddress() should not work on this, toAddress(Version) should.

	t.ipv4test2(false, "*%", false, true)  //because the string could represent ipv6, and we are allowing zone, we treat the % as ipv6 zone, and then we invalidate because no zone for ipv4
	t.ipv4test2(false, "*%x", false, true) //no zone for ipv4
	t.ipv4test(true, "**")                 // toAddress() should not work on this, toAddress(Version) should.
	t.ipv6test(true, "**")                 // toAddress() should not work on this, toAddress(Version) should.
	t.ipv6test(true, "*%x")                //ipv6 which allows zone

	t.ipv4test(true, "*.*.*.*") //toAddress() should work on this

	t.ipv4test(true, "1.*.3")

	t.ipv4test(false, "a.*.3.4")
	t.ipv4test(false, "*.a.3.4")
	t.ipv4test(false, "1.*.a.4")
	t.ipv4test(false, "1.*.3.a")

	t.ipv4test(false, ".2.3.*")
	t.ipv4test(false, "1..*.4")
	t.ipv4test(false, "1.*..4")
	t.ipv4test(false, "*.2.3.")

	t.ipv4test(false, "256.*.3.4")
	t.ipv4test(false, "1.256.*.4")
	t.ipv4test(false, "*.2.256.4")
	t.ipv4test(false, "1.*.3.256")

	t.ipv4test(true, "0.0.*.0")
	t.ipv4test(true, "00.*.0.0")
	t.ipv4test(true, "0.00.*.0")
	t.ipv4test(true, "0.*.00.0")
	t.ipv4test(true, "*.0.0.00")
	t.ipv4test(true, "000.0.*.0")
	t.ipv4test(true, "0.000.0.*")
	t.ipv4test(true, "*.0.000.0")
	t.ipv4test(true, "0.0.*.000")

	t.ipv4test(true, "0.0.*.0")
	t.ipv4test(true, "00.*.0.0")
	t.ipv4test(true, "0.00.*.0")
	t.ipv4test(true, "0.*.00.0")
	t.ipv4test(true, "*.0.0.00")
	t.ipv4test(true, "000.0.*.0")
	t.ipv4test(true, "0.000.0.*")
	t.ipv4test(true, "*.0.000.0")
	t.ipv4test(true, "0.0.*.000")

	t.ipv4test(true, "000.000.000.*")

	t.ipv4test(t.isLenient(), "0000.0.*.0")
	t.ipv4test(t.isLenient(), "*.0000.0.0")
	t.ipv4test(t.isLenient(), "0.*.0000.0")
	t.ipv4test(t.isLenient(), "*.0.0.0000")

	t.ipv4test(false, ".0.*.0")
	t.ipv4test(false, "0..*.0")
	t.ipv4test(false, "0.*..0")
	t.ipv4test(false, "*.0.0.")

	t.ipv4test(true, "1.*.3.4/255.1.0.0")
	t.ipv4test(false, "1.*.3.4/255.1.0.0/16")
	t.ipv4test(false, "1.*.3.4/255.*.0.0")   //range in mask
	t.ipv4test(false, "1.*.3.4/255.1-2.0.0") //range in mask
	t.ipv4test(false, "1.*.3.4/1::1")        //mask mismatch
	t.ipv6test(false, "1:*::/1.2.3.4")       //mask mismatch

	t.ipv4test(false, "1.2.3.4/255.*.0.0")   //range in mask
	t.ipv4test(false, "1.2.3.4/255.1-2.0.0") //range in mask
	t.ipv6test(false, "1:2::/1:*::")         //range in mask
	t.ipv6test(false, "1:2::/1:1-2::")       //range in mask

	t.ipv4testOnly(false, "1:2:3:4:5:*:7:8") //fixed
	t.ipv4testOnly(false, "*::1")            //fixed

	t.ipv6test(true, "*")  //toAddress() should not work on this, toAddress(version) should
	t.ipv6test(true, "*%") //toAddress() should not work on this, toAddress(version) should

	t.ipv6test(true, "*:*:*:*:*:*:*:*") //toAddress() should work on this

	t.ipv6test(true, "*::1") // loopback, compressed, non-routable

	t.ipv4test(true, "1.0-0.3.0")
	t.ipv4test(true, "1.0-3.3.0")
	t.ipv4test(true, "1.1-3.3.0")
	t.ipv4test(true, "1-8.1-3.2-4.0-5")

	t.ipv6test(true, "1:0-0:2:0::")
	t.ipv6test(true, "1:0-3:2:0::")
	t.ipv6test(true, "1:1-3:2:0::")
	t.ipv6test(true, "1-fff:1-3:2-4:0-5::")

	t.ipv6test(false, "-:0:0:0:0:0:0:0:0")
	t.ipv6test(true, "-:0:0:0:0:0:0:0") // this is actually equivalent to 0-ffff:0:0:0:0:0:0:0 or 0-:0:0:0:0:0:0:0 or -ffff:0:0:0:0:0:0:0
	t.ipv6test(false, "-:0:0:0:0:0:0")
	t.ipv6test(false, "-:0:0:0:0:0")
	t.ipv6test(false, "-:0:0:0:0")
	t.ipv6test(false, "-:0:0:0")
	t.ipv6test(false, "-:0:0")
	t.ipv6test(false, "-:0")

	t.ipv6test(false, ":-0:0:0:0:0:0:0")
	t.ipv6test(false, ":-0:0:0:0:0:0")
	t.ipv6test(false, ":-0:0:0:0:0")
	t.ipv6test(false, ":-0:0:0:0")
	t.ipv6test(false, ":-0:0:0")
	t.ipv6test(false, ":-0:0")
	t.ipv6test(false, ":-0")

	t.ipv6test(false, "-:1:1:1:1:1:1:1:1")
	t.ipv6test(true, "-:1:1:1:1:1:1:1") // this is actually equivalent to 0-ffff:0:0:0:0:0:0:0 or 0-:0:0:0:0:0:0:0 or -ffff:0:0:0:0:0:0:0
	t.ipv6test(false, "-:1:1:1:1:1:1")
	t.ipv6test(false, "-:1:1:1:1:1")
	t.ipv6test(false, "-:1:1:1:1")
	t.ipv6test(false, "-:1:1:1")
	t.ipv6test(false, "-:1:1")
	t.ipv6test(false, "-:1")

	t.ipv6test(false, ":-1:1:1:1:1:1:1")
	t.ipv6test(false, ":-1:1:1:1:1:1")
	t.ipv6test(false, ":-1:1:1:1:1")
	t.ipv6test(false, ":-1:1:1:1")
	t.ipv6test(false, ":-1:1:1")
	t.ipv6test(false, ":-1:1")
	t.ipv6test(false, ":-1")

	t.ipv6test(true, "::*")                             // unspecified, compressed, non-routable
	t.ipv6test(true, "0:0:*:0:0:0:0:1")                 // loopback, full
	t.ipv6test(true, "0:0:*:0:0:0:0:0")                 // unspecified, full
	t.ipv6test(true, "2001:*:0:0:8:800:200C:417A")      // unicast, full
	t.ipv6test(true, "FF01:*:0:0:0:0:0:101")            // multicast, full
	t.ipv6test(true, "2001:DB8::8:800:200C:*")          // unicast, compressed
	t.ipv6test(true, "FF01::*:101")                     // multicast, compressed
	t.ipv6test(false, "2001:DB8:0:0:8:*:200C:417A:221") // unicast, full
	t.ipv6test(false, "FF01::101::*")                   // multicast, compressed
	t.ipv6test(true, "fe80::217:f2ff:*:ed62")

	t.ipv6test(true, "2001:*:1234:0000:0000:C1C0:ABCD:0876")
	t.ipv6test(true, "3ffe:0b00:0000:0000:0001:0000:*:000a")
	t.ipv6test(true, "FF02:0000:0000:0000:0000:0000:*:0001")
	t.ipv6test(true, "*:0000:0000:0000:0000:0000:0000:0001")
	t.ipv6zerotest(false, "0000:0000:0000:0000:*0000:0000:0000:*0")
	t.ipv6test(t.isLenient(), "02001:*:1234:0000:0000:C1C0:ABCD:0876") // extra 0 not allowed!
	t.ipv6test(false, "2001:0000:1234:0000:0*:C1C0:ABCD:0876")         // extra 0 not allowed!
	t.ipv6test(true, "2001:0000:1234:0000:*:C1C0:ABCD:0876")

	//t.ipv6test(true," 2001:0000:1234:0000:0000:C1C0:ABCD:0876"); // leading space
	//t.ipv6test(true,"2001:0000:1234:0000:0000:C1C0:ABCD:0876 "); // trailing space
	//t.ipv6test(true," 2001:0000:1234:0000:0000:C1C0:ABCD:0876  "); // leading and trailing space

	t.ipv6test(false, "2001:0000:1234:0000:0000:C1C0*:ABCD:0876  0") // junk after valid address
	t.ipv6test(false, "0 2001:0000:123*:0000:0000:C1C0:ABCD:0876")   // junk before valid address
	t.ipv6test(false, "2001:0000:1234: 0000:0000:C1C0:*:0876")       // internal space

	t.ipv6test(true, "3ffe:0b00:*:0001:0000:0000:000a")
	t.ipv6test(false, "3ffe:0b00:1:0001:0000:0000:000a")           // seven segments
	t.ipv6test(false, "FF02:0000:0000:0000:0000:0000:0000:*:0001") // nine segments
	t.ipv6test(false, "3ffe:*::1::a")                              // double "::"
	t.ipv6test(false, "::1111:2222:3333:4444:5555:*::")            // double "::"
	t.ipv6test(true, "2::10")
	t.ipv6test(true, "ff02::1")
	t.ipv6test(true, "fe80:*::")
	t.ipv6test(true, "2002:*::")
	t.ipv6test(true, "2001:*::")
	t.ipv6test(true, "*:0db8:1234::")
	t.ipv6test(true, "::ffff:*:0")
	t.ipv6test(true, "*::1")
	t.ipv6test(true, "1:2:3:4:*:6:7:8")
	t.ipv6test(true, "1:2:*:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::*")
	t.ipv6test(true, "1:2:3:*::8")
	t.ipv6test(true, "1:2:3::8")
	t.ipv6test(true, "*:2::8")
	t.ipv6test(true, "1::*")
	t.ipv6test(true, "*::2:3:4:5:6:7")
	t.ipv6test(true, "*::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:*")
	t.ipv6test(true, "1::2:*:4")
	t.ipv6test(true, "1::*:3")
	t.ipv6test(true, "1::*")

	t.ipv6test(true, "::*:3:4:5:6:7:8")
	t.ipv6test(true, "*::2:3:4:5:6:7")
	t.ipv6test(true, "::*:3:4:5:6")
	t.ipv6test(true, "::*:3:4:5")
	t.ipv6test(true, "::2:3:*")
	t.ipv6test(true, "*::2:3")
	t.ipv6test(true, "::*")
	t.ipv6test(true, "1:*:3:4:5:6::")
	t.ipv6test(true, "1:2:3:4:*::")
	t.ipv6test(true, "1:2:3:*::")
	t.ipv6test(true, "1:2:3::*")
	t.ipv6test(true, "*:2::")
	t.ipv6test(true, "*::")
	t.ipv6test(true, "*:2:3:4:5::7:8")
	t.ipv6test(false, "1:2:3::4:5::7:*") // Double "::"
	t.ipv6test(false, "12345::6:7:*")
	t.ipv6test(true, "1:2:3:4::*:*")
	t.ipv6test(true, "1:*:3::7:8")
	t.ipv6test(true, "*:*::7:8")
	t.ipv6test(true, "*::*:8")

	// Testing IPv4 addresses represented as dotted-quads
	// Leading zero's in IPv4 addresses not allowed: some systems treat the leading "0" in ".086" as the start of an octal number
	// Update: The BNF in RFC-3986 explicitly defines the dec-octet (for IPv4 addresses) not to have a leading zero
	//t.ipv6test(false,"fe80:0000:0000:*:0204:61ff:254.157.241.086");
	t.ipv6test(!t.isLenient(), "fe80:0000:0000:*:0204:61ff:254.157.241.086")
	t.ipv6test(true, "::*:192.0.128.*")
	t.ipv6test(false, "XXXX:XXXX:XXXX:XXXX:XXXX:XXXX:1.2.3.4")
	t.ipv6test(true, "1111:2222:*:4444:5555:6666:00.00.00.00")
	t.ipv6test(true, "1111:2222:3333:4444:5555:6666:000.*.000.000")
	t.ipv6test(false, "*:2222:3333:4444:5555:6666:256.256.256.256")

	t.ipv6test(true, "*:2222:3333:4444:5555:6666:123.123.123.123")
	t.ipv6test(true, "1111:*:3333:4444:5555::123.123.123.123")
	t.ipv6test(true, "1111:2222:*:4444::123.123.123.123")
	t.ipv6test(true, "1111:2222:3333::*.*.123.123")
	t.ipv6test(true, "1111:2222::123.123.*.*")
	t.ipv6test(true, "1111:2222::123.123.123.*")
	t.ipv6test(true, "1111::123.*.123.123")
	t.ipv6test(true, "::123.123.123.*")
	t.ipv6test(true, "1111:2222:3333:4444::*:123.123.123.123")
	t.ipv6test(true, "1111:2222:*::6666:123.123.123.123")
	t.ipv6test(true, "*:2222::6666:123.123.123.123")
	t.ipv6test(true, "1111::6666:*.*.*.*")
	t.ipv6test(true, "::6666:123.123.2.123")
	t.ipv6test(true, "1111:*:3333::5555:6666:123.*.123.123")
	t.ipv6test(true, "1111:2222::*:6666:123.123.*.*")
	t.ipv6test(true, "1111::*:6666:*.*.123.123")
	t.ipv6test(true, "1111::*:6666:*.0-255.123.123") //1111::*:6666:*.123.123
	t.ipv6test(true, "::5555:6666:123.123.123.123")
	t.ipv6test(true, "1111:2222::4444:5555:*:123.123.123.123")
	t.ipv6test(true, "1111::4444:5555:6666:123.*.123.123")
	t.ipv6test(true, "*::4444:5555:6666:123.123.123.123")
	t.ipv6test(true, "1111::*:4444:5555:6666:123.123.123.123")
	t.ipv6test(true, "::2222:*:4444:5555:6666:123.123.123.123")
	t.ipv6test(true, "::*:*:*:*:*:*.*.*.*")
	t.ipv6test(true, "*::*:*:*:*:*.*.*.*")
	t.ipv6test(false, "*:::*:*:*:*.*.*.*")
	t.ipv6test(false, "*:::*:*:*:*:*.*.*.*")
	t.ipv6test(true, "*::*:*:*:*:*.*.*.*")
	t.ipv6test(false, "*::*:*:*:*:*:*.*.*.*")
	t.ipv6test(false, "*:*:*:*:*:*:*:*:*.*.*.*")
	t.ipv6test(false, "*:*:*:*:*:*:*::*.*.*.*")
	t.ipv6test(false, "*:*:*:*:*:*::*:*.*.*.*")
	t.ipv6test(true, "*:*:*:*:*:*:*.*.*.*")
	t.ipv6test(true, "*:*:*:*:*::*.*.*.*")
	t.ipv6test(true, "*:*:*:*::*:*.*.*.*")

	t.ipv6test(true, "::*")
	t.ipv6test(true, "*:0:0:0:0:0:0:*")

	// Additional cases: http://crisp.tweakblogs.net/blog/2031/ipv6-validation-%28and-caveats%29.html
	t.ipv6test(true, "0:a:b:*:d:e:f::")
	t.ipv6test(true, "::0:a:*:*:d:e:f") // syntactically correct, but bad form (::0:... could be combined)
	t.ipv6test(true, "a:b:c:*:*:f:0::")
	t.ipv6test(false, "':10.*.0.1")

	t.ipv4test(true, "1.*.4")
	t.ipv4test(true, "1.2.*")
	t.ipv4test(true, "*.1")
	t.ipv4test(true, "1.*")
	t.ipv4test(true, "1.*.1")
	t.ipv4test(true, "1.*.*")
	t.ipv4test(true, "*.*.1")
	t.ipv4test(true, "*.1.*")
	t.ipv4test(t.isLenient(), "1")
	t.ipv4test(t.isLenient(), "1.1")
	t.ipv4test(t.isLenient(), "1.1.1")

	t.ipv4test(true, "*.1.2.*")
	t.ipv4test(true, "*.1.*.2")
	t.ipv4test(true, "*.*.*.2")
	t.ipv4test(true, "*.*.*.*")
	t.ipv4test(true, "1.*.2.*")
	t.ipv4test(true, "1.2.*.*")

	t.ipv4test(true, "*.*")
	t.ipv6test(true, "1::1.2.*")
	t.ipv6test(true, "1::1.2.**")
	t.ipv6test(false, "1::1.2.**z")
	t.ipv6test(true, "1::1.2.3.4")
	t.ipv6test(true, "1:*:1")
	t.ipv4test(true, "1.2.*")

	t.ipv4test(false, "%.%")
	t.ipv6test(false, "1::1.2.%")
	t.ipv6test(true, "1::1.2.*%")
	t.ipv6test(true, "1::1.2.*%z")
	t.ipv6test(false, "1:%:1")

	t.ipv6test(true, "1::%-.1")
	t.ipv6test(true, "1::%-.1/16") //that is a zone of "-." and a prefix of 16
	t.ipv6test(true, "1::%-1/16")  //that is a zone of "-" and a prefix of 16
	t.ipv6test(true, "1::-1:16")   //that is just an address with a ranged segment 0-1

	t.ipv6test(true, "1::%-.1-16")  // -.1-16 is the zone
	t.ipv6test(true, "1::%-.1/16")  //we treat /16 as prefix length
	t.ipv6test(false, "1::%-.1:16") //we reject ':' as part of zone
	t.ipv6test(false, "1::%-.1/1a") //prefix has 'a'
	t.ipv6test(false, "1::%-1/1a")  //prefix has 'a'
	t.ipv6test(true, "1::%%1")      //zone has '%'
	t.ipv6test(true, "1::%%1/16")   //zone has '%'
	t.ipv6test(true, "1::%%ab")     //zone has '%'
	t.ipv6test(true, "1::%%ab/16")  //zone has '%'
	t.ipv6test(true, "1::%$1")      //zone has '$'
	t.ipv6test(true, "1::%$1/16")   //zone has '$'

	t.ipv4test(true, "1.2.%") //we allow this now, the % is seen as a wildcard because we are ipv4 - if we allow zone and the string can be interpreted as ipv6 then % is seen as zone character

	t.ipv6test(true, "1:*")
	t.ipv6test(true, "*:1:*")
	t.ipv6test(true, "*:1")

	//t.ipv6test(true, "*:1:1.*.1");//cannot be converted to ipv6 range
	t.ipv6test(true, "*:1:1.*.*")
	//t.ipv6test(true, "*:1:*.1");//cannot be converted to ipv6 range
	t.ipv6test(true, "*:1:*.0-255.1.1")
	t.ipv6test(true, "*:1:1.*")

	t.ipv6test(false, "1:1:1.*.1")
	t.ipv6test(false, "1:1:1.*.1.1")
	t.ipv6test(true, "1:1:*.*")
	t.ipv6test(true, "1:2:3:4:5:*.*")
	t.ipv6test(true, "1:2:3:4:5:6:*.*")
	t.ipv6test(false, "1:1:1.*")

	t.ipv6test(true, "1::1:1.*.*")
	t.ipv6test(true, "1::1:*.*.1.1")
	t.ipv6test(true, "1::1:1.*")

	t.ipv6test(true, "1:*.*.*.*") //in this one, the wildcard covers both ipv6 and ipv4 parts
	t.ipv6test(true, "1::*.*.*.*")
	t.ipv6test(true, "1:*.*.1.2")    //in this one, the wildcard covers both ipv6 and ipv4 parts
	t.ipv6test(true, "1::*.*.1.2")   //compression takes precedence so the wildcard does not cover both ipv6 and ipv4 parts
	t.ipv6test(true, "1::2:*.*.1.2") //compression takes precedence so the wildcard does not cover both ipv6 and ipv4 parts
	t.ipv6test(true, "::2:*.*.1.2")  //compression takes precedence so the wildcard does not cover both ipv6 and ipv4 parts
	t.ipv6test(false, "1:1.*.2")
	t.ipv6test(false, "1:1.*.2.2")
	t.ipv6test(t.isLenient(), "1:*:1.2")

	t.ipv6test(true, "*:1:1.*")
	t.ipv6test(t.isLenient(), "*:1:1.2.3")
	t.ipv6test(true, "::1:1.*")
	t.ipv6test(t.isLenient(), "::1:1.2.3")

	t.ipv6test(true, "1:*:1")
	t.ipv6test(true, "1:*:1:1.1.*")
	t.ipv6test(true, "1:*:1:1.1.*.*")
	t.ipv6test(true, "1:*:1:*")
	t.ipv6test(true, "1:*:1:*.*.1.2")
	t.ipv6test(true, "1:*:1:1.*")
	t.ipv6test(t.isLenient(), "1:*:1:1.2.3")

	t.ipv6test(false, "1:*:1:2:3:4:5:6:7")
	t.ipv6test(false, "1:*:1:2:3:4:5:1.2.3.4")
	t.ipv6test(true, "1:*:2:3:4:5:1.2.3.4")
	t.ipv6test(false, "1:*:2:3:4:5:1.2.3.4.5")
	t.ipv6test(false, "1:1:2:3:4:5:1.2.3.4.5")
	t.ipv6test(false, "1:1:2:3:4:5:6:1.2.3.4")
	t.ipv6test(false, "1:1:2:3:4:5:6:1.*.3.4")
	t.ipv6test(true, "1:2:3:4:5:6:1.2.3.4")
	t.ipv6test(true, "1:2:3:4:5:6:1.*.3.4")

	t.ipv4test(true, "255._.3.4")
	t.ipv4test(true, "1.255._.4")
	t.ipv4test(true, "_.2.255.4")
	t.ipv4test(true, "1._.3.255")

	t.ipv4test(true, "255.__.3.4")
	t.ipv4test(true, "1.255.__.4")
	t.ipv4test(true, "__.2.255.4")
	t.ipv4test(true, "1.__.3.255")

	t.ipv4test(true, "255.___.3.4")
	t.ipv4test(true, "1.255.___.4")
	t.ipv4test(true, "___.2.255.4")
	t.ipv4test(true, "1.___.3.255")

	t.ipv4test(t.isLenient(), "255.____.3.4")
	t.ipv4test(t.isLenient(), "1.255.____.4")
	t.ipv4test(t.isLenient(), "____.2.255.4")
	t.ipv4test(t.isLenient(), "1.____.3.255")

	t.ipv4test(false, "255._2_.3.4")
	t.ipv4test(false, "1.255._2_.4")
	t.ipv4test(false, "_2_.2.255.4")
	t.ipv4test(false, "1._2_.3.255")

	t.ipv4test(true, "255.2__.3.4")
	t.ipv4test(true, "1.255.2__.4")
	t.ipv4test(true, "2__.2.255.4")
	t.ipv4test(true, "1.2__.3.255")

	t.ipv4test(true, "255.2_.3.4")
	t.ipv4test(true, "1.255.2_.4")
	t.ipv4test(true, "2_.2.255.4")
	t.ipv4test(true, "1.2_.3.255")

	t.ipv4test(false, "255.__2.3.4")
	t.ipv4test(false, "1.255.__2.4")
	t.ipv4test(false, "__2.2.255.4")
	t.ipv4test(false, "1.__2.3.255")

	t.ipv4test(true, "25_.__.3.4")
	t.ipv4test(true, "1.255.2__._")
	t.ipv4test(true, "2_.2_.255.__")
	t.ipv4test(false, "1.2__.3__.25_")
	t.ipv4test(true, "1.2__.3_.25_")
	t.ipv4test(true, "1.2__.2__.25_")

	t.ipv4test(false, "1.1--2.1.1")
	t.ipv4test(false, "1.1-2-3.1.1")
	t.ipv4test(false, "1.1-2-.1.1")
	t.ipv4test(false, "1.-1-2.1.1")

	t.ipv4test(false, "1.1_2_.1.1")
	t.ipv4test(false, "1.1_2.1.1")
	t.ipv4test(true, "1.1_.1.1")
	t.ipv4test(false, "1.1_-2.1.1")
	t.ipv4test(false, "1.1-2_.1.1")
	t.ipv4test(false, "1.1*-2.1.1")
	t.ipv4test(false, "1.1-2*.1.1")
	t.ipv4test(false, "1.*1-2.1.1")
	t.ipv4test(false, "1.1-*2.1.1")
	t.ipv4test(false, "1.*-2.1.1")
	t.ipv4test(false, "1.1-*.1.1")

	t.ipv6test(false, "1:1--2:1:1::")
	t.ipv6test(false, "1:1-2-3:1:1::")
	t.ipv6test(false, "1:1-2-:1:1::")
	t.ipv6test(false, "1:-1-2:1:1::")

	t.ipv6test(false, "1:1_2_:1.1::")
	t.ipv6test(false, "1:1_2:1:1::")
	t.ipv6test(true, "1:1_:1:1::")

	t.ipv6test(false, "1:1_-2:1:1::")
	t.ipv6test(false, "1:1-2_:1:1::")
	t.ipv6test(false, "1:1-_2:1:1::")
	t.ipv6test(false, "1:1*-2:1:1::")
	t.ipv6test(false, "1:1-2*:1:1::")
	t.ipv6test(false, "1:*-2:1:1::")
	t.ipv6test(false, "1:1-*:1:1::")
	t.ipv6test(false, "1:*1-2:1:1::")
	t.ipv6test(false, "1:1-*2:1:1::")

	//double -
	// _4_ single char wildcards not in trailing position

	t.ipv6test(true, "::ffff:_:0")
	t.ipv6test(true, "_::1")
	t.ipv6test(true, "1:2:3:4:_:6:7:8")
	t.ipv6test(true, "1:2:_:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::_")
	t.ipv6test(true, "1:2:3:_::8")
	t.ipv6test(true, "_:2::8")
	t.ipv6test(true, "1::_")
	t.ipv6test(true, "_::2:3:4:5:6:7")
	t.ipv6test(true, "_::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:_")
	t.ipv6test(true, "1::2:_:4")
	t.ipv6test(true, "1::_:3")
	t.ipv6test(true, "1::_")

	t.ipv6test(true, "::ffff:__:0")
	t.ipv6test(true, "__::1")
	t.ipv6test(true, "1:2:3:4:__:6:7:8")
	t.ipv6test(true, "1:2:__:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::__")
	t.ipv6test(true, "1:2:3:__::8")
	t.ipv6test(true, "__:2::8")
	t.ipv6test(true, "1::__")
	t.ipv6test(true, "__::2:3:4:5:6:7")
	t.ipv6test(true, "__::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:__")
	t.ipv6test(true, "1::2:__:4")
	t.ipv6test(true, "1::__:3")
	t.ipv6test(true, "1::__")

	t.ipv6test(true, "::ffff:___:0")
	t.ipv6test(true, "___::1")
	t.ipv6test(true, "1:2:3:4:___:6:7:8")
	t.ipv6test(true, "1:2:___:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::___")
	t.ipv6test(true, "1:2:3:___::8")
	t.ipv6test(true, "___:2::8")
	t.ipv6test(true, "1::___")
	t.ipv6test(true, "___::2:3:4:5:6:7")
	t.ipv6test(true, "___::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:___")
	t.ipv6test(true, "1::2:___:4")
	t.ipv6test(true, "1::___:3")
	t.ipv6test(true, "1::___")

	t.ipv6test(true, "::ffff:____:0")
	t.ipv6test(true, "____::1")
	t.ipv6test(true, "1:2:3:4:____:6:7:8")
	t.ipv6test(true, "1:2:____:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::____")
	t.ipv6test(true, "1:2:3:____::8")
	t.ipv6test(true, "____:2::8")
	t.ipv6test(true, "1::____")
	t.ipv6test(true, "____::2:3:4:5:6:7")
	t.ipv6test(true, "____::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:____")
	t.ipv6test(true, "1::2:____:4")
	t.ipv6test(true, "1::____:3")
	t.ipv6test(true, "1::____")

	t.ipv6test(false, "::ffff:_____:0")
	t.ipv6test(false, "_____::1")
	t.ipv6test(false, "1:2:3:4:_____:6:7:8")
	t.ipv6test(false, "1:2:_____:4:5:6::8")
	t.ipv6test(false, "1:2:3:4:5::_____")
	t.ipv6test(false, "1:2:3:_____::8")
	t.ipv6test(false, "_____:2::8")
	t.ipv6test(false, "1::_____")
	t.ipv6test(false, "_____::2:3:4:5:6:7")
	t.ipv6test(false, "_____::2:3:4:5:6")
	t.ipv6test(false, "1::2:3:4:_____")
	t.ipv6test(false, "1::2:_____:4")
	t.ipv6test(false, "1::_____:3")
	t.ipv6test(false, "1::_____")

	t.ipv6test(false, "::ffff:ff___:0")
	t.ipv6test(false, "f____::1")
	t.ipv6test(false, "1:2:3:4:ffff_:6:7:8")
	t.ipv6test(false, "1:2:ffff_:4:5:6::8")
	t.ipv6test(false, "1:2:3:4:5::f_f__")
	t.ipv6test(false, "1:2:3:fff__::8")
	t.ipv6test(false, "f___f:2::8")
	t.ipv6test(false, "1::ff_ff")
	t.ipv6test(false, "ff_ff::2:3:4:5:6:7")
	t.ipv6test(false, "f____::2:3:4:5:6")
	t.ipv6test(false, "1::2:3:4:F____")
	t.ipv6test(false, "1::2:FF___:4")
	t.ipv6test(false, "1::FFF__:3")
	t.ipv6test(false, "1::FFFF_")

	t.ipv6test(false, "::ffff:_2_:0")
	t.ipv6test(false, "_2_::1")
	t.ipv6test(false, "1:2:3:4:_2_:6:7:8")
	t.ipv6test(false, "1:2:_2_:4:5:6::8")
	t.ipv6test(false, "1:2:3:4:5::_2_")
	t.ipv6test(false, "1:2:3:_2_::8")
	t.ipv6test(false, "_2_:2::8")
	t.ipv6test(false, "1::_2_")
	t.ipv6test(false, "_2_::2:3:4:5:6:7")
	t.ipv6test(false, "_2_::2:3:4:5:6")
	t.ipv6test(false, "1::2:3:4:_2_")
	t.ipv6test(false, "1::2:_2_:4")
	t.ipv6test(false, "1::_2_:3")
	t.ipv6test(false, "1::_2_")

	t.ipv6test(false, "::ffff:_2:0")
	t.ipv6test(false, "_2::1")
	t.ipv6test(false, "1:2:3:4:_2:6:7:8")
	t.ipv6test(false, "1:2:_2:4:5:6::8")
	t.ipv6test(false, "1:2:3:4:5::_2")
	t.ipv6test(false, "1:2:3:_2::8")
	t.ipv6test(false, "_2:2::8")
	t.ipv6test(false, "1::_2")
	t.ipv6test(false, "_2::2:3:4:5:6:7")
	t.ipv6test(false, "_2::2:3:4:5:6")
	t.ipv6test(false, "1::2:3:4:_2")
	t.ipv6test(false, "1::2:_2:4")
	t.ipv6test(false, "1::_2:3")
	t.ipv6test(false, "1::_2")

	t.ipv6test(true, "::ffff:2_:0")
	t.ipv6test(true, "2_::1")
	t.ipv6test(true, "1:2:3:4:2_:6:7:8")
	t.ipv6test(true, "1:2:2_:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::2_")
	t.ipv6test(true, "1:2:3:2_::8")
	t.ipv6test(true, "2_:2::8")
	t.ipv6test(true, "1::2_")
	t.ipv6test(true, "2_::2:3:4:5:6:7")
	t.ipv6test(true, "2_::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:2_")
	t.ipv6test(true, "1::2:2_:4")
	t.ipv6test(true, "1::2_:3")
	t.ipv6test(true, "1::2_")

	t.ipv6test(true, "::ffff:2___:0")
	t.ipv6test(true, "2___::1")
	t.ipv6test(true, "1:2:3:4:2___:6:7:8")
	t.ipv6test(true, "1:2:2___:4:5:6::8")
	t.ipv6test(true, "1:2:3:4:5::2___")
	t.ipv6test(true, "1:2:3:2___::8")
	t.ipv6test(true, "2___:2::8")
	t.ipv6test(true, "1::2___")
	t.ipv6test(true, "2___::2:3:4:5:6:7")
	t.ipv6test(true, "2___::2:3:4:5:6")
	t.ipv6test(true, "1::2:3:4:2___")
	t.ipv6test(true, "1::2:2___:4")
	t.ipv6test(true, "1::2___:3")
	t.ipv6test(true, "1::2___")

	t.ipv6test(true, "::fff_:2___:0")
	t.ipv6test(true, "2___::_")
	t.ipv6test(true, "1:2:3:4:2___:6_:7_:8")
	t.ipv6test(true, "1:2:2___:4:5:6::8__")
	t.ipv6test(true, "1:2:3_:4:5::2___")
	t.ipv6test(true, "1:2:3:2___::8")
	t.ipv6test(true, "2___:2::8")
	t.ipv6test(true, "1::2___")
	t.ipv6test(true, "2___::2_:3__:4:5:6:7")
	t.ipv6test(true, "2___::2:3_:4:5:6")
	t.ipv6test(true, "1::2:3:4_:2___")
	t.ipv6test(true, "1::2:2___:4f__")
	t.ipv6test(true, "1___::2___:3___")
	t.ipv6test(true, "1_::2___")

	t.ipv6test(t.isLenient(), "*:1:1._.__")
	t.ipv6test(true, "*:1:1._.__.___")
	//t.ipv6test(false, "*:_:1:_.1.1._");//this passes validation but conversion to mask fails because the ipv4 ranges cannot be converted to ipv6 ranges
	t.ipv6test(true, "*:_:1:1._.1._")
	t.ipv6test(true, "*:_:1:_.___.1._")
	t.ipv6test(true, "*:_:1:_.___._.___")
	t.ipv6test(true, "1:*:1_:1:1.1_.1.1")

	t.ipv6test(false, "1:1:1.2_.1")
	t.ipv6test(false, "1:1:1.2__.1.1")
	t.ipv6test(false, "1:1:_.*")
	t.ipv6test(false, "1:1:1._")

	t.ipv6test(true, "a-f:b:c:d:e:f:a:bb")
	t.ipv6test(true, "-f:b:c:d:e:f:a:bb")

	t.testCIDRSubnets("9.*.237.26/0", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/1", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/4", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/5", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/7", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/8", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/9", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/16", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/30", "9.*.237.26")
	t.testCIDRSubnets("9.*.237.26/31", "9.*.237.26-27")
	t.testCIDRSubnets("9.*.237.26/32", "9.*.237.26")

	t.testContains("0.0.0.0/0", "1-2.*.3.*", false)

	t.testContains("0-127.0.0.0/8", "127-127.*.3.*", false)
	t.testContains("0.0.0.0/4", "13-15.*.3.*", false)
	t.testContains("0-15.*.*.*/4", "13-15.*.3.*", false)
	t.testContains("0.0.0.0/4", "9.*.237.*/16", false)
	t.testContains("0.0.0.0/4", "8-9.*.237.*/16", false)

	t.testNotContains("1-2.0.0.0/4", "9.*.237.*/16")
	t.testNotContains("1-2.0.0.0/4", "8-9.*.237.*/16")

	t.testNotContains("1-2.0.0.0/4", "9-17.*.237.*/16")
	t.testContains("8.0.0.0/5", "15.2.3.4", false)
	t.testContains("8.0.0.0/7", "8-9.*.3.*", false)
	t.testContains("9.0.0.0/8", "9.*.3.*", false)
	t.testContains("9.128.0.0/9", "9.128-255.*.0", false)
	t.testContains("9.128.0.0/15", "9.128-129.3.*", false)
	t.testContains("9.129.0.0/16", "9.129.3.*", false)
	t.testNotContains("9.129.0.0/16", "9.128-129.3.*")
	t.testNotContains("9.129.0.0/16", "9.128.3.*")
	t.testContains("9.129.237.24/30", "9.129.237.24-27", true)
	t.testContains("9.129.237.24/30", "9.129.237.24-27/31", true)
	t.testContains("9.129.237.24-27/30", "9.129.237.24-27/31", true)

	t.testContains("*.*.*.*/0", "9.129.237.26/0", false)
	t.testContains("0.0.0.0/0", "*.*.*.*/0", true)
	t.testContains("0.0.0.0/4", "0-15.0.0.*/4", false)
	t.testNotContains("192.0.0.0/4", "0-15.0.0.*/4")

	t.testNotContains("0-127.129.237.26/1", "0-127.0.*.0/1")
	t.testNotContains("9.129.237.26/0", "*.*.*.1/0")
	t.testNotContains("9.129.237.26/4", "0-15.0.1.*/4")
	t.testNotContains("1-16.0.0.*/4", "9.129.237.26/4")
	t.testNotContains("9.129.237.26/5", "8-15.0.0.0/5")
	t.testNotContains("9.129.237.26/7", "8-9.0.0.1-3/7")
	t.testNotContains("7-9.0.0.1-3/7", "9.129.237.26/7")
	t.testNotContains("9.129.237.26/8", "9.*.0.0/8")
	t.testNotContains("9.129.237.26/9", "9.128-255.0.0/9")
	t.testNotContains("9.129.237.26/15", "9.128-129.0.*/15")
	t.testNotContains("9.129.237.26/16", "9.129.*.1/16")
	t.testNotContains("9.129.237.26/30", "9.129.237.27/30")

	t.testContains("0.0.0.0/4", "9.129.237.26/4", false)
	t.testContains("8.0.0.0/5", "8-15.0.0.0/5", false)
	t.testContains("8.0.0.0/7", "8-9.0.0.1-3/7", false)
	t.testContains("7-9.*.*.*/7", "9.129.237.26/7", false)
	t.testContains("9.0.0.0/8", "9.*.0.0/8", false)
	t.testContains("9.128.0.0/9", "9.128-255.0.0/9", false)
	t.testContains("9.128.0.0/15", "9.128-129.0.*/15", false)
	t.testContains("9.128.0.0/15", "9.128.0.0/15", true)
	t.testContains("9.129.0.0/16", "9.129.*.*/16", true)
	t.testContains("9.128-129.*.*/15", "9.128.0.0/15", true)
	t.testContains("9.128.*.*/16", "9.128.0.0/16", true)
	t.testContains("9.129.*.*/16", "9.129.*.*/16", true)
	t.testContains("9.129.*.*/16", "9.129.*.0/16", false)
	t.testContains("9.129.237.24/30", "9.129.237.24-27/30", true)
	t.testContains("9.128-129.*.26/32", "9.128-129.*.26/32", true)

	t.testNotContains("1-16.0.0.0/4", "9.129.237.26/4")
	t.testNotContains("9.129.237.26/5", "8-15.0.0.0/5")
	t.testNotContains("9.129.237.26/7", "8-9.0.0.1-3/7")
	t.testNotContains("7-9.0.0.1-3/7", "9.129.237.26/7")
	t.testNotContains("9.129.237.26/8", "9.*.0.0/8")
	t.testNotContains("9.129.237.26/9", "9.128-255.0.0/9")
	t.testNotContains("9.129.237.26/15", "9.128-129.0.*/15")
	t.testNotContains("9.129.237.26/16", "9.129.*.1/16")
	t.testNotContains("9.129.237.26/16", "9.129.1.*/16")
	t.testNotContains("9.129.237.25/30", "9.129.237.26/30")

	t.testContains("1-16.0.0.*/4", "9.0.0.*/4", false)
	t.testNotContains("1-16.0.0.0-254/4", "9.0.0.*/4")
	t.testContains("0-16.0.0.0/4", "9.0.0.*/4", false)
	t.testContains("8-15.129.237.26/5", "9.129.237.26/5", false)
	t.testContains("8-9.129.237.26/7", "9.129.237.26/7", false)
	t.testContains("7-9.0.0.1-3/7", "9.0.0.2/7", false)
	t.testContains("9.*.237.26/8", "9.*.237.26/8", true)
	t.testContains("9.128-255.237.26/9", "9.129.237.26/9", false)
	t.testContains("9.128-129.237.26/15", "9.129.237.26/15", false)
	t.testContains("9.129.*.*/16", "9.129.237.26/16", false)
	t.testContains("9.129.237.24-27/30", "9.129.237.26/30", false)
	t.testContains("9.128-129.*.26/32", "9.128-129.*.26/32", true)

	t.testNotContains("9.129.237.26/4", "16-17.0.0.*/4")
	t.testNotContains("9.129.237.26/7", "2.0.0.1-3/7")

	t.testContains("::ffff:1.*.3.4", "1.2.3.4", false) //ipv4 mapped

	t.testNotContains("::ffff:1.2-4.3.4/112", "1.2-3.3.*")
	t.testNotContains("ffff:0:0:0:0:0:*:0/32", "ffff:0:ffff:1-d:e:f:*:b")
	t.testNotContains("fffc-ffff::ffff/30", "fffd-fffe:0:0:0:0:0:0:0/30")
	t.testNotContains("ffff:0-d::ffff/32", "ffff:a-c:0:0:0:0:0:0/32")
	t.testNotContains("ffff::ffff/0", "a-b:0:b:0:c:d-e:*:0/0")
	t.testNotContains("ffff::ffff/1", "8000-8fff:0:0:0:0:*:a-b:0/1")
	t.testNotContains("ffff:*::fffb/126", "ffff:*:0:0:0:0:0:fffc-ffff/126")
	t.testNotContains("ffff:1-2::fffb/126", "ffff:1-2:0:0:0:0:0:fffc-ffff/126")

	t.testContains("::ffff:1.2-4.0.0/112", "1.2-3.3.*", false)

	t.testContains("0:0:0:0:0:0:0:0/0", "a:*:c:d:e:1-ffff:a:b", false)
	t.testContains("8000:0:0:0:0:0:0:0/1", "8000-8fff:b:c:d:e:f:*:b", false)
	t.testNotContains("8000:0:0:0:0:0:0:0/1", "7fff-8fff:b:c:d:e:f:*:b")
	t.testContains("ffff:0:0:0:0:0:0:0/30", "ffff:0-3:c:d:e:f:a:b", false)
	t.testNotContains("ffff:0:0:0:0:0:0:0/30", "ffff:0-4:c:d:e:f:a:b")

	t.testContains("ffff:0:0:0:0:0:0:0/32", "ffff:0:ffff:1-d:e:f:*:b", false)
	t.testContains("fffc-ffff::/30", "fffd-fffe:0:0:0:0:0:0:0/30", false)
	t.testContains("ffff:0-d::/32", "ffff:a-c:0:0:0:0:0:0/32", false)

	t.testNotContains("ffff:0:0:0:0:1-2:0:0/32", "ffff:0-1:ffff:d:e:f:a:b")
	t.testContains("ffff:0:0:0:0:4-ffff:0:fffc-ffff", "ffff:0:0:0:0:4-ffff:0:fffd-ffff", false)
	t.testContains("ffff:0:0:0:0:4-ffff:0:fffc/126", "ffff:0:0:0:0:4-ffff:0:fffd-ffff", false)
	t.testContains("ffff:0:0:0:0:4-ffff:0:fffc/126", "ffff:0:0:0:0:4-ffff:0:fffc-ffff", true)
	t.testContains("ffff:0:*:0:0:4-ffff:0:ffff/128", "ffff:0:*:0:0:4-ffff:0:ffff", true)

	t.testContains("ffff:*:0:0:0:0:0:fffa-ffff/126", "ffff:*::ffff/126", false)

	t.testContains("::/0", "a-b:0:b:0:c:d-e:*:0/0", false)
	t.testContains("8000::/1", "8000-8fff:0:0:0:0:*:a-b:0/1", false)
	t.testContains("ffff:*::fffc/126", "ffff:*:0:0:0:0:0:fffc-ffff/126", true)
	t.testContains("ffff:1-2::fffc/126", "ffff:1-2:0:0:0:0:0:fffc-ffff/126", true)

	t.testContains("10.162.155.1-255", "10.162.155.1-51", false)
	t.testContains("10.162.155.1-51", "10.162.155.1-51", true)
	t.testContains("10.162.1-51.155", "10.162.1-51.155", true)
	t.testContains("10.162.1-255.155", "10.162.1-51.155", false)
	t.testContains("1-255.10.162.155", "1-51.10.162.155", false)

	t.testContains("10.162.155.0-255", "10.162.155.0-51", false)
	t.testContains("10.162.155.0-51", "10.162.155.0-51", true)
	t.testContains("10.162.0-51.155", "10.162.0-51.155", true)
	t.testContains("10.162.0-255.155", "10.162.0-51.155", false)
	t.testContains("0-255.10.162.155", "0-51.10.162.155", false)

	t.testNotContains("192.13.1.0/25", "192.13.1.1-255")
	t.testNotContains("192.13.1.1-255", "192.13.1.0/25")

	t.testContains("192.13.1.0/25", "192.13.1.1-127", false)
	t.testContains("192.13.1.0/25", "192.13.1.0-127", true)

	t.testContains("192.13.1.0-127", "192.13.1.0/25", true)

	t.testContains("ffff:1-3::/32", "ffff:2::", false)
	t.testContains("ffff:2-3::/32", "ffff:2::", false)
	t.testContains("ffff:1-3::/32", "ffff:3::", false)

	t.testNotContains("ffff:1-3::/32", "ffff:4::")

	t.testContains("ffff:1000-3000::/20", "ffff:2000::", false)
	t.testContains("ffff:2000-3000::/20", "ffff:2000::", false)
	t.testContains("ffff:1000-3000::/20", "ffff:3000::", false)

	t.testNotContains("ffff:1000-3000::/20", "ffff:4000::")
	t.testNotContains("ffff:2000-3000::/20", "ffff:4000::")

	t.testContains("ffff:1000::/20", "ffff:1111-1222::", false)
	t.testNotContains("ffff:1000::/20", "ffff:1-::")

	t.testContains("ffff:1-:*", "ffff:1000::/20", false)
	t.testNotContains("ffff:1000::/20", "ffff:1111-2222::")
	t.testNotContains("ffff:1000::/20", "ffff:1-10::")
	t.testNotContains("ffff:1000::/20", "ffff:1-1::")

	t.testContains("::/64", "::", false)
	t.testNotContains("1:2::/64", "::")
	t.testContains("1:2::/64", "1:2::", false)

	t.testNotContains("5.62.62-63.*", "5.62.64.1")
	t.testNotContains("5.62.62-63.*", "5.62.68.1")
	t.testNotContains("5.62.62-63.*", "5.62.78.1")

	t.testContains("192.13.1.0/25", "192.13.1.1-127", false)

	t.testNotContains("192.13.1.0/25", "192.13.1.1-255")
	//testContainsNonZeroHosts("192.13.1.1-127", "192.13.1.0/25")
	//testContainsNonZeroHosts("192.13.1.1-255", "192.13.1.0/24")
	//testNotContainsNonZeroHosts("192.13.1.1-255", "192.13.1.0/23")
	//
	//testContainsNonZeroHosts("192.13.1.0-255", "192.13.1.0/23")

	t.testContains("192.13.1.0-255", "192.13.1.0/23", false)

	t.testContains("192.13.0-1.0-255", "192.13.1.0/23", false)
	t.testContains("192.13.0-1.0-255", "192.13.0.0/23", true)

	//testContainsNonZeroHosts("::192:13:1:1-7fff", "::192:13:1:0/113")
	//testContainsNonZeroHosts("::192:13:1:1-ffff", "::192:13:1:0/112")
	//testNotContainsNonZeroHosts("::192:13:1:1-ffff", "::192:13:1:0/111")

	t.testSubnet("1.2-4.3.4", "255.255.254.255", 24, "1.2-4.2.4/24", "1.2-4.2.4", "1.2-4.3.4/24")
	t.testSubnet("1.2-4.3.4", "255.248.254.255", 24, "1.0.2.4/24", "1.0.2.4", "1.2-4.3.4/24")

	t.testSubnet("__::", "ffff::", 128, "0-ff:0:0:0:0:0:0:0/128", "0-ff:0:0:0:0:0:0:0", "0-ff:0:0:0:0:0:0:0/128")
	t.testSubnet("0-ff::", "fff0::", 128, "", "", "0-ff:0:0:0:0:0:0:0/128")

	t.testSubnet("0-ff::", "fff0::", 12, "0-ff:0:0:0:0:0:0:0/12", "", "0-ff:0:0:0:0:0:0:0/12")
	//testSubnet("0-f0::", "fff0::", 12, "0-f0:0:0:0:0:0:0:0/12", "0-f0:0:0:0:0:0:0:0", "0-f0:0:0:0:0:0:0:0/12");
	t.testSubnet("0-f0::", "fff0::", 12, "0-f0:0:0:0:0:0:0:0/12", "", "0-f0:0:0:0:0:0:0:0/12")
	t.testSubnet("0-f::", "fff0::", 12, "0-f:0:0:0:0:0:0:0/12", "0:0:0:0:0:0:0:0", "0-f:0:0:0:0:0:0:0/12")
	t.testSubnet("0-f::*", "fff0::ffff", 12, "0-f:0:0:0:0:0:0:*/12", "0:0:0:0:0:0:0:*", "0-f:0:0:0:0:0:0:*/12")

	t.testSubnet("::1:__", "::1:ffff", 128, "0:0:0:0:0:0:1:0-ff/128", "0:0:0:0:0:0:1:0-ff", "0:0:0:0:0:0:1:0-ff/128")
	t.testSubnet("::1:__", "::1:ffff", 126, "0:0:0:0:0:0:1:0-fc/126", "0:0:0:0:0:0:1:0-ff", "0:0:0:0:0:0:1:0-fc/126")
	t.testSubnet("::1:0-ff", "::1:fff0", 128, "", "", "0:0:0:0:0:0:1:0-ff/128")
	t.testSubnet("::1:0-ff", "::1:fff0", 124, "0:0:0:0:0:0:1:0-f0/124", "", "0:0:0:0:0:0:1:0-f0/124")
	t.testSubnet("*::1:0-f", "ffff::1:fff0", 124, "*:0:0:0:0:0:1:0/124", "*:0:0:0:0:0:1:0", "*:0:0:0:0:0:1:0/124")

	t.testReverseHostAddress("*.*.0-240.0/20")
	t.testReverseHostAddress("*.*.0.0/16")
	t.testReverseHostAddress("*:0-f000::/20")

	t.testResolved("8.*.27.26", "8.*.27.26")

	t.testResolved("2001:*:0:0:8:800:200C:417A", "2001:*:0:0:8:800:200C:417A")

	t.testNormalized("ABCD:EF12:*:*:***:A:*:BBBB", "abcd:ef12:*:*:*:a:*:bbbb")
	t.testNormalized("ABCD:EF12:*:*:**:A:***:BBBB%g", "abcd:ef12:*:*:*:a:*:bbbb%g")

	t.testNormalized("1.*", "1.*.*.*")
	t.testNormalized("*.1.*", "*.1.*.*")
	t.testNormalized("*:1::*", "*:1::*")
	t.testNormalized("*:1:*", "*:1:*:*:*:*:*:*")
	t.testNormalized("001-002:0001-0002:01-2:1-02:01-02:*", "1-2:1-2:1-2:1-2:1-2:*:*:*")

	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/0", "1:2:3:4:5:6:7:8/0", []ipaddr.BitCount{0, 0, 0, 0, 0, 0, 0, 0, 0})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb", "1:2:3:4:5:6:7:8/0", []ipaddr.BitCount{0, 16, 32, 48, 64, 80, 96, 112, 128})
	t.testInsertAndAppendPrefs("a:b:c:d:e:f:aa:bb/0", "1:2:3:4:5:6:7:8", []ipaddr.PrefixLen{nil, p0, p0, p0, p0, p0, p0, p0, p0})

	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/64", []ipaddr.BitCount{64, 64, 64, 64, 64, 64, 64, 64, 64})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb", "1:2:3:4:5:6:7:8/64", []ipaddr.BitCount{64, 64, 64, 64, 64, 80, 96, 112, 128})
	t.testInsertAndAppendPrefs("a:b:c:d:e:f:aa:bb/63", "1:2:3:4:5:6:7:8", []ipaddr.PrefixLen{nil, nil, nil, nil, p63, p63, p63, p63, p63})
	t.testInsertAndAppendPrefs("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8", []ipaddr.PrefixLen{nil, nil, nil, nil, p64, p64, p64, p64, p64})
	t.testInsertAndAppendPrefs("a:b:c:d:e:f:aa:bb/65", "1:2:3:4:5:6:7:8", []ipaddr.PrefixLen{nil, nil, nil, nil, nil, p65, p65, p65, p65})

	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/128", "1:2:3:4:5:6:7:8/128", []ipaddr.BitCount{128, 128, 128, 128, 128, 128, 128, 128, 128})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb", "1:2:3:4:5:6:7:8/128", []ipaddr.BitCount{128, 128, 128, 128, 128, 128, 128, 128, 128})
	t.testInsertAndAppendPrefs("a:b:c:d:e:f:aa:bb/128", "1:2:3:4:5:6:7:8", []ipaddr.PrefixLen{nil, nil, nil, nil, nil, nil, nil, nil, p128})

	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/32", "1:2:3:4:5:6:7:8/64", []ipaddr.BitCount{64, 64, 32, 32, 32, 32, 32, 32, 32})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/32", []ipaddr.BitCount{32, 32, 32, 48, 64, 64, 64, 64, 64})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/0", "1:2:3:4:5:6:7:8/64", []ipaddr.BitCount{64, 0, 0, 0, 0, 0, 0, 0, 0})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/0", []ipaddr.BitCount{0, 16, 32, 48, 64, 64, 64, 64, 64})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/128", []ipaddr.BitCount{128, 128, 128, 128, 64, 64, 64, 64, 64})
	t.testInsertAndAppend("a:b:c:d:e:f:aa:bb/128", "1:2:3:4:5:6:7:8/64", []ipaddr.BitCount{64, 64, 64, 64, 64, 80, 96, 112, 128})

	t.testInsertAndAppend("1.2.3.4/0", "5.6.7.8/0", []ipaddr.BitCount{0, 0, 0, 0, 0})
	t.testInsertAndAppend("1.2.3.4", "5.6.7.8/0", []ipaddr.BitCount{0, 8, 16, 24, 32})
	t.testInsertAndAppendPrefs("1.2.3.4/0", "5.6.7.8", []ipaddr.PrefixLen{nil, p0, p0, p0, p0})

	t.testInsertAndAppend("1.2.3.4/16", "5.6.7.8/16", []ipaddr.BitCount{16, 16, 16, 16, 16})
	t.testInsertAndAppend("1.2.3.4", "5.6.7.8/16", []ipaddr.BitCount{16, 16, 16, 24, 32})
	t.testInsertAndAppendPrefs("1.2.3.4/16", "5.6.7.8", []ipaddr.PrefixLen{nil, nil, p16, p16, p16})

	t.testInsertAndAppend("1.2.3.4/32", "5.6.7.8/32", []ipaddr.BitCount{32, 32, 32, 32, 32})
	t.testInsertAndAppend("1.2.3.4", "5.6.7.8/32", []ipaddr.BitCount{32, 32, 32, 32, 32})
	t.testInsertAndAppendPrefs("1.2.3.4/31", "5.6.7.8", []ipaddr.PrefixLen{nil, nil, nil, nil, p31})
	t.testInsertAndAppendPrefs("1.2.3.4/32", "5.6.7.8", []ipaddr.PrefixLen{nil, nil, nil, nil, p32})

	t.testInsertAndAppend("1.2.3.4/16", "5.6.7.8/24", []ipaddr.BitCount{24, 24, 16, 16, 16})
	t.testInsertAndAppend("1.2.3.4/24", "5.6.7.8/7", []ipaddr.BitCount{7, 8, 16, 24, 24})
	t.testInsertAndAppend("1.2.3.4/24", "5.6.7.8/16", []ipaddr.BitCount{16, 16, 16, 24, 24})
	t.testInsertAndAppend("1.2.3.4/0", "5.6.7.8/16", []ipaddr.BitCount{16, 0, 0, 0, 0})
	t.testInsertAndAppend("1.2.3.4/16", "5.6.7.8/0", []ipaddr.BitCount{0, 8, 16, 16, 16})
	t.testInsertAndAppend("1.2.3.4/17", "5.6.7.8/0", []ipaddr.BitCount{0, 8, 16, 17, 17})
	t.testInsertAndAppend("1.2.3.4/16", "5.6.7.8/32", []ipaddr.BitCount{32, 32, 16, 16, 16})
	t.testInsertAndAppend("1.2.3.4/32", "5.6.7.8/16", []ipaddr.BitCount{16, 16, 16, 24, 32})

	t.testReplace("a:b:c:d:e:f:aa:bb/0", "1:2:3:4:5:6:7:8/0")
	t.testReplace("a:b:c:d:e:f:aa:bb", "1:2:3:4:5:6:7:8/0")
	t.testReplace("a:b:c:d:e:f:aa:bb/0", "1:2:3:4:5:6:7:8")

	t.testReplace("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/64")
	t.testReplace("a:b:c:d:e:f:aa:bb", "1:2:3:4:5:6:7:8/64")
	t.testReplace("a:b:c:d:e:f:aa:bb/63", "1:2:3:4:5:6:7:8")
	t.testReplace("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8")
	t.testReplace("a:b:c:d:e:f:aa:bb/65", "1:2:3:4:5:6:7:8")

	t.testReplace("a:b:c:d:e:f:aa:bb/128", "1:2:3:4:5:6:7:8/128")
	t.testReplace("a:b:c:d:e:f:aa:bb", "1:2:3:4:5:6:7:8/128")
	t.testReplace("a:b:c:d:e:f:aa:bb/128", "1:2:3:4:5:6:7:8")

	t.testReplace("a:b:c:d:e:f:aa:bb/32", "1:2:3:4:5:6:7:8/64")
	t.testReplace("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/32")
	t.testReplace("a:b:c:d:e:f:aa:bb/0", "1:2:3:4:5:6:7:8/64")
	t.testReplace("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/0")
	t.testReplace("a:b:c:d:e:f:aa:bb/64", "1:2:3:4:5:6:7:8/128")
	t.testReplace("a:b:c:d:e:f:aa:bb/128", "1:2:3:4:5:6:7:8/64")

	t.testReplace("1.2.3.4/0", "5.6.7.8/0")
	t.testReplace("1.2.3.4", "5.6.7.8/0")
	t.testReplace("1.2.3.4/0", "5.6.7.8")

	t.testReplace("1.2.3.4/16", "5.6.7.8/16")
	t.testReplace("1.2.3.4", "5.6.7.8/16")
	t.testReplace("1.2.3.4/16", "5.6.7.8")

	t.testReplace("1.2.3.4/32", "5.6.7.8/32")
	t.testReplace("1.2.3.4", "5.6.7.8/32")
	t.testReplace("1.2.3.4/31", "5.6.7.8")
	t.testReplace("1.2.3.4/32", "5.6.7.8")

	t.testReplace("1.2.3.4/16", "5.6.7.8/24")
	t.testReplace("1.2.3.4/24", "5.6.7.8/7")
	t.testReplace("1.2.3.4/24", "5.6.7.8/16")
	t.testReplace("1.2.3.4/0", "5.6.7.8/16")
	t.testReplace("1.2.3.4/16", "5.6.7.8/0")
	t.testReplace("1.2.3.4/17", "5.6.7.8/0")
	t.testReplace("1.2.3.4/16", "5.6.7.8/32")
	t.testReplace("1.2.3.4/32", "5.6.7.8/16")

	t.testSub("1:1::/32", "1:1:1:1:1:1:1:1", []string{
		"1:1:0:0:0:0:0:0/48",
		"1:1:2-fffe:0:0:0:0:0/47",
		"1:1:1:0:0:0:0:0/64",
		"1:1:1:2-fffe:0:0:0:0/63",
		"1:1:1:1:0:0:0:0/80",
		"1:1:1:1:2-fffe:0:0:0/79",
		"1:1:1:1:1:0:0:0/96",
		"1:1:1:1:1:2-fffe:0:0/95",
		"1:1:1:1:1:1:0:0/112",
		"1:1:1:1:1:1:2-fffe:0/111",
		"1:1:1:1:1:1:1:0",
		"1:1:1:1:1:1:1:2-fffe/127",
	})
	t.testSub("1:1::/32", "1:1::/16", []string{
		"1:1:1-ffff:0:0:0:0:0/48",
		"1:1:0:1-ffff:0:0:0:0/64",
		"1:1:0:0:1-ffff:0:0:0/80",
		"1:1:0:0:0:1-ffff:0:0/96",
		"1:1:0:0:0:0:1-ffff:0/112",
		"1:1:0:0:0:0:0:1-ffff"},
	)
	t.testSub("1:1::/32", "1:1::/48", []string{"1:1:1-ffff:0:0:0:0:0/48"})
	t.testSub("1:1::/32", "1:1::/64", []string{
		"1:1:1-ffff:0:0:0:0:0/48",
		"1:1:0:1-ffff:0:0:0:0/64",
	})
	t.testSub("1:1::/32", "1:1:2:2::/64", []string{
		"1:1:0:0:0:0:0:0/47",
		"1:1:3-ffff:0:0:0:0:0/48",
		"1:1:2:0:0:0:0:0/63",
		"1:1:2:3-ffff:0:0:0:0/64",
	})
	t.testSub("10.0.0.0/22", "10.0.0.0/24", []string{"10.0.1-3.0/24"}) //[10.0.1-3.0/24]

	t.testIntersect("1:1:1-3:1:1:1:1:1", "1:1:2-4:1:1:1:1:1", "1:1:2-3:1:1:1:1:1")
	t.testIntersect("1:1:1-3:1:0:1:1:1", "1:1:2-4:1:1:1:1:1", "")

	t.testToPrefixBlock("1.3.*.*", "1.3.*.*")
	t.testToPrefixBlock("1.2-3.*.*", "1.2-3.*.*")
	t.testToPrefixBlock("1.3.3.4/15", "1.2-3.*.*/15")
	t.testToPrefixBlock("*.3.3.4/15", "*.2-3.*.*/15")
	t.testToPrefixBlock("1.3.3.4/16", "1.3.*.*/16")

	t.testToPrefixBlock("1:3:3:4::/15", "0-1:*/15")
	t.testToPrefixBlock("*:3:3:4::/15", "0-fffe::/15")
	t.testToPrefixBlock("1:3:3:4::/16", "1:*/16")

	t.testMaxHost("1.*.255.255/16", "1.*.255.255/16")
	t.testMaxHost("1.2.*.*/16", "1.2.255.255/16")
	t.testMaxHost("1.*.*.*/16", "1.*.255.255/16")
	t.testMaxHost("1.2.*.1/16", "1.2.255.255/16")
	t.testMaxHost("1.*.*.1/16", "1.*.255.255/16")

	t.testZeroHost("1.*.0.0/16", "1.*.0.0/16")
	t.testZeroHost("1.2.*.*/16", "1.2.0.0/16")
	t.testZeroHost("1.*.*.*/16", "1.*.0.0/16")
	t.testZeroHost("1.2.*.1/16", "1.2.0.0/16")
	t.testZeroHost("1.*.*.1/16", "1.*.0.0/16")

	t.testZeroNetwork("1.*.0.0/16", "0.0.0.0/16")
	t.testZeroNetwork("1.2.*.*/16", "0.0.*.*/16")
	t.testZeroNetwork("1.*.*.*/16", "0.0.*.*/16")
	t.testZeroNetwork("1.2.*.1/16", "0.0.*.1/16")
	t.testZeroNetwork("1.*.*.1/16", "0.0.*.1/16")

	t.testMaxHost("1:*::ffff:ffff:ffff:ffff/64", "1:*::ffff:ffff:ffff:ffff/64")
	t.testMaxHost("1:2::ffff:ffff:ffff:ffff/64", "1:2::ffff:ffff:ffff:ffff/64")
	t.testMaxHost("1:*::*:*:*:*/64", "1:*::ffff:ffff:ffff:ffff/64")
	t.testMaxHost("1:2::*:*:*:*/64", "1:2::ffff:ffff:ffff:ffff/64")
	t.testMaxHost("1:2::*:*:*:1/64", "1:2::ffff:ffff:ffff:ffff/64")
	t.testMaxHost("1:*:1/64", "1:*:ffff:ffff:ffff:ffff/64")
	t.testMaxHost("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/0", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/0")
	t.testMaxHost("*:*/0", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/0")
	t.testMaxHost("::/128", "::/128")

	t.testZeroHost("1:*::/64", "1:*::/64")
	t.testZeroHost("1:2::/64", "1:2::/64")
	t.testZeroHost("1:*::*:*:*:*/64", "1:*::/64")
	t.testZeroHost("1:2::*:*:*:*/64", "1:2::/64")
	t.testZeroHost("1:2::*:*:*:1/64", "1:2::/64")
	t.testZeroHost("1:*:1/64", "1:*:*:*::/64")
	t.testZeroHost("::/0", "::/0")
	t.testZeroHost("*:*/0", "::/0")
	t.testZeroHost("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128")

	t.testZeroHost("1:2:3:4::/64", "1:2:3:4::/64")
	t.testZeroHost("1:2:*/64", "1:2:*:*::/64")
	t.testZeroHost("1:2:3:4:*:1/64", "1:2:3:4::/64")
	t.testZeroHost("1:*:1/64", "1:*:*:*::/64")

	t.testZeroNetwork("1:*::/64", "::/64")
	t.testZeroNetwork("1:2::/64", "::/64")
	t.testZeroNetwork("1:*::*:*:*:*/64", "::*:*:*:*/64")
	t.testZeroNetwork("1:2::*:*:*:*/64", "::*:*:*:*/64")
	t.testZeroNetwork("1:2::*:*:*:1/64", "::*:*:*:1/64")
	t.testZeroNetwork("1:*:1/64", "::*:*:*:1/64")
	t.testZeroNetwork("::/0", "::/0")
	t.testZeroNetwork("*:*/0", "*:*/0")
	t.testZeroNetwork("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128", "::/128")

	t.testZeroNetwork("1:2:3:4::/64", "::/64")
	t.testZeroNetwork("1:2:3:4:*/64", "0:0:0:0:*/64")
	t.testZeroNetwork("1:2:*/64", "0:0:0:0:*/64")
	t.testZeroNetwork("1:2:3:4:*:1/64", "0:0:0:0:*:1/64")
	t.testZeroNetwork("1:*:1/64", "0:0:0:0:*:1/64")

	t.testIsPrefixBlock("1.2.*.*", false, false)
	t.testIsPrefixBlock("1.2.3.*", false, false)
	t.testIsPrefixBlock("1.*.*.*", false, false)
	t.testIsPrefixBlock("1.2-3.*.*", false, false)
	t.testIsPrefixBlock("1.2.128-255.*", false, false)
	t.testIsPrefixBlock("*.*/0", true, true)
	t.testIsPrefixBlock("1.2.*.*/16", true, true)
	t.testIsPrefixBlock("1.2.3.*/16", false, false)
	t.testIsPrefixBlock("1.*.*.*/16", true, false)
	t.testIsPrefixBlock("1.2-3.*.*/16", true, false)
	t.testIsPrefixBlock("1.2.128-255.*/16", false, false)

	t.testPrefixBlocks("1.2.*.*", 8, false, false)
	t.testPrefixBlocks("1.2.3.*", 8, false, false)
	t.testPrefixBlocks("1.*.*.*", 8, true, true)
	t.testPrefixBlocks("1.2-3.*.*", 8, false, false)
	t.testPrefixBlocks("1.2.128-255.*", 8, false, false)
	t.testPrefixBlocks("*.*/0", 8, true, false)
	t.testPrefixBlocks("1.2.*.*/16", 8, false, false)
	t.testPrefixBlocks("1.2.3.*/16", 8, false, false)
	t.testPrefixBlocks("1.*.*.*/16", 8, true, true)
	t.testPrefixBlocks("1.2-3.*.*/16", 8, false, false)
	t.testPrefixBlocks("1.2.128-255.*/16", 8, false, false)

	t.testPrefixBlocks("1.2.*.*", 24, true, false)
	t.testPrefixBlocks("1.2.3.*", 24, true, true)
	t.testPrefixBlocks("1.*.*.*", 24, true, false)
	t.testPrefixBlocks("1.2-3.*.*", 24, true, false)
	t.testPrefixBlocks("1.2.128-255.*", 24, true, false)
	t.testPrefixBlocks("*.*/0", 24, true, false)
	t.testPrefixBlocks("1.2.*.*/16", 24, true, false)
	t.testPrefixBlocks("1.2.3.*/16", 24, true, !false)
	t.testPrefixBlocks("1.*.*.*/16", 24, true, false)
	t.testPrefixBlocks("1.2-3.*.*/16", 24, true, false)
	t.testPrefixBlocks("1.2.128-255.*/16", 24, true, false)

	t.testIsPrefixBlock("a:b:c:d:*/64", true, true)
	t.testIsPrefixBlock("a:b:c:*/64", true, false)
	t.testIsPrefixBlock("a:b:c:d-e:*/64", true, false)
	t.testIsPrefixBlock("a:b:c:d:e:*/64", false, false)
	t.testIsPrefixBlock("a:b:c:d:0-ffff:*/64", true, true)
	t.testIsPrefixBlock("a:b:c:d:8000-ffff:*/64", false, false)

	t.testPrefixBlocks("a:b:c:d:*/64", 0, false, false)
	t.testPrefixBlocks("a:b:c:*/64", 0, false, false)
	t.testPrefixBlocks("a:b:c:d-e:*/64", 0, false, false)
	t.testPrefixBlocks("*:*/64", 0, true, true)
	t.testPrefixBlocks("a:b:c:d:e:*/64", 0, false, false)
	t.testPrefixBlocks("a:b:c:d:0-ffff:*/64", 0, false, false)

	t.testPrefixBlocks("a:b:c:d:*/64", 63, false, false)
	t.testPrefixBlocks("a:b:c:*/64", 63, true, false)
	t.testPrefixBlocks("a:b:c:d-e:*/64", 63, false, false)
	t.testPrefixBlocks("a:b:c:e-f:*/64", 63, true, true)
	t.testPrefixBlocks("a:b:c:d:e:*/64", 63, false, false)
	t.testPrefixBlocks("a:b:c:d:0-ffff:*/64", 63, false, false)

	t.testPrefixBlocks("a:b:c:d:*/64", 64, true, true)
	t.testPrefixBlocks("a:b:c:*/64", 64, true, false)
	t.testPrefixBlocks("a:b:c:d-e:*/64", 64, true, false)
	t.testPrefixBlocks("a:b:c:d:e:*/64", 64, false, false)
	t.testPrefixBlocks("a:b:c:d:0-ffff:*/64", 64, true, true)
	t.testPrefixBlocks("a:b:c:d:8000-ffff:*/64", 64, false, false)

	t.testPrefixBlocks("a:b:c:d:*/64", 65, true, false)
	t.testPrefixBlocks("a:b:c:*/64", 65, true, false)
	t.testPrefixBlocks("a:b:c:d-e:*/64", 65, true, false)
	t.testPrefixBlocks("a:b:c:d:e:*/64", 65, false, false)
	t.testPrefixBlocks("a:b:c:d:0-ffff:*/64", 65, true, !true)
	t.testPrefixBlocks("a:b:c:d:8000-ffff:*/64", 65, true, !false)
	t.testPrefixBlocks("a:b:c:d:0-ffff:*/64", 65, true, false)

	t.testPrefixBlocks("a:b:c:d:*/64", 128, true, false)
	t.testPrefixBlocks("a:b:c:*/64", 128, true, false)
	t.testPrefixBlocks("a:b:c:d-e:*/64", 128, true, false)
	t.testPrefixBlocks("a:b:c:d:e:*/64", 128, true, false)
	t.testPrefixBlocks("a:b:c:d:0-ffff:*/64", 128, true, false)

	t.testSplitBytes("1.2.*.4")
	t.testSplitBytes("1.2-4.3.4/16")
	t.testSplitBytes("1.2.3.4-5/0")
	t.testSplitBytes("1.2.*/32")
	t.testSplitBytes("ffff:2:3:4:eeee:dddd:cccc-dddd:bbbb")
	t.testSplitBytes("ffff:2:3:4:eeee:dddd:cccc:bbbb/64")
	t.testSplitBytes("ffff:2:3:4:*:dddd:cccc:bbbb/0")
	t.testSplitBytes("*:*/128")
	t.testSplitBytes("*:*")

	t.testIncrement("1.2.*.*/16", 0, "1.2.0.0")
	t.testIncrement("1.2.*.*/16", 1, "1.2.0.1")
	t.testIncrement("1.2.*.*/16", 65535, "1.2.255.255")
	t.testIncrement("1.2.*.*/16", 65536, "1.3.0.0")
	t.testIncrement("1.2.*.*/16", -1, "1.1.255.255")
	t.testIncrement("1.2.*.*/16", -65536, "1.1.0.0")
	t.testIncrement("1.2.*.*/16", -65537, "1.0.255.255")

	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", 0, "ffff:ffff:ffff:ffff:ffff:1:2:ffff")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", 1, "ffff:ffff:ffff:ffff:ffff:1:3:ffff")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", 3, "ffff:ffff:ffff:ffff:ffff:2:3:ffff")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", 4, "ffff:ffff:ffff:ffff:ffff:2:4::")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", 5, "ffff:ffff:ffff:ffff:ffff:2:4:1")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:fffe-ffff:fffe-ffff:ffff", 5, "")

	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", -0x10002ffff, "ffff:ffff:ffff:ffff:ffff::")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", -0x100030000, "ffff:ffff:ffff:ffff:fffe:ffff:ffff:ffff")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", -0x100030003, "ffff:ffff:ffff:ffff:fffe:ffff:ffff:fffc")
	t.testIncrement("ffff:ffff:ffff:ffff:ffff:1-2:2-3:ffff", -0x100030004, "ffff:ffff:ffff:ffff:fffe:ffff:ffff:fffb")

	t.testIncrement("::1-2:2-3:ffff", -0x100030000, "")

	t.testIncrement("ffff:3-4:ffff:ffff:ffff:1-2:2-3::", 7, "ffff:4:ffff:ffff:ffff:2:3::")
	t.testIncrement("ffff:3-4:ffff:ffff:ffff:1-2:2-3::", 9, "ffff:4:ffff:ffff:ffff:2:3:2")

	t.testLeadingZeroAddr("00-1.1.2.3", true)
	t.testLeadingZeroAddr("1.00-1.2.3", true)
	t.testLeadingZeroAddr("1.2.00-1.3", true)
	t.testLeadingZeroAddr("1.2.3.00-1", true)
	t.testLeadingZeroAddr("1-01.1.2.3", true)
	t.testLeadingZeroAddr("1.01-1.2.3", true)
	t.testLeadingZeroAddr("1.2.1-01.3", true)
	t.testLeadingZeroAddr("1.2.3.01-1", true)
	t.testLeadingZeroAddr("0-1.1.2.3", false)
	t.testLeadingZeroAddr("1.0-1.2.3", false)
	t.testLeadingZeroAddr("1.2.0-1.3", false)
	t.testLeadingZeroAddr("1.2.3.0-1", false)

	t.testLeadingZeroAddr("00-1:1:2:3::", true)
	t.testLeadingZeroAddr("1:00-1:2:3::", true)
	t.testLeadingZeroAddr("1:2:00-1:3::", true)
	t.testLeadingZeroAddr("1:2:3:00-1::", true)
	t.testLeadingZeroAddr("1-01:1:2:3::", true)
	t.testLeadingZeroAddr("1:1-01:2:3::", true)
	t.testLeadingZeroAddr("1:2:1-01:3::", true)
	t.testLeadingZeroAddr("1:2:3:1-01::", true)
	t.testLeadingZeroAddr("0-1:1:2:3::", false)
	t.testLeadingZeroAddr("1:0-1:2:3::", false)
	t.testLeadingZeroAddr("1:2:0-1:3::", false)
	t.testLeadingZeroAddr("1:2:3:0-1::", false)

	t.testRangeExtend("1.2.3.4-5", "1.2.4.3", "1.2.3-5.6", "", "1.2.3.4", "1.2.5.6")
	t.testRangeExtend("1.2.3.4-5", "1.2.4.3", "1.2.1-5.6", "", "1.2.1.6", "1.2.5.6")

	t.testIncompatibleAddress2("a:b:c:d:e:f:1.2.*.4", "a:b:c:d:e:f:1.2.0.4", "a:b:c:d:e:f:1.2.255.4", []interface{}{0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 1, 2, []uint{0, 0xff}, 4}) //[a, b, c, d, e, f, 1, 2, 0-ff, 4]
	t.testIncompatibleAddress2("::ffff:0.0.*.0", "::ffff:0.0.0.0", "::ffff:0.0.255.0", []interface{}{0, 0xffff, 0, 0, []uint{0, 0xff}, 0})                                   //[0, ffff, 0, 0, 0-ff, 0]
	t.testIncompatibleAddress2("::ffff:*.0.0.0", "::ffff:0.0.0.0", "::ffff:255.0.0.0", []interface{}{0, 0xffff, []uint{0, 0xff}, 0, 0, 0})                                   //[0, ffff, 0-ff, 0, 0, 0]
	t.testMaskedIncompatibleAddress("0-ffff::1/f000::10", "::", "f000::")
	t.testSubnetStringRange("0-ffff::1/f000::", "::1", "ffff::1", []interface{}{[]uint{0, 0xffff}, 0, 1}, p4)
	t.testSubnetStringRange("0-ffff::/f000::", "::", "ffff::", []interface{}{[]uint{0, 0xffff}, 0}, p4)
	t.testSubnetStringRange("0-f000::/f000::", "::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", []interface{}{[]uint{0, 0xffff}, []*big.Int{bigZeroConst(), setBigString("ffffffffffffffffffffffffffff", 16)}}, p4) //[0-f000, 0]

	t.testSubnetStringRange2("0-ffff::/0fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "::", "fff::", []interface{}{[]uint{0, 0xfff}, 0}) // [0-fff, 0]  // /8 prefix?

	t.testSubnetStringRange2("1.*.*.*", "1.0.0.0", "1.255.255.255", []interface{}{1, []uint{0, 0xff}, []uint{0, 0xff}, []uint{0, 0xff}})                                                       //[1, 0-255, 0-255, 0-255]
	t.testSubnetStringRange2("1.*.*", "1.0.0.0", "1.255.255.255", []interface{}{1, []uint{0, 0xff}, []uint{0, 0xffff}})                                                                        //[1, 0-255, 0-65535]
	t.testSubnetStringRange2("1.*", "1.0.0.0", "1.255.255.255", []interface{}{1, []uint{0, 0xffffff}})                                                                                         //[1, 0-16777215]
	t.testSubnetStringRange2("a:b:c:*:cc:d:e:f", "a:b:c:0:cc:d:e:f", "a:b:c:ffff:cc:d:e:f", []interface{}{0xa, 0xb, 0xc, []uint{0, 0xffff}, 0xcc, 0xd, 0xe, 0xf})                              //[a, b, c, 0-ffff, cc, d, e, f]
	t.testSubnetStringRange2("a:*:cc:d:e:f", "a::cc:d:e:f", "a:ffff:ffff:ffff:cc:d:e:f", []interface{}{0xa, []uint64{0, 0xffffffffffff}, 0xcc, 0xd, 0xe, 0xf})                                 //[a, 0-ffffffffffff, cc, d, e, f]
	t.testSubnetStringRange2("*:cc:d:e:f", "::cc:d:e:f", "ffff:ffff:ffff:ffff:cc:d:e:f", []interface{}{[]*big.Int{bigZeroConst(), setBigString("ffffffffffffffff", 16)}, 0xcc, 0xd, 0xe, 0xf}) //[0-ffffffffffffffff, cc, d, e, f]

	t.testSubnetStringRange2("a:b:c:*:cc:d:1.255.3.128", "a:b:c:0:cc:d:1.255.3.128", "a:b:c:ffff:cc:d:1.255.3.128", []interface{}{0xa, 0xb, 0xc, []uint{0, 0xffff}, 0xcc, 0xd, 1, 255, 3, 128})                              //[a, b, c, 0-ffff, cc, d, e, f]
	t.testSubnetStringRange2("a:*:cc:d:1.255.3.128", "a::cc:d:1.255.3.128", "a:ffff:ffff:ffff:cc:d:1.255.3.128", []interface{}{0xa, []uint64{0, 0xffffffffffff}, 0xcc, 0xd, 1, 255, 3, 128})                                 //[a, 0-ffffffffffff, cc, d, e, f]
	t.testSubnetStringRange2("*:cc:d:1.255.3.128", "::cc:d:1.255.3.128", "ffff:ffff:ffff:ffff:cc:d:1.255.3.128", []interface{}{[]*big.Int{bigZeroConst(), setBigString("ffffffffffffffff", 16)}, 0xcc, 0xd, 1, 255, 3, 128}) //[0-ffffffffffffffff, cc, d, e, f]

	if t.isLenient() {
		// inet_aton
		t.testSubnetStringRange2("1.*.1", "1.0.0.1", "1.255.0.1", []interface{}{1, []uint{0, 0xff}, 1})                                                                                                       //[1, 0-255, 1]
		t.testSubnetStringRange2("*.1", "0.0.0.1", "255.0.0.1", []interface{}{[]uint{0, 0xff}, 1})                                                                                                            //[0-255, 1]
		t.testIncompatibleAddress2("a:b:cc:*.4", "a:b:cc:0:0:0:0.0.0.4", "a:b:cc:ffff:ffff:ffff:255.0.0.4", []interface{}{0xa, 0xb, 0xcc, []*big.Int{bigZeroConst(), setBigString("ffffffffffffff", 16)}, 4}) //[a, b, cc, 0-ffffffffffffff, 4]
		t.testIncompatibleAddress2("1:2:3:4:*.3.4", "1:2:3:4::0.3.0.4", "1:2:3:4:ffff:ffff:255.3.0.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffff}, 3, 4})                                            //[1, 2, 3, 4, 0-ffffffffff, 3, 4]
		t.testIncompatibleAddress2("1:2:3:4:*.4", "1:2:3:4::0.0.0.4", "1:2:3:4:ffff:ffff:255.0.0.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffff}, 4})                                                 //[1, 2, 3, 4, 0-ffffffffff, 4]
	} else {
		// not inet_aton
		t.testSubnetStringRange2("1.*.1", "1.0.0.1", "1.255.255.1", []interface{}{1, []uint{0, 0xffff}, 1})
		t.testSubnetStringRange2("*.1", "0.0.0.1", "255.255.255.1", []interface{}{[]uint{0, 0xffffff}, 1})
		t.testIncompatibleAddress2("a:b:cc:*.4", "a:b:cc:0:0:0:0.0.0.4", "a:b:cc:ffff:ffff:ffff:255.255.255.4", []interface{}{0xa, 0xb, 0xcc, []*big.Int{bigZeroConst(), setBigString("ffffffffffffffffff", 16)}, 4}) //[a, b, cc, 0-ffffffffffffffffff, 4]
		t.testSubnetStringRange2("1:2:3:4:*.3.4", "1:2:3:4::0.0.3.4", "1:2:3:4:ffff:ffff:255.255.3.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffffff}, 3, 4})                                                  //[1, 2, 3, 4, 0-ffffffffffff, 3, 4]
		t.testIncompatibleAddress2("1:2:3:4:*.4", "1:2:3:4::0.0.0.4", "1:2:3:4:ffff:ffff:255.255.255.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffffffff}, 4})                                                 //[1, 2, 3, 4, 0-ffffffffffffff, 4]
	}
	t.testSubnetStringRange1("1-2.3.4-5.6", "1.3.4.6", "2.3.5.6", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, 6}, nil, false)                                                         //[1-2, 3, 4-5, 6]
	t.testSubnetStringRange1("1-2:3:4-5:6::", "1:3:4:6::", "2:3:5:6::", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, 6, 0}, nil, false)                                                //[1-2, 3, 4-5, 6, 0]
	t.testIncompatibleAddress1("1:2:3:4:5:6:1-3.2.0.4-5", "1:2:3:4:5:6:1.2.0.4", "1:2:3:4:5:6:3.2.0.5", []interface{}{1, 2, 3, 4, 5, 6, []uint{1, 3}, 2, 0, []uint{4, 5}}, nil, false) //[1, 2, 3, 4, 5, 6, 1-3, 2, 0, 4-5]
	t.testMaskedIncompatibleAddress("0.0.0.*/0.0.0.128", "0.0.0.0", "0.0.0.128")                                                                                                       //iae

	t.testSubnetStringRange1("1.2-3.4.5", "1.2.4.5", "1.3.4.5", []interface{}{1, []uint{2, 3}, 4, 5}, nil, false)                                                                                  //[1, 2-3, 4, 5]
	t.testSubnetStringRange1("1:2-3:4:5::", "1:2:4:5::", "1:3:4:5::", []interface{}{1, []uint{2, 3}, 4, 5, 0}, nil, false)                                                                         //[1, 2-3, 4, 5, 0]
	t.testSubnetStringRange1("1:2:4:5:6-9:7:8:f", "1:2:4:5:6:7:8:f", "1:2:4:5:9:7:8:f", []interface{}{1, 2, 4, 5, []uint{6, 9}, 7, 8, 0xf}, nil, false)                                            //[1, 2, 4, 5, 6-9, 7, 8, f]
	t.testIncompatibleAddress1("a:b:cc:dd:e:*.2.3.4", "a:b:cc:dd:e:0:0.2.3.4", "a:b:cc:dd:e:ffff:255.2.3.4", []interface{}{0xa, 0xb, 0xcc, 0xdd, 0xe, []uint{0, 0xffffff}, 2, 3, 4}, nil, false)   // [a, b, cc, dd, e, 0-ffffff, 2, 3, 4]
	t.testIncompatibleAddress1("a:b:cc:dd:*.2.3.4", "a:b:cc:dd:0:0:0.2.3.4", "a:b:cc:dd:ffff:ffff:255.2.3.4", []interface{}{0xa, 0xb, 0xcc, 0xdd, []uint64{0, 0xffffffffff}, 2, 3, 4}, nil, false) // [a, b, cc, dd, 0-ffffffffff, 2, 3, 4]
	t.testIncompatibleAddress1("a:b:cc:*.2.3.4", "a:b:cc:0:0:0:0.2.3.4", "a:b:cc:ffff:ffff:ffff:255.2.3.4", []interface{}{0xa, 0xb, 0xcc, []uint64{0, 0xffffffffffffff}, 2, 3, 4}, nil, false)     // [a, b, cc, 0-ffffffffffffff, 2, 3, 4]

	t.testSubnetStringRange1("1:2:4:5:6-9:7:8:f/ffff:0:ffff:0:ffff:0:ffff:0", "1:0:4:0:6:0:8:0", "1:0:4:0:9:0:8:0", []interface{}{1, 0, 4, 0, []uint{6, 9}, 0, 8, 0}, nil, false) //[1, 2, 4, 5, 6-9, 7, 8, f]
	t.testSubnetStringRange1("1:2:4:5-6:6:7:8:f/ffff:0:ffff:0:ffff:0:ffff:0", "1:0:4:0:6:0:8:0", "1:0:4:0:6:0:8:0", []interface{}{1, 0, 4, 0, 6, 0, 8, 0}, nil, true)             //[1, 2, 4, 5, 6-9, 7, 8, f]

	t.testSubnetStringRange1("1.*.*.*/11", "1.0.0.0", "1.255.255.255", []interface{}{1, []uint{0, 0xff}, []uint{0, 0xff}, []uint{0, 0xff}}, p11, true) //[1, 0-255, 0-255, 0-255]
	t.testSubnetStringRange1("1.*.*/32", "1.0.0.0", "1.255.255.255", []interface{}{1, []uint{0, 0xff}, []uint{0, 0xffff}}, p32, true)                  //[1, 0-255, 0-65535]
	t.testSubnetStringRange1("1.*/24", "1.0.0.0", "1.255.255.255", []interface{}{1, []uint{0, 0xffffff}}, p24, true)                                   //[1, 0-16777215]

	t.testSubnetStringRange("a:b:c:*:cc:d:e:f/64", "a:b:c:0:cc:d:e:f", "a:b:c:ffff:cc:d:e:f", []interface{}{0xa, 0xb, 0xc, []uint{0, 0xffff}, 0xcc, 0xd, 0xe, 0xf}, p64)                              //[a, b, c, 0-ffff, cc, d, e, f]
	t.testSubnetStringRange("a:*:cc:d:e:f/64", "a::cc:d:e:f", "a:ffff:ffff:ffff:cc:d:e:f", []interface{}{0xa, []uint64{0, 0xffffffffffff}, 0xcc, 0xd, 0xe, 0xf}, p64)                                 //[a, 0-ffffffffffff, cc, d, e, f]
	t.testSubnetStringRange("*:cc:d:e:f/64", "::cc:d:e:f", "ffff:ffff:ffff:ffff:cc:d:e:f", []interface{}{[]*big.Int{bigZeroConst(), setBigString("ffffffffffffffff", 16)}, 0xcc, 0xd, 0xe, 0xf}, p64) //[0-ffffffffffffffff, cc, d, e, f]

	//prefix subnets
	t.testSubnetStringRange("a:*::/64", "a::", "a:ffff::ffff:ffff:ffff:ffff", []interface{}{0xa, []uint{0, 0xffff}, []*big.Int{bigZeroConst(), setBigString("ffffffffffffffff", 16)}}, p64) //[a, 0-ffffffffffff, cc, d, e, f]
	t.testSubnetStringRange("1.128.0.0/11", "1.128.0.0", "1.159.255.255", []interface{}{1, []uint{128, 159}, []uint{0, 0xff}, []uint{0, 0xff}}, p11)                                        //[1, 0-255, 0-255, 0-255]

	if t.isLenient() {
		// inet_aton

		t.testSubnetStringRange("1.*.1/16", "1.0.0.1", "1.255.0.1", []interface{}{1, []uint{0, 0xff}, 1}, p16) //[1, 0-255, 1]
		t.testSubnetStringRange("*.1/16", "0.0.0.1", "255.0.0.1", []interface{}{[]uint{0, 0xff}, 1}, p16)      //[0-255, 1]
		t.testIncompatibleAddress("a:b:cc:*.4/112", "a:b:cc:0:0:0:0.0.0.4", "a:b:cc:ffff:ffff:ffff:255.0.0.4",
			[]interface{}{0xa, 0xb, 0xcc, []*big.Int{bigZeroConst(), setBigString("ffffffffffffff", 16)}, 4}, p112) //[a, b, cc, 0-ffffffffffffff, 4]
		t.testIncompatibleAddress("1:2:3:4:*.3.4/112", "1:2:3:4::0.3.0.4", "1:2:3:4:ffff:ffff:255.3.0.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffff}, 3, 4}, p112) //[1, 2, 3, 4, 0-ffffffffff, 3, 4]
		t.testIncompatibleAddress("1:2:3:4:*.4/112", "1:2:3:4::0.0.0.4", "1:2:3:4:ffff:ffff:255.0.0.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffff}, 4}, p112)      //[1, 2, 3, 4, 0-ffffffffff, 4]

		// prefix subnet

		t.testIncompatibleAddress("a:b:cc:*.0/112", "a:b:cc:0:0:0:0.0.0.0", "a:b:cc:ffff:ffff:ffff:255.0.255.255", []interface{}{0xa, 0xb, 0xcc, []*big.Int{bigZeroConst(), setBigString("ffffffffffffff", 16)}, []uint{0, 0xffff}}, p112) //[a, b, cc, 0-ffffffffffffff, 4]
	} else {
		// not inet_aton

		t.testSubnetStringRange("1.*.1/16", "1.0.0.1", "1.255.255.1", []interface{}{1, []uint{0, 0xffff}, 1}, p16)
		t.testSubnetStringRange("*.1/16", "0.0.0.1", "255.255.255.1", []interface{}{[]uint{0, 0xffffff}, 1}, p16)
		t.testIncompatibleAddress("a:b:cc:*.4/112", "a:b:cc:0:0:0:0.0.0.4", "a:b:cc:ffff:ffff:ffff:255.255.255.4", []interface{}{0xa, 0xb, 0xcc, []*big.Int{bigZeroConst(), setBigString("ffffffffffffffffff", 16)}, 4}, p112) //[a, b, cc, 0-ffffffffffffffffff, 4]
		t.testSubnetStringRange("1:2:3:4:*.3.4/112", "1:2:3:4::0.0.3.4", "1:2:3:4:ffff:ffff:255.255.3.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffffff}, 3, 4}, p112)                                                  //[1, 2, 3, 4, 0-ffffffffffff, 3, 4]
		t.testIncompatibleAddress("1:2:3:4:*.4/112", "1:2:3:4::0.0.0.4", "1:2:3:4:ffff:ffff:255.255.255.4", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffffffff}, 4}, p112)                                                 //[1, 2, 3, 4, 0-ffffffffffffff, 4]

		// prefix subnet

		t.testSubnetStringRange("1:2:3:4:*.0.0/112", "1:2:3:4::0.0.0.0", "1:2:3:4:ffff:ffff:255.255.255.255", []interface{}{1, 2, 3, 4, []uint64{0, 0xffffffffffff}, []uint{0, 0xff}, []uint{0, 0xff}}, p112) //[1, 2, 3, 4, 0-ffffffffffffff, 4]

	}
	// prefix subnet

	t.testSubnetStringRange("a:b:cc::0.0.0.0/64", "a:b:cc:0:0:0:0.0.0.0", "a:b:cc::ffff:ffff:255.255.255.255",
		[]interface{}{0xa, 0xb, 0xcc, []uint64{0, 0xffffffff}, []uint{0, 0xff}, []uint{0, 0xff}, []uint{0, 0xff}, []uint{0, 0xff}}, p64) //[a, b, cc, 0-ffffffffffffff, 4]

	t.testSubnetStringRange("1-2.3.4-5.6/16", "1.3.4.6", "2.3.5.6", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, 6}, p16) //[1-2, 3, 4-5, 6]
	t.testSubnetStringRange("1-2.3.4-5.0/23", "1.3.4.0", "2.3.5.0", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, 0}, p23) //[1-2, 3, 4-5, 6]

	t.testSubnetStringRange("1-2.3.4-5.0/24", "1.3.4.0", "2.3.5.255", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, []uint{0, 0xff}}, p24) //[1-2, 3, 4-5, 6]

	t.testSubnetStringRange("1-2:3:4-5:6::/48", "1:3:4:6::", "2:3:5:6::", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, 6, 0}, p48) //[1-2, 3, 4-5, 6, 0]

	t.testSubnetStringRange("1-2:3:4-5::/48", "1:3:4::", "2:3:5:ffff:ffff:ffff:ffff:ffff", []interface{}{[]uint{1, 2}, 3, []uint{4, 5}, []*big.Int{bigZeroConst(), setBigString("ffffffffffffffffffff", 16)}}, p48) //[1-2, 3, 4-5, 6, 0]

	t.testIncompatibleAddress("1:2:3:4:5:6:1-3.2.0.0/112", "1:2:3:4:5:6:1.2.0.0", "1:2:3:4:5:6:3.2.255.255", []interface{}{1, 2, 3, 4, 5, 6, []uint{1, 3}, 2, []uint{0, 0xff}, []uint{0, 0xff}}, p112) //[1, 2, 3, 4, 5, 6, 1-3, 2, 0, 4-5]

	t.testIncompatibleAddress("1:2:3:4:5:6:1-3.2.0.4-5/112", "1:2:3:4:5:6:1.2.0.4", "1:2:3:4:5:6:3.2.0.5", []interface{}{1, 2, 3, 4, 5, 6, []uint{1, 3}, 2, 0, []uint{4, 5}}, p112) //[1, 2, 3, 4, 5, 6, 1-3, 2, 0, 4-5]

	t.testSubnetStringRange1("1-3.1-3.1-3.1-3/175.80.81.83",
		"1.0.0.1", "3.0.1.3",
		[]interface{}{[]int{1, 3}, 0, []int{0, 1}, []int{1, 3}},
		nil, false)

	t.testMaskedIncompatibleAddress("*.*/202.63.240.51", "0.0.0.0", "202.63.240.51") //10101010 00111111 11110000 00110011
	t.testMaskedIncompatibleAddress("*.*/63.240.51.202", "0.0.0.0", "63.240.51.202")
	t.testMaskedIncompatibleAddress("*.*/240.51.202.63", "0.0.0.0", "240.51.202.63")
	t.testMaskedIncompatibleAddress("*.*/51.202.63.240", "0.0.0.0", "51.202.63.240")

	t.testMaskedIncompatibleAddress("*.*.*.*/202.63.240.51", "0.0.0.0", "202.63.240.51")
	t.testMaskedIncompatibleAddress("*.*.*.*/63.240.51.202", "0.0.0.0", "63.240.51.202")
	t.testMaskedIncompatibleAddress("*.*.*.*/240.51.202.63", "0.0.0.0", "240.51.202.63")
	t.testMaskedIncompatibleAddress("*.*.*.*/51.202.63.240", "0.0.0.0", "51.202.63.240")

	t.testMaskedIncompatibleAddress("*:aaaa:bbbb:cccc/abcd:dcba:aaaa:bbbb:cccc::dddd",
		"::cccc", "abcd:dcba:aaaa:bbbb:cccc::cccc")
	t.testMaskedIncompatibleAddress("aaaa:bbbb:*:cccc/abcd:dcba:aaaa:bbbb:cccc::dddd",
		"aa88:98ba::cccc", "aa88:98ba:aaaa:bbbb:cccc::cccc")
	t.testMaskedIncompatibleAddress("aaaa:bbbb:*/abcd:dcba:aaaa:bbbb:cccc::dddd",
		"aa88:98ba::", "aa88:98ba:aaaa:bbbb:cccc::dddd")

	t.testMaskedIncompatibleAddress("*.*/63.255.15.0", "0.0.0.0", "63.255.15.0")

	t.testSubnetStringRange1("*.*/63.15.255.255",
		"0.0.0.0", "63.15.255.255",
		[]interface{}{[]int{0, 63}, []int{0, 0xfffff}},
		nil, false)

	t.testPrefix("25:51:27:*:*:*:*:*", nil, 48, p48)
	t.testPrefix("25:51:27:*:*:*:*:*/48", p48, 48, p48)
	t.testPrefix("25:50-51:27::/48", p48, 48, nil)
	t.testPrefix("25:50-51:27:*:*:*:*:*", nil, 48, nil)
	t.testPrefix("25:51:27:12:82:55:2:2", nil, 128, p128)
	t.testPrefix("*:*:*:*:*:*:*:*", nil, 0, p0)
	t.testPrefix("*:*:*:*:*:*:0-fe:*", nil, 112, nil)
	t.testPrefix("*:*:*:*:*:*:0-ff:*", nil, 104, nil)
	t.testPrefix("*:*:*:*:*:*:0-ffff:*", nil, 0, p0)
	t.testPrefix("*:*:*:*:*:*:0-7fff:*", nil, 97, nil)
	t.testPrefix("*:*:*:*:*:*:8000-ffff:*", nil, 97, nil)
	t.testPrefix("*.*.*.*", nil, 0, p0)
	t.testPrefix("3.*.*.*", nil, 8, p8)
	t.testPrefix("3.*.*.1-3", nil, 32, nil)
	t.testPrefix("3.0-127.*.*", nil, 9, p9)
	t.testPrefix("3.128-255.*.*", nil, 9, p9)

	//TODO soon
	//testMasked("1.*.3.4", null, null, "1.*.3.4");
	//testMasked("1.*.3.4/255.255.1.0", "255.255.1.0", null, "1.*.1.0");
	//testMasked("1.*.3.4/255.255.254.0", "255.255.254.0", 23, false ? "1.*.2.0/23" : "1.*.3.4/23");
	//
	//testMasked("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", null, null, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff");
	//testMasked("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/0:101:0:101:0:101:0:101", "0:101:0:101:0:101:0:101", null, "0:101:0:101:0:101:0:101");
	//testMasked("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/ffff:ffff:8000::", "ffff:ffff:8000::", 33, false ? "ffff:ffff:8000::/33" : "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/33");
	//testMasked("ffff:ffff::/ffff:ffff:8000::", "ffff:ffff:8000::", 33, "ffff:ffff::/33");

	t.ipAddressTester.run()
}

func setBigString(str string, base int) *big.Int {
	res, b := new(big.Int).SetString(str, base)
	if !b {
		panic("bad string for big int")
	}
	return res
}

func (t ipAddressRangeTester) ipv4rangestest(pass bool, x string, ipv4RangeOptions, ipv6RangeOptions ipaddr.RangeParameters) {
	t.iprangestest(pass, x, false, false, true, ipv4RangeOptions, ipv6RangeOptions)
}

func (t ipAddressRangeTester) ipv4rangetest(pass bool, x string, rangeOptions ipaddr.RangeParameters) {
	t.iprangetest(pass, x, false, false, true, rangeOptions)
}

func (t ipAddressRangeTester) ipv6rangestest(pass bool, x string, ipv4Options, ipv6Options ipaddr.RangeParameters) {
	t.iprangestest(pass, x, false, false, false, ipv4Options, ipv6Options)
}

func (t ipAddressRangeTester) ipv6rangetest(pass bool, x string, options ipaddr.RangeParameters) {
	t.iprangetest(pass, x, false, false, false, options)
}

func (t ipAddressRangeTester) iprangestest(pass bool, x string, isZero, notBoth, ipv4Test bool, ipv4RangeOptions, ipv6RangeOptions ipaddr.RangeParameters) {

	addr := t.createDoubleParametrizedAddress(x, ipv4RangeOptions, ipv6RangeOptions)
	if t.iptest(pass, addr, isZero, notBoth, ipv4Test) {
		//do it a second time to test the caching
		t.iptest(pass, addr, isZero, notBoth, ipv4Test)
	}
}

func (t ipAddressRangeTester) iprangetest(pass bool, x string, isZero, notBoth, ipv4Test bool, rangeOptions ipaddr.RangeParameters) {
	addr := t.createParametrizedAddress(x, rangeOptions)
	if t.iptest(pass, addr, isZero, notBoth, ipv4Test) {
		//do it a second time to test the caching
		t.iptest(pass, addr, isZero, notBoth, ipv4Test)
	}
}

func (t ipAddressRangeTester) testPrefix(original string, prefixLength ipaddr.PrefixLen, minPrefix ipaddr.BitCount, equivalentPrefix ipaddr.PrefixLen) {
	addr := t.createAddress(original).GetAddress()
	t.testBase.testPrefix(addr, prefixLength, minPrefix, equivalentPrefix)
	t.incrementTestCount()
}

func (t ipAddressRangeTester) testTrees() {

	t.testTree("1.2.3.4", []string{
		"1.2.3.4",
		"1.2.3.*",
		"1.2.*.*",
		"1.*.*.*",
		"*.*.*.*",
		"*",
	})

	t.testTree("1.2.3.*", []string{
		"1.2.3.*",
		"1.2.*.*",
		"1.*.*.*",
		"*.*.*.*",
		"*",
	})

	t.testTree("1.2.*.*", []string{
		"1.2.*.*",
		"1.*.*.*",
		"*.*.*.*",
		"*",
	})

	t.testTree("a:b:c:d:e:f:a:b", []string{
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a::/112",
		"a:b:c:d:e:f::/96",
		"a:b:c:d:e::/80",
		"a:b:c:d::/64",
		"a:b:c::/48",
		"a:b::/32",
		"a::/16",
		"::/0",
		"*",
	})

	t.testTree("1.2.3.4/28", []string{
		"1.2.3.4/28",
		"1.2.3.4/24",
		"1.2.0.4/16",
		"1.0.0.4/8",
		"0.0.0.4/0",
	})
	t.testTree("1.2.3.4/17", []string{
		"1.2.3.4/17",
		"1.2.3.4/16",
		"1.0.3.4/8",
		"0.0.3.4/0",
	})
	t.testTree("a:b:c:d:e:f:a:b/97", []string{
		"a:b:c:d:e:f:a:b/97",
		"a:b:c:d:e:f:a:b/96",
		"a:b:c:d:e::a:b/80",
		"a:b:c:d::a:b/64",
		"a:b:c::a:b/48",
		"a:b::a:b/32",
		"a::a:b/16",
		"::a:b/0",
	})
	t.testTree("a:b:c:d:e:f:ffff:b/97", []string{
		"a:b:c:d:e:f:ffff:b/97",
		"a:b:c:d:e:f:7fff:b/96",
		"a:b:c:d:e::7fff:b/80",
		"a:b:c:d::7fff:b/64",
		"a:b:c::7fff:b/48",
		"a:b::7fff:b/32",
		"a::7fff:b/16",
		"::7fff:b/0",
	})
	t.testTree("a:b:c:d:e:f:a:b/96", []string{
		"a:b:c:d:e:f:a:b/96",
		"a:b:c:d:e::a:b/80",
		"a:b:c:d::a:b/64",
		"a:b:c::a:b/48",
		"a:b::a:b/32",
		"a::a:b/16",
		"::a:b/0",
	})

	t.testTree("a:b:c:d::a:b", []string{
		"a:b:c:d::a:b",
		"a:b:c:d:0:0:a::/112",
		"a:b:c:d::/96",
		"a:b:c:d::/80",
		"a:b:c:d::/64",
		"a:b:c::/48",
		"a:b::/32",
		"a::/16",
		"::/0",
		"*",
	})
	t.testTree("::c:d:e:f:a:b", []string{
		"::c:d:e:f:a:b",
		"0:0:c:d:e:f:a::/112",
		"0:0:c:d:e:f::/96",
		"0:0:c:d:e::/80",
		"0:0:c:d::/64",
		"0:0:c::/48",
		"::/32",
		"::/16",
		"::/0",
		"*",
	})
}

func (t ipAddressRangeTester) testTree(start string, parents []string) {
	str := t.createAddress(start)
	originaLabelStr := str
	labelStr := str
	originalPrefixed := str.IsPrefixed()
	if !originalPrefixed {
		address := str.GetAddress()
		//convert 1.2.3.* to 1.2.3.*/24 which is needed by adjustPrefixBySegment
		address = address.AssignPrefixForSingleBlock()
		str = address.ToAddressString()
	}

	original := str
	i := 0
	var last *ipaddr.IPAddressString
	for {
		label := getLabel(labelStr)
		expected := parents[i]
		if label != expected {
			t.addFailure(newFailure("failed expected: "+expected+" actual: "+label, str))
			break
		}
		last = str
		str = enlargeSubnetStr(str)
		if str == nil || last == str {
			break
		}
		labelStr = str
		i++
	}

	//now do the same thing but use the IPAddress objects instead
	labelStr = originaLabelStr
	str = original
	i = 0
	for {
		label := getLabel(labelStr)
		expected := parents[i]
		if label != expected {
			t.addFailure(newFailure("failed expected: "+expected+" actual: "+label, str))
			break
		}
		labelAddr := enlargeSubnet(str.GetAddress())
		//IPAddress labelAddr = str.getAddress().adjustPrefixBySegment(false);
		//IPAddress subnetAddr = labelAddr.toPrefixBlock(labelAddr.getNetworkPrefixLength());
		//if(labelAddr != subnetAddr) {
		//addFailure(new Failure("not already a subnet " + labelAddr + " expected: " + subnetAddr, labelAddr));
		//}
		str = labelAddr.ToAddressString()
		labelStr = str
		if *str.GetNetworkPrefixLen() == 0 { //when network prefix is 0, IPAddress.adjustPrefixBySegment() returns the same address
			break
		}
		i++
	}
	t.incrementTestCount()
}

func (t ipAddressRangeTester) testIPv4IPAddrStrings(w *ipaddr.IPAddressString, ipAddr *ipaddr.IPAddress, normalizedString, normalizedWildcardString, sqlString, fullString, octalString, hexString, reverseDNSString,
	singleHex,
	singleOctal string) {
	t.testBase.testStrings(w, ipAddr, normalizedString, normalizedWildcardString, normalizedWildcardString, sqlString, fullString,
		normalizedString, normalizedString, normalizedWildcardString, normalizedString, normalizedWildcardString, reverseDNSString, normalizedString,
		singleHex, singleOctal)

	//now test some IPv4-only strings
	t.testIPv4OnlyStrings(w, ipAddr.ToIPv4Address(), octalString, hexString)
	t.testInetAtonCombos(w, ipAddr.ToIPv4Address())
}

func (t ipAddressRangeTester) testIPv4OnlyStrings(w *ipaddr.IPAddressString, ipAddr *ipaddr.IPv4Address, octalString, hexString string) {
	oct := ipAddr.ToInetAtonString(ipaddr.Inet_aton_radix_octal)
	hex := ipAddr.ToInetAtonString(ipaddr.Inet_aton_radix_hex)
	octMatch := oct == octalString
	if !octMatch {
		t.addFailure(newFailure("failed expected: "+octalString+" actual: "+oct, w))
	} else {
		hexMatch := hex == hexString
		if !hexMatch {
			t.addFailure(newFailure("failed expected: "+hexString+" actual: "+hex, w))
		}
	}
	t.incrementTestCount()
}

func (t ipAddressRangeTester) testInetAtonCombos(w *ipaddr.IPAddressString, ipAddr *ipaddr.IPv4Address) {
	vals := []ipaddr.Inet_aton_radix{ipaddr.Inet_aton_radix_octal, ipaddr.Inet_aton_radix_hex, ipaddr.Inet_aton_radix_decimal}
	for _, radix := range vals {
		for i := 0; i < ipaddr.IPv4SegmentCount; i++ {
			//try {
			str, e := ipAddr.ToInetAtonJoinedString(radix, i)
			if e != nil {
				//verify this case: joining segments results in a joined segment that is not a contiguous range
				section := ipAddr.GetSection()
				verifiedIllegalJoin := false
				for j := section.GetSegmentCount() - i - 1; j < section.GetSegmentCount()-1; j++ {
					if section.GetSegment(j).IsMultiple() {
						for j++; j < section.GetSegmentCount(); j++ {
							if !section.GetSegment(j).IsFullRange() {
								verifiedIllegalJoin = true
								break
							}
						}
					}
				}
				if !verifiedIllegalJoin {
					t.addFailure(newFailure("failed expected: "+ipAddr.String()+" actual: "+e.Error(), w))
				}
			} else {
				parsed := ipaddr.NewIPAddressStringParams(str, inetAtonwildcardAndRangeOptions)
				// try{
				parsedValue := parsed.GetAddress()
				if !ipAddr.Equals(parsedValue) {
					t.addFailure(newFailure("failed expected: "+ipAddr.String()+" actual: "+parsedValue.String(), w))
				} else {
					//int pos;
					origStr := str
					count := 0
					pos := -1
					for pos = strings.IndexByte(str, ipaddr.IPv4SegmentSeparator); pos >= 0 && pos < len(str); {
						//for ((pos = str.indexOf(ipaddr.IPv4SegmentSeparator)) >= 0){
						str = str[pos+1:]
						pos = strings.IndexByte(str, ipaddr.IPv4SegmentSeparator)
						count++
					}
					if ipaddr.IPv4SegmentCount-1-i != count {
						failStr := "failed expected separator count in " + origStr + ": " + strconv.Itoa(ipaddr.IPv4SegmentCount-1-i) + " actual separator count: " + strconv.Itoa(count)
						t.addFailure(newFailure(failStr, w))

						//str = origStr
						//count = 0
						//pos := -1
						//for pos = strings.IndexByte(str, ipaddr.IPv4SegmentSeparator); pos >= 0 && pos < len(str); {
						//	//for ((pos = str.indexOf(ipaddr.IPv4SegmentSeparator)) >= 0){
						//	str = str[pos+1:]
						//	count++
						//}
						//fmt.Println("WTF")
					}
				}
			}
			t.incrementTestCount()
		}
	}
}

func (t ipAddressRangeTester) testIPv4Strings(addr, normalizedString, normalizedWildcardString, sqlString, fullString, octalString, hexString, reverseDNSString, singleHex, singleOctal string) {
	w := t.createAddress(addr)
	ipAddr := w.GetAddress()
	//createList(w);
	t.testIPv4IPAddrStrings(w, ipAddr, normalizedString, normalizedWildcardString, sqlString, fullString, octalString, hexString, reverseDNSString, singleHex, singleOctal)
}

func (t ipAddressRangeTester) testIPv6Strings(addr,
	normalizedString,
	normalizedWildcardString,
	canonicalWildcardString,
	sqlString,
	fullString,
	compressedString,
	canonicalString,
	subnetString,
	compressedWildcardString,
	mixedStringNoCompressMixed,
	mixedStringNoCompressHost,
	mixedStringCompressCoveredHost,
	mixedString,
	reverseDNSString,
	uncHostString,
	base85String,
	singleHex,
	singleOctal string) {
	w := t.createAddress(addr)
	ipAddr := w.GetAddress()

	//createList(w);

	t.testBase.testIPv6Strings(w,
		ipAddr,
		normalizedString,
		normalizedWildcardString,
		canonicalWildcardString,
		sqlString,
		fullString,
		compressedString,
		canonicalString,
		subnetString,
		compressedWildcardString,
		mixedStringNoCompressMixed,
		mixedStringNoCompressHost,
		mixedStringCompressCoveredHost,
		mixedString,
		reverseDNSString,
		uncHostString,
		base85String,
		singleHex,
		singleOctal)
}

//each ipv4 failure is 6, each ipv6 is 10, current total is 520

func (t ipAddressRangeTester) testStrings() {

	//boolean allPrefixesAreSubnets = prefixConfiguration.allPrefixedAddressesAreSubnets();
	//boolean isNoAutoSubnets = prefixConfiguration.prefixedSubnetsAreExplicit();

	t.testIPv4Strings("1.2.3.4", "1.2.3.4", "1.2.3.4", "1.2.3.4", "001.002.003.004", "01.02.03.04", "0x1.0x2.0x3.0x4", "4.3.2.1.in-addr.arpa", "0x01020304", "000100401404")

	t.testIPv4Strings("1.2.3.4/16", "1.2.3.4/16", "1.2.3.4", "1.2.3.4", "001.002.003.004/16", "01.02.03.04/16", "0x1.0x2.0x3.0x4/16", "4.3.2.1.in-addr.arpa", "0x01020304", "000100401404")

	t.testIPv4Strings("1.2.*.*", "1.2.*.*", "1.2.*.*", "1.2.%.%", "001.002.000-255.000-255", "01.02.*.*", "0x1.0x2.*.*", "*.*.2.1.in-addr.arpa", "0x01020000-0x0102ffff", "000100400000-000100577777") //note that wildcards are never converted to CIDR.
	t.testIPv4Strings("1.2.*", "1.2.*.*", "1.2.*.*", "1.2.%.%", "001.002.000-255.000-255", "01.02.*.*", "0x1.0x2.*.*", "*.*.2.1.in-addr.arpa", "0x01020000-0x0102ffff", "000100400000-000100577777")

	t.testIPv4Strings("1.2.*.*/16", "1.2.0.0/16", "1.2.*.*", "1.2.%.%", "001.002.000.000/16", "01.02.00.00/16", "0x1.0x2.0x0.0x0/16", "*.*.2.1.in-addr.arpa", "0x01020000-0x0102ffff", "000100400000-000100577777")
	t.testIPv4Strings("1.2.*/16", "1.2.0.0/16", "1.2.*.*", "1.2.%.%", "001.002.000.000/16", "01.02.00.00/16", "0x1.0x2.0x0.0x0/16", "*.*.2.1.in-addr.arpa", "0x01020000-0x0102ffff", "000100400000-000100577777")
	t.testIPv4Strings("1.*.*/16", "1.*.0.0/16", "1.*.*.*", "1.%.%.%", "001.000-255.000.000/16", "01.*.00.00/16", "0x1.*.0x0.0x0/16", "*.*.*.1.in-addr.arpa", "0x01000000-0x01ffffff", "000100000000-000177777777")

	t.testIPv4Strings("0.0.0.0", "0.0.0.0", "0.0.0.0", "0.0.0.0", "000.000.000.000", "00.00.00.00", "0x0.0x0.0x0.0x0", "0.0.0.0.in-addr.arpa", "0x00000000", "000000000000")
	t.testIPv4Strings("9.63.127.254", "9.63.127.254", "9.63.127.254", "9.63.127.254", "009.063.127.254", "011.077.0177.0376", "0x9.0x3f.0x7f.0xfe", "254.127.63.9.in-addr.arpa", "0x093f7ffe", "001117677776")

	t.testIPv4Strings("9.63.127.254/16", "9.63.127.254/16", "9.63.127.254", "9.63.127.254", "009.063.127.254/16", "011.077.0177.0376/16", "0x9.0x3f.0x7f.0xfe/16", "254.127.63.9.in-addr.arpa", "0x093f7ffe", "001117677776")

	t.testIPv4Strings("9.63.*.*", "9.63.*.*", "9.63.*.*", "9.63.%.%", "009.063.000-255.000-255", "011.077.*.*", "0x9.0x3f.*.*", "*.*.63.9.in-addr.arpa", "0x093f0000-0x093fffff", "001117600000-001117777777") //note that wildcards are never converted to CIDR.
	t.testIPv4Strings("9.63.*", "9.63.*.*", "9.63.*.*", "9.63.%.%", "009.063.000-255.000-255", "011.077.*.*", "0x9.0x3f.*.*", "*.*.63.9.in-addr.arpa", "0x093f0000-0x093fffff", "001117600000-001117777777")

	t.testIPv4Strings("9.63.*.*/16", "9.63.0.0/16", "9.63.*.*", "9.63.%.%", "009.063.000.000/16", "011.077.00.00/16", "0x9.0x3f.0x0.0x0/16", "*.*.63.9.in-addr.arpa", "0x093f0000-0x093fffff", "001117600000-001117777777")
	t.testIPv4Strings("9.63.*/16", "9.63.0.0/16", "9.63.*.*", "9.63.%.%", "009.063.000.000/16", "011.077.00.00/16", "0x9.0x3f.0x0.0x0/16", "*.*.63.9.in-addr.arpa", "0x093f0000-0x093fffff", "001117600000-001117777777")
	t.testIPv4Strings("9.*.*/16", "9.*.0.0/16", "9.*.*.*", "9.%.%.%", "009.000-255.000.000/16", "011.*.00.00/16", "0x9.*.0x0.0x0/16", "*.*.*.9.in-addr.arpa", "0x09000000-0x09ffffff", "001100000000-001177777777")

	t.testIPv4Strings("1.2.3.250-255", "1.2.3.250-255", "1.2.3.250-255", "1.2.3.25_", "001.002.003.250-255", "01.02.03.0372-0377", "0x1.0x2.0x3.0xfa-0xff", "250-255.3.2.1.in-addr.arpa", "0x010203fa-0x010203ff", "000100401772-000100401777")
	t.testIPv4Strings("1.2.3.200-255", "1.2.3.200-255", "1.2.3.200-255", "1.2.3.2__", "001.002.003.200-255", "01.02.03.0310-0377", "0x1.0x2.0x3.0xc8-0xff", "200-255.3.2.1.in-addr.arpa", "0x010203c8-0x010203ff", "000100401710-000100401777")
	t.testIPv4Strings("1.2.3.100-199", "1.2.3.100-199", "1.2.3.100-199", "1.2.3.1__", "001.002.003.100-199", "01.02.03.0144-0307", "0x1.0x2.0x3.0x64-0xc7", "100-199.3.2.1.in-addr.arpa", "0x01020364-0x010203c7", "000100401544-000100401707")
	t.testIPv4Strings("100-199.2.3.100-199", "100-199.2.3.100-199", "100-199.2.3.100-199", "1__.2.3.1__", "100-199.002.003.100-199", "0144-0307.02.03.0144-0307", "0x64-0xc7.0x2.0x3.0x64-0xc7", "100-199.3.2.100-199.in-addr.arpa", "", "")
	t.testIPv4Strings("100-199.2.3.100-198", "100-199.2.3.100-198", "100-199.2.3.100-198", "1__.2.3.100-198", "100-199.002.003.100-198", "0144-0307.02.03.0144-0306", "0x64-0xc7.0x2.0x3.0x64-0xc6", "100-198.3.2.100-199.in-addr.arpa", "", "")
	t.testIPv4Strings("1.2.3.0-99", "1.2.3.0-99", "1.2.3.0-99", "1.2.3.0-99", "001.002.003.000-099", "01.02.03.00-0143", "0x1.0x2.0x3.0x0-0x63", "0-99.3.2.1.in-addr.arpa", "0x01020300-0x01020363", "000100401400-000100401543")
	t.testIPv4Strings("1.2.3.100-155", "1.2.3.100-155", "1.2.3.100-155", "1.2.3.100-155", "001.002.003.100-155", "01.02.03.0144-0233", "0x1.0x2.0x3.0x64-0x9b", "100-155.3.2.1.in-addr.arpa", "0x01020364-0x0102039b", "000100401544-000100401633")
	t.testIPv4Strings("1.2.3.100-255", "1.2.3.100-255", "1.2.3.100-255", "1.2.3.100-255", "001.002.003.100-255", "01.02.03.0144-0377", "0x1.0x2.0x3.0x64-0xff", "100-255.3.2.1.in-addr.arpa", "0x01020364-0x010203ff", "000100401544-000100401777")

	t.testIPv4Strings("1.129-254.5.5/12", "1.129-254.5.5/12", "1.129-254.5.5", "1.129-254.5.5", "001.129-254.005.005/12", "01.0201-0376.05.05/12", "0x1.0x81-0xfe.0x5.0x5/12", "5.5.129-254.1.in-addr.arpa", "", "")
	t.testIPv4Strings("1.2__.5.5/14", "1.200-255.5.5/14", "1.200-255.5.5", "1.2__.5.5", "001.200-255.005.005/14", "01.0310-0377.05.05/14", "0x1.0xc8-0xff.0x5.0x5/14", "5.5.200-255.1.in-addr.arpa", "", "")
	t.testIPv4Strings("1.*.5.5/12", "1.*.5.5/12", "1.*.5.5", "1.%.5.5", "001.000-255.005.005/12", "01.*.05.05/12", "0x1.*.0x5.0x5/12", "5.5.*.1.in-addr.arpa", "", "")
	//OK we are testing 01.*.02405/12 and our bounds check for inet_aton does not work because later when creating address it is not treated as inet_aton due to the *
	//so when we do the bounds checking for inet_aton we need to check for * and only test with single segment boundaries
	//also check for that setting where * extends beyond single segment

	t.testIPv6Strings("::",
		"0:0:0:0:0:0:0:0",
		"0:0:0:0:0:0:0:0",
		"::",
		"0:0:0:0:0:0:0:0",
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"::",
		"::",
		"::",
		"::",
		"::0.0.0.0",
		"::",
		"::",
		"::",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-0-0-0-0-0-0.ipv6-literal.net",
		"00000000000000000000",
		"0x00000000000000000000000000000000",
		"00000000000000000000000000000000000000000000")

	t.testIPv6Strings("::2",
		"0:0:0:0:0:0:0:2",
		"0:0:0:0:0:0:0:2",
		"::2",
		"0:0:0:0:0:0:0:2",
		"0000:0000:0000:0000:0000:0000:0000:0002",
		"::2",
		"::2",
		"::2",
		"::2",
		"::0.0.0.2",
		"::0.0.0.2",
		"::0.0.0.2",
		"::0.0.0.2",
		"2.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-0-0-0-0-0-2.ipv6-literal.net",
		"00000000000000000002",
		"0x00000000000000000000000000000002",
		"00000000000000000000000000000000000000000002")

	t.testIPv6Strings("::7fff:ffff:ffff:ffff",
		"0:0:0:0:7fff:ffff:ffff:ffff",
		"0:0:0:0:7fff:ffff:ffff:ffff",
		"::7fff:ffff:ffff:ffff",
		"0:0:0:0:7fff:ffff:ffff:ffff",
		"0000:0000:0000:0000:7fff:ffff:ffff:ffff",
		"::7fff:ffff:ffff:ffff",
		"::7fff:ffff:ffff:ffff",
		"::7fff:ffff:ffff:ffff",
		"::7fff:ffff:ffff:ffff",
		"::7fff:ffff:255.255.255.255",
		"::7fff:ffff:255.255.255.255",
		"::7fff:ffff:255.255.255.255",
		"::7fff:ffff:255.255.255.255",
		"f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.7.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-0-0-7fff-ffff-ffff-ffff.ipv6-literal.net",
		"0000000000d*-h_{Y}sg",
		"0x00000000000000007fffffffffffffff",
		"00000000000000000000000777777777777777777777")

	t.testIPv6Strings("0:0:0:1::",
		"0:0:0:1:0:0:0:0",
		"0:0:0:1:0:0:0:0",
		"0:0:0:1::",
		"0:0:0:1:0:0:0:0",
		"0000:0000:0000:0001:0000:0000:0000:0000",
		"0:0:0:1::",
		"0:0:0:1::",
		"0:0:0:1::",
		"0:0:0:1::",
		"::1:0:0:0.0.0.0",
		"0:0:0:1::",
		"0:0:0:1::",
		"0:0:0:1::",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-0-1-0-0-0-0.ipv6-literal.net",
		"0000000000_sw2=@*|O1",
		"0x00000000000000010000000000000000",
		"00000000000000000000002000000000000000000000")

	t.testIPv6Strings("::8fff:ffff:ffff:ffff",
		"0:0:0:0:8fff:ffff:ffff:ffff",
		"0:0:0:0:8fff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff",
		"0:0:0:0:8fff:ffff:ffff:ffff",
		"0000:0000:0000:0000:8fff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff",
		"::8fff:ffff:255.255.255.255",
		"::8fff:ffff:255.255.255.255",
		"::8fff:ffff:255.255.255.255",
		"::8fff:ffff:255.255.255.255",
		"f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.8.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-0-0-8fff-ffff-ffff-ffff.ipv6-literal.net",
		"0000000000i(`c)xypow",
		"0x00000000000000008fffffffffffffff",
		"00000000000000000000001077777777777777777777")

	t.testIPv6Strings("::8fff:ffff:ffff:ffff:ffff",
		"0:0:0:8fff:ffff:ffff:ffff:ffff",
		"0:0:0:8fff:ffff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff:ffff",
		"0:0:0:8fff:ffff:ffff:ffff:ffff",
		"0000:0000:0000:8fff:ffff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:ffff:ffff",
		"::8fff:ffff:ffff:255.255.255.255",
		"::8fff:ffff:ffff:255.255.255.255",
		"::8fff:ffff:ffff:255.255.255.255",
		"::8fff:ffff:ffff:255.255.255.255",
		"f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.f.8.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-0-8fff-ffff-ffff-ffff-ffff.ipv6-literal.net",
		"00000004&U-n{rbbza$w",
		"0x0000000000008fffffffffffffffffff",
		"00000000000000000217777777777777777777777777")

	t.testIPv6Strings("a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"000a:000b:000c:000d:000e:000f:000a:000b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:0.10.0.11",
		"a:b:c:d:e:f:0.10.0.11",
		"a:b:c:d:e:f:0.10.0.11",
		"a:b:c:d:e:f:0.10.0.11",
		"b.0.0.0.a.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-e-f-a-b.ipv6-literal.net",
		"00|N0s0$ND2DCD&%D3QB",
		"0x000a000b000c000d000e000f000a000b",
		"00000240001300006000032000160000740002400013")

	t.testIPv6Strings("a:b:c:d:e:f:a:b/64",
		"a:b:c:d:e:f:a:b/64",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:a:b",
		"000a:000b:000c:000d:000e:000f:000a:000b/64",
		"a:b:c:d:e:f:a:b/64",
		"a:b:c:d:e:f:a:b/64",
		"a:b:c:d:e:f:a:b/64",
		"a:b:c:d:e:f:a:b",
		"a:b:c:d:e:f:0.10.0.11/64",
		"a:b:c:d:e:f:0.10.0.11/64",
		"a:b:c:d:e:f:0.10.0.11/64",
		"a:b:c:d:e:f:0.10.0.11/64",
		"b.0.0.0.a.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-e-f-a-b.ipv6-literal.net/64",
		"00|N0s0$ND2DCD&%D3QB/64",
		"0x000a000b000c000d000e000f000a000b",
		"00000240001300006000032000160000740002400013")
	t.testIPv6Strings("::c:d:e:f:a:b/64",
		"0:0:c:d:e:f:a:b/64",
		"0:0:c:d:e:f:a:b",
		"::c:d:e:f:a:b",
		"0:0:c:d:e:f:a:b",
		"0000:0000:000c:000d:000e:000f:000a:000b/64",
		"::c:d:e:f:a:b/64",
		"::c:d:e:f:a:b/64",
		"::c:d:e:f:a:b/64",
		"::c:d:e:f:a:b",
		"::c:d:e:f:0.10.0.11/64",
		"::c:d:e:f:0.10.0.11/64",
		"::c:d:e:f:0.10.0.11/64",
		"::c:d:e:f:0.10.0.11/64",
		"b.0.0.0.a.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-c-d-e-f-a-b.ipv6-literal.net/64",
		"0000001G~Ie^C9jXExx>/64",
		"0x00000000000c000d000e000f000a000b",
		"00000000000000006000032000160000740002400013")

	t.testIPv6Strings("::c:d:e:f:a:b",
		"0:0:c:d:e:f:a:b",
		"0:0:c:d:e:f:a:b",
		"::c:d:e:f:a:b",
		"0:0:c:d:e:f:a:b",
		"0000:0000:000c:000d:000e:000f:000a:000b",
		"::c:d:e:f:a:b",
		"::c:d:e:f:a:b",
		"::c:d:e:f:a:b",
		"::c:d:e:f:a:b",
		"::c:d:e:f:0.10.0.11",
		"::c:d:e:f:0.10.0.11",
		"::c:d:e:f:0.10.0.11",
		"::c:d:e:f:0.10.0.11",
		"b.0.0.0.a.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
		"0-0-c-d-e-f-a-b.ipv6-literal.net",
		"0000001G~Ie^C9jXExx>",
		"0x00000000000c000d000e000f000a000b",
		"00000000000000006000032000160000740002400013")

	t.testIPv6Strings("a:b:c:d::",
		"a:b:c:d:0:0:0:0",
		"a:b:c:d:0:0:0:0",
		"a:b:c:d::",
		"a:b:c:d:0:0:0:0",
		"000a:000b:000c:000d:0000:0000:0000:0000",
		"a:b:c:d::",
		"a:b:c:d::",
		"a:b:c:d::",
		"a:b:c:d::",
		"a:b:c:d::0.0.0.0",
		"a:b:c:d::",
		"a:b:c:d::",
		"a:b:c:d::",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-0-0-0-0.ipv6-literal.net",
		"00|N0s0$ND2BxK96%Chk",
		"0x000a000b000c000d0000000000000000",
		"00000240001300006000032000000000000000000000")

	t.testIPv6Strings("a:b:c:d::/64",
		"a:b:c:d:0:0:0:0/64",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:%:%:%:%",
		"000a:000b:000c:000d:0000:0000:0000:0000/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d::0.0.0.0/64",
		"a:b:c:d::0.0.0.0/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-0-0-0-0.ipv6-literal.net/64",
		"00|N0s0$ND2BxK96%Chk/64",
		"0x000a000b000c000d0000000000000000-0x000a000b000c000dffffffffffffffff",
		"00000240001300006000032000000000000000000000-00000240001300006000033777777777777777777777")

	t.testIPv6Strings("a::d:*:*:*:*/65",
		"a:0:0:d:0-8000:0:0:0/65",
		"a:0:0:d:*:*:*:*",
		"a::d:*:*:*:*",
		"a:0:0:d:%:%:%:%",
		"000a:0000:0000:000d:0000-8000:0000:0000:0000/65",
		"a:0:0:d:0-8000::/65",
		"a:0:0:d:0-8000::/65",
		"a:0:0:d:0-8000::/65",
		"a::d:*:*:*:*",
		"a::d:0-8000:0:0.0.0.0/65",
		"a::d:0-8000:0:0.0.0.0/65",
		"a:0:0:d:0-8000::/65",
		"a:0:0:d:0-8000::/65",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-0"+ipaddr.IPv6AlternativeRangeSeparatorStr+"8000-0-0-0.ipv6-literal.net/65",
		"00|M>t|tt+WbKhfd5~qN"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|M>t|tt+;M72aZe}L&/65",
		"0x000a00000000000d0000000000000000-0x000a00000000000dffffffffffffffff",
		"00000240000000000000032000000000000000000000-00000240000000000000033777777777777777777777")

	t.testIPv6Strings("a::d:0-7fff:*:*:*/65",
		"a:0:0:d:0:0:0:0/65",
		"a:0:0:d:0-7fff:*:*:*",
		"a::d:0-7fff:*:*:*",
		"a:0:0:d:0-7fff:%:%:%",
		"000a:0000:0000:000d:0000:0000:0000:0000/65",
		"a:0:0:d::/65",
		"a:0:0:d::/65",
		"a:0:0:d::/65",
		"a::d:0-7fff:*:*:*",
		"a::d:0:0:0.0.0.0/65",
		"a::d:0:0:0.0.0.0/65",
		"a:0:0:d::/65",
		"a:0:0:d::/65",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.0-7.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-0-0-0-0.ipv6-literal.net/65",
		"00|M>t|tt+WbKhfd5~qN/65",
		"0x000a00000000000d0000000000000000-0x000a00000000000d7fffffffffffffff",
		"00000240000000000000032000000000000000000000-00000240000000000000032777777777777777777777")

	t.testIPv6Strings("a::d:0:0:0:0/65",
		"a:0:0:d:0:0:0:0/65",
		"a:0:0:d:0-7fff:*:*:*",
		"a::d:0-7fff:*:*:*",
		"a:0:0:d:0-7fff:%:%:%",
		"000a:0000:0000:000d:0000:0000:0000:0000/65",
		"a:0:0:d::/65",
		"a:0:0:d::/65",
		"a:0:0:d::/65",
		"a::d:0-7fff:*:*:*",
		"a::d:0:0:0.0.0.0/65",
		"a::d:0:0:0.0.0.0/65",
		"a:0:0:d::/65",
		"a:0:0:d::/65",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.0-7.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-0-0-0-0.ipv6-literal.net/65",
		"00|M>t|tt+WbKhfd5~qN/65",
		"0x000a00000000000d0000000000000000-0x000a00000000000d7fffffffffffffff",
		"00000240000000000000032000000000000000000000-00000240000000000000032777777777777777777777")

	t.testIPv6Strings("a::d:*:*:*:0/65",
		"a:0:0:d:*:*:*:0/65",
		"a:0:0:d:*:*:*:0",
		"a::d:*:*:*:0",
		"a:0:0:d:%:%:%:0",
		"000a:0000:0000:000d:0000-ffff:0000-ffff:0000-ffff:0000/65",
		"a::d:*:*:*:0/65",
		"a::d:*:*:*:0/65",
		"a:0:0:d:*:*:*::/65",
		"a::d:*:*:*:0",
		"a::d:*:*:*.*.0.0/65",
		"a::d:*:*:*.*.0.0/65",
		"a::d:*:*:*.*.0.0/65",
		"a::d:*:*:*.*.0.0/65",
		"0.0.0.0.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-*-*-*-0.ipv6-literal.net/65",
		"",
		"",
		"")

	t.testIPv6Strings("a::d:0:*:0:*/65",
		"a:0:0:d:0:*:0:*/65",
		"a:0:0:d:0:*:0:*",
		"a::d:0:*:0:*",
		"a:0:0:d:0:%:0:%",
		"000a:0000:0000:000d:0000:0000-ffff:0000:0000-ffff/65",
		"a::d:0:*:0:*/65",
		"a::d:0:*:0:*/65",
		"a:0:0:d:0:*::*/65",
		"a::d:0:*:0:*",
		"a::d:0:*:0.0.*.*/65",
		"a::d:0:*:0.0.*.*/65",
		"a::d:0:*:0.0.*.*/65",
		"a::d:0:*:0.0.*.*/65",
		"*.*.*.*.0.0.0.0.*.*.*.*.0.0.0.0.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-0-*-0-*.ipv6-literal.net/65",
		"",
		"",
		"")

	t.testIPv6Strings("a::d:*:0:0:0/65",
		"a:0:0:d:*:0:0:0/65",
		"a:0:0:d:*:0:0:0",
		"a:0:0:d:*::",
		"a:0:0:d:%:0:0:0",
		"000a:0000:0000:000d:0000-ffff:0000:0000:0000/65",
		"a:0:0:d:*::/65",
		"a:0:0:d:*::/65",
		"a:0:0:d:*::/65",
		"a:0:0:d:*::",
		"a::d:*:0:0.0.0.0/65",
		"a::d:*:0:0.0.0.0/65",
		"a:0:0:d:*::/65",
		"a:0:0:d:*::/65",
		"0.0.0.0.0.0.0.0.0.0.0.0.*.*.*.*.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-*-0-0-0.ipv6-literal.net/65",
		"",
		"",
		"")

	t.testIPv6Strings("a:b:c:d:*::/64",
		"a:b:c:d:*:0:0:0/64",
		"a:b:c:d:*:0:0:0",
		"a:b:c:d:*::",
		"a:b:c:d:%:0:0:0",
		"000a:000b:000c:000d:0000-ffff:0000:0000:0000/64",
		"a:b:c:d:*::/64",
		"a:b:c:d:*::/64",
		"a:b:c:d:*::/64",
		"a:b:c:d:*::",
		"a:b:c:d:*::0.0.0.0/64",
		"a:b:c:d:*::0.0.0.0/64",
		"a:b:c:d:*::/64",
		"a:b:c:d:*::/64",
		"0.0.0.0.0.0.0.0.0.0.0.0.*.*.*.*.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-*-0-0-0.ipv6-literal.net/64",
		"",
		"",
		"")

	t.testIPv6Strings("a:0:c:d:e:f:0:0/97",
		"a:0:c:d:e:f:0:0/97",
		"a:0:c:d:e:f:0-7fff:*",
		"a:0:c:d:e:f:0-7fff:*",
		"a:0:c:d:e:f:0-7fff:%",
		"000a:0000:000c:000d:000e:000f:0000:0000/97",
		"a:0:c:d:e:f::/97",
		"a:0:c:d:e:f::/97",
		"a:0:c:d:e:f::/97",
		"a::c:d:e:f:0-7fff:*",
		"a::c:d:e:f:0.0.0.0/97",
		"a::c:d:e:f:0.0.0.0/97",
		"a::c:d:e:f:0.0.0.0/97",
		"a:0:c:d:e:f::/97",
		"*.*.*.*.*.*.*.0-7.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-e-f-0-0.ipv6-literal.net/97",
		"00|M>t};s?v~hFl`j3_$/97",
		"0x000a0000000c000d000e000f00000000-0x000a0000000c000d000e000f7fffffff",
		"00000240000000006000032000160000740000000000-00000240000000006000032000160000757777777777")

	t.testIPv6Strings("a:0:c:d:e:f:0:0/96",
		"a:0:c:d:e:f:0:0/96",
		"a:0:c:d:e:f:*:*",
		"a:0:c:d:e:f:*:*",
		"a:0:c:d:e:f:%:%",
		"000a:0000:000c:000d:000e:000f:0000:0000/96",
		"a:0:c:d:e:f::/96",
		"a:0:c:d:e:f::/96",
		"a:0:c:d:e:f::/96",
		"a::c:d:e:f:*:*",
		"a::c:d:e:f:0.0.0.0/96",
		"a::c:d:e:f:0.0.0.0/96",
		"a:0:c:d:e:f::/96",
		"a:0:c:d:e:f::/96",
		"*.*.*.*.*.*.*.*.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-e-f-0-0.ipv6-literal.net/96",
		"00|M>t};s?v~hFl`j3_$/96",
		"0x000a0000000c000d000e000f00000000-0x000a0000000c000d000e000fffffffff",
		"00000240000000006000032000160000740000000000-00000240000000006000032000160000777777777777")

	t.testIPv6Strings("a:0:c:d:e:f:1:0/112",
		"a:0:c:d:e:f:1:0/112",
		"a:0:c:d:e:f:1:*",
		"a:0:c:d:e:f:1:*",
		"a:0:c:d:e:f:1:%",
		"000a:0000:000c:000d:000e:000f:0001:0000/112",
		"a::c:d:e:f:1:0/112",     //compressed
		"a:0:c:d:e:f:1:0/112",    //canonical (only zeros are single so not compressed)
		"a:0:c:d:e:f:1::/112",    //subnet
		"a::c:d:e:f:1:*",         //compressed wildcard
		"a::c:d:e:f:0.1.0.0/112", //mixed, no compress
		"a::c:d:e:f:0.1.0.0/112", //mixed, no compress host
		"a::c:d:e:f:0.1.0.0/112",
		"a::c:d:e:f:0.1.0.0/112",
		"*.*.*.*.1.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-e-f-1-0.ipv6-literal.net/112",
		"00|M>t};s?v~hFl`jD0%/112",
		"0x000a0000000c000d000e000f00010000-0x000a0000000c000d000e000f0001ffff",
		"00000240000000006000032000160000740000200000-00000240000000006000032000160000740000377777") //mixed

	t.testIPv6Strings("a:0:c:d:0:0:1:0/112",
		"a:0:c:d:0:0:1:0/112", //normalized
		"a:0:c:d:0:0:1:*",     //normalized wildcard
		"a:0:c:d::1:*",        //canonical wildcard
		"a:0:c:d:0:0:1:%",     //sql
		"000a:0000:000c:000d:0000:0000:0001:0000/112", //full
		"a:0:c:d::1:0/112",                            //compressed
		"a:0:c:d::1:0/112",                            //canonical
		"a:0:c:d:0:0:1::/112",                         //subnet
		"a:0:c:d::1:*",                                //compressed wildcard
		"a:0:c:d::0.1.0.0/112",                        //mixed, no compress
		"a:0:c:d::0.1.0.0/112",                        //mixed, no compress host
		"a:0:c:d::0.1.0.0/112",
		"a:0:c:d::0.1.0.0/112",
		"*.*.*.*.1.0.0.0.0.0.0.0.0.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-0-0-1-0.ipv6-literal.net/112",
		"00|M>t};s?v}5L>MDR^a/112",
		"0x000a0000000c000d0000000000010000-0x000a0000000c000d000000000001ffff",
		"00000240000000006000032000000000000000200000-00000240000000006000032000000000000000377777") //mixed

	t.testIPv6Strings("a:b:c:*::/64",
		"a:b:c:*:0:0:0:0/64",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*:*",
		"a:b:c:%:%:%:%:%",
		"000a:000b:000c:0000-ffff:0000:0000:0000:0000/64",
		"a:b:c:*::/64",
		"a:b:c:*::/64",
		"a:b:c:*::/64",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*::0.0.0.0/64",
		"a:b:c:*::0.0.0.0/64",
		"a:b:c:*::/64",
		"a:b:c:*::/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-*-0-0-0-0.ipv6-literal.net/64",
		"00|N0s0$N0-%*(tF5l-X"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|N0s0;%Z;E{Rk+ZU@X/64",
		"0x000a000b000c00000000000000000000-0x000a000b000cffffffffffffffffffff",
		"00000240001300006000000000000000000000000000-00000240001300006377777777777777777777777777")

	t.testIPv6Strings("a::/64",
		"a:0:0:0:0:0:0:0/64",
		"a:0:0:0:*:*:*:*",
		"a::*:*:*:*",
		"a:0:0:0:%:%:%:%",
		"000a:0000:0000:0000:0000:0000:0000:0000/64",
		"a::/64",
		"a::/64",
		"a::/64",
		"a::*:*:*:*",
		"a::0.0.0.0/64",
		"a::0.0.0.0/64",
		"a::/64",
		"a::/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.0.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-0-0-0-0-0.ipv6-literal.net/64",
		"00|M>t|ttwH6V62lVY`A/64",
		"0x000a0000000000000000000000000000-0x000a000000000000ffffffffffffffff",
		"00000240000000000000000000000000000000000000-00000240000000000000001777777777777777777777")

	t.testIPv6Strings("a:0:0:d:e:f:0:0/112",
		"a:0:0:d:e:f:0:0/112",
		"a:0:0:d:e:f:0:*",
		"a::d:e:f:0:*",
		"a:0:0:d:e:f:0:%",
		"000a:0000:0000:000d:000e:000f:0000:0000/112",
		"a::d:e:f:0:0/112",
		"a::d:e:f:0:0/112",
		"a:0:0:d:e:f::/112",
		"a::d:e:f:0:*",
		"a::d:e:f:0.0.0.0/112",
		"a::d:e:f:0.0.0.0/112",
		"a::d:e:f:0.0.0.0/112",
		"a:0:0:d:e:f::/112",
		"*.*.*.*.0.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-e-f-0-0.ipv6-literal.net/112",
		"00|M>t|tt+WcwbECb*xq/112",
		"0x000a00000000000d000e000f00000000-0x000a00000000000d000e000f0000ffff",
		"00000240000000000000032000160000740000000000-00000240000000000000032000160000740000177777")

	t.testIPv6Strings("a:0:c:d:e:f:0:0/112",
		"a:0:c:d:e:f:0:0/112",
		"a:0:c:d:e:f:0:*",
		"a:0:c:d:e:f:0:*",
		"a:0:c:d:e:f:0:%",
		"000a:0000:000c:000d:000e:000f:0000:0000/112",
		"a:0:c:d:e:f::/112",
		"a:0:c:d:e:f::/112",
		"a:0:c:d:e:f::/112",
		"a::c:d:e:f:0:*",
		"a::c:d:e:f:0.0.0.0/112",
		"a::c:d:e:f:0.0.0.0/112",
		"a::c:d:e:f:0.0.0.0/112",
		"a:0:c:d:e:f::/112",
		"*.*.*.*.0.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-e-f-0-0.ipv6-literal.net/112",
		"00|M>t};s?v~hFl`j3_$/112",
		"0x000a0000000c000d000e000f00000000-0x000a0000000c000d000e000f0000ffff",
		"00000240000000006000032000160000740000000000-00000240000000006000032000160000740000177777")

	t.testIPv6Strings("a:0:c:d:e:f:a:0/112",
		"a:0:c:d:e:f:a:0/112",
		"a:0:c:d:e:f:a:*",
		"a:0:c:d:e:f:a:*",
		"a:0:c:d:e:f:a:%",
		"000a:0000:000c:000d:000e:000f:000a:0000/112",
		"a::c:d:e:f:a:0/112",
		"a:0:c:d:e:f:a:0/112",
		"a:0:c:d:e:f:a::/112",
		"a::c:d:e:f:a:*",
		"a::c:d:e:f:0.10.0.0/112",
		"a::c:d:e:f:0.10.0.0/112",
		"a::c:d:e:f:0.10.0.0/112",
		"a::c:d:e:f:0.10.0.0/112",
		"*.*.*.*.a.0.0.0.f.0.0.0.e.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-e-f-a-0.ipv6-literal.net/112",
		"00|M>t};s?v~hFl`k9s=/112",
		"0x000a0000000c000d000e000f000a0000-0x000a0000000c000d000e000f000affff",
		"00000240000000006000032000160000740002400000-00000240000000006000032000160000740002577777")

	t.testIPv6Strings("a:0:c:d:0:0:0:100/120",
		"a:0:c:d:0:0:0:100/120",                       //normalized
		"a:0:c:d:0:0:0:100-1ff",                       //normalized wildcard
		"a:0:c:d::100-1ff",                            //canonical wildcard
		"a:0:c:d:0:0:0:1__",                           //sql
		"000a:0000:000c:000d:0000:0000:0000:0100/120", //full
		"a:0:c:d::100/120",                            //compressed
		"a:0:c:d::100/120",                            //canonical
		"a:0:c:d::100/120",                            //subnet
		"a:0:c:d::100-1ff",                            //compressed wildcard
		"a:0:c:d::0.0.1.0/120",                        //mixed, no compress
		"a:0:c:d::0.0.1.0/120",                        //mixed, no compress host
		"a:0:c:d::0.0.1.0/120",
		"a:0:c:d::0.0.1.0/120",
		"*.*.1.0.0.0.0.0.0.0.0.0.0.0.0.0.d.0.0.0.c.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-c-d-0-0-0-100.ipv6-literal.net/120",
		"00|M>t};s?v}5L>MDI>a/120",
		"0x000a0000000c000d0000000000000100-0x000a0000000c000d00000000000001ff",
		"00000240000000006000032000000000000000000400-00000240000000006000032000000000000000000777") //mixed

	t.testIPv6Strings("a:b:c:d:*/64",
		"a:b:c:d:0:0:0:0/64",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:%:%:%:%",
		"000a:000b:000c:000d:0000:0000:0000:0000/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d::0.0.0.0/64",
		"a:b:c:d::0.0.0.0/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-0-0-0-0.ipv6-literal.net/64",
		"00|N0s0$ND2BxK96%Chk/64",
		"0x000a000b000c000d0000000000000000-0x000a000b000c000dffffffffffffffff",
		"00000240001300006000032000000000000000000000-00000240001300006000033777777777777777777777")

	t.testIPv6Strings("a:b:c:d:*:*:*:*/64",
		"a:b:c:d:0:0:0:0/64",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:%:%:%:%",
		"000a:000b:000c:000d:0000:0000:0000:0000/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d::0.0.0.0/64",
		"a:b:c:d::0.0.0.0/64",
		"a:b:c:d::/64",
		"a:b:c:d::/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-0-0-0-0.ipv6-literal.net/64",
		"00|N0s0$ND2BxK96%Chk/64",
		"0x000a000b000c000d0000000000000000-0x000a000b000c000dffffffffffffffff",
		"00000240001300006000032000000000000000000000-00000240001300006000033777777777777777777777")

	t.testIPv6Strings("a::d:*:*:*:*/64",
		"a:0:0:d:0:0:0:0/64",
		"a:0:0:d:*:*:*:*",
		"a::d:*:*:*:*",
		"a:0:0:d:%:%:%:%",
		"000a:0000:0000:000d:0000:0000:0000:0000/64",
		"a:0:0:d::/64",
		"a:0:0:d::/64",
		"a:0:0:d::/64",
		"a::d:*:*:*:*",
		"a::d:0:0:0.0.0.0/64",
		"a::d:0:0:0.0.0.0/64",
		"a:0:0:d::/64",
		"a:0:0:d::/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-0-0-0-0.ipv6-literal.net/64",
		"00|M>t|tt+WbKhfd5~qN/64",
		"0x000a00000000000d0000000000000000-0x000a00000000000dffffffffffffffff",
		"00000240000000000000032000000000000000000000-00000240000000000000033777777777777777777777")

	t.testIPv6Strings("1::/32",
		"1:0:0:0:0:0:0:0/32",
		"1:0:*:*:*:*:*:*",
		"1:0:*:*:*:*:*:*",
		"1:0:%:%:%:%:%:%",
		"0001:0000:0000:0000:0000:0000:0000:0000/32",
		"1::/32",
		"1::/32",
		"1::/32",
		"1::*:*:*:*:*:*",
		"1::0.0.0.0/32",
		"1::0.0.0.0/32",
		"1::/32",
		"1::/32",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.0.0.0.0.1.0.0.0.ip6.arpa",
		"1-0-0-0-0-0-0-0.ipv6-literal.net/32",
		"008JOm8Mm5*yBppL!sg1/32",
		"0x00010000000000000000000000000000-0x00010000ffffffffffffffffffffffff",
		"00000020000000000000000000000000000000000000-00000020000077777777777777777777777777777777")

	t.testIPv6Strings("ffff::/104",
		"ffff:0:0:0:0:0:0:0/104",
		"ffff:0:0:0:0:0:0-ff:*",
		"ffff::0-ff:*",
		"ffff:0:0:0:0:0:0-ff:%",
		"ffff:0000:0000:0000:0000:0000:0000:0000/104",
		"ffff::/104",
		"ffff::/104",
		"ffff::/104",
		"ffff::0-ff:*",
		"ffff::0.0.0.0/104",
		"ffff::0.0.0.0/104",
		"ffff::0.0.0.0/104",
		"ffff::/104",
		"*.*.*.*.*.*.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-0-0.ipv6-literal.net/104",
		"=q{+M|w0(OeO5^EGP660/104",
		"0xffff0000000000000000000000000000-0xffff0000000000000000000000ffffff",
		"03777760000000000000000000000000000000000000-03777760000000000000000000000000000077777777")

	t.testIPv6Strings("ffff::/108",
		"ffff:0:0:0:0:0:0:0/108",
		"ffff:0:0:0:0:0:0-f:*",
		"ffff::0-f:*",
		"ffff:0:0:0:0:0:_:%",
		"ffff:0000:0000:0000:0000:0000:0000:0000/108",
		"ffff::/108",
		"ffff::/108",
		"ffff::/108",
		"ffff::0-f:*",
		"ffff::0.0.0.0/108",
		"ffff::0.0.0.0/108",
		"ffff::0.0.0.0/108",
		"ffff::/108",
		"*.*.*.*.*.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-0-0.ipv6-literal.net/108",
		"=q{+M|w0(OeO5^EGP660/108",
		"0xffff0000000000000000000000000000-0xffff00000000000000000000000fffff",
		"03777760000000000000000000000000000000000000-03777760000000000000000000000000000003777777")

	t.testIPv6Strings("ffff::1000:0/108",
		"ffff:0:0:0:0:0:1000:0/108",
		"ffff:0:0:0:0:0:1000-100f:*",
		"ffff::1000-100f:*",
		"ffff:0:0:0:0:0:100_:%",
		"ffff:0000:0000:0000:0000:0000:1000:0000/108",
		"ffff::1000:0/108",
		"ffff::1000:0/108",
		"ffff:0:0:0:0:0:1000::/108",
		"ffff::1000-100f:*",
		"ffff::16.0.0.0/108",
		"ffff::16.0.0.0/108",
		"ffff::16.0.0.0/108",
		"ffff::16.0.0.0/108",
		"*.*.*.*.*.0.0.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-1000-0.ipv6-literal.net/108",
		"=q{+M|w0(OeO5^ELbE%G/108",
		"0xffff0000000000000000000010000000-0xffff00000000000000000000100fffff",
		"03777760000000000000000000000000002000000000-03777760000000000000000000000000002003777777")

	t.testIPv6Strings("ffff::a000:0/108",
		"ffff:0:0:0:0:0:a000:0/108",
		"ffff:0:0:0:0:0:a000-a00f:*",
		"ffff::a000-a00f:*",
		"ffff:0:0:0:0:0:a00_:%",
		"ffff:0000:0000:0000:0000:0000:a000:0000/108",
		"ffff::a000:0/108",
		"ffff::a000:0/108",
		"ffff:0:0:0:0:0:a000::/108",
		"ffff::a000-a00f:*",
		"ffff::160.0.0.0/108",
		"ffff::160.0.0.0/108",
		"ffff::160.0.0.0/108",
		"ffff::160.0.0.0/108",
		"*.*.*.*.*.0.0.a.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-a000-0.ipv6-literal.net/108",
		"=q{+M|w0(OeO5^E(z82>/108",
		"0xffff00000000000000000000a0000000-0xffff00000000000000000000a00fffff",
		"03777760000000000000000000000000024000000000-03777760000000000000000000000000024003777777")

	t.testIPv6Strings("ffff::/107",
		"ffff:0:0:0:0:0:0:0/107",
		"ffff:0:0:0:0:0:0-1f:*",
		"ffff::0-1f:*",
		"ffff:0:0:0:0:0:0-1f:%",
		"ffff:0000:0000:0000:0000:0000:0000:0000/107",
		"ffff::/107",
		"ffff::/107",
		"ffff::/107",
		"ffff::0-1f:*",
		"ffff::0.0.0.0/107",
		"ffff::0.0.0.0/107",
		"ffff::0.0.0.0/107",
		"ffff::/107",
		"*.*.*.*.*.0-1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-0-0.ipv6-literal.net/107",
		"=q{+M|w0(OeO5^EGP660/107",
		"0xffff0000000000000000000000000000-0xffff00000000000000000000001fffff",
		"03777760000000000000000000000000000000000000-03777760000000000000000000000000000007777777")

	t.testIPv6Strings("abcd::/107",
		"abcd:0:0:0:0:0:0:0/107",
		"abcd:0:0:0:0:0:0-1f:*",
		"abcd::0-1f:*",
		"abcd:0:0:0:0:0:0-1f:%",
		"abcd:0000:0000:0000:0000:0000:0000:0000/107",
		"abcd::/107",
		"abcd::/107",
		"abcd::/107",
		"abcd::0-1f:*",
		"abcd::0.0.0.0/107",
		"abcd::0.0.0.0/107",
		"abcd::0.0.0.0/107",
		"abcd::/107",
		"*.*.*.*.*.0-1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.d.c.b.a.ip6.arpa",
		"abcd-0-0-0-0-0-0-0.ipv6-literal.net/107",
		"o6)n`s#^$cP5&p^H}p=a/107",
		"0xabcd0000000000000000000000000000-0xabcd00000000000000000000001fffff",
		"02536320000000000000000000000000000000000000-02536320000000000000000000000000000007777777")

	t.testIPv6Strings("1:2:3:4::/80",
		"1:2:3:4:0:0:0:0/80", //normalized
		"1:2:3:4:0:*:*:*",    //normalizedWildcards
		"1:2:3:4:0:*:*:*",    //canonicalWildcards
		"1:2:3:4:0:%:%:%",    //sql
		"0001:0002:0003:0004:0000:0000:0000:0000/80",
		"1:2:3:4::/80", //compressed
		"1:2:3:4::/80",
		"1:2:3:4::/80",
		"1:2:3:4::*:*:*",
		"1:2:3:4::0.0.0.0/80", //mixed no compress
		"1:2:3:4::0.0.0.0/80", //mixedNoCompressHost
		"1:2:3:4::/80",
		"1:2:3:4::/80",
		"*.*.*.*.*.*.*.*.*.*.*.*.0.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-0-0-0.ipv6-literal.net/80",
		"008JQWOV7Skb)C|ve)jA/80",
		"0x00010002000300040000000000000000-0x00010002000300040000ffffffffffff",
		"00000020000200001400010000000000000000000000-00000020000200001400010000007777777777777777")

	t.testIPv6Strings("a:b:c:*:*:*:*:*", //as noted above, addresses are not converted to prefix if starting as wildcards.
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*:*",
		"a:b:c:%:%:%:%:%",
		"000a:000b:000c:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*:*",
		"a:b:c:*:*:*:*.*.*.*",
		"a:b:c:*:*:*:*.*.*.*",
		"a:b:c:*:*:*:*.*.*.*",
		"a:b:c:*:*:*:*.*.*.*",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-*-*-*-*-*.ipv6-literal.net",
		"00|N0s0$N0-%*(tF5l-X"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|N0s0;%a&*sUa#KSGX",
		"0x000a000b000c00000000000000000000-0x000a000b000cffffffffffffffffffff",
		"00000240001300006000000000000000000000000000-00000240001300006377777777777777777777777777")

	t.testIPv6Strings("a:b:c:d:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:%:%:%:%",
		"000a:000b:000c:000d:0000-ffff:0000-ffff:0000-ffff:0000-ffff",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*.*.*.*",
		"a:b:c:d:*:*:*.*.*.*",
		"a:b:c:d:*:*:*.*.*.*",
		"a:b:c:d:*:*:*.*.*.*",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-*-*-*-*.ipv6-literal.net",
		"00|N0s0$ND2BxK96%Chk"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|N0s0$ND{&WM}~o9(k",
		"0x000a000b000c000d0000000000000000-0x000a000b000c000dffffffffffffffff",
		"00000240001300006000032000000000000000000000-00000240001300006000033777777777777777777777")

	t.testIPv6Strings("a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:%:%:%:%",
		"000a:000b:000c:000d:0000-ffff:0000-ffff:0000-ffff:0000-ffff",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*:*",
		"a:b:c:d:*:*:*.*.*.*",
		"a:b:c:d:*:*:*.*.*.*",
		"a:b:c:d:*:*:*.*.*.*",
		"a:b:c:d:*:*:*.*.*.*",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.c.0.0.0.b.0.0.0.a.0.0.0.ip6.arpa",
		"a-b-c-d-*-*-*-*.ipv6-literal.net",
		"00|N0s0$ND2BxK96%Chk"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|N0s0$ND{&WM}~o9(k",
		"0x000a000b000c000d0000000000000000-0x000a000b000c000dffffffffffffffff",
		"00000240001300006000032000000000000000000000-00000240001300006000033777777777777777777777")

	t.testIPv6Strings("a::c:d:*",
		"a:0:0:0:0:c:d:*",
		"a:0:0:0:0:c:d:*",
		"a::c:d:*",
		"a:0:0:0:0:c:d:%",
		"000a:0000:0000:0000:0000:000c:000d:0000-ffff",
		"a::c:d:*",
		"a::c:d:*",
		"a::c:d:*",
		"a::c:d:*",
		"a::c:0.13.*.*",
		"a::c:0.13.*.*",
		"a::c:0.13.*.*",
		"a::c:0.13.*.*",
		"*.*.*.*.d.0.0.0.c.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-0-0-c-d-*.ipv6-literal.net",
		"00|M>t|ttwH6V6EEzblZ"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|M>t|ttwH6V6EEzkrZ",
		"0x000a0000000000000000000c000d0000-0x000a0000000000000000000c000dffff",
		"00000240000000000000000000000000600003200000-00000240000000000000000000000000600003377777")

	t.testIPv6Strings("a::d:*:*:*:*",
		"a:0:0:d:*:*:*:*",
		"a:0:0:d:*:*:*:*",
		"a::d:*:*:*:*",
		"a:0:0:d:%:%:%:%",
		"000a:0000:0000:000d:0000-ffff:0000-ffff:0000-ffff:0000-ffff",
		"a::d:*:*:*:*",
		"a::d:*:*:*:*",
		"a::d:*:*:*:*",
		"a::d:*:*:*:*",
		"a::d:*:*:*.*.*.*",
		"a::d:*:*:*.*.*.*",
		"a::d:*:*:*.*.*.*",
		"a::d:*:*:*.*.*.*",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.d.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-d-*-*-*-*.ipv6-literal.net",
		"00|M>t|tt+WbKhfd5~qN"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|M>t|tt-R6^kVV>{?N",
		"0x000a00000000000d0000000000000000-0x000a00000000000dffffffffffffffff",
		"00000240000000000000032000000000000000000000-00000240000000000000033777777777777777777777")

	t.testIPv6Strings("a::c:d:*/64",
		"a:0:0:0:0:c:d:*/64",
		"a:0:0:0:0:c:d:*",
		"a::c:d:*",
		"a:0:0:0:0:c:d:%",
		"000a:0000:0000:0000:0000:000c:000d:0000-ffff/64",
		"a::c:d:*/64",
		"a::c:d:*/64",
		"a::c:d:*/64",
		"a::c:d:*",
		"a::c:0.13.*.*/64",
		"a::c:0.13.*.*/64",
		"a::c:0.13.*.*/64",
		"a::c:0.13.*.*/64",
		"*.*.*.*.d.0.0.0.c.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-0-0-c-d-*.ipv6-literal.net/64",
		"00|M>t|ttwH6V6EEzblZ"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|M>t|ttwH6V6EEzkrZ/64",
		"0x000a0000000000000000000c000d0000-0x000a0000000000000000000c000dffff",
		"00000240000000000000000000000000600003200000-00000240000000000000000000000000600003377777")

	t.testIPv6Strings("a::c:d:*/80", //similar to above, but allows us to test the base 85 string with non-64 bit prefix
		"a:0:0:0:0:c:d:*/80",
		"a:0:0:0:0:c:d:*",
		"a::c:d:*",
		"a:0:0:0:0:c:d:%",
		"000a:0000:0000:0000:0000:000c:000d:0000-ffff/80",
		"a::c:d:*/80",
		"a::c:d:*/80",
		"a::c:d:*/80",
		"a::c:d:*",
		"a::c:0.13.*.*/80",
		"a::c:0.13.*.*/80",
		"a::c:0.13.*.*/80",
		"a::c:0.13.*.*/80",
		"*.*.*.*.d.0.0.0.c.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-0-0-c-d-*.ipv6-literal.net/80",
		"00|M>t|ttwH6V6EEzblZ"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|M>t|ttwH6V6EEzkrZ/80",
		"0x000a0000000000000000000c000d0000-0x000a0000000000000000000c000dffff",
		"00000240000000000000000000000000600003200000-00000240000000000000000000000000600003377777")

	t.testIPv6Strings("a::c:d:*/48", //similar to above, but allows us to test the base 85 string with non-64 bit prefix
		"a:0:0:0:0:c:d:*/48",
		"a:0:0:0:0:c:d:*",
		"a::c:d:*",
		"a:0:0:0:0:c:d:%",
		"000a:0000:0000:0000:0000:000c:000d:0000-ffff/48",
		"a::c:d:*/48",
		"a::c:d:*/48",
		"a::c:d:*/48",
		"a::c:d:*",
		"a::c:0.13.*.*/48",
		"a::c:0.13.*.*/48",
		"a::c:0.13.*.*/48",
		"a::c:0.13.*.*/48",
		"*.*.*.*.d.0.0.0.c.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.ip6.arpa",
		"a-0-0-0-0-c-d-*.ipv6-literal.net/48",
		"00|M>t|ttwH6V6EEzblZ"+ipaddr.IPv6AlternativeRangeSeparatorStr+"00|M>t|ttwH6V6EEzkrZ/48",
		"0x000a0000000000000000000c000d0000-0x000a0000000000000000000c000dffff",
		"00000240000000000000000000000000600003200000-00000240000000000000000000000000600003377777")

	t.testIPv4Strings("1.2.0.4/16", "1.2.0.4/16", "1.2.0.4", "1.2.0.4", "001.002.000.004/16", "01.02.00.04/16", "0x1.0x2.0x0.0x4/16", "4.0.2.1.in-addr.arpa", "0x01020004", "000100400004")
	t.testIPv4Strings("1.2.3.0/16", "1.2.3.0/16", "1.2.3.0", "1.2.3.0", "001.002.003.000/16", "01.02.03.00/16", "0x1.0x2.0x3.0x0/16", "0.3.2.1.in-addr.arpa", "0x01020300", "000100401400")
	t.testIPv4Strings("1.2.0.0/14", "1.2.0.0/14", "1.2.0.0", "1.2.0.0", "001.002.000.000/14", "01.02.00.00/14", "0x1.0x2.0x0.0x0/14", "0.0.2.1.in-addr.arpa", "0x01020000", "000100400000")

	t.testIPv4Strings("1.2.*.4/16", "1.2.*.4/16", "1.2.*.4", "1.2.%.4", "001.002.000-255.004/16", "01.02.*.04/16", "0x1.0x2.*.0x4/16", "4.*.2.1.in-addr.arpa", "", "")
	t.testIPv4Strings("1.2.3.*/16", "1.2.3.*/16", "1.2.3.*", "1.2.3.%", "001.002.003.000-255/16", "01.02.03.*/16", "0x1.0x2.0x3.*/16", "*.3.2.1.in-addr.arpa", "0x01020300-0x010203ff", "000100401400-000100401777")
	t.testIPv4Strings("1.2.*.*/14", "1.2.*.*/14", "1.2.*.*", "1.2.%.%", "001.002.000-255.000-255/14", "01.02.*.*/14", "0x1.0x2.*.*/14", "*.*.2.1.in-addr.arpa", "0x01020000-0x0102ffff", "000100400000-000100577777") //000100400000-000100400000/14"

	t.testIPv6Strings("ffff::/8",
		"ffff:0:0:0:0:0:0:0/8",
		"ffff:0:0:0:0:0:0:0",
		"ffff::",
		"ffff:0:0:0:0:0:0:0",
		"ffff:0000:0000:0000:0000:0000:0000:0000/8",
		"ffff::/8",
		"ffff::/8",
		"ffff::/8",
		"ffff::",
		"ffff::0.0.0.0/8",
		"ffff::0.0.0.0/8",
		"ffff::/8",
		"ffff::/8",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-0-0.ipv6-literal.net/8",
		"=q{+M|w0(OeO5^EGP660/8",
		"0xffff0000000000000000000000000000",
		"03777760000000000000000000000000000000000000")

	t.testIPv6Strings("ffff::eeee:eeee/108",
		"ffff:0:0:0:0:0:eeee:eeee/108",
		"ffff:0:0:0:0:0:eeee:eeee",
		"ffff::eeee:eeee",
		"ffff:0:0:0:0:0:eeee:eeee",
		"ffff:0000:0000:0000:0000:0000:eeee:eeee/108",
		"ffff::eeee:eeee/108",
		"ffff::eeee:eeee/108",
		"ffff::eeee:eeee/108",
		"ffff::eeee:eeee",
		"ffff::238.238.238.238/108",
		"ffff::238.238.238.238/108",
		"ffff::238.238.238.238/108",
		"ffff::238.238.238.238/108",
		"e.e.e.e.e.e.e.e.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.f.f.f.f.ip6.arpa",
		"ffff-0-0-0-0-0-eeee-eeee.ipv6-literal.net/108",
		"=q{+M|w0(OeO5^F87dpH/108",
		"0xffff00000000000000000000eeeeeeee",
		"03777760000000000000000000000000035673567356")

	t.testIPv6Strings("1:2:3:4::%x%x%", //Note: % is the zone character (not sql wildcard), so this is handled as 1:2:3:4:: with zone x%x%
		"1:2:3:4:0:0:0:0%x%x%", //normalized
		"1:2:3:4:0:0:0:0%x%x%", //normalizedWildcards
		"1:2:3:4::%x%x%",       //canonicalWildcards
		"1:2:3:4:0:0:0:0%x%x%", //sql
		"0001:0002:0003:0004:0000:0000:0000:0000%x%x%",
		"1:2:3:4::%x%x%",        //compressed
		"1:2:3:4::%x%x%",        //canonical
		"1:2:3:4::%x%x%",        //subnet
		"1:2:3:4::%x%x%",        //compressed wildcard
		"1:2:3:4::0.0.0.0%x%x%", //mixed no compress
		"1:2:3:4::%x%x%",        //mixedNoCompressHost
		"1:2:3:4::%x%x%",
		"1:2:3:4::%x%x%",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-0-0-0sxsxs.ipv6-literal.net",
		"008JQWOV7Skb)C|ve)jA"+ipaddr.IPv6AlternativeZoneSeparatorStr+"x%x%",
		"0x00010002000300040000000000000000%x%x%",
		"00000020000200001400010000000000000000000000%x%x%") //mixed

	t.testIPv6Strings("1:2:3:4:5:6:7:8%a/64", //Note: % is the zone character (not sql wildcard), so this is handled as 1:2:3:4:: with zone :%:%
		"1:2:3:4:5:6:7:8%a/64", //normalized
		"1:2:3:4:5:6:7:8%a",    //normalizedWildcards
		"1:2:3:4:5:6:7:8%a",    //canonicalWildcards
		"1:2:3:4:5:6:7:8%a",    //sql
		"0001:0002:0003:0004:0005:0006:0007:0008%a/64",
		"1:2:3:4:5:6:7:8%a/64",     //compressed
		"1:2:3:4:5:6:7:8%a/64",     //canonical
		"1:2:3:4:5:6:7:8%a/64",     //subnet
		"1:2:3:4:5:6:7:8%a",        //compressed wildcard
		"1:2:3:4:5:6:0.7.0.8%a/64", //mixed no compress
		"1:2:3:4:5:6:0.7.0.8%a/64", //mixedNoCompressHost
		"1:2:3:4:5:6:0.7.0.8%a/64",
		"1:2:3:4:5:6:0.7.0.8%a/64",
		"8.0.0.0.7.0.0.0.6.0.0.0.5.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-5-6-7-8sa.ipv6-literal.net/64",
		"008JQWOV7SkcR4tS1R_a"+ipaddr.IPv6AlternativeZoneSeparatorStr+"a/64",
		"0x00010002000300040005000600070008%a",
		"00000020000200001400010000050000300001600010%a")

	t.testIPv6Strings("1:2:3:4::%a/64", //Note: % is the zone character (not sql wildcard), so this is handled as 1:2:3:4:: with zone :%:%
		"1:2:3:4:0:0:0:0%a/64", //normalized
		"1:2:3:4:*:*:*:*%a",    //normalizedWildcards
		"1:2:3:4:*:*:*:*%a",    //canonicalWildcards
		"1:2:3:4:%:%:%:%%a",    //sql
		"0001:0002:0003:0004:0000:0000:0000:0000%a/64",
		"1:2:3:4::%a/64",        //compressed
		"1:2:3:4::%a/64",        //canonical
		"1:2:3:4::%a/64",        //subnet
		"1:2:3:4:*:*:*:*%a",     //compressed wildcard
		"1:2:3:4::0.0.0.0%a/64", //mixed no compress
		"1:2:3:4::0.0.0.0%a/64", //mixedNoCompressHost
		"1:2:3:4::%a/64",
		"1:2:3:4::%a/64",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-0-0-0sa.ipv6-literal.net/64",
		"008JQWOV7Skb)C|ve)jA"+ipaddr.IPv6AlternativeZoneSeparatorStr+"a/64",
		"0x00010002000300040000000000000000-0x0001000200030004ffffffffffffffff%a",
		"00000020000200001400010000000000000000000000-00000020000200001400011777777777777777777777%a")

	t.testIPv6Strings("1:2:3:4::%.a.a", //Note: % is the zone character (not sql wildcard), so this is handled as 1:2:3:4:: with zone .a.a
		"1:2:3:4:0:0:0:0%.a.a", //normalized
		"1:2:3:4:0:0:0:0%.a.a", //normalizedWildcards
		"1:2:3:4::%.a.a",       //canonicalWildcards
		"1:2:3:4:0:0:0:0%.a.a", //sql
		"0001:0002:0003:0004:0000:0000:0000:0000%.a.a",
		"1:2:3:4::%.a.a",        //compressed
		"1:2:3:4::%.a.a",        //canonical
		"1:2:3:4::%.a.a",        //subnet
		"1:2:3:4::%.a.a",        //compressed wildcard
		"1:2:3:4::0.0.0.0%.a.a", //mixed no compress
		"1:2:3:4::%.a.a",        //mixedNoCompressHost
		"1:2:3:4::%.a.a",
		"1:2:3:4::%.a.a",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-0-0-0s.a.a.ipv6-literal.net",
		"008JQWOV7Skb)C|ve)jA"+ipaddr.IPv6AlternativeZoneSeparatorStr+".a.a",
		"0x00010002000300040000000000000000%.a.a",
		"00000020000200001400010000000000000000000000%.a.a") //mixed
	t.testIPv6Strings("1:2:3:4::*:*:*",
		"1:2:3:4:0:*:*:*", //normalized
		"1:2:3:4:0:*:*:*", //normalizedWildcards
		"1:2:3:4:0:*:*:*", //canonicalWildcards
		"1:2:3:4:0:%:%:%", //sql
		"0001:0002:0003:0004:0000:0000-ffff:0000-ffff:0000-ffff",
		"1:2:3:4::*:*:*",     //compressed
		"1:2:3:4:0:*:*:*",    //canonical
		"1:2:3:4::*:*:*",     //subnet
		"1:2:3:4::*:*:*",     //compressed wildcard
		"1:2:3:4::*:*.*.*.*", //mixed no compress
		"1:2:3:4::*:*.*.*.*", //mixedNoCompressHost
		"1:2:3:4::*:*.*.*.*",
		"1:2:3:4::*:*.*.*.*",
		"*.*.*.*.*.*.*.*.*.*.*.*.0.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-*-*-*.ipv6-literal.net",
		"008JQWOV7Skb)C|ve)jA"+ipaddr.IPv6AlternativeRangeSeparatorStr+"008JQWOV7Skb?_P3;X#A",
		"0x00010002000300040000000000000000-0x00010002000300040000ffffffffffff",
		"00000020000200001400010000000000000000000000-00000020000200001400010000007777777777777777")

	t.testIPv6Strings("1:2:3:4::",
		"1:2:3:4:0:0:0:0", //normalized
		"1:2:3:4:0:0:0:0", //normalizedWildcards
		"1:2:3:4::",       //canonicalWildcards
		"1:2:3:4:0:0:0:0", //sql
		"0001:0002:0003:0004:0000:0000:0000:0000",
		"1:2:3:4::", //compressed
		"1:2:3:4::",
		"1:2:3:4::",
		"1:2:3:4::",
		"1:2:3:4::0.0.0.0", //mixed no compress
		"1:2:3:4::",        //mixedNoCompressHost
		"1:2:3:4::",
		"1:2:3:4::",
		"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-0-0-0.ipv6-literal.net",
		"008JQWOV7Skb)C|ve)jA",
		"0x00010002000300040000000000000000",
		"00000020000200001400010000000000000000000000") //mixed

	t.testIPv6Strings("1:2:3:4:0:6::",
		"1:2:3:4:0:6:0:0", //normalized
		"1:2:3:4:0:6:0:0", //normalizedWildcards
		"1:2:3:4:0:6::",   //canonicalWildcards
		"1:2:3:4:0:6:0:0", //sql
		"0001:0002:0003:0004:0000:0006:0000:0000",
		"1:2:3:4:0:6::", //compressed
		"1:2:3:4:0:6::",
		"1:2:3:4:0:6::",      //subnet
		"1:2:3:4:0:6::",      //compressedWildcard
		"1:2:3:4::6:0.0.0.0", //mixed no compress
		"1:2:3:4:0:6::",      //mixedNoCompressHost
		"1:2:3:4:0:6::",
		"1:2:3:4:0:6::",
		"0.0.0.0.0.0.0.0.6.0.0.0.0.0.0.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-4-0-6-0-0.ipv6-literal.net",
		"008JQWOV7Skb)D3fCrWG",
		"0x00010002000300040000000600000000",
		"00000020000200001400010000000000300000000000")
	t.testIPv6Strings("1:2:3:0:0:6::",
		"1:2:3:0:0:6:0:0", //normalized
		"1:2:3:0:0:6:0:0", //normalizedWildcards
		"1:2:3::6:0:0",    //canonicalWildcards
		"1:2:3:0:0:6:0:0", //sql
		"0001:0002:0003:0000:0000:0006:0000:0000",
		"1:2:3::6:0:0", //compressed
		"1:2:3::6:0:0",
		"1:2:3::6:0:0",     //subnet
		"1:2:3::6:0:0",     //compressedWildcard
		"1:2:3::6:0.0.0.0", //mixed no compress
		"1:2:3::6:0.0.0.0", //mixedNoCompressHost
		"1:2:3::6:0.0.0.0",
		"1:2:3:0:0:6::",
		"0.0.0.0.0.0.0.0.6.0.0.0.0.0.0.0.0.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.ip6.arpa",
		"1-2-3-0-0-6-0-0.ipv6-literal.net",
		"008JQWOV7O(=61h*;$LC",
		"0x00010002000300000000000600000000",
		"00000020000200001400000000000000300000000000")
}

/*
Integer prefix = address.getNetworkPrefixLength();
		if(!nextSegment && prefix != null && prefix == 0 && address.isMultiple() && address.isPrefixBlock()) {
			return new IPAddressString(IPAddress.SEGMENT_WILDCARD_STR, validationOptions);
		}


zeroed is true
*/

//xxxx ok, so I find the code I have in java is jsut too convoluted for ajust by segment xxxxx
//too many  corner cases
//but then that leaves the question, how to express the transition from 0 prefix to all address
//It sort of hinges on the notion of adjusting the prefix as enlarging or reducing a subnet
//BUT that is a bit of a stretch, which I have mover away from
//So, you could perhaps use "enlarge" or "shrink" subnet, although shrink we already have as setPrefixLenZeroed
//Maybe increasePrefixBlockSize?  Or nothing, since it is just a AdjustPrefix along with a ToPrefixBlock?
//Can you generalize increasePrefixBlockSize to go to '*'?  Not really.  IncreaseSubnetSize?
//That still does not make sense since it involves expanding to IPv6.
//	Just do your own.

func enlargeSubnetStr(str *ipaddr.IPAddressString /*boolean nextSegment  false , int bitsPerSegment, /* boolean skipBitCountPrefix false */) *ipaddr.IPAddressString {
	addr := str.GetAddress()
	if addr == nil {
		return nil
	}
	res := enlargeSubnet(addr)
	if res == addr {
		if !res.IsPrefixBlock() {
			return nil
		}
		return ipaddr.NewIPAddressString(ipaddr.SegmentWildcardStr)
	}
	return res.ToAddressString()
	//prefix := str.GetNetworkPrefixLen()
	//	addr := str.GetAddress()
	//	if(prefix == nil) {
	//		return addr.SetPrefixLen(addr.GetBitCount()).ToAddressString()
	//	}
	//	prefLen := *prefix
	//	if prefLen == 0 {
	//		return ipaddr.NewIPAddressString(ipaddr.SegmentWildcardStr)
	//	}
	//	adjustment := ((prefLen - 1) % addr.GetBitsPerSegment()) + 1
	//	return addr.SetPrefixLen(prefLen + adjustment).ToPrefixBlock().ToAddressString()
}

func enlargeSubnet(addr *ipaddr.IPAddress /*boolean nextSegment  false , int bitsPerSegment, /* boolean skipBitCountPrefix false */) *ipaddr.IPAddress {
	prefix := addr.GetNetworkPrefixLen()
	if prefix == nil {
		return addr.SetPrefixLen(addr.GetBitCount())
	}
	prefLen := *prefix
	if prefLen == 0 {
		return addr
	}
	adjustment := ((prefLen - 1) % addr.GetBitsPerSegment()) + 1
	addr, _ = addr.SetPrefixLenZeroed(prefLen - adjustment)
	if addr.GetLower().IsZeroHost() {
		//if addr.IsZeroHost() {xxxx
		addr = addr.ToPrefixBlock()
	}
	return addr
}

func getLabel(addressString *ipaddr.IPAddressString) string {
	address := addressString.GetAddress()
	if address == nil {
		return addressString.String()
	}
	if !address.IsMultiple() {
		return address.ToPrefixLenString()
	}
	return address.ToSubnetString()
}
