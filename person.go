package main

type Person struct {
	entity
	state           int
	stimuli         []Stimulus
	currentSate     PersonState
	stimulusRange   int
	currentStimulus Stimulus
}

type PersonState func(*Game, int) PersonState

func NewPerson() *Person {
	person := &Person{
		entity: entity{
			id: get_next_uuid(),
			x:  1,
			y:  2,
		},
		stimulusRange: 20,
	}
	person.currentSate = person.idle

	return person
}

func (person *Person) idle(_ *Game, _ int) PersonState {
	stimIndex := getHighestStimuliIndex(person, person.stimulusRange)

	if stimIndex > -1 {
		person.currentStimulus = person.stimuli[stimIndex]
		person.stimuli = []Stimulus{}
		return person.moving
	}

	return person.idle
}

func (person *Person) moving(game *Game, time int) PersonState {
	stimIndex := getHighestStimuliIndex(person, person.stimulusRange)

	if stimIndex > 0 && person.stimuli[stimIndex].intensity > person.currentStimulus.intensity {
		person.currentStimulus = person.stimuli[stimIndex]
		person.stimuli = []Stimulus{}
	}

	if time-person.lastUpdated > 10 {
		person.lastUpdated = time
		person.MoveTowards(Coord{person.currentStimulus.x, person.currentStimulus.y})
		if person.x == person.currentStimulus.x && person.y == person.currentStimulus.y {
			log.Info("Person has reached the destination")
			return person.idle
		}
	}

	return person.moving
}

func (person *Person) update(time int, game *Game) []entity {

	game.process_solid_collisions(person)
	newState := person.currentSate(game, time)
	person.currentSate = newState

	return []entity{}
}

func (person *Person) get_xy() (x, y int) {
	return person.x, person.y
}

func (person *Person) stimulate(stimulus Stimulus) {
	person.stimuli = append(person.stimuli, stimulus)
}

func (person *Person) trigger_collision(_ Updatable) {
	log.Info("NOT IMPLEMENTED: trigger_collision")
}

func (person *Person) get_collision_type() CollisionType {
	return person.collisionType
}

func (person *Person) get_stimuli() []Stimulus {
	return person.stimuli
}

func (person *Person) get_id() int {
	return person.id
}
