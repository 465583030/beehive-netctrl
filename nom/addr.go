package nom

import (
	"encoding/gob"
	"fmt"
)

// MACAddr represents a MAC address.
type MACAddr [6]byte

var (
	BroadcastMAC           MACAddr   = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	CDPMulticastMAC        MACAddr   = [6]byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCC}
	CiscoSTPMulticastMAC   MACAddr   = [6]byte{0x01, 0x00, 0x0C, 0xCC, 0xCC, 0xCD}
	IEEE802MulticastPrefix MACAddr   = [6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x00}
	IPv4MulticastPrefix    MACAddr   = [6]byte{0x01, 0x00, 0x5E, 0x00, 0x00, 0x00}
	IPv6MulticastPrefix    MACAddr   = [6]byte{0x33, 0x33, 0x00, 0x00, 0x00, 0x00}
	LLDPMulticastMACs      []MACAddr = []MACAddr{
		[6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x0E},
		[6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x0C},
		[6]byte{0x01, 0x80, 0xC2, 0x00, 0x00, 0x00},
	}
)

func (m MACAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		m[0], m[1], m[2], m[3], m[4], m[5])
}

// Key returns an string represtation of the MAC address suitable to store in
// dictionaries. It is more efficient compared to MACAddr.String().
func (m MACAddr) Key() string {
	return string(m[:])
}

// IsBroadcast returns whether the MAC address is a broadcast address.
func (m MACAddr) IsBroadcast() bool {
	return m == BroadcastMAC
}

// IsMulticast returns whether the MAC address is a multicast address.
func (m MACAddr) IsMulticast() bool {
	return m == CDPMulticastMAC || m == CiscoSTPMulticastMAC ||
		m.hasPrefix(IEEE802MulticastPrefix, 3) ||
		m.hasPrefix(IPv4MulticastPrefix, 3) ||
		m.hasPrefix(IPv6MulticastPrefix, 2)
}

// IsLLDP returns whether the mac address is a multicast address used for LLDP.
func (m MACAddr) IsLLDP() bool {
	for _, lm := range LLDPMulticastMACs {
		if m == lm {
			return true
		}
	}
	return false
}

func (m MACAddr) hasPrefix(p MACAddr, l int) bool {
	for i := 0; i < l; i++ {
		if p[i] != m[i] {
			return false
		}
	}
	return true
}

func (m MACAddr) Mask(mask MACAddr) MACAddr {
	masked := m
	for i := range mask {
		masked[i] &= mask[i]
	}
	return masked
}

func (m MACAddr) Less(thatm MACAddr) bool {
	for i := range thatm {
		switch {
		case m[i] < thatm[i]:
			return true
		case m[i] > thatm[i]:
			return false
		}
	}
	return false
}

// MaskedMACAddr is a MAC address that is wildcarded with a mask.
type MaskedMACAddr struct {
	Addr MACAddr // The MAC address.
	Mask MACAddr // The mask of the MAC address.
}

// Match returns whether the masked mac address matches mac.
func (mm MaskedMACAddr) Match(mac MACAddr) bool {
	return mm.Mask.Mask(mm.Addr) == mm.Mask.Mask(mac)
}

// Subsumes returns whether this mask address includes all the addresses matched
// by thatmm.
func (mm MaskedMACAddr) Subsumes(thatmm MaskedMACAddr) bool {
	if thatmm.Mask.Less(mm.Mask) {
		return false
	}
	return mm.Match(thatmm.Addr.Mask(thatmm.Mask))
}

// IPv4Addr represents an IP version 4 address in big endian byte order.
// For example, 127.0.0.1 is represented as [4]byte{127, 0, 0, 1}.
type IPv4Addr [4]byte

// Uint32 converts the IP version 4 address into a 32-bit integer in little
// endian byte order.
func (ip IPv4Addr) Uint32() uint32 {
	return uint32(ip[0]<<24 | ip[1]<<16 | ip[2]<<8 | ip[3])
}

// Mask masked the IP address with mask.
func (ip IPv4Addr) Mask(mask IPv4Addr) IPv4Addr {
	masked := ip
	for i := range masked {
		masked[i] &= mask[i]
	}
	return masked
}

// Less returns whether ip is less than thatip.
func (ip IPv4Addr) Less(thatip IPv4Addr) bool {
	for i := range ip {
		switch {
		case ip[i] < thatip[i]:
			return true
		case ip[i] > thatip[i]:
			return false
		}
	}
	return false
}

// MaskedIPv4Addr represents a masked IP address (ie, an IPv4 prefix)
type MaskedIPv4Addr struct {
	Addr IPv4Addr
	Mask IPv4Addr
}

// Match returns whether the masked IP address matches ip.
func (mi MaskedIPv4Addr) Match(ip IPv4Addr) bool {
	return mi.Addr.Mask(mi.Mask) == ip.Mask(mi.Mask)
}

func (mi MaskedIPv4Addr) Subsumes(thatmi MaskedIPv4Addr) bool {
	if thatmi.Mask.Less(mi.Mask) {
		return false
	}
	return mi.Addr.Mask(mi.Mask) == thatmi.Addr.Mask(mi.Mask)
}

// IPv6Addr represents an IP version 6 address in big-endian byte order.
type IPv6Addr [16]byte

// Mask masked the IP address with mask.
func (ip IPv6Addr) Mask(mask IPv6Addr) IPv6Addr {
	masked := ip
	for i := range masked {
		masked[i] &= mask[i]
	}
	return masked
}

// Less returns whether ip is less than thatip.
func (ip IPv6Addr) Less(thatip IPv6Addr) bool {
	for i := range ip {
		switch {
		case ip[i] < thatip[i]:
			return true
		case ip[i] > thatip[i]:
			return false
		}
	}
	return false
}

// MaskedIPv6Addr represents a masked IPv6 address.
type MaskedIPv6Addr struct {
	Addr IPv6Addr
	Mask IPv6Addr
}

// Match returns whether the masked IP address matches ip.
func (mi MaskedIPv6Addr) Match(ip IPv6Addr) bool {
	return mi.Addr.Mask(mi.Mask) == ip.Mask(mi.Mask)
}

func (mi MaskedIPv6Addr) Subsumes(thatmi MaskedIPv6Addr) bool {
	if thatmi.Mask.Less(mi.Mask) {
		return false
	}
	return mi.Addr.Mask(mi.Mask) == thatmi.Addr.Mask(mi.Mask)
}

func init() {
	gob.Register(IPv4Addr{})
	gob.Register(IPv6Addr{})
	gob.Register(MACAddr{})
	gob.Register(MaskedIPv4Addr{})
	gob.Register(MaskedIPv6Addr{})
	gob.Register(MaskedMACAddr{})
}
