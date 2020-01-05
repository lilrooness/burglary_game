package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var id_counter = 0

func start_loop(screen tcell.Screen, game *Game) {
	quit := make(chan struct{})

	go func() {
		for {
			ev := screen.PollEvent()
			switch evType := ev.(type) {
			case *tcell.EventKey:
				if evType.Key() == tcell.KeyEscape || evType.Key() == tcell.KeyEnter {
					close(quit)
					return
				}
			}
		}
	}()

	stop := false
	game_time := 0
	for !stop {
		select {
		case <-quit:
			stop = true
		case <-time.After(time.Millisecond * 50):
			updateGame(game, game_time)
			render(screen, game)
		}
		game_time++
	}

	screen.Fini()
}

func get_next_uuid() int {
	this_id := id_counter
	id_counter = id_counter + 1
	return this_id
}

func updateGame(game *Game, time int) {
	log.WithFields(logrus.Fields{
		"updating updatables": len(game.updatables),
		"fields":              game.updatables,
	}).Info("Updating game")
	game.update(time)
}

func render(screen tcell.Screen, game *Game) {
	style := tcell.StyleDefault

	screen.Clear()

	style = style.Background(tcell.NewHexColor(int32(0x000000)))
	for _, room := range game.rooms {
		for roomX := room.x; roomX < room.w; roomX++ {
			for roomY := room.y; roomY < room.h; roomY++ {
				screen.SetContent(roomX, roomY, ' ', []rune{}, style)
			}
		}
	}

	glyph := '#'
	style = style.Background(tcell.NewHexColor(int32(0xffffff)))
	var xpos, ypos int
	for _, entity := range game.updatables {
		switch entity.(type) {
		case *Person:
			glyph = '&'
			person, _ := entity.(*Person)
			xpos = person.x
			ypos = person.y
			style = style.Background(tcell.NewHexColor(int32(0x000000))).Foreground(tcell.NewHexColor(int32(0xff00ff)))
		case *Cat:
			glyph = '$'
			cat, _ := entity.(*Cat)
			xpos = cat.x
			ypos = cat.y

			if len(cat.dirtyWith) > 0 {
				style = style.Background(tcell.NewHexColor(int32(0x00ffff))).Foreground(tcell.NewHexColor(int32(0xff00ff)))
			} else {
				style = style.Background(tcell.NewHexColor(int32(0x000000))).Foreground(tcell.NewHexColor(int32(0xff00ff)))
			}
		case *CatSick:
			glyph = '*'
			catSick, _ := entity.(*CatSick)
			xpos = catSick.x
			ypos = catSick.y
			style = style.Background(tcell.NewHexColor(int32(0x00ff00))).Foreground(tcell.NewHexColor(int32(0xff00ff)))
		case *SpiltMilk:
			glyph = 'M'
			milk, _ := entity.(*SpiltMilk)
			xpos = milk.x
			ypos = milk.y
			style = style.Background(tcell.NewHexColor(int32(0xffffff))).Foreground(tcell.NewHexColor(int32(0xff00ff)))
		}
		screen.SetContent(xpos, ypos, glyph, []rune{}, style)
	}

	screen.Show()
}

func main() {

	file, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	game := NewGame()

	s, e := tcell.NewScreen()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	start_loop(s, &game)

}
