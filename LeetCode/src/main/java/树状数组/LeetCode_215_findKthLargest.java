package 树状数组;

import java.util.PriorityQueue;
import java.util.Random;

public class LeetCode_215_findKthLargest {
    //值域映射 + 树状数组 + 二分
    int M = 100010, N = 2 * M;
    int[] tr = new int[N];

    int lowbit(int x) {
        return x & -x;
    }

    int query(int x) {
        int ans = 0;
        for (int i = x; i > 0; i -= lowbit(i)) ans += tr[i];
        return ans;
    }

    void add(int x) {
        for (int i = x; i < N; i += lowbit(i)) tr[i]++;
    }

    public int findKthLargest(int[] nums, int k) {
        for (int x : nums) add(x + M);
        int l = 0, r = N - 1;
        while (l < r) {
            int mid = l + r + 1 >> 1;
            if (query(N - 1) - query(mid - 1) >= k) l = mid;
            else r = mid - 1;
        }
        return r - M;
    }


    //维护一个k个元素的小顶堆 -> 由于找第 K 大元素，其实就是整个数组排序以后后半部分最小的那个元素。
    /**
     * 如果当前堆不满，直接添加；
     * 堆满的时候，如果新读到的数小于等于堆顶，肯定不是我们要找的元素，只有新遍历到的数大于堆顶的时候，才将堆顶拿出，然后放入新读到的数，进而让堆自己去调整内部结构。
     *
   */
    public int findKthLargest2(int[] nums, int k) {
        PriorityQueue<Integer> q = new PriorityQueue<>(k, (a, b)->a-b);
        for (int x : nums) {
            if (q.size() < k || x > q.peek()) q.offer(x);
            if (q.size() > k) q.poll();
        }
        return q.peek();
    }

    //快排的思想，分治
    private final static Random random = new Random(System.currentTimeMillis());

    public int findKthLargest3(int[] nums, int k) {
        // 第 1 大的数，下标是 len - 1;
        // 第 2 大的数，下标是 len - 2;
        // ...
        // 第 k 大的数，下标是 len - k;
        int len = nums.length;
        int target = len - k;

        int left = 0;
        int right = len - 1;

        while (true) {
            int pivotIndex = partition(nums, left, right);
            if (pivotIndex == target) {
                return nums[pivotIndex];
            } else if (pivotIndex < target) {
                left = pivotIndex + 1;
            } else {
                // pivotIndex > target
                right = pivotIndex - 1;
            }
        }
    }

    private int partition(int[] nums, int left, int right) {
        int randomIndex = left + random.nextInt(right - left + 1);
        swap(nums, left, randomIndex);


        // all in nums[left + 1..le) <= pivot;
        // all in nums(ge..right] >= pivot;
        int pivot = nums[left];
        int le = left + 1;
        int ge = right;

        while (true) {
            while (le <= ge && nums[le] < pivot) {
                le++;
            }

            while (le <= ge && nums[ge] > pivot) {
                ge--;
            }

            if (le >= ge) {
                break;
            }
            swap (nums, le, ge);
            le++;
            ge--;
        }

        swap(nums, left, ge);
        return ge;
    }

    private void swap(int[] nums, int index1, int index2) {
        int temp = nums[index1];
        nums[index1] = nums[index2];
        nums[index2] = temp;
    }
}

