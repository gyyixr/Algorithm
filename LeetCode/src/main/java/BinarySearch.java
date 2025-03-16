public class BinarySearch {
  public static void main(String[] args) {
    int[] arr = {1, 2, 3, 5, 5, 5, 11, 33, 66, 88, 100};
    System.out.println(getBinarySearch(arr, 65));
  }

  public static int getBinarySearch(int[] arr, int target) {
    if (arr.length <= 0) return -1;
    int n = arr.length;
    int l = -1, r = n;
    while (l + 1 != r) {
      int m = (r + l) / 2;
      if (target >= arr[m]) {
        l = m;
      } else {
        r = m;
      }
    }
    return arr[l];
  }
}
