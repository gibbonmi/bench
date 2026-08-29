package target

// Symbol is the declaration the loader test queries.
func Symbol() {}

// T carries a func field, so a fixture can plant a registry-shaped consumer.
type T struct{ Run func() }
