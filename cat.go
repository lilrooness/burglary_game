package main

type State int

const (
	CAT_NORMAL   State = 0
	CAT_DIRTY    State = 1
	CAT_BATHING  State = 2
	CAT_VOMITING State = 3
)

type CatState func (*Game, int) CatState

type Cat struct {
	entity
	state         State
	stimuli       []Stimulus
	currentTarget Coord
	walkingSpeed  int
	bathingSince  int
	dirtyWith     []string
	stimulusRange int
	currentState  CatState
	currentStimulus Stimulus
}

func NewCat() *Cat {
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
			id:            2,
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
	if time-cat.bathingSince >= 5 {
		for _, dirt := range cat.dirtyWith {
			switch dirt {
			case "milk":
				cat.ingest(&SpiltMilk{})
			}
		}
		log.Info("FINISHED BATHING")
		return cat.idle
	}

	log.Info("STILL BATHING")
	return cat.bathing
}

func (cat *Cat) idle(_ *Game, time int) CatState {
	stimIndex := getHighestStimuliIndex(cat, cat.stimulusRange)

	if stimIndex > -1 {
		cat.currentStimulus = cat.stimuli[stimIndex]
		cat.stimuli = []Stimulus{}
		return cat.moving
	}

	if cat.state == CAT_DIRTY {
		log.Info("CAT_DIRTY")
		cat.bathingSince = time
		return cat.bathing
	} else if cat.state == CAT_VOMITING {
		log.Info("CAT_DIRTY")
		return cat.vomit
	}

	return cat.idle
}

func (cat *Cat) vomit(game *Game, _ int) CatState {
	cat.state = CAT_NORMAL
	game.updatables = append(game.updatables, &CatSick{
		entity{
			x: cat.x,
			y: cat.y,
		},
	})

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

func (cat *Cat) ingest(ingestable Ingestable) {
	for _, effect := range ingestable.get_ingestion_effects() {
		switch effect {
		case LACTOSE:
			cat.state = CAT_VOMITING
		default:
			cat.state = CAT_VOMITING
		}
	}
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
		cat.state = CAT_DIRTY
		cat.dirtyWith = append(cat.dirtyWith, "milk")
	}
}

func (cat *Cat) get_collision_type() CollisionType {
	return cat.collisionType
}

func (cat *Cat) get_stimuli() []Stimulus {
	return cat.stimuli
}
