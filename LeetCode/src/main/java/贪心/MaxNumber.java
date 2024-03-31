package 贪心;

import java.util.Arrays;

/**
 * 数组 A 中给定可以使用的 1~9 的数，返回由数组 A 中的元素组成的小于等于n的最大数。
 * 示例 1：A={1, 2, 9, 4}，n=2533，返回 2499。
 * 示例 2：A={1, 2, 5, 4}，n=2543，返回 2542。
 * 示例 3：A={1, 2, 5, 4}，n=2541，返回 2525。
 * 示例 4：A={1, 2, 9, 4}，n=2111，返回 1999。
 * 示例 5：A={5, 9}，n=5555，返回 999。
 */
public class MaxNumber {
    public static void main(String[] args) {
        int num = 2533;
        int[] input = new int[]{1,2,4,9};
        System.out.println(getMax(num, input));
    }

    public static int max = Integer.MIN_VALUE;
    public static int len;
    public static int targetNum;
    public static int[] nums;

    public static int getMax(int target, int[] array) {
        //记得先排序！
        Arrays.sort(array);
        int length = (target + "").length();
        targetNum = target;
        len = length;
        nums = array;
        dfs(0, 0,0);
        return max;
    }

    private static void dfs(int cur, int curLen, int startIndex) {
        max = Math.max(max, cur);
        if (curLen > len) return;//剪枝
        for (int i = startIndex; i < nums.length; i++) {//startIndex也是为了剪枝
            int tem = cur * 10 + nums[i];
            if (tem > targetNum) {//最开始排序是为了在这里剪枝
                continue;
            }
            cur = cur * 10 + nums[i];
            curLen += 1;
            dfs(cur, curLen, i);
            cur /= 10;// cur = (cur - num)/10
            curLen -= 1;
        }
    }
}
