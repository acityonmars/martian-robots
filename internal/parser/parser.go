package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/acityonmars/martian-robots/internal/engine"
)

type robotFieldValue int

const (
	XPos robotFieldValue = iota
	YPos
	RobotOrientation
)

// maxCoordinate is the possible maximum coordinate.
const maxCoordinate int = 50

// minInputLineCnt is the minimum number of lines for input data. i.e. 1 line of world coordinate + 2 lines for a robot and instructions
const minInputLineCnt int = 3

// maxInstructionCnt is the maximum number of instructions letters.
const maxInstructionCnt int = 100

// ParseInputData parses input data in a string slice, and returns mars.World and mars.Mission.
func ParseInputData(inputData []string, logger *slog.Logger) (*engine.World, *engine.Mission, error) {
	inputDataLineCnt := len(inputData)
	if inputDataLineCnt < minInputLineCnt || inputDataLineCnt%2 == 0 {
		return nil, nil, fmt.Errorf("invalid input data with incorrect line cnt: %d", inputDataLineCnt)
	}

	upperRightPos, err := parseWorldCoordinate(inputData[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse world coordinate: %w", err)
	}

	world := engine.NewWorld(*upperRightPos, logger)

	mission, err := parseMission(inputData[1:], world)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse mission data: %w", err)
	}

	return world, mission, nil
}

func parseWorldCoordinate(coordinate string) (*engine.Position, error) {
	fields := strings.Fields(coordinate)
	fieldCnt := len(fields)
	if fieldCnt != 2 {
		return nil, fmt.Errorf("expected two fields for coordinate but got %d", fieldCnt)
	}

	tmpPositions := make([]int, 2)
	for i, field := range fields {
		num, err := parseNumString(field)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a field to number: %w", err)
		}

		if num > maxCoordinate {
			return nil, fmt.Errorf("exceeded maximum coordinate: %d", num)
		}

		tmpPositions[i] = num
	}

	upperRightPos := engine.Position{
		X: tmpPositions[0],
		Y: tmpPositions[1],
	}

	return &upperRightPos, nil
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
func parseMission(missionRawData []string, world *engine.World) (*engine.Mission, error) {
	missionLineCnt := len(missionRawData)
	if missionLineCnt%2 != 0 {
		return nil, fmt.Errorf("incorrect line number for raw mission data: %d", missionLineCnt)
	}

	mission := engine.Mission{}
	robotPlan := engine.RobotPlan{}
	for i, line := range missionRawData {
		// Odd line is for robot input. (e.g. "1 1 E")
		if i%2 == 0 {
			robot, err := parseRobot(line)
			if err != nil {
				return nil, fmt.Errorf("failed to parse robot info: %w", err)
			}
			robot.ID = i / 2
			robotPlan.Robot = robot
			robotPlan.World = world
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

func parseRobot(rawRobotInfo string) (*engine.Robot, error) {
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

	robot := engine.Robot{
		Pos: engine.Position{X: x, Y: y}, Ori: orientation,
	}

	return &robot, nil
}

// parseOrientation parses raw input orientation such as "N" to an Orientation.
func parseOrientation(rawInput string) (engine.Orientation, error) {
	switch rawInput {
	case "N":
		return engine.North, nil
	case "E":
		return engine.East, nil
	case "S":
		return engine.South, nil
	case "W":
		return engine.West, nil
	default:
		return 0, fmt.Errorf("invalid orientation: %q", rawInput)
	}
}

// parseCommands parses raw input commands such as "RFRFRFRF" to a mars.Command slice.
func parseCommands(rawInput string) (*[]engine.Command, error) {
	trimmedInput := strings.TrimSpace(rawInput)
	if len(trimmedInput) == 0 {
		return nil, errors.New("no command is given or found")
	}

	if len(trimmedInput) > maxInstructionCnt {
		return nil, fmt.Errorf("instruction letters exceed the maximum length: %d", maxInstructionCnt)
	}

	var commands []engine.Command
	for i, commandLetter := range trimmedInput {
		command, ok := engine.CommandFor(commandLetter)
		if !ok {
			return nil, fmt.Errorf("unsupported command letter: %q at position: %d", commandLetter, i)
		}
		commands = append(commands, command)
	}

	return &commands, nil
}
