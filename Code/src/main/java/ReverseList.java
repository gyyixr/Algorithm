public class ReverseList {
  /**
   * 全部链表反转
   *
   * @param head
   * @return
   */
  public static ListNode reverse(ListNode head) {
    if (head.next == null) return head;
    ListNode last = reverse(head.next);
    head.next.next = head;
    head.next = null;
    return last;
  }

  /** 反转前n个节点 */
  static ListNode successor = null; // 后驱节点

  public static ListNode reverseN(ListNode head, int n) {
    if (n == 1) {
      // 记录第 n + 1 个节点
      successor = head.next;
      return head;
    }
    // 以 head.next 为起点，需要反转前 n - 1 个节点
    ListNode last = reverseN(head.next, n - 1);
    head.next.next = head;
    // 让反转之后的 head 节点和后面的节点连起来
    head.next = successor;
    return last;
  }

  /**
   * 反转[m,n]区间的链表
   *
   * @param head
   * @param m
   * @param n
   * @return
   */
  public static ListNode reverseBetween(ListNode head, int m, int n) {
    // base case
    if (m == 1) {
      return reverseN(head, n);
    }
    // 前进到反转的起点触发 base case
    head.next = reverseBetween(head.next, m - 1, n - 1);
    return head;
  }

  public static void main(String[] args) {
    ListNode node = new ListNode(0);
    node.next = new ListNode(1);
    node.next.next = new ListNode(2);
    node.next.next.next = new ListNode(3);
    node.next.next.next.next = new ListNode(4);
    ListNode.printNodeList(node);
    /*ListNode newHead = reverse(node);
    ListNode.printNodeList(newHead);*/
    // ListNode node1= reverseN(node, 2);
    // ListNode.printNodeList(node1);
    //    ListNode listNode = reverseBetween(node, 2, 4);
    //    ListNode.printNodeList(listNode);
    ReverseKGroup reverseKGroup = new ReverseKGroup();
    ListNode listNode = reverseKGroup.reverseKGroup(node, 2);
    ListNode.printNodeList(listNode);
  }
}
