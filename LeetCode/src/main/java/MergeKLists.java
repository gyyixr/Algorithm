public class MergeKLists {
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

        curl.next = l1 == null? l2: l1;

        return dummyHead.next;
    }

    public static void main(String[] args) {
        ListNode a = new ListNode(1);
        a.next = new ListNode(3);

        ListNode b = new ListNode(2);
        b.next = new ListNode(4);
        MergeKLists mergeKLists = new MergeKLists();
        ListNode[] input = new ListNode[]{a,b,};

        ListNode.printNodeList(mergeKLists.mergeKLists(input));

    }
}
