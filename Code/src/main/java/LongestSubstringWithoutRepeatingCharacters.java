import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

public class LongestSubstringWithoutRepeatingCharacters {
    public static int lengthOfLongestSubstring(String s) {
        int maxLen = 0;
        Map<Character, Integer> charIndexMap = new HashMap<>();
        int left = 0; // 窗口的起始位置

        for (int right = 0; right < s.length(); right++) {
            char currentChar = s.charAt(right);
            // 如果当前字符在哈希表中，且索引大于等于左指针位置，移动左指针
            if (charIndexMap.containsKey(currentChar) && charIndexMap.get(currentChar) >= left) {
                left = charIndexMap.get(currentChar) + 1;
            }
            // 更新当前字符的索引
            charIndexMap.put(currentChar, right);
            // 更新最长子串长度
            maxLen = Math.max(maxLen, right - left + 1);
        }

        return maxLen;
    }

    public static void main(String[] args) {
        String s = "abcdabcbb";
        s.contains("");
        System.out.println(lengthOfLongestSubstring(s));
        System.out.println(lengthOfLongestSubstring1(s));// 输出: 3
    }


    public static int lengthOfLongestSubstring1(String s) {
        if (s == null || s.isEmpty()) return 0;
        int max = 1, i = 0, j = 1;
        Set<Character> set = new HashSet<>();
        set.add(s.charAt(i));
        while (j < s.length()) {
            if (set.contains(s.charAt(j)))
            {set.remove(s.charAt(i));
            i++;}
            else {
                set.add(s.charAt(j++));
            }
            max = Math.max(max, j - i);
        }
        return max;
    }

}