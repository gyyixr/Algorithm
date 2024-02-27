import java.util.HashMap;
import java.util.Map;

public class CheckInclusion {

    public static boolean checkInclusion(String s1, String s2) {
        int m = s1.length(), n = s2.length();
        if (m > n) return false;
        int[] cnt = new int[26];
        for (char c : s1.toCharArray()) cnt[c - 'a']++;
        int[] cur = new int[26];
        for (int i = 0; i < m; i++) cur[s2.charAt(i) - 'a']++;
        if (check(cnt, cur)) return true;
        for (int i = m; i < n; i++) {
            cur[s2.charAt(i) - 'a']++;
            cur[s2.charAt(i - m) - 'a']--;
            if (check(cnt, cur)) return true;
        }
        return false;
    }

    public static boolean check(int[] cnt1, int[] cnt2) {
        for (int i = 0; i < 26; i++) {
            if (cnt1[i] != cnt2[i]) return false;
        }
        return true;
    }

    public static void main(String[] args) {
        String s1 = "hello";
        String s2 = "ooolleoooleh";
        System.out.println(checkInclusion(s1, s2));
    }
}
