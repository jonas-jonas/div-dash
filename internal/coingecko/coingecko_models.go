package coingecko

import "time"

type CoingeckoCoin struct {
	CoingeckoID string `json:"id"`
	SymbolID    string `json:"symbol"`
	Name        string `json:"name"`
}

// Example:
// {
//     "polkadot": {
//         "eur": 23.7
//     }
// }
type CoingeckoPriceResponse map[string]map[string]float64

type CoingeckoDetails struct {
	ID                 string            `json:"id"`
	Symbol             string            `json:"symbol"`
	Name               string            `json:"name"`
	Platforms          map[string]string `json:"platforms"`
	BlockTimeInMinutes int               `json:"block_time_in_minutes"`
	HashingAlgorithm   string            `json:"hashing_algorithm"`
	Categories         []string          `json:"categories"`
	Localization       struct {
		En   string `json:"en"`
		De   string `json:"de"`
		Es   string `json:"es"`
		Fr   string `json:"fr"`
		It   string `json:"it"`
		Pl   string `json:"pl"`
		Ro   string `json:"ro"`
		Hu   string `json:"hu"`
		Nl   string `json:"nl"`
		Pt   string `json:"pt"`
		Sv   string `json:"sv"`
		Vi   string `json:"vi"`
		Tr   string `json:"tr"`
		Ru   string `json:"ru"`
		Ja   string `json:"ja"`
		Zh   string `json:"zh"`
		ZhTw string `json:"zh-tw"`
		Ko   string `json:"ko"`
		Ar   string `json:"ar"`
		Th   string `json:"th"`
		ID   string `json:"id"`
	} `json:"localization"`
	Description struct {
		En   string `json:"en"`
		De   string `json:"de"`
		Es   string `json:"es"`
		Fr   string `json:"fr"`
		It   string `json:"it"`
		Pl   string `json:"pl"`
		Ro   string `json:"ro"`
		Hu   string `json:"hu"`
		Nl   string `json:"nl"`
		Pt   string `json:"pt"`
		Sv   string `json:"sv"`
		Vi   string `json:"vi"`
		Tr   string `json:"tr"`
		Ru   string `json:"ru"`
		Ja   string `json:"ja"`
		Zh   string `json:"zh"`
		ZhTw string `json:"zh-tw"`
		Ko   string `json:"ko"`
		Ar   string `json:"ar"`
		Th   string `json:"th"`
		ID   string `json:"id"`
	} `json:"description"`
	Links struct {
		Homepage                    []string    `json:"homepage"`
		BlockchainSite              []string    `json:"blockchain_site"`
		OfficialForumURL            []string    `json:"official_forum_url"`
		ChatURL                     []string    `json:"chat_url"`
		AnnouncementURL             []string    `json:"announcement_url"`
		TwitterScreenName           string      `json:"twitter_screen_name"`
		FacebookUsername            string      `json:"facebook_username"`
		BitcointalkThreadIdentifier interface{} `json:"bitcointalk_thread_identifier"`
		TelegramChannelIdentifier   string      `json:"telegram_channel_identifier"`
		SubredditURL                string      `json:"subreddit_url"`
		ReposURL                    struct {
			Github    []string      `json:"github"`
			Bitbucket []interface{} `json:"bitbucket"`
		} `json:"repos_url"`
	} `json:"links"`
	Image struct {
		Thumb string `json:"thumb"`
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"image"`
	CountryOrigin                string  `json:"country_origin"`
	GenesisDate                  string  `json:"genesis_date"`
	SentimentVotesUpPercentage   float64 `json:"sentiment_votes_up_percentage"`
	SentimentVotesDownPercentage float64 `json:"sentiment_votes_down_percentage"`
	MarketCapRank                int     `json:"market_cap_rank"`
	CoingeckoRank                int     `json:"coingecko_rank"`
	CoingeckoScore               float64 `json:"coingecko_score"`
	DeveloperScore               float64 `json:"developer_score"`
	CommunityScore               float64 `json:"community_score"`
	LiquidityScore               float64 `json:"liquidity_score"`
	PublicInterestScore          float64 `json:"public_interest_score"`
	CommunityData                struct {
		FacebookLikes            interface{} `json:"facebook_likes"`
		TwitterFollowers         int         `json:"twitter_followers"`
		RedditAveragePosts48H    float64     `json:"reddit_average_posts_48h"`
		RedditAverageComments48H float64     `json:"reddit_average_comments_48h"`
		RedditSubscribers        int         `json:"reddit_subscribers"`
		RedditAccountsActive48H  int         `json:"reddit_accounts_active_48h"`
		TelegramChannelUserCount interface{} `json:"telegram_channel_user_count"`
	} `json:"community_data"`
	DeveloperData struct {
		Forks                        int `json:"forks"`
		Stars                        int `json:"stars"`
		Subscribers                  int `json:"subscribers"`
		TotalIssues                  int `json:"total_issues"`
		ClosedIssues                 int `json:"closed_issues"`
		PullRequestsMerged           int `json:"pull_requests_merged"`
		PullRequestContributors      int `json:"pull_request_contributors"`
		CodeAdditionsDeletions4Weeks struct {
			Additions int `json:"additions"`
			Deletions int `json:"deletions"`
		} `json:"code_additions_deletions_4_weeks"`
		CommitCount4Weeks              int   `json:"commit_count_4_weeks"`
		Last4WeeksCommitActivitySeries []int `json:"last_4_weeks_commit_activity_series"`
	} `json:"developer_data"`
	PublicInterestStats struct {
		AlexaRank   int         `json:"alexa_rank"`
		BingMatches interface{} `json:"bing_matches"`
	} `json:"public_interest_stats"`
	LastUpdated time.Time `json:"last_updated"`
	Tickers     []struct {
		Base   string `json:"base"`
		Target string `json:"target"`
		Market struct {
			Name                string `json:"name"`
			Identifier          string `json:"identifier"`
			HasTradingIncentive bool   `json:"has_trading_incentive"`
		} `json:"market"`
		Last          float64 `json:"last"`
		Volume        float64 `json:"volume"`
		ConvertedLast struct {
			Btc float64 `json:"btc"`
			Eth float64 `json:"eth"`
			Usd float64 `json:"usd"`
		} `json:"converted_last"`
		ConvertedVolume struct {
			Btc float64 `json:"btc"`
			Eth float64 `json:"eth"`
			Usd float64 `json:"usd"`
		} `json:"converted_volume"`
		TrustScore             string      `json:"trust_score"`
		BidAskSpreadPercentage float64     `json:"bid_ask_spread_percentage"`
		Timestamp              time.Time   `json:"timestamp"`
		LastTradedAt           time.Time   `json:"last_traded_at"`
		LastFetchAt            time.Time   `json:"last_fetch_at"`
		IsAnomaly              bool        `json:"is_anomaly"`
		IsStale                bool        `json:"is_stale"`
		TradeURL               string      `json:"trade_url"`
		TokenInfoURL           interface{} `json:"token_info_url"`
		CoinID                 string      `json:"coin_id"`
		TargetCoinID           string      `json:"target_coin_id,omitempty"`
	} `json:"tickers"`
}

type CoingeckoChart struct {
	Prices [][]float64
}
