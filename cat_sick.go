package main

type CatSick struct {
	entity
	cleanTurns int
}

func NewCatSick(x, y int) *CatSick {
	id := get_next_uuid()
	return &CatSick{
		entity: entity{
			id: id,
			x:  x,
			y:  y,
			collisionType: SOLID,
		},
		cleanTurns: 10,
	}
}

func (catSick *CatSick) update(_ int, _ *Game) []entity {
	return []entity{}
}

func (catSick *CatSick) get_xy() (x, y int) {
	return catSick.x, catSick.y
}

func (catSick *CatSick) get_collision_type() CollisionType {
	return catSick.collisionType
}

func (catSick *CatSick) get_ingestion_effects() []IngestionEffect {
	return []IngestionEffect{
		HEMETIC,
	}
}

func (catSick *CatSick) get_id() int {
	return catSick.id
}

func (catSick *CatSick) clean() bool {
	catSick.cleanTurns--
	return catSick.cleanTurns == 0
}
