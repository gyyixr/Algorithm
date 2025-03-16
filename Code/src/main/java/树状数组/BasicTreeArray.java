package 树状数组;


public class BasicTreeArray {
    static int maxLength = 100;
    static int[] tree = new int[maxLength];

    //求某一个数的二进制表示中最低的一位1所表示的数，比如6（110）返回2（10）
    static int lowBit(int x) {
        return x & (-x);
    }

    static void add(int i, int value) {
        while (i <= maxLength) {
            //单点更新的执行过程是从左向右，每步的执行逻辑是i += lowBit(i),同时更新端点的值。
            tree[i] += value;
            i += lowBit(i);
        }
    }

    static int getSum(int i) {
        int sum = 0;
        while (i > 0) {
            //区间查询的执行过程是从右向左，每步的执行逻辑是 i -= lowBit(i)，同时累加端点的值
            sum += tree[i];
            i -= lowBit(i);
        }
        return sum;
    }

    public static void main(String[] args) {
        for (int i = 1; i <= 6; i++) {
            add(i, i);
        }
        System.out.println(getSum(3));
        //输入：1 2 3 4 5 6
        //输出：6
    }
}
