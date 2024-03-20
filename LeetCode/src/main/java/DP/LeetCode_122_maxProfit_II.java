package DP;
//股票问题2:可以多次买卖
public class LeetCode_122_maxProfit_II {
    //贪心
    public int maxProfit(int[] prices) {
        int result = 0;
        for (int i = 1; i < prices.length; i++) {
            result += Math.max(prices[i] - prices[i - 1], 0);
        }
        return result;
    }

    //DP更通用
    public int maxProfitDp(int[] prices) {
        // [天数][是否持有股票]
        int[][] dp = new int[prices.length][2];

        // base case
        dp[0][0] = 0;
        dp[0][1] = -prices[0];

        for (int i = 1; i < prices.length; i++) {
            // dp公式
            dp[i][0] = Math.max(dp[i - 1][0], dp[i - 1][1] + prices[i]);
            dp[i][1] = Math.max(dp[i - 1][1], dp[i - 1][0] - prices[i]);
        }

        return dp[prices.length - 1][0];
    }

    public static void main(String[] args) {
        int[] input = new int[]{7,1,5,3,6,4};
        LeetCode_122_maxProfit_II solution = new LeetCode_122_maxProfit_II();
        System.out.println(solution.maxProfitDp(input));
        System.out.println(solution.maxProfit(input));
    }
}
