package types

type Weapon struct {
	Name  string
	TP    string
	Skill string

	RNG   string
	DMG   int
	DLS   int
	Hand1 string
	Hand2 string

	Rarity int
	Price  int
	Curr   string

	Qualities  string
	Additional string

	Source string
	Pic    string
}

type Armor struct {
	Name string
	TP   string

	Phys  int
	Super int

	Rarity int
	Price  int
	Curr   string

	Qualities  string
	Additional string

	Source string
	Pic    string
}
