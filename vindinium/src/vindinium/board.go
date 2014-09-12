package vindinium

import (
	"reflect"
	"strconv"
	"fmt"
)

const (
	WALL = iota - 2
	AIR
	TAVERN

	AIR_TILE    = " "
	WALL_TILE   = "#"
	TAVERN_TILE = "["
	MINE_TILE   = "$"
	HERO_TILE   = "@"
)

var (
	AIM = map[Direction]*Position{
		"North": &Position{-1, 0},
		"East":  &Position{0, 1},
		"South": &Position{1, 0},
		"West":  &Position{0, -1},
	}
)

type Board struct {
	Size    int    `json:"size"`
	Tiles   string `json:"tiles"`
	Tileset [][]interface{}
	MinesLocs map[Position]int
	TavernsLocs map[Position]int
	HeroesLocs map[Position]int
	PlayerId int
}

type Position struct {
	X, Y int
}

func tileToInt(tiles string, index int) int {
	tile := []rune(tiles)[index]
	str, _ := strconv.Atoi(string(tile))

	return str
}

func (board *Board) parseTile(tile string) interface{} {
	switch string([]rune(tile)[0]) {
	case AIR_TILE:
		return AIR
	case WALL_TILE:
		return WALL
	case TAVERN_TILE:
		return TAVERN
	case MINE_TILE:
		id := string([]rune(tile)[1])
		return &MineTile{id}
	case HERO_TILE:
		char := string([]rune(tile)[1])
		id, _ := strconv.Atoi(char)
		return &HeroTile{id}
	default:
		return -3
	}
}

func (board *Board) parseTiles() {
	var vector [][]rune
	var matrix [][][]rune
	ts := make([][]interface{}, board.Size)
	board.MinesLocs = make(map[Position]int)
	board.TavernsLocs = make(map[Position]int)
	board.HeroesLocs = make(map[Position]int)

	for i := 0; i <= len(board.Tiles)-2; i = i + 2 {
		vector = append(vector, []rune(board.Tiles)[i:i+2])
	}

	for i := 0; i < len(vector); i = i + board.Size {
		matrix = append(matrix, vector[i:i+board.Size])
	}

	for xi, x := range matrix {
		innerList := make([]interface{}, board.Size)
		for xsi, xs := range x {

			tmp := board.parseTile(string(xs))
			switch i := tmp.(type) {
				case *MineTile:
					tmpHid, _ := strconv.Atoi(i.HeroId)
					board.MinesLocs[Position{xi , xsi}] = tmpHid					
					//fmt.Printf("\n %v , %v \n", xi, xsi)
				case *HeroTile:
					board.HeroesLocs[Position{xi, xsi}] = i.Id
				case int: 
					if i == TAVERN {
						board.TavernsLocs[Position{xi, xsi}] = 0 
					}

			}
			innerList[xsi] = tmp 
		}
		ts[xi] = innerList
	}

	board.Tileset = ts
}

func (board *Board) Passable(loc Position) bool {
	if _, ok := board.HeroesLocs[loc]; ok {
		return false
	}
	tile := board.Tileset[loc.X][loc.Y]
	return tile != WALL && tile != TAVERN && reflect.TypeOf(tile).String() != "MineTile"
}

func (board *Board) To(loc Position, direction Direction) *Position {
	row := loc.X
	col := loc.Y
	dLoc := AIM[direction]
	nRow := row + dLoc.X
	if nRow < 0 {
		nRow = 0
	}
	if nRow > board.Size {
		nRow = board.Size
	}
	nCol := col + dLoc.Y
	if nCol < 0 {
		nCol = 0
	}
	if nCol > board.Size {
		nCol = board.Size
	}
	
	newPos := &Position{nRow, nCol}

	if !board.Passable(*newPos) {
		if p, ok := board.MinesLocs[*newPos]; ok {
			if p != (board.PlayerId + 1) {
				fmt.Println("\n\n GETTING A MINE BITCH \n\n")
				return newPos
			}
		}
		return &loc
	 }
	 
		

	return newPos 
}
