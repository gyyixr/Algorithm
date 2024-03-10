package fenzhi;
/**
 * 实现 pow(x, n) ，即计算 x 的整数 n 次幂函数（即，xn ）。
 * 示例 1：
 * 输入：x = 2.00000, n = 10
 * 输出：1024.00000
 */
public class Leetcode_50_pow {
    public double myPow(double x, int n) {
        if(x == 0.0d) return 0.0d;
        long b = n;
        double res = 1.0;
        if(b < 0) {
            x = 1 / x;
            b = -b;
        }
        while(b > 0) {
            //最后一位是1
            if((b%2) == 1) res *= x;
            x *= x;
            b >>= 1;
        }
        return res;
    }
}
