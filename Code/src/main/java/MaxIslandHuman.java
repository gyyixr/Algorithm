import java.util.*;

//最大人工岛（最多可以把一个海洋变成陆地，求最大岛屿面积）
class MaxIslandHuman {
    int[] directions = new int[]{-1, 0, 1, 0, -1};
    int n;
    public int largestIsland(int[][] grid) {
        n = grid.length;
        int res = 0, idx = 2;   // 0 和 1都用了，所以从2开始编号

        //将(岛屿索引，岛屿面积)存入map中。
        //即求，编号后，每个岛屿的面积。
        Map<Integer, Integer> area = new HashMap<>();
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == 1) {
                    area.put(idx, calculateArea(idx, grid, i, j));
                    idx++;
                }
            }
        }

        //当遇到一个非岛屿时，则判断上下左右是否连通。
        //使用set去重，防止重复计算。
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                int sum = 0;
                Set<Integer> set = new HashSet<>(); //去重
                if (grid[i][j] == 0) {
                    sum += 1;

                    //上
                    if (i - 1 >= 0 && grid[i - 1][j] != 0 && set.add(grid[i - 1][j])) {
                        sum += area.get(grid[i - 1][j]);
                    }
                    //下
                    if (i + 1 < n && grid[i + 1][j] != 0 && set.add(grid[i + 1][j])) {
                        sum += area.get(grid[i + 1][j]);
                    }
                    //左
                    if (j - 1 >= 0 && grid[i][j - 1] != 0 && set.add(grid[i][j - 1])) {
                        sum += area.get(grid[i][j - 1]);
                    }
                    //右
                    if (j + 1 < n && grid[i][j + 1] != 0 && set.add(grid[i][j + 1])) {
                        sum += area.get(grid[i][j + 1]);
                    }
                }
                res = Math.max(res, sum);
            }
        }
        return res == 0 ? area.get(2) : res;
    }

    //BFS计算岛屿面积
    //将连通岛屿格子值赋值为岛屿编号，防止重复计算。
    public int calculateArea(int idx, int[][] grid, int i, int j) {
        Queue<int[]> q = new LinkedList<>();
        q.offer(new int[]{i, j});
        grid[i][j] = idx;
        int sum = 0;
        while (!q.isEmpty()) {
            int[] cur = q.poll();
            sum++;      //计算岛屿面积

            //for循环判断该岛屿上下左右是否是岛屿。如果是岛屿则入队。
            //directions数组中，存的是方向。
            for (int k = 0; k < directions.length - 1; k++) {
                int x = cur[0] + directions[k];
                int y = cur[1] + directions[k + 1];
                if (x < 0 || x >= n || y < 0 || y >= n || grid[x][y] != 1) continue;    //如果越界或者不是岛屿则跳过
                grid[x][y] = idx;
                q.offer(new int[]{x, y});
            }
        }
        return sum;
    }
}