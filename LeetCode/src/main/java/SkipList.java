import java.util.Random;

class Node {
    int value;
    Node[] forward;

    public Node(int value, int level) {
        this.value = value;
        this.forward = new Node[level];
    }
}

public class SkipList {
    private static final int MAX_LEVEL = 16;
    private Node header;
    private int level;

    public SkipList() {
        this.header = new Node(0, MAX_LEVEL);
        this.level = 1;
    }

    private int randomLevel() {
        int level = 1;
        Random rand = new Random();
        while (rand.nextDouble() < 0.5 && level < MAX_LEVEL) {
            level++;
        }
        return level;
    }

    public void insert(int value) {
        Node[] update = new Node[MAX_LEVEL];
        Node current = header;

        for (int i = level - 1; i >= 0; i--) {
            while (current.forward[i] != null && current.forward[i].value < value) {
                current = current.forward[i];
            }
            update[i] = current;
        }

        int newLevel = randomLevel();
        if (newLevel > level) {
            for (int i = level; i < newLevel; i++) {
                update[i] = header;
            }
            level = newLevel;
        }

        Node newNode = new Node(value, newLevel);
        for (int i = 0; i < newLevel; i++) {
            newNode.forward[i] = update[i].forward[i];
            update[i].forward[i] = newNode;
        }
    }

    public void delete(int value) {
        Node[] update = new Node[MAX_LEVEL];
        Node current = header;

        for (int i = level - 1; i >= 0; i--) {
            while (current.forward[i] != null && current.forward[i].value < value) {
                current = current.forward[i];
            }
            update[i] = current;
        }

        if (current.forward[0] != null && current.forward[0].value == value) {
            for (int i = 0; i < level; i++) {
                if (update[i].forward[i] != null && update[i].forward[i].value == value) {
                    update[i].forward[i] = update[i].forward[i].forward[i];
                }
            }

            // Update the level if necessary
            while (level > 1 && header.forward[level - 1] == null) {
                level--;
            }
        }
    }

    public boolean search(int value) {
        Node current = header;
        for (int i = level - 1; i >= 0; i--) {
            while (current.forward[i] != null && current.forward[i].value < value) {
                current = current.forward[i];
            }
            if (current.forward[i] != null && current.forward[i].value == value) {
                return true;
            }
        }
        return false;
    }

    public void display() {
        System.out.println("Skip List:");
        for (int i = level - 1; i >= 0; i--) {
            System.out.print("Level " + i + ": ");
            Node node = header.forward[i];
            while (node != null) {
                System.out.print(node.value + " ");
                node = node.forward[i];
            }
            System.out.println();
        }
        System.out.println();
    }

    public static void main(String[] args) {
        SkipList skipList = new SkipList();

        // Insert elements
        skipList.insert(3);
        skipList.insert(6);
        skipList.insert(7);
        skipList.insert(9);
        skipList.insert(12);
        skipList.insert(19);
        skipList.insert(17);
        skipList.insert(26);
        skipList.insert(21);
        skipList.insert(25);
        skipList.display();

        // Search for an element
        int searchValue = 19;
        System.out.println("Search " + searchValue + ": " + skipList.search(searchValue));

        // Delete an element
        int deleteValue = 17;
        System.out.println("Delete " + deleteValue);
        skipList.delete(deleteValue);
        skipList.display();

        // Search for the deleted element
        System.out.println("Search " + deleteValue + ": " + skipList.search(deleteValue));
    }
}
