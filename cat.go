package main

type CatState int

const (
	CAT_NORMAL   CatState = 0
	CAT_DIRTY    CatState = 1
	CAT_BATHING  CatState = 2
	CAT_VOMITING CatState = 3
)

type Cat struct {
	entity
	state          CatState
	stimuli        []Stimulus
	currentTarget  Coord
	hasTarget      bool
	walkingSpeed   int
	bathingSince   int
	dirtyWith      []string
	needs_to_vomit bool
}

func (cat *Cat) update(time int, game *Game) []entity {

	generatedEntities := []entity{}

	stimIndex := getHighestStimuliIndex(cat.stimuli)
	if stimIndex > -1 {
		cat.hasTarget = true
		cat.currentTarget = Coord{
			x: cat.stimuli[stimIndex].x,
			y: cat.stimuli[stimIndex].y,
		}
	}

	cat.stimuli = []Stimulus{}

	if time-cat.lastUpdated > 10 && cat.hasTarget {
		cat.lastUpdated = time
		cat.MoveTowards(cat.currentTarget)
		if cat.x == cat.currentTarget.x && cat.y == cat.currentTarget.y {
			cat.hasTarget = false
		}
	}

	game.process_solid_collisions(cat)

	if !cat.hasTarget {
		switch cat.state {
		case CAT_DIRTY:
			cat.state = CAT_BATHING
			cat.bathingSince = time
		case CAT_BATHING:
			cat.bath(time)

		case CAT_VOMITING:
			cat.vomit(game)
		}
	}

	return generatedEntities
}

func (cat *Cat) vomit(game *Game) {
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
}

func (cat *Cat) bath(time int) {
	if time-cat.bathingSince >= 20 {
		cat.state = CAT_NORMAL
		for _, dirt := range cat.dirtyWith {
			switch dirt {
			case "milk":
				cat.ingest(&SpiltMilk{})
			}
		}
	}
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
		cat.state = CAT_DIRTY
		cat.dirtyWith = append(cat.dirtyWith, "milk")
	}
}

func (cat *Cat) get_collision_type() CollisionType {
	return cat.collisionType
}
