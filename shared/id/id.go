package shared_id

// Indentifier Type design mode
type AccountID string

func (a AccountID) String() string {
	return string(a)
}

type TripID string

func (a TripID) String() string {
	return string(a)
}