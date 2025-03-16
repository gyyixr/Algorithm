package 二叉树;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.Stack;

class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    TreeNode(int x) {
        val = x;
    }
}

//https://www.codenong.com/cs106522343/
public class BinaryTreeBottomUpLevelOrderTraversal {
    public static void main(String[] args) {
        // 创建一个示例二叉树
        //       1
        //      / \
        //     2   3
        //    / \
        //   4   5
        TreeNode root = new TreeNode(1);
        root.left = new TreeNode(2);
        root.right = new TreeNode(3);
        root.left.left = new TreeNode(4);
        root.left.right = new TreeNode(5);

//        // 调用层次遍历方法
//        List<List<Integer>> result = bottomUpLevelOrder(root);
//        // 打印结果
//        for (List<Integer> level : result) {
//            System.out.println(level);
//        }

        List<Integer> integers = levelOrderBottom2Top(root);
        for(Integer item: integers) {
            System.out.println(item);
        }
    }

    // 从底到上的层次遍历二叉树并返回结果
    public static List<List<Integer>> bottomUpLevelOrder(TreeNode root) {
        List<List<Integer>> result = new ArrayList<>();
        if (root == null) return result;

        Stack<TreeNode> currentLevel = new Stack<>();
        Stack<TreeNode> nextLevel = new Stack<>();
        currentLevel.push(root);

        while (!currentLevel.isEmpty()) {
            int levelSize = currentLevel.size();
            List<Integer> currentList = new ArrayList<>();
            for (int i = 0; i < levelSize; i++) {
                TreeNode node = currentLevel.pop();
                currentList.add(node.val);
                if (node.right != null) {
                    nextLevel.push(node.right);
                }
                if (node.left != null) {
                    nextLevel.push(node.left);
                }
            }
            result.add(currentList);
            while (!nextLevel.isEmpty()) {
                currentLevel.push(nextLevel.pop());
            }
        }

        return result;
    }

    public static List<Integer> levelOrderBottom2Top(TreeNode root) {
        LinkedList<TreeNode> stack = new LinkedList<>();
        LinkedList<Integer> output = new LinkedList<>();
        if (root == null) {
            return output;
        }

        stack.add(root);
        while (!stack.isEmpty()) {
            TreeNode node = stack.pollFirst();
            //之前是pollLast()，这里从栈变成队列了。
            output.addFirst(node.val);
            if (node.right != null) {
                stack.add(node.right);
            }
            if (node.left != null) {
                stack.add(node.left);
            }
            //栈变成队列之后，想要保证先左后右的顺序，自然要把左右节点代码交换
        }
        return output;
    }
}
