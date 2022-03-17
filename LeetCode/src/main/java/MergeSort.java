import java.util.Arrays;

public class MergeSort {
  /**
   * 将arr[l...mid] 和arr[mid+1....r] 两部分进行归并
   * @param arr
   * @param l
   * @param mid
   * @param r
   */
  private void merge(Comparable[] arr, int l, int mid, int r) {
    // 复制arr的[l,r+1)区间，左闭右开，其实就是[l,r]
    Comparable[] aux = Arrays.copyOfRange(arr, l, r + 1);
    // 初始化，i指向左半部分的起始；j指向右半部分其实索引位置mid+1
    int i = l, j = mid + 1;
    for (int k = l; k <= r; k++) {
      //
      if (i > mid) {
        // 左半部分元素已经全部处理完毕
        arr[k] = aux[j - l];
        j++;
      } else if (j > r) {
        // 右半部分元素已经全部处理完毕
        arr[k] = aux[i - l];
        i++;
      } else if (aux[i - l].compareTo(aux[j - l]) < 0) {
        // 左半部分所指元素<右半部分所指元素
        arr[k] = aux[i - l];
        i++;
      } else {
        arr[k] = aux[j - l];
        j++;
      }
    }
  }

  /**
   * 对arr的[l,r]范围排序
   *
   * @param arr
   * @param l
   * @param r
   */
  public void sort(Comparable[] arr, int l, int r) {
    if (l >= r) return;
    int mid = (r + l) / 2;
    sort(arr, l, mid);
    sort(arr, mid + 1, r);
    merge(arr, l, mid, r);
  }

  public void sort(Comparable[] arr) {
    int n = arr.length;
    sort(arr, 0, n - 1);
  }

  public static void main(String[] args) {
    //
    MergeSort mergeSort = new MergeSort();
    Integer[] arr = {2, 1, 11, 5, 7, 3, 10, 7, 8, 3,0};
    System.out.println(Arrays.toString(arr));
    mergeSort.sort(arr);
    System.out.println(Arrays.toString(arr));

  }
}
