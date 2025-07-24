package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], []byte(name))
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.xCoord = int32(x)
		person.yCord = int32(y)
		person.zCord = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(mana) << mpOffset
		person.meta = person.meta | mask
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(health) << hpOffset
		person.meta = person.meta | mask
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint8(respect) << 4
		person.respectAndLvl = person.respectAndLvl | mask
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint8(strength) << 4
		person.strengthAndXp = person.strengthAndXp | mask
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.strengthAndXp = person.strengthAndXp | uint8(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectAndLvl = person.respectAndLvl | uint8(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(1 << houseOffset)
		person.meta = person.meta | mask
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(1 << gunOffset)
		person.meta = person.meta | mask
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		mask := uint32(1 << familyOffset)
		person.meta = person.meta | mask
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.meta = person.meta | uint32(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	mpOffset     = 22
	hpOffset     = 12
	houseOffset  = 5
	gunOffset    = 4
	familyOffset = 3
)

type GamePerson struct {
	strengthAndXp uint8
	respectAndLvl uint8
	name          [42]byte
	meta          uint32
	xCoord        int32
	yCord         int32
	zCord         int32
	gold          uint32
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}

	for _, o := range options {
		o(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.xCoord)
}

func (p *GamePerson) Y() int {
	return int(p.yCord)
}

func (p *GamePerson) Z() int {
	return int(p.zCord)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	mask := uint32(1023 << mpOffset)
	val := p.meta & mask
	return int(val >> mpOffset)
}

func (p *GamePerson) Health() int {
	mask := uint32(1023 << hpOffset)
	val := p.meta & mask
	return int(val >> hpOffset)
}

func (p *GamePerson) Respect() int {
	mask := uint8(0b1111 << 4)
	val := p.respectAndLvl & mask
	return int(val >> 4)
}

func (p *GamePerson) Strength() int {
	mask := uint8(0b1111 << 4)
	val := p.strengthAndXp & mask
	return int(val >> 4)
}

func (p *GamePerson) Experience() int {
	mask := uint8(0b00001111)
	val := p.strengthAndXp & mask
	return int(val)
}

func (p *GamePerson) Level() int {
	mask := uint8(0b00001111)
	val := p.respectAndLvl & mask
	return int(val)
}

func (p *GamePerson) HasHouse() bool {
	mask := uint32(1 << houseOffset)
	flag := (p.meta & mask) >> houseOffset
	return flag == 1
}

func (p *GamePerson) HasGun() bool {
	mask := uint32(1 << gunOffset)
	flag := (p.meta & mask) >> gunOffset
	return flag == 1
}

func (p *GamePerson) HasFamily() bool {
	mask := uint32(1 << familyOffset)
	flag := (p.meta & mask) >> familyOffset
	return flag == 1
}

func (p *GamePerson) Type() int {
	mask := uint32(3)
	return int(p.meta & mask)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
