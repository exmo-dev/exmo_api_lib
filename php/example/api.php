<?php

    class api {
        private $key = '';
        private $secret = '';
        private $url = '';
        public function __construct ($key, $secret, $url) {
            $this -> key = $key;
            $this -> secret = $secret;
            $this -> url = $url;
        }
        public function query ($name, $request = []) {
            try {
                return $this -> exicute ($name, $request);
            } catch (Exception $e) {
                return $e -> getMessage ();
            }
        }
        private function exicute ($name, $request = []) {
            $microtime = explode(' ', microtime ());
            $request['nonce'] = $microtime[1] . substr ($microtime[0], 2, 6);
            $post = http_build_query ($request, '', '&');
            $sign = hash_hmac('sha512', $post, $this -> secret);
            $headers = ['Sign: ' . $sign, 'Key: ' . $this -> key];
            
            $curl = curl_init ();
            curl_setopt ($curl, CURLOPT_RETURNTRANSFER, true);
            curl_setopt ($curl, CURLOPT_USERAGENT, 'Mozilla/4.0 (compatible; PHP client; ' . php_uname('s') . '; PHP/' . phpversion() . ')');
            curl_setopt ($curl, CURLOPT_URL, $this -> url . $name);
            curl_setopt ($curl, CURLOPT_POSTFIELDS, $post);
            curl_setopt ($curl, CURLOPT_HTTPHEADER, $headers);
            curl_setopt ($curl, CURLOPT_SSL_VERIFYPEER, FALSE);
            $response = curl_exec ($curl);
            
            if ($response === false) {
                throw new Exception ('Failed to get answer' . curl_error ($curl));
            } else {
                $result = json_decode ($response, true);
                if ($result === null) {
                    throw new Exception ('Invalid data received');
                } else {
                    return $result;
                }
            }
        }
    }
