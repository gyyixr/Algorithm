import java.util.LinkedList;
import java.util.List;

public class Btree {
    static List<List<Integer>> result;
    static LinkedList<Integer> path;

    static class BtreeNode {
        public int value;
        public BtreeNode left;
        public BtreeNode right;

        public BtreeNode(int value) {
            this.value = value;
        }
    }

    public static void main(String[] args) {
        int sum = 0;
        // 构造一颗二叉树
        BtreeNode root = new BtreeNode(1);
        root.right = new BtreeNode(3);
        root.left = new BtreeNode(2);
        root.right.left = new BtreeNode(6);
        root.right.right = new BtreeNode(7);
        root.left.left = new BtreeNode(4);
        root.left.right = new BtreeNode(5);

        pathSum(root, 0);
        System.out.println(result);

        for (List<Integer> list : result) {
            for (Integer v : list) {
                sum += v;
            }
        }
        System.out.println("路径和 -> " + sum);
    }


    public static List<List<Integer>> pathSum(BtreeNode root, int count) {
        result = new LinkedList<>();
        path = new LinkedList<>();
        travesal(root, count);
        return result;
    }

    private static void travesal(BtreeNode root, int count) {
        if (root == null) return;
        path.offer(root.value);
        count += root.value;
        if (root.left == null && root.right == null) {
            result.add(new LinkedList<>(path));
        }
        travesal(root.left, count);
        travesal(root.right, count);
        path.removeLast();
    }
}
