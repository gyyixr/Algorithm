package 二叉树;

//给定一个二叉搜索树的根节点 root ，和一个整数 k ，请你设计一个算法查找其中第 k 个最小元素（从 1 开始计数）。
public class KthSmallest {
    public static class TreeNode {
        int val;
        TreeNode left;
        TreeNode right;

        TreeNode() {
        }

        TreeNode(int val) {
            this.val = val;
        }

        TreeNode(int val, TreeNode left, TreeNode right) {
            this.val = val;
            this.left = left;
            this.right = right;
        }
    }

    private int count = 0;

    //BST的中序遍历正好是排好序的
    public TreeNode KthNode(TreeNode pRoot, int k) {
        if (pRoot == null) {
            return null;
        }
        //中序-左子树先遍历-最小；如果递归左子树已经有结果了就不用继续遍历了
        TreeNode nodeLeft = KthNode(pRoot.left, k);
        if (nodeLeft != null) {
            return nodeLeft;
        }

        //中序-中间root节点
        count++;
        if (count == k) {
            return pRoot;
        }

        //中序-右子树先遍历-最大
        TreeNode nodeRight = KthNode(pRoot.right, k);
        return nodeRight;
    }

    public static void main(String[] args) {
        KthSmallest solution = new KthSmallest();
        TreeNode root = new TreeNode(3);
        root.left = new TreeNode(1);
        root.right = new TreeNode(4);
        root.left.right = new TreeNode(2);
        System.out.println(solution.KthNode(root, 3).val);
    }
}
