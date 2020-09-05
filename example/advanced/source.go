package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/example/advanced/event"
	"github.com/google/uuid"
	"github.com/infobloxopen/protoc-gen-gorm/types"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// one low frequency high staking customer and one high frequency low stake one
func customer() string {
	if randA.Intn(100) < 5 {
		return a
	}
	return b
}

func gameid() uint64 {
	if randB.Intn(100) < 50 {
		return 123456
	}
	return 654321
}

func customerStake(customer string) string {
	if customer == a {
		return "100.00"
	}
	return "1.00"
}

func sampleBet() *event.Bet {
	c := customer()
	return &event.Bet{
		Id:           &types.UUIDValue{uuid.Must(uuid.NewRandom()).String()},
		CreateTime:   timestamppb.New(time.Now().UTC()),
		CustomerUuid: c,
		Stake:        &event.Bet_Stake{Value: customerStake(c), Currency: "USD", ExchangeRate: "1.0"},
		Selections:   []*event.Selection{{GameId: gameid(), Market: "moneyline/home", Odds: "1.23"}},
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
