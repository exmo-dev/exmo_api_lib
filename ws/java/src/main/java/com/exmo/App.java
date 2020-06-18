package com.exmo;

import java.net.URI;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

import org.java_websocket.client.WebSocketClient;
import org.java_websocket.handshake.ServerHandshake;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class App {
    private static final String HMAC_SHA512 = "HmacSHA512";

    public static void main(String[] args) {
        publicApiUsageExample();
    }

    public static void publicApiUsageExample() {
        WebSocketClient ws = null;
        try {
            ws = startExmoClient("wss://ws-api.exmo.com:443/v1/public", new String[]{"{\"id\":1,\"method\":\"subscribe\",\"topics\":[\"spot/trades:BTC_USD\",\"spot/ticker:LTC_USD\"]}"});
            // wait for pressing any key
            System.in.read();
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            if (ws != null) {
                try {
                    ws.close();
                } catch (Exception e) {

                }
            }
        }
    }

    public static void privateApiUsageExample() {
        WebSocketClient ws = null;
        try {
            final String apiKey = "";
            final String secretKey = "";
            final long nonce = System.currentTimeMillis();

            SecretKeySpec keySpec = new SecretKeySpec(secretKey.getBytes(StandardCharsets.UTF_8), HMAC_SHA512);
            Mac mac = Mac.getInstance(HMAC_SHA512);
            mac.init(keySpec);
            byte[] macData = mac.doFinal((apiKey + nonce).getBytes(StandardCharsets.UTF_8));
            String sign = Base64.getEncoder().encodeToString(macData);

            String loginCommand = String.format("{\"id\":1,\"method\":\"login\",\"api_key\":\"%s\",\"sign\":\"%s\",\"nonce\":%d}", apiKey, sign, nonce);

            ws = startExmoClient("wss://ws-api.exmo.com:443/v1/private", new String[]{
                    loginCommand,
                    "{\"id\":2,\"method\":\"subscribe\",\"topics\":[\"spot/orders\",\"spot/user_trades\"]}"
            });

            // wait for pressing any key
            System.in.read();
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            if (ws != null) {
                try {
                    ws.close();
                } catch (Exception e) {

                }
            }
        }
    }

    public static WebSocketClient startExmoClient(String url, String[] initMessages) throws URISyntaxException, InterruptedException {
        WebSocketClient ws = newWsClient(url);
        ws.connectBlocking();
        for (String message : initMessages) {
            ws.send(message);
            System.out.println("sent: " + message);
        }
        return ws;
    }

    private static WebSocketClient newWsClient(String url) throws URISyntaxException {
        return new WebSocketClient(new URI(url)) {

            @Override
            public void onOpen(ServerHandshake serverHandshake) {
            }

            @Override
            public void onMessage(String s) {
                System.out.println("read: " + s);
            }

            @Override
            public void onClose(int i, String s, boolean b) {
                System.out.println("close: " + i + " " + s);
            }

            @Override
            public void onError(Exception e) {
                System.out.println("error: " + e.toString());
            }
        };
    }
}
