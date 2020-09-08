local tb = require("tributary")

-- setup network
-- source --> bets_window --> window_query --> printer
tb.create_window("bets_window")
tb.link("streaming_ingest", "bets_window")

-- select customer_id from bets where sport='soccer
local query = [[
    select
	round(sum(stake*exchange_rate*(b.odds -1)), 2) as liability,
	customer_uuid,
	game_id
from
	bets b
join selections s on
	b.uuid = s.bet_uuid
where
	datetime(b.create_time / 1000000000, 'unixepoch') >= datetime('now', '-10 seconds')
group by
	customer_uuid,
	game_id
having
	liability > 130
]]
tb.query_window("window_query", query)
tb.link("bets_window", "window_query")

tb.create_filter("dedupe_liability", 60)
tb.link("window_query", "dedupe_liability")
tb.link("dedupe_liability", "printer")
