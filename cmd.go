package main

import (
	"flag"
	"os"
	"bufio"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Simulation struct {
	Rows, Columns, Vehicles, Rides, Bonus, Steps int
}

type Ride struct {
	Start, End       Coordinates
	Earliest, Latest int
	Id               int
}

type Vehicle struct {
	Position       Coordinates
	CurrentRide    int
	CompletedRides []int
	Id             int
	AvailableAt    int
}

type Coordinates struct {
	X, Y int
}

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "in", "", "The path to the input file")

	var outPath string
	flag.StringVar(&outPath, "out", "", "The path to the output file")

	flag.Parse()

	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	firstLine := scanner.Text()
	simulation := createSimulation(firstLine)

	rides := createRides(scanner)
	initialRides := make([]Ride, len(rides))
	copy(initialRides, rides)

	bestScore := 0
	bestVehicles := make([]Vehicle, simulation.Vehicles)

	for x := 0; true; x++ {
		for i := range rides {
			rand.Seed(time.Now().UnixNano())
			j := rand.Intn(i + 1)
			rides[i], rides[j] = rides[j], rides[i]
		}

		totalDistance := 0

		vehicles := make([]Vehicle, simulation.Vehicles)
		for i := range vehicles {
			vehicles[i] = Vehicle{
				Position: Coordinates{
					X: 0,
					Y: 0,
				},
				CurrentRide:    -1,
				CompletedRides: []int{},
				Id:             i,
			}
		}

		for step := 0; step < simulation.Steps; step++ {
			for i, vehicle := range vehicles {
				if step >= vehicle.AvailableAt {
					closestRideIndex := -1
					shortestFinishTime := -1
					closestRide := Ride{}
					distanceToClosestRide := 0
					for j, ride := range rides {
						distanceToClosestRide = CalculateDistance(vehicle.Position, ride.Start)

						rideDistance := CalculateDistance(ride.Start, ride.End)
						timeUntilStart := ride.Earliest - step
						finishTime := step + int(math.Max(float64(timeUntilStart), float64(distanceToClosestRide))) + rideDistance

						if shortestFinishTime == -1 || (finishTime < ride.Latest && finishTime < shortestFinishTime) {
							closestRideIndex = j
							closestRide = ride
							shortestFinishTime = finishTime
						}
					}
					rideDistance := CalculateDistance(closestRide.Start, closestRide.End)
					timeUntilStart := closestRide.Earliest - step
					vehicle.AvailableAt += int(math.Max(float64(timeUntilStart), float64(distanceToClosestRide))) + rideDistance
					vehicle.CompletedRides = append(vehicle.CompletedRides, closestRide.Id)
					vehicle.Position = closestRide.End
					if len(rides) == 0 {
						break
					}
					rides = append(rides[:closestRideIndex], rides[closestRideIndex+1:]...)

					vehicles[i] = vehicle

					totalDistance += rideDistance
					if closestRide.Earliest == step {
						totalDistance += simulation.Bonus
					}
				}
			}
		}
		rides = make([]Ride, len(initialRides))
		copy(rides, initialRides)
		if totalDistance > bestScore {
			bestVehicles = vehicles
			bestScore = totalDistance
			fmt.Println(bestScore)
			ioutil.WriteFile(outPath, []byte(CreateOutput(bestVehicles)), 0644)
		}
		fmt.Println(x)
	}

	ioutil.WriteFile(outPath, []byte(CreateOutput(bestVehicles)), 0644)
}

func createSimulation(line string) Simulation {
	intSlice := toIntSlice(line)

	return Simulation{
		intSlice[0],
		intSlice[1],
		intSlice[2],
		intSlice[3],
		intSlice[4],
		intSlice[5],
	}
}

func createRides(s *bufio.Scanner) []Ride {
	rides := []Ride{}

	i := 0
	for s.Scan() {
		currentLine := s.Text()

		rides = append(rides, createRide(currentLine, i))
		i++
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	return rides
}

func createRide(line string, id int) Ride {
	intSlice := toIntSlice(line)

	return Ride{
		Coordinates{intSlice[0], intSlice[1]},
		Coordinates{intSlice[2], intSlice[3]},
		intSlice[4],
		intSlice[5],
		id,
	}
}

func toIntSlice(s string) []int {
	parts := strings.Split(s, " ")
	intSlice := make([]int, len(parts))

	for i, firstLinePart := range parts {
		firstLineInt, err := strconv.Atoi(firstLinePart)
		if err != nil {
			panic(err)
		}

		intSlice[i] = firstLineInt
	}

	return intSlice
}

func CalculateDistance(a, b Coordinates) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

func CreateOutput(vehicles []Vehicle) string {
	lines := make([]string, len(vehicles))

	for i, vehicle := range vehicles {
		lines[i] = strconv.Itoa(len(vehicle.CompletedRides)) + " "

		strs := make([]string, len(vehicle.CompletedRides))
		for j, completedRide := range vehicle.CompletedRides {
			strs[j] = strconv.Itoa(completedRide)
		}
		lines[i] += strings.Join(strs, " ")
	}

	return strings.Join(lines, "\n") + "\n"
}
