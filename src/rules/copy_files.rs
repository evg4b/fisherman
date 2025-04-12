use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

pub struct CopyFiles {
    name: String,
    glob: TemplateString,
    destination: TemplateString,
}

impl CopyFiles {
    pub fn new(name: String, glob: TemplateString, destination: TemplateString) -> CopyFiles {
        CopyFiles {
            name,
            glob,
            destination,
        }
    }
}

impl CompiledRule for CopyFiles {
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
