import java.util.Arrays;

public class Demo {
  public static void main(String[] args) {
    int[] arr = new int[1];
    arr[0] = 0;
    System.out.println(Arrays.toString(arr));
    changeArr(arr);
    System.out.println(Arrays.toString(arr));
  }

  public static void changeArr(int[] arr) {
    arr[0] = 1;
  }
}
