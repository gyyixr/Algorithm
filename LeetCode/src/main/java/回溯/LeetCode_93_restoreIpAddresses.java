package 回溯;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;

/**
 * 有效 IP 地址 正好由四个整数（每个整数位于 0 到 255 之间组成，且不能含有前导 0），整数之间用 '.' 分隔。
 * <p>
 * 例如："0.1.2.201" 和 "192.168.1.1" 是 有效 IP 地址，但是 "0.011.255.245"、"192.168.1.312" 和 "192.168@1.1" 是 无效 IP 地址。
 * 给定一个只包含数字的字符串 s ，用以表示一个 IP 地址，返回所有可能的有效 IP 地址，这些地址可以通过在 s 中插入 '.' 来形成。你 不能 重新排序或删除 s 中的任何数字。你可以按 任何 顺序返回答案。
 * <p>
 * 示例 1：
 * 输入：s = "25525511135"
 * 输出：["255.255.11.135","255.255.111.35"]
 * 示例 2：
 * 输入：s = "0000"
 * 输出：["0.0.0.0"]
 * 示例 3：
 * 输入：s = "101023"
 * 输出：["1.0.10.23","1.0.102.3","10.1.0.23","10.10.2.3","101.0.2.3"]
 */
public class LeetCode_93_restoreIpAddresses {
    List<String> res = new LinkedList<>();
    LinkedList<String> segment = new LinkedList<>();

    public List<String> restoreIpAddresses(String s) {
        helper(s, 0);
        return res;
    }

    void helper(String s, int start) {
        if (start == s.length() && segment.size() == 4) {
            StringBuilder t = new StringBuilder();
            for (String se : segment) t.append(se).append(".");
            t.deleteCharAt(t.length() - 1);
            res.add(t.toString());
            return;
        }
        if (start < s.length() && segment.size() == 4) return;
        for (int i = 1; i <= 3; i++) {
            if (start + i - 1 >= s.length()) return;
            if (i != 1 && s.charAt(start) == '0') return;
            String str = s.substring(start, start + i);
            if (i == 3 && Integer.parseInt(str) > 255) return;
            segment.addLast(str);
            helper(s, start + i);
            segment.removeLast();
        }
    }

    //====================================================================================
    List<String> result = new ArrayList<>();

    public List<String> restoreIpAddressesII(String s) {
        StringBuilder sb = new StringBuilder(s);
        backTracking(sb, 0, 0);
        return result;
    }

    private void backTracking(StringBuilder s, int startIndex, int dotCount) {
        if (dotCount == 3) {
            if (isValid(s, startIndex, s.length() - 1)) {
                result.add(s.toString());
            }
            return;
        }
        for (int i = startIndex; i < s.length(); i++) {
            if (isValid(s, startIndex, i)) {
                s.insert(i + 1, '.');
                //这里+2是因为当前的segment已经判断符合条件，后面要跟一个.
                backTracking(s, i + 2, dotCount + 1);
                s.deleteCharAt(i + 1);
            } else {
                break;
            }
        }
    }

    //[start, end]
    private boolean isValid(StringBuilder s, int start, int end) {
        if (start > end) return false;
        //注意前导0的判断,IP只能有0.0.0.0这一种
        if (s.charAt(start) == '0' && start != end) return false;
        return Integer.parseInt(s.substring(start, end + 1)) <= 255;
    }


    public static void main(String[] args) {
        LeetCode_93_restoreIpAddresses leetCode_93_restoreIpAddresses = new LeetCode_93_restoreIpAddresses();
        System.out.println(leetCode_93_restoreIpAddresses.restoreIpAddresses("101023"));
        System.out.println(leetCode_93_restoreIpAddresses.restoreIpAddressesII("0000"));
    }

}

/**
 * class Solution {
 * List<String> result = new ArrayList<>();
 * <p>
 * public List<String> restoreIpAddresses(String s) {
 * if (s.length() > 12) return result; // 算是剪枝了
 * backTrack(s, 0, 0);
 * return result;
 * }
 * <p>
 * // startIndex: 搜索的起始位置， pointNum:添加逗点的数量
 * private void backTrack(String s, int startIndex, int pointNum) {
 * if (pointNum == 3) {// 逗点数量为3时，分隔结束
 * // 判断第四段⼦字符串是否合法，如果合法就放进result中
 * if (isValid(s,startIndex,s.length()-1)) {
 * result.add(s);
 * }
 * return;
 * }
 * for (int i = startIndex; i < s.length(); i++) {
 * if (isValid(s, startIndex, i)) {
 * s = s.substring(0, i + 1) + "." + s.substring(i + 1);    //在str的后⾯插⼊⼀个逗点
 * pointNum++;
 * backTrack(s, i + 2, pointNum);// 插⼊逗点之后下⼀个⼦串的起始位置为i+2
 * pointNum--;// 回溯
 * s = s.substring(0, i + 1) + s.substring(i + 2);// 回溯删掉逗点
 * } else {
 * break;
 * }
 * }
 * }
 * <p>
 * // 判断字符串s在左闭⼜闭区间[start, end]所组成的数字是否合法
 * private Boolean isValid(String s, int start, int end) {
 * if (start > end) {
 * return false;
 * }
 * if (s.charAt(start) == '0' && start != end) { // 0开头的数字不合法
 * return false;
 * }
 * int num = 0;
 * for (int i = start; i <= end; i++) {
 * if (s.charAt(i) > '9' || s.charAt(i) < '0') { // 遇到⾮数字字符不合法
 * return false;
 * }
 * num = num * 10 + (s.charAt(i) - '0');
 * if (num > 255) { // 如果⼤于255了不合法
 * return false;
 * }
 * }
 * return true;
 * }
 * }
 */
