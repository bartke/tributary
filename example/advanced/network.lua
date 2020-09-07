local tb = require("tributary")

-- setup network
-- source --> bets_window --> window_query --> printer
tb.create_window("bets_window")
tb.link("streaming_ingest", "bets_window")
-- select customer_id from bets where sport='soccer
local query = [[
    select
	sum(stake*exchange_rate*(b.odds -1)) as liability,
	customer_uuid,
	game_id
from
	bets b
join selections s on
	b.uuid = s.bet_uuid
where
	datetime(b.create_time / 1000000000, 'unixepoch') >= datetime('now', '-1 minute')
group by
	customer_uuid,
	game_id
having
	liability > 100
]]
tb.query_window("window_query", query)
tb.link("bets_window", "window_query")
tb.link("window_query", "printer")

-- run all network ndoes
tb.run("streaming_ingest")
tb.run("bets_window")
tb.run("window_query")
tb.run("printer")
