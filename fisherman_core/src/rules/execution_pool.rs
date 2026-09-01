use anyhow::Result;
use std::future::Future;
use std::sync::Arc;
use tokio::sync::Semaphore;
use tokio::task::JoinHandle;

/// Maximum number of async rules allowed to run at the same time.
pub const MAX_CONCURRENT_ASYNC_RULES: usize = 8;

/// Runs rule tasks concurrently, capping how many are in flight at once.
///
/// Async rules do file I/O or spawn processes, so letting every one of them
/// loose at the same time would swamp the machine on a large hook. Tasks handed
/// to [`execute`](Self::execute) are spawned immediately but only start once the
/// pool has a free slot.
pub struct RuleExecutionPool {
    semaphore: Arc<Semaphore>,
}

impl Default for RuleExecutionPool {
    fn default() -> Self {
        Self::new()
    }
}

impl RuleExecutionPool {
    pub fn new() -> Self {
        Self::with_limit(MAX_CONCURRENT_ASYNC_RULES)
    }

    /// Builds a pool that runs at most `limit` tasks at the same time.
    pub fn with_limit(limit: usize) -> Self {
        Self {
            semaphore: Arc::new(Semaphore::new(limit)),
        }
    }

    /// Spawns `task`, holding it until the pool has room for it.
    pub fn execute<T, F>(&self, task: F) -> JoinHandle<Result<T>>
    where
        F: Future<Output = Result<T>> + Send + 'static,
        T: Send + 'static,
    {
        let semaphore = self.semaphore.clone();

        tokio::spawn(async move {
            let _permit = semaphore.acquire().await?;
            task.await
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use anyhow::anyhow;
    use std::sync::atomic::{AtomicUsize, Ordering};

    #[derive(Default)]
    struct Tracker {
        in_flight: AtomicUsize,
        peak_in_flight: AtomicUsize,
    }

    impl Tracker {
        /// Occupies a slot long enough for any other admitted task to overlap.
        async fn work(&self) {
            let in_flight = self.in_flight.fetch_add(1, Ordering::SeqCst) + 1;
            self.peak_in_flight.fetch_max(in_flight, Ordering::SeqCst);

            // Hand the scheduler enough turns for every other admitted task to
            // reach this point, so genuine overlap is observable without timing.
            for _ in 0..8 {
                tokio::task::yield_now().await;
            }

            self.in_flight.fetch_sub(1, Ordering::SeqCst);
        }

        fn peak(&self) -> usize {
            self.peak_in_flight.load(Ordering::SeqCst)
        }
    }

    async fn run(pool: &RuleExecutionPool, tasks: usize, tracker: Arc<Tracker>) -> Result<()> {
        let handles: Vec<_> = (0..tasks)
            .map(|index| {
                let tracker = tracker.clone();
                pool.execute(async move {
                    tracker.work().await;
                    anyhow::Ok(index)
                })
            })
            .collect();

        for handle in handles {
            handle.await??;
        }

        Ok(())
    }

    #[tokio::test]
    async fn runs_tasks_concurrently_up_to_the_limit() -> Result<()> {
        let tracker = Arc::new(Tracker::default());
        run(&RuleExecutionPool::with_limit(4), 4, tracker.clone()).await?;

        assert_eq!(tracker.peak(), 4);

        Ok(())
    }

    #[tokio::test]
    async fn caps_how_many_tasks_run_at_once() -> Result<()> {
        let tracker = Arc::new(Tracker::default());
        run(&RuleExecutionPool::with_limit(2), 6, tracker.clone()).await?;

        assert!(tracker.peak() <= 2);

        Ok(())
    }

    #[tokio::test]
    async fn every_task_still_runs_when_queued_behind_the_limit() -> Result<()> {
        let tracker = Arc::new(Tracker::default());
        let pool = RuleExecutionPool::with_limit(1);
        let counter = Arc::new(AtomicUsize::new(0));

        let handles: Vec<_> = (0..5)
            .map(|_| {
                let tracker = tracker.clone();
                let counter = counter.clone();
                pool.execute(async move {
                    tracker.work().await;
                    anyhow::Ok(counter.fetch_add(1, Ordering::SeqCst))
                })
            })
            .collect();

        for handle in handles {
            handle.await??;
        }

        assert_eq!(counter.load(Ordering::SeqCst), 5);
        assert_eq!(tracker.peak(), 1);

        Ok(())
    }

    #[tokio::test]
    async fn propagates_task_errors() -> Result<()> {
        let handle = RuleExecutionPool::new().execute(async { Err::<(), _>(anyhow!("boom")) });

        let result = handle.await?;
        assert_eq!(result.unwrap_err().to_string(), "boom");

        Ok(())
    }

    #[tokio::test]
    async fn default_pool_uses_the_configured_limit() {
        assert_eq!(
            RuleExecutionPool::default().semaphore.available_permits(),
            MAX_CONCURRENT_ASYNC_RULES
        );
    }
}
