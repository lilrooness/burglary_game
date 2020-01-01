package main

type SpiltMilk struct {
	entity
}

func (spiltMilk *SpiltMilk) update(_ int, _ *Game) []entity {
	return []entity{}
}

func (spiltMilk *SpiltMilk) get_xy() (x, y int) {
	return spiltMilk.x, spiltMilk.y
}

func (spiltMilk *SpiltMilk) get_collision_type() CollisionType {
	return spiltMilk.collisionType
}

func (spiltMilk *SpiltMilk) get_ingestion_effects() []IngestionEffect {
	return []IngestionEffect{
		LACTOSE,
	}
}

func (spiltMilk *SpiltMilk) get_id() int {
	return spiltMilk.id
}
