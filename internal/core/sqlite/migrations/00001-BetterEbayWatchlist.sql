CREATE TABLE ebay_watchlist_auctions (
  id         NUMBER PRIMARY KEY NOT NULL,
  url        TEXT NOT NULL,
  title      TEXT NOT NULL,
  price      NUMBER NOT NULL,
  ends_at    NUMBER NOT NULL,
  updated_at NUMBER NOT NULL,
  created_at NUMBER NOT NULL
);

CREATE INDEX ix_ebay_watchlist_auctions_id ON ebay_watchlist_auctions (id);
CREATE INDEX ix_ebay_watchlist_auctions_url_ends_at ON ebay_watchlist_auctions (url, ends_at);

CREATE TABLE ebay_watchlist_auction_notifications (
  id                                   NUMBER PRIMARY KEY NOT NULL,
  ebay_watchlist_auction_id            NUMBER NOT NULL,
  minutes_before_auction_end_to_notify NUMBER NOT NULL,
  sent_at                              NUMBER,
  updated_at                           NUMBER NOT NULL,
  created_at                           NUMBER NOT NULL
);

CREATE INDEX ix_ebay_watchlist_auction_notifications_ebay_watchlist_auction_id ON ebay_watchlist_auction_notifications (ebay_watchlist_auction_id);
CREATE INDEX ix_ebay_watchlist_auction_notifications_ebay_watchlist_auction_id_sent_at ON ebay_watchlist_auction_notifications (ebay_watchlist_auction_id, sent_at);
