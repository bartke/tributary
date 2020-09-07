## Tributary

Simple Go event stream processor.
- flow based, isolates concurrently running network nodes
- network nodes have in and out ports that can be connected

Simple Windowing using SQL
- alert too slack or telegram on query return
- clear window after alert successfully sent from next network node

Types
- `[]byte` json/octet-stream payloads
- Type registry?
- multi output
- filter reported, uniqueness

TODO: example with one windower DB, multiple networks.

As per [example/scripted](example/scripted/network.lua), we can set up networks at runtime with
lua scripts, such as

```lua
local tb = require("tributary")

tb.link("ticker_1s", "filter_even")
tb.link("filter_even", "printer")
```

If we register the network nodes on the lua module:

```go
import "github.com/bartke/tributary"
```

```go
	m := tributary.New()

	// register sources, pipelines, sinks
	tickerNode := common.NewTicker()
	m.RegisterSource("ticker_1s", tickerNode)
	filterNode := common.NewFilter()
	m.RegisterPipeline("filter_even", filterNode)
	printerNode := common.NewPrinter()
	m.RegisterSink("printer", printerNode)
```

Event
- payload flowing through the system

Sources
- Google PubSub
- RabbitMQ
- Ticker

Filter
- Matcher

Pipeline
- Matcher
- Window based
  - ToSQL
  - WindowSQL
  - TriggerSQL
  - JoinSQL

Sinks:
- Slack
- Callback
- Google PubSub
- RabbitMQ

```go

// payload: {"key":"test","value":5}

source := pubsub.New()
pipeline :=

```

```js
var sportsbet = JSON.parse(`{"selection":"123456/soccer.match_odds/home","sport":"soccer","odds":1.23,"stake":{"value":200,"currency":"USD","exchange_rate":1},"customer":88888888}`)
sportsbet.stake.value
```

```sql
select selection from sports_bets where sport='soccer'
select sum(stake.value*stake.exchange_rate) from sports_bets.window:time(30s) where sport='soccer' group by customer output when sum > 10000
-- uuid, stake_value, stake_exchange_rate, sport, customer, timestamp
select avg(odds) from sports_bets.window:length(100) where sport='soccer'
select sum(odds) from sports_bets.window:length(100) where sport='soccer' output when odds > 100
select sum(odds) from sports_bets.window:length(100) where sport='soccer' output every 5 minutes
select sum(odds) from sports_bets.window:length(100) where sport='soccer' output every 5 events
select customer from sports_bets.window:time(10m) join casino_bets.window:time(10m) on sports_bets.customer = casino_bets.customer output every 1 events
-- sports_bets: uuid, customer, timestamp
-- casino_bets: uuid, customer, timestamp
```
