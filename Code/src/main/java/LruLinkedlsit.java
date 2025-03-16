import java.util.HashMap;
import java.util.LinkedList;
import java.util.Map;

class LRULinkedList<K, V> {
    private final int capacity;
    private final Map<K, V> cacheMap;
    private final LinkedList<K> lruList;

    public LRULinkedList(int capacity) {
        this.capacity = capacity;
        this.cacheMap = new HashMap<>();
        this.lruList = new LinkedList<>();
    }

    public V get(K key) {
        if (cacheMap.containsKey(key)) {
            // 移动访问过的元素到链表头部，表示最近使用
            lruList.remove(key);
            lruList.addFirst(key);
            return cacheMap.get(key);
        }
        return null;
    }

    public void put(K key, V value) {
        if (cacheMap.size() >= capacity) {
            // 达到容量上限，淘汰最久未使用的元素（链表尾部）
            K leastRecentlyUsed = lruList.removeLast();
            cacheMap.remove(leastRecentlyUsed);
        }
        // 将新元素添加到链表头部表示最近使用
        lruList.addFirst(key);
        cacheMap.put(key, value);
    }

    public void printCacheState() {
        System.out.println("Cache State:");
        for (K key : lruList) {
            System.out.println(key + ": " + cacheMap.get(key));
        }
        System.out.println("---------------");
    }
    public static void main(String[] args) {
        LRULinkedList<Integer, String> lruCache = new LRULinkedList<>(3);

        lruCache.put(1, "One");
        lruCache.printCacheState();

        lruCache.put(2, "Two");
        lruCache.printCacheState();

        lruCache.put(3, "Three");
        lruCache.printCacheState();

        System.out.println("Value for key 2: " + lruCache.get(2));
        lruCache.printCacheState();

        lruCache.put(4, "Four"); // This will trigger eviction of the least recently used element (key 1)
        lruCache.printCacheState();
    }
}