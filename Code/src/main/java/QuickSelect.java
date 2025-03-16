public class QuickSelect {
    public static int quickSelect(int[] nums, int left, int right, int k) {
        if (left == right) {
            // 如果数组只有一个元素，直接返回这个元素
            return nums[left];
        }

        int pivotIndex = partition(nums, left, right);

        if (pivotIndex == k - 1) {
            // 如果枢轴的位置正好是K-1，返回枢轴元素
            return nums[pivotIndex];
        } else if (pivotIndex < k - 1) {
            // 如果枢轴位置小于K-1，继续在右边分区查找
            return quickSelect(nums, pivotIndex + 1, right, k);
        } else {
            // 如果枢轴位置大于K-1，继续在左边分区查找
            return quickSelect(nums, left, pivotIndex - 1, k);
        }
    }

    private static int partition(int[] nums, int left, int right) {
        int pivot = nums[right]; // 选择最右边的元素作为枢轴
        int i = left - 1; // 小于枢轴的元素计数

        for (int j = left; j < right; j++) {
            if (nums[j] > pivot) {
                i++; // 找到大于枢轴的元素，交换位置
                swap(nums, i, j);
            }
        }
        swap(nums, i + 1, right); // 将枢轴放到正确的位置
        return i + 1; // 返回枢轴的位置
    }

    private static void swap(int[] nums, int i, int j) {
        int temp = nums[i];
        nums[i] = nums[j];
        nums[j] = temp;
    }

    public static void main(String[] args) {
        int[] nums = {1, 2, 3, 4, 5, 6};
        int k = 6; // 假设我们要找到第2大的元素
        int kthLargest = quickSelect(nums, 0, nums.length - 1, k);
        System.out.println("The " + k + "th largest element is: " + kthLargest);
    }
}