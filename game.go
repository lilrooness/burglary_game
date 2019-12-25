package main

import (
	"github.com/sirupsen/logrus"
)

type entity struct {
	x, y          int
	lastUpdated   int
	collisionType CollisionType
}

func (e *entity) MoveTowards(coord Coord) {
	dx := coord.x - e.x
	dy := coord.y - e.y

	if dx > 0 {
		dx = 1
	} else if dx < 0 {
		dx = -1
	}

	if dy > 0 {
		dy = 1
	} else if dy < 0 {
		dy = -1
	}

	e.x += dx
	e.y += dy
}

func (game *Game) process_solid_collisions(collidable DynamicCollidable) {
	xpos, ypos := collidable.get_xy()
	for _, updatable := range(game.updatables) {
		u_xpos, u_ypos := updatable.get_xy()
		if updatable.get_collision_type() == SOLID && u_xpos == xpos && u_ypos == ypos {
			log.WithFields(logrus.Fields{
    		"u_xpos": u_xpos,
    		"u_ypos": u_ypos,
    		"xpos": xpos,
    		"ypos": ypos,
  		}).Info("collision!")
			collidable.trigger_collision(updatable)
		}
	}
}

type Coord struct {
	x, y int
}

type Updatable interface {
	update(time int, game *Game) []entity
	get_xy() (x, y int)
	get_collision_type() CollisionType
}

type CollisionType int

const (
	NO_COLLISION CollisionType = 0
	SOLID        CollisionType = 1
	DYNAMIC      CollisionType = 2
)

type DynamicCollidable interface {
	trigger_collision(updatable Updatable)
	get_xy() (x, y int)
}

type StimulusIntensity int

const (
	STIMULUS_LOW    StimulusIntensity = 1
	STIMULUS_MEDIUM StimulusIntensity = 2
	STIMULUS_HIGH   StimulusIntensity = 3
)

type Stimulatable interface {
	stimulate(stimulus Stimulus)
}
type Stimulus struct {
	isScary   bool
	intensity StimulusIntensity
	x, y      int
}

type IngestionEffect int

const (
	LACTOSE IngestionEffect = 0
	ALCOHOL IngestionEffect = 1
	HEMETIC IngestionEffect = 2
)

type Ingestable interface {
	get_ingestion_effects() []IngestionEffect
}

type Person struct {
	entity
}

type Burgler struct {
	entity
}

type Room struct {
	w, h, x, y int
}

type Game struct {
	updatables []Updatable
	rooms      []Room
}

func (game *Game) update(time int) {
	for _, v := range game.updatables {
		v.update(time, game)
	}
}

func NewGame() Game {
	return Game{
		rooms: []Room{
			Room{
				x: 1,
				y: 1,
				w: 20,
				h: 20,
			},
		},
		updatables: []Updatable{
			&SpiltMilk{
				entity: entity{
					x:             5,
					y:             5,
					collisionType: 1,
				},
			},
			&Cat{
				stimuli: []Stimulus{
					Stimulus{
						intensity: STIMULUS_MEDIUM,
						x:         10,
						y:         10,
					},
				},
				entity: entity{
					x:             1,
					y:             1,
					collisionType: 2,
				},
			},
		},
	}
}

func getHighestStimuliIndex(stimuli []Stimulus) int {
	highestStimulusIndex := -1
	for i, v := range stimuli {
		if highestStimulusIndex == -1 {
			highestStimulusIndex = i
		} else {
			if v.intensity > stimuli[highestStimulusIndex].intensity {
				highestStimulusIndex = i
			}
		}
	}

	return highestStimulusIndex
}
