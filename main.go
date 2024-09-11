package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

type Post struct {
	ID         int    `json:"id"`
	Date       string `json:"date"`
	Slug       string `json:"slug"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Categories []int  `json:"categories"`
	Tags       []int  `json:"tags"`
	Title      struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered  string `json:"rendered"`
		Protected bool   `json:"protected"`
	} `json:"content"`
}

func main() {
	// MySQL connection
	db, err := sql.Open("mysql", "root:kebersamaan@@tcp(127.0.0.1:3306)/fetch")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	url := "https://funbahasa.com/wp-json/wp/v2/posts?_fields=id,title,date,slug,status,type,categories,tags,content&order=desc&per_page=100&page=3"

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("Referer", url) // Use the variable for consistency
	req.Header.Set("Cookie", "wordpress_test_cookie=WP%20Cookie%20check; wp_lang=en_US; wp-settings-1=libraryContent%3Dbrowse%26editor%3Dtinymce%26hidetb%3D1%26editor_plain_text_paste_warning%3D2; wp-settings-time-1=1719470724; __gads=ID=9b8d7da93571d6a1:T=1721036096:RT=1721036096:S=ALNI_Mb7d3SDuARZNtJOPO6s5ygh3yJPoA; __gpi=UID=00000e9293890e3c:T=1721036096:RT=1721036096:S=ALNI_MbDc-cQdinvoSvPchk86BA4wGYPLQ; __eoi=ID=6f67eba16e4445ee:T=1721036096:RT=1721036096:S=AA-AfjbM38fonwgMsizutJlaXMJK; cf_clearance=L_HlFVQBgzM_jL_VhGecSckoHbKRbB7FmP_5snMMUgI-1726030703-1.2.1.1-C8ETgT1JDY.U9b1aPpl1kXxU6IWS7937ViJsfINRuDIBk3patXh1cvXnjhNy1LPrE_ZZugDQsEOLYjOKBc5kw8QhUvvcNUcJ7JGW51f0GPkVsOQSQA7tKroONloas_CUAsaVt1Sj6APDcfgxaPNALT3XoWq4Qn8hcOFsiswUII3ZNgMAgDQxQJU_Le7IZTGiyrQTB1Rs3WCP7u.PFc9cIP_I3Ntby3DIdVLpgJHomVs8YgPhIRD.5TAy4fRaJo0TdS5e_ug4x8ezBYVg4awSvn_l_HBnWo2UtbUmrNFFMb8w.xl1sidyDRRe7IvWnVArKCSGnFONJT344VRPzrj27Q0wweJu0VcQnxe8I6DBAGOQE58cwkm.tnNPurta0KzDu1HNzRR.ALVyLymqhhw8JVeCB4UO7I5uSrKI3YAi5gS2QqxRnFoBngGzMzJ3kj1d")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Unmarshal the JSON response
	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Prepare statement for insertion
	stmt, err := db.Prepare("INSERT INTO posts (id, title, slug, status, type, categories,tags,content, created_at) VALUES (?, ?, ?, ?, ?, ?,?,?, ?)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()

	// Loop through the posts and insert into MySQL
	for _, post := range posts {
		// Attempt to parse the date with and without timezone
		createdAt, err := time.Parse(time.RFC3339, post.Date)
		if err != nil {
			// Try parsing without timezone
			createdAt, err = time.Parse("2006-01-02T15:04:05", post.Date)
			if err != nil {
				log.Fatalf("Error parsing date: %v", err)
			}
		}

		// Safely handle Categories and Tags
		var category, tag int
		if len(post.Categories) > 0 {
			category = post.Categories[0]
		} else {
			category = 0 // or a default value
		}
		if len(post.Tags) > 0 {
			tag = post.Tags[0]
		} else {
			tag = 0 // or a default value
		}

		if _, err = stmt.Exec(post.ID, post.Title.Rendered, post.Slug, post.Status, post.Type, category, tag, post.Content.Rendered, createdAt); err != nil {
			log.Fatalf("Error inserting into MySQL: %v", err)
		}

		fmt.Printf("Inserted post ID %d\n", post.ID)
	}
}
