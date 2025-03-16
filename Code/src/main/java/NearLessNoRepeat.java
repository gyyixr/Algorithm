import java.util.*;

public class NearLessNoRepeat {
    public int[][] getNearLessNoRepeat(int[] arr) {
        int[][] ans = new int[arr.length][2];
        Stack<Integer> stack = new Stack<>();
        // 遍历数组，入栈
        for (int i = 0; i < arr.length; ++i) {
            while (!stack.isEmpty() && arr[i] < arr[stack.peek()]) {
                int popIndex = stack.pop();
                int leftLessIndex = stack.isEmpty() ? -1 : stack.peek();
                ans[popIndex][0] = leftLessIndex;
                ans[popIndex][1] = i;
            }
            stack.push(i);
        }

        while (!stack.isEmpty()) {
            int popIndex = stack.pop();
            int leftLessIndex = stack.isEmpty() ? -1 : stack.peek();
            ans[popIndex][0] = leftLessIndex;
            // 说明该索引右边没有比当前小的元素，有的话该索引在上边循环就弹出了
            ans[popIndex][1] = -1;
        }
        return ans;
    }

    public static void main(String[] args) {
        int[] input = new int[]{1,2,3,4,5,6,7};
        int[] ints = Arrays.copyOfRange(input, 0, input.length);
        System.out.println(Arrays.toString(ints));
        NearLessNoRepeat nearLessNoRepeat = new NearLessNoRepeat();
        int[][] nearLessNoRepeat1 = nearLessNoRepeat.getNearLessNoRepeat(input);
        for(int i = 0; i < nearLessNoRepeat1.length; i++) {
            System.out.println(Arrays.toString(nearLessNoRepeat1[i]));
        }



    }
}
