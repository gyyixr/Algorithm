import java.util.Arrays;

public class QuickSort {
  public static void sort(int[] a) {
    if (a.length > 0) {
      sort(a, 0, a.length - 1);
    }
  }

  public static void sort(int[] a, int low, int height) {
    int i = low;
    int j = height;
    if (i > j) { // 放在k之前，防止下标越界，这里要严格大于，因为i和j一定要会和，即i=j会发生
      return;
    }
    int k = a[i];

    while (i < j) {//这里要<
      while (i < j && a[j] >= k) { // 找出小的数,这里要 >=
        j--;
      }
      while (i < j && a[i] <= k) { // 找出大的数,这里要 <=
        i++;
      }
      if (i < j) { // 交换
        int swap = a[i];
        a[i] = a[j];
        a[j] = swap;
      }
    }
    /** k = a[i]; a[i] = a[low]; a[low] = k; * */
    // 交换a[low]和a[i]/a[j]
    a[low] = a[i];
    a[i] = k;

    // 对左边进行排序,递归算法
    sort(a, low, i - 1);
    // 对右边进行排序
    sort(a, i + 1, height);
  }

  public static void main(String[] args) {
    int[] arr = new int[] {5,5, 9, 7, 4, 3, 7, 6, 1, 9, 9, 7, 4};
    System.out.println(Arrays.toString(arr));
    sort(arr);
    System.out.println(Arrays.toString(arr));
  }
}
