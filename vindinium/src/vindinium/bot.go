package vindinium

import (
	"math/rand"
	"strings"
	path "pathfinding"
)

type Direction string

var DIRS = []Direction{"Stay", "North", "South", "East", "West"}

var MyID = "wt8tmbq6"
var Name = "kwilson"

var heroIndex = -1

func randDir() Direction {
	dir := DIRS[rand.Intn(len(DIRS))]

	return dir
}

type Bot interface {
	Move(state *State) Direction
}

type RandomBot struct{}

func (b *RandomBot) Move(state *State) Direction {
	return randDir()
}

type FighterBot struct{}

func pathValue(nodes []*path.Node, isChest bool, b *Board, g *Game) int {

	val := len(nodes)

	if isChest && len(nodes) > 0 {
		//greater value for uncontrolled mines
		if v, ok := b.MinesLocs[Position{nodes[len(nodes) - 1].X, nodes[len(nodes) - 1].Y}]; ok {
			if v == 0 {
				val = val - 2
				if len(nodes) > 6 {
					val = val - 1
				}
			}
		} 
	}

	//plus one for each hero in path
	for _, n := range nodes {
		if h, ok := b.HeroesLocs[Position{n.X, n.Y}]; ok {
			if g.Heroes[h - 1].Life < int(g.Heroes[heroIndex].Life / 2) {
				val = val + 1
			} else {
				val = val + 2
			}
		}
	}

	return val	

}

func dist(p1, p2 *Position, b *Board) []*path.Node {

	mapStr := buildGraph(*p1, *p2, b)
	map_data := Read_Map(mapStr)
	graph := path.NewGraph(map_data)
	nodes_path := path.Astar(graph)
	return nodes_path

}

var lastPos = Position{0,0}
var lastDir = DIRS[0]
var minesGone map[Position]int
var mineCount = -1

var (
	onPath = false
	nodeCount = 0
	nodes []*path.Node

)

func (b *FighterBot) Move(state *State) Direction {
	 g := NewGame(state)

	if mineCount == -1 {
		minesGone = make(map[Position]int)
		mineCount = len(g.HeroesLocs)
	}	
	
	for i, h := range g.Heroes {
		if h.Name == Name {
			heroIndex = i
		}
	}
	g.Board.PlayerId = heroIndex
	//fmt.Printf("\n myHero: %d, %d\n", g.Heroes[heroIndex].Pos.X, g.Heroes[heroIndex].Pos.Y)
	curPos := g.Heroes[heroIndex].Pos
	var tmpVal = 1000000
	var tmpPath []*path.Node
	var newPos Position
	//fmt.Printf("Mineloc count : %d , TavernLoc count : %d , heroes count: %d\n", len(g.MinesLocs), len(g.TavernsLocs), len(g.HeroesLocs))
	if nodeCount >= len(nodes) || !onPath {
	  	if g.Heroes[heroIndex].Life > 35  {   
			for pos, play := range g.MinesLocs {
				tP := dist(curPos, &pos, g.Board)
				if tV := pathValue(tP, true, g.Board, g); tV < tmpVal {
					if play != (heroIndex + 1){
							tmpVal = tV 
							tmpPath = tP 
					}
				}
			}	
			onPath = true	
			nodes = tmpPath 
			nodeCount = 1
			if nodeCount <= len(nodes) - 1{
				newPos = Position{nodes[nodeCount].X, nodes[nodeCount].Y}
			} else {
				onPath = false
			}
		} else {
			for pos, _ := range g.TavernsLocs {
				tP := dist(curPos, &pos, g.Board);
				if tV := pathValue(tP, false, g.Board, g); tV < tmpVal {
					tmpVal = tV 
					tmpPath = tP
				}	
			}
			onPath = true	
			nodes = tmpPath 
			nodeCount = 1
			if nodeCount <= len(nodes) - 1 {
				newPos = Position{nodes[nodeCount].X, nodes[nodeCount].Y}
			} else {
				onPath = false
			}
		}		

	} else {
		newPos = Position{nodes[nodeCount].X, nodes[nodeCount].Y}
		nodeCount++
	}

	tmpPos := HeroAround(*curPos, g.Board)
	if tmpPos.X != curPos.X || tmpPos.Y != curPos.Y {
		if g.Heroes[heroIndex].Life > 60 {
			newPos = tmpPos
			onPath = false
		}  
	}

	tmpPos2 := TavernAround(*curPos, g.Board)
	if tmpPos2.X != curPos.X || tmpPos2.Y != curPos.Y {
		if g.Heroes[heroIndex].Life <= 80 {
			newPos = tmpPos2
			onPath = false
		}
	}	

	if _, ok := g.MinesLocs[newPos]; ok {
		minesGone[newPos] = minesGone[newPos] + 1
	}
	
	if _, ok := g.HeroesLocs[newPos]; ok {
		if g.Heroes[heroIndex].Life < 40 {
			onPath = false
			if newPos.X > curPos.X {
				if g.Board.Passable(Position{curPos.X - 1, curPos.Y}) {
					return "North"	
				}	
			} else if newPos.X < curPos.X {
				if g.Board.Passable(Position{curPos.X + 1, curPos.Y}) {
					return "South"
				}

			} else if newPos.Y > curPos.Y {
				if g.Board.Passable(Position{curPos.X, curPos.Y - 1}){
					return "West"
				}	
			} else if newPos.Y < curPos.Y {
				if g.Board.Passable(Position{curPos.X, curPos.Y + 1}){
					return "East"
				}
			}	
		}
	}

	if newPos.X > curPos.X {
		return "South"
	} else if newPos.X < curPos.X {
		return "North"
	} else if newPos.Y > curPos.Y {
		return "East"	
	} else if newPos.Y < curPos.Y {
		return "West"
	}


	return	"Stay" 
}

func buildGraph(curPos, chestPos Position, b *Board) string {

	var mapStr string
	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			if curPos.X == x && curPos.Y == y{
				mapStr = mapStr + "s"
				continue
			} else if chestPos.X == x && chestPos.Y == y {
				mapStr = mapStr + "e"
				continue
			} else if _, ok := b.MinesLocs[Position{x, y}]; ok {
				mapStr = mapStr + "#"
				continue
			}
			
			if b.Passable(Position{x,y}) {
				mapStr = mapStr + "."
			} else {
				mapStr = mapStr + "#"
			}
		}
		if x != b.Size - 1 {
			mapStr = mapStr + "\n"
		}
	}

	//fmt.Println(mapStr)
	return mapStr
}

func HeroAround(curPos Position, b *Board) Position {

	if _, ok := b.HeroesLocs[Position{curPos.X + 1, curPos.Y}]; ok {
		return Position{curPos.X + 1, curPos.Y}
	} 
	if _, ok := b.HeroesLocs[Position{curPos.X - 1, curPos.Y}]; ok {
		return Position{curPos.X - 1, curPos.Y}
	}
 	if _, ok := b.HeroesLocs[Position{curPos.X, curPos.Y + 1}]; ok {
		return Position{curPos.X, curPos.Y + 1}
	}
 	if _, ok := b.HeroesLocs[Position{curPos.X + 1, curPos.Y - 1}]; ok {
		return Position{curPos.X, curPos.Y - 1}
	}

	return curPos 
}

func TavernAround(curPos Position, b *Board) Position {

	if _, ok := b.TavernsLocs[Position{curPos.X + 1, curPos.Y}]; ok {
		return Position{curPos.X + 1, curPos.Y}
	} 
	if _, ok := b.TavernsLocs[Position{curPos.X - 1, curPos.Y}]; ok {
		return Position{curPos.X - 1, curPos.Y}
	}
 	if _, ok := b.TavernsLocs[Position{curPos.X, curPos.Y + 1}]; ok {
		return Position{curPos.X, curPos.Y + 1}
	}
 	if _, ok := b.TavernsLocs[Position{curPos.X + 1, curPos.Y - 1}]; ok {
		return Position{curPos.X, curPos.Y - 1}
	}

	return curPos 

}

func Read_Map(map_str string) *path.MapData {
        rows := strings.Split(map_str, "\n")
        if len(rows) == 0 {
                panic("The map needs to have at least 1 row")
        }
        row_count := len(rows)
        col_count := len(rows[0])

        result := *path.NewMapData(row_count, col_count)
        for i := 0; i < row_count; i++ {
                for j := 0; j < col_count; j++ {
                        char := rows[i][j]
                        switch char {
                        case '.':
                                result[i][j] = path.LAND
                        case '#':
                                result[i][j] = path.WALL
                        case 's':
                                result[i][j] = path.START
                        case 'e':
                                result[i][j] = path.STOP
                        }
                }
        }
        return &result
}
