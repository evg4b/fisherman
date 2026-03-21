
use crate::rules::compiled_rule::CompiledRule;
use crate::scripting::Expression;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct RuleOLD {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
}


#[cfg(test)]
mod tests {}
