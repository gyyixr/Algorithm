class ListNode {
    public int val;
    public ListNode next;

    public ListNode(int value) {
        this.val = value;
    }

    public ListNode() {
    }

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
     *
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
     *
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

    public ListNode mergeKLists(ListNode[] lists) {
        if (lists.length == 0) {
            return null;
        }
        return merge(lists, 0, lists.length - 1);
    }

    private ListNode merge(ListNode[] lists, int lo, int hi) {
        if (lo == hi) {
            return lists[lo];
        }
        int mid = lo + (hi - lo) / 2;
        ListNode l1 = merge(lists, lo, mid);
        ListNode l2 = merge(lists, mid + 1, hi);
        return merge2Lists(l1, l2);
    }

    private ListNode merge2Lists(ListNode l1, ListNode l2) {
        ListNode dummyHead = new ListNode(0);
        ListNode curl = dummyHead;
        while (l1 != null && l2 != null) {
            if (l1.val < l2.val) {
                curl.next = l1;
                l1 = l1.next;
            } else {
                curl.next = l2;
                l2 = l2.next;
            }
            curl = curl.next;
        }

        curl.next = l1 == null ? l2 : l1;

        return dummyHead.next;
    }

    public static void main(String[] args) {
        ListNode a = new ListNode(1);
        a.next = new ListNode(3);

        ListNode b = new ListNode(2);
        b.next = new ListNode(4);
      MergeTwoList mergeKLists = new MergeTwoList();
        ListNode[] input = new ListNode[]{a, b,};

        ListNode.printNodeList(mergeKLists.mergeKLists(input));
    }
}
