import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

public class LCP_14_Split_Arrays {
    public int splitArray(int[] nums) {
        // 预置 minPrime
        int[] minPrime = new int[1000000 + 1];
        for (int i = 2; i < minPrime.length; i++) {
            if (minPrime[i] < 2) {
                for (int j = i; j < minPrime.length; j += i) {
                    minPrime[j] = i;
                }
            }
        }

        // 记录执行到位置的次数
        int[] dp = new int[nums.length];
        // 记录因子位置
        Map<Integer,Integer> posMap=new HashMap<Integer,Integer>();

        for (int i = 0; i < nums.length; i++) {
            // 预设次数
            dp[i] = i > 0 ? dp[i - 1] + 1 : 1;

            // 分解 $nums[$i] 的因子;
            int n = nums[i];
            while (n > 1) {
                int factor = minPrime[n];
                int minIndex = i;
                if (posMap.containsKey(factor)) {
                    minIndex = posMap.get(factor);
                }else {
                    posMap.put(factor, minIndex);
                }

                if (minIndex == 0) {
                    dp[i] = 1;
                } else {
                    // 和 已记录处理的位置+1 对比去一个低的
                    if ( ( dp[minIndex - 1] + 1 ) < dp[i] ){
                        dp[i] = dp[minIndex - 1] + 1;
                    }
                }
                // 更新posMap
                if (dp[i] < dp[minIndex]) {
                    posMap.put(factor, i);
                }
                n = n / factor;
            }
        }
        return dp[nums.length - 1];
    }

    public static void main(String[] args) {
        int[] input = new int[]{2,3,3,2,3,3};
        LCP_14_Split_Arrays lcp_14_split_arrays = new LCP_14_Split_Arrays();
        System.out.println(lcp_14_split_arrays.splitArray(input));

    }
}
