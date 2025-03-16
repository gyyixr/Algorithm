public class MergeSortedArray {
  /**
   * 给你两个按 非递减顺序 排列的整数数组 nums1 和 nums2，另有两个整数 m 和 n ，分别表示 nums1 和 nums2 中的元素数目。
   * 最终，合并后数组不应由函数返回，而是存储在数组 nums1 中。为了应对这种情况，nums1 的初始长度为 m + n，其中前 m 个元素表示应合并的元素，后 n 个元素为 0
   * ，应忽略。nums2 的长度为 n 。
   *
   * @param nums1
   * @param m
   * @param nums2
   * @param n
   */
  public void merge(int[] nums1, int m, int[] nums2, int n) {
      //从后往前遍历更好
      //从前往后遍历，则比较后选择两者较小的；从后往前遍历则比较后选择两者较大的
        int i = m - 1, j = n - 1, k = m + n - 1;
        while(i >= 0 && j >= 0) {
            if (nums1[i] < nums2[j]) {
                nums1[k] = nums2[j--];
            } else {
                nums1[k] = nums1[i--];
            }
            k--;
        }
        while(i >= 0) nums1[k--] = nums1[i--];//因为结果放在num1里，所以这一行可以不写，从后往前遍历的好处，
        while(j >= 0) nums1[k--] = nums2[j--];
    }
}
