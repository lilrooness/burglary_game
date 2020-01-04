package main

type State int

const (
	CAT_NORMAL State = 0
)

type CatState func(*Game, int) CatState

type Cat struct {
	entity
	state           State
	stimuli         []Stimulus
	currentTarget   Coord
	walkingSpeed    int
	bathingSince    int
	dirtyWith       []string
	stimulusRange   int
	currentState    CatState
	currentStimulus Stimulus
}

func NewCat() *Cat {
	id := get_next_uuid()
	cat := &Cat{
		stimulusRange: 15,
		stimuli: []Stimulus{
			Stimulus{
				intensity: STIMULUS_MEDIUM,
				x:         10,
				y:         10,
			},
		},
		entity: entity{
			id:            id,
			x:             1,
			y:             1,
			collisionType: 2,
		},
	}

	cat.currentState = cat.idle

	return cat
}

func (cat *Cat) update(time int, game *Game) []entity {
	game.process_solid_collisions(cat)
	newState := cat.currentState(game, time)
	cat.currentState = newState

	return []entity{}
}

func (cat *Cat) bathing(_ *Game, time int) CatState {
	if time-cat.bathingSince >= 15 {
		var newState CatState
		for _, dirt := range cat.dirtyWith {
			switch dirt {
			case "milk":
				newState = cat.ingest(&SpiltMilk{})
			}
		}

		cat.dirtyWith = []string{}

		if newState == nil {
			return cat.idle
		} else {
			return newState
		}

	}

	return cat.bathing
}

func (cat *Cat) idle(_ *Game, time int) CatState {
	stimIndex := getHighestStimuliIndex(cat, cat.stimulusRange)

	if stimIndex > -1 {
		cat.currentStimulus = cat.stimuli[stimIndex]
		cat.stimuli = []Stimulus{}
		return cat.moving
	}

	if len(cat.dirtyWith) > 0 {
		cat.bathingSince = time
		return cat.bathing
	}

	return cat.idle
}

func (cat *Cat) vomit(game *Game, _ int) CatState {
	cat.state = CAT_NORMAL
	game.updatables = append(game.updatables, NewCatSick(cat.x, cat.y))

	for _, stimulatable := range game.updatables {
		stimulatable, ok := stimulatable.(Stimulatable)
		if ok && stimulatable.get_id() != cat.get_id() {
			stimulatable.stimulate(Stimulus{
				intensity: STIMULUS_HIGH,
				x:         cat.x,
				y:         cat.y,
			})
		}
	}

	cat.stimulate(Stimulus{
		x:         cat.x + 5,
		y:         cat.y + 1,
		intensity: STIMULUS_HIGH,
	})

	return cat.idle
}

func (cat *Cat) moving(game *Game, time int) CatState {
	stimIndex := getHighestStimuliIndex(cat, cat.stimulusRange)

	if stimIndex > 0 && cat.stimuli[stimIndex].intensity > cat.currentStimulus.intensity {
		cat.currentStimulus = cat.stimuli[stimIndex]
		cat.stimuli = []Stimulus{}
	}

	if time-cat.lastUpdated > 10 {
		cat.lastUpdated = time
		cat.MoveTowards(Coord{cat.currentStimulus.x, cat.currentStimulus.y})
		if cat.x == cat.currentStimulus.x && cat.y == cat.currentStimulus.y {
			log.Info("Cat has reached the destination")
			return cat.idle
		}
	}

	return cat.moving
}

func (cat *Cat) ingest(ingestable Ingestable) CatState {
	for _, effect := range ingestable.get_ingestion_effects() {
		switch effect {
		case LACTOSE:
			return cat.vomit
		}
	}

	return nil
}

func (cat *Cat) get_xy() (x, y int) {
	return cat.x, cat.y
}

func (cat *Cat) stimulate(stimulus Stimulus) {
	cat.stimuli = append(cat.stimuli, stimulus)
}

func (cat *Cat) trigger_collision(updatable Updatable) {
	if _, ok := updatable.(*SpiltMilk); ok {
		log.Info("The cat is now dirty")
		cat.dirtyWith = append(cat.dirtyWith, "milk")
	}
}

func (cat *Cat) get_collision_type() CollisionType {
	return cat.collisionType
}

func (cat *Cat) get_stimuli() []Stimulus {
	return cat.stimuli
}

func (cat *Cat) get_id() int {
	return cat.id
}
