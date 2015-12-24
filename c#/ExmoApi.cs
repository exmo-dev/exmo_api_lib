using System;
using System.Collections.Generic;
using System.Collections.Specialized;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Security.Cryptography;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Web;

namespace ExmoAPI
{
    public class ExmoApi
    {
        private static long _nounce;
        // API settings
        private string _key; 
        private string _secret; 
        private string _url = "http://api.exmo.com/v1/{0}";

        static ExmoApi()
        {
            _nounce = Helpers.GetTimestamp();
        }

        public ExmoApi(string key, string secret)
        {
            _key = key;
            _secret = secret;
        }

        public async Task<string> ApiQueryAsync(string apiName, IDictionary<string, string> req)
        {
            using (var client = new HttpClient())
            {
                var n = Interlocked.Increment(ref _nounce);
                req.Add("nonce", Convert.ToString(n));
                var message = ToQueryString(req);

                var sign = Sign(_secret, message);

                var content = new FormUrlEncodedContent(req);
                content.Headers.Add("Sign", sign);
                content.Headers.Add("Key", _key);

                var response = await client.PostAsync(string.Format(_url, apiName), content);

                return await response.Content.ReadAsStringAsync();
            }
        }

        public async Task<HttpStatusCode> ApiQueryAsyncEx(string apiName, IDictionary<string, string> req)
        {
            using (var client = new HttpClient())
            {
                var n = Interlocked.Increment(ref _nounce);
                req.Add("nonce", Convert.ToString(n));
                var message = ToQueryString(req);

                var sign = Sign(_secret, message);

                var content = new FormUrlEncodedContent(req);
                content.Headers.Add("Sign", sign);
                content.Headers.Add("Key", _key);

                var response = await client.PostAsync(string.Format(_url, apiName), content);
                await Task.Factory.StartNew(async () =>
                {
                    var data = await response.Content.ReadAsStringAsync();
                    Console.WriteLine(data);
                });
                
                return response.StatusCode;
            }
        }

        public string ApiQuery(string apiName, IDictionary<string,string> req)
        {
            using (var wb = new WebClient())
            {
                req.Add("nonce", Convert.ToString(_nounce++));
                var message = ToQueryString(req);
                
                var sign = Sign(_secret, message);
               
                wb.Headers.Add("Sign", sign);
                wb.Headers.Add("Key", _key);

                var data = req.ToNameValueCollection();
                
                var response = wb.UploadValues(string.Format(_url, apiName), "POST", data);
                return Encoding.UTF8.GetString(response);
            }
        }

        private string ToQueryString(IDictionary<string, string> dic)
        {
            var array = (from key in dic.Keys
                         select string.Format("{0}={1}", HttpUtility.UrlEncode(key), HttpUtility.UrlEncode(dic[key])))
                .ToArray();
            return string.Join("&", array);
        }

        public static string Sign(string key, string message)
        {
            using (HMACSHA512 hmac = new HMACSHA512(Encoding.UTF8.GetBytes(key)))
            {
                byte[] b = hmac.ComputeHash(Encoding.UTF8.GetBytes(message));
                return ByteToString(b);
            }
        }

        public static string ByteToString(byte[] buff)
        {
            string sbinary = "";

            for (int i = 0; i < buff.Length; i++)
            {
                sbinary += buff[i].ToString("X2"); // hex format
            }
            return (sbinary).ToLowerInvariant();
        }

    }

    public static class Helpers
    {
        public static NameValueCollection ToNameValueCollection<TKey, TValue>(this IDictionary<TKey, TValue> dict)
        {
            var nameValueCollection = new NameValueCollection();

            foreach (var kvp in dict)
            {
                string value = string.Empty;
                if (kvp.Value != null)
                    value = kvp.Value.ToString();

                nameValueCollection.Add(kvp.Key.ToString(), value);
            }

            return nameValueCollection;
        }

        public static long GetTimestamp()
        {
            var d = (DateTime.Now - new DateTime(1970, 1, 1)).TotalMilliseconds;
            return (long)d;
        }
    }
}
