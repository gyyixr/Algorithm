public class MaxAreaOfIsland {
    //  岛屿问题思路
//    void dfs(int[][] grid, int r, int c) {
//        // 判断 base case
//        if (!inArea(grid, r, c)) {
//            return;
//        }
//        // 如果这个格子不是岛屿，直接返回
//        if (grid[r][c] != 1) {
//            return;
//        }
//        grid[r][c] = 2; // 将格子标记为「已遍历过」
//
//        // 访问上、下、左、右四个相邻结点
//        dfs(grid, r - 1, c);
//        dfs(grid, r + 1, c);
//        dfs(grid, r, c - 1);
//        dfs(grid, r, c + 1);
//    }
//
//    // 判断坐标 (r, c) 是否在网格中
//    boolean inArea(int[][] grid, int r, int c) {
//        return 0 <= r && r < grid.length
//                && 0 <= c && c < grid[0].length;
//    }


    //岛屿最大面积（dfs）
    public int maxAreaOfIsland(int[][] grid) {
        int res = 0;
        for (int r = 0; r < grid.length; r++) {
            for (int c = 0; c < grid[0].length; c++) {
                if (grid[r][c] == 1) {
                    int a = area(grid, r, c);
                    res = Math.max(res, a);
                }
            }
        }
        return res;
    }

    int area(int[][] grid, int r, int c) {
        if (!inArea(grid, r, c)) {
            return 0;
        }
        if (grid[r][c] != 1) {
            return 0;
        }
        grid[r][c] = 2;

        return 1
                + area(grid, r - 1, c)
                + area(grid, r + 1, c)
                + area(grid, r, c - 1)
                + area(grid, r, c + 1);
    }

    boolean inArea(int[][] grid, int r, int c) {
        return 0 <= r && r < grid.length
                && 0 <= c && c < grid[0].length;
    }


}
