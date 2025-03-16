package 贪心;

/**
 * 例如，如果 s = "()))" ，你可以插入一个开始括号为 "(()))" 或结束括号为 "())))" 。
 * 返回 为使结果字符串 s 有效而必须添加的最少括号数。
 */
public class LeetCode_921_minAddToMakeValid {
    public int minAddToMakeValid(String S) {
        int res=0,record=0;
        for(int i=0;i<S.length();i++){
            if(S.charAt(i)=='(')record++;
            else{
                if(record!=0)record--;
                else res++;
            }
        }
        return res+record;
    }

}
