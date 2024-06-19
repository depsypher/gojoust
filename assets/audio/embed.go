package audio

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Sound int

const (
	BumpSound       Sound = 0
	EggSound        Sound = 1
	EnergizeSound   Sound = 2
	FlapDnSound     Sound = 3
	FlapUpSound     Sound = 4
	HitSound        Sound = 5
	JoustFreSound   Sound = 6
	JoustLavSound   Sound = 7
	PteroSound      Sound = 8
	SkidSound       Sound = 9
	SpawnSound      Sound = 10
	SpawnEnemySound Sound = 11
	Walk1Sound      Sound = 12
	Walk2Sound      Sound = 13
	WhompSound      Sound = 14
)

var (
	//go:embed bump.ogg
	Bump []byte

	//go:embed egg.ogg
	Egg []byte

	//go:embed energize.ogg
	Energize []byte

	//go:embed flap-dn.ogg
	FlapDn []byte

	//go:embed flap-up.ogg
	FlapUp []byte

	//go:embed hit.ogg
	Hit []byte

	//go:embed joustfre.ogg
	JoustFre []byte

	//go:embed joustlav.ogg
	JoustLav []byte

	//go:embed ptero.ogg
	Ptero []byte

	//go:embed skid.ogg
	Skid []byte

	//go:embed spawn.ogg
	Spawn []byte

	//go:embed spawn-enemy.ogg
	SpawnEnemy []byte

	//go:embed walk1.ogg
	Walk1 []byte

	//go:embed walk2.ogg
	Walk2 []byte

	//go:embed whomp.ogg
	Whomp []byte

	soundFiles = map[Sound][]byte{
		BumpSound:     Bump,
		EggSound:      Egg,
		EnergizeSound: Energize,
		FlapDnSound:   FlapDn,
		FlapUpSound:   FlapUp,
		Walk1Sound:    Walk1,
		Walk2Sound:    Walk2,
		HitSound:      Hit,
	}
	audioContext = audio.NewContext(44100)
)

type SoundPlayer struct {
	player *audio.Player
}

type GameSounds map[Sound]*SoundPlayer

func LoadSounds() (GameSounds, error) {
	sounds := map[Sound]*SoundPlayer{}
	for name, file := range soundFiles {
		reader := bytes.NewReader(file)
		decoded, err := vorbis.DecodeWithSampleRate(44100, reader)
		if err != nil {
			return nil, err
		}
		player, err := audioContext.NewPlayer(decoded)
		if err != nil {
			return nil, err
		}
		sounds[name] = &SoundPlayer{player: player}
	}
	return sounds, nil
}

func (gs GameSounds) StopSounds() {
	for _, s := range gs {
		s.Stop()
	}
}

func (s *SoundPlayer) Play(soundOn bool) error {
	if soundOn && !s.player.IsPlaying() {
		err := s.player.Rewind()
		if err != nil {
			return err
		}
		s.player.Play()
	}
	return nil
}

func (s *SoundPlayer) Stop() {
	if s.player.IsPlaying() {
		s.player.Pause()
	}
}
