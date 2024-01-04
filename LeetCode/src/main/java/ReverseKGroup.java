public class ReverseKGroup {
  /**
   * 反转以 a 为头结点的链表
   *
   * @param a
   * @return
   */
  public ListNode reverse(ListNode a) {
    ListNode pre, cur, nxt;
    pre = null;
    cur = a;
    nxt = null;
    while (cur != null) {
      nxt = cur.next;
      // 逐个结点反转
      cur.next = pre; // 更新指针位置
      pre = cur;
      cur = nxt;
    }
    return pre;
    // 返回反转后的头结点
  }

  /** 反转区间 [a, b) 的元素，注意是左闭右开 */
  public ListNode reverse(ListNode a, ListNode b) {
    ListNode pre, cur, nxt;
    pre = null;
    cur = a;
    nxt = null;
    // while 终止的条件改一下就行了
    while (cur != b) {
      nxt = cur.next;
      cur.next = pre;
      pre = cur;
      cur = nxt;
    }
    // 返回反转后的头结点
    return pre;
  }

  /**
   *  K个一组反转
   * @param head
   * @param k
   * @return
   */
  public ListNode reverseKGroup(ListNode head, int k) {
    if (head == null) return null;
    // 区间 [a, b) 包含 k 个待反转元素
    ListNode a, b;
    a = b = head;
    for (int i = 0; i < k; i++) {
      // 不足 k 个，不需要反转，base case
      if (b == null) return head;
      b = b.next;
    }
    // 反转前 k 个元素
    ListNode newHead = reverse(a, b);
    // 递归反转后续链表并连接起来
    head.next = reverseKGroup(b, k);
    return newHead;
  }
}
