package main

import (
	"gemini-webscrapping/common_review_scrapping"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func bench() {
	type benchTest struct {
		DirectoriiUrl string
		SiteUrls      string
	}
	tests := []benchTest{
		{DirectoriiUrl: "https://directorii.com/scam-alerts/41/", SiteUrls: "https://www.bbb.org/us/oh/powell/profile/roofing-contractors/liberty-restoration-llc-0302-70090036"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/42/", SiteUrls: "https://www.woodtv.com/news/target-8/target-8-alert-unlicensed-roofer-stole-customers-credit-card/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/43/", SiteUrls: "https://www.wsbtv.com/news/local/cobb-county/elderly-woman-says-salesman-pressured-her-into-paying-98k-new-roof/GKCD6G22DBAMBKD5UUP6D53JBY/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/44/", SiteUrls: "https://www.wwltv.com/article/news/investigations/david-hammer/state-police-investigate-large-insurance-fraud-case-apex-mcclenny-moseley-asociates/289-601378ac-005f-4e2d-9d43-87d94dee7cdf"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/46/", SiteUrls: "https://www.youtube.com/watch?v=i7pmERRiNyM"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/47/", SiteUrls: "https://www.nbc4i.com/news/local-news/columbus/columbus-roofing-contractor-sued-for-taking-payments-despite-not-completing-work/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/48/", SiteUrls: "https://www.wtrf.com/wetzel-county/local-company-suing-roofer-for-failure-to-honor-warranty/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/49/", SiteUrls: "https://roofinginsights.com/roof-rejuvenation-shingle-magic/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/50/", SiteUrls: "https://www.wate.com/investigations/still-havent-got-my-roof-knoxville-widow-pursuing-legal-action-against-roofer/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/51/", SiteUrls: "https://www.atlantanewsfirst.com/2023/11/14/tired-excuses-customer-demands-gwinnett-roofer-return-insurance-check/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/52/", SiteUrls: "https://www.justice.gov/usao-nj/pr/middlesex-county-construction-company-admits-causing-death-employee-who-fell-roof-during"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/53/", SiteUrls: "https://www.cbsnews.com/texas/news/north-texas-woman-believes-late-stepmother-was-exploited-by-solar-panel-contractor-lender/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/55/", SiteUrls: "https://www.bbb.org/us/ca/roseville/profile/loans/goodleap-llc-1156-33013909 https://www.cbsnews.com/texas/news/north-texas-woman-believes-late-stepmother-was-exploited-by-solar-panel-contractor-lender/"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/56/", SiteUrls: "https://wgxa.tv/news/local/shocked-grandmother-faces-26000-bill-for-unauthorized-roof-replacement-family-speaks-out-about-company-alexus-peavy-c-and-c-home-improvements-georgia-hawkinsville-tifton"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/57/", SiteUrls: "https://www.dol.gov/newsroom/releases/osha/osha20231212"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/58/", SiteUrls: "https://www.mitchellrepublic.com/news/local/one-person-in-custody-following-federal-investigation-thursday-at-mitchell-roofing-business"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/59/", SiteUrls: "https://www.wfmj.com/story/50293473/poland-roofing-contractor-charged-with-theft"},
		{DirectoriiUrl: "https://directorii.com/scam-alerts/60/", SiteUrls: "https://www.news-journal.com/shreveport-man-says-he-was-cheated-by-a-local-business/article_f311ac08-dd8b-5190-8364-117d6d8ce4fe.html"},
	}

	file, err := os.OpenFile("benchmark.md", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()

	for _, test := range tests {
		review, errors := common_review_scrapping.ScrapReview(test.SiteUrls)
		if len(errors) > 0 {
			for _, err := range errors {
				log.Printf("Error scraping review from %s: %v", test.SiteUrls, err)
			}
		} else {
			log.Println(test.SiteUrls, "was scrapped")
		}
		_, err = file.WriteString("\n----------------------\n# " + test.DirectoriiUrl + "\n\n" + test.SiteUrls + "\n\n" + review + "\n\n")
		if err != nil {
			log.Fatal("Failed to write to file: ", err)
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
