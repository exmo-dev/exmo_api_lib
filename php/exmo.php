function api_query($api_name, array $req = array())
{
    $mt = explode(' ', microtime());
    $NONCE = $mt[1] . substr($mt[0], 2, 6);

    // API settings
    $key = "K-c8e39f2516b19586214b02c3ce4ab70981e1d7b9"; //TODO replace with your api key from profile page
    $secret = "S-6f4f6576f90814edfc0cc14f5ec6dc5c55a45af1"; //TODO replace with your api secret from profile page

    $url = "http://api.exmo.com/v1/$api_name";

    $req['nonce'] = $NONCE;

    // generate the POST data string
    $post_data = http_build_query($req, '', '&');

    $sign = hash_hmac('sha512', $post_data, $secret);

    // generate the extra headers
    $headers = array(
        'Sign: ' . $sign,
        'Key: ' . $key,
    );

    // our curl handle (initialize if required)
    static $ch = null;
    if (is_null($ch)) {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_USERAGENT, 'Mozilla/4.0 (compatible; PHP client; ' . php_uname('s') . '; PHP/' . phpversion() . ')');
    }
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $post_data);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE);

    // run the query
    $res = curl_exec($ch);
    if ($res === false) throw new Exception('Could not get reply: ' . curl_error($ch));
   
    $dec = json_decode($res, true);
    if ($dec === null)
        throw new Exception('Invalid data received, please make sure connection is working and requested API exists');

    return $dec;
}
