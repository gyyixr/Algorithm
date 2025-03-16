import java.util.*;
class MonotonicStack {
    public static int[] dailyTemperatures(int[] temperatures) {
        int lens = temperatures.length;
        int []res = new int[lens];
        Deque<Integer> stack = new LinkedList<>();
        for(int i=0;i <lens; i++) {
            while(!stack.isEmpty() && temperatures[i] > temperatures[stack.peek()]){
                res[stack.peek()] = i - stack.peek();
                stack.pop();
            }
            stack.push(i);
        }

        return  res;
    }

    public static void main(String[] args) {
        int[] temperatures = new int[]{73, 72, 71, 70, 69, 68, 67, 66};
        System.out.println(Arrays.toString(dailyTemperatures(temperatures)));
    }



}
