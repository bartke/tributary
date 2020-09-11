local tb = require("tributary")

-- setup network
-- source --> bets_window --> window_query --> printer
--                        --> stream_split (script 2)
tb.create_window("bets_window")
tb.link("streaming_ingest", "bets_window")

-- select aggregate customer liability if > 130
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

-- here we create a stream split to reutilize what was set up before
tb.create_forwarder("stream_split")
tb.fanout("bets_window", "window_query", "stream_split")

tb.create_filter("dedupe_liability", 10)
tb.link("window_query", "dedupe_liability")
tb.link("dedupe_liability", "printer")
