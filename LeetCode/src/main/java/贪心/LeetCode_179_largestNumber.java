package 贪心;

import java.util.Arrays;

/**
 * 给定一组非负整数 nums，重新排列每个数的顺序（每个数不可拆分）使之组成一个最大的整数。
 * 注意：输出结果可能非常大，所以你需要返回一个字符串而不是整数。
 * 示例 1：
 * 输入：nums = [10,2]
 * 输出："210"
 * 示例 2：
 * 输入：nums = [3,30,34,5,9]
 * 输出："9534330"
 */
public class LeetCode_179_largestNumber {
    public String largestNumber(int[] nums) {
        int n = nums.length;
        String[] ss = new String[n];
        //用字符串的比较(ascii码)来代替数字的比较
        for (int i = 0; i < n; i++) ss[i] = "" + nums[i];
        Arrays.sort(ss, (a, b) -> {
            String sa = a + b, sb = b + a ;
            //后者compareTo前者 -> 从大到小，降序
            return sb.compareTo(sa);
        });
        System.out.println(Arrays.toString(ss));

        StringBuilder sb = new StringBuilder();
        for (String s : ss) sb.append(s);
        int len = sb.length();
        int k = 0;
        //处理前导0， [0,0,0]这样的情况输出是"000"，去除前导0后可以得到正确答案"0"
        while (k < len - 1 && sb.charAt(k) == '0') k++;
        return sb.substring(k);
    }

    public static void main(String[] args) {
        LeetCode_179_largestNumber leetCode179LargestNumber = new LeetCode_179_largestNumber();
        System.out.println(leetCode179LargestNumber.largestNumber(new int[]{3,30,34,5,9}));
    }

}
