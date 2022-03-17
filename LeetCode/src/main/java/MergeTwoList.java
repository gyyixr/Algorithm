class ListNode {
  public int val;
  public ListNode next;

  public ListNode(int value) {
    this.val = value;
  }

  public ListNode() {}

  public static void printNodeList(ListNode head) {
    if (head.next == null) {
      System.out.println("head: " + head.val);
    }
    System.out.print("head: ");
    while (head.next != null) {
      System.out.print(head.val + " -> ");
      head = head.next;
    }
    System.out.print(head.val);
    System.out.println(" \n");
  }
}

public class MergeTwoList {
  /**
   * 迭代法
   * @param list1
   * @param list2
   * @return
   */
  public ListNode mergeTwoLists2(ListNode list1, ListNode list2) {
    ListNode preHead = new ListNode(-1);
    ListNode pre = preHead;
    while (list1 != null && list2 != null) {
      if (list1.val < list2.val) {
        pre.next = list1;
        list1 = list1.next;
      } else {
        pre.next = list2;
        list2 = list2.next;
      }
      pre = pre.next;
    }
    pre.next = list1 == null ? list2 : list1;
    return preHead.next;
  }

  /**
   * 递归法
   * @param list1
   * @param list2
   * @return
   */
  public ListNode mergeTwoLists(ListNode list1, ListNode list2) {
    if (list1 == null) {
      return list2;
    }
    if (list2 == null) {
      return list1;
    }
    if (list1.val < list2.val) {
      list1.next = mergeTwoLists(list1.next, list2);
      return list1;
    } else {
      list2.next = mergeTwoLists(list1, list2.next);
      return list2;
    }
  }
}
