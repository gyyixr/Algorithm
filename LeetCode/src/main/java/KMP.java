public class KMP {
    public static boolean kmp(String text, String pattern) {
        int n = text.length();
        int m = pattern.length();
        if (n < m) { // 特判，文本串长度小于模式串长度，直接返回false
            return false;
        }
        // 计算next数组
        int[] next = getNext(pattern);
        // 从前往后匹配
        int i = 0, j = 0;
        while (i < n && j < m) {
            if (j == -1 || text.charAt(i) == pattern.charAt(j)) {
                i++;
                j++;
            } else {
                j = next[j]; // 跳过一部分不必要的比较操作
            }
        }
        return j == m; // 如果j等于m，说明匹配成功，返回true，否则返回false
    }

    // 计算next数组
    private static int[] getNext(String pattern) {
        int m = pattern.length();
        int[] next = new int[m];
        int i = 0, j = -1;
        next[0] = -1;
        while (i < m - 1) {
            if (j == -1 || pattern.charAt(i) == pattern.charAt(j)) {
                i++;
                j++;
                next[i] = j;
            } else {
                j = next[j];
            }
        }
        return next;
    }
}
