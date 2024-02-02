import java.util.Arrays;

public class MergeSort {

  public static void main(String[] args) {
    int[] arr = {3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5};
    System.out.println("Original array: " + Arrays.toString(arr));

    int[] sortedArr = mergeSort(arr);

    System.out.println("Sorted array: " + Arrays.toString(sortedArr));
  }

  public static int[] mergeSort(int[] arr) {
    if (arr.length <= 1) {
      return arr;
    }

    int mid = arr.length / 2;
    int[] left = Arrays.copyOfRange(arr, 0, mid);
    int[] right = Arrays.copyOfRange(arr, mid, arr.length);

    return merge(mergeSort(left), mergeSort(right));
  }

  public static int[] merge(int[] left, int[] right) {
    int[] result = new int[left.length + right.length];
    int i = 0, j = 0, k = 0;

    while (i < left.length && j < right.length) {
      if (left[i] <= right[j]) {
        result[k++] = left[i++];
      } else {
        result[k++] = right[j++];
      }
    }

    // 处理剩余元素
    while (i < left.length) {
      result[k++] = left[i++];
    }
    while (j < right.length) {
      result[k++] = right[j++];
    }

    return result;
  }
}
