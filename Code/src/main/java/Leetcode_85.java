import java.util.LinkedList;
import java.util.Stack;

public class Leetcode_85 {

    // 主函数，计算最大矩形面积
    public static int maximalRectangle(char[][] matrix) {
        // 边界检查
        if (matrix == null || matrix.length == 0) {
            return 0;
        }

        int rows = matrix.length;
        int cols = matrix[0].length;
        int[] heights = new int[cols];  // 辅助数组，记录每个位置上方连续1的高度
        int maxArea = 0;

        // 遍历每一行
        for (int i = 0; i < rows; i++) {
            // 更新辅助数组
            for (int j = 0; j < cols; j++) {
                if (matrix[i][j] == '1') {
                    heights[j]++;
                } else {
                    heights[j] = 0;
                }
            }

            // 计算以当前行为底的最大矩形面积
            maxArea = Math.max(maxArea, largestRectangleArea(heights));
        }

        return maxArea;
    }

    // 计算柱形图的最大矩形面积
    public static int largestRectangleArea(int[] heights) {
       // Stack<Integer> stack = new Stack<>();  // 栈用于存储柱形的索引
        LinkedList<Integer> stack = new LinkedList<>();
        int maxArea = 0;

        // 遍历每个柱形的高度
        for (int i = 0; i < heights.length; i++) {
            // 如果当前高度小于栈顶柱形的高度，出栈并计算面积
            while (!stack.isEmpty() && heights[i] < heights[stack.peek()]) {
                int height = heights[stack.pop()];
                int width = stack.isEmpty() ? i : i - stack.peek() - 1;
                maxArea = Math.max(maxArea, height * width);
            }
            // 将当前柱形的索引入栈
            stack.push(i);
        }

        // 处理栈中剩余的柱形
        while (!stack.isEmpty()) {
            int height = heights[stack.pop()];
            int width = stack.isEmpty() ? heights.length : heights.length - stack.peek() - 1;
            maxArea = Math.max(maxArea, height * width);
        }

        return maxArea;
    }

    // 主函数，用于测试
    public static void main(String[] args) {
        // 示例用法
        char[][] matrix = {
                {'1', '0', '1', '0', '0'},
                {'1', '0', '1', '1', '1'},
                {'1', '1', '1', '1', '1'},
                {'1', '0', '0', '1', '0'}
        };

        int result = maximalRectangle(matrix);
        System.out.println("Maximal Rectangle Area: " + result);
    }
}
