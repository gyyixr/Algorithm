package fenzhi;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
// 思路：（1）根据操作符分割，做分治递归（2）要考虑纯数字的情况 （3）使用字典保存计算过的表达式
/**
 * 给你一个由数字和运算符组成的字符串 expression ，按不同优先级组合数字和运算符，计算并返回所有可能组合的结果。你可以 按任意顺序 返回答案。
 * 示例 1：
 * 输入：expression = "2-1-1"
 * 输出：[0,2]
 * 解释：
 * ((2-1)-1) = 0
 * (2-(1-1)) = 2
 * 示例 2：
 * 输入：expression = "2*3-4*5"
 * 输出：[-34,-14,-10,-10,10]
 * 解释：
 * (2*(3-(4*5))) = -34
 * ((2*3)-(4*5)) = -14
 * ((2*(3-4))*5) = -10
 * (2*((3-4)*5)) = -10
 * (((2*3)-4)*5) = 10
 */
public class LeectCode_241_express_compute {
    //添加一个 map
    HashMap<String, List<Integer>> map = new HashMap<>();

    public List<Integer> diffWaysToCompute(String input) {
        if (input.length() == 0) {
            return new ArrayList<>();
        }
        //如果已经有当前解了，直接返回
        if (map.containsKey(input)) {
            return map.get(input);
        }
        List<Integer> result = new ArrayList<>();
        int num = 0;
        int index = 0;
        while (index < input.length() && !isOperation(input.charAt(index))) {
            num = num * 10 + input.charAt(index) - '0';
            index++;
        }
        if (index == input.length()) {
            result.add(num);
            //存到 map
            map.put(input, result);
            return result;
        }
        for (int i = 0; i < input.length(); i++) {
            if (isOperation(input.charAt(i))) {
                List<Integer> resultLeft = diffWaysToCompute(input.substring(0, i));
                List<Integer> resultRight = diffWaysToCompute(input.substring(i + 1));
                // 题目要求计算所有可能的情况，所以这里要两层for遍历所有的情况
                for (int j = 0; j < resultLeft.size(); j++) {
                    for (int k = 0; k < resultRight.size(); k++) {
                        char op = input.charAt(i);
                        result.add(calculate(resultLeft.get(j), op, resultRight.get(k)));
                    }
                }
            }
        }
        //存到 map
        map.put(input, result);
        return result;
    }

    private int calculate(int num1, char c, int num2) {
        switch (c) {
            case '+':
                return num1 + num2;
            case '-':
                return num1 - num2;
            case '*':
                return num1 * num2;
        }
        return -1;
    }

    private boolean isOperation(char c) {
        return c == '+' || c == '-' || c == '*';
    }

    public static void main(String[] args) {
        LeectCode_241_express_compute leectCode_241_express_compute = new LeectCode_241_express_compute();
        System.out.println(leectCode_241_express_compute.diffWaysToCompute("2*3-4*5"));

    }
}
