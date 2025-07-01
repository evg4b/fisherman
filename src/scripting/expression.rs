use anyhow::Result;
use rhai::{Engine, Scope};
use serde::{Deserialize, Deserializer};
use std::cell::RefCell;
use std::collections::HashMap;

thread_local! {
    static ENGINE: RefCell<Engine> = RefCell::new(Engine::new());
}

#[derive(Debug, Clone)]
pub struct Expression {
    pub condition: String,
}

impl<'de> Deserialize<'de> for Expression {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        Ok(Expression::new(s.as_str()))
    }
}

impl Expression {
    pub fn new(condition: &str) -> Self {
        Self {
            condition: condition.to_string(),
        }
    }

    pub fn check(&self, variables: &HashMap<String, String>) -> Result<bool> {
        let engine = ENGINE.take();

        let mut scope = Scope::with_capacity(variables.len());

        variables.iter().for_each(|(key, value)| {
            scope.push(key.to_owned(), value.to_owned());
        });

        match engine.eval_expression_with_scope::<bool>(&mut scope, &self.condition) {
            Ok(result) => Ok(result),
            Err(err) => Err(anyhow::anyhow!("Error: {}", err)),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_check_expression() {
        let a = Expression::new("1 > 0").check(&HashMap::new()).unwrap();
        assert!(a);
    }

    #[test]
    fn test_check_expression2() {
        let a = Expression::new("is_def_var(\"xx\") && xx > 10")
            .check(&HashMap::new())
            .unwrap();
        assert!(!a);
    }

    #[test]
    fn test_check_expression_error() {
        let a = Expression::new("1 >").check(&HashMap::new()).unwrap_err();
        assert_eq!(
            a.to_string(),
            "Error: Syntax error: Script is incomplete (line 1, position 4)"
        );
    }

    #[test]
    fn test_check_expression3() {
        let mut variables = HashMap::new();
        variables.insert("xx".to_string(), "20".to_string());
        let a = Expression::new("parse_int(xx) > 10")
            .check(&variables)
            .unwrap();
        assert!(a);
    }

    #[test]
    fn test_check_expression23() {
        let mut variables = HashMap::new();
        variables.insert("xx".to_string(), "91".to_string());
        let a = Expression::new("(is_def_var(\"yy\") && parse_int(yy) > 10) || (is_def_var(\"xx\") && parse_int(xx) > 10)")
            .check(&variables)
            .unwrap();
        assert!(a);
    }

    #[test]
    fn expression_deserialize() {
        let rule = r#"
            "is_def_var(\"test\")"
        "#;

        let rule: Expression = serde_json::from_str(rule).unwrap();

        assert_eq!(rule.condition, "is_def_var(\"test\")".to_string());
    }

    #[test]
    fn expression_deserialize2() {
        let rule = r#"
            { "condition": "1 > 0" }
        "#;

        let err = serde_json::from_str::<Expression>(rule).unwrap_err();

        assert_eq!(
            err.to_string(),
            "invalid type: map, expected a string at line 2 column 12"
        );
    }
}
