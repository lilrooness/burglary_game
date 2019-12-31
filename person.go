package main

type Person struct {
	entity
	state int
	stimuli []Stimulus
}

func (person *Person) update(_ int, _ *Game) []entity {
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