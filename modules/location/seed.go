package location

// Seed returns the default location document used to populate an empty store.
func Seed() Location {
	return Location{
		AddressLine1: "123 Market St",
		AddressLine2: "San Francisco, CA 94102",
		TravelRadius: 20,
	}
}

