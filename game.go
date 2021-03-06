package main

import (
	"github.com/sirupsen/logrus"
	"math"
)

type entity struct {
	id            int
	x, y          int
	lastUpdated   int
	collisionType CollisionType
}

type Coord struct {
	x, y int
}

type Job interface {
	do(employable Employable, game *Game, time int) bool
	get_coord() Coord
}

type Updatable interface {
	update(time int, game *Game) []entity
	get_xy() (x, y int)
	get_collision_type() CollisionType
	get_id() int
}

type Cleanable interface {
	Updatable
	clean() bool
}

type Employable interface {
	Updatable
	queue_job(job Job) bool
}

type CollisionType int

const (
	NO_COLLISION CollisionType = 0
	SOLID        CollisionType = 1
	DYNAMIC      CollisionType = 2
)

type DynamicCollidable interface {
	get_id() int
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
	Updatable
	stimulate(stimulus Stimulus)
	get_stimuli() (stimuli []Stimulus)
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

type Burgler struct {
	entity
}

type Room struct {
	w, h, x, y int
}

type Game struct {
	updatables    []Updatable
	rooms         []Room
	deletable_ids []int
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
	for _, updatable := range game.updatables {
		u_xpos, u_ypos := updatable.get_xy()

		if updatable.get_id() != collidable.get_id() && updatable.get_collision_type() == SOLID && u_xpos == xpos && u_ypos == ypos {
			log.WithFields(logrus.Fields{
				"u_xpos": u_xpos,
				"u_ypos": u_ypos,
				"xpos":   xpos,
				"ypos":   ypos,
			}).Info("collision!")
			collidable.trigger_collision(updatable)
		}
	}
}

func (game *Game) update(time int) {

	for _, id := range game.deletable_ids {
		ok, index := game.get_updatable_index_by_id(id)
		if ok {
			log.Info("found a deletable!")
			game.updatables = append(game.updatables[:index], game.updatables[index+1:]...)
		}
	}

	for _, v := range game.updatables {
		if !game.is_updatable_deleted(v.get_id()) {
			v.update(time, game)
		}
	}
}

func (game *Game) is_updatable_deleted(id int) bool {
	if ok, _ := game.get_updatable_index_by_id(id); !ok {
		return true
	}

	for _, v := range game.deletable_ids {
		if v == id {
			return true
		}
	}

	return false
}

func (game *Game) delete_updatable(id int) {
	game.deletable_ids = append(game.deletable_ids, id)
}

func (game *Game) get_updatable_index_by_id(id int) (bool, int) {
	for i, v := range game.updatables {
		if v.get_id() == id {
			return true, i
		}
	}

	return false, -1
}

func (game *Game) get_updatable_by_id(id int) (bool, Updatable) {
	if ok, index := game.get_updatable_index_by_id(id); ok {
		return true, game.updatables[index]
	}

	return false, nil
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
					id:            get_next_uuid(),
					x:             5,
					y:             5,
					collisionType: 1,
				},
			},
			NewCat(),
			NewPerson(),
		},
	}
}

func getHighestStimuliIndex(stimulatable Stimulatable, stimulusRange int) int {
	highestStimulusIndex := -1
	stimuli := stimulatable.get_stimuli()
	if len(stimuli) == 0 {
		return highestStimulusIndex
	}
	x, y := stimulatable.get_xy()
	for i, v := range stimuli {
		d := dist(x, y, v.x, v.y)

		log.WithFields(logrus.Fields{
			"distance": d,
		}).Info("checking distance to stimuli ...")

		if d <= stimulusRange && highestStimulusIndex == -1 {
			highestStimulusIndex = i
		} else {
			if d <= stimulusRange && v.intensity > stimuli[highestStimulusIndex].intensity {
				highestStimulusIndex = i
			}
		}
	}

	return highestStimulusIndex
}

func dist(x1, y1, x2, y2 int) int {
	result := math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
	return int(result)
}
