local tb = require("tributary")

-- select aggregate customer liability if > 130
local query = [[
select
	round(sum(stake*exchange_rate*(b.odds -1)), 2) as liability,
	customer_uuid,
	game_id
from
	example_bets b
join example_selections s on
	b.uuid = s.bet_uuid
where
	FROM_UNIXTIME(b.create_time) >= now() - interval 10 second
group by
	customer_uuid,
	game_id
having
	liability > 130
]]
tb.sliding_window_time("customer_liability", "window_out", query, "example_bets", "10s", "create_time")

-- here we create a stream split to reutilize what was set up before
tb.create_forwarder("stream_split")
tb.fanout("streaming_ingest", "customer_liability", "stream_split")

tb.create_filter("dedupe_liability", "10s")

-- setup network
tb.link("window_out", "dedupe_liability")
tb.link("dedupe_liability", "liability_printer")
