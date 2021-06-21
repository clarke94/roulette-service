package bet

// Bet is a domain model.
type Bet struct {
	ID       string
	TableID  string
	Bet      string
	Type     string
	Amount   int64
	Currency string
}

// Result is the round result from a game.
type Result struct {
	Number  int
	Color   string
	Winners []Winner
}

// Winner is a winning bet from a round.
type Winner struct {
	BetID    string
	Amount   int64
	Currency string
}

const (
	colorRed   = "red"
	colorBlack = "black"
	colorGreen = "green"
)

// Type is the supported Bet type.
var (
	TypeRedBlack = "red/black"
	TypeStraight = "straight"
)

// TypeMultiplierMap is the Bet type that is available and the associated multiplier for that bet.
var TypeMultiplierMap = map[string]int64{
	// Outside bets
	TypeRedBlack: 1,
	// Inside bets
	TypeStraight: 35,
}

func betToWinner(b Bet) Winner {
	return Winner{
		BetID:    b.ID,
		Amount:   b.Amount,
		Currency: b.Currency,
	}
}

func betListToWinner(b []Bet) []Winner {
	winners := make([]Winner, len(b))

	for i := range b {
		winners[i] = betToWinner(b[i])
	}

	return winners
}
