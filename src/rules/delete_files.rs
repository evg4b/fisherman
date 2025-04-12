use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

pub struct DeleteFiles {
    name: String,
    glob: TemplateString,
}

impl DeleteFiles {
    pub fn new(name: String, glob: TemplateString) -> DeleteFiles {
        DeleteFiles { name, glob }
    }
}

impl CompiledRule for DeleteFiles {
    fn sync(&self) -> bool {
        false
    }

    fn check(&self, _: &dyn Context) -> Result<RuleResult> {
        Ok(RuleResult::Success {
            name: self.name.clone(),
            output: None,
        })
    }
}
