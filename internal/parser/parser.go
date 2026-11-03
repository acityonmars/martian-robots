package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/acityonmars/martian-robots/internal/mars"
)

type fieldValue int

const (
	XPos fieldValue = iota
	YPos
	RobotOrientation
)

const minInputLineCnt int = 3

// TODO: document - inputdata must have an odd number of lines
func ParseMissionData(inputData []string) (*mars.Position, *mars.Mission, error) {
	inputDataLineCnt := len(inputData)
	if inputDataLineCnt < minInputLineCnt || inputDataLineCnt%2 == 0 {
		return nil, nil, fmt.Errorf("invalid input data with incorrect line cnt: %d", inputDataLineCnt)
	}

	upperRightPos, err := parseWorldCoordinate(inputData[0]) // TODO: use enum for 0
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse world coordinate: %w", err)
	}

	// TODO: parse mission
	mission, err := parseMission(inputData[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse mission data: %w", err)
	}

	return &upperRightPos, mission, nil
}

func parseWorldCoordinate(coordinate string) (mars.Position, error) {
	fields := strings.Fields(coordinate)
	if len(fields) != 2 {
		panic("there must be two numbers")
	}

	tmpPositions := make([]int, 2)
	for i, field := range fields {
		num, err := parseNumString(field)
		if err != nil {
			panic("")
		}
		tmpPositions[i] = num
	}

	upperRightPos := mars.Position{
		X: tmpPositions[0],
		Y: tmpPositions[1],
	}

	return upperRightPos, nil
}

func parseNumString(strNum string) (int, error) {
	num, err := strconv.Atoi(strNum)
	if err != nil {
		return 0, fmt.Errorf("invalid coordinate value %d: %w", num, err)
	}

	return num, nil
}

// parseMission parses raw mission data to mars.Mission.
//
// Raw data example:
//
//	1 1 E
//	RFRFRFRF
//	3 2 N
//	FRRFLLFFRRFLL
func parseMission(missionRawData []string) (*mars.Mission, error) {
	missionLineCnt := len(missionRawData)
	if missionLineCnt%2 != 0 {
		return nil, fmt.Errorf("incorrect line number for raw mission data: %d", missionLineCnt)
	}

	mission := mars.Mission{}
	robotPlan := mars.RobotPlan{}
	for i, line := range missionRawData {
		// Odd line is for robot input. (e.g. "1 1 E")
		if i%2 == 0 {
			robot, err := parseRobot(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse robot info: %w", err)
			}
			robot.ID = i
			robotPlan.Robot = *robot
		} else {
			// Even line is for the commands for the robot. (e.g. "RFRFRF")
			commands, err := parseCommands(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse commands: %w", err)
			}
			robotPlan.Commands = *commands
			mission.Plans = append(mission.Plans, robotPlan)
		}
	}

	return &mission, nil
}

func parseRobot(rawRobotInfo string) (*mars.Robot, error) {
	fields := strings.Fields(rawRobotInfo)
	fieldsCnt := len(fields)
	if fieldsCnt != 3 {
		return nil, fmt.Errorf("expected 3 fields for robot info but got %d - raw info: %s", fieldsCnt, rawRobotInfo)
	}

	x, err := parseNumString(fields[XPos])
	if err != nil {
		return nil, fmt.Errorf("failed to parse x coordinate for robot: %w", err)
	}
	y, err := parseNumString(fields[YPos])
	if err != nil {
		return nil, fmt.Errorf("failed to parse y coordinate for robot: %w", err)
	}

	orientation, err := parseOrientation(fields[RobotOrientation])
	if err != nil {
		return nil, fmt.Errorf("failed to parse orientation for robot: %w", err)
	}

	robot := mars.Robot{
		Position: mars.Position{X: x, Y: y}, Orientation: orientation,
	}

	return &robot, nil
}

// parseOrientation parses raw input orientation such as "N" to an Orientation.
func parseOrientation(rawInput string) (mars.Orientation, error) {
	switch rawInput {
	case "N":
		return mars.North, nil
	case "E":
		return mars.East, nil
	case "S":
		return mars.South, nil
	case "W":
		return mars.West, nil
	default:
		return 0, fmt.Errorf("invalid orientation: %q", rawInput)
	}
}

// parseCommands parses raw input commands such as "RFRFRFRF" to a mars.Command slice.
func parseCommands(rawInput string) (*[]mars.Command, error) {
	trimmedInput := strings.TrimSpace(rawInput)
	if len(trimmedInput) == 0 {
		return nil, errors.New("no command is given or found")
	}

	var commands []mars.Command
	for i, commandLetter := range trimmedInput {
		command, ok := mars.CommandFor(commandLetter)
		if !ok {
			return nil, fmt.Errorf("unsupported command letter: %q at position: %d", commandLetter, i)
		}
		commands = append(commands, command)
	}

	return &commands, nil
}
