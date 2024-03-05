import java.util.Stack;

public class LeetCode_227_calculate {
    public int calculate(String s) {
        s = s + "+0";
        char preOperation = '+';
        int num = 0;
        Stack<Integer> stack = new Stack<>();
        for (char c : s.toCharArray()) {
            if (c == ' ') continue;
            if (Character.isDigit(c)) {
                num = num * 10 + (c-'0');
            } else {
                switch (preOperation) {
                    case '+': stack.push(num); break;
                    case '-': stack.push(-1 * num); break;
                    case '*': stack.push(stack.pop() * num); break;
                    case '/': stack.push(stack.pop() / num); break;
                    default:
                }
                preOperation = c;
                num = 0;
            }
        }
        System.out.println(stack);
        return stack.stream().mapToInt(Integer::intValue).sum();
    }

    public static void main(String[] args) {
        LeetCode_227_calculate leetCode_227_calculate = new LeetCode_227_calculate();
        System.out.println(leetCode_227_calculate.calculate("3+5 / 2 "));
    }
}
