package main

type Serie struct {
	ID                 uint   `json:"id"`
	Title              string `json:"title"`
	Status             string `json:"status"`
	LastEpisodeWatched int    `json:"lastEpisodeWatched"`
	TotalEpisodes      int    `json:"totalEpisodes"`
	Ranking            int    `json:"ranking"`
}
