var exmo = require("./exmo");

exmo.init_exmo({key:'your_key', secret:'your_secret'});
//request version
exmo.api_query("user_info", { }, function(result){
	console.log(result);
	exmo.api_query("user_cancelled_orders", { "limit":1, "offset":0 }, function(result){
    console.log(result);

    //http nodejs version
    exmo.api_query2("user_info", { }, function(result){
      console.log(result);
      exmo.api_query2("user_cancelled_orders", { "limit":2, "offset":0 }, function(result){
        console.log(result);
      });
    });
  });
	
});
