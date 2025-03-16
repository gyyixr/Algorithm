import java.util.Stack;

public class LeetCode_224_calculate2 {
    public int calculate(String s) {
        int ans = 0;
        char[] str = s.toCharArray();
        int len = str.length;
        Stack<Integer> st_num = new Stack<>();
        Stack<Integer> st_opt = new Stack<>();
        int sign = 1;//正负号,运算符号
        for (int i = 0; i < len; i++) {
            if (str[i] == ' ') continue;
            if (str[i] == '+' || str[i] == '-') {
                sign = str[i] == '+' ? 1 : -1;
            } else if (Character.isDigit(str[i])) {//数字
                int num = str[i] - '0';
                while (i < len - 1 && Character.isDigit(str[i + 1])) {//将这个数字找完
                    num = num * 10 + (str[++i] - '0');
                }
                ans += sign * num;
            } else if (str[i] == '(') {//左括号，暂存结果
                st_num.push(ans);
                st_opt.push(sign);
                ans = 0;
                sign = 1;
            } else ans = st_num.pop() + ans * st_opt.pop();//右括号更新结果
        }
        return ans;
    }

    public static void main(String[] args) {
        LeetCode_224_calculate2 leetCode_224_calculate2 = new LeetCode_224_calculate2();
        System.out.println(leetCode_224_calculate2.calculate("(1+(4+5+2)-3)+(6+8)"));
    }
}
