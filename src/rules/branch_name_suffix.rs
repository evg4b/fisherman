use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;

pub struct BranchNameSuffix {
    name: String,
    suffix: TemplateString,
}

impl BranchNameSuffix {
    pub fn new(name: String, suffix: TemplateString) -> Self {
        Self { name, suffix }
    }
}

impl CompiledRule for BranchNameSuffix {
    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        let processed_prefix = self.suffix.to_string()?;
        let branch_name = ctx.current_branch()?;
        if branch_name.ends_with(&processed_prefix) {
            Ok(RuleResult::Success {
                name: self.name.clone(),
                output: processed_prefix,
            })
        } else {
            Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "Branch name does not end with suffix: {}",
                    processed_prefix
                ),
            })
        }
    }
}
