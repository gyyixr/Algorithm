import java.util.Arrays;
import java.util.List;


class Solution {
    public int monotoneIncreasingDigits(int n) {
        String s = String.valueOf(n);
        char[] chars = s.toCharArray();
        int start = s.length();
        for (int i = s.length() - 2; i >= 0; i--) {
            if (chars[i] > chars[i + 1]) {
                chars[i]--;
                start = i+1;
            }
        }
        for (int i = start; i < s.length(); i++) {
            chars[i] = '9';
        }
        return Integer.parseInt(String.valueOf(chars));
    }
}

public class Demo {
  public static void main(String[] args) {
    System.out.println(new Solution().monotoneIncreasingDigits(4321));
    String s = "123456";
    System.out.println(s.substring(4,5));

      System.out.println(Integer.parseInt(s));
    Arrays.asList();
  }
}
