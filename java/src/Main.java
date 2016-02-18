import java.util.HashMap;

public class Main {

    public static void main(String[] args) {

        Exmo e = new Exmo("your_key","your_secret");
        String result = e.Request("user_info", null);
        System.out.println(result);
        String result2 = e.Request("user_cancelled_orders", new HashMap<String, String>() {{
            put("limit", "2");
            put("offset", "0");
        }});
        System.out.println(result2);
    }
}
