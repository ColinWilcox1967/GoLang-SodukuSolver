package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
        "flag"
)

const gridSize = 9

type (
	Sudoku [gridSize][gridSize]int
)

var (
	 puzzleFile string
)

func parseArguments() {
	flag.StringVar(&puzzleFile, "puzzle", "puzzle.txt", "Name of the Soduku puzzle to be completed.")
	flag.Parse()
}

func main() {
	parseArguments()

	puzzle, err := loadPuzzle(puzzleFile)
	if err != nil {
		fmt.Println("Error loading puzzle:", err)
		return
	}

	fmt.Println("Original Puzzle:")
	printPuzzle(puzzle)

	if solve(&puzzle) {
		fmt.Println("\nSolved Puzzle:")
		printPuzzle(puzzle)
	} else {
		fmt.Println("\nNo solution exists.")
	}
}

// Load the puzzle from a file
func loadPuzzle(filename string) (Sudoku, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Sudoku{}, err
	}
	defer file.Close()

	var puzzle Sudoku
	scanner := bufio.NewScanner(file)
	row := 0

	for scanner.Scan() {
		if row >= gridSize {
			break
		}
		line := scanner.Text()
		if len(line) < gridSize {
			return Sudoku{}, fmt.Errorf("invalid line: %s", line)
		}
		for col := 0; col < gridSize; col++ {
			ch := line[col]
			if ch == '.' || ch == '0' {
				puzzle[row][col] = 0
			} else {
				val, err := strconv.Atoi(string(ch))
				if err != nil || val < 1 || val > 9 {
					return Sudoku{}, fmt.Errorf("invalid number at row %d, col %d", row, col)
				}
				puzzle[row][col] = val
			}
		}
		row++
	}

	if err := scanner.Err(); err != nil {
		return Sudoku{}, err
	}
	return puzzle, nil
}

// Print the puzzle
func printPuzzle(p Sudoku) {
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			fmt.Printf("%d ", p[i][j])
		}
		fmt.Println()
	}
}

// Solve the puzzle using backtracking
func solve(p *Sudoku) bool {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if p[row][col] == 0 {
				for num := 1; num <= 9; num++ {
					if isSafe(*p, row, col, num) {
						p[row][col] = num
						if solve(p) {
							return true
						}
						p[row][col] = 0
					}
				}
				return false // trigger backtracking
			}
		}
	}
	return true // puzzle solved
}

// Check if num can be placed at p[row][col]
func isSafe(p Sudoku, row, col, num int) bool {
	// Row and column check
	for i := 0; i < gridSize; i++ {
		if p[row][i] == num || p[i][col] == num {
			return false
		}
	}
	// 3x3 box check
	startRow := row / 3 * 3
	startCol := col / 3 * 3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if p[startRow+i][startCol+j] == num {
				return false
			}
		}
	}
	return true
}
