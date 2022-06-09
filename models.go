package railwayNetwork

type User struct {
	Login        string
	PasswordHash string
	FName        string
	LName        string
	PassportNum  string
}

type Route struct {
	RouteName   string
	RouteID     int
	arrivalTime string
	Stops       string
}

type StopPlaces struct {
	Stop string
}

type Ticket struct {
	train     string
	departure string
	arrival   string
}
