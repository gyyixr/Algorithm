import java.util.Arrays;

public class GenerateMatrix {
    public static int[][] generateMatrix(int n) {
        int[][] res = new int[n][n];
        int i = 0, j = 0, cur = 2;
        res[0][0] = 1;
        while (cur <= n * n) {
            while (j < n - 1 && res[i][j + 1] == 0) res[i][++j] = cur++; // 右
            while (i < n - 1 && res[i + 1][j] == 0) res[++i][j] = cur++; // 下
            while (j > 0 && res[i][j - 1] == 0) res[i][--j] = cur++; // 左
            while (i > 0 && res[i - 1][j] == 0) res[--i][j] = cur++; // 上
        }
        return res;
    }

    public static void main(String[] args) {
        System.out.println(Arrays.deepToString(generateMatrix(3)));
    }
}
