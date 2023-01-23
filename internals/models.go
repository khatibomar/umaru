package internals

import "codeberg.org/omarkhatib/umaru/database"

type GithubRepos []struct {
	Name        string `json:"name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
}

type ViewData struct {
	Data any
	Year string
}

type MyAnimeList struct {
	Data []struct {
		Node struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			MainPicture struct {
				Medium string `json:"medium"`
				Large  string `json:"large"`
			} `json:"main_picture"`
		} `json:"node"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type Book struct {
	Name  string
	Image string
	Link  string
}

type Books struct {
	Reading  []database.GetReadingBookRow
	Done     []database.GetDoneBooksRow
	WantRead []database.GetWantToReadBooksRow
}
