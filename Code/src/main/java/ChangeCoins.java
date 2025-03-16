import java.util.*;

public class ChangeCoins {
  public static long[] kthPalindrome(int[] queries, int intLength) {
    long[] res = new long[queries.length];
    long max = (long) (Math.pow(10, intLength) - 1);
    long min = (long) Math.pow(10, intLength - 1);
    Map<Long, Long> map = new HashMap<>();
    int count = 0;
    for (long j = min; j <= max; j++) {
      if (isPalindrome(j)) {
        count++;
        map.put((long) count, j);
      }
    }
    for (int i = 0; i < queries.length; i++) {
      res[i] = map.get(queries[i]);
    }

    return res;
  }

  static boolean isPalindrome(long num) {
    String s = new String(String.valueOf(num));
    int i = 0;
    int j = s.length() - 1;
    while (i < j) {
      if (s.charAt(i) != s.charAt(j)) {
        return false;
      }
      i++;
      j--;
    }
    return true;
  }

  public static void main(String[] args) {
    int[] queries = new int[] {1, 2, 3, 4, 5, 90};
    System.out.println(Arrays.toString(kthPalindrome(queries, 3)));
  }
}
