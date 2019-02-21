library(httr)
library(jsonlite)
library(nanotime)
library(digest)

api_url <- "https://api.exmo.com/v1/user_info/"
api_key <- "K-..."
api_secret <- "S-..."

api_query <- function(url, key, secret){
	nonce <- (as.numeric(as.POSIXct(Sys.time()))*10000000)%/%1

	body <- 
		list(
			nonce = nonce
		)

	data <- paste(names(body), body, sep = "=", collapse = "&")

	signature <- 
		hmac( 
			key 	= secret, 
			object 	= data, 
			algo 	= "sha512" 
		)
	body
	responce <- POST( 
		url  = url,
		accept_json(),
		add_headers(
			"Key" = key, 
			"Sign" = signature,
			"Content-Type" = "application/x-www-form-urlencoded"
		),
		body = data,
		verbose()
	)

	exmo_content <- content(responce, as="text")# 
	json_exmo_content <- fromJSON(exmo_content)
	return(json_exmo_content)
}

api_query(api_url, api_key, api_secret)