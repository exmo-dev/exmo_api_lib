#include "exmo_api.hpp"
#include <iostream>

int main()
{
	exmo_api api("your_key", "your_secret");

	json_data response = api.call("user_info", "");
	std::clog << ">> user_info: " << response << "\n\n";
	
	response = api.call("user_cancelled_orders", api.build({ "limit=100", "offset=0" }));
	std::clog << ">> user_cancelled_orders: " << response << "\n\n";
	
	return 0;
}

