package nightoftheskysunveiling

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.NightOfTheSkysUnveiling, NewSet)
}

type Set struct {
	char   *character.CharWrapper
	core   *core.Core
	stacks int
	Index  int
	Count  int
	buff   []float64
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		Count: count,
	}

	// Apply 2-piece effect
	if s.Count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80 // +80 Elemental Mastery
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("notsu-2pc", -1), // permanent
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	return &s, nil
}

// 4-piece effect
func (s *Set) MoonsignLvl(newMoonsignLvl int) {
	if s.Count < 4 {
		return
	}

	const buffKey = "notsu-4pc"

	// If GleamLv is 0 → reset stacks
	if newMoonsignLvl == 0 {
		s.stacks = 0
	} else {
		// Increase stack if GleamLv increased (max 2)
		if newMoonsignLvl > s.stacks && s.stacks < 2 {
			s.stacks++
		}
	}

	// Apply permanent CR buff (15% per stack)
	s.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(buffKey, -1), // permanent
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			for i := range s.buff {
				s.buff[i] = 0
			}
			s.buff[attributes.CR] = 0.15 * float64(s.stacks)
			return s.buff, true
		},
	})
}
