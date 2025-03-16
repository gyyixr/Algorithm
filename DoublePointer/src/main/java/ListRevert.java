import java.util.Comparator;
import java.util.PriorityQueue;

public class ListRevert {


    public static void main(String[] args) {
        System.out.println("222222");
        PriorityQueue q = new PriorityQueue<String>(Comparator.reverseOrder());
        q.add("1");
        q.add("2");
        q.add("3");
        System.out.println(q.poll());

    }
}