local tb = require("tributary")

-- setup network
-- source --> parser --> bets_window --> window_query --> printer
tb.parse("parser")
tb.link("streaming_ingest", "parser")
tb.create_window("bets_window", "sqlite", "memory", "json_parser")
tb.link("parser", "bets_window")
-- select customer_id from bets where sport='soccer
tb.query_window("window_query", "select customer_id from bets where sport='soccer'")
tb.link("bets_window", "window_query")
tb.link("window_query", "printer")

-- run all network ndoes
tb.run("streaming_ingest")
tb.run("parser")
tb.run("bets_window")
tb.run("window_query")
tb.run("printer")
