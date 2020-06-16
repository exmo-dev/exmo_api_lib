#ifndef EXMO_API_HPP
#define EXMO_API_HPP

#include "hmac_sha512.hpp"
#include "connection.hpp"
#include <ctime>
#include <cassert>
#include <string>
#include <map>


class exmo_api {
public:
	exmo_api(const std::string& key, const std::string& secret)
	{
		key_ = key;
		secret_ = secret;
		url_ = "api.exmo.com/v1/";
		connection_ = http::connection();
		nonce_ = ::time(nullptr);
	}

	json_data call(const std::string& method, const std::string& p) {
		std::string params = "nonce=";
		nonce_++;
		params.append(std::to_string(nonce_));
		
		if (p.size() != 0) {
			params.append("&");
		}
		params.append(p);

		std::map<std::string, std::string> headers;
		headers["Content-type"] = "application/x-www-form-urlencoded";
		headers["Key"] = key_;
		headers["Sign"] = this->signature(params);

		connection_.request(url_ + method, http::post(), params, headers);
		return connection_.get_response();
	}

	std::string build(std::vector<std::string> params_) {

		std::string params = "";
		for (auto i : params_) {
			params += "&" + i;
		}
		return params;
	}

private:
	std::string signature(const std::string& params) {
		HMAC_SHA512 hmac_sha512(secret_, params);
		return hmac_sha512.hex_digest();
	}

private:
	unsigned long nonce_;
	std::string key_;
	std::string secret_;
	std::string url_;
	http::connection connection_;
};
#endif