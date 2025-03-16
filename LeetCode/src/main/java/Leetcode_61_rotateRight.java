 class LeetCode_61_rotateRight {
    static class ListNode61 {
        int val;
        ListNode61 next;
        ListNode61 right;
        ListNode61() {}
        ListNode61(int val) { this.val = val; }
    }

    public  ListNode61 rotateRight(ListNode61 head, int k) {
        if (head == null || k == 0) {
            return head;
        }

        ListNode61 tail = head;
        int n = 1;
        while (tail.next != null) {
            tail = tail.next;
            n ++;
        }
        //连成一个环状
        tail.next = head;

        // 具体转多少下要调整一下, 链表只能往左转, 想要往右转就要转n-k+1下, k优化一下起始是k%n下
        for (int i = 0; i < n - k % n - 1; i ++) {
            head = head.next;
        }
        ListNode61 newHead = head.next;
        head.next = null;
        return newHead;

    }

    public static void main(String[] args) {
        LeetCode_61_rotateRight leetcode_61_rotateRight = new LeetCode_61_rotateRight();
        ListNode61 node = new ListNode61(1);
        node.next = new ListNode61(2);
        node.next.next = new ListNode61(3);
        System.out.println(leetcode_61_rotateRight.rotateRight(node, 2000000000).val);

    }

}
