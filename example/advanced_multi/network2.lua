local tb = require("tributary")

-- last one's stake >= 100 on game 654321
local query = [[
    select
	round(stake*exchange_rate, 2) as stake,
	customer_uuid,
	game_id
from
	example_bets b
join example_selections s on
	b.uuid = s.bet_uuid
where
	stake*exchange_rate >= 100
	and game_id = 654321
order by
	b.create_time desc
limit 1
]]
tb.query_window("window_query2", query)
tb.link("stream_split", "window_query2")

tb.create_filter("dedupe_stake", "10s")
tb.link("window_query2", "dedupe_stake")
tb.link("dedupe_stake", "printer2")
