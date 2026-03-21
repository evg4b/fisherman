use crate::context::Context;
use anyhow::Result;

#[derive(Debug)]
pub enum RuleResultOld {
    Success {
        name: String,
        output: Option<String>,
    },
    Failure {
        name: String,
        message: String,
    },
}

pub trait CompiledRule: Send + Sync {
    /// Returns `true` if this rule must run sequentially (one at a time),
    /// or `false` if it is self-contained and may run in the parallel pool.
    fn is_sequential(&self) -> bool;
    fn check(&self, ctx: &dyn Context) -> Result<RuleResultOld>;
}
