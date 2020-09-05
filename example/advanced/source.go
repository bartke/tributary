package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/bartke/tributary"
	"github.com/google/uuid"
)

const (
	a = "76d17c9e-734b-452c-a5ee-852d1e6261bd"
	b = "66440a98-a92a-4ce0-a5aa-d851a6f288d9"
)

var randA, randB *rand.Rand

func init() {
	randA = rand.New(rand.NewSource(time.Now().UnixNano()))
	randB = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type selection struct {
	GameID   string `json:"game_id"`
	MarketID string `json:"market_id"`
	Odds     string `json:"odds"`
}

type stake struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type bet struct {
	RequestUUID  string      `json:"request_id"`
	CustomerUUID string      `json:"customer_uuid"`
	Stake        stake       `json:"stake"`
	Selections   []selection `json:"selections"`
	Odds         string      `json:"odds"`
}

// one low frequency high staking customer and one high frequency low stake one
func customer() string {
	if randA.Intn(100) < 5 {
		return a
	}
	return b
}

func gameid() string {
	if randB.Intn(100) < 50 {
		return "123456"
	}
	return "654321"
}

func customerStake(customer string) string {
	if customer == a {
		return "100.00"
	}
	return "1.00"
}

func sampleBet() *bet {
	c := customer()
	return &bet{
		RequestUUID:  uuid.Must(uuid.NewRandom()).String(),
		CustomerUUID: c,
		Stake:        stake{customerStake(c), "USD"},
		Selections:   []selection{{gameid(), "moneyline/home", "1.23"}},
		Odds:         "1.23",
	}
}

type msg struct {
	payload []byte
	ctx     context.Context
}

func Msg(p []byte, ctx context.Context) *msg {
	return &msg{payload: p, ctx: ctx}
}

func (m msg) Payload() []byte {
	return m.payload
}

func (m msg) Context() context.Context {
	return m.ctx
}

type stream struct {
	ticker *time.Ticker
	out    chan tributary.Event
}

func NewStream() *stream {
	return &stream{
		ticker: time.NewTicker(50 * time.Millisecond),
		out:    make(chan tributary.Event),
	}
}

func (s *stream) Run() {
	for {
		<-s.ticker.C
		m, _ := json.Marshal(sampleBet())
		s.out <- Msg(m, context.Background())
	}
}

func (s *stream) Out() <-chan tributary.Event {
	return s.out
}
