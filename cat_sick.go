package main

type CatSick struct {
	entity
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