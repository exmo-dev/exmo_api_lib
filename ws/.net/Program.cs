using System;
using System.IO;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Security.Cryptography;

namespace WebSocketClient
{
    class Program
    {
        private static readonly string Connection = "wss://ws-api.exmo.com:443/v1/public";
        private static readonly string ConnectionPrivate = "wss://ws-api.exmo.com:443/v1/private";

        static async Task Main(string[] args)
        {
            TestPublic();
            TestPrivate();

            Console.ReadLine();
        }

        static void TestPublic()
        {
            Task.Factory.StartNew(async () =>
            {
                using (var socket = new ClientWebSocket())
                    try
                    {
                        await socket.ConnectAsync(new Uri(Connection), CancellationToken.None);
                        await SendPublic(socket);
                        await Receive(socket);

                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine($"ERROR - {ex.Message}");
                    }

            });
        }

        static void TestPrivate()
        {
            Task.Factory.StartNew(async () =>
            {
                using (var socket = new ClientWebSocket())
                    try
                    {
                        await socket.ConnectAsync(new Uri(ConnectionPrivate), CancellationToken.None);
                        await SendPrivate(socket);
                        await Receive(socket);

                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine($"ERROR - {ex.Message}");
                    }

            });
        }

        static async Task SendPublic(ClientWebSocket socket) => await Send(socket, "{\"id\":1,\"method\":\"subscribe\",\"topics\":[\"spot/trades:BTC_USD\",\"spot/ticker:LTC_USD\"]}");
        static async Task SendPrivate(ClientWebSocket socket) {
            var apiKey = "K-";
            var secretKey = "S-";
            var nonce = (long)(DateTime.Now - new DateTime(1970, 1, 1)).TotalMilliseconds;
            var sign = "";
            using (HMACSHA512 hmac = new HMACSHA512(Encoding.UTF8.GetBytes(secretKey)))
            {
                byte[] b = hmac.ComputeHash(Encoding.UTF8.GetBytes(apiKey + Convert.ToString(nonce)));
                sign = Convert.ToBase64String(b);
            }

            var loginCommand = string.Format("{{\"id\":1,\"method\":\"login\",\"api_key\":\"{0}\",\"sign\":\"{1}\",\"nonce\":{2:D}}}", apiKey, sign, nonce);
            await Send(socket, loginCommand);

            var subCommand = "{ \"id\":2,\"method\":\"subscribe\",\"topics\":[\"spot/orders\",\"spot/user_trades\"]}";
            await Send(socket, subCommand);
            return;
        } 

        static async Task Send(ClientWebSocket socket, string data) =>
            await socket.SendAsync(Encoding.UTF8.GetBytes(data), WebSocketMessageType.Text, true, CancellationToken.None);

        static async Task Receive(ClientWebSocket socket)
        {
            var buffer = new ArraySegment<byte>(new byte[2048]);
            do
            {
                WebSocketReceiveResult result;
                using (var ms = new MemoryStream())
                {
                    do
                    {
                        result = await socket.ReceiveAsync(buffer, CancellationToken.None);
                        ms.Write(buffer.Array, buffer.Offset, result.Count);
                    } while (!result.EndOfMessage);

                    if (result.MessageType == WebSocketMessageType.Close)
                        break;

                    ms.Seek(0, SeekOrigin.Begin);
                    using (var reader = new StreamReader(ms, Encoding.UTF8))
                        Console.WriteLine(await reader.ReadToEndAsync());
                }
            } while (true);
        }
    }
}
