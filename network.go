package main

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

type segs struct {
	i [4]uint8
}

type Mask struct {
	*segs
}
type IP struct {
	*segs
	mask *Mask
}
type Net struct {
	*segs
	mask *Mask
}

func newSegsByStr(s string) (*segs, error) {
	segs := new(segs)
	_segs := strings.Split(s, ".")
	if len(_segs) != 4 {
		segs.i[0] = uint8(0)
		segs.i[1] = uint8(0)
		segs.i[2] = uint8(0)
		segs.i[3] = uint8(0)
		return segs, errors.New("")
	}
	for i := 0; i < 4; i++ {
		segUint8, err := strconv.Atoi(_segs[i])
		if err != nil || segUint8 > 255 {
			segs.i[0] = uint8(0)
			segs.i[1] = uint8(0)
			segs.i[2] = uint8(0)
			segs.i[3] = uint8(0)
			return segs, errors.New("")
		}
		segs.i[i] = uint8(segUint8)
	}
	return segs, nil
}
func (s *segs) string() string {
	return strconv.Itoa(int(s.i[0])) + "." + strconv.Itoa(int(s.i[1])) + "." + strconv.Itoa(int(s.i[2])) + "." + strconv.Itoa(int(s.i[3]))
}
func (s *segs) equal(ss *segs) bool {
	for i := 0; i < 4; i++ {
		if s.i[i] != ss.i[i] {
			return false
		}
	}
	return true
}
func (s *segs) not() {
	s.i[0] = ^s.i[0]
	s.i[1] = ^s.i[1]
	s.i[2] = ^s.i[2]
	s.i[3] = ^s.i[3]
}
func (s *segs) uint32() uint32 {
	return uint32(s.i[3]) | (uint32(s.i[2]) << 8) | (uint32(s.i[1]) << 16) | (uint32(s.i[0]) << 24)
}
func (s *segs) prev() *segs {
	ss := new(segs)
	suint32 := s.uint32() - 1
	ss.i[3] = uint8(suint32 & uint32(255))
	ss.i[2] = uint8((suint32 >> 8) & uint32(255))
	ss.i[1] = uint8((suint32 >> 16) & uint32(255))
	ss.i[0] = uint8((suint32 >> 24) & uint32(255))

	return ss
}
func (s *segs) next() *segs {
	ss := new(segs)
	suint32 := s.uint32() + 1
	ss.i[3] = uint8(suint32 & uint32(255))
	ss.i[2] = uint8((suint32 >> 8) & uint32(255))
	ss.i[1] = uint8((suint32 >> 16) & uint32(255))
	ss.i[0] = uint8((suint32 >> 24) & uint32(255))

	return ss
}
func and(a, b *segs) *segs {
	c := new(segs)
	c.i[0] = a.i[0] & b.i[0]
	c.i[1] = a.i[1] & b.i[1]
	c.i[2] = a.i[2] & b.i[2]
	c.i[3] = a.i[3] & b.i[3]
	return c
}
func or(a, b *segs) *segs {
	c := new(segs)
	c.i[0] = a.i[0] | b.i[0]
	c.i[1] = a.i[1] | b.i[1]
	c.i[2] = a.i[2] | b.i[2]
	c.i[3] = a.i[3] | b.i[3]
	return c
}
func not(a *segs) *segs {
	b := new(segs)
	b.i[0] = ^a.i[0]
	b.i[1] = ^a.i[1]
	b.i[2] = ^a.i[2]
	b.i[3] = ^a.i[3]

	return b
}

func maskTest(m Mask) bool {
	for i := 0; i < 33; i++ {
		nm, _ := NewMaskByLen(i)
		if m.segs.equal(nm.segs) {
			return true
		}
	}
	return false
}

func NewMaskByLen(l int) (*Mask, error) {
	if l > -1 && l < 33 {
		full := ^uint32(0)
		maskbit := full << uint(32-l)

		mask := new(Mask)
		mask.segs = new(segs)
		mask.segs.i[3] = uint8(maskbit & uint32(255))
		mask.segs.i[2] = uint8((maskbit >> 8) & uint32(255))
		mask.segs.i[1] = uint8((maskbit >> 16) & uint32(255))
		mask.segs.i[0] = uint8((maskbit >> 24) & uint32(255))
		return mask, nil
	}
	return nil, errors.New("")
}
func NewMaskByStr(s string) (*Mask, error) {
	segs, err := newSegsByStr(s)
	mask := new(Mask)
	mask.segs = segs
	if err != nil {
		return mask, err
	}
	if maskTest(*mask) == false {
		return nil, errors.New("")
	}
	return mask, nil
}
func (m *Mask) Equal(mm *Mask) bool {
	return m.segs.equal(mm.segs)
}
func (m *Mask) Len() int {
	for i := 0; i < 33; i++ {
		nm, _ := NewMaskByLen(i)
		if m.Equal(nm) {
			return i
		}
	}
	return 0
}
func (m *Mask) String() string {
	return m.segs.string()
}
func (m *Mask) StringLen() string {
	for i := 0; i < 33; i++ {
		nm, _ := NewMaskByLen(i)
		if m.Equal(nm) {
			return strconv.Itoa(i)
		}
	}
	return "0"
}

func NewIPByStr(s string) (*IP, error) {
	ip := new(IP)
	if net.ParseIP(s) == nil {
		return nil, errors.New("")
	}
	ip.segs, _ = newSegsByStr(s)
	ip.mask, _ = NewMaskByLen(32)

	return ip, nil
}
func (i *IP) SetMask(m *Mask) {
	i.mask = m
}
func (i *IP) SetMaskByStr(s string) {
	mask, err := NewMaskByStr(s)
	if err == nil {
		i.mask = mask
	}
}
func (i *IP) SetMaskByLen(l int) {
	mask, err := NewMaskByLen(l)
	if err == nil {
		i.mask = mask
	}
}
func (i *IP) GetMask() *Mask {
	return i.mask
}
func (i *IP) String() string {
	return i.segs.string()
}
func (i *IP) NextIP() *IP {
	ip := new(IP)
	ip.segs = i.segs.next()
	ip.mask, _ = NewMaskByLen(32)

	return ip
}

func (i *IP) PrevIP() *IP {
	ip := new(IP)
	ip.segs = i.segs.prev()
	ip.mask, _ = NewMaskByLen(32)

	return ip
}
func (i *IP) NetIP() *IP {
	ip := new(IP)
	ip.SetMaskByLen(32)
	ip.segs = and(i.segs, i.mask.segs)
	return ip
}
func (i *IP) BroadcastIP() *IP {
	ip := new(IP)
	ip.SetMaskByLen(32)
	ip.segs = or(i.NetIP().segs, not(i.mask.segs))
	return ip
}
func (i *IP) Equal(ii *IP) bool {
	if i.segs.equal(ii.segs) && i.mask.segs.equal(ii.mask.segs) {
		return true
	}
	return false
}
