package DP;

/**
 * 给定一个长度为 n 的 0 索引整数数组 nums。初始位置为 nums[0]。
 * 每个元素 nums[i] 表示从索引 i 向前跳转的最大长度。换句话说，如果你在 nums[i] 处，你可以跳转到任意 nums[i + j] 处:
 * 返回到达 nums[n - 1] 的最小跳跃次数。
 */
public class LeetCode_45_jump2 {
    public int jump(int[] nums) {
        // left 记录当前覆盖的最远距离【下标】
        int left = 0;
        // right 记录下一次覆盖的最远距离【下标】
        int right = 0;
        // 记录跳跃了多少次
        int result = 0;
        for (int i = 0; i <= right && right < nums.length - 1; i++) {
            left = Math.max(left, nums[i] + i);
            // 循环等于右坐标时，跳跃次数+1
            if (i == right) {
                right = left;
                result++;
            }
        }
        return result;
    }
}
