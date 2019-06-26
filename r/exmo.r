library(httr)
library(jsonlite)
library(nanotime)
library(digest)

api_url <- "https://api.exmo.com/v1/"
api_key <- "K-..."
api_secret <- "S-..."

api_query <- function(method, key, secret, params = list()){
	
	nonce <- (as.numeric(as.POSIXct(Sys.time()))*10000000)%/%1
	
	params <- c(params, nonce = nonce)
	
	data <- paste(names(params), params, sep = "=", collapse = "&")

	signature <- 
		hmac( 
			key 	= secret, 
			object 	= data, 
			algo 	= "sha512" 
		)
	
	response <- POST( 
		url  = paste(api_url, method, sep = ""),
		accept_json(),
		add_headers(
			"Key" = key, 
			"Sign" = signature,
			"Content-Type" = "application/x-www-form-urlencoded"
		),
		body = data
	)

	exmo_content <- content(response, as="text")
	json_exmo_content <- fromJSON(exmo_content)
	return(json_exmo_content)
}

params <- list(
	pair 	= "ETH_BTC",
	limit 	= 100,
	offset 	= 0

)

api_query("user_trades", api_key, api_secret, params)
api_query("user_info", api_key, api_secret)