use anyhow::Result;

#[derive(Debug)]
pub enum RuleResult {
    Success { name: String, output: String, },
    Failure { name: String, message: String },
}

pub trait CompiledRule {
    fn check(&self) -> Result<RuleResult>;
}