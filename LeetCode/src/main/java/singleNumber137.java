import java.util.HashMap;
import java.util.Map;

/**
 * 给你一个整数数组nums ，除某个元素仅出现 一次 外，其余每个元素都恰出现 三次 。请你找出并返回那个只出现了一次的元素。
 *
 * <p>
 *
 * <p>示例 1：
 *
 * <p>输入：nums = [2,2,3,2] 输出：3 示例 2：
 *
 * <p>输入：nums = [0,1,0,1,0,1,99] 输出：99
 *
 * <p>来源：力扣（LeetCode） 链接：https://leetcode-cn.com/problems/single-number-ii
 * 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
public class singleNumber137 {
    /**
     * 哈希表法
     * @param nums
     * @return
     */
    public int singleNumber(int[] nums) {
        Map<Integer, Integer> map = new HashMap<>();
        for (int x : nums) {
            map.put(x, map.getOrDefault(x, 0) + 1);
        }
        for (int x : map.keySet()) {
            if (map.get(x) == 1) return x;
        }
        return -1;
    }

    /**
     * 位数统计法
     * @param nums
     * @return
     */
    public int singleNumberBitCount(int[] nums) {
        int[] cnt = new int[32];
        for (int x : nums) {
            for (int i = 0; i < 32; i++) {
                if (((x >> i) & 1) == 1) {
                    //统计这一个位上所有的数字贡献的1的总和
                    cnt[i]++;
                }
            }
        }
        int ans = 0;
        for (int i = 0; i < 32; i++) {
            if ((cnt[i] % 3 & 1) == 1) {
                //因为数字都出现三次，所以每个位上出现3次的被忽略，出现一次的就是被查找到的数字
                ans += (1 << i);
            }
        }
        return ans;
    }



}
