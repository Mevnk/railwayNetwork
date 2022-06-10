package railwayNetwork

import "time"

func TimeDiff(departure string, actualDeparture string) float64 {
	departure1, _ := time.Parse("15:04", departure)
	actualDeparture1, _ := time.Parse("15:04", actualDeparture)

	difference := departure1.Sub(actualDeparture1)

	return difference.Minutes()

}
