# Gemini webscrapper

It`s required to extract info from BBB, YouTube and unstructured US news sites.

[Example](/benchmark.md)

## Setup

Add `GEMINI_API_KEY` and `YOUTUBE_API_KEY` to /.env file.

## Run benchmark

It is required to launch program with USA IP address.

```bash
go mod tidy
go run .
```

## Architecture

![](/img/components.svg)

```puml
@startuml
package "System" {

  component "URL Input" as urlInput
  component "URL Classifier" as urlClass
  component "YouTube API Client" as ytClient
  component "BBB API Client" as bbbClient
  component "Rate Limiter for Gemini" as geminiLimiter
  component "Gemini API Client" as geminiClient
  component "Output Formatter" as output

  urlInput -right-> urlClass
  urlClass ..> ytClient : <<youtube url>>
  urlClass ..> bbbClient : <<bbb.org url>> 
  urlClass ..> geminiLimiter : <<other url>>
  geminiLimiter -right-> geminiClient
  ytClient -right-> output
  bbbClient -right-> output
  geminiClient -right-> output
  
}
@enduml
```
