package main

import (
	"github.com/sirupsen/logrus"
)

type Person struct {
	entity
	state           int
	stimuli         []Stimulus
	currentSate     PersonState
	stimulusRange   int
	currentStimulus Stimulus
	jobQueue        []Job
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

	if len(person.jobQueue) > 0 {
		return person.doing_job
	}

	return person.idle
}

func (person *Person) doing_job(game *Game, time int) PersonState {
	job := person.jobQueue[0]

	jobCoord := job.get_coord()

	if x, y := person.get_xy(); x != jobCoord.x && y != jobCoord.y {
		person.move_to(jobCoord, time)
		return person.doing_job
	}

	if done := job.do(person, game, time); done {
		person.jobQueue[0] = nil // dereference the pointer for garbage collection
		person.jobQueue = person.jobQueue[1:]
		log.WithFields(logrus.Fields{
			"jobQueue": person.jobQueue,
		}).Info("Person has finished the job!")

		return person.idle
	}

	return person.doing_job
}

func (person *Person) moving(game *Game, time int) PersonState {
	stimIndex := getHighestStimuliIndex(person, person.stimulusRange)

	if stimIndex > 0 && person.stimuli[stimIndex].intensity > person.currentStimulus.intensity {
		person.currentStimulus = person.stimuli[stimIndex]
		person.stimuli = []Stimulus{}
	}

	if done := person.move_to(Coord{person.currentStimulus.x, person.currentStimulus.y}, time); done {
		return person.idle
	}

	return person.moving
}

func (person *Person) update(time int, game *Game) []entity {

	game.process_solid_collisions(person)
	newState := person.currentSate(game, time)
	person.currentSate = newState

	return []entity{}
}

// returns true if the destination has been reached
func (person *Person) move_to(destination Coord, time int) bool {
	if time-person.lastUpdated > 10 {
		person.lastUpdated = time
		person.MoveTowards(destination)
		if person.x == person.currentStimulus.x && person.y == person.currentStimulus.y {
			log.Info("Person has reached the destination")
			return true
		}
	}

	return false
}

func (person *Person) get_xy() (x, y int) {
	return person.x, person.y
}

func (person *Person) stimulate(stimulus Stimulus) {
	person.stimuli = append(person.stimuli, stimulus)
}

func (person *Person) trigger_collision(updatable Updatable) {
	log.Info("NOT IMPLEMENTED: trigger_collision")
}

func (person *Person) get_collision_type() CollisionType {
	return person.collisionType
}

func (person *Person) get_stimuli() []Stimulus {
	return person.stimuli
}

func (person *Person) queue_job(job Job) bool {
	switch job.(type) {
	case *CleaningJob:
		person.jobQueue = append(person.jobQueue, job)
		return true
	}
	return false
}

func (person *Person) get_id() int {
	return person.id
}
