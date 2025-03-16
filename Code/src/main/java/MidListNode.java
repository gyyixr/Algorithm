public class MidListNode {
    //偶数个节点时返回靠右的那个
    public static ListNode middleLeftNode(ListNode head) {
        if (head == null) {
            return null;
        }
        ListNode slow = head;
        ListNode fast = head;

        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;
        }
        return slow;
    }

    //偶数个节点时返回靠左的那个
    public static ListNode middleRightNode(ListNode head) {
        if (head == null) {
            return null;
        }
        ListNode slow = head;
        ListNode fast = head;

        while (fast.next != null && fast.next.next != null) {
            slow = slow.next;
            fast = fast.next.next;
        }
        return slow;
    }

    public static void main(String[] args) {
        ListNode root = new ListNode(1);
        root.next = new ListNode(2);
        root.next.next = new ListNode(3);
        root.next.next.next = new ListNode(4);
        System.out.println(middleLeftNode(root).val);
        System.out.println(middleRightNode(root).val);
//        root.next.next.next.next = new ListNode(5);
//        System.out.println(middleLeftNode(root).val);
//        System.out.println(middleRightNode(root).val);
    }
}
