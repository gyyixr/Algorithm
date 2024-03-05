import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

public class LeetCode_21_generateParenthesis {
    List<String> res = new ArrayList<>();
    // 回溯过程中的路径
    String track = "";

    public List<String> generateParenthesis(int n) {
        if (n == 0) return null;
// 记录所有合法的括号组合

// 可用的左括号和右括号数量初始化为 n backtrack(n, n, track, res); return res;
        backtrack(n, n);
        return res;
    }

    public void backtrack(int left, int right) {
        //左括号使用一定大于等于右括号，这里是去掉非法的分支
        if (right < left) return;
        if (left < 0 || right < 0) return;// 当所有括号都恰好用完时，得到一个合法的括号组合
        if (left == 0 && right == 0) {
            res.add(track);
            return;
        }
        track = track + "(";
        backtrack(left - 1, right);
        track = track.substring(0, track.length() - 1);
        track = track + ")";
        backtrack(left, right - 1);
        track = track.substring(0, track.length() - 1);

    }

    public static void main(String[] args) {
        LeetCode_21_generateParenthesis leetCode_21_generateParenthesis = new LeetCode_21_generateParenthesis();
        System.out.println(Arrays.toString(leetCode_21_generateParenthesis.generateParenthesis(3).toArray()));

    }

}
