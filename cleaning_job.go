package main

type CleaningJob struct {
	tl, br       Coord
	cleanable_id int
}

func NewCleaningJob(cleanable Cleanable) *CleaningJob {
	x, y := cleanable.get_xy()
	return &CleaningJob{
		tl:           Coord{x, y},
		br:           Coord{x, y},
		cleanable_id: cleanable.get_id(),
	}
}

func (job *CleaningJob) do(employable Employable, game *Game, time int) bool {

	result := false

	if ok, updatable := game.get_updatable_by_id(job.cleanable_id); ok {
		if cleanable, ok := updatable.(Cleanable); ok {
			result = cleanable.clean()
		} else {
			log.Error("updatable could not be cast to cleanable")
		}

	} else {
		log.Error("updatable not found!!")
	}

	if result {
		for i, v := range game.updatables {
			if job.cleanable_id == v.get_id() {
				game.delete_updatable(i)
				break
			}
		}

	}
	return result
}

func (cleanable *CleaningJob) get_coord() Coord {
	return Coord{
		x: cleanable.tl.x + (cleanable.tl.x-cleanable.br.x)/2,
		y: cleanable.tl.y + (cleanable.tl.y-cleanable.br.y)/2,
	}
}
