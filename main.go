package main

import (
	"fmt"
	"gemini-webscrapping/common"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func bench() {
	type BenchTest struct {
		DirectoriiUrl string
		SiteUrl       string
	}
	tests := []BenchTest{
		{DirectoriiUrl: "https://directorii.com/scam-alerts/41/", SiteUrl: "https://www.bbb.org/us/oh/powell/profile/roofing-contractors/liberty-restoration-llc-0302-70090036"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/42/", SiteUrl: "https://www.woodtv.com/news/target-8/target-8-alert-unlicensed-roofer-stole-customers-credit-card/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/43/", SiteUrl: "https://www.wsbtv.com/news/local/cobb-county/elderly-woman-says-salesman-pressured-her-into-paying-98k-new-roof/GKCD6G22DBAMBKD5UUP6D53JBY/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/44/", SiteUrl: "https://www.wwltv.com/article/news/investigations/david-hammer/state-police-investigate-large-insurance-fraud-case-apex-mcclenny-moseley-asociates/289-601378ac-005f-4e2d-9d43-87d94dee7cdf"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/46/", SiteUrl: "https://www.youtube.com/watch?v=i7pmERRiNyM"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/47/", SiteUrl: "https://www.nbc4i.com/news/local-news/columbus/columbus-roofing-contractor-sued-for-taking-payments-despite-not-completing-work/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/48/", SiteUrl: "https://www.wtrf.com/wetzel-county/local-company-suing-roofer-for-failure-to-honor-warranty/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/49/", SiteUrl: "https://roofinginsights.com/roof-rejuvenation-shingle-magic/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/50/", SiteUrl: "https://www.wate.com/investigations/still-havent-got-my-roof-knoxville-widow-pursuing-legal-action-against-roofer/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/51/", SiteUrl: "https://www.atlantanewsfirst.com/2023/11/14/tired-excuses-customer-demands-gwinnett-roofer-return-insurance-check/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/52/", SiteUrl: "https://www.justice.gov/usao-nj/pr/middlesex-county-construction-company-admits-causing-death-employee-who-fell-roof-during"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/53/", SiteUrl: "https://www.cbsnews.com/texas/news/north-texas-woman-believes-late-stepmother-was-exploited-by-solar-panel-contractor-lender/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/55/", SiteUrl: "https://www.bbb.org/us/ca/roseville/profile/loans/goodleap-llc-1156-33013909"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/56/", SiteUrl: "https://wgxa.tv/news/local/shocked-grandmother-faces-26000-bill-for-unauthorized-roof-replacement-family-speaks-out-about-company-alexus-peavy-c-and-c-home-improvements-georgia-hawkinsville-tifton"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/57/", SiteUrl: "https://www.dol.gov/newsroom/releases/osha/osha20231212"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/58/", SiteUrl: "https://www.mitchellrepublic.com/news/local/one-person-in-custody-following-federal-investigation-thursday-at-mitchell-roofing-business"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/59/", SiteUrl: "https://www.wfmj.com/story/50293473/poland-roofing-contractor-charged-with-theft"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/60/", SiteUrl: "https://www.news-journal.com/shreveport-man-says-he-was-cheated-by-a-local-business/article_f311ac08-dd8b-5190-8364-117d6d8ce4fe.html"},
	}

	file, err := os.OpenFile("benchmark.md", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()

	for _, test := range tests {
		reviews, err := common.ScrapReviews(test.SiteUrl)
		log.Println(test.SiteUrl, "was scrapped")
		if err != nil {
			log.Printf("Error scraping reviews from %s: %v", test.SiteUrl, err)
			continue
		}
		_, err = file.WriteString("\n\n----------------------\n# " + test.DirectoriiUrl + "\n\n" + test.SiteUrl + "\n")
		if err != nil {
			log.Fatal("Failed to write to file: ", err)
		}
		for i, review := range reviews {
			output := fmt.Sprintf("\nName: %s\n\nMessage: %s\n\nDate: %s\n\n", review.Name, review.Message, review.Date)
			if review.Mark != nil {
				output += fmt.Sprintf("Mark: %v\n\n", *review.Mark)
			}
			if i < len(reviews)-1 {
				output += fmt.Sprintf("*************************\n")
			}
			_, err := file.WriteString(output)
			if err != nil {
				log.Fatal("Failed to write to file: ", err)
			}
		}
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bench()
}
