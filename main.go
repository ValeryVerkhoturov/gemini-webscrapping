package main

import (
	"fmt"
	"gemini-webscrapping/common"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//reviews, err := common.ScrapReviews("https://www.bbb.org/us/oh/columbus/profile/roofing-contractors/liberty-restoration-llc-0302-70090036")
	//reviews, err := common.ScrapReviews("https://www.wsbtv.com/news/local/cobb-county/elderly-woman-says-salesman-pressured-her-into-paying-98k-new-roof/GKCD6G22DBAMBKD5UUP6D53JBY/")
	reviews, err := common.ScrapReviews("https://www.wwltv.com/article/news/investigations/david-hammer/state-police-investigate-large-insurance-fraud-case-apex-mcclenny-moseley-asociates/289-601378ac-005f-4e2d-9d43-87d94dee7cdf")

	//reviews, err := common.ScrapReviews("https://www.youtube.com/watch?v=i7pmERRiNyM")
	//reviews, err := common.ScrapReviews("https://youtu.be/YxYeMEpup9Q?si=yCu7nSxJnIPO_wDD")
	//reviews, err := common.ScrapReviews("https://www.youtube.com/shorts/zfcg-5Tw5qA")
	if err != nil {
		log.Fatal(err)
	}
	for _, review := range reviews {
		fmt.Println("\n*************************")
		fmt.Println(
			"Name:", review.Name,
			"\nMessage:", review.Message,
			"\nDate:", review.Date,
		)
		if review.Mark != nil {
			fmt.Println(
				"Mark:", *review.Mark,
			)
		}
	}
}
