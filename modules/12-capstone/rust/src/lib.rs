//! The capstone cache in Rust: a bounded concurrent cache over a slow backend,
//! with a counting semaphore for backpressure. No async runtime -- std threads,
//! Mutex, and a hand-rolled Semaphore (std has no counting semaphore).

use std::collections::HashMap;
use std::sync::atomic::{AtomicI64, Ordering};
use std::sync::{Condvar, Mutex};
use std::time::Duration;

pub struct Semaphore {
    count: Mutex<usize>,
    cv: Condvar,
}

impl Semaphore {
    pub fn new(n: usize) -> Self {
        Semaphore {
            count: Mutex::new(n),
            cv: Condvar::new(),
        }
    }
    fn acquire(&self) {
        let mut c = self.count.lock().unwrap();
        while *c == 0 {
            c = self.cv.wait(c).unwrap();
        }
        *c -= 1;
    }
    fn release(&self) {
        *self.count.lock().unwrap() += 1;
        self.cv.notify_one();
    }
}

pub struct Cache {
    entries: Mutex<HashMap<i64, i64>>,
    capacity: usize,
    backend_calls: AtomicI64,
    latency: Duration,
    sem: Semaphore,
}

impl Cache {
    pub fn new(capacity: usize, max_inflight: usize, latency: Duration) -> Self {
        Cache {
            entries: Mutex::new(HashMap::with_capacity(capacity)),
            capacity,
            backend_calls: AtomicI64::new(0),
            latency,
            sem: Semaphore::new(max_inflight),
        }
    }

    pub fn get(&self, key: i64) -> i64 {
        {
            let m = self.entries.lock().unwrap();
            if let Some(&v) = m.get(&key) {
                return v;
            }
        }
        self.sem.acquire();
        self.backend_calls.fetch_add(1, Ordering::Relaxed);
        std::thread::sleep(self.latency);
        let v = key * key;
        self.sem.release();

        let mut m = self.entries.lock().unwrap();
        if m.len() >= self.capacity
            && let Some(&k) = m.keys().next()
        {
            m.remove(&k);
        }
        m.insert(key, v);
        v
    }

    pub fn backend_calls(&self) -> i64 {
        self.backend_calls.load(Ordering::Relaxed)
    }
    pub fn len(&self) -> usize {
        self.entries.lock().unwrap().len()
    }
    pub fn is_empty(&self) -> bool {
        self.len() == 0
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::Arc;

    #[test]
    fn correct_and_bounded() {
        let c = Arc::new(Cache::new(16, 4, Duration::ZERO));
        let handles: Vec<_> = (0..8)
            .map(|w| {
                let c = Arc::clone(&c);
                std::thread::spawn(move || {
                    for i in 0..500i64 {
                        let key = (w * 500 + i) % 100;
                        assert_eq!(c.get(key), key * key);
                    }
                })
            })
            .collect();
        for h in handles {
            h.join().unwrap();
        }
        assert!(c.len() <= 16);
        assert!(c.backend_calls() < 8 * 500);
    }
}
