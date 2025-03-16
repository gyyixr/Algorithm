import java.util.HashMap;
import java.util.LinkedHashSet;

class LFUCache {
    private int capacity;
    private HashMap<Integer, Integer> keyToVal;
    private HashMap<Integer, Integer> keyToFreq;
    private HashMap<Integer, LinkedHashSet<Integer>> freqToKeys;
    private int minFreq;

    public LFUCache(int capacity) {
        this.capacity = capacity;
        this.keyToVal = new HashMap<>();
        this.keyToFreq = new HashMap<>();
        this.freqToKeys = new HashMap<>();
        this.minFreq = 0;
    }

    public int get(int key) {
        if (!keyToVal.containsKey(key)) {
            return -1;
        }

        int freq = keyToFreq.get(key);
        freqToKeys.get(freq).remove(key);

        if (freq == minFreq && freqToKeys.get(freq).isEmpty()) {
            minFreq++;
        }

        keyToFreq.put(key, freq + 1);
        freqToKeys.computeIfAbsent(freq + 1, k -> new LinkedHashSet<>()).add(key);

        return keyToVal.get(key);
    }

    public void put(int key, int value) {
        if (capacity <= 0) {
            return;
        }

        if (keyToVal.containsKey(key)) {
            keyToVal.put(key, value);
            get(key);  // Increase frequency
            return;
        }

        if (keyToVal.size() >= capacity) {
            removeLFU();
        }

        keyToVal.put(key, value);
        keyToFreq.put(key, 1);
        freqToKeys.computeIfAbsent(1, k -> new LinkedHashSet<>()).add(key);
        minFreq = 1;
    }

    private void removeLFU() {
        LinkedHashSet<Integer> minFreqSet = freqToKeys.get(minFreq);
        int keyToRemove = minFreqSet.iterator().next();
        minFreqSet.remove(keyToRemove);

        if (minFreqSet.isEmpty()) {
            freqToKeys.remove(minFreq);
        }

        keyToVal.remove(keyToRemove);
        keyToFreq.remove(keyToRemove);
    }

    public static void main(String[] args) {
        LFUCache lfuCache = new LFUCache(3);

        lfuCache.put(1, 10);
        lfuCache.printCacheState();

        lfuCache.put(2, 20);
        lfuCache.printCacheState();

        lfuCache.put(3, 30);
        lfuCache.printCacheState();

        System.out.println("Value for key 2: " + lfuCache.get(2));
        lfuCache.printCacheState();

        lfuCache.put(4, 40); // This will trigger eviction of the least frequently used element (key 1)
        lfuCache.printCacheState();
    }

    private void printCacheState() {
        System.out.println("LFU Cache State:");
        for (int freq : freqToKeys.keySet()) {
            System.out.print("Freq " + freq + ": ");
            for (int key : freqToKeys.get(freq)) {
                System.out.print("(" + key + ", " + keyToVal.get(key) + ") ");
            }
            System.out.println();
        }
        System.out.println("---------------");
    }
}
